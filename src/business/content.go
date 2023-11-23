package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

type ribbonProps struct {
	Props struct {
		InitialState struct {
			Menu struct {
				MenuList []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Seo  struct {
						ShareURL string `json:"share_url"`
					} `json:"seo"`
					SubMenu []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
						Seo  struct {
							ShareURL string `json:"share_url"`
						} `json:"seo"`
					} `json:"subMenu"`
				} `json:"menuList"`
			} `json:"Menu"`
		} `json:"initialState"`
	} `json:"props"`
}

var (
	ErrGetRibbonProps          = errors.New("get ribbonProps failure")
	ErrGetRibbonId             = errors.New("get ribbonId failure")
	ErrGetContentId            = errors.New("get contentId failure")
	ErrCrawlingContentByRibbon = errors.New("get contentId by ribbonId failure")
	ErrCrawlingContentById     = errors.New("get content by contentId failure")
)

var (
	currentDir, _ = os.Getwd()
)

var (
	badRibbonIds  = []string{"data-rm1st-becauseyouwatched-DTT", "data-recommendforyou-DTT", "data-suggestforyou-DTT", "data-cache-VieONselection", "WATCHING_RIBBON"}
	badContentIds = []string{}
)

type RibbonStorage interface {
}

type ribbonBusiness struct {
	ribbonStorage RibbonStorage
}

func NewRibbonBusiness(ribbonStorage RibbonStorage) *ribbonBusiness {
	return &ribbonBusiness{ribbonStorage: ribbonStorage}
}

type ContentStorage interface {
	SaveContent(ctx context.Context, content entity.Content) (uuid *string, err error)
}

type contentBusiness struct {
	contentStorage ContentStorage
}

func NewContentBusiness(contentStorage ContentStorage) *contentBusiness {
	return &contentBusiness{contentStorage: contentStorage}
}

func (business *contentBusiness) GetRibbonSlug(url string) ([]string, error) {
	var ribbonProps ribbonProps
	var shareUrls []string
	c := colly.NewCollector()

	isSuccess := true
	c.OnHTML("script#__NEXT_DATA__", func(h *colly.HTMLElement) {
		if err := json.Unmarshal([]byte(h.Text), &ribbonProps); err != nil {
			fmt.Println("Error while ribbonProps: " + err.Error())
			isSuccess = false
		} else {
			for _, menu := range ribbonProps.Props.InitialState.Menu.MenuList {
				isFirst := true
				shareUrls = append(shareUrls, menu.Seo.ShareURL)
				if isFirst {
					for _, submenu := range menu.SubMenu {
						shareUrls = append(shareUrls, submenu.Seo.ShareURL)
					}
					isFirst = false
				}
			}
		}
	})
	c.Visit(url)

	if isSuccess {
		result := utils.NewSliceUtil().RemoveDuplicatedStringSlice(shareUrls, badRibbonIds)
		return result, nil
	} else {
		return nil, ErrGetRibbonProps
	}
}

