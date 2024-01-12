package utils

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      uuid.UUID
	Credentials map[string]bool
	Expires     int64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		fmt.Println("the error was 1")
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			fmt.Println("teh errrow is 2")
			return nil, err
		}

		// Expires time.
		expires := int64(claims["expires"].(float64))
		// User credentials.
		credentials := map[string]bool{
			"post:create": claims["post:create"].(bool),
		}

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}

	return nil, err
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)
	// fmt.Println("tokenstring", tokenString)
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	// fmt.Println("onlytoken", len(onlyToken))
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {

	envFile, _ := godotenv.Read(".env")
	return []byte(envFile["JWT_SECRET_KEY"]), nil
}
