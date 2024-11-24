package mocks

import (
	"bytes"
	"io"
	"net/http"
)

// HTTPClient интерфейс для мока
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// MockClient это структура для мока HTTP-клиента
type MockClient struct {
	Response *http.Response
	Err      error
}

// Get это мок-метод для HTTP-запросов
func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

// Создание мока
func NewMockClient(body string, statusCode int, err error) *MockClient {
	return &MockClient{
		Response: &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		},
		Err: err,
	}
}
