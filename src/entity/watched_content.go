package entity

import "gorm.io/gorm"

type WatchedContent struct {
	gorm.Model
	UUID string `json:"-" gorm:"uuid;unique"`
}

type WatchedContents []WatchedContent

func (WatchedContent) TableName() string  { return "watched_contents" }
func (WatchedContents) TableName() string { return WatchedContent{}.TableName() }
