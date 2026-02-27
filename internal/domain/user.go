package domain

import "time"

// User maps to the Supabase "profiles" table.
// TypeScript equivalent: interface User { id: string; email: string; ... }
// Keep types plain — string, int, time.Time, []string — for TypeScript portability.
type User struct {
	ID        string    `json:"id"         db:"id"`
	Email     string    `json:"email"      db:"email"`
	Name      string    `json:"name"       db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
