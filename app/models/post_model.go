package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Title     string    `db:"title" json:"title" validate:"required,min=10"`
	Content   string    `db:"content" json:"content" validate:"required"`
	Published bool      `db:"published" json:"published"`
}
