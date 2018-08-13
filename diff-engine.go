package main

import (
	"github.com/mmcdole/gofeed"
	"fmt"
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/database"
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
	fp := gofeed.NewParser()
	for k,v := range *conf.Feeds {

		feed, err := fp.ParseURL(v)
		if err != nil {
			log.Printf("Unable to parse '%s' feed\n", k)
			continue
		}
		for i := 0; i < len(feed.Items); i++ {
			item := feed.Items[i]
			fmt.Println(item.Title)
		}
	}
	database.ConnectToDb(conf.Database)
}

