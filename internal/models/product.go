package models

// TODO: change get rid of array and create new struct(probably) product_categories
type Product struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}
