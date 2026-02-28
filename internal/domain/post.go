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
	Photo1File    string    `json:"photo1_file_name" db:"photo1_file_name"`
	Photo2File    string    `json:"photo2_file_name" db:"photo2_file_name"`
	Photo3File    string    `json:"photo3_file_name" db:"photo3_file_name"`
	Photo4File    string    `json:"photo4_file_name" db:"photo4_file_name"`
	ImageSource1  string    `json:"image_source1" db:"image_source1"`
	ImageSource2  string    `json:"image_source2" db:"image_source2"`
	ImageSource3  string    `json:"image_source3" db:"image_source3"`
	ImageSource4  string    `json:"image_source4" db:"image_source4"`
	Status        int       `json:"status" db:"status"`
	TimePosted    int64     `json:"time_posted" db:"time_posted"`
	TimePostedAt  time.Time `json:"time_posted_at" db:"time_posted_at"`
	Price         float64   `json:"price" db:"price"`
	HasPrice      bool      `json:"has_price" db:"has_price"`
	HasImage      bool      `json:"has_image" db:"has_image"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
