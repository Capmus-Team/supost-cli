package domain

// Category maps to public.category.
type Category struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ShortName string `json:"short_name" db:"short_name"`
}

// Subcategory maps to public.subcategory.
type Subcategory struct {
	ID         int64  `json:"id" db:"id"`
	CategoryID int64  `json:"category_id" db:"category_id"`
	Name       string `json:"name" db:"name"`
}

// CategoryWithSubcategories is the utility output view for category browsing.
type CategoryWithSubcategories struct {
	ID            int64         `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	ShortName     string        `json:"short_name" db:"short_name"`
	Subcategories []Subcategory `json:"subcategories" db:"-"`
}
