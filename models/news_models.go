package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type NewsArticle struct {
	gorm.Model
	Id int64 `gorm:"type:bigint;PRIMARY_KEY;AUTO_INCREMENT"`
	Title string `gorm:"type:text;NOT NULL"`
	Link string `gorm:"type:text"`
	Date *time.Time `gorm:"type:timestamp"`
	Body string `gorm:"type:text"`
}