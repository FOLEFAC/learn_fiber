package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"

	swagger "github.com/arsmn/fiber-swagger"
)

type Post struct {
	Id        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Title     string    `db:"title" json:"title" validate:"required,min=10"`
	Content   string    `db:"content" json:"content" validate:"required"`
	Published bool      `db:"published" json:"published"`
}

// Queries struct for collect all app queries.
type Queries struct {
	*PostQueries // load queries from Book model
}

// BookQueries struct for queries from Book model.
type PostQueries struct {
	*sqlx.DB
}

// GetBook method for getting one book by given ID.
func (q *PostQueries) GetPost(id uuid.UUID) (Post, error) {
	// Define book variable.
	post := Post{}

	// Define query string.
	query := `SELECT * FROM posts WHERE id = $1`

	// Send query to database.
	err := q.Get(&post, query, id)
	if err != nil {
		// Return empty object and error.
		return post, err
	}

	// Return query result.
	return post, nil
}

// GetUsers method for getting all users.
func (q *PostQueries) GetPosts() ([]Post, error) {
	// Define users variable.
	posts := []Post{}

	// Define query string.
	query := `SELECT * FROM posts`

	// Send query to database.
	err := q.Select(&posts, query)
	if err != nil {
		// Return empty object and error.
		return posts, err
	}

	// Return query result.
	return posts, nil
}

// CreateBook method for creating book by given Book object.
func (q *PostQueries) CreatePost(p *Post) error {
	// Define query string.
	query := `INSERT INTO posts VALUES ($1, $2, $3, $4, $5)`
	//fmt.Println(strconv.FormatBool(p.Published))

	// Send query to database.
	_, err := q.Exec(query, p.Id, p.CreatedAt, p.Title, p.Content, strconv.FormatBool(p.Published))
	if err != nil {
		// Return only error.
		fmt.Println("okay", err)
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateBook method for updating book by given Book object.
func (q *PostQueries) UpdatePost(id uuid.UUID, p *Post) error {
	// Define query string.
	query := `UPDATE posts SET title = $2, content = $3, published = $4 WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id, p.Title, p.Content, p.Published)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		PostQueries: &PostQueries{DB: db}, // from Post model
	}, nil
}

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*sqlx.DB, error) {

	envFile, err := godotenv.Read(".env")
	if err != nil {
		return nil, fmt.Errorf("can't read env file")
	}

	// Define database connection settings.
	maxConn, _ := strconv.Atoi(envFile["DB_MAX_CONNECTIONS"])
	maxIdleConn, _ := strconv.Atoi(envFile["DB_MAX_IDLE_CONNECTIONS"])
	maxLifetimeConn, _ := strconv.Atoi(envFile["DB_MAX_LIFETIME_CONNECTIONS"])

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", envFile["DB_SERVER_URL"])
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings.
	db.SetMaxOpenConns(maxConn)                           // the default is 0 (unlimited)
	db.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

func getSinglePostHandler(c *fiber.Ctx) error {
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
func getPostsHandler(ctx *fiber.Ctx) error {

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

func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := make(map[string]string) // or map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		//fmt.Println("field", err.Field(), err.Tag(), err.Param(),err.Error())
		fields[err.Field()] = err.Error()
	}
	fmt.Println(fields)
	return fields
}

func createPostHandler(c *fiber.Ctx) error {

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

// @title Fiber Swagger API
// @version 2.0
// @description This is an auto-generated API docs.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {

	app := fiber.New()
	app.Get("/posts", getPostsHandler)
	app.Get("/post/:id", getSinglePostHandler)
	app.Post("/user/new", createPostHandler)
	app.Patch("/user/update", UpdatePostHandler)
	//app.Delete("/user/delete/:id", deleteUserHandler)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World")
	})

	app.Get("/swagger/*", swagger.Handler) // default

	app.Listen(":4000")

}
