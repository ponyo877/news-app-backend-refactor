package entyty

import "time"

type Article struct {
	ImageURL  string
	SiteID    ID
	Sitetitle string
	Title     string
	URL       string
	UpdatedAt time.Time
	CreatedAt time.Time
}
