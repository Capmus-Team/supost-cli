package domain

import "time"

// Message maps to app_private.message.
type Message struct {
	ID        int64     `json:"id" db:"id"`
	PostID    int64     `json:"post_id" db:"post_id"`
	Message   string    `json:"message" db:"message"`
	IP        string    `json:"ip" db:"ip"`
	Email     string    `json:"email" db:"email"`
	RawEmail  string    `json:"raw_email" db:"raw_email"`
	Source    string    `json:"source" db:"source"`
	Status    string    `json:"status" db:"status"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Scammed   bool      `json:"scammed" db:"scammed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
