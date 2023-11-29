package queries

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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
	// Define posts variable.
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
