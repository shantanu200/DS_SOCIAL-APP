package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetLocalUser(c *fiber.Ctx) (float64, error) {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)

	fmt.Println(claims);

	id := claims["id"].(float64)

	return id, nil
}
