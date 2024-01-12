package models

import (
	"github.com/google/uuid"
)

type Vote struct {
	UserId uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	PostId uuid.UUID `db:"post_id" json:"post_id" validate:"required,uuid"`
}
