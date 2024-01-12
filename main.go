package main

// cd /etc/postgresql/15/main
// sudo vi pg_hba.conf

// usermod -aG sudo martins

// migrate -path ./migrations -database "postgres://postgres:ladder-99@localhost:5433/cleeroute?sslmode=disable" up 1
// export PATH=$(go env GOPATH)/bin:$PATH// export PATH=$(go env GOPATH)/bin:$PATH
// export PATH=$PATH:/usr/local/go/bin
import (
	"fmt"
	"log"

	"github.com/FOLEFAC/learn_fiber/pkg/middleware"
	"github.com/FOLEFAC/learn_fiber/pkg/routes"

	"github.com/gofiber/fiber/v2"

	_ "github.com/FOLEFAC/learn_fiber/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title Fiber Swagger API
// @version 2.0
// @description This is an auto-generated API docs.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /
// @schemes http
func main() {
	fmt.Println("ehllo world")
	// Define a new Fiber app with config.
	app := fiber.New()

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register public routes for app.
	routes.PrivateRoutes(app) // Register private routes for app
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start Server
	if err := app.Listen(":4000"); err != nil {
		log.Fatal(err)
		fmt.Println("the error", err)
	}
}
