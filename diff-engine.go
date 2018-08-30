package main

import (
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/service"
	"github.com/sergi/go-diff/diffmatchpatch"
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
	diffSvc.FetchAndSaveFeed()
	changes, err := diffSvc.CheckForDiffs()
	diffLib := diffmatchpatch.New()
	if err != nil {
		log.Printf("%s\n",err)
	} else {
		for _, change := range *changes {
			log.Printf("REEROO! STEALTH EDIT!\n")
			diffs := diffLib.DiffMain(change.Original.Body, change.Changed.Body, false)
			for _, d := range diffs {
				log.Printf("%s \n", d.Text)
			}
		}
	}
	articleSvc.Destroy()
}