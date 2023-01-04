package entity

import "time"

type Comment struct {
	ID         ID
	UserName   string
	AvatarURL  string
	DeviceHash string
	Message    string
	Article    Article
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

// NewComment create a new article
func NewComment(userName, avatarURL, deviceHash, message string, article Article) (*Comment, error) {
	comment := &Comment{
		ID:         NewID(),
		UserName:   userName,
		AvatarURL:  avatarURL,
		DeviceHash: deviceHash,
		Message:    message,
		Article:    article,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
	if err := comment.Validate(); err != nil {
		return nil, ErrInvalidEntity
	}
	return comment, nil
}

// Validate validate data
func (c *Comment) Validate() error {
	if c.UserName == "" || c.AvatarURL == "" || c.DeviceHash == "" || c.Message == "" {
		return ErrInvalidEntity
	}
	return nil
}
