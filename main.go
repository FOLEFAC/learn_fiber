package main

// migrate -path ./migrations -database "postgres://postgres:ladder-99@localhost:5433/cleeroute?sslmode=disable" up 1
// export PATH=$(go env GOPATH)/bin:$PATH// export PATH=$(go env GOPATH)/bin:$PATH

//"github.com/FOLEFAC/learn_fiber/tree/main/pkg/routes"
import (
	"github.com/FOLEFAC/learn_fiber/pkg/routes"
	"github.com/gofiber/fiber/v2"

	_ "https://github.com/FOLEFAC/learn_fiber/tree/main/docs" // load API Docs files (Swagger)

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

	// Define a new Fiber app with config.
	app := fiber.New()

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

}
