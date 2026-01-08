package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	// Normalize baseURL: remove trailing slash if present
	normalizedBaseURL := strings.TrimSuffix(baseURL, "/")

	return &RancherClient{
		baseURL: normalizedBaseURL,
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

// Generic list/get helpers
func (c *RancherClient) listResource(ctx context.Context, apiPath string) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", apiPath, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

func (c *RancherClient) getResource(ctx context.Context, apiPath string) (interface{}, error) {
	data, err := c.doRequest(ctx, "GET", apiPath, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

// auditlogCattleIo_v1 - AuditPolicy
func (c *RancherClient) ListAuditPolicies(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/auditlog.cattle.io/v1/auditpolicies")
}

func (c *RancherClient) GetAuditPolicy(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/auditlog.cattle.io/v1/auditpolicies/%s", name))
}

// extCattleIo_v1 - Kubeconfig
func (c *RancherClient) ListKubeconfigs(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/ext.cattle.io/v1/kubeconfigs")
}

func (c *RancherClient) GetKubeconfig(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/kubeconfigs/%s", name))
}

// extCattleIo_v1 - Token
func (c *RancherClient) ListTokens(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/ext.cattle.io/v1/tokens")
}

func (c *RancherClient) GetToken(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/tokens/%s", name))
}

// managementCattleIo_v3 - GlobalRole
func (c *RancherClient) ListGlobalRoles(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/globalroles")
}

func (c *RancherClient) GetGlobalRole(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalroles/%s", name))
}

// managementCattleIo_v3 - GlobalRoleBinding
func (c *RancherClient) ListGlobalRoleBindings(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/globalrolebindings")
}

func (c *RancherClient) GetGlobalRoleBinding(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalrolebindings/%s", name))
}

// managementCattleIo_v3 - RoleTemplate
func (c *RancherClient) ListRoleTemplates(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/roletemplates")
}

func (c *RancherClient) GetRoleTemplate(ctx context.Context, name string) (interface{}, error) {
	return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/roletemplates/%s", name))
}

// managementCattleIo_v3 - ClusterRoleTemplateBinding (all namespaces)
func (c *RancherClient) ListClusterRoleTemplateBindings(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/clusterroletemplatebindings")
}

// managementCattleIo_v3 - ClusterRoleTemplateBinding (namespaced)
func (c *RancherClient) ListNamespacedClusterRoleTemplateBindings(ctx context.Context, namespace string) (interface{}, error) {
	if namespace == "" {
		return c.ListClusterRoleTemplateBindings(ctx)
	}
	return c.listResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings", namespace))
}

func (c *RancherClient) GetClusterRoleTemplateBinding(ctx context.Context, name, namespace string) (interface{}, error) {
	if namespace == "" {
		return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusterroletemplatebindings/%s", name))
	}
	return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings/%s", namespace, name))
}

// managementCattleIo_v3 - ProjectRoleTemplateBinding (all namespaces)
func (c *RancherClient) ListProjectRoleTemplateBindings(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/projectroletemplatebindings")
}

// managementCattleIo_v3 - ProjectRoleTemplateBinding (namespaced)
func (c *RancherClient) ListNamespacedProjectRoleTemplateBindings(ctx context.Context, namespace string) (interface{}, error) {
	if namespace == "" {
		return c.ListProjectRoleTemplateBindings(ctx)
	}
	return c.listResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings", namespace))
}

func (c *RancherClient) GetProjectRoleTemplateBinding(ctx context.Context, name, namespace string) (interface{}, error) {
	if namespace == "" {
		return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projectroletemplatebindings/%s", name))
	}
	return c.getResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings/%s", namespace, name))
}

// managementCattleIo_v3 - Project (all namespaces)
func (c *RancherClient) ListProjectsAllNamespaces(ctx context.Context) (interface{}, error) {
	return c.listResource(ctx, "/apis/management.cattle.io/v3/projects")
}

// managementCattleIo_v3 - Project (namespaced)
func (c *RancherClient) ListNamespacedProjects(ctx context.Context, namespace string) (interface{}, error) {
	if namespace == "" {
		return c.ListProjectsAllNamespaces(ctx)
	}
	return c.listResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects", namespace))
}

// Status subresources
// Note: Rancher API doesn't expose status as a subresource endpoint.
// The status is part of the main resource object, so we get the resource and extract the status field.

// getStatusFromResource extracts the status field from a resource object
func (c *RancherClient) getStatusFromResource(ctx context.Context, getFunc func() (interface{}, error)) (interface{}, error) {
	resource, err := getFunc()
	if err != nil {
		return nil, err
	}

	// Convert to map to extract status field
	resourceMap, ok := resource.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert resource to map")
	}

	status, exists := resourceMap["status"]
	if !exists {
		return nil, fmt.Errorf("status field not found in resource")
	}

	return status, nil
}

// GetClusterStatus gets the status of a specific cluster
func (c *RancherClient) GetClusterStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetCluster(ctx, name)
	})
}

