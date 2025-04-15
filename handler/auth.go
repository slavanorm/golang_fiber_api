package handler

import (
	"rentincome/config"
	"rentincome/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	var userData UserData

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error on login request", "data": err.Error()})
	}

	identity := input.Identity
	pass := input.Password
	userModel, err := new(model.User), *new(error)
	if isEmail(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(identity)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err.Error() })
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err.Error() })
	} else {
		userData = UserData{
			ID:       userModel.ID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid identity or password", "data": nil})
	}
    // handle JWT
    claims := jwt.MapClaims{
		//"username": userData.Username,
		"user_id": userData.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.Config("JwtKey")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	
	return c.JSON(fiber.Map{"message": "Success", "data": t})
}
