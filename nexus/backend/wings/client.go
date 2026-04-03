package wings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nexus/backend/models"
)

// Client for interacting with Wings API
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new Wings API client
func NewClient(node *models.Node) *Client {
	baseURL := fmt.Sprintf("%s://%s:%d", node.Scheme, node.FQDN, node.WingsPort)

	return &Client{
		BaseURL: baseURL,
		Token:   node.Token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// doRequest performs HTTP requests with authentication
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// GetServerDetails retrieves server details from Wings
func (c *Client) GetServerDetails(serverUUID string) (*ServerDetails, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/servers/%s", serverUUID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("wings error: %s", errResp.ErrMsg)
	}

	var details ServerDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &details, nil
}

// CreateServer creates a new server on Wings
func (c *Client) CreateServer(serverUUID string, startOnCompletion bool) error {
	payload := CreateServerPayload{
		UUID:              serverUUID,
		StartOnCompletion: startOnCompletion,
	}

	resp, err := c.doRequest("POST", "/api/servers", payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("wings error: %s", errResp.Error)
	}

	return nil
}

// DeleteServer deletes a server from Wings
func (c *Client) DeleteServer(serverUUID string) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/api/servers/%s", serverUUID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("wings error: %s", errResp.Error)
	}

	return nil
}

// SendPowerAction sends a power action to the server
func (c *Client) SendPowerAction(serverUUID, action string) error {
	payload := PowerActionPayload{Action: action}

	resp, err := c.doRequest("POST", fmt.Sprintf("/api/servers/%s/power", serverUUID), payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("wings error: %s", errResp.Error)
	}

	return nil
}

// GetServerResources retrieves resource usage from Wings
func (c *Client) GetServerResources(serverUUID string) (*ServerResources, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/servers/%s/resources", serverUUID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("wings error: %s", errResp.ErrMsg)
	}

	var resources ServerResources
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resources, nil
}
