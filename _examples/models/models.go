package models

// User Model Struct (Berfungsi sebagai DB Table, OpenAPI Spec, dan DTO Validation)
type User struct {
	ID    string `json:"id" path:"id" doc:"User Unique ID"`
	Name  string `json:"name" validate:"required" doc:"User Full Name"`
	Email string `json:"email" validate:"required,email" doc:"User Email Address"`
}

// Product Model Struct
type Product struct {
	ID    string  `json:"id" path:"id" doc:"Product Unique ID"`
	Name  string  `json:"name" validate:"required" doc:"Product Item Name"`
	Price float64 `json:"price" validate:"required" doc:"Price in USD"`
}
