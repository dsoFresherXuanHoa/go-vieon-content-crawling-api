package entity

import "gorm.io/gorm"

type WatchedRibbon struct {
	gorm.Model
	UUID string `json:"id" gorm:"unique"`
}

type WatchedRibbons []WatchedRibbon

func (WatchedRibbon) TableName() string  { return "watched_ribbons" }
func (WatchedRibbons) TableName() string { return WatchedRibbon{}.TableName() }
