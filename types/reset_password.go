package types

import (
	"time"

	"github.com/google/uuid"
)

type ResetPassword struct {
	Token uuid.UUID `json:"token"`
	Expiration time.Time `json:"expiration"`
	Email     string `json:"email"`
}