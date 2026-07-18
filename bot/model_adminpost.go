package bot

import (
	"time"

	"gorm.io/gorm"
)

type AdminPost struct {
	gorm.Model
	Text          string    `gorm:"type:text"`
	Published     bool      `gorm:"default:false"`
	TimePublished time.Time `gorm:"index"`
}
