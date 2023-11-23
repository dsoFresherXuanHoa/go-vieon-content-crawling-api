package storage

import (
	"context"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
)

var (
	ErrSaveContent2DB = errors.New("save contents information failure")
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
	return &content.UUID, nil
}
