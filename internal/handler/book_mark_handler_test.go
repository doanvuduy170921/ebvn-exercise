package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"strings"

	"lesson01-ebvn/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookMarkHandler_HealthCheck(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name             string
		setUpRequest     func(ctx *gin.Context)
		setUpMockService func(ctx context.Context) *mocks.BookMarkService

		expectedStatus int

		expectResponse string
	}{
		{
			name: "success",
			setUpRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			},
			setUpMockService: func(ctx context.Context) *mocks.BookMarkService {
				serviceMock := mocks.NewBookMarkService(t)
				serviceMock.On("GetHealthInfo").Return("service bookmark", "core instance_id")
				return serviceMock
			},
			expectedStatus: http.StatusOK,
			expectResponse: `{"instance_id":"core instance_id","message":"OK","service_name":"service bookmark"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			tc.setUpRequest(ctx)

			mockService := tc.setUpMockService(ctx)

			testHandler := NewBookMarkHandler(mockService)

			testHandler.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectResponse, recorder.Body.String())

		})
	}
}

func TestBookMarkHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setUpRequest     func(ctx *gin.Context)
		setUpMockService func(ctx context.Context) *mocks.BookMarkService
		expectedStatus   int
		expectResponse   string
	}{
		{
			name: "success",
			setUpRequest: func(ctx *gin.Context) {
				body := `{"url":"https://google.com","exp":604800}`
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setUpMockService: func(ctx context.Context) *mocks.BookMarkService {
				serviceMock := mocks.NewBookMarkService(t)
				serviceMock.On("GenerateKey", ctx, "https://google.com", 604800).
					Return("abc1234", nil)
				return serviceMock
			},
			expectedStatus: http.StatusOK,
			expectResponse: `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "bad request - missing url",
			setUpRequest: func(ctx *gin.Context) {
				body := `{"exp":604800}`
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setUpMockService: func(ctx context.Context) *mocks.BookMarkService {
				return mocks.NewBookMarkService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectResponse: `{"code":400,"message":"Key: 'ShortenReq.Url' Error:Field validation for 'Url' failed on the 'required' tag"}`,
		},
		{
			name: "internal server error - service error",
			setUpRequest: func(ctx *gin.Context) {
				body := `{"url":"https://google.com","exp":604800}`
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setUpMockService: func(ctx context.Context) *mocks.BookMarkService {
				serviceMock := mocks.NewBookMarkService(t)
				serviceMock.On("GenerateKey", ctx, "https://google.com", 604800).
					Return("", errors.New("redis error"))
				return serviceMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectResponse: `{"code":500,"message":"redis error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			tc.setUpRequest(ctx)
			mockService := tc.setUpMockService(ctx)
			testHandler := NewBookMarkHandler(mockService)
			testHandler.ShortenURL(ctx)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectResponse, recorder.Body.String())
		})
	}
}
