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

type SearchProductParameter struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	SortBy   string  `json:"sort_by"`  // e.g., "price", "name"
	OrderBy  string  `json:"order_by"` // e.g., "asc", "desc"
	Limit    int64   `json:"limit"`    // Number of products to return
	Page     int64   `json:"page"`     // Page number for pagination
}

type SearchProductResponse struct {
	Products    []Product `json:"products"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
	TotalCount  int64     `json:"total_count"`
	TotalPages  int       `json:"total_pages"`
	NextPage    bool      `json:"next_page"`
	NextPageURL *string   `json:"next_page_url"`
}
