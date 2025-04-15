package handler

import (
	"rentincome/database"
	"rentincome/model"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.Find(&user, id)
	if user.Username == "" {
		return ReturnError(c, 404)
	}
	return ReturnOK(c, user)
}

func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ID       uint   `json:"id"`
	}

	db := database.DB
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return ReturnError(c, err)
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return ReturnError(c, err, "bad password")
	}

	user.Password = hash
	err = db.Create(&user).Error
	if err != nil {
		return ReturnError(c, err)
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
		ID:       user.ID,
	}
	return ReturnOK(c, newUser)
}

func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Username string `json:"username" validate:"omitempty,min=1" `
		Names    string `json:"names" validate:"omitempty,min=1" `
		Email    string `json:"email" validate:"omitempty,min=1" `
		Phone    string `json:"phone" validate:"omitempty,min=1" `
		Password string `json:"password" validate:"omitempty,min=1" `
	}
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	user := new(model.User)
	db := database.DB
	err = db.First(&user, id).Error
	if err != nil {
		ReturnError(c, 404, err)
	}

	var input UpdateUserInput
	err = c.BodyParser(&input)
	if err != nil {
		return ReturnError(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ReturnError(c, fiber.StatusUnprocessableEntity, err)
	}

	// Update the user fields
	user.Username = input.Username
	user.Names = input.Names
	user.Email = input.Email
	user.Phone = input.Phone

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ReturnError(c, "Failed to hash password")
	}
	user.Password = string(hashedPassword)

	err = db.Save(&user).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB
	var user model.User

	db.First(&user, id)

	db.Delete(&user)
	return ReturnOK(c)
}
