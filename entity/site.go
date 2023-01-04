package entity

import "time"

type Site struct {
	ID            ID
	Title         string
	RSSURL        string
	ImageURL      string
	LastUpdatedAt time.Time
}

// NewSite create a new site
func NewSite(title, RSSURL, ImageURL string) (Site, error) {
	site := Site{
		ID:            NewID(),
		Title:         title,
		RSSURL:        RSSURL,
		ImageURL:      ImageURL,
		LastUpdatedAt: time.Time{},
	}
	if err := site.Validate(); err != nil {
		return Site{}, ErrInvalidEntity
	}
	return site, nil
}

// IsNewerLastUpdatedAt
func (s Site) IsNewerLastUpdatedAt(lastUpdatedAt time.Time) bool {
	return s.LastUpdatedAt.After(lastUpdatedAt) || s.LastUpdatedAt.Equal(lastUpdatedAt)
}

// UpdateLastUpdatedAt
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
