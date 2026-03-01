package domain

import "time"

// UserSignupSubmission is the validated input payload for Supabase Auth signup.
type UserSignupSubmission struct {
	DisplayName string `json:"display_name" db:"-"`
	Email       string `json:"email" db:"email"`
	Phone       string `json:"phone" db:"phone"`
	Password    string `json:"password" db:"-"`
}

// UserSignupResult is the command output for a completed signup request.
type UserSignupResult struct {
	UserID                string    `json:"user_id" db:"id"`
	DisplayName           string    `json:"display_name" db:"-"`
	Email                 string    `json:"email" db:"email"`
	Phone                 string    `json:"phone" db:"phone"`
	EmailConfirmationSent bool      `json:"email_confirmation_sent" db:"-"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}
