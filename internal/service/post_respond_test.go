package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockPostRespondRepo struct {
	post         domain.Post
	savedMessage domain.Message
	saveCalled   bool
}

func (m *mockPostRespondRepo) GetPostByID(_ context.Context, _ int64) (domain.Post, error) {
	return m.post, nil
}

func (m *mockPostRespondRepo) CreateResponseMessage(_ context.Context, _ int64, replyToEmail, message, _ string) (domain.Message, error) {
	m.saveCalled = true
	m.savedMessage.Email = replyToEmail
	m.savedMessage.RawEmail = replyToEmail
	m.savedMessage.Message = message
	if m.savedMessage.ID == 0 {
		m.savedMessage.ID = 77
	}
	return m.savedMessage, nil
}

type mockPostRespondSender struct {
	last domain.ResponseEmailMessage
	sent bool
}

func (m *mockPostRespondSender) SendResponseEmail(_ context.Context, msg domain.ResponseEmailMessage) error {
	m.last = msg
	m.sent = true
	return nil
}

func TestPostRespondService_DryRun(t *testing.T) {
	repo := &mockPostRespondRepo{
		post: domain.Post{
			ID:           130031908,
			Email:        "wientjes@alumni.stanford.edu",
			Name:         "Looking for a buddy to go to the movies",
			AccessToken:  "dfc6dbef",
			TimePostedAt: time.Date(2026, time.February, 26, 17, 32, 0, 0, time.UTC),
		},
	}
	sender := &mockPostRespondSender{}
	svc := NewPostRespondService(repo)

	result, err := svc.Respond(context.Background(), domain.PostRespondSubmission{
		PostID:  130031908,
		Message: "looks good",
		ReplyTo: "gwientjes@gmail.com",
	}, true, "https://supost.com", "response@mg.supost.com", sender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EmailSent || result.MessageSaved {
		t.Fatalf("dry run should not send or save")
	}
	if sender.sent || repo.saveCalled {
		t.Fatalf("dry run should not invoke sender/repo save")
	}
	if !strings.Contains(result.Body, "https://supost.com/post/publish/dfc6dbef") {
		t.Fatalf("expected publish URL in body")
	}
}

func TestPostRespondService_SendAndPersist(t *testing.T) {
	repo := &mockPostRespondRepo{
		post: domain.Post{
			ID:           130031908,
			Email:        "owner@stanford.edu",
			Name:         "Looking for a buddy to go to the movies",
			AccessToken:  "dfc6dbef",
			TimePostedAt: time.Date(2026, time.February, 26, 17, 32, 0, 0, time.UTC),
		},
	}
	sender := &mockPostRespondSender{}
	svc := NewPostRespondService(repo)

	result, err := svc.Respond(context.Background(), domain.PostRespondSubmission{
		PostID:  130031908,
		Message: "Hello, I want to buy your bike",
		ReplyTo: "gwientjes@gmail.com",
	}, false, "https://supost.com", "response@mg.supost.com", sender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.EmailSent || !result.MessageSaved {
		t.Fatalf("expected send+save true")
	}
	if !sender.sent || !repo.saveCalled {
		t.Fatalf("expected sender and repo save calls")
	}
	if sender.last.To != "owner@stanford.edu" {
		t.Fatalf("unexpected email destination %q", sender.last.To)
	}
	if sender.last.ReplyTo != "gwientjes@gmail.com" {
		t.Fatalf("unexpected reply-to %q", sender.last.ReplyTo)
	}
}

func TestPostRespondService_Validation(t *testing.T) {
	svc := NewPostRespondService(&mockPostRespondRepo{})
	_, err := svc.Respond(context.Background(), domain.PostRespondSubmission{
		PostID:  1,
		Message: "",
		ReplyTo: "invalid",
	}, true, "https://supost.com", "response@mg.supost.com", &mockPostRespondSender{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "message is required") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestPostRespondService_ValidationReplyToRequired(t *testing.T) {
	svc := NewPostRespondService(&mockPostRespondRepo{})
	_, err := svc.Respond(context.Background(), domain.PostRespondSubmission{
		PostID:  130031802,
		Message: "Hello, I want to buy your bike",
		ReplyTo: "",
	}, true, "https://supost.com", "response@mg.supost.com", &mockPostRespondSender{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "reply_to is required") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestPostRespondService_ValidationReplyToEmailFormat(t *testing.T) {
	svc := NewPostRespondService(&mockPostRespondRepo{})
	_, err := svc.Respond(context.Background(), domain.PostRespondSubmission{
		PostID:  130031802,
		Message: "Hello, I want to buy your bike",
		ReplyTo: "gwientjesgmail.com",
	}, true, "https://supost.com", "response@mg.supost.com", &mockPostRespondSender{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "reply_to must be a valid email") {
		t.Fatalf("unexpected error %v", err)
	}
}