// GetUserStatus gets the status of a specific user
func (c *RancherClient) GetUserStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetUser(ctx, name)
	})
}

// GetProjectStatus gets the status of a specific project
func (c *RancherClient) GetProjectStatus(ctx context.Context, name, namespace string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetProject(ctx, name, namespace)
	})
}

// GetGlobalRoleStatus gets the status of a specific global role
func (c *RancherClient) GetGlobalRoleStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetGlobalRole(ctx, name)
	})
}

// GetGlobalRoleBindingStatus gets the status of a specific global role binding
func (c *RancherClient) GetGlobalRoleBindingStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetGlobalRoleBinding(ctx, name)
	})
}

// GetRoleTemplateStatus gets the status of a specific role template
func (c *RancherClient) GetRoleTemplateStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetRoleTemplate(ctx, name)
	})
}

// GetClusterRoleTemplateBindingStatus gets the status of a specific cluster role template binding
func (c *RancherClient) GetClusterRoleTemplateBindingStatus(ctx context.Context, name, namespace string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetClusterRoleTemplateBinding(ctx, name, namespace)
	})
}

// GetProjectRoleTemplateBindingStatus gets the status of a specific project role template binding
func (c *RancherClient) GetProjectRoleTemplateBindingStatus(ctx context.Context, name, namespace string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetProjectRoleTemplateBinding(ctx, name, namespace)
	})
}

// GetAuditPolicyStatus gets the status of a specific audit policy
func (c *RancherClient) GetAuditPolicyStatus(ctx context.Context, name string) (interface{}, error) {
	return c.getStatusFromResource(ctx, func() (interface{}, error) {
		return c.GetAuditPolicy(ctx, name)
	})
}

