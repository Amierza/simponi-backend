package middleware

import (
	"net/http"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/response"
	"github.com/gin-gonic/gin"
)

func Authentication(jwtService jwt.IJWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		tokenString := parts[1]

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_EXPIRED)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role_id", claims.RoleID)

		ctx.Next()
	}
}
