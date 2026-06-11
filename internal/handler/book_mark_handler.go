package handler

import (
	"github.com/gin-gonic/gin"
	"lesson01-ebvn/dto"
	"lesson01-ebvn/internal/service"
	"net/http"
)

// BookMarkHandler handles HTTP requests related to bookmarks.
type BookMarkHandler struct {
	service service.BookMarkService
}

// NewBookMarkHandler creates a new BookMarkHandler with the given service.
func NewBookMarkHandler(service service.BookMarkService) *BookMarkHandler {
	return &BookMarkHandler{
		service: service,
	}
}

// @Summary HealthCheck
// @Description Get serviceName and instance_id
// @Tags healthCheck
// @Success 200 {object} map[string]interface{} "Success"
// @Router /health-check [get]
func (b *BookMarkHandler) HealthCheck(ctx *gin.Context) {
	serviceName, instanceId := b.service.GetHealthInfo()
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "OK",
		"service_name": serviceName,
		"instance_id":  instanceId,
	})
}

func (b *BookMarkHandler) ShortenURL(ctx *gin.Context) {
	var input dto.ShortenReq
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}
	code, err := b.service.GenerateKey(ctx.Request.Context(), input.Url, input.Exp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "Shorten URL generated successfully!",
	})
}
