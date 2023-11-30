package business

import (
	"context"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/utils"
	"strconv"

	"golang.org/x/exp/slices"
)

var (
	ErrGetContentId            = errors.New("get contentId failure")
	ErrCrawlingContentByRibbon = errors.New("get contentId by ribbonId failure")
	ErrCrawlingContentById     = errors.New("get content by contentId failure")
)

var (
	badContentIds = []string{}
)

type ContentStorage interface {
	SaveContent(ctx context.Context, content entity.Content) (uuid *string, err error)
	SaveWatchedContent(ctx context.Context, content entity.WatchedContent) (uuid *string, err error)
	FindAllWatchedContentIds(ctx context.Context) ([]string, error)
}

type contentBusiness struct {
	contentStorage ContentStorage
}

func NewContentBusiness(contentStorage ContentStorage) *contentBusiness {
	return &contentBusiness{contentStorage: contentStorage}
}

func (business *contentBusiness) GetContentId(ctx context.Context, ribbonIds []string) ([]string, error) {
	var contentIds []string

	// Start New Crawl
	for _, ribbonId := range ribbonIds {
		url := "https://api.vieon.vn/backend/cm/v5/ribbon/" + ribbonId + "?platform=web&ui=012021"
		if ribbon, err := utils.NewNetUtil().CrawlRibbon(url); err != nil {
			fmt.Println("Error first crawling ribbon by id " + ribbonId + ": " + err.Error())
			return nil, ErrCrawlingContentByRibbon
		} else {
			total := ribbon.Metadata.Total
			page := total/30 + 1
			for i := 0; i < int(page); i++ {
				url := "https://api.vieon.vn/backend/cm/v5/ribbon/" + ribbonId + "?limit=30&page=" + strconv.Itoa(i) + "&platform=web&ui=012021"
				if ribbon, err := utils.NewNetUtil().CrawlRibbon(url); err != nil {
					fmt.Println("Error crawling ribbon by id " + ribbonId + ": " + err.Error())
					return nil, ErrCrawlingContentByRibbon
				} else {
					contents := ribbon.Items
					for _, content := range contents {
						contentIds = append(contentIds, content.UUID)
					}
				}
			}
		}
	}

	// Save Watched Content
	targetContentIds := utils.NewSliceUtil().RemoveDuplicatedStringSlice(contentIds, badContentIds)
	if watchedContentIds, err := business.contentStorage.FindAllWatchedContentIds(ctx); err != nil {
		return nil, err
	} else {
		var result []string
		fmt.Println("Watched Content: ", len(watchedContentIds), len(contentIds))
		for _, contentId := range targetContentIds {
			if !slices.Contains(watchedContentIds, contentId) {
				if _, err := business.contentStorage.SaveWatchedContent(ctx, entity.WatchedContent{UUID: contentId}); err != nil {
					return nil, err
				}
				result = append(result, contentId)
			}
		}
		return result, nil
	}
}

func (business *contentBusiness) SyncCrawlContent(ctx context.Context, contentIds []string) error {
	for _, contentId := range contentIds {
		url := "https://api.vieon.vn/backend/cm/v5/content/" + contentId
		if content, err := utils.NewNetUtil().CrawlContent(url); err != nil {
			fmt.Println("Error crawling content by id " + contentId + ": " + err.Error())
			return ErrCrawlingContentById
		} else {
			content.Mark()
			if _, err := business.contentStorage.SaveContent(ctx, *content); err != nil {
				fmt.Println("Error while save content: " + contentId + ": " + err.Error())
				return err
			}
		}
	}
	return nil
}
