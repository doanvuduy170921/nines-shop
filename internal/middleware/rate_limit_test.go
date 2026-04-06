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
		mockReturnVal  int64 // Phải là int64 ở đây
		expectedStatus int
	}{
		{
			name:           "Request hợp lệ - chưa vượt giới hạn",
			limit:          int64(5), // Ép kiểu cụ thể
			mockReturnVal:  int64(1), // Ép kiểu cụ thể
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Request bị chặn - vượt quá giới hạn",
			limit:          int64(5),
			mockReturnVal:  int64(6),
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(cache.MockRedis)

			// Truyền tt.mockReturnVal (đã là int64) vào Return
			mockRedis.On("Incr", mock.Anything, mock.Anything).Return(tt.mockReturnVal, nil)

			// 3. Setup Gin Router
			r := gin.New()
			// Gắn middleware với tham số limit truyền từ table
			r.Use(RateLimitMiddleware(mockRedis, tt.limit, time.Minute))

			r.GET("/test-api", func(c *gin.Context) {
				c.String(http.StatusOK, "success")
			})

			// 4. Tạo request giả lập
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test-api", nil)

			// 5. Thực thi
			r.ServeHTTP(w, req)

			// 6. Kiểm tra kết quả bằng testify/assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Đảm bảo hàm Incr thực sự được gọi đúng 1 lần
			mockRedis.AssertExpectations(t)
		})
	}
}
