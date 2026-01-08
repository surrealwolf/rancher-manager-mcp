package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterGlobalRoleCreateUpdateTools registers global role create and update tools
func RegisterGlobalRoleCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_global_role", "Create a new global role", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"role": map[string]interface{}{
				"type":        "object",
				"description": "Global role object to create",
			},
		},
		"required": []string{"role"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createGlobalRole(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_global_role", "Update/replace a global role", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role",
			},
			"role": map[string]interface{}{
				"type":        "object",
				"description": "Updated global role object",
			},
		},
		"required": []string{"name", "role"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateGlobalRole(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_global_role", "Partially update a global role", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchGlobalRole(ctx, args, rancherClient)
	})
}

func createGlobalRole(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	role, ok := args["role"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("role parameter is required and must be an object")
	}
	return rancherClient.CreateGlobalRole(ctx, role)
}

func updateGlobalRole(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	role, ok := args["role"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("role parameter is required and must be an object")
	}
	return rancherClient.UpdateGlobalRole(ctx, name, role)
}

func patchGlobalRole(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	patch, ok := args["patch"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("patch parameter is required and must be an object")
	}
	return rancherClient.PatchGlobalRole(ctx, name, patch)
}
