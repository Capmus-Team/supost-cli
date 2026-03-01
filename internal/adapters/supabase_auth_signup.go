package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const defaultSupabaseAuthTimeout = 15 * time.Second

// SupabaseAuthSignupClient creates users through Supabase Auth signup endpoint.
type SupabaseAuthSignupClient struct {
	baseURL  string
	apiKey   string
	adminKey string
	client   *http.Client
}

type supabaseSignupRequest struct {
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Data     map[string]string `json:"data,omitempty"`
}

type supabaseSignupResponse struct {
	User *supabaseSignupUser `json:"user"`

	// Some Supabase/Auth versions return the user fields at the top level.
	ID           string         `json:"id"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	CreatedAtRaw string         `json:"created_at"`
	UserMetadata map[string]any `json:"user_metadata"`

	// Error/message shapes seen from Auth responses.
	Msg              string `json:"msg"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        string `json:"error_code"`
}

type supabaseSignupUser struct {
	ID           string         `json:"id"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	CreatedAtRaw string         `json:"created_at"`
	UserMetadata map[string]any `json:"user_metadata"`
}

type supabaseAdminCreateUserRequest struct {
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	EmailConfirm bool              `json:"email_confirm"`
	UserMetadata map[string]string `json:"user_metadata,omitempty"`
}

// NewSupabaseAuthSignupClient builds a signup client for Supabase Auth.
func NewSupabaseAuthSignupClient(baseURL string, apiKey string, adminKey string) (*SupabaseAuthSignupClient, error) {
	trimmedBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmedBaseURL == "" {
		return nil, fmt.Errorf("supabase_url is required")
	}
	trimmedAPIKey := strings.TrimSpace(apiKey)
	if trimmedAPIKey == "" {
		return nil, fmt.Errorf("supabase publishable/anon key is required")
	}
	return &SupabaseAuthSignupClient{
		baseURL:  trimmedBaseURL,
		apiKey:   trimmedAPIKey,
		adminKey: strings.TrimSpace(adminKey),
		client:   &http.Client{Timeout: defaultSupabaseAuthTimeout},
	}, nil
}

// SignUp creates a user through POST /auth/v1/signup.
func (c *SupabaseAuthSignupClient) SignUp(ctx context.Context, submission domain.UserSignupSubmission) (domain.UserSignupResult, error) {
	payload := supabaseSignupRequest{
		Email:    submission.Email,
		Password: submission.Password,
		Data: map[string]string{
			"display_name": submission.DisplayName,
			"phone":        submission.Phone,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("encoding signup payload: %w", err)
	}

	raw, statusCode, err := c.sendSupabaseRequest(ctx, "/auth/v1/signup", c.apiKey, body)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("sending signup request: %w", err)
	}

	if statusCode < 200 || statusCode >= 300 {
		if statusCode == http.StatusTooManyRequests && c.adminKey != "" && isEmailSendRateLimit(raw) {
			return c.createUserViaAdmin(ctx, submission)
		}
		return domain.UserSignupResult{}, fmt.Errorf("supabase signup failed: status %d: %s", statusCode, strings.TrimSpace(string(raw)))
	}

	user, err := parseSupabaseSignupUser(raw)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("decoding signup response: %w", err)
	}
	return signupResultFromUser(submission, user, true), nil
}

func (c *SupabaseAuthSignupClient) createUserViaAdmin(ctx context.Context, submission domain.UserSignupSubmission) (domain.UserSignupResult, error) {
	payload := supabaseAdminCreateUserRequest{
		Email:        submission.Email,
		Password:     submission.Password,
		EmailConfirm: true,
		UserMetadata: map[string]string{
			"display_name": submission.DisplayName,
			"phone":        submission.Phone,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("encoding admin signup payload: %w", err)
	}

	raw, statusCode, err := c.sendSupabaseRequest(ctx, "/auth/v1/admin/users", c.adminKey, body)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("sending admin signup request: %w", err)
	}
	if statusCode < 200 || statusCode >= 300 {
		return domain.UserSignupResult{}, fmt.Errorf("supabase admin create user failed: status %d: %s", statusCode, strings.TrimSpace(string(raw)))
	}

	user, err := parseSupabaseSignupUser(raw)
	if err != nil {
		return domain.UserSignupResult{}, fmt.Errorf("decoding admin signup response: %w", err)
	}
	return signupResultFromUser(submission, user, false), nil
}

func (c *SupabaseAuthSignupClient) sendSupabaseRequest(ctx context.Context, path string, authKey string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", authKey)
	req.Header.Set("Authorization", "Bearer "+authKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if readErr != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response body: %w", readErr)
	}
	return raw, resp.StatusCode, nil
}

func parseSupabaseSignupUser(raw []byte) (*supabaseSignupUser, error) {
	var decoded supabaseSignupResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}

	user := decoded.User
	if user == nil && strings.TrimSpace(decoded.ID) != "" {
		user = &supabaseSignupUser{
			ID:           decoded.ID,
			Email:        decoded.Email,
			Phone:        decoded.Phone,
			CreatedAtRaw: decoded.CreatedAtRaw,
			UserMetadata: decoded.UserMetadata,
		}
	}
	if user == nil {
		reason := strings.TrimSpace(decoded.Msg)
		if reason == "" {
			reason = strings.TrimSpace(decoded.ErrorDescription)
		}
		if reason == "" {
			reason = strings.TrimSpace(decoded.Error)
		}
		if reason == "" {
			reason = strings.TrimSpace(string(raw))
		}
		if reason == "" {
			reason = "empty response body"
		}
		return nil, fmt.Errorf("supabase signup returned no user: %s", reason)
	}
	return user, nil
}

func signupResultFromUser(submission domain.UserSignupSubmission, user *supabaseSignupUser, emailConfirmationSent bool) domain.UserSignupResult {
	createdAt := time.Time{}
	if user.CreatedAtRaw != "" {
		parsed, err := time.Parse(time.RFC3339Nano, user.CreatedAtRaw)
		if err == nil {
			createdAt = parsed
		}
	}

	displayName := submission.DisplayName
	phone := submission.Phone
	if user.UserMetadata != nil {
		if value := metadataString(user.UserMetadata, "display_name"); value != "" {
			displayName = value
		}
		if value := metadataString(user.UserMetadata, "phone"); value != "" {
			phone = value
		}
	}

	return domain.UserSignupResult{
		UserID:                strings.TrimSpace(user.ID),
		DisplayName:           displayName,
		Email:                 strings.TrimSpace(user.Email),
		Phone:                 phone,
		EmailConfirmationSent: emailConfirmationSent,
		CreatedAt:             createdAt,
	}
}

func isEmailSendRateLimit(raw []byte) bool {
	var decoded supabaseSignupResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(decoded.ErrorCode), "over_email_send_rate_limit")
}

func metadataString(metadata map[string]any, key string) string {
	raw, ok := metadata[key]
	if !ok {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}
