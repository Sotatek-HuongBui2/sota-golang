package middlewares

import (
	"errors"
	"net/http"
	"strings"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/services"

	"github.com/gin-gonic/gin"
)

// check authentication
func RequireAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		arr := strings.Split(token, " ")

		if len(arr) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
			return

		}

		if arr[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
			return
		}
		token = arr[1]
		claims, err := services.VerifyToken(token)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
			return
		}
		ctx.Set("user_id", claims.UserId)

		ctx.Next()
	}
}
