package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-vieon-content-crawling-api/src/entity"
	"io"
	"net/http"
)

var (
	ErrCrawlingTargetURL      = errors.New("crawling target url failure")
	ErrReadResponseBytes      = errors.New("read bytes from response failure")
	ErrBindingResponse2Struct = errors.New("binding response to struct failure")
)

type netUtil struct{}

func NewNetUtil() *netUtil {
	return &netUtil{}
}

func (netUtil) CrawlContent(url string) (content *entity.Content, err error) {
	if resp, err := http.Get(url); err != nil {
		fmt.Println("Error while crawling target url: " + err.Error())
		return nil, ErrCrawlingTargetURL
	} else if bytes, err := io.ReadAll(resp.Body); err != nil {
		defer resp.Body.Close()
		fmt.Println("Error while get bytes from res: " + err.Error())
		return nil, ErrReadResponseBytes
	} else {
		defer resp.Body.Close()
		var content entity.Content
		if err := json.Unmarshal(bytes, &content); err != nil {
			fmt.Println("Error while binding response to struct: " + err.Error())
			return nil, ErrBindingResponse2Struct
		} else {
			return &content, nil
		}
	}
}

func (netUtil) CrawlRibbon(url string) (content *entity.Ribbon, err error) {
	if resp, err := http.Get(url); err != nil {
		fmt.Println("Error while crawling target url: " + err.Error())
		return nil, ErrCrawlingTargetURL
	} else if bytes, err := io.ReadAll(resp.Body); err != nil {
		defer resp.Body.Close()
		fmt.Println("Error while get bytes from res: " + err.Error())
		return nil, ErrReadResponseBytes
	} else {
		defer resp.Body.Close()
		var ribbon entity.Ribbon
		if err := json.Unmarshal(bytes, &ribbon); err != nil {
			fmt.Println("Error while binding response to struct: " + err.Error())
			return nil, ErrBindingResponse2Struct
		} else {
			return &ribbon, nil
		}
	}
}
