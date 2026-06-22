package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IProductReviewHandler interface {
		CreateProductReview(ctx *gin.Context)
		GetProductReviews(ctx *gin.Context)
	}

	productReviewHandler struct {
		productReviewService service.IProductReviewService
		logger               *zap.Logger
	}
)

func NewProductReviewHandler(productReviewService service.IProductReviewService, logger *zap.Logger) *productReviewHandler {
	return &productReviewHandler{
		productReviewService: productReviewService,
		logger:               logger,
	}
}

// CreateProductReview godoc
//
//	@Summary		Create product review (auto-tagging)
//	@Description	Create a product review. The review text is sent to the ML service which returns tags that are stored with the review (Requires permission: CreateProductReview)
//	@Tags			Product Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string							true	"Store ID (UUID)"
//	@Param			product_id	path		string							true	"Product ID (UUID)"
//	@Param			payload		body		dto.CreateProductReviewRequest	true	"Create product review request"
//	@Success		201			{object}	dto.ProductReviewResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse				"Invalid input / UUID"
//	@Failure		401			{object}	dto.ErrorResponse				"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse				"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse				"Product not found"
//	@Failure		500			{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/stores/{store_id}/products/{product_id}/reviews [post]
func (prh *productReviewHandler) CreateProductReview(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		prh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(dto.FAILED_CREATE_PRODUCT_REVIEW, dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		prh.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(dto.FAILED_CREATE_PRODUCT_REVIEW, dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.CreateProductReviewRequest
	payload.StoreID = &storeID
	payload.ProductID = &productID
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		prh.logger.Error("invalid create product review request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_CREATE_PRODUCT_REVIEW, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := prh.productReviewService.CreateProductReview(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_CREATE_PRODUCT_REVIEW, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_CREATE_PRODUCT_REVIEW, result)
	ctx.JSON(http.StatusCreated, res)
}

// GetProductReviews godoc
//
//	@Summary		Get product reviews
//	@Description	Get paginated reviews of a product (Requires permission: GetProductReviews)
//	@Tags			Product Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string								true	"Store ID (UUID)"
//	@Param			product_id	path		string								true	"Product ID (UUID)"
//	@Param			page		query		int									false	"Page number"
//	@Param			per_page	query		int									false	"Items per page"
//	@Success		200			{object}	dto.ProductReviewsResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse					"Invalid UUID"
//	@Failure		401			{object}	dto.ErrorResponse					"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse					"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse					"Product not found"
//	@Failure		500			{object}	dto.ErrorResponse					"Internal Server Error"
//	@Router			/stores/{store_id}/products/{product_id}/reviews [get]
func (prh *productReviewHandler) GetProductReviews(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		prh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_REVIEWS, dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		prh.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_REVIEWS, dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		prh.logger.Error("invalid get product reviews query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_REVIEWS, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := prh.productReviewService.GetProductReviews(ctx, &payload, &storeID, &productID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_REVIEWS, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.SUCCESS_GET_PRODUCT_REVIEWS,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}
