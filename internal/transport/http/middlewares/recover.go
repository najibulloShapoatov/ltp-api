package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rvr := recover(); rvr != nil {
				zap.S().Error(rvr)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, rvr)
			}
		}()
		ctx.Next()
	}
}
