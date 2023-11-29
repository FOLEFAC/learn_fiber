package database

import (
	_ "github.com/FOLEFAC/learn_fiber/pkg/routes"
)

// Queries struct for collect all app queries.
type Queries struct {
	*PostQueries // load queries from Post Model
	*UserQueries // load queries from User Model
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
		UserQueries: &UserQueries{DB: db}, // from Post model

	}, nil
}
