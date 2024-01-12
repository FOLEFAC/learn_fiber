package database

import "github.com/FOLEFAC/learn_fiber/app/queries"

// Queries struct for collect all app queries.
type Queries struct {
	*queries.PostQueries // load queries from Post Model
	*queries.UserQueries // load queries from User Model
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
		PostQueries: &queries.PostQueries{DB: db}, // from Post model
		UserQueries: &queries.UserQueries{DB: db}, // from Post model

	}, nil
}
