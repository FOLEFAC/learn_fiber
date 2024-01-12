package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/FOLEFAC/learn_fiber/app/models"
	"github.com/FOLEFAC/learn_fiber/pkg/utils"
	"github.com/FOLEFAC/learn_fiber/platform/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type RefreshStore struct {
	refresh_token string
}

// Converts from []byte to a json object according to the User struct.
func toJson(val []byte) RefreshStore {
	refresh_store := RefreshStore{}
	err := json.Unmarshal(val, &refresh_store)
	if err != nil {
		panic(err)
	}
	return refresh_store
}

// RenewTokens method for renew access and refresh tokens.
// @Description Renew access and refresh tokens.
// @Summary renew access and refresh tokens
// @Tags Token
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func RenewTokens(c *fiber.Ctx) error {

	// Create a new renew refresh token struct.
	renewal := &models.Renew{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renewal); err != nil {

		// Return, if JSON data is not correct.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	// expiresAccessToken := claims.Expires

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renewal.Refresh)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now < expiresRefreshToken {
		// Define user ID.
		userID := claims.UserID
		envFile, err := godotenv.Read(".env")

		cache := redis.NewClient(&redis.Options{
			Addr: envFile["REDIS_ADDRESS"] + ":" + envFile["REDIS_PORT"],
		})

		val := cache.Get(c.Context(), userID.String())

		string_val := strings.Split(val.String(), ":")

		if len(string_val) > 2 {
			// Return, if user not found.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "user with the given ID needs to sign in",
			})
		}

		// Create database connection.
		db, err := database.OpenDBConnection()
		if err != nil {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Get user by ID.
		selectedUser, err := db.GetUserByID(userID)
		if err != nil {
			// Return, if user not found.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "user with the given ID is not found",
			})
		}

		// Generate JWT Access & Refresh tokens.
		tokens, err := utils.GenerateNewTokens(userID.String(), selectedUser.Email)
		if err != nil {
			// Return status 500 and token generation error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		refresh_expiry, _ := strconv.Atoi(envFile["JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"])

		cacheErr := cache.Set(
			context.Background(), userID.String(), tokens.Refresh, time.Duration(refresh_expiry)*time.Hour).Err()

		if cacheErr != nil {
			return cacheErr
		}

		return c.JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}
}
