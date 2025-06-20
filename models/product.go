package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
}

type ProductCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductManagementParameter struct {
	Action string `json:"action"`
	Product
}

type ProductCategoryManagementParameter struct {
	Action string `json:"action"`
	ProductCategory
}
