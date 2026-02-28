package domain

import "time"

// PostCreateSubmission is the validated input payload for creating a post.
type PostCreateSubmission struct {
	CategoryID    int64     `json:"category_id" db:"category_id"`
	SubcategoryID int64     `json:"subcategory_id" db:"subcategory_id"`
	Name          string    `json:"name" db:"name"`
	Body          string    `json:"body" db:"body"`
	Email         string    `json:"email" db:"email"`
	Price         float64   `json:"price" db:"price"`
	PriceProvided bool      `json:"price_provided" db:"-"`
	AccessToken   string    `json:"access_token" db:"access_token"`
	PostedAt      time.Time `json:"posted_at" db:"time_posted_at"`
}

// PostCreatePersisted is the DB return payload after insert.
type PostCreatePersisted struct {
	PostID      int64     `json:"post_id" db:"id"`
	AccessToken string    `json:"access_token" db:"access_token"`
	PostedAt    time.Time `json:"posted_at" db:"time_posted_at"`
}

// PublishEmailMessage is the Mailgun payload for publish-link emails.
type PublishEmailMessage struct {
	From    string `json:"from" db:"-"`
	To      string `json:"to" db:"-"`
	Subject string `json:"subject" db:"-"`
	Text    string `json:"text" db:"-"`
}

// PostCreateSubmitResult is the command output for submission mode.
type PostCreateSubmitResult struct {
	DryRun      bool      `json:"dry_run" db:"-"`
	PostID      int64     `json:"post_id" db:"-"`
	AccessToken string    `json:"access_token" db:"-"`
	PublishURL  string    `json:"publish_url" db:"-"`
	PostedAt    time.Time `json:"posted_at" db:"-"`
	EmailTo     string    `json:"email_to" db:"-"`
	EmailSent   bool      `json:"email_sent" db:"-"`
	Subject     string    `json:"subject" db:"-"`
	Body        string    `json:"body" db:"-"`
}
