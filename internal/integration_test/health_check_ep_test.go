package intergration_test

import (
	"github.com/stretchr/testify/assert"
	"lesson01-ebvn/internal/api"
	"lesson01-ebvn/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookMarkEndpoint(t *testing.T) {
	t.Parallel()

	testConfig := &config.Config{
		ServiceName: "bookmark-service-test",
		InstanceID:  "test-instance-id",
		Port:        "8080",
	}

	tests := []struct {
		name string

		setUpTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedCode int

		expectedResp string
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
			testAPI := api.NewEngine(testConfig)

			rec := tc.setUpTestHTTP(testAPI)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
