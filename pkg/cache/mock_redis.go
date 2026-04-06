package cache

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Incr(key string, ttl time.Duration) (int64, error) {
	args := m.Called(key, ttl)

	// Lấy giá trị đầu tiên trả về từ mock
	res := args.Get(0)

	// Kiểm tra nếu giá trị là kiểu int (thường gặp khi truyền số trực tiếp trong test)
	if v, ok := res.(int); ok {
		return int64(v), args.Error(1)
	}

	// Kiểm tra nếu giá trị là kiểu int64
	if v, ok := res.(int64); ok {
		return v, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *MockRedis) Set(k string, v interface{}, t time.Duration) error { return nil }
func (m *MockRedis) Get(k string, d any) error                          { return nil }
func (m *MockRedis) Clear(k string) error                               { return nil }
func (m *MockRedis) Exists(k string) (bool, error)                      { return false, nil }
