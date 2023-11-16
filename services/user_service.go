package services

import (
	"product/models"

	"gorm.io/gorm"
)

type UserService interface {
	GetAll() ([]models.User, error)
	GetByCondition(key string, value string) (models.User, error)
	Create(userRequest models.User) (models.User, error)
	Update(id string, userRequest models.User) (models.User, error)
}

func NewUserService(gormDB *gorm.DB) UserService {
	return &UserServiceImpl{
		db: gormDB,
	}
}

type UserServiceImpl struct {
	db *gorm.DB
}

func (us *UserServiceImpl) GetAll() ([]models.User, error) {
	var users []models.User

	err := us.db.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserServiceImpl) GetByCondition(key string, value string) (models.User, error) {
	var user models.User

	err := us.db.First(&user, key, value).Error

	if err != nil {
		return models.User{}, err
	}

	return user, err
}

func (us *UserServiceImpl) Create(userRequest models.User) (models.User, error) {
	var user models.User

	rec := us.db.Create(&userRequest)

	if rec.Error != nil {
		return models.User{}, rec.Error
	}

	rec.Last(&user)

	return user, nil
}

func (us *UserServiceImpl) Update(id string, user models.User) (models.User, error) {
	rec := us.db.Save(&user)

	if rec.Error != nil {
		return models.User{}, rec.Error
	}

	rec.Last(&user)

	return user, nil
}
