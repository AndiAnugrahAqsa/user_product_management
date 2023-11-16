package services

import (
	"product/models"

	"gorm.io/gorm"
)

type ProductService interface {
	GetAll() ([]models.Product, error)
	GetByCondition(key string, value string) (models.Product, error)
	Create(productRequest models.Product) (models.Product, error)
	Update(id string, productRequest models.Product) (models.Product, error)
	Delete(product models.Product) error
}

func NewProductService(gormDB *gorm.DB) ProductService {
	return &ProductServiceImpl{
		db: gormDB,
	}
}

type ProductServiceImpl struct {
	db *gorm.DB
}

func (ps *ProductServiceImpl) GetAll() ([]models.Product, error) {
	var products []models.Product

	err := ps.db.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductServiceImpl) GetByCondition(key string, value string) (models.Product, error) {
	var product models.Product

	err := ps.db.First(&product, key, value).Error

	if err != nil {
		return models.Product{}, err
	}

	return product, err
}

func (ps *ProductServiceImpl) Create(productRequest models.Product) (models.Product, error) {
	var product models.Product

	rec := ps.db.Create(&productRequest)

	if rec.Error != nil {
		return models.Product{}, rec.Error
	}

	rec.Last(&product)

	return product, nil
}

func (ps *ProductServiceImpl) Update(id string, product models.Product) (models.Product, error) {
	rec := ps.db.Save(&product)

	if rec.Error != nil {
		return models.Product{}, rec.Error
	}

	rec.Last(&product)

	return product, nil
}

func (ps *ProductServiceImpl) Delete(product models.Product) error {
	rec := ps.db.Delete(&product)

	if rec.Error != nil {
		return rec.Error
	}

	return nil
}
