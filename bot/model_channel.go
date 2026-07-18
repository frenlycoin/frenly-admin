package bot

import (
	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	TelegramId int64  `gorm:"uniqueIndex"`
	Type       uint8  `gorm:"default:1"`
	Name       string `gorm:"size:255"`
	Link       string `gorm:"size:255"`
}

func getChannelOrCreate(c telebot.Context) (*Channel, error) {
	ch := &Channel{}

	res := db.Where("telegram_id = ?", c.Chat().ID).Attrs(
		&Channel{}).FirstOrCreate(ch)

	if res.Error != nil {
		loge(res.Error)
		return ch, res.Error
	}

	ch.Name = c.Chat().Title
	ch.Link = c.Chat().Username
	db.Save(ch)

	return ch, nil
}

func getChannel(id int) *Channel {
	c := &Channel{}

	db.First(c, id)

	return c
}

func getChannelByLink(link string) *Channel {
	c := &Channel{}

	db.Where("link = ?", link).First(c)

	return c
}
