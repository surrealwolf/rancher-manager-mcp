package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterProjectTools registers all project management tools
func RegisterProjectTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_projects", "List all Rancher projects", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace to filter projects",
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listProjects(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_project", "Get details of a specific project", map[string]interface{}{
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
		return getProject(ctx, args, rancherClient)
	})
}

func listProjects(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.ListNamespacedProjects(ctx, namespace)
}

func getProject(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetProject(ctx, name, namespace)
}
