package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const defaultMailgunAPIBase = "https://api.mailgun.net"

// MailgunSender sends publish-link emails through Mailgun.
type MailgunSender struct {
	apiBase     string
	domain      string
	apiKey      string
	defaultFrom string
	client      *http.Client
}

// NewMailgunSender constructs a Mailgun adapter.
func NewMailgunSender(apiBase, domain, apiKey, defaultFrom string, timeout time.Duration) (*MailgunSender, error) {
	base := strings.TrimRight(strings.TrimSpace(apiBase), "/")
	if base == "" {
		base = defaultMailgunAPIBase
	}
	if strings.TrimSpace(domain) == "" {
		return nil, fmt.Errorf("mailgun domain is required")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("mailgun api key is required")
	}
	if strings.TrimSpace(defaultFrom) == "" {
		defaultFrom = "response@mg.supost.com"
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &MailgunSender{
		apiBase:     base,
		domain:      strings.TrimSpace(domain),
		apiKey:      strings.TrimSpace(apiKey),
		defaultFrom: strings.TrimSpace(defaultFrom),
		client:      &http.Client{Timeout: timeout},
	}, nil
}

// SendPublishEmail sends one plain-text publish message.
func (m *MailgunSender) SendPublishEmail(ctx context.Context, msg domain.PublishEmailMessage) error {
	to := strings.TrimSpace(msg.To)
	subject := strings.TrimSpace(msg.Subject)
	text := strings.TrimSpace(msg.Text)
	if to == "" || subject == "" || text == "" {
		return fmt.Errorf("mailgun message to/subject/text are required")
	}

	from := strings.TrimSpace(msg.From)
	if from == "" {
		from = m.defaultFrom
	}

	form := url.Values{}
	form.Set("from", from)
	form.Set("to", to)
	form.Set("subject", subject)
	form.Set("text", text)

	endpoint := fmt.Sprintf("%s/v3/%s/messages", m.apiBase, m.domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("creating mailgun request: %w", err)
	}
	req.SetBasicAuth("api", m.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending mailgun request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("mailgun send failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
