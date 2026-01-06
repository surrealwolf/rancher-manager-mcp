package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type RancherClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewRancherClient(baseURL, token string, insecureSkipVerify bool) *RancherClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	return &RancherClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

func (c *RancherClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logrus.Debugf("Making request: %s %s", method, url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ListClusters lists all clusters
func (c *RancherClient) ListClusters(ctx context.Context) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", "/apis/management.cattle.io/v3/clusters", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// GetCluster gets a specific cluster
func (c *RancherClient) GetCluster(ctx context.Context, name string) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/apis/management.cattle.io/v3/clusters/%s", name), nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// ListUsers lists all users
func (c *RancherClient) ListUsers(ctx context.Context) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", "/apis/management.cattle.io/v3/users", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// GetUser gets a specific user
func (c *RancherClient) GetUser(ctx context.Context, name string) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/apis/management.cattle.io/v3/users/%s", name), nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// ListProjects lists all projects
func (c *RancherClient) ListProjects(ctx context.Context) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", "/apis/management.cattle.io/v3/projects", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// GetProject gets a specific project
func (c *RancherClient) GetProject(ctx context.Context, name, namespace string) (interface{}, error) {
	path := fmt.Sprintf("/apis/management.cattle.io/v3/projects/%s", name)
	if namespace != "" {
		path = fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects/%s", namespace, name)
	}
	
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// VerifyToken verifies the Rancher API token
func (c *RancherClient) VerifyToken(ctx context.Context) error {
	// Try to get user info or list clusters to verify token
	_, err := c.doRequest(ctx, "GET", "/apis/management.cattle.io/v3/users", nil)
	if err != nil {
		return fmt.Errorf("token verification failed: %w", err)
	}
	return nil
}
