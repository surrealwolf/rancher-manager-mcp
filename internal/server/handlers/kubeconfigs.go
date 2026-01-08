package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterKubeconfigTools registers all kubeconfig management tools
func RegisterKubeconfigTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_kubeconfigs", "List all kubeconfigs", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listKubeconfigs(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_kubeconfig", "Get details of a specific kubeconfig", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the kubeconfig",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getKubeconfig(ctx, args, rancherClient)
	})
}

func listKubeconfigs(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return rancherClient.ListKubeconfigs(ctx)
}

func getKubeconfig(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetKubeconfig(ctx, name)
}
