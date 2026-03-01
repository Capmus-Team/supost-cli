package service

import (
	"context"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockUserSignupProvider struct {
	lastSubmission domain.UserSignupSubmission
	result         domain.UserSignupResult
	called         bool
}

func (m *mockUserSignupProvider) SignUp(_ context.Context, submission domain.UserSignupSubmission) (domain.UserSignupResult, error) {
	m.called = true
	m.lastSubmission = submission
	return m.result, nil
}

func TestUserSignupService_SignUp_Success(t *testing.T) {
	provider := &mockUserSignupProvider{
		result: domain.UserSignupResult{
			UserID:                "user-123",
			DisplayName:           "Greg Wientjes",
			Email:                 "wientjes@alumni.stanford.edu",
			Phone:                 "+16505551234",
			EmailConfirmationSent: true,
			CreatedAt:             time.Now(),
		},
	}
	svc := NewUserSignupService(provider)

	result, err := svc.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "  Greg Wientjes  ",
		Email:       "WIENTJES@ALUMNI.STANFORD.EDU",
		Phone:       " +16505551234 ",
		Password:    "supost123!",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !provider.called {
		t.Fatalf("expected provider call")
	}
	if provider.lastSubmission.Email != "wientjes@alumni.stanford.edu" {
		t.Fatalf("expected normalized lowercase email, got %q", provider.lastSubmission.Email)
	}
	if result.UserID == "" {
		t.Fatalf("expected user id")
	}
}

func TestUserSignupService_SignUp_ValidationErrors(t *testing.T) {
	svc := NewUserSignupService(&mockUserSignupProvider{})
	_, err := svc.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "",
		Email:       "bad-email",
		Phone:       "not-a-phone",
		Password:    "short",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
