package handler

import (
	"github.com/gin-gonic/gin"
	"lesson01-ebvn/internal/service"
	"net/http"
)

type BookMarkHandler struct {
	service service.BookMarkService
}

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
