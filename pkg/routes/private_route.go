package routes

import (
	"fmt"

	"github.com/FOLEFAC/learn_fiber/app/handlers"
	"github.com/FOLEFAC/learn_fiber/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {

	fmt.Println("i got here and i am happy")
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/post/new", middleware.JWTProtected(), handlers.CreatePostHandler) // create a new book
	// route.Post("/user/sign/out", middleware.JWTProtected(), handlers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), handlers.RenewTokens) // renew Access & Refresh tokens

	// // Routes for PUT method:
	// route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// // Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
}

// // PublicRoutes func for describe group of public routes.
// func PrivateRoutes(a *fiber.App) {
// 	// Create routes group.
// 	route := a.Group("/api/v1")

// 	// Routes for GET method:
// 	route.Get("/posts", handlers.GetPostsHandler)
// 	route.Get("/post/:id", handlers.GetSinglePostHandler)
// 	route.Post("/post/new", handlers.CreatePostHandler)
// 	route.Patch("/post/update", handlers.UpdatePostHandler)

// 	// Routes for POST method:
// 	route.Post("/user/login", handlers.UserSignIn)
// 	route.Post("/user/new", handlers.CreateUserHandler) // register a new user

// }
