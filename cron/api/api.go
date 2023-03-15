package api

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

// API 封装了对某个接口的调用操作。
type API struct {
	endpoint string
	client   *http.Client
}

// NewAPI 创建一个新的 API 实例。
func NewAPI(endpoint string) *API {
	return &API{
		endpoint: endpoint,
		client:   http.DefaultClient,
	}
}

// Call 调用指定的 API 接口并返回响应结果。
func (a *API) Call(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)
