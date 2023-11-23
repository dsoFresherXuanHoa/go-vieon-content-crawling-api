package rest

import (
	"fmt"
	"go-vieon-content-crawling-api/src/business"
	"go-vieon-content-crawling-api/src/constants"
	"go-vieon-content-crawling-api/src/entity"
	"go-vieon-content-crawling-api/src/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	GetRibbonSlugFailure     = "Cannot finished get ribbonSlug: some many or all slug has been missing!"
	GetRibbonIdFailure       = "Cannot finished get ribbonId: some many or all id has been missing!"
	GetContentIdFailure      = "Cannot finished get contentId: some many or all id has been missing!"
	SyncCrawlContentsFailure = "Cannot finished crawl content information: some many or all contents has been missing!"
	SyncCrawlContentsSuccess = "All content has been crawl: congrats!"
)

func SyncCrawlContent(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		contentStorage := storage.NewContentStore(sqlStorage)
		contentBusiness := business.NewContentBusiness(contentStorage)

		if ribbonSlug, err := contentBusiness.GetRibbonSlug("https://vieon.vn"); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetRibbonSlugFailure))
		} else if ribbonIds, ribbonIdsFilePath, err := contentBusiness.GetRibbonId(ribbonSlug); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetRibbonIdFailure))
		} else if contentIds, contentIdsFilePath, err := contentBusiness.GetContentId(ribbonIds); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetContentIdFailure))
		} else if err := contentBusiness.SyncCrawlContent(ctx, contentIds[0:10]); err != nil {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(nil, http.StatusOK, constants.StatusOK, err.Error(), SyncCrawlContentsFailure))
		} else {
			fmt.Println("Save ribbonIds in: ", *ribbonIdsFilePath)
			fmt.Println("Save contentIds in: ", *contentIdsFilePath)
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(len(contentIds), http.StatusOK, constants.StatusOK, "", SyncCrawlContentsSuccess))
		}
	}
}
