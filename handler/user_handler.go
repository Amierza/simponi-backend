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
	IUserHandler interface {
		GetProfile(ctx *gin.Context)
	}

	userHandler struct {
		userService service.IUserService
		logger      *zap.Logger
	}
)

func NewUserHandler(userService service.IUserService, logger *zap.Logger) *userHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (uh *userHandler) GetProfile(ctx *gin.Context) {
	result, err := uh.userService.GetProfile(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_PROFILE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_PROFILE), result)
	ctx.JSON(http.StatusOK, res)
}
