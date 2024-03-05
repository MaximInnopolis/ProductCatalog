package models

type Product struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}
