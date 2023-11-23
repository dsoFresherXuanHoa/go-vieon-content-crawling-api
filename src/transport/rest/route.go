package rest

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type routeConfig struct {
	router *gin.Engine
}

func NewRouteConfig(router *gin.Engine) *routeConfig {
	return &routeConfig{router: router}
}

func (cfg routeConfig) Config(db *gorm.DB) {
	v1 := cfg.router.Group("/api/v1")
	{
		contents := v1.Group("/contents/crawl")
		{
			contents.GET("/sync", SyncCrawlContent(db))
		}
	}
}
