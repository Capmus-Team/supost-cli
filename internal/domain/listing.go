package domain

import "time"

// Listing maps to the Supabase "listings" table.
// TypeScript equivalent: interface Listing { id: string; user_id: string; ... }
type Listing struct {
	ID          string    `json:"id"          db:"id"`
	UserID      string    `json:"user_id"     db:"user_id"`
	Title       string    `json:"title"       db:"title"`
	Description string    `json:"description" db:"description"`
	Price       int       `json:"price"       db:"price"`       // cents
	Category    string    `json:"category"    db:"category"`
	Status      string    `json:"status"      db:"status"`      // active, sold, expired
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  db:"updated_at"`
}

func (l *Listing) Validate() error {
	if l.Title == "" {
		return ErrMissingTitle
	}
	if l.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
