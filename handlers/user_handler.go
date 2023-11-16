package handlers

import (
	"product/middleware"
	"product/models"
	"product/services"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{
		userService,
	}
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	userRequest := models.UserRequest{}
	c.BodyParser(&userRequest)

	user, _ := uh.userService.GetByCondition("email", userRequest.Email)

	if user.ID == 0 {
		return response(c, fiber.StatusBadRequest, "email is not registered", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		return response(c, fiber.StatusBadRequest, "password invalid", nil)
	}

	token, err := middleware.GenerateToken(user, 6)

	if err != nil {
		return response(c, fiber.StatusInternalServerError, "failed to generate token", nil)
	}

	return response(c, fiber.StatusOK, "login success", fiber.Map{"token": token})
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	userRequest := models.UserRequest{}
	c.BodyParser(&userRequest)

	if err := userRequest.Validate(); err != nil {
		return response(c, fiber.StatusBadRequest, "invalid request", nil)
	}

	if len(userRequest.Password) < 6 {
		return response(c, fiber.StatusBadRequest, "password must be at least 6 character", nil)
	}

	user, _ := uh.userService.GetByCondition("email", userRequest.Email)

	if user.ID != 0 {
		return response(c, fiber.StatusConflict, "email has been registered", nil)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)

	userRequest.Password = string(password)

	user, err := uh.userService.Create(userRequest.ConvertToUser())

	if user.ID == 0 || err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	return response(c, fiber.StatusOK, "successfully regist user", user.ConvertToResponse())
}

func (uh *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := uh.userService.GetAll()

	if err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	if len(users) == 0 {
		return response(c, fiber.StatusNoContent, "", nil)
	}

	var usersResponse []models.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, user.ConvertToResponse())
	}

	return response(c, fiber.StatusOK, "successfully get all users", usersResponse)
}

func (uh *UserHandler) Update(c *fiber.Ctx) error {
	userRequest := models.UserRequest{}
	id := c.Params("id")

	c.BodyParser(&userRequest)

	if err := userRequest.Validate(); err != nil || len(userRequest.Password) < 6 {
		return response(c, fiber.StatusBadRequest, "invalid request", nil)
	}

	user, err := uh.userService.GetByCondition("id", id)

	if err != nil {
		return response(c, fiber.StatusBadRequest, "user is not found", nil)
	}

	user.Email = userRequest.Email
	user.Name = userRequest.Name

	user, err = uh.userService.Update(id, user)

	if err != nil {
		return response(c, fiber.StatusInternalServerError, "Upps Sorry, There is something wrong in server", nil)
	}

	return response(c, fiber.StatusOK, "successfully update user", user.ConvertToResponse())
}
