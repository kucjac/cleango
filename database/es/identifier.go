package es

import (
	"github.com/google/uuid"
)

// IdGenerator is the interface used by identity generators.
type IdGenerator interface {
	GenerateId() string
}

// UUIDGenerator implements IdGenerator interface. Generates UUID V4 identifier.
type UUIDGenerator struct{}

// GenerateId generates identified.
// Implements IdGenerator interface.
func (u UUIDGenerator) GenerateId() string {
	return uuid.New().String()
}
