package service

import (
	"log"
	//"github.com/threetoes/diff-engine/models"
	"github.com/threetoes/diff-engine/models"
)

type DiffService interface {
	FetchAndSaveDiffs() ([]*models.NewsArticle, error)
}

type DiffServiceImpl struct {
	DiffService
	feedService FeedService
	newsArticleService models.NewsArticlePersistenceLayer
}

func (svc *DiffServiceImpl) FetchAndSaveDiffs() ([]*models.NewsArticle, error) {
	feeds, err := svc.feedService.GetCurrentFeeds()
	if err != nil {
		return nil, err
	}
	diffs := make([]*models.NewsArticle, 0)
	for site, v := range feeds {
		log.Printf("Checking %s...\n", site)
		for _, item := range v {
			loaded, _ := svc.newsArticleService.SearchByLink(item.Link)
			if len(*loaded) == 0 {
				log.Printf("\tSaving article '%s' from %s\n", item.Title, site)
				svc.newsArticleService.Save(item)
				continue
			}
			changed := true
			for _, la := range *loaded {
				if item.Title == la.Title  && item.Body == la.Body {
					changed = false
				}
			}
			if changed {
				diffs = append(diffs, item)
				log.Printf("\tArticle '%s' from %s has changed, saving a copy", item.Title, site)
				svc.newsArticleService.Save(item)
			}
		}
	}
	return diffs, nil
}

func NewDiff(feed *GoFeedService, layer *NewsArticleService) *DiffServiceImpl {
	diffService := DiffServiceImpl{
		newsArticleService:layer,
		feedService:feed,
	}
	return &diffService
}