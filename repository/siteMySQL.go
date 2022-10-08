package repository

import (
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
)

// SiteMySQL mysql repository
type SiteMySQL struct {
	db *gorm.DB
}

type SiteMySQLPresenter struct {
	ID            string    `gorm:"column:id;primary_key"`
	Title         string    `gorm:"column:title"`
	RSSURL        string    `gorm:"column:rss_url"`
	ImageURL      string    `gorm:"column:image_url"`
	LastUpdatedAt time.Time `gorm:"column:last_updated_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	CreatedA      time.Time `gorm:"column:created_at"`
}

func (s SiteMySQLPresenter) TableName() string {
	return "sites"
}

// NewSiteMySQL create new repository
func NewSiteMySQL(db *gorm.DB) *SiteMySQL {
	return &SiteMySQL{
		db: db,
	}
}

// Get
func (r *SiteMySQL) Get(id entity.ID) (entity.Site, error) {
	return entity.Site{}, nil
}

// Search
func (r *SiteMySQL) Search(query string) ([]entity.Site, error) {
	return []entity.Site{}, nil
}

// List
func (r *SiteMySQL) List() ([]entity.Site, error) {
	var siteMySQLList []SiteMySQLPresenter
	if err := r.db.Find(&siteMySQLList).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Infof("DBの接続に失敗しました: %v", err)
		return nil, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return []entity.Site{}, nil
	}
	var siteList []entity.Site
	for _, siteMySQL := range siteMySQLList {
		id, err := entity.StringToID(siteMySQL.ID)
		if err != nil {
			return nil, err
		}
		site := entity.Site{
			ID:            id,
			Title:         siteMySQL.Title,
			RSSURL:        siteMySQL.RSSURL,
			ImageURL:      siteMySQL.ImageURL,
			LastUpdatedAt: siteMySQL.LastUpdatedAt,
		}
		siteList = append(siteList, site)
	}
	return siteList, nil
}

// Create
func (r *SiteMySQL) Create(e entity.Site) (entity.ID, error) {
	return e.ID, nil
}

// Update
func (r *SiteMySQL) Update(e entity.Site) error {
	var siteMySQL SiteMySQLPresenter
	if err := r.db.Model(&siteMySQL).Where("id = ?", e.ID.String()).Update("last_updated_at", e.LastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

// Delete
func (r *SiteMySQL) Delete(id entity.ID) error {
	return nil
}
