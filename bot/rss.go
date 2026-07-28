package bot

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// RSS feed XML structures
type RssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title string    `xml:"title"`
	Link  string    `xml:"link"`
	Items []RssItem `xml:"item"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

const (
	ChinaDailyWorldRSS = "https://www.chinadaily.com.cn/rss/world_rss.xml"
	SourceChinaDaily   = "chinadaily"
	RssFetchInterval   = time.Hour
)

func fetchAndStoreRssFeed() {
	feed, err := fetchRss(ChinaDailyWorldRSS)
	if err != nil {
		logs(fmt.Sprintf("rss: failed to fetch feed from %s: %v", ChinaDailyWorldRSS, err))
		return
	}

	saved := 0
	for _, item := range feed.Channel.Items {
		guid := item.Guid
		if guid == "" {
			guid = item.Link
		}
		if guid == "" {
			continue
		}

		rp := RssPost{
			Source:  SourceChinaDaily,
			Guid:    guid,
			Title:   item.Title,
			Content: item.Description,
			Link:    item.Link,
		}

		res := db.Where("source = ? AND guid = ?", SourceChinaDaily, guid).FirstOrCreate(&rp)
		if res.Error != nil {
			logs(fmt.Sprintf("rss: failed to save post %s: %v", guid, res.Error))
			continue
		}
		if res.RowsAffected > 0 {
			saved++
		}
	}

	if saved > 0 {
		logs(fmt.Sprintf("rss: saved %d new posts from %s", saved, ChinaDailyWorldRSS))
	}
}

func fetchRss(url string) (*RssFeed, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var feed RssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

func startRssMonitor() {
	fetchAndStoreRssFeed()

	for {
		time.Sleep(RssFetchInterval)
		fetchAndStoreRssFeed()
	}
}
