package service

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	responseSafetyLine  = "Safety: If someone sends you a check, do not send them any money back. https://supost.com/safety"
	responseContactLine = "Report responses to contact@supost.com"
)

// PostRespondRepository defines post lookup + message persistence operations.
type PostRespondRepository interface {
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
	CreateResponseMessage(ctx context.Context, postID int64, replyToEmail, message, ip, userAgent string) (domain.Message, error)
}

// PostRespondEmailSender defines response-email side effects.
type PostRespondEmailSender interface {
	SendResponseEmail(ctx context.Context, msg domain.ResponseEmailMessage) error
}

// PostRespondService orchestrates post response sends.
type PostRespondService struct {
	repo PostRespondRepository
}

// NewPostRespondService constructs PostRespondService.
func NewPostRespondService(repo PostRespondRepository) *PostRespondService {
	return &PostRespondService{repo: repo}
}

// Respond validates, optionally sends, and optionally persists a response message.
func (s *PostRespondService) Respond(
	ctx context.Context,
	input domain.PostRespondSubmission,
	dryRun bool,
	baseURL string,
	fromEmail string,
	sender PostRespondEmailSender,
) (domain.PostRespondResult, error) {
	normalized, err := normalizePostRespondInput(input)
	if err != nil {
		return domain.PostRespondResult{}, err
	}

	post, err := s.repo.GetPostByID(ctx, normalized.PostID)
	if err != nil {
		return domain.PostRespondResult{}, err
	}
	if strings.TrimSpace(post.Email) == "" {
		return domain.PostRespondResult{}, fmt.Errorf("post %d has no destination email", normalized.PostID)
	}
	if strings.TrimSpace(post.AccessToken) == "" {
		return domain.PostRespondResult{}, fmt.Errorf("post %d has no access token", normalized.PostID)
	}

	subject, body := buildResponseEmailContent(post, normalized, baseURL)
	result := domain.PostRespondResult{
		DryRun:       dryRun,
		PostID:       post.ID,
		PostEmail:    strings.TrimSpace(post.Email),
		ReplyTo:      normalized.ReplyTo,
		MessageSaved: false,
		EmailSent:    false,
		Subject:      subject,
		Body:         body,
		SentAt:       time.Now(),
	}

	if dryRun {
		return result, nil
	}
	if sender == nil {
		return domain.PostRespondResult{}, fmt.Errorf("response email sender is required")
	}

	msg := domain.ResponseEmailMessage{
		From:    strings.TrimSpace(fromEmail),
		To:      result.PostEmail,
		ReplyTo: result.ReplyTo,
		Subject: subject,
		Text:    body,
	}
	if err := sender.SendResponseEmail(ctx, msg); err != nil {
		return domain.PostRespondResult{}, err
	}
	result.EmailSent = true

	saved, err := s.repo.CreateResponseMessage(ctx, post.ID, result.ReplyTo, normalized.Message, normalized.IP, normalized.UserAgent)
	if err != nil {
		return domain.PostRespondResult{}, err
	}
	if saved.ID > 0 {
		result.MessageID = saved.ID
	}
	result.MessageSaved = true
	return result, nil
}

func normalizePostRespondInput(input domain.PostRespondSubmission) (domain.PostRespondSubmission, error) {
	normalized := input
	normalized.Message = strings.TrimSpace(input.Message)
	normalized.ReplyTo = strings.ToLower(strings.TrimSpace(input.ReplyTo))
	normalized.IP = strings.TrimSpace(input.IP)
	normalized.UserAgent = strings.TrimSpace(input.UserAgent)

	problems := make([]string, 0, 5)
	if normalized.PostID <= 0 {
		problems = append(problems, "post_id is required")
	}
	if normalized.Message == "" {
		problems = append(problems, "message is required")
	}
	if normalized.ReplyTo == "" {
		problems = append(problems, "reply_to is required")
	} else if !isValidEmail(normalized.ReplyTo) {
		problems = append(problems, "reply_to must be a valid email")
	}
	if normalized.IP != "" {
		if _, err := netip.ParseAddr(normalized.IP); err != nil {
			problems = append(problems, "ip must be a valid IPv4 or IPv6 address")
		}
	}

	if len(problems) > 0 {
		return domain.PostRespondSubmission{}, fmt.Errorf("validation failed: %s", strings.Join(problems, "; "))
	}
	return normalized, nil
}

func isValidEmail(value string) bool {
	at := strings.LastIndex(value, "@")
	if at <= 0 || at >= len(value)-1 {
		return false
	}
	domainPart := strings.TrimSpace(value[at+1:])
	return strings.Contains(domainPart, ".")
}

func buildResponseEmailContent(post domain.Post, input domain.PostRespondSubmission, baseURL string) (string, string) {
	root := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if root == "" {
		root = defaultSupostBaseURL
	}
	title := strings.TrimSpace(post.Name)
	if title == "" {
		title = "(untitled post)"
	}

	postedAt := postTimestamp(post)
	if postedAt.IsZero() {
		postedAt = time.Now()
	}
	postedLine := postedAt.Format("Mon, Jan 2, 2006 03:04 PM")
	publishURL := root + "/post/publish/" + strings.TrimSpace(post.AccessToken)
	postURL := fmt.Sprintf("%s/post/index/%d", root, post.ID)

	subject := fmt.Sprintf("SUpost - %s response: %s", input.ReplyTo, title)
	body := strings.Join([]string{
		"Reply to: " + input.ReplyTo,
		"",
		input.Message,
		"",
		responseSafetyLine,
		"",
		fmt.Sprintf("%s - Posted: %s", title, postedLine),
		"",
		"To delete your post, use this link and click 'Delete your post.'",
		publishURL,
		"",
		postURL,
		"",
		responseContactLine,
	}, "\n")
	return subject, body
}

func postTimestamp(post domain.Post) time.Time {
	if !post.TimePostedAt.IsZero() {
		return post.TimePostedAt
	}
	if post.TimePosted > 0 {
		return time.Unix(post.TimePosted, 0)
	}
	return time.Time{}
}
