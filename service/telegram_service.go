package service

import (
	"fmt"
	"github.com/threetoes/diff-engine/config"
	"github.com/threetoes/diff-engine/models"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type TelegramApi interface {
	SendDiff(text *models.DiffText) error
	TelegramUpdateRunLoop()
	DiffSenderLoop()
}

type TelegramService struct {
	TelegramApi
	botapi *tgbotapi.BotAPI
	userIds []int64
	diffs chan *models.DiffText
}

func (svc TelegramService) SendDiff(text *models.DiffText) error{
	sendText := fmt.Sprintf("WOOWOOWOO! WE HAVE AN EDIT!\n\n%s\n\n%s", text.Url, text.DiffText)
	for _, usr := range svc.userIds {
		svc.botapi.Send(tgbotapi.NewMessage(usr, sendText))
	}
	return nil
}

func (svc TelegramService)  TelegramUpdateRunLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := svc.botapi.GetUpdatesChan(u)
	if err != nil {
		log.Println("Telegram init failed!")
		return
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "register" {
			svc.userIds = append(svc.userIds, update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Registered you, homeslice")
			svc.botapi.Send(msg)
		}
	}
}

func (svc TelegramService) DiffSenderLoop(){
	for elem := range svc.diffs {
		svc.SendDiff(elem)
	}
}

func NewTelegramBot(conf *config.ConfigFileOptions, diffChannel chan *models.DiffText) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(conf.ServiceSettings.TelegramToken)
	if err != nil {
		return nil, err
	}
	svc := TelegramService{
		botapi: bot,
		diffs:diffChannel,
	}
	return &svc, nil
}