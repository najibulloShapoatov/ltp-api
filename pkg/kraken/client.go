package kraken

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Config struct {
	Url     string
	WithLog bool
}

type Client interface {
	GetTicker(ctx context.Context, pair string) (*TickerResponse, error)
}
type client struct {
	cfg *Config
}

func New(cfg *Config) Client {
	return &client{cfg: cfg}
}

func (c *client) GetTicker(ctx context.Context, pair string) (*TickerResponse, error) {
	return requestGet[TickerResponse](
		ctx,
		fmt.Sprintf("%s/0/public/Ticker?pair=%s", c.cfg.Url, pair),
		c.cfg.WithLog,
		nil,
	)
}

func requestGet[T any](ctx context.Context, uri string, log bool, headers map[string]string) (*T, error) {
	cl := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if log {
		zap.L().Info("Request",
			zap.String("url", uri),
			zap.String("header", fmt.Sprintf("%+v", req.Header)),
		)
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if log {
		zap.L().Info("Response",
			zap.String("requestUrl", uri),
			zap.String("requestHeader", fmt.Sprintf("%+v", req.Header)),
			zap.String("responseStatus", resp.Status),
			zap.String("response", string(body)),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	var data = new(T)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
