package utils

import "github.com/google/uuid"

// NewUUID creates a new random UUID.
func NewUUID() string {
	return uuid.New().String()
}
