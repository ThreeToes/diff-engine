package main

import (
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/database"
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
	feeds, err := feedSvc.GetCurrentFeeds()
	for k,v := range feeds {
		log.Printf("%s:\n", k)
		for _, item := range v {
			log.Printf("\t%s\n",item.Title)
		}
	}
	database.ConnectToDb(conf.Database)
}

