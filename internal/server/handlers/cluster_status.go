package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterStatusTools registers cluster status tools
func RegisterClusterStatusTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("get_cluster_status", "Get status of a specific cluster", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the cluster",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getClusterStatus(ctx, args, rancherClient)
	})
}

func getClusterStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetClusterStatus(ctx, name)
}
