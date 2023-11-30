package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/utils"

	"github.com/gocolly/colly"
	"golang.org/x/exp/slices"
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
	ErrGetRibbonProps = errors.New("get ribbonProps failure")
	ErrGetRibbonId    = errors.New("get ribbonId failure")
)

var (
	badRibbonIds = []string{"data-rm1st-becauseyouwatched-DTT", "data-recommendforyou-DTT", "data-suggestforyou-DTT", "data-cache-VieONselection", "WATCHING_RIBBON"}
)

type RibbonStorage interface {
	SaveWatchedRibbon(ctx context.Context, content entity.WatchedRibbon) (uuid *string, err error)
	FindAllWatchedRibbonIds(ctx context.Context) ([]string, error)
}

type ribbonBusiness struct {
	ribbonStorage RibbonStorage
}

func NewRibbonBusiness(ribbonStorage RibbonStorage) *ribbonBusiness {
	return &ribbonBusiness{ribbonStorage: ribbonStorage}
}

func (business *ribbonBusiness) GetRibbonSlug(url string) ([]string, error) {
	var ribbonProps ribbonProps
	var shareUrls []string
	isSuccess := true

	c := colly.NewCollector()
	c.OnHTML("#__NEXT_DATA__", func(h *colly.HTMLElement) {
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

func (business *ribbonBusiness) GetRibbonId(ctx context.Context, urls []string) ([]string, error) {
	var ribbonIds []string

	// Start New Crawl
	c := colly.NewCollector()
	for _, url := range urls {
		c.OnHTML("#SECTION_HOME", func(h *colly.HTMLElement) {
			h.ForEach(".rocopa", func(_ int, h *colly.HTMLElement) {
				ribbonIds = append(ribbonIds, h.Attr("id"))
			})
		})
		c.Visit(url)
	}

	// Save Watched Ribbon
	targetRibbonIds := utils.NewSliceUtil().RemoveDuplicatedStringSlice(ribbonIds, badRibbonIds)
	if watchedRibbonIds, err := business.ribbonStorage.FindAllWatchedRibbonIds(ctx); err != nil {
		return nil, err
	} else {
		for _, ribbonId := range targetRibbonIds {
			if !slices.Contains(watchedRibbonIds, ribbonId) {
				if _, err := business.ribbonStorage.SaveWatchedRibbon(ctx, entity.WatchedRibbon{UUID: ribbonId}); err != nil {
					return nil, err
				}
			}
		}
		return targetRibbonIds, nil
	}
}
