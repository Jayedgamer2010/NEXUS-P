package wings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"nexus/backend/models"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *Client) buildURL(node models.Node, path string) string {
	base := fmt.Sprintf("%s://%s:%d", node.Scheme, node.FQDN, node.DaemonListen)
	return base + path
}

func (c *Client) doRequest(method, url string, node models.Node, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+node.DaemonToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			return nil, fmt.Errorf("wings unavailable: %w", err)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("wings request timed out: %w", err)
		}
		return nil, fmt.Errorf("wings request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) GetServerDetails(node models.Node, uuid string) (*ServerDetails, error) {
	url := c.buildURL(node, fmt.Sprintf("/api/servers/%s", uuid))
	resp, err := c.doRequest("GET", url, node, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.ErrMsg != "" {
			return nil, fmt.Errorf("wings error: %s", errResp.ErrMsg)
		}
		return nil, fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	var details ServerDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &details, nil
}

func (c *Client) CreateServer(node models.Node, payload CreateServerPayload) error {
	url := c.buildURL(node, "/api/servers")
	resp, err := c.doRequest("POST", url, node, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.ErrMsg != "" {
			return fmt.Errorf("wings error: %s", errResp.ErrMsg)
		}
		return fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteServer(node models.Node, uuid string) error {
	url := c.buildURL(node, fmt.Sprintf("/api/servers/%s", uuid))
	resp, err := c.doRequest("DELETE", url, node, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.ErrMsg != "" {
			return fmt.Errorf("wings error: %s", errResp.ErrMsg)
		}
		return fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SendPowerAction(node models.Node, uuid string, action string) error {
	url := c.buildURL(node, fmt.Sprintf("/api/servers/%s/power", uuid))
	payload := PowerActionPayload{Action: action}

	resp, err := c.doRequest("POST", url, node, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.ErrMsg != "" {
			return fmt.Errorf("wings error: %s", errResp.ErrMsg)
		}
		return fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetServerResources(node models.Node, uuid string) (*ServerResources, error) {
	url := c.buildURL(node, fmt.Sprintf("/api/servers/%s/resources/utilization", uuid))
	resp, err := c.doRequest("GET", url, node, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	var resources ServerResources
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resources, nil
}

func (c *Client) GetSystemInfo(node models.Node) (*SystemInfo, error) {
	url := c.buildURL(node, "/api/system")
	resp, err := c.doRequest("GET", url, node, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wings returned status %d", resp.StatusCode)
	}

	var info SystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}
