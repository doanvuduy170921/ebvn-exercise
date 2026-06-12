package dto

// ShortenReq represents the request body for the shorten URL endpoint.

type ShortenReq struct {
	Exp int    `json:"exp"`
	Url string `json:"url" binding:"required"`
}
