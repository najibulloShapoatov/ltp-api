package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	Register(router *gin.RouterGroup)
}

func NotFoundHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
