package bot

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var conf *Config

var b *telebot.Bot

var db *gorm.DB

// Package init function
func init() {
	conf = initConfig()

	db = initDb()

	b = initTelegram(conf.TelegramKey)

	initMonitor()

	prices = initPrices()

	time.Sleep(time.Second * 3)

	go startRssMonitor()
}

func Start() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logs("Bot started successfully. 🚀")

	b.Start()
}