// Generic create/update helpers
func (c *RancherClient) createResource(ctx context.Context, apiPath string, body interface{}) (interface{}, error) {
	data, err := c.doRequest(ctx, "POST", apiPath, body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

func (c *RancherClient) updateResource(ctx context.Context, apiPath string, body interface{}) (interface{}, error) {
	data, err := c.doRequest(ctx, "PUT", apiPath, body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

func (c *RancherClient) patchResource(ctx context.Context, apiPath string, body interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, apiPath)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/merge-patch+json")
	req.Header.Set("Accept", "application/json")

	logrus.Debugf("Making request: PATCH %s", url)
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

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

// Create methods

// CreateCluster creates a new cluster
func (c *RancherClient) CreateCluster(ctx context.Context, cluster map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/management.cattle.io/v3/clusters", cluster)
}

// CreateUser creates a new user
func (c *RancherClient) CreateUser(ctx context.Context, user map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/management.cattle.io/v3/users", user)
}

// CreateProject creates a new project
func (c *RancherClient) CreateProject(ctx context.Context, project map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.createResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects", namespace), project)
	}
	return c.createResource(ctx, "/apis/management.cattle.io/v3/projects", project)
}

// CreateAuditPolicy creates a new audit policy
func (c *RancherClient) CreateAuditPolicy(ctx context.Context, policy map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/auditlog.cattle.io/v1/auditpolicies", policy)
}

// CreateKubeconfig creates a new kubeconfig
func (c *RancherClient) CreateKubeconfig(ctx context.Context, kubeconfig map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/ext.cattle.io/v1/kubeconfigs", kubeconfig)
}

// CreateToken creates a new token
func (c *RancherClient) CreateToken(ctx context.Context, token map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/ext.cattle.io/v1/tokens", token)
}

// CreateGlobalRole creates a new global role
func (c *RancherClient) CreateGlobalRole(ctx context.Context, role map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/management.cattle.io/v3/globalroles", role)
}

// CreateGlobalRoleBinding creates a new global role binding
func (c *RancherClient) CreateGlobalRoleBinding(ctx context.Context, binding map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/management.cattle.io/v3/globalrolebindings", binding)
}

// CreateRoleTemplate creates a new role template
func (c *RancherClient) CreateRoleTemplate(ctx context.Context, template map[string]interface{}) (interface{}, error) {
	return c.createResource(ctx, "/apis/management.cattle.io/v3/roletemplates", template)
}

// CreateClusterRoleTemplateBinding creates a new cluster role template binding
func (c *RancherClient) CreateClusterRoleTemplateBinding(ctx context.Context, binding map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.createResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings", namespace), binding)
	}
	return c.createResource(ctx, "/apis/management.cattle.io/v3/clusterroletemplatebindings", binding)
}

// CreateProjectRoleTemplateBinding creates a new project role template binding
func (c *RancherClient) CreateProjectRoleTemplateBinding(ctx context.Context, binding map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.createResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings", namespace), binding)
	}
	return c.createResource(ctx, "/apis/management.cattle.io/v3/projectroletemplatebindings", binding)
}

// Update methods (PUT - replace)

// UpdateCluster updates/replaces a cluster
func (c *RancherClient) UpdateCluster(ctx context.Context, name string, cluster map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusters/%s", name), cluster)
}

// UpdateUser updates/replaces a user
func (c *RancherClient) UpdateUser(ctx context.Context, name string, user map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/users/%s", name), user)
}

// UpdateProject updates/replaces a project
func (c *RancherClient) UpdateProject(ctx context.Context, name string, project map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects/%s", namespace, name), project)
	}
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projects/%s", name), project)
}

// UpdateAuditPolicy updates/replaces an audit policy
func (c *RancherClient) UpdateAuditPolicy(ctx context.Context, name string, policy map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/auditlog.cattle.io/v1/auditpolicies/%s", name), policy)
}

// UpdateKubeconfig updates/replaces a kubeconfig
func (c *RancherClient) UpdateKubeconfig(ctx context.Context, name string, kubeconfig map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/kubeconfigs/%s", name), kubeconfig)
}

// UpdateToken updates/replaces a token
func (c *RancherClient) UpdateToken(ctx context.Context, name string, token map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/tokens/%s", name), token)
}

// UpdateGlobalRole updates/replaces a global role
func (c *RancherClient) UpdateGlobalRole(ctx context.Context, name string, role map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalroles/%s", name), role)
}

// UpdateGlobalRoleBinding updates/replaces a global role binding
func (c *RancherClient) UpdateGlobalRoleBinding(ctx context.Context, name string, binding map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalrolebindings/%s", name), binding)
}

// UpdateRoleTemplate updates/replaces a role template
func (c *RancherClient) UpdateRoleTemplate(ctx context.Context, name string, template map[string]interface{}) (interface{}, error) {
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/roletemplates/%s", name), template)
}

// UpdateClusterRoleTemplateBinding updates/replaces a cluster role template binding
func (c *RancherClient) UpdateClusterRoleTemplateBinding(ctx context.Context, name string, binding map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings/%s", namespace, name), binding)
	}
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusterroletemplatebindings/%s", name), binding)
}

// UpdateProjectRoleTemplateBinding updates/replaces a project role template binding
func (c *RancherClient) UpdateProjectRoleTemplateBinding(ctx context.Context, name string, binding map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings/%s", namespace, name), binding)
	}
	return c.updateResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projectroletemplatebindings/%s", name), binding)
}

// Patch methods (PATCH - partial update)

// PatchCluster partially updates a cluster
func (c *RancherClient) PatchCluster(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusters/%s", name), patch)
}

// PatchUser partially updates a user
func (c *RancherClient) PatchUser(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/users/%s", name), patch)
}

// PatchProject partially updates a project
func (c *RancherClient) PatchProject(ctx context.Context, name string, patch map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects/%s", namespace, name), patch)
	}
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projects/%s", name), patch)
}

