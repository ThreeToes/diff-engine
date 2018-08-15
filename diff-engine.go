package main

import (
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/service"
)

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
	articleSvc.Destroy()
}

