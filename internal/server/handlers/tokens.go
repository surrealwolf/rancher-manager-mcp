package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterTokenTools registers all token management tools
func RegisterTokenTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_tokens", "List all API tokens", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listTokens(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_token", "Get details of a specific API token", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the token",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getToken(ctx, args, rancherClient)
	})
}

func listTokens(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return rancherClient.ListTokens(ctx)
}

func getToken(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetToken(ctx, name)
}
