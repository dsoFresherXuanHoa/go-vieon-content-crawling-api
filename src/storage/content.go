package storage

import (
	"context"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
)

var (
	ErrFindAllWatchedContents = errors.New("find all watched contents information failure")
	ErrSaveContent2DB         = errors.New("save contents information failure")
	ErrSaveWatchedContent2DB  = errors.New("save watched contents information failure")
)

type contentStorage struct {
	sql *sqlStorage
}

func NewContentStore(sql *sqlStorage) *contentStorage {
	return &contentStorage{sql: sql}
}

func (s *contentStorage) SaveContent(ctx context.Context, content entity.Content) (uuid *string, err error) {
	if err := s.sql.db.Table(entity.Content{}.TableName()).Create(&content).Error; err != nil {
		fmt.Println("Error while save content to database: " + err.Error())
		return nil, ErrSaveContent2DB
	}
	return &content.UUID, err
}

func (s *contentStorage) BatchSaveContent(ctx context.Context, contents entity.Contents) (*int, error) {
	totalContent := len(contents)
	if err := s.sql.db.Table(entity.Content{}.TableName()).CreateInBatches(contents, 10000).Error; err != nil {
		fmt.Println("Error while batch save content to database: " + err.Error())
		return nil, ErrSaveContent2DB
	}
	return &totalContent, nil
}

func (s *contentStorage) SaveWatchedContent(ctx context.Context, content entity.WatchedContent) (uuid *string, err error) {
	if err := s.sql.db.Table(entity.WatchedContent{}.TableName()).Create(&content).Error; err != nil {
		fmt.Println("Error while save watched content to database: " + err.Error())
		return nil, ErrSaveWatchedContent2DB
	}
	return &content.UUID, nil
}

func (s *contentStorage) BatchSaveWatchedContent(ctx context.Context, contents entity.WatchedContents) (*int, error) {
	totalContent := len(contents)
	if err := s.sql.db.Table(entity.WatchedContent{}.TableName()).CreateInBatches(&contents, 10000).Error; err != nil {
		fmt.Println("Error while save watched content to database: " + err.Error())
		return nil, ErrSaveWatchedContent2DB
	}
	return &totalContent, nil
}

func (s *contentStorage) FindAllWatchedContentIds(ctx context.Context) ([]string, error) {
	var watchedContents entity.WatchedContents
	if err := s.sql.db.Table(entity.WatchedContent{}.TableName()).Find(&watchedContents).Error; err != nil {
		fmt.Println("Error while get all watched contents: " + err.Error())
		return nil, ErrFindAllWatchedContents
	} else {
		var watchedContentIds []string
		for _, watchedContent := range watchedContents {
			watchedContentIds = append(watchedContentIds, watchedContent.UUID)
		}
		return watchedContentIds, nil
	}
}
