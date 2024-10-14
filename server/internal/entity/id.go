package entity

import "github.com/google/uuid"

type ID = uuid.UUID

// NewID returns a new, unique ID.
func NewID() ID {
	return ID(uuid.New())
}

// ParseID parses the given string as an ID and returns the result. If the string is
// invalid, it returns an error.
func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}
