package entity

import (
	"gorm.io/gorm"
)

type Ribbon struct {
	gorm.Model `json:"-"`

	UUID        string `json:"id"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	IsPremium   int    `json:"is_premium"`
	Description string `json:"description"`
	Slug        string `json:"-"`

	Items Contents `json:"items" sql:"-" gorm:"-"`
	Seo   struct {
		Slug string `json:"slug"`
	} `json:"seo" sql:"-" gorm:"-"`
	Metadata struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"metadata" sql:"-" gorm:"-"`
}

type Ribbons []Ribbon

func (Ribbon) TableName() string  { return "ribbons" }
func (Ribbons) TableName() string { return Content{}.TableName() }

func (r *Ribbon) Mark() {
	r.Slug = r.Seo.Slug
}
