package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterTools registers all cluster management tools
func RegisterClusterTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_clusters", "List all Rancher clusters", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listClusters(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_cluster", "Get details of a specific cluster", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the cluster",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getCluster(ctx, args, rancherClient)
	})
}

func listClusters(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return rancherClient.ListClusters(ctx)
}

func getCluster(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetCluster(ctx, name)
}
