package service

import (
	"log"
	//"github.com/threetoes/diff-engine/models"
	"github.com/threetoes/diff-engine/models"
	"github.com/threetoes/diff-engine/config"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ChangedArticle struct {
	Original *models.NewsArticle
	Changed *models.NewsArticle
}

type DiffService interface {
	// Fetch the RSS articles and save them if we don't already have a copy
	FetchAndSaveFeed() ([]*models.NewsArticle, error)
	// Check for any changes
	CheckForDiffs() (*[]*ChangedArticle, error)
}

type DiffServiceImpl struct {
	DiffService
	feedService FeedService
	newsArticleService models.NewsArticlePersistenceLayer
	config config.FeedLocations
}

func (svc *DiffServiceImpl) FetchAndSaveFeed() ([]*models.NewsArticle, error) {
	feeds, err := svc.feedService.GetCurrentFeeds()
	if err != nil {
		return nil, err
	}
	articles := make([]*models.NewsArticle, 0)
	for site, v := range feeds {
		log.Printf("Checking %s...\n", site)
		for _, item := range v {
			loaded, _ := svc.newsArticleService.SearchByLink(item.Link)
			if len(*loaded) == 0 {
				log.Printf("\tSaving article '%s' from %s\n", item.Title, site)
				svc.newsArticleService.Save(item)
				articles = append(articles, item)
			}
		}
	}
	return articles, nil
}

func (svc *DiffServiceImpl) CheckForDiffs() (*[]*ChangedArticle, error){
	diffs := make([]*ChangedArticle, 0)
	for feedName, feed := range svc.config {
		watchList, err := svc.newsArticleService.GetWatchListBySource(feedName)
		if err != nil {
			return nil, err
		}
		for _, article := range *watchList {
			doc, err := goquery.NewDocument(article.Link)
			if err != nil {
				return nil, err
			}
			oldDoc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Body))
			if err != nil {
				return nil, err
			}
			oldArticleText := oldDoc.Text()
			newArticleText := doc.Find(feed.ArticleSelector).Text()
			newTitle := doc.Find(feed.TitleSelector).Text()
			if article.Title != newTitle || oldArticleText != newArticleText {
				diffs = append(diffs, &ChangedArticle{
					Original: &article,
					Changed: &models.NewsArticle{
						Title:newTitle,
						Body: newArticleText,
						Link: article.Link,
						Source: feedName,
						Author: article.Author,
						Date: article.Date,
					},
				})
			}
		}
	}
	return &diffs, nil
}

func NewDiff(feed *GoFeedService, layer *NewsArticleService) *DiffServiceImpl {
	diffService := DiffServiceImpl{
		newsArticleService:layer,
		feedService:feed,
	}
	return &diffService
}