package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestSupabaseAuthSignupClient_SignUp_SendsExpectedPayload(t *testing.T) {
	var (
		gotMethod string
		gotPath   string
		gotAPIKey string
		gotAuth   string
		gotBody   supabaseSignupRequest
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAPIKey = r.Header.Get("apikey")
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":{"id":"abc123","email":"user@example.com","user_metadata":{"display_name":"Greg","phone":"+16505551234"},"created_at":"2026-02-28T21:00:00Z"}}`))
	}))
	defer server.Close()

	client, err := NewSupabaseAuthSignupClient(server.URL, "sb_publishable_test", "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	result, err := client.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "Greg",
		Email:       "user@example.com",
		Phone:       "+16505551234",
		Password:    "password123",
	})
	if err != nil {
		t.Fatalf("signup request failed: %v", err)
	}

	if gotMethod != http.MethodPost || gotPath != "/auth/v1/signup" {
		t.Fatalf("unexpected request %s %s", gotMethod, gotPath)
	}
	if gotAPIKey != "sb_publishable_test" {
		t.Fatalf("unexpected apikey header: %q", gotAPIKey)
	}
	if gotAuth != "Bearer sb_publishable_test" {
		t.Fatalf("unexpected authorization header: %q", gotAuth)
	}
	if gotBody.Data["display_name"] != "Greg" || gotBody.Data["phone"] != "+16505551234" {
		t.Fatalf("unexpected metadata payload: %+v", gotBody.Data)
	}
	if gotBody.Password != "password123" {
		t.Fatalf("expected provided password in signup payload, got %q", gotBody.Password)
	}
	if result.UserID != "abc123" {
		t.Fatalf("unexpected user id: %q", result.UserID)
	}
}

func TestSupabaseAuthSignupClient_SignUp_ReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_signup"}`))
	}))
	defer server.Close()

	client, err := NewSupabaseAuthSignupClient(server.URL, "sb_publishable_test", "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	_, err = client.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "Greg",
		Email:       "user@example.com",
		Phone:       "+16505551234",
		Password:    "password123",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "supabase signup failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSupabaseAuthSignupClient_SignUp_AcceptsTopLevelUserShape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"top-level-1","email":"top@example.com","phone":"+16505550000","created_at":"2026-02-28T21:00:00Z","user_metadata":{"display_name":"Top Level","phone":"+16505550000"}}`))
	}))
	defer server.Close()

	client, err := NewSupabaseAuthSignupClient(server.URL, "sb_publishable_test", "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	result, err := client.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "Top Level",
		Email:       "top@example.com",
		Phone:       "+16505550000",
		Password:    "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UserID != "top-level-1" {
		t.Fatalf("unexpected user id: %q", result.UserID)
	}
}

func TestSupabaseAuthSignupClient_SignUp_NoUserIncludesMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":null,"session":null,"msg":"User already registered"}`))
	}))
	defer server.Close()

	client, err := NewSupabaseAuthSignupClient(server.URL, "sb_publishable_test", "")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	_, err = client.SignUp(context.Background(), domain.UserSignupSubmission{
		DisplayName: "Greg",
		Email:       "user@example.com",
		Phone:       "+16505551234",
		Password:    "password123",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "returned no user") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "User already registered") {
		t.Fatalf("expected detailed message in error, got: %v", err)
	}
}
