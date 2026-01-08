package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

// TestCRUDOperations tests full CRUD operations with test objects
func TestCRUDOperations(t *testing.T) {
	if testClient == nil {
		t.Skip("Skipping CRUD tests: test client not initialized")
	}

	ctx := context.Background()
	cleanupItems := []cleanupItem{}

	// Cleanup function
	defer func() {
		for _, item := range cleanupItems {
			item.cleanup(ctx, t)
		}
	}()

	// Test Global Role CRUD
	t.Run("GlobalRole", func(t *testing.T) {
		testGlobalRoleCRUD(ctx, t, &cleanupItems)
	})

	// Test Role Template CRUD
	t.Run("RoleTemplate", func(t *testing.T) {
		testRoleTemplateCRUD(ctx, t, &cleanupItems)
	})
}

type cleanupItem struct {
	resourceType string
	name         string
	cleanup      func(ctx context.Context, t *testing.T)
}

func testGlobalRoleCRUD(ctx context.Context, t *testing.T, cleanupItems *[]cleanupItem) {
	testName := fmt.Sprintf("%s-global-role", testPrefix)

	// Create
	t.Run("Create", func(t *testing.T) {
		globalRole := map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "GlobalRole",
			"metadata": map[string]interface{}{
				"name": testName,
			},
			"displayName": "MCP Test Global Role",
			"rules":       []interface{}{},
		}

		result, err := testClient.CreateGlobalRole(ctx, globalRole)
		if err != nil {
			t.Fatalf("CreateGlobalRole failed: %v", err)
		}
		if result == nil {
			t.Fatal("CreateGlobalRole returned nil")
		}

		*cleanupItems = append(*cleanupItems, cleanupItem{
			resourceType: "globalrole",
			name:         testName,
			cleanup: func(ctx context.Context, t *testing.T) {
				_, err := testClient.DeleteGlobalRole(ctx, testName)
				if err != nil {
					t.Logf("Cleanup: DeleteGlobalRole failed: %v", err)
				}
			},
		})

		t.Logf("CreateGlobalRole: created %s", testName)
	})

	// Get
	t.Run("Get", func(t *testing.T) {
		result, err := testClient.GetGlobalRole(ctx, testName)
		if err != nil {
			t.Fatalf("GetGlobalRole failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetGlobalRole returned nil")
		}
		t.Logf("GetGlobalRole: retrieved %s", testName)
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		// Get existing resource to get resourceVersion
		existing, err := testClient.GetGlobalRole(ctx, testName)
		if err != nil {
			t.Fatalf("GetGlobalRole failed before update: %v", err)
		}

		existingData, err := json.Marshal(existing)
		if err != nil {
			t.Fatalf("Failed to marshal existing resource: %v", err)
		}

		var existingObj map[string]interface{}
		if err := json.Unmarshal(existingData, &existingObj); err != nil {
			t.Fatalf("Failed to unmarshal existing resource: %v", err)
		}

		// Update with new displayName
		existingObj["displayName"] = "MCP Test Global Role - Updated"

		result, err := testClient.UpdateGlobalRole(ctx, testName, existingObj)
		if err != nil {
			t.Fatalf("UpdateGlobalRole failed: %v", err)
		}
		if result == nil {
			t.Fatal("UpdateGlobalRole returned nil")
		}
		t.Logf("UpdateGlobalRole: updated %s", testName)
	})

	// Patch
	t.Run("Patch", func(t *testing.T) {
		patch := map[string]interface{}{
			"displayName": "MCP Test Global Role - Patched",
		}

		result, err := testClient.PatchGlobalRole(ctx, testName, patch)
		if err != nil {
			t.Fatalf("PatchGlobalRole failed: %v", err)
		}
		if result == nil {
			t.Fatal("PatchGlobalRole returned nil")
		}
		t.Logf("PatchGlobalRole: patched %s", testName)
	})

	// Status
	t.Run("Status", func(t *testing.T) {
		result, err := testClient.GetGlobalRoleStatus(ctx, testName)
		if err != nil {
			t.Fatalf("GetGlobalRoleStatus failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetGlobalRoleStatus returned nil")
		}
		t.Logf("GetGlobalRoleStatus: retrieved status for %s", testName)
	})
}

func testRoleTemplateCRUD(ctx context.Context, t *testing.T, cleanupItems *[]cleanupItem) {
	testName := fmt.Sprintf("%s-role-template", testPrefix)

	// Create
	t.Run("Create", func(t *testing.T) {
		roleTemplate := map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "RoleTemplate",
			"metadata": map[string]interface{}{
				"name": testName,
			},
			"displayName": "MCP Test Role Template",
			"rules":       []interface{}{},
		}

		result, err := testClient.CreateRoleTemplate(ctx, roleTemplate)
		if err != nil {
			t.Fatalf("CreateRoleTemplate failed: %v", err)
		}
		if result == nil {
			t.Fatal("CreateRoleTemplate returned nil")
		}

		*cleanupItems = append(*cleanupItems, cleanupItem{
			resourceType: "roletemplate",
			name:         testName,
			cleanup: func(ctx context.Context, t *testing.T) {
				_, err := testClient.DeleteRoleTemplate(ctx, testName)
				if err != nil {
					t.Logf("Cleanup: DeleteRoleTemplate failed: %v", err)
				}
			},
		})

		t.Logf("CreateRoleTemplate: created %s", testName)
	})

	// Get
	t.Run("Get", func(t *testing.T) {
		result, err := testClient.GetRoleTemplate(ctx, testName)
		if err != nil {
			t.Fatalf("GetRoleTemplate failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetRoleTemplate returned nil")
		}
		t.Logf("GetRoleTemplate: retrieved %s", testName)
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		// Get existing resource to get resourceVersion
		existing, err := testClient.GetRoleTemplate(ctx, testName)
		if err != nil {
			t.Fatalf("GetRoleTemplate failed before update: %v", err)
		}

		existingData, err := json.Marshal(existing)
		if err != nil {
			t.Fatalf("Failed to marshal existing resource: %v", err)
		}

		var existingObj map[string]interface{}
		if err := json.Unmarshal(existingData, &existingObj); err != nil {
			t.Fatalf("Failed to unmarshal existing resource: %v", err)
		}

		// Update with new displayName
		existingObj["displayName"] = "MCP Test Role Template - Updated"

		result, err := testClient.UpdateRoleTemplate(ctx, testName, existingObj)
		if err != nil {
			t.Fatalf("UpdateRoleTemplate failed: %v", err)
		}
		if result == nil {
			t.Fatal("UpdateRoleTemplate returned nil")
		}
		t.Logf("UpdateRoleTemplate: updated %s", testName)
	})

	// Patch
	t.Run("Patch", func(t *testing.T) {
		patch := map[string]interface{}{
			"displayName": "MCP Test Role Template - Patched",
		}

		result, err := testClient.PatchRoleTemplate(ctx, testName, patch)
		if err != nil {
			t.Fatalf("PatchRoleTemplate failed: %v", err)
		}
		if result == nil {
			t.Fatal("PatchRoleTemplate returned nil")
		}
		t.Logf("PatchRoleTemplate: patched %s", testName)
	})

	// Status
	t.Run("Status", func(t *testing.T) {
		result, err := testClient.GetRoleTemplateStatus(ctx, testName)
		if err != nil {
			t.Fatalf("GetRoleTemplateStatus failed: %v", err)
		}
		if result == nil {
			t.Fatal("GetRoleTemplateStatus returned nil")
		}
		t.Logf("GetRoleTemplateStatus: retrieved status for %s", testName)
	})
}