// PatchAuditPolicy partially updates an audit policy
func (c *RancherClient) PatchAuditPolicy(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/auditlog.cattle.io/v1/auditpolicies/%s", name), patch)
}

// PatchKubeconfig partially updates a kubeconfig
func (c *RancherClient) PatchKubeconfig(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/kubeconfigs/%s", name), patch)
}

// PatchToken partially updates a token
func (c *RancherClient) PatchToken(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/tokens/%s", name), patch)
}

// PatchGlobalRole partially updates a global role
func (c *RancherClient) PatchGlobalRole(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalroles/%s", name), patch)
}

// PatchGlobalRoleBinding partially updates a global role binding
func (c *RancherClient) PatchGlobalRoleBinding(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalrolebindings/%s", name), patch)
}

// PatchRoleTemplate partially updates a role template
func (c *RancherClient) PatchRoleTemplate(ctx context.Context, name string, patch map[string]interface{}) (interface{}, error) {
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/roletemplates/%s", name), patch)
}

// PatchClusterRoleTemplateBinding partially updates a cluster role template binding
func (c *RancherClient) PatchClusterRoleTemplateBinding(ctx context.Context, name string, patch map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings/%s", namespace, name), patch)
	}
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusterroletemplatebindings/%s", name), patch)
}

// PatchProjectRoleTemplateBinding partially updates a project role template binding
func (c *RancherClient) PatchProjectRoleTemplateBinding(ctx context.Context, name string, patch map[string]interface{}, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings/%s", namespace, name), patch)
	}
	return c.patchResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projectroletemplatebindings/%s", name), patch)
}

// Generic delete helper
func (c *RancherClient) deleteResource(ctx context.Context, apiPath string) (interface{}, error) {
	data, err := c.doRequest(ctx, "DELETE", apiPath, nil)
	if err != nil {
		return nil, err
	}
	// DELETE requests might return empty body or the deleted resource
	if len(data) == 0 {
		return map[string]interface{}{"deleted": true}, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return result, nil
}

// Delete methods

// DeleteCluster deletes a cluster
func (c *RancherClient) DeleteCluster(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusters/%s", name))
}

// DeleteUser deletes a user
func (c *RancherClient) DeleteUser(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/users/%s", name))
}

// DeleteProject deletes a project
func (c *RancherClient) DeleteProject(ctx context.Context, name string, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projects/%s", namespace, name))
	}
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projects/%s", name))
}

// DeleteAuditPolicy deletes an audit policy
func (c *RancherClient) DeleteAuditPolicy(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/auditlog.cattle.io/v1/auditpolicies/%s", name))
}

// DeleteKubeconfig deletes a kubeconfig
func (c *RancherClient) DeleteKubeconfig(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/kubeconfigs/%s", name))
}

// DeleteToken deletes a token
func (c *RancherClient) DeleteToken(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/ext.cattle.io/v1/tokens/%s", name))
}

// DeleteGlobalRole deletes a global role
func (c *RancherClient) DeleteGlobalRole(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalroles/%s", name))
}

// DeleteGlobalRoleBinding deletes a global role binding
func (c *RancherClient) DeleteGlobalRoleBinding(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/globalrolebindings/%s", name))
}

// DeleteRoleTemplate deletes a role template
func (c *RancherClient) DeleteRoleTemplate(ctx context.Context, name string) (interface{}, error) {
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/roletemplates/%s", name))
}

// DeleteClusterRoleTemplateBinding deletes a cluster role template binding
func (c *RancherClient) DeleteClusterRoleTemplateBinding(ctx context.Context, name string, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/clusterroletemplatebindings/%s", namespace, name))
	}
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/clusterroletemplatebindings/%s", name))
}

// DeleteProjectRoleTemplateBinding deletes a project role template binding
func (c *RancherClient) DeleteProjectRoleTemplateBinding(ctx context.Context, name string, namespace string) (interface{}, error) {
	if namespace != "" {
		return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/namespaces/%s/projectroletemplatebindings/%s", namespace, name))
	}
	return c.deleteResource(ctx, fmt.Sprintf("/apis/management.cattle.io/v3/projectroletemplatebindings/%s", name))
}
