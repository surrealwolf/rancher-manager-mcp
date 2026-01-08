package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterGlobalRoleTools registers all global role management tools
func RegisterGlobalRoleTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_global_roles", "List all global roles", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listGlobalRoles(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_global_role", "Get details of a specific global role", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getGlobalRole(ctx, args, rancherClient)
	})
}

func listGlobalRoles(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return rancherClient.ListGlobalRoles(ctx)
}

func getGlobalRole(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetGlobalRole(ctx, name)
}
