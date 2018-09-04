package main

import (
	"flag"
	"log"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/service"
	"github.com/threetoes/diff-engine/models"
	"sync"
	"bufio"
	"os"
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
	cntrl := make(chan service.Signal)
	dffs := make(chan *models.DiffText)
	var grp sync.WaitGroup
	grp.Add(1)
	go service.Differ(conf, dffs, cntrl, &grp)
	log.Printf("Press Enter to quit")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	log.Printf("Shutting down")
	cntrl <- service.STOP
	grp.Wait()
	log.Printf("Goodbye!")
}