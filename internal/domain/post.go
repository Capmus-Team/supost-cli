package domain

import "time"

const (
	// PostStatusActive matches public.post.status = 1.
	PostStatusActive = 1
)

// Post maps to the Supabase public.post table.
type Post struct {
	ID            int64     `json:"id" db:"id"`
	CategoryID    int64     `json:"category_id" db:"category_id"`
	SubcategoryID int64     `json:"subcategory_id" db:"subcategory_id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	Body          string    `json:"body" db:"body"`
	Status        int       `json:"status" db:"status"`
	TimePosted    int64     `json:"time_posted" db:"time_posted"`
	TimePostedAt  time.Time `json:"time_posted_at" db:"time_posted_at"`
	Price         float64   `json:"price" db:"price"`
	HasPrice      bool      `json:"has_price" db:"has_price"`
	HasImage      bool      `json:"has_image" db:"has_image"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
