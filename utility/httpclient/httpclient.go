package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HttpClient struct chứa thông tin cấu hình
type HttpClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

// NewHttpClient tạo một HttpClient mới
func NewHttpClient(baseURL string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Headers: make(map[string]string),
	}
}

// SetHeader thêm hoặc cập nhật header
func (c *HttpClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// makeRequest tạo yêu cầu HTTP chung
func (c *HttpClient) makeRequest(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	// Tạo URL với endpoint
	fullURL, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, err
	}

	// Thêm query params vào URL nếu có
	if params != nil {
		query := fullURL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		fullURL.RawQuery = query.Encode()
	}

	// Xử lý body nếu không nil
	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}

	// Tạo yêu cầu
	req, err := http.NewRequest(method, fullURL.String(), requestBody)
	if err != nil {
		return nil, err
	}

	// Gắn header
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Nếu body là JSON
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Gửi yêu cầu
	return c.HTTPClient.Do(req)
}

// GET yêu cầu GET
func (c *HttpClient) GET(endpoint string, params map[string]string) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, endpoint, nil, params)
}

// POST yêu cầu POST
func (c *HttpClient) POST(endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	return c.makeRequest(http.MethodPost, endpoint, body, params)
}

// PUT yêu cầu PUT
func (c *HttpClient) PUT(endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	return c.makeRequest(http.MethodPut, endpoint, body, params)
}

// DELETE yêu cầu DELETE
func (c *HttpClient) DELETE(endpoint string, params map[string]string) (*http.Response, error) {
	return c.makeRequest(http.MethodDelete, endpoint, nil, params)
}

// ParseJSONResponse chuyển đổi phản hồi thành JSON
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("API trả về mã lỗi: " + resp.Status)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}
