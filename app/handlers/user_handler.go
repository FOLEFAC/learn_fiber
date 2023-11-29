package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func createUserHandler(c *fiber.Ctx) error {
	// Create new User struct
	signUp := &User{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Create database connection.
	db, err := OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Create a new validator for a User model.
	validate := NewValidator()
	user := &User{}
	// Set initialized default data for user:
	user.Email = signUp.Email
	user.PasswordHash = GeneratePassword(signUp.PasswordHash)
	user.CreatedAt = time.Now()
	user.Id = uuid.New()
	// Validate book fields.
	fmt.Println(user)
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   ValidatorErrors(err),
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
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  user,
	})
}
