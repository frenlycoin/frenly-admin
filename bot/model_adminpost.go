package bot

import (
	"time"

	"gorm.io/gorm"
)

type AdminPost struct {
	gorm.Model
	Channel       int64     `gorm:"index"`
	Text          string    `gorm:"type:text"`
	Published     bool      `gorm:"default:false"`
	TimePublished time.Time `gorm:"index"`
}
