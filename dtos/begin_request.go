package dtos

import "time"

type BeginRequest struct {
	Key       string
	Timestamp time.Time
}
