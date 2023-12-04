package main

import (
	"go-vieon-content-crawling-api/src/configs"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/transport/rest"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if db, err := configs.NewGormClient().Instance(); err != nil {
		panic("Can't connect to database via GORM: " + err.Error())
	} else {
		port := os.Getenv("PORT")
		models := []interface{}{
			&entity.WatchedRibbon{},
			&entity.WatchedContent{},
			&entity.Content{},
		}
		db.AutoMigrate(models...)
		router := gin.Default()
		rest.NewRouteConfig(router).Config(db)

		router.Run(":" + port)
	}
}
