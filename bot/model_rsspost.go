package bot

import (
	"gorm.io/gorm"
)

type RssPost struct {
	gorm.Model
	Source  string `gorm:"uniqueIndex:idx_rss_source_guid"`
	Guid    string `gorm:"uniqueIndex:idx_rss_source_guid"`
	Title   string `gorm:"type:text"`
	Content string `gorm:"type:text"`
	Link    string `gorm:"type:text"`
}
