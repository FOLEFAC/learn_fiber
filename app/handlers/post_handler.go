package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSinglePostHandler(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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

	// Get book by ID.
	post, err := db.GetPost(id)
	if err != nil {
		// Return, if book not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with the given ID is not found",
			"post":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  post,
	})
}

// GetBooks func gets all existing posts.
// @Description Get all existing posts.
// @Summary get all existing posts
// @Tags Posts
// @Accept */*
// @Produce json
// @Success 200 {array} Post
// @Router / [get]
func GetPostsHandler(ctx *fiber.Ctx) error {

	// Create database connection.
	db, err := OpenDBConnection()

	if err != nil {
		fmt.Println(err.Error())
		// Return status 500 and database connection error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all users.
	posts, err := db.GetPosts()
	fmt.Println("erroror", err)
	if err != nil {
		// Return, if users not found.
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "posts were not found",
			"count": 0,
			"users": nil,
		})
	}

	// Return status 200 OK.
	return ctx.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(posts),
		"users": posts,
	})

}

func CreatePostHandler(c *fiber.Ctx) error {

	// Create new Book struct
	post := &Post{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(post); err != nil {
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

	// Create a new validator for a Book model.
	validate := NewValidator()

	// Set initialized default data for book:
	post.Id = uuid.New()
	post.CreatedAt = time.Now()
	post.Published = false
	// Validate book fields.
	if err := validate.Struct(post); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   ValidatorErrors(err),
		})
	}

	// Delete book by given ID.
	if err := db.CreatePost(post); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  post,
	})
}

func UpdatePostHandler(c *fiber.Ctx) error {

	// Create new Book struct
	post := &Post{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(post); err != nil {
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

	// Checking, if book with given ID is exists.
	foundPost, err := db.GetPost(post.Id)
	if err != nil {
		// Return status 404 and book not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with this ID not found",
		})
	}

	// Create a new validator for a Book model.
	validate := NewValidator()

	// Validate book fields.
	if err := validate.Struct(post); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   ValidatorErrors(err),
		})
	}

	// Update book by given ID.
	if err := db.UpdatePost(foundPost.Id, post); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201.
	return c.SendStatus(fiber.StatusOK)
}
