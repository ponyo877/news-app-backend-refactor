package entity

import "github.com/google/uuid"

// ID entity ID
type ID struct {
	value uuid.UUID
}

// NewID create a new entity ID
func NewID() ID {
	return ID{
		value: uuid.New(),
	}
}

// StringToID convert a string to an entity ID
func StringToID(s string) (ID, error) {
	UUID, err := uuid.Parse(s)
	return ID{
		value: UUID,
	}, err
}

func (i ID) String() string {
	return i.value.String()
}
