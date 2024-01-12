package models

// Renew struct to describe refresh token object.
type Renew struct {
	Refresh string `json:"refresh" validate:"required"`
}
