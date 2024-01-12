package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	jwtware "github.com/gofiber/contrib/jwt"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/contrib/jwt
func JWTProtected() func(*fiber.Ctx) error {

	envFile, _ := godotenv.Read(".env")

	// Create config for JWT authentication middleware.
	config := jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(envFile["JWT_SECRET_KEY"])},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		fmt.Println("here we go")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	fmt.Println("there is an error", err.Error())
	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
