package rest

import (
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
		ribbonStorage := storage.NewRibbonStore(sqlStorage)
		contentStorage := storage.NewContentStore(sqlStorage)
		ribbonBusiness := business.NewRibbonBusiness(ribbonStorage)
		contentBusiness := business.NewContentBusiness(contentStorage)

		if ribbonSlug, err := ribbonBusiness.GetRibbonSlug("https://vieon.vn"); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetRibbonSlugFailure))
		} else if ribbonIds, err := ribbonBusiness.GetRibbonId(ctx, ribbonSlug); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetRibbonIdFailure))
		} else if contentIds, err := contentBusiness.GetContentId(ctx, ribbonIds); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), GetContentIdFailure))
		} else if err := contentBusiness.SyncCrawlContent(ctx, contentIds); err != nil {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(nil, http.StatusOK, constants.StatusOK, err.Error(), SyncCrawlContentsFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(gin.H{
				"ribbonIds":  len(ribbonIds),
				"contentIds": len(contentIds),
			}, http.StatusOK, constants.StatusOK, "", SyncCrawlContentsSuccess))
		}
	}
}
