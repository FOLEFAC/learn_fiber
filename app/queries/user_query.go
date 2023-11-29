package queries

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

// CreateUser method for creating book by given User object.
func (q *UserQueries) CreateUser(u *User) error {
	// Define query string.
	query := `INSERT INTO users VALUES ($1, $2, $3, $4)`

	// Send query to database.
	_, err := q.Exec(query, u.Id, u.CreatedAt, u.Email, u.PasswordHash)
	if err != nil {
		// Return only error.
		fmt.Println("okay", err)
		return err
	}

	// This query returns nothing.
	return nil
}
