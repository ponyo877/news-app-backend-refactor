package presenter

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type SiteResponce struct {
	Data []*Site `json:"data"`
}

type Site struct {
	ID            string    `json:"id"`
	Title         string    `json:"titles"`
	RSSURL        string    `json:"url"`
	ImageURL      string    `json:"image"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type SiteList []*Site

func pickSite(site entity.Site) (Site, error) {
	sitePresenter := Site{
		ID:            site.ID.String(),
		Title:         site.Title,
		RSSURL:        site.RSSURL,
		ImageURL:      site.ImageURL,
		LastUpdatedAt: site.LastUpdatedAt,
	}
	return sitePresenter, nil
}

func PickSiteList(siteList []entity.Site) (SiteList, error) {
	var sitePresenterList SiteList
	for _, article := range siteList {
		var sitePresenter Site
		sitePresenter, err := pickSite(article)
		if err != nil {
			return SiteList{}, err
		}
		sitePresenterList = append(sitePresenterList, &sitePresenter)
	}
	return sitePresenterList, nil
}
