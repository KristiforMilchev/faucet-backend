package dtos

import "github.com/google/uuid"

type EVMAuthRequest struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	Signature string    `json:"signature" binding:"required"`
}
