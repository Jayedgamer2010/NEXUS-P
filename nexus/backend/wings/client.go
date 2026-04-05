package wings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"nexus/backend/models"
)

var (
	ErrWingsOffline   = errors.New("wings daemon is offline or unreachable")
	ErrWingsTimeout   = errors.New("wings daemon request timed out")
	ErrWingsAPI       = errors.New("wings API error")
	ErrInvalidAction  = errors.New("invalid power action")
)

type WingsClient struct {
	httpClient *http.Client
}

func NewWingsClient() *WingsClient {
	return &WingsClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WingsClient) buildURL(node models.Node, p string) string {
	if node.Scheme == "" {
		node.Scheme = "https"
	}
	if node.DaemonListen == 0 {
		node.DaemonListen = 8080
	}
	return fmt.Sprintf("%s://%s:%d%s", node.Scheme, node.FQDN, node.DaemonListen, p)
}

func (w *WingsClient) buildHeaders(node models.Node) http.Header {
	h := make(http.Header)
	h.Set("Authorization", fmt.Sprintf("Bearer %s.%s", node.DaemonTokenID, node.DaemonToken))
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	return h
}

func (w *WingsClient) doRequest(node models.Node, method, path string, body interface{}) (*http.Response, error) {
	url := w.buildURL(node, path)
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header = w.buildHeaders(node)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		var netErr netError
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				return nil, ErrWingsTimeout
			}
		}
		return nil, ErrWingsOffline
	}

	return resp, nil
}

type netError interface {
	error
	Timeout() bool
}

func (w *WingsClient) GetSystemInfo(node models.Node) (map[string]interface{}, error) {
	resp, err := w.doRequest(node, http.MethodGet, "/api/system", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (w *WingsClient) GetServerDetails(node models.Node, uuid string) (map[string]interface{}, error) {
	resp, err := w.doRequest(node, http.MethodGet, fmt.Sprintf("/api/servers/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("%w: status %d", ErrWingsAPI, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (w *WingsClient) CreateServer(node models.Node, payload CreateServerPayload) error {
	resp, err := w.doRequest(node, http.MethodPost, "/api/servers", payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("%w: status %d, response: %v", ErrWingsAPI, resp.StatusCode, errResp)
	}
	return nil
}

func (w *WingsClient) DeleteServer(node models.Node, uuid string, force bool) error {
	path := fmt.Sprintf("/api/servers/%s", uuid)
	if force {
		path += "?force=true"
	}
	resp, err := w.doRequest(node, http.MethodDelete, path, nil)
	if err != nil {
		// If Wings is offline, we still want to delete the server from DB
		return nil
	}
	defer resp.Body.Close()
	return nil
}

func (w *WingsClient) SendPowerAction(node models.Node, uuid string, action string) error {
	switch action {
	case "start", "stop", "restart", "kill":
	default:
		return ErrInvalidAction
	}

	resp, err := w.doRequest(node, http.MethodPost, fmt.Sprintf("/api/servers/%s/power", uuid), PowerPayload{Action: action})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%w: status %d", ErrWingsAPI, resp.StatusCode)
	}
	return nil
}

func (w *WingsClient) GetServerResources(node models.Node, uuid string) (*ServerResources, error) {
	resp, err := w.doRequest(node, http.MethodGet, fmt.Sprintf("/api/servers/%s/resources", uuid), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status %d", ErrWingsAPI, resp.StatusCode)
	}

	var result ServerResources
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (w *WingsClient) GetConsoleToken(node models.Node, uuid string) (*ServerConsoleToken, error) {
	resp, err := w.doRequest(node, http.MethodGet, fmt.Sprintf("/api/servers/%s/ws/console", uuid), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status %d", ErrWingsAPI, resp.StatusCode)
	}

	var result ServerConsoleToken
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
