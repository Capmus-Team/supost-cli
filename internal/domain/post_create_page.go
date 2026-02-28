package domain

const (
	PostCreateStageChooseCategory    = "choose_category"
	PostCreateStageChooseSubcategory = "choose_subcategory"
	PostCreateStageForm              = "form"
)

// PostCreatePage is the staged view model for post creation.
type PostCreatePage struct {
	Stage           string        `json:"stage" db:"-"`
	CategoryID      int64         `json:"category_id" db:"-"`
	SubcategoryID   int64         `json:"subcategory_id" db:"-"`
	CategoryName    string        `json:"category_name" db:"-"`
	SubcategoryName string        `json:"subcategory_name" db:"-"`
	Categories      []Category    `json:"categories" db:"-"`
	Subcategories   []Subcategory `json:"subcategories" db:"-"`
}