func (business *contentBusiness) GetRibbonId(urls []string) ([]string, *string, error) {
	isSuccess := true

	// Read Previous Crawl
	var previousIds []string
	latestRibbonCrawl := filepath.Join(currentDir, "data", "csv", "latest", "ribbon.txt")
	if previousRibbonIdsLogPath, err := utils.NewTextUtil().File2String(latestRibbonCrawl); err == nil {
		if ribbonIds, err := utils.NewCSVUtil().CSV2ODSlice(*previousRibbonIdsLogPath); err == nil {
			previousIds = ribbonIds
		}
	}

	// Start New Crawl
	var ribbonIds []string
	c := colly.NewCollector()
	for _, url := range urls {
		c.OnHTML("#SECTION_HOME", func(h *colly.HTMLElement) {
			h.ForEach(".rocopa", func(_ int, h *colly.HTMLElement) {
				ribbonIds = append(ribbonIds, h.Attr("id"))
			})
		})
		c.Visit(url)
	}

	filePath := filepath.Join(currentDir, "data", "csv", "ribbon", fmt.Sprint(time.Now().Unix(), ".csv"))
	diffFilePath := filepath.Join(currentDir, "data", "csv", "diff", "ribbon", fmt.Sprint(time.Now().Unix(), ".csv"))
	if isSuccess {
		// Remove Duplicated and Write to CSV File
		latestRibbonIds := utils.NewSliceUtil().RemoveDuplicatedStringSlice(ribbonIds, badRibbonIds)
		result := utils.NewSliceUtil().Diff(previousIds, latestRibbonIds)
		if err := utils.NewCSVUtil().ODSlice2CSV(latestRibbonIds, filePath); err != nil {
			return latestRibbonIds, nil, err
		} else if err := utils.NewCSVUtil().ODSlice2CSV(result, diffFilePath); err != nil {
			return result, nil, err
		} else {
			// Write Latest Log
			if err := utils.NewTextUtil().String2File(filePath, latestRibbonCrawl); err != nil {
				return result, &filePath, err
			}
			return result, &filePath, nil
		}
	} else {
		return nil, nil, ErrGetRibbonId
	}
}

func (business *contentBusiness) GetContentId(ribbonIds []string) ([]string, *string, error) {
	// Read Previous Crawl
	var previousIds []string
	latestContentCrawl := filepath.Join(currentDir, "data", "csv", "latest", "content.txt")

	if previousContentIdsLogPath, err := utils.NewTextUtil().File2String(latestContentCrawl); err == nil {
		if contentIds, err := utils.NewCSVUtil().CSV2ODSlice(*previousContentIdsLogPath); err == nil {
			previousIds = contentIds
		}
	}

	// Start New Crawl
	var contentIds []string
	for _, ribbonId := range ribbonIds {
		url := "https://api.vieon.vn/backend/cm/v5/ribbon/" + ribbonId + "?platform=web&ui=012021"
		if ribbon, err := utils.NewNetUtil().CrawlRibbon(url); err != nil {
			fmt.Println("Error crawling ribbon by id " + ribbonId + ": " + err.Error())
			return nil, nil, ErrCrawlingContentByRibbon
		} else {
			total := ribbon.Metadata.Total
			page := total/30 + 1
			for i := 0; i < int(page); i++ {
				url := "https://api.vieon.vn/backend/cm/v5/ribbon/" + ribbonId + "?limit=30&page=" + strconv.Itoa(i) + "&platform=web&ui=012021"
				if ribbon, err := utils.NewNetUtil().CrawlRibbon(url); err != nil {
					fmt.Println("Error crawling ribbon by uuid " + ribbonId + ": " + err.Error())
					return nil, nil, ErrCrawlingContentByRibbon
				} else {
					contents := ribbon.Items
					for _, content := range contents {
						contentIds = append(contentIds, content.UUID)
					}
				}
			}
		}
	}

	filePath := filepath.Join(currentDir, "data", "csv", "content", fmt.Sprint(time.Now().Unix(), ".csv"))
	diffFilePath := filepath.Join(currentDir, "data", "csv", "diff", "content", fmt.Sprint(time.Now().Unix(), ".csv"))

	// Remove Duplicated and Write to CSV File
	latestContentIds := utils.NewSliceUtil().RemoveDuplicatedStringSlice(contentIds, badContentIds)
	result := utils.NewSliceUtil().Diff(previousIds, latestContentIds)
	if err := utils.NewCSVUtil().ODSlice2CSV(latestContentIds, filePath); err != nil {
		return latestContentIds, nil, err
	} else if err := utils.NewCSVUtil().ODSlice2CSV(result, diffFilePath); err != nil {
		return result, nil, err
	} else {
		// Write Latest Log
		if err := utils.NewTextUtil().String2File(filePath, latestContentCrawl); err != nil {
			return result, &filePath, err
		}
	}
	return result, &filePath, nil
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
