package entity

import (
	"github.com/google/uuid"
)

// Image entity Image
type Image struct {
	File []byte
	Name string
}

// NewImage create a new entity Image
func NewImage() ID {
	return ID{
		value: uuid.New(),
	}
}

// FileName
func (i Image) FileName() string {
	return i.Name
}
