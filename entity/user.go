package entity

import "time"

type User struct {
	ID         ID
	Name       string
	AvatarURL  string
	DeviceHash string
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

// NewUser create a new user
func NewUser(name, avatarURL, deviceHash string) (User, error) {
	user := User{
		ID:         NewID(),
		Name:       name,
		AvatarURL:  avatarURL,
		DeviceHash: deviceHash,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
	if err := user.Validate(); err != nil {
		return User{}, ErrInvalidEntity
	}
	return user, nil
}

// Validate validate data
func (c *User) Validate() error {
	if c.Name == "" || c.AvatarURL == "" || c.DeviceHash == "" {
		return ErrInvalidEntity
	}
	return nil
}
