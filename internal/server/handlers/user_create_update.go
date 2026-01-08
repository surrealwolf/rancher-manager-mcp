package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterUserCreateUpdateTools registers user create and update tools
func RegisterUserCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_user", "Create a new user", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"user": map[string]interface{}{
				"type":        "object",
				"description": "User object to create",
			},
		},
		"required": []string{"user"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createUser(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_user", "Update/replace a user", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the user",
			},
			"user": map[string]interface{}{
				"type":        "object",
				"description": "Updated user object",
			},
		},
		"required": []string{"name", "user"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateUser(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_user", "Partially update a user", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the user",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchUser(ctx, args, rancherClient)
	})
}

func createUser(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	user, ok := args["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user parameter is required and must be an object")
	}
	return rancherClient.CreateUser(ctx, user)
}

func updateUser(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	user, ok := args["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user parameter is required and must be an object")
	}
	return rancherClient.UpdateUser(ctx, name, user)
}

func patchUser(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchUser(ctx, name, patch)
}
