package middleware

import (
	"net/http"
	"net/http/httptest"
	"nineshop-be/pkg/cache"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		limit          int64
		mockReturnVal  int64
		expectedStatus int
	}{
		{
			name:           "Valid request - limit not exceeded",
			limit:          int64(5),
			mockReturnVal:  int64(1),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Request blocked - limit exceeded",
			limit:          int64(5),
			mockReturnVal:  int64(6),
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(cache.MockRedis)
			mockRedis.On("Incr", mock.Anything, mock.Anything).Return(tt.mockReturnVal, nil)

			r := gin.New()
			r.Use(RateLimitMiddleware(mockRedis, tt.limit, time.Minute))

			r.GET("/test-api", func(c *gin.Context) {
				c.String(http.StatusOK, "success")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test-api", nil)

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRedis.AssertExpectations(t)
		})
	}
}
