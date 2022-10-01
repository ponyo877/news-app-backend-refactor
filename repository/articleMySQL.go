package repository

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleRepositoryPresenter struct {
	ID          string    `gorm:"column:id;primary_key"`
	Title       string    `gorm:"column:title"`
	URL         string    `gorm:"column:url"`
	ImageURL    string    `gorm:"column:image_url"`
	SiteID      string    `gorm:"column:site_id"`
	PublishedAt time.Time `gorm:"column:published_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

type SiteArticleRepositoryPresenter struct {
	ArticleRepositoryPresenter
	SiteMySQLPresenter
}

type SiteArticleRepositoryPresenterList []SiteArticleRepositoryPresenter

func (a ArticleRepositoryPresenter) TableName() string {
	return "articles"
}

// pickArticle
func (p *SiteArticleRepositoryPresenter) pickArticle() (entity.Article, error) {
	siteID, err := entity.StringToID(p.SiteMySQLPresenter.ID)
	if err != nil {
		return entity.Article{}, err
	}
	site := entity.Site{
		ID:            siteID,
		Title:         p.SiteMySQLPresenter.Title,
		RSSURL:        p.SiteMySQLPresenter.RSSURL,
		LastUpdatedAt: p.SiteMySQLPresenter.LastUpdatedAt,
	}
	articleID, err := entity.StringToID(p.ArticleRepositoryPresenter.ID)
	if err != nil {
		return entity.Article{}, err
	}
	articleTitle := entity.NewArticleTitle(p.ArticleRepositoryPresenter.Title)
	imageURL, err := entity.NewImageURL(p.ArticleRepositoryPresenter.ImageURL)
	if err != nil {
		return entity.Article{}, err
	}
	article := entity.Article{
		ID:          articleID,
		Title:       articleTitle,
		URL:         p.ArticleRepositoryPresenter.URL,
		ImageURL:    imageURL,
		Site:        site,
		PublishedAt: p.ArticleRepositoryPresenter.PublishedAt,
		UpdatedAt:   p.ArticleRepositoryPresenter.UpdatedAt,
		CreatedAt:   p.ArticleRepositoryPresenter.CreatedAt,
	}
	return article, nil
}

// pickArticleList
func (s *SiteArticleRepositoryPresenterList) pickArticleList() ([]entity.Article, error) {
	var articleList []entity.Article
	for _, siteArticleRepositoryPresenter := range *s {
		article, err := siteArticleRepositoryPresenter.pickArticle()
		if err != nil {
			return nil, err
		}
		articleList = append(articleList, article)
	}
	return articleList, nil
}

// Get
func (r *ArticleRepository) Get(ID entity.ID) (entity.Article, error) {
	var siteArticleRepositoryPresenter SiteArticleRepositoryPresenter
	if err := r.db.
		Model(&ArticleRepositoryPresenter{}).
		Select("articles.*, sites.*").
		Joins("LEFT JOIN sites on sites.id = articles.site_id").
		Where("articles.id = ?", ID.String()).
		Find(&siteArticleRepositoryPresenter).Error; err != nil {
		log.Infof("DBの接続に失敗しました: %v", err)
	}
	article, err := siteArticleRepositoryPresenter.pickArticle()
	if err == gorm.ErrRecordNotFound {
		return entity.Article{}, entity.ErrNotFound
	}
	if err != nil {
		log.Infof("pickArticleに失敗しました: %v", err)
		return entity.Article{}, err
	}
	return article, nil
}

// ListOption
func (r *ArticleRepository) ListOption(basePublishedAt time.Time, invisibleIDSet entity.IDSet) ([]entity.Article, error) {
	var siteArticleRepositoryPresenterList SiteArticleRepositoryPresenterList
	if err := r.db.
		Model(&ArticleRepositoryPresenter{}).
		Select("articles.*, sites.*").
		Joins("LEFT JOIN sites on sites.id = articles.site_id").
		Where("sites.id NOT IN ?", invisibleIDSet.Strings()).
		Where("articles.published_at < ?", basePublishedAt).
		Find(&siteArticleRepositoryPresenterList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Infof("DBの接続に失敗しました: %v", err)
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return []entity.Article{}, nil
	}
	articleList, err := siteArticleRepositoryPresenterList.pickArticleList()
	if err != nil {
		log.Infof("pickArticleListに失敗しました: %v", err)
		return nil, err
	}
	return articleList, nil
}

// Create
func (r *ArticleRepository) Create(e entity.Article) (entity.ID, time.Time, error) {
	imageURL, err := e.ImageURL.URL()
	if err != nil {
		return entity.NewID(), time.Now(), err
	}
	ArticleRepositoryPresenter := &ArticleRepositoryPresenter{
		ID:          e.ID.String(),
		Title:       e.Title.String(),
		URL:         e.URL,
		ImageURL:    imageURL,
		SiteID:      e.Site.ID.String(),
		PublishedAt: e.PublishedAt,
		UpdatedAt:   e.UpdatedAt,
		CreatedAt:   e.CreatedAt,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "site_id"}, {Name: "title"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "created_at"}),
	}).Create(&ArticleRepositoryPresenter).Error; err != nil {
		return entity.NewID(), time.Now(), err
	}
	return e.ID, e.CreatedAt, nil
}

// Update
func (r *ArticleRepository) Update(e entity.Article) error {
	return nil
}

// DeleteByID
func (r *ArticleRepository) DeleteByID(id entity.ID) error {
	return nil
}

// List
func (r *ArticleRepository) List(IDList []entity.ID) ([]entity.Article, error) {
	var articleList []entity.Article
	for _, ID := range IDList {
		article, err := r.Get(ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Infof("Getが失敗しました: %v", err)
			return nil, err
		} else if err == gorm.ErrRecordNotFound {
			continue
		}
		articleList = append(articleList, article)
	}
	if len(articleList) == 0 {
		return []entity.Article{}, nil
	}
	return articleList, nil
}
