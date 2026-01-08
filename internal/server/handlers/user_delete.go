package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterUserDeleteTools registers user delete tools
func RegisterUserDeleteTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("delete_user", "Delete a user", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the user to delete",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return deleteUser(ctx, args, rancherClient)
	})
}

func deleteUser(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.DeleteUser(ctx, name)
}
