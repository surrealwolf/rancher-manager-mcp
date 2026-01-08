package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterProjectStatusTools registers project status tools
func RegisterProjectStatusTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("get_project_status", "Get status of a specific project", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the project",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the project",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getProjectStatus(ctx, args, rancherClient)
	})
}

func getProjectStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetProjectStatus(ctx, name, namespace)
}
