// Package client provides a Go SDK for the Holiday API Indonesia
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// Client is the Holiday API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// Option is a functional option for configuring the Client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// New creates a new Holiday API client
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Holiday represents an Indonesian holiday
type Holiday struct {
	ID          int       `json:"id"`
	Date        string    `json:"date"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// APIResponse represents the standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta represents pagination metadata
type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %s: %s", e.Code, e.Message)
}

// GetHolidays retrieves holidays with optional filters
func (c *Client) GetHolidays(ctx context.Context, year, month *int, holidayType string) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if year != nil {
		q.Add("year", fmt.Sprintf("%d", *year))
	}
	if month != nil {
		q.Add("month", fmt.Sprintf("%d", *month))
	}
	if holidayType != "" {
		q.Add("type", holidayType)
	}
	req.URL.RawQuery = q.Encode()

	return c.executeHolidayRequest(req)
}

// GetHolidaysByYear retrieves all holidays for a specific year
func (c *Client) GetHolidaysByYear(ctx context.Context, year int) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays/year/%d", c.baseURL, year)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.executeHolidayRequest(req)
}

// GetHolidaysByMonth retrieves holidays for a specific year and month
func (c *Client) GetHolidaysByMonth(ctx context.Context, year, month int) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays/month/%d/%d", c.baseURL, year, month)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.executeHolidayRequest(req)
}

// GetTodayHoliday retrieves today's holiday if any
func (c *Client) GetTodayHoliday(ctx context.Context) (*Holiday, error) {
	url := fmt.Sprintf("%s/holidays/today", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	holidays, err := c.executeHolidayRequest(req)
	if err != nil {
		return nil, err
	}

	if len(holidays) == 0 {
		return nil, nil
	}

	return &holidays[0], nil
}

// GetUpcomingHolidays retrieves upcoming holidays
func (c *Client) GetUpcomingHolidays(ctx context.Context, limit int) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays/upcoming?limit=%d", c.baseURL, limit)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.executeHolidayRequest(req)
}

// GetHolidaysThisYear retrieves holidays for the current year
func (c *Client) GetHolidaysThisYear(ctx context.Context) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays/this-year", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.executeHolidayRequest(req)
}

// GetHolidaysThisMonth retrieves holidays for the current month
func (c *Client) GetHolidaysThisMonth(ctx context.Context) ([]Holiday, error) {
	url := fmt.Sprintf("%s/holidays/this-month", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.executeHolidayRequest(req)
}

// HealthCheck checks if the API is healthy
func (c *Client) HealthCheck(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// executeHolidayRequest executes a request and returns holidays
func (c *Client) executeHolidayRequest(req *http.Request) ([]Holiday, error) {
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}
		if apiResp.Error != nil {
			return nil, apiResp.Error
		}
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return nil, apiResp.Error
		}
		return nil, fmt.Errorf("API request failed")
	}

	// Convert the data to holidays
	data, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var holidays []Holiday
	if err := json.Unmarshal(data, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}
