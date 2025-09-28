package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) ExecuteRequest(ctx context.Context, request *models.QueryRequest) (*QueryResponse, error) {
	fullURL := c.baseURL + request.URL

	var body io.Reader
	if request.Method == "POST" || request.Method == "PUT" || request.Method == "PATCH" {
		if len(request.Args) > 0 {
			jsonBody, err := json.Marshal(request.Args)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			body = bytes.NewReader(jsonBody)
		}
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &QueryResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       responseBody,
		Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
	}, nil
}

type QueryResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Success    bool
}

func (r *QueryResponse) String() string {
	var result strings.Builder

	result.WriteString(strconv.Itoa(r.StatusCode))

	if len(r.Body) > 0 {
		result.WriteString(fmt.Sprintf(":%s", string(r.Body)))
	}

	return result.String()
}

func (r *QueryResponse) GetBodyAsString() string {
	return string(r.Body)
}

func (r *QueryResponse) GetBodyAsJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
