package service

import (
	"github.com/jinzhu/gorm"
	"github.com/threetoes/diff-engine/models"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/database"
	"log"
)

type NewsArticleService struct {
	models.NewsArticlePersistenceLayer
	database *gorm.DB
}

func NewArticleService(conf *config.DatabaseOptions) (*NewsArticleService, error) {
	var svc NewsArticleService
	var err error
	svc.database, err = database.ConnectToDb(conf)
	if err != nil {
		log.Printf("failed to connect to db!\n")
		log.Println(err)
		return nil, err
	}
	return &svc, nil
}

func (svc *NewsArticleService) Initialise() error{
	err := svc.database.AutoMigrate(&models.NewsArticle{}).Error
	if err != nil {
		log.Println("Could not migrate database!")
		log.Fatal(err)
	}
	return err
}

func (svc *NewsArticleService) DropAll() error {
	return svc.database.DropTableIfExists(&models.NewsArticle{}).Error
}

func (svc *NewsArticleService) Destroy() error {
	return svc.database.Close()
}

func (svc *NewsArticleService) Save(article *models.NewsArticle) (*models.NewsArticle, error) {
	var returned models.NewsArticle
	returned.ID = 0
	q := svc.database.FirstOrInit(&returned, article)
	if returned.ID == 0 {
		q = q.Create(&returned)
	} else {
		q = q.Save(&article)
	}
	return &returned, q.Error
}

func (svc *NewsArticleService) SearchByLink(link string) (*[]*models.NewsArticle, error) {
	var newsArticles [] *models.NewsArticle
	q := svc.database.Order("id").Where("link = ?", link).Find(&newsArticles)
	return &newsArticles, q.Error
}

func (svc *NewsArticleService) Delete(article *models.NewsArticle) error {
	return svc.database.Delete(article,"id = ?", article.ID).Error
}

func (svc *NewsArticleService) GetById(id uint) (*models.NewsArticle, error) {
	var article models.NewsArticle
	q := svc.database.First(&article)
	return &article, q.Error
}

func (svc *NewsArticleService) GetWatchList() (*[]models.NewsArticle, error){
	var articles []models.NewsArticle
	db := svc.database.Where("date > (now() - interval '1 week')").
		Order("id desc").
		Select("DISTINCT link, id, title, link, date, created_at, body, author").
		Find(&articles)
	if db.Error != nil {
		return nil, db.Error
	}
	return &articles, nil
}

func (svc *NewsArticleService) GetWatchListBySource(source string) (*[]models.NewsArticle, error){
	var articles []models.NewsArticle
	db := svc.database.Where("date > (now() - interval '1 week') AND source = ?", source).
		Order("id desc").
		Select("DISTINCT link, id, title, link, date, created_at, body, author").
		Find(&articles)
	if db.Error != nil {
		return nil, db.Error
	}
	return &articles, nil
}