package domain

import "time"

// PostRespondSubmission is the input payload for responding to a post.
type PostRespondSubmission struct {
	PostID    int64  `json:"post_id" db:"post_id"`
	Message   string `json:"message" db:"message"`
	ReplyTo   string `json:"reply_to" db:"reply_to"`
	IP        string `json:"ip" db:"ip"`
	UserAgent string `json:"user_agent" db:"user_agent"`
}

// ResponseEmailMessage is the Mailgun payload for post responses.
type ResponseEmailMessage struct {
	From    string `json:"from" db:"-"`
	To      string `json:"to" db:"-"`
	ReplyTo string `json:"reply_to" db:"-"`
	Subject string `json:"subject" db:"-"`
	Text    string `json:"text" db:"-"`
}

// PostRespondResult is the command output for post response sends.
type PostRespondResult struct {
	DryRun       bool      `json:"dry_run" db:"-"`
	PostID       int64     `json:"post_id" db:"-"`
	PostEmail    string    `json:"post_email" db:"-"`
	ReplyTo      string    `json:"reply_to" db:"-"`
	MessageID    int64     `json:"message_id" db:"-"`
	MessageSaved bool      `json:"message_saved" db:"-"`
	EmailSent    bool      `json:"email_sent" db:"-"`
	Subject      string    `json:"subject" db:"-"`
	Body         string    `json:"body" db:"-"`
	SentAt       time.Time `json:"sent_at" db:"-"`
}
