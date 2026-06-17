package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"lesson01-ebvn/internal/handler/dto"
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

// @Summary ShorttenURL
// @Description Accepts a long URL and an expiration time, then generates a shortened key.
// @Param        request  body      dto.ShortenReq  true  "URL to be shortened and expiration time"
// @Tags shortenURL
// @Success 200 {object} map[string]interface{} "Success"
// @Router /v1/links/shorten [post]
func (b *BookMarkHandler) ShortenURL(ctx *gin.Context) {
	var input dto.ShortenReq
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}
	code, err := b.service.GenerateKey(ctx, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
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

// Redirect godoc
// @Summary      Redirect to original URL
// @Description  Retrieve the original URL from the short code and redirect the client to it
// @Tags         Redirect
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "Short URL Code"
// @Success      307   {string}  string  "Redirect to original URL"
// @Failure      400   {object}  map[string]string "code is required"
// @Failure      404   {object}  map[string]string "url not found"
// @Failure      500   {object}  map[string]string "internal server error"
// @Router       /v1/links/redirect/{code} [get]
func (b *BookMarkHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Code is required",
		})
		return
	}
	url, err := b.service.GetURL(ctx, code)
	if err != nil {
		log.Debug().Msg(err.Error())
		if errors.Is(err, service.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "url not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, url)

}
