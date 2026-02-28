package domain

import "time"

const (
	CategoryCampusJob     int64 = 1
	CategoryJobsOffCampus int64 = 2
	CategoryHousing       int64 = 3
	CategoryForSale       int64 = 5
	CategoryServices      int64 = 7
	CategoryPersonals     int64 = 8
	CategoryCommunity     int64 = 9
)

// HomeCategorySection is the sidebar category view model for home.
type HomeCategorySection struct {
	CategoryID       int64     `json:"category_id" db:"category_id"`
	CategoryName     string    `json:"category_name" db:"category_name"`
	SubcategoryNames []string  `json:"subcategory_names" db:"subcategory_names"`
	LastPostedAt     time.Time `json:"last_posted_at" db:"last_posted_at"`
}
