package routes

import (
	"github.com/FOLEFAC/learn_fiber/app/handlers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/posts", handlers.GetPostsHandler)
	route.Get("/post/:id", handlers.GetSinglePostHandler)

	// Routes for POST method:
	route.Post("/user/login", handlers.UserSignIn)
	route.Post("/user/new", handlers.CreateUserHandler) // register a new user

}
