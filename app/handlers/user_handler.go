package handlers

//export PATH=$PATH:/usr/local/go/bin
import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FOLEFAC/learn_fiber/app/models"
	"github.com/FOLEFAC/learn_fiber/pkg/utils"
	"github.com/FOLEFAC/learn_fiber/platform/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &models.SignIn{}
	fmt.Println("email get before", signIn.Email)
	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	fmt.Println("email get", signIn.Email)

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	fmt.Println("get me that email boyyyyyyyyyy", signIn.Email)
	// Get user by email.
	selectedUser, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		// Return, if user not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given email is not found",
		})
	}
	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(selectedUser.PasswordHash, signIn.Password)
	fmt.Println(compareUserPassword)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "wrong user email address or password",
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(selectedUser.Id.String(), selectedUser.Email)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userID := selectedUser.Id.String()

	envFile, err := godotenv.Read(".env")

	cache := redis.NewClient(&redis.Options{
		Addr: envFile["REDIS_ADDRESS"] + ":" + envFile["REDIS_PORT"],
	})

	// // Create a new Redis connection.
	// connRedis, err := cache.RedisConnection()
	// if err != nil {
	// 	// Return status 500 and Redis connection error.
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   err.Error(),
	// 	})
	// }
	refresh_expiry, _ := strconv.Atoi(envFile["JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"])

	cacheErr := cache.Set(context.Background(), userID, tokens.Refresh, time.Duration(refresh_expiry)*time.Hour).Err()
	if cacheErr != nil {
		return cacheErr
	}

	// // Save refresh token to Redis.
	// errSaveToRedis := connRedis.Set(context.Background(), userID, tokens.Refresh, 0).Err()
	// if errSaveToRedis != nil {
	// 	// Return status 500 and Redis connection error.
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   errSaveToRedis.Error(),
	// 	})
	// }

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

func CreateUserHandler(c *fiber.Ctx) error {
	// Create new User struct
	signUp := &models.SignUp{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
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

	// Create a new validator for a User model.
	validate := utils.NewValidator()
	user := &models.User{}

	// Set initialized default data for user:
	user.Email = signUp.Email
	user.PasswordHash = utils.GeneratePassword(signUp.Password)
	user.CreatedAt = time.Now()
	user.Id = uuid.New()

	// Validate book fields.
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}
	// Create User .
	if err := db.CreateUser(user); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.PasswordHash = ""
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  user,
	})
}
