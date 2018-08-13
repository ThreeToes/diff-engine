package database

import (
	"github.com/threetoes/diff-engine/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectToDb(dbOptions *config.DatabaseOptions) (*gorm.DB, error) {
	confString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		dbOptions.Host, dbOptions.Port, dbOptions.Username, dbOptions.Password, dbOptions.SslMode)
	return gorm.Open("postgres", confString)
}