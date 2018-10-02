package service

import (
	"github.com/threetoes/diff-engine/config"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"github.com/threetoes/diff-engine/models"
	"sync"
	"time"
)

func Differ(conf *config.ConfigFileOptions, articleChannel chan *models.DiffText, controlChannel chan Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	feedSvc := NewFeedService(conf.Feeds)
	articleSvc, err := NewArticleService(conf.Database)
	if err != nil {
		log.Printf("Article service failed to initialise!")
		return
	}
	articleSvc.Initialise()
	defer articleSvc.Destroy()
	diffSvc := NewDiff(feedSvc, articleSvc)
	diffSvc.FetchAndSaveFeed()
	for ;; {
		changes, err := diffSvc.CheckForDiffs()
		diffLib := diffmatchpatch.New()
		if err != nil {
			log.Printf("%s\n", err)
		} else {
			for _, change := range *changes {
				log.Printf("STEALTH EDIT!\n")
				diffs := diffLib.DiffMain(change.Original.Body, change.Changed.Body, false)
				diff := &models.DiffText{
					OriginalText: change.Original.Body,
					NewText:      change.Changed.Body,
					DiffText:     diffLib.DiffPrettyText(diffs),
					Url:		  change.Original.Source,
				}
				articleChannel <- diff
			}
		}
		for i := 0; i < conf.ServiceSettings.ArticleGrabCooldown; i++ {
			select {
			case msg := <-controlChannel:
				if msg == STOP {
					log.Printf("Stopping diff routine")
					return
				}
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}
}
