package handler

import (
	"fmt"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	IProductCategoriesHandler interface {
		GetProductCategories(ctx *gin.Context)
	}

	productCategoriesHandler struct {
		productCategoriesService 	service.IProductCategoriesService
		logger         				*zap.Logger
	}
)

func NewProductCategoriesHandler(productCategoriesService service.IProductCategoriesService, logger *zap.Logger) *productCategoriesHandler {
	return &productCategoriesHandler{
		productCategoriesService: productCategoriesService,
		logger: logger,
	}
}

func (ph *productCategoriesHandler) GetProductCategories(ctx *gin.Context) {
	result, err := ph.productCategoriesService.GetProductCategories(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product categories", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product categories", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}