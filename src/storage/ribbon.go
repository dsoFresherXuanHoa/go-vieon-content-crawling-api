package storage

import (
	"context"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
)

var (
	ErrSaveWatchedRibbon2DB = errors.New("save watched ribbon information failure")
	ErrFindAllWatchedRibbon = errors.New("find all watched ribbon failure")
)

type ribbonStorage struct {
	sql *sqlStorage
}

func NewRibbonStore(sql *sqlStorage) *ribbonStorage {
	return &ribbonStorage{sql: sql}
}

func (s *ribbonStorage) SaveWatchedRibbon(ctx context.Context, content entity.WatchedRibbon) (uuid *string, err error) {
	if err := s.sql.db.Table(entity.WatchedRibbon{}.TableName()).Create(&content).Error; err != nil {
		fmt.Println("Error while save watched ribbon to database: " + err.Error())
		return nil, ErrSaveWatchedRibbon2DB
	}
	return &content.UUID, nil
}

func (s *ribbonStorage) FindAllWatchedRibbonIds(ctx context.Context) ([]string, error) {
	var watchedRibbons entity.WatchedRibbons
	if err := s.sql.db.Table(entity.WatchedRibbon{}.TableName()).Find(&watchedRibbons).Error; err != nil {
		fmt.Println("Error while get all watched ribbon: " + err.Error())
		return nil, ErrFindAllWatchedRibbon
	} else {
		var watchedRibbonIds []string
		for _, watchedRibbon := range watchedRibbons {
			watchedRibbonIds = append(watchedRibbonIds, watchedRibbon.UUID)
		}
		return watchedRibbonIds, nil
	}
}
