package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// UserSignupProvider defines Supabase Auth signup side effects where consumed.
type UserSignupProvider interface {
	SignUp(ctx context.Context, submission domain.UserSignupSubmission) (domain.UserSignupResult, error)
}

// UserSignupService validates signup inputs and delegates to Supabase Auth.
type UserSignupService struct {
	provider UserSignupProvider
}

// NewUserSignupService constructs UserSignupService.
func NewUserSignupService(provider UserSignupProvider) *UserSignupService {
	return &UserSignupService{provider: provider}
}

// SignUp validates a signup payload and creates a new Supabase Auth user.
func (s *UserSignupService) SignUp(ctx context.Context, submission domain.UserSignupSubmission) (domain.UserSignupResult, error) {
	if s.provider == nil {
		return domain.UserSignupResult{}, fmt.Errorf("signup provider is required")
	}

	normalized := domain.UserSignupSubmission{
		DisplayName: strings.TrimSpace(submission.DisplayName),
		Email:       strings.ToLower(strings.TrimSpace(submission.Email)),
		Phone:       strings.TrimSpace(submission.Phone),
		Password:    submission.Password,
	}

	problems := make([]string, 0, 4)
	if normalized.DisplayName == "" {
		problems = append(problems, "display_name is required")
	}
	if normalized.Email == "" || !strings.Contains(normalized.Email, "@") {
		problems = append(problems, "email must be valid")
	}
	if !looksLikePhone(normalized.Phone) {
		problems = append(problems, "phone must be a valid international number (example: +16505551234)")
	}
	if len(normalized.Password) < 8 {
		problems = append(problems, "password must be at least 8 characters")
	}
	if len(problems) > 0 {
		return domain.UserSignupResult{}, errors.New(strings.Join(problems, "; "))
	}

	return s.provider.SignUp(ctx, normalized)
}

func looksLikePhone(phone string) bool {
	value := strings.TrimSpace(phone)
	if value == "" {
		return false
	}

	if value[0] == '+' {
		value = value[1:]
	}
	if len(value) < 7 || len(value) > 15 {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
