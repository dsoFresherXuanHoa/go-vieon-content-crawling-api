package entity

import (
	"fmt"

	"gorm.io/gorm"
)

type Content struct {
	gorm.Model `json:"-"`

	UUID      string `json:"id" gorm:"uuid"`
	Type      int    `json:"type"`
	Category  string `json:"-"`
	Title     string `json:"title"`
	ShowName  string `json:"-"`
	EnableAds int    `json:"enable_ads"`
	IsPremium int    `json:"is_premium"`
	People    struct {
		Director []struct {
			Name string `json:"name"`
		} `json:"director"`
		Actor []struct {
			Name string `json:"name"`
		} `json:"actor"`
	} `json:"people" sql:"-" gorm:"-"`
	Director []string `json:"-" sql:"-" gorm:"-"`
	Actor    []string `json:"-" sql:"-" gorm:"-"`
	Tags     []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"tags" sql:"-" gorm:"-"`
	Country     string `json:"-"`
	Genre       string `json:"-"`
	ReleaseYear int    `json:"release_year"`
	Images      struct {
		ThumbnailHot string `json:"thumbnail_hot_v4"`
		Thumbnail    string `json:"thumbnail_v4"`
		Poster       string `json:"poster_v4"`
	} `json:"images" sql:"-" gorm:"-"`
	Thumbnail   string  `json:"-"`
	Poster      string  `json:"-"`
	Rating      float64 `json:"avg_rate"`
	TotalRating int     `json:"total_rate"`
	Seo         struct {
		Slug string `json:"slug"`
	} `json:"seo" sql:"-" gorm:"-"`
	IsEnd          bool `json:"is_end"`
	IsDownloadable int  `json:"is_downloadable"`
	Movie          struct {
		Title          string `json:"title"`
		Episode        int    `json:"episode"`
		CurrentEpisode string `json:"current_episode"`
	} `json:"movie" sql:"-" gorm:"-"`
	Episode          int    `json:"-"`
	CurrentEpisode   string `json:"-"`
	PublicDate       int    `json:"created_at"`
	IsComingSoon     int    `json:"is_coming_soon"`
	TrialDuration    int    `json:"trial_duration"`
	ForceLogin       int    `json:"force_login"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Resolution       int    `json:"resolution"`
	IsPremiumType    string `json:"is_premium_display"`
	Slug             string `json:"slug"`
	AllowKid         bool   `json:"allows_kid"`
}

type Contents []Content

func (Content) TableName() string  { return "contents" }
func (Contents) TableName() string { return Content{}.TableName() }

func (c *Content) Mark() {
	if len(c.Tags) != 0 {
		for _, tag := range c.Tags {
			if tag.Type == "country" {
				c.Country = tag.Name
			} else if tag.Type == "category" {
				c.Category = tag.Name
			} else if tag.Type == "genre" {
				c.Genre += fmt.Sprint(tag.Name, " ")
			}
		}
	}
	if len(c.Images.Thumbnail) == 0 {
		c.Thumbnail = c.Images.ThumbnailHot
	} else {
		c.Thumbnail = c.Images.Thumbnail
	}
	c.Poster = c.Images.Poster
	c.Episode = c.Movie.Episode
	c.CurrentEpisode = c.Movie.CurrentEpisode
	c.ShowName = c.Movie.Title
}
