package main

import (
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/service"
	"github.com/PuerkitoBio/goquery"
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

type Signal int

const STOP = 1

func main(){
	var cmdLine config.CommandLineOptions
	err := cmdLine.GetOptions()
	if err != nil {
		log.Println("Error parsing command line!")
		log.Println(err)
		flag.Usage()
		return
	}
	conf, err := config.GetConfig(*cmdLine.ConfigLocation)
	if err != nil {
		log.Println("Error getting config!")
		log.Println(err)
		return
	}
	feedSvc := service.NewFeedService(conf.Feeds)
	articleSvc,err := service.NewArticleService(conf.Database)
	articleSvc.Initialise()
	diffSvc := service.NewDiff(feedSvc, articleSvc)
	diffSvc.FetchAndSaveDiffs()
	articles, err := articleSvc.GetWatchList()
	for _, a := range *articles {
		s, _ := CheckDiff(a.Link, a.Body)
		for _, d := range s {
			if d.Type != diffmatchpatch.DiffEqual {
				log.Printf("REE ROO STEALTH EDIT!\n")
				log.Printf("%s", diffmatchpatch.New().DiffText1(s))
				break
			}
		}
	}
	articleSvc.Destroy()
}

func FeedDiffs(articleSvc *service.NewsArticleService, signals chan Signal){

}

func CheckDiff(url string, body string) ([]diffmatchpatch.Diff, error){
	doc, err := goquery.NewDocument(url)
	doc2, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	content := doc.Find(".article__body").Text()
	content2 := doc2.Text()
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(content, content2, false)
	return diffs, err
}
