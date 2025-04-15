package handler

import (
	"errors"
	"net/mail"
	"rentincome/database"
	"rentincome/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// auth

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// User

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	print(p)
	print(user.Password)
	return bool(CheckPasswordHash(p, user.Password))
}

func ReturnError(c *fiber.Ctx, options ...interface{}) error {
	var data interface{} = nil
	var message string = ""
	var status int = 500

	// Parse options
	for _, option := range options {
		switch v := option.(type) {
		case error:
			message = v.Error()
		case interface{}:
			data = v
		case string:
			message = v
		case int:
			status = v
		}
	}

	response := fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	}

	return c.Status(status).JSON(response)
}

func ReturnOK(c *fiber.Ctx, options ...interface{}) error {
	var data interface{} = nil
	var message string = ""
	var status int = 200

	// Parse options
	for _, option := range options {
		switch v := option.(type) {
		case string:
			message = v
		case int:
			status = v
		case struct{}:
			data = v
		case interface{}:
			// its default basically
			data = v
		}
	}

	response := fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	}

	return c.Status(status).JSON(response)
}

func GetFromJwt(c *fiber.Ctx, claimKey string) (interface{}, error) {
	// Retrieve the parsed JWT token from the context
	j, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT token")
	}

	// Extract the claims from the token
	claims, ok := j.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	// Retrieve the specific claim by key
	value, exists := claims[claimKey]
	if !exists {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Claim not found in JWT")
	}
	return value, nil
}
