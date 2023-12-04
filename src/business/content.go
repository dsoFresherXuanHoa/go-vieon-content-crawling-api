package business

import (
	"context"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/utils"
	"math"
	"strconv"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
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
	BatchSaveContent(ctx context.Context, contents entity.Contents) (*int, error)
	SaveWatchedContent(ctx context.Context, content entity.WatchedContent) (uuid *string, err error)
	BatchSaveWatchedContent(ctx context.Context, contents entity.WatchedContents) (*int, error)
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
		var targetWatchedContents []entity.WatchedContent
		fmt.Println("Watched Content: ", len(watchedContentIds), len(contentIds), len(targetContentIds))
		for _, contentId := range targetContentIds {
			if !slices.Contains(watchedContentIds, contentId) && contentId != "" {
				targetWatchedContents = append(targetWatchedContents, entity.WatchedContent{Model: gorm.Model{}, UUID: contentId})
				result = append(result, contentId)
			}

		}
		fmt.Println(len(targetWatchedContents))
		fmt.Println(len(result))
		if _, err := business.contentStorage.BatchSaveWatchedContent(ctx, targetWatchedContents); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (business *contentBusiness) SyncCrawlContent(ctx context.Context, contentIds []string) error {
	var targetContents []entity.Content
	for _, contentId := range contentIds {
		url := "https://api.vieon.vn/backend/cm/v5/content/" + contentId
		if content, err := utils.NewNetUtil().CrawlContent(url); err != nil {
			fmt.Println("Error crawling content by id " + contentId + ": " + err.Error())
			return ErrCrawlingContentById
		} else if content.Title != "" {
			content.Mark()
			targetContents = append(targetContents, *content)
		}
	}
	totalContent := len(targetContents)
	totalSaveTime := (totalContent / 2000) + 1
	fmt.Println(totalSaveTime, totalContent)
	for i := 0; i < totalSaveTime; i++ {
		startIndex := i * 2000
		endIndex := int(math.Min(float64(startIndex+2000), float64(len(targetContents))))
		fmt.Println(startIndex, endIndex)
		if _, err := business.contentStorage.BatchSaveContent(ctx, targetContents[startIndex:endIndex]); err != nil {
			fmt.Println("Error while save batch content: " + err.Error())
			return err
		}
	}
	return nil
}
