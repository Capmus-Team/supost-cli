package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockPostCreateSubmitRepo struct {
	categories    []domain.Category
	subcategories []domain.Subcategory
	submission    domain.PostCreateSubmission
	persisted     domain.PostCreatePersisted
	createCalled  bool
}

func (m *mockPostCreateSubmitRepo) ListCategories(_ context.Context) ([]domain.Category, error) {
	return m.categories, nil
}

func (m *mockPostCreateSubmitRepo) ListSubcategories(_ context.Context) ([]domain.Subcategory, error) {
	return m.subcategories, nil
}

func (m *mockPostCreateSubmitRepo) CreatePendingPost(_ context.Context, submission domain.PostCreateSubmission) (domain.PostCreatePersisted, error) {
	m.createCalled = true
	m.submission = submission
	if m.persisted.PostID == 0 {
		m.persisted.PostID = 130031999
	}
	if m.persisted.AccessToken == "" {
		m.persisted.AccessToken = submission.AccessToken
	}
	if m.persisted.PostedAt.IsZero() {
		m.persisted.PostedAt = submission.PostedAt
	}
	return m.persisted, nil
}

type mockPublishSender struct {
	last domain.PublishEmailMessage
	sent bool
}

func (m *mockPublishSender) SendPublishEmail(_ context.Context, msg domain.PublishEmailMessage) error {
	m.last = msg
	m.sent = true
	return nil
}

func TestPostCreateService_Submit_DryRunDoesNotPersistOrSend(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)
	sender := &mockPublishSender{}

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", sender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry run result")
	}
	if repo.createCalled {
		t.Fatalf("expected no repository insert in dry run")
	}
	if sender.sent {
		t.Fatalf("expected no email send in dry run")
	}
	if !strings.Contains(result.PublishURL, "/post/publish/") {
		t.Fatalf("missing publish URL in result: %q", result.PublishURL)
	}
}

func TestPostCreateService_Submit_PersistsAndSends(t *testing.T) {
	now := time.Date(2026, time.February, 26, 17, 32, 0, 0, time.UTC)
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
		persisted: domain.PostCreatePersisted{
			PostID:      130031999,
			AccessToken: "abcdef",
			PostedAt:    now,
		},
	}
	sender := &mockPublishSender{}
	svc := NewPostCreateService(repo)

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
	}, false, "https://supost.com", "response@mg.supost.com", sender)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.createCalled {
		t.Fatalf("expected repository insert to be called")
	}
	if !sender.sent {
		t.Fatalf("expected publish email send")
	}
	if result.PostID != 130031999 {
		t.Fatalf("unexpected post id %d", result.PostID)
	}
	if sender.last.To != "wientjes@alumni.stanford.edu" {
		t.Fatalf("unexpected recipient %q", sender.last.To)
	}
	if !strings.Contains(sender.last.Text, "https://supost.com/post/publish/abcdef") {
		t.Fatalf("missing publish URL in email body")
	}
}

func TestPostCreateService_Submit_InvalidEmailRejected(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "user@gmail.com",
		Price:         100,
		PriceProvided: true,
	}, false, "https://supost.com", "response@mg.supost.com", &mockPublishSender{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "Stanford email") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_StanfordSubdomainAccepted(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@cs.stanford.edu",
		Price:         100,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{})
	if err != nil {
		t.Fatalf("expected cs.stanford.edu to pass validation, got %v", err)
	}
}
