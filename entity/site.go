package entity

import "time"

type Site struct {
	ID            ID
	Title         string
	RSSURL        string
	LastUpdatedAt time.Time
}

// NewSite create a new site
func NewSite(title, RSSURL string) (*Site, error) {
	site := &Site{
		ID:            NewID(),
		Title:         title,
		RSSURL:        RSSURL,
		LastUpdatedAt: time.Time{},
	}

	if err := site.Validate(); err != nil {
		return nil, ErrInvalidEntity
	}
	return site, nil
}

func (s Site) IsNewerLastUpdatedAt(lastUpdatedAt time.Time) bool {
	return s.LastUpdatedAt.After(lastUpdatedAt) || s.LastUpdatedAt.Equal(lastUpdatedAt)
}

func (s Site) UpdateLastUpdatedAt(lastUpdatedAt time.Time) Site {
	s.LastUpdatedAt = lastUpdatedAt
	return s
}

// Validate validate data
func (s Site) Validate() error {
	if s.Title == "" || s.RSSURL == "" {
		return ErrInvalidEntity
	}

	return nil
}
