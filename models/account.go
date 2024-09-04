package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id         uuid.UUID
	Name       string
	Address    string
	LastLease  time.Time
	TotalLease float64
}
