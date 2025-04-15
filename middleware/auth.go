package middleware

import (
	"fmt"
	"rentincome/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func jwtError(c *fiber.Ctx, err error) error {
	fmt.Println(err.Error())
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"message": "Invalid or expired JWT", "data": nil})
}

func JwtProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("JwtKey"))},
		ErrorHandler: jwtError,
	})
}
