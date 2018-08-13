package main

import (
	"github.com/mmcdole/gofeed"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"github.com/jinzhu/gorm"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
)

type CommandLineOptions struct {
	ConfigLocation *string
}

type DatabaseOptions struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	SslMode string `json:"ssl_mode"`
}

type FeedLocations map[string] string

type ConfigFileOptions struct {
	Database *DatabaseOptions `json:"database"`
	Feeds *FeedLocations `json:"feeds"`
}

type CommandLineError struct {
	MissingConfigs []string
}

func (c CommandLineError) Error() string {
	return fmt.Sprintf("Missing args: %v", c.MissingConfigs)
}

func GetConfig(path string) (*ConfigFileOptions, error){
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf ConfigFileOptions
	err = json.Unmarshal(byteValue, &conf)
	return &conf, err
}

func (c *CommandLineOptions) GetOptions() error {
	c.ConfigLocation = flag.String("config", "", "Configuration file location")
	flag.Parse()
	if *c.ConfigLocation == "" {
		return CommandLineError{
			MissingConfigs: []string{"config"},
		}
	}
	return nil
}

func main(){
	var cmdLine CommandLineOptions
	err := cmdLine.GetOptions()
	if err != nil {
		log.Println("Error parsing command line!")
		log.Println(err)
		flag.Usage()
		return
	}
	conf, err := GetConfig(*cmdLine.ConfigLocation)
	if err != nil {
		log.Println("Error getting config!")
		log.Println(err)
		return
	}
	fp := gofeed.NewParser()
	for k,v := range *conf.Feeds {

		feed, err := fp.ParseURL(v)
		if err != nil {
			log.Println("Unable to parse '%s' feed", k)
			continue
		}
		for i := 0; i < len(feed.Items); i++ {
			item := feed.Items[i]
			fmt.Println(item.Title)
		}
	}
	ConnectToDb(conf.Database)
}

func ConnectToDb(dbOptions *DatabaseOptions) {
	confString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		dbOptions.Host, dbOptions.Port, dbOptions.Username, dbOptions.Password, dbOptions.SslMode)
	db, err := gorm.Open("postgres", confString)
	if err != nil {
		fmt.Println("failed to connect")
		fmt.Println(err)
	} else {
		fmt.Println("connected ok")
	}
	db.Close()
}