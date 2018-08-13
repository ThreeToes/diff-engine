package test

import (
	"testing"
	"github.com/threetoes/diff-engine/service"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/models"
	"time"
)

func getConfig() *config.DatabaseOptions {
	return &config.DatabaseOptions{
		Name: "diff-engine-test",
		Host: "localhost",
		Port: 5432,
		Password: "diff-engine",
		SslMode: "disable",
		Username: "diff-engine",
	}
}

func setup() (*service.NewsArticleService, error) {
	svc, err := service.NewArticleService(getConfig())
	if err != nil {
		return nil, err
	}
	err = svc.Initialise()
	return svc, err
}

func destroy(svc *service.NewsArticleService) {
	svc.DropAll()
	svc.Destroy()
}

func TestInitialise(t *testing.T) {
	svc,err := setup()
	if err != nil {
		t.Log(err)
		destroy(svc)
		t.Fail()
		return
	}
	destroy(svc)
}

func TestCreate(t *testing.T) {
	svc, err := setup()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	now := time.Now()
	article := models.NewsArticle{
		Body: "herpderp",
		Link: "http://test.com/",
		Date: &now,
		Title: "Great title, well done",
	}
	returned, err := svc.Save(&article)
	if err != nil {
		t.Log(err)
		destroy(svc)
		t.Fail()
		return
	}
	if returned.ID == 0 {
		t.Log("ID was not set! Probably failed to create")
		destroy(svc)
		t.Fail()
		return
	}
	destroy(svc)
}

func TestGetById(t *testing.T) {
	svc, err := setup()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	now := time.Now()
	article := models.NewsArticle{
		Body: "herpderp",
		Link: "http://test.com/",
		Date: &now,
		Title: "Great title, well done",
	}
	returned, err := svc.Save(&article)
	if err != nil {
		Fail(t, err, svc)
		return
	}
	if returned.ID == 0 {
		Fail(t, "ID not set! Insert probably failed", svc)
		return
	}
	got, err := svc.GetById(returned.ID)
	if err != nil {
		Fail(t, err, svc)
		return
	}
	if got.Title != article.Title {
		Fail(t, "Titles didn't match!", svc)
		return
	}
	if got.ID != returned.ID {
		Fail(t, "IDs didn't match!", svc)
		return
	}
	if got.Link != article.Link {
		Fail(t, "Links didn't match!", svc)
		return
	}
	destroy(svc)
}

func TestSearchByLink(t *testing.T){
	svc, err := setup()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	now := time.Now()
	article := models.NewsArticle{
		Body: "herpderp",
		Link: "http://test.com/",
		Date: &now,
		Title: "Great title, well done",
	}
	returned, err := svc.Save(&article)
	if returned.ID == 0 || err != nil {
		t.Log("Failed to save!")
		Fail(t, err, svc)
		return
	}
	article2 := models.NewsArticle{
		Body: "herpderp",
		Link: "http://test.com/",
		Date: &now,
		Title: "Great title, well done",
	}
	returned, err = svc.Save(&article2)
	if returned.ID == 0 || err != nil {
		t.Log("Failed to save!")
		Fail(t, err, svc)
		return
	}

	article3 := models.NewsArticle{
		Body: "herpderp",
		Link: "http://test2.com/",
		Date: &now,
		Title: "Great title, well done",
	}
	returned, err = svc.Save(&article3)
	if returned.ID == 0 || err != nil {
		t.Log("Failed to save!")
		Fail(t, err, svc)
		return
	}

	articles, err := svc.SearchByLink(article.Link)
	if err != nil {
		Fail(t, err, svc)
		return
	}
	if len(*articles) != 2 {
		Fail(t, "Incorrect number of returned articles!", svc)
		return
	}
	for i := 0; i < len(*articles); i++ {
		a := (*articles)[i]
		if a.Link != article.Link {
			Fail(t, "Picked up an incorrect link!", svc)
			return
		}
	}
	destroy(svc)
}

func TestDelete(t *testing.T) {

}

func Fail(t *testing.T, msg interface{}, svc *service.NewsArticleService) {
	destroy(svc)
	t.Log(msg)
	t.Fail()
}
