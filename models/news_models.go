package models

import (
	"time"
)

type NewsArticle struct {
	ID uint `gorm:"primary_key; auto_increment"`
	Title string `gorm:"type:text;NOT NULL"`
	Link string `gorm:"type:text"`
	Date *time.Time `gorm:"type:timestamp;"`
	CreatedAt *time.Time `gorm:"type:timestamp;DEFAULT:NOW()"`
	Body string `gorm:"type:text"`
}

type NewsArticlePersistenceLayer interface {
	Initialise() error
	Destroy() error
	Save(article *NewsArticle) (*NewsArticle, error)
	SearchByLink(link string) (*[]*NewsArticle, error)
	Delete(article *NewsArticle) error
	GetById(id uint) (*NewsArticle, error)
}