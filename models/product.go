package models

import "github.com/go-playground/validator/v10"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(250)"`
	Price       int
	Stock       int
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

type ProductRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	Price       int    `json:"price" form:"price" validate:"required"`
	Stock       int    `json:"stock" form:"stock" validate:"required"`
}

func (p *ProductRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(p)

	return err
}

func (p *Product) ConvertToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}
}

func (p *ProductRequest) ConvertToProduct() Product {
	return Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}
}
