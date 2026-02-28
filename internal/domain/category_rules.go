package domain

// CategoryPriceRequired reports whether posts in this category require a price.
func CategoryPriceRequired(categoryID int64) bool {
	return categoryID == CategoryHousing || categoryID == CategoryForSale
}

// CategoryPriceAllowed reports whether posts in this category may include a price.
func CategoryPriceAllowed(categoryID int64) bool {
	return CategoryPriceRequired(categoryID)
}
