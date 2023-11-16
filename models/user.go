package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRequest struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password"`
}

func (u *UserRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(u)

	return err
}

func (u *User) ConvertToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *UserRequest) ConvertToUser() User {
	return User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
