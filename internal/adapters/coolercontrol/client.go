package coolercontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

type Client struct {
	BaseURL string
	Path    string
	Token   string
	HTTP    *http.Client
}

type Observation struct {
	Payload json.RawMessage
}

var _ provider.Adapter = Client{}

func New(baseURL, path, token string, timeout time.Duration) Client {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Path:    path,
		Token:   token,
		HTTP:    &http.Client{Timeout: timeout},
	}
}

func (c Client) Name() string {
	return "coolercontrol"
}

func (c Client) Capabilities() []domain.Capability {
	return []domain.Capability{{
		Resource:   domain.ResourcePWM,
		Operations: []domain.Operation{domain.OperationObserve},
	}}
}

func (c Client) Health(context.Context) domain.ProviderHealth {
	return domain.ProviderHealth{Provider: c.Name()}
}

func (c Client) Observe(ctx context.Context) (any, error) {
	if c.BaseURL == "" || c.Path == "" {
		return nil, fmt.Errorf("CoolerControl base URL and read-only path are required")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/"+strings.TrimLeft(c.Path, "/"), nil)
	if err != nil {
		return nil, fmt.Errorf("create CoolerControl request: %w", err)
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request CoolerControl API: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("CoolerControl API returned HTTP %d", response.StatusCode)
	}
	payload, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("read CoolerControl response: %w", err)
	}
	if !json.Valid(payload) {
		return nil, fmt.Errorf("CoolerControl response is not valid JSON")
	}
	return Observation{Payload: append(json.RawMessage(nil), payload...)}, nil
}
