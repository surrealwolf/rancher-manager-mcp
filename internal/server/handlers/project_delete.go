package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterProjectDeleteTools registers project delete tools
func RegisterProjectDeleteTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("delete_project", "Delete a project", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the project to delete",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the project",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return deleteProject(ctx, args, rancherClient)
	})
}

func deleteProject(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.DeleteProject(ctx, name, namespace)
}
