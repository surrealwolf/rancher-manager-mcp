package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	testClient *RancherClient
	testPrefix string
)

func TestMain(m *testing.M) {
	rancherURL := os.Getenv("RANCHER_URL")
	rancherToken := os.Getenv("RANCHER_TOKEN")
	insecureSkipVerify := os.Getenv("RANCHER_INSECURE_SKIP_VERIFY") == "true" || os.Getenv("RANCHER_INSECURE_SKIP_VERIFY") == "1"

	if rancherURL == "" || rancherToken == "" {
		fmt.Println("Skipping integration tests: RANCHER_URL and RANCHER_TOKEN must be set")
		os.Exit(0)
	}

	testClient = NewRancherClient(rancherURL, rancherToken, insecureSkipVerify)
	testPrefix = fmt.Sprintf("mcp-test-%d", time.Now().Unix())

	// Verify connection
	ctx := context.Background()
	if err := testClient.VerifyToken(ctx); err != nil {
		fmt.Printf("Failed to verify token: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestListClusters(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListClusters(ctx)
	if err != nil {
		t.Fatalf("ListClusters failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListClusters returned nil")
	}
	t.Logf("ListClusters: found clusters")
}

func TestGetCluster(t *testing.T) {
	ctx := context.Background()
	clusters, err := testClient.ListClusters(ctx)
	if err != nil {
		t.Skipf("Skipping GetCluster: ListClusters failed: %v", err)
	}

	clusterName := extractFirstItemName(clusters)
	if clusterName == "" {
		t.Skip("Skipping GetCluster: No clusters found")
	}

	result, err := testClient.GetCluster(ctx, clusterName)
	if err != nil {
		t.Fatalf("GetCluster failed: %v", err)
	}
	if result == nil {
		t.Fatal("GetCluster returned nil")
	}
	t.Logf("GetCluster: retrieved cluster %s", clusterName)
}

func TestGetClusterStatus(t *testing.T) {
	ctx := context.Background()
	clusters, err := testClient.ListClusters(ctx)
	if err != nil {
		t.Skipf("Skipping GetClusterStatus: ListClusters failed: %v", err)
	}

	clusterName := extractFirstItemName(clusters)
	if clusterName == "" {
		t.Skip("Skipping GetClusterStatus: No clusters found")
	}

	result, err := testClient.GetClusterStatus(ctx, clusterName)
	if err != nil {
		t.Fatalf("GetClusterStatus failed: %v", err)
	}
	if result == nil {
		t.Fatal("GetClusterStatus returned nil")
	}
	t.Logf("GetClusterStatus: retrieved status for cluster %s", clusterName)
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListUsers returned nil")
	}
	t.Logf("ListUsers: found users")
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	users, err := testClient.ListUsers(ctx)
	if err != nil {
		t.Skipf("Skipping GetUser: ListUsers failed: %v", err)
	}

	userName := extractFirstItemName(users)
	if userName == "" {
		t.Skip("Skipping GetUser: No users found")
	}

	result, err := testClient.GetUser(ctx, userName)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if result == nil {
		t.Fatal("GetUser returned nil")
	}
	t.Logf("GetUser: retrieved user %s", userName)
}

func TestGetUserStatus(t *testing.T) {
	ctx := context.Background()
	users, err := testClient.ListUsers(ctx)
	if err != nil {
		t.Skipf("Skipping GetUserStatus: ListUsers failed: %v", err)
	}

	userName := extractFirstItemName(users)
	if userName == "" {
		t.Skip("Skipping GetUserStatus: No users found")
	}

	result, err := testClient.GetUserStatus(ctx, userName)
	if err != nil {
		t.Fatalf("GetUserStatus failed: %v", err)
	}
	if result == nil {
		t.Fatal("GetUserStatus returned nil")
	}
	t.Logf("GetUserStatus: retrieved status for user %s", userName)
}

func TestListProjects(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListProjects(ctx)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListProjects returned nil")
	}
	t.Logf("ListProjects: found projects")
}

func TestListRoleTemplates(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListRoleTemplates(ctx)
	if err != nil {
		t.Fatalf("ListRoleTemplates failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListRoleTemplates returned nil")
	}
	t.Logf("ListRoleTemplates: found role templates")
}

func TestListGlobalRoles(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListGlobalRoles(ctx)
	if err != nil {
		t.Fatalf("ListGlobalRoles failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListGlobalRoles returned nil")
	}
	t.Logf("ListGlobalRoles: found global roles")
}

func TestListGlobalRoleBindings(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListGlobalRoleBindings(ctx)
	if err != nil {
		t.Fatalf("ListGlobalRoleBindings failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListGlobalRoleBindings returned nil")
	}
	t.Logf("ListGlobalRoleBindings: found global role bindings")
}

func TestListClusterRoleTemplateBindings(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListClusterRoleTemplateBindings(ctx)
	if err != nil {
		t.Fatalf("ListClusterRoleTemplateBindings failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListClusterRoleTemplateBindings returned nil")
	}
	t.Logf("ListClusterRoleTemplateBindings: found cluster role template bindings")
}

func TestListProjectRoleTemplateBindings(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListProjectRoleTemplateBindings(ctx)
	if err != nil {
		t.Fatalf("ListProjectRoleTemplateBindings failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListProjectRoleTemplateBindings returned nil")
	}
	t.Logf("ListProjectRoleTemplateBindings: found project role template bindings")
}

func TestListTokens(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListTokens(ctx)
	if err != nil {
		t.Fatalf("ListTokens failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListTokens returned nil")
	}
	t.Logf("ListTokens: found tokens")
}

func TestListKubeconfigs(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListKubeconfigs(ctx)
	if err != nil {
		t.Fatalf("ListKubeconfigs failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListKubeconfigs returned nil")
	}
	t.Logf("ListKubeconfigs: found kubeconfigs")
}

func TestListAuditPolicies(t *testing.T) {
	ctx := context.Background()
	result, err := testClient.ListAuditPolicies(ctx)
	if err != nil {
		t.Fatalf("ListAuditPolicies failed: %v", err)
	}
	if result == nil {
		t.Fatal("ListAuditPolicies returned nil")
	}
	t.Logf("ListAuditPolicies: found audit policies")
}

// Helper function to extract first item name from list response
func extractFirstItemName(result interface{}) string {
	data, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return ""
	}

	// Try items array first
	if items, ok := obj["items"].([]interface{}); ok && len(items) > 0 {
		if item, ok := items[0].(map[string]interface{}); ok {
			if metadata, ok := item["metadata"].(map[string]interface{}); ok {
				if name, ok := metadata["name"].(string); ok {
					return name
				}
			}
		}
	}

	// Try data array
	if data, ok := obj["data"].([]interface{}); ok && len(data) > 0 {
		if item, ok := data[0].(map[string]interface{}); ok {
			if metadata, ok := item["metadata"].(map[string]interface{}); ok {
				if name, ok := metadata["name"].(string); ok {
					return name
				}
			}
		}
	}

	return ""
}
