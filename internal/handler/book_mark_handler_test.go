package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

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
