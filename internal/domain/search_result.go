package domain

// SearchResultPage is the paginated post result contract for search views.
type SearchResultPage struct {
	Query         string `json:"query" db:"-"`
	CategoryID    int64  `json:"category_id" db:"-"`
	SubcategoryID int64  `json:"subcategory_id" db:"-"`
	Page          int    `json:"page" db:"-"`
	PerPage       int    `json:"per_page" db:"-"`
	HasMore       bool   `json:"has_more" db:"-"`
	Posts         []Post `json:"posts" db:"-"`
}
