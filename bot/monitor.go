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
	m.generateLifePostsIfNeeded()
	m.generateHoroscopePostsIfNeeded()
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

func (m *Monitor) generateLifePostsIfNeeded() {
	lastDay, err := getLastPostDay("lastFrenlyLifePost")
	if err != nil {
		logs(fmt.Sprintf("failed to load last FrenlyLife post day: %v", err))
		return
	}

	today := time.Now().UTC().Format("2006-01-02")
	if lastDay == today {
		return
	}

	advice, err := getLifeAdvice()
	if err != nil {
		logs(fmt.Sprintf("failed to generate life advice: %v", err))
		return
	}

	posts := strings.Split(advice, "\n\n")
	for _, post := range posts {
		trimmed := strings.TrimSpace(post)
		if trimmed == "" {
			continue
		}
		if err := db.Create(&AdminPost{Channel: FrenlyLife, Text: trimmed}).Error; err != nil {
			logs(fmt.Sprintf("failed to save life admin post: %v", err))
		}
	}

	if err := saveLastPostDay("lastFrenlyLifePost", today); err != nil {
		logs(fmt.Sprintf("failed to save last FrenlyLife post day: %v", err))
	}
}

func (m *Monitor) generateHoroscopePostsIfNeeded() {
	lastDay, err := getLastPostDay("lastFrenlyAstroPost")
	if err != nil {
		logs(fmt.Sprintf("failed to load last FrenlyAstro post day: %v", err))
		return
	}

	today := time.Now().UTC().Format("2006-01-02")
	if lastDay == today {
		return
	}

	horoscopes, err := getDailyHoroscope()
	if err != nil {
		logs(fmt.Sprintf("failed to generate daily horoscope: %v", err))
		return
	}

	posts := strings.Split(horoscopes, "\n\n")
	for _, post := range posts {
		trimmed := strings.TrimSpace(post)
		if trimmed == "" {
			continue
		}
		if err := db.Create(&AdminPost{Channel: FrenlyAstro, Text: trimmed}).Error; err != nil {
			logs(fmt.Sprintf("failed to save horoscope admin post: %v", err))
		}
	}

	if err := saveLastPostDay("lastFrenlyAstroPost", today); err != nil {
		logs(fmt.Sprintf("failed to save last FrenlyAstro post day: %v", err))
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

	channels := []int64{FrenlyDevs, FrenlyTon, FrenlyLife, FrenlyAstro}
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

	channels := []int64{FrenlyDevs, FrenlyTon, FrenlyLife}
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

func (m *Monitor) publishHoroscopeRoutine() {
	// Sleep until the next 2-hour mark at a random minute between 30-59
	now := time.Now()
	nextHour := now.Hour()
	nextHour = nextHour - (nextHour % 2) + 2
	rand.Seed(now.UnixNano())
	randomMinute := 30 + rand.Intn(30) // 30-59
	next := time.Date(now.Year(), now.Month(), now.Day(), nextHour, randomMinute, 0, 0, time.UTC)
	if !now.Before(next) {
		next = next.Add(2 * time.Hour)
	}
	time.Sleep(time.Until(next))

	for {
		m.publishHoroscopeIfNeeded()
		rand.Seed(time.Now().UnixNano())
		randomMinute = 30 + rand.Intn(30)
		next = next.Add(2 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), randomMinute, 0, 0, time.UTC)
		time.Sleep(time.Until(next))
	}
}

func (m *Monitor) publishHoroscopeIfNeeded() {
	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	var posts []AdminPost
	if err := db.Where("channel = ? AND published = ? AND created_at >= ? AND created_at < ?", FrenlyAstro, false, dayStart, dayEnd).Find(&posts).Error; err != nil {
		logs(fmt.Sprintf("failed to load pending horoscope posts: %v", err))
		return
	}

	if len(posts) == 0 {
		return
	}

	rand.Seed(now.UnixNano())
	chosen := posts[rand.Intn(len(posts))]
	if err := m.publishAdminPost(chosen, now); err != nil {
		logs(fmt.Sprintf("failed to publish horoscope post: %v", err))
	}
}

func (m *Monitor) publishPositiveNewsRoutine() {
	rand.Seed(time.Now().UnixNano())
	delay := time.Hour + time.Duration(rand.Int63n(int64(time.Hour)))
	time.Sleep(delay)

	for {
		m.publishPositiveNews()
		rand.Seed(time.Now().UnixNano())
		delay = time.Hour + time.Duration(rand.Int63n(int64(time.Hour)))
		time.Sleep(delay)
	}
}

func (m *Monitor) publishPositiveNews() {
	news, err := getPositiveNewsPost()
	if err != nil {
		logs(fmt.Sprintf("failed to generate positive news: %v", err))
		return
	}

	rec := &telebot.Chat{ID: FrenlyNews}
	if _, err := b.Send(rec, news, telebot.NoPreview); err != nil {
		logs(fmt.Sprintf("failed to publish positive news: %v", err))
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	go m.publishAdminPostsRoutine()
	go m.publishHoroscopeRoutine()
	go m.publishPositiveNewsRoutine()
	return m
}
