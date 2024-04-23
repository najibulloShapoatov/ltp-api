package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"ltp-api/internal/services"
	"net/http"
)

type LTPHandler struct {
	svc services.KrakenService
}

func NewLTPHandler(i *do.Injector) *LTPHandler {
	return &LTPHandler{
		svc: do.MustInvoke[services.KrakenService](i),
	}
}

func (h *LTPHandler) Register(router *gin.RouterGroup) {
	faq := router.Group("ltp")
	{
		faq.GET("", h.ltp)
	}
}

// ltp
// @Summary ltp
// @Schemes
// @Description ltp
// @Tags LTPHandler
// @Accept json
// @Produce json
// @Success 200 {object} dto.LTPResponse
// @Failure      400  string  "Bad Request"
// @Router /ltp [get]
func (h *LTPHandler) ltp(ctx *gin.Context) {

	data, err := h.svc.LTP(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, data)
		return
	}

	ctx.JSON(http.StatusOK, data)
}
