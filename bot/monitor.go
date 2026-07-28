package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Monitor struct {
}

func (m *Monitor) start() {
	m.generateAdminPostsIfNeeded()
	m.publishPendingTodayPostsOnStartup()

	for {
		m.generateAdminPostsIfNeeded()
		time.Sleep(time.Second * MonitorTick)
	}
}

func (m *Monitor) generateAdminPostsIfNeeded() {
	m.generateDevPostsIfNeeded()
	m.generateTonPostsIfNeeded()
}

func (m *Monitor) generateDevPostsIfNeeded() {
	lastDay, err := getLastPostDay("lastFrenlyDevPost")
	if err != nil {
		logs(fmt.Sprintf("failed to load last FrenlyDev post day: %v", err))
		return
	}

	today := time.Now().UTC().Format("2006-01-02")
	if lastDay == today {
		return
	}

	advice, err := getProgrammingAdvice()
	if err != nil {
		logs(fmt.Sprintf("failed to generate programming advice: %v", err))
		return
	}

	posts := strings.Split(advice, "\n\n")
	for _, post := range posts {
		trimmed := strings.TrimSpace(post)
		if trimmed == "" {
			continue
		}
		if err := db.Create(&AdminPost{Channel: FrenlyDevs, Text: trimmed}).Error; err != nil {
			logs(fmt.Sprintf("failed to save admin post: %v", err))
		}
	}

	if err := saveLastPostDay("lastFrenlyDevPost", today); err != nil {
		logs(fmt.Sprintf("failed to save last FrenlyDev post day: %v", err))
	}
}

func (m *Monitor) generateTonPostsIfNeeded() {
	lastDay, err := getLastPostDay("lastFrenlyTonPost")
	if err != nil {
		logs(fmt.Sprintf("failed to load last FrenlyTon post day: %v", err))
		return
	}

	today := time.Now().UTC().Format("2006-01-02")
	if lastDay == today {
		return
	}

	news, err := getTonNews()
	if err != nil {
		logs(fmt.Sprintf("failed to generate TON news: %v", err))
		return
	}

	posts := strings.Split(news, "\n\n")
	for _, post := range posts {
		trimmed := strings.TrimSpace(post)
		if trimmed == "" {
			continue
		}
		if err := db.Create(&AdminPost{Channel: FrenlyTon, Text: trimmed}).Error; err != nil {
			logs(fmt.Sprintf("failed to save TON admin post: %v", err))
		}
	}

	if err := saveLastPostDay("lastFrenlyTonPost", today); err != nil {
		logs(fmt.Sprintf("failed to save last FrenlyTon post day: %v", err))
	}
}

func getLastPostDay(key string) (string, error) {
	var kv KeyValue
	res := db.Where("key = ?", key).First(&kv)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", res.Error
	}
	return kv.ValueStr, nil
}

func saveLastPostDay(key string, day string) error {
	var kv KeyValue
	res := db.Where("key = ?", key).First(&kv)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
		kv = KeyValue{Key: key}
	}

	kv.ValueStr = day
	return db.Save(&kv).Error
}

func getLastFrenlyDevPostDay() (string, error) {
	var kv KeyValue
	res := db.Where("key = ?", "lastFrenlyDevPost").First(&kv)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", res.Error
	}
	return kv.ValueStr, nil
}

func saveLastFrenlyDevPostDay(day string) error {
	var kv KeyValue
	res := db.Where("key = ?", "lastFrenlyDevPost").First(&kv)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
		kv = KeyValue{Key: "lastFrenlyDevPost"}
	}

	kv.ValueStr = day
	return db.Save(&kv).Error
}

func (m *Monitor) publishPendingTodayPostsOnStartup() {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	channels := []int64{FrenlyDevs, FrenlyTon}
	for _, ch := range channels {
		var post AdminPost
		res := db.Where("channel = ? AND published = ? AND created_at >= ? AND created_at < ?", ch, false, dayStart, dayEnd).First(&post)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				continue
			}
			logs(fmt.Sprintf("failed to load pending admin post for channel %d: %v", ch, res.Error))
			continue
		}

		if err := m.publishAdminPost(post, now); err != nil {
			logs(fmt.Sprintf("failed to publish startup admin post for channel %d: %v", ch, err))
		}
	}
}

func (m *Monitor) publishAdminPostsIfNeeded() {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	channels := []int64{FrenlyDevs, FrenlyTon}
	for _, ch := range channels {
		var posts []AdminPost
		if err := db.Where("channel = ? AND published = ? AND created_at >= ? AND created_at < ?", ch, false, dayStart, dayEnd).Find(&posts).Error; err != nil {
			logs(fmt.Sprintf("failed to load pending admin posts for channel %d: %v", ch, err))
			continue
		}

		if len(posts) == 0 {
			continue
		}

		rand.Seed(now.UnixNano())
		chosen := posts[rand.Intn(len(posts))]
		if err := m.publishAdminPost(chosen, now); err != nil {
			logs(fmt.Sprintf("failed to publish admin post for channel %d: %v", ch, err))
		}
	}
}

func (m *Monitor) publishAdminPost(post AdminPost, now time.Time) error {
	rec := &telebot.Chat{ID: post.Channel}
	if _, err := b.Send(rec, post.Text, telebot.NoPreview); err != nil {
		return err
	}

	post.Published = true
	post.TimePublished = now
	if err := db.Save(&post).Error; err != nil {
		return err
	}

	return nil
}

func (m *Monitor) publishAdminPostsRoutine() {
	rand.Seed(time.Now().UnixNano())
	delay := (time.Hour * 2) + time.Duration(rand.Int63n(int64(time.Hour)))
	time.Sleep(delay)

	for {
		m.publishAdminPostsIfNeeded()
		delay = (time.Hour * 2) + time.Duration(rand.Int63n(int64(time.Hour)))
		time.Sleep(delay)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	go m.publishAdminPostsRoutine()
	return m
}
