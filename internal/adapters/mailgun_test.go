package adapters

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestMailgunSender_SendPublishEmail(t *testing.T) {
	var capturedBody string
	sender, err := NewMailgunSender("https://api.mailgun.net", "mg.supost.com", "test-key", "response@mg.supost.com", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}
	sender.client = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			if !strings.Contains(r.URL.Path, "/v3/mg.supost.com/messages") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			user, pass, ok := r.BasicAuth()
			if !ok || user != "api" || pass != "test-key" {
				t.Fatalf("unexpected basic auth")
			}
			body, _ := io.ReadAll(r.Body)
			capturedBody = string(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	err = sender.SendPublishEmail(context.Background(), domain.PublishEmailMessage{
		To:      "wientjes@alumni.stanford.edu",
		Subject: "SUpost - Publish your post! Test",
		Text:    "Publish your post by pressing:\n\nhttps://supost.com/post/publish/token",
	})
	if err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	values, err := url.ParseQuery(capturedBody)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if values.Get("from") != "response@mg.supost.com" {
		t.Fatalf("unexpected from %q", values.Get("from"))
	}
	if values.Get("to") != "wientjes@alumni.stanford.edu" {
		t.Fatalf("unexpected to %q", values.Get("to"))
	}
}

func TestMailgunSender_SendPublishEmail_Non2xx(t *testing.T) {
	sender, err := NewMailgunSender("https://api.mailgun.net", "mg.supost.com", "test-key", "response@mg.supost.com", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}
	sender.client = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("bad request")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	err = sender.SendPublishEmail(context.Background(), domain.PublishEmailMessage{
		To:      "wientjes@alumni.stanford.edu",
		Subject: "subject",
		Text:    "body",
	})
	if err == nil {
		t.Fatalf("expected error for non-2xx response")
	}
}
