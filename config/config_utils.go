package config

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type CommandLineOptions struct {
	ConfigLocation *string
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

type CommandLineError struct {
	MissingConfigs []string
}

func (c CommandLineError) Error() string {
	return fmt.Sprintf("Missing args: %v", c.MissingConfigs)
}

type DatabaseOptions struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	SslMode  string `json:"ssl_mode"`
}

type FeedLocations map[string] *FeedEntry

type FeedEntry struct {
	Url					string `json:"url"`
	ArticleSelector 	string `json:"article_selector"`
}

type ConfigFileOptions struct {
	Database *DatabaseOptions `json:"database"`
	Feeds    *FeedLocations   `json:"feeds"`
}

func GetConfig(path string) (*ConfigFileOptions, error) {
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

