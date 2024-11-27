package navclient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type (
	Client struct {
		timeout time.Duration
		host    string
		token   string
	}
)

func NewClient(host string, token string, timeout time.Duration) *Client {
	return &Client{
		host:    host,
		timeout: timeout,
		token:   token,
	}
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, params any) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.host, endpoint)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	payload := bytes.NewReader(jsonBytes)
	req, err := http.NewRequestWithContext(ctx, method, baseURL, payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "ApiKey "+c.token)

	httpClient := &http.Client{
		Timeout: c.timeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

func (c *Client) TestConnection(ctx context.Context) error {
	resp, err := c.doJSON(ctx, http.MethodGet, "status", nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couldn't reach navigator,status code: %d", resp.StatusCode)
	}

	return nil
}
