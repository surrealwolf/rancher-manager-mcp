package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterKubeconfigDeleteTools registers kubeconfig delete tools
func RegisterKubeconfigDeleteTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("delete_kubeconfig", "Delete a kubeconfig", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the kubeconfig to delete",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return deleteKubeconfig(ctx, args, rancherClient)
	})
}

func deleteKubeconfig(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.DeleteKubeconfig(ctx, name)
}
