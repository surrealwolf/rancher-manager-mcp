package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterGlobalRoleBindingCreateUpdateTools registers global role binding create and update tools
func RegisterGlobalRoleBindingCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_global_role_binding", "Create a new global role binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Global role binding object to create",
			},
		},
		"required": []string{"binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createGlobalRoleBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_global_role_binding", "Update/replace a global role binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role binding",
			},
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Updated global role binding object",
			},
		},
		"required": []string{"name", "binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateGlobalRoleBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_global_role_binding", "Partially update a global role binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role binding",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchGlobalRoleBinding(ctx, args, rancherClient)
	})
}

func createGlobalRoleBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	binding, ok := args["binding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("binding parameter is required and must be an object")
	}
	return rancherClient.CreateGlobalRoleBinding(ctx, binding)
}

func updateGlobalRoleBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	binding, ok := args["binding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("binding parameter is required and must be an object")
	}
	return rancherClient.UpdateGlobalRoleBinding(ctx, name, binding)
}

func patchGlobalRoleBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchGlobalRoleBinding(ctx, name, patch)
}
