package service

import (
	"github.com/mmcdole/gofeed"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/models"
)

type FeedMap map[string] []*models.NewsArticle

type FeedService interface {
	GetCurrentFeeds() (FeedMap, error)
}

type GoFeedService struct {
	FeedService
	parser *gofeed.Parser
	feeds *config.FeedLocations
}

func (svc *GoFeedService) GetCurrentFeeds() (FeedMap, error) {
	fm := FeedMap{}
	for k,v := range *svc.feeds {
		feed, _ := svc.parser.ParseURL(v)
		transformed := make([]*models.NewsArticle,0)
		for _,item := range feed.Items {
			transformed = append(transformed, transform(item))
		}
		fm[k] = transformed
	}
	return fm, nil
}

func NewFeedService(locations *config.FeedLocations) *GoFeedService {
	feedService := GoFeedService{}
	feedService.parser = gofeed.NewParser()
	feedService.feeds = locations
	return &feedService
}

func transform(item *gofeed.Item) *models.NewsArticle {
	var a models.NewsArticle
	a.Link = item.Link
	a.Title = item.Title
	a.Author = item.Author.Name
	a.Date = item.PublishedParsed
	a.Body = item.Content
	return &a
}
