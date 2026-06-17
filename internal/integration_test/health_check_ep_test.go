package intergration_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"lesson01-ebvn/internal/api"
	"lesson01-ebvn/internal/config"
	"lesson01-ebvn/pkg/redis"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestAPI(t *testing.T) api.Engine {
	t.Helper()

	redisClient, err := redis.NewRedisClient()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
	}

	testConfig := &config.Config{
		ServiceName: "bookmark-service-test",
		InstanceID:  "test-instance-id",
		Port:        "8080",
	}

	return api.NewEngine(testConfig, redisClient)
}

func TestBookMarkEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	tests := []struct {
		name          string
		setUpTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		expectedCode  int
		expectedResp  string
	}{
		{
			name: "normal case",
			setUpTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"instance_id":"test-instance-id","message":"OK","service_name":"bookmark-service-test"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testAPI := setupTestAPI(t)
			rec := tc.setUpTestHTTP(testAPI)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}

func TestShortenURLEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	tests := []struct {
		name          string
		setUpTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		expectedCode  int
		expectedResp  string
		checkResp     bool
	}{
		{
			name: "success",
			setUpTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body, _ := json.Marshal(map[string]interface{}{
					"url": "https://google.com",
					"exp": 604800,
				})
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedCode: http.StatusOK,
			checkResp:    false,
		},
		{
			name: "bad request - missing url",
			setUpTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body, _ := json.Marshal(map[string]interface{}{
					"exp": 604800,
				})
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"code":400,"message":"Key: 'ShortenReq.Url' Error:Field validation for 'Url' failed on the 'required' tag"}`,
			checkResp:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testAPI := setupTestAPI(t)
			rec := tc.setUpTestHTTP(testAPI)
			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.checkResp {
				assert.Equal(t, tc.expectedResp, rec.Body.String())
			}
		})
	}
}

func TestRedirectEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	tests := []struct {
		name          string
		setUpTestHTTP func(api api.Engine) *httptest.ResponseRecorder
		expectedCode  int
		expectedResp  string
		checkResp     bool
	}{
		{
			name: "success - redirect to original url",
			setUpTestHTTP: func(a api.Engine) *httptest.ResponseRecorder {
				body, _ := json.Marshal(map[string]interface{}{
					"url": "https://google.com",
					"exp": 604800,
				})
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				a.ServeHTTP(rec, req)

				var shortenResp map[string]string
				json.Unmarshal(rec.Body.Bytes(), &shortenResp)
				code := shortenResp["code"]

				req2 := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+code, nil)
				rec2 := httptest.NewRecorder()
				a.ServeHTTP(rec2, req2)
				return rec2
			},
			expectedCode: http.StatusTemporaryRedirect,
			checkResp:    false,
		},
		{
			name: "not found - code does not exist",
			setUpTestHTTP: func(a api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notexist999", nil)
				rec := httptest.NewRecorder()
				a.ServeHTTP(rec, req)
				return rec
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"message":"url not found"}`,
			checkResp:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testAPI := setupTestAPI(t)
			rec := tc.setUpTestHTTP(testAPI)
			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.checkResp {
				assert.Equal(t, tc.expectedResp, rec.Body.String())
			}
		})
	}
}
