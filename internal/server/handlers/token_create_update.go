package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterTokenCreateUpdateTools registers token create and update tools
func RegisterTokenCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_token", "Create a new API token", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"token": map[string]interface{}{
				"type":        "object",
				"description": "Token object to create",
			},
		},
		"required": []string{"token"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createToken(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_token", "Update/replace an API token", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the token",
			},
			"token": map[string]interface{}{
				"type":        "object",
				"description": "Updated token object",
			},
		},
		"required": []string{"name", "token"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateToken(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_token", "Partially update an API token", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the token",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchToken(ctx, args, rancherClient)
	})
}

func createToken(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	token, ok := args["token"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("token parameter is required and must be an object")
	}
	return rancherClient.CreateToken(ctx, token)
}

func updateToken(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	token, ok := args["token"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("token parameter is required and must be an object")
	}
	return rancherClient.UpdateToken(ctx, name, token)
}

func patchToken(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchToken(ctx, name, patch)
}
