package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterKubeconfigCreateUpdateTools registers kubeconfig create and update tools
func RegisterKubeconfigCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_kubeconfig", "Create a new kubeconfig", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"kubeconfig": map[string]interface{}{
				"type":        "object",
				"description": "Kubeconfig object to create",
			},
		},
		"required": []string{"kubeconfig"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createKubeconfig(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_kubeconfig", "Update/replace a kubeconfig", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the kubeconfig",
			},
			"kubeconfig": map[string]interface{}{
				"type":        "object",
				"description": "Updated kubeconfig object",
			},
		},
		"required": []string{"name", "kubeconfig"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateKubeconfig(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_kubeconfig", "Partially update a kubeconfig", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the kubeconfig",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchKubeconfig(ctx, args, rancherClient)
	})
}

func createKubeconfig(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	kubeconfig, ok := args["kubeconfig"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("kubeconfig parameter is required and must be an object")
	}
	return rancherClient.CreateKubeconfig(ctx, kubeconfig)
}

func updateKubeconfig(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	kubeconfig, ok := args["kubeconfig"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("kubeconfig parameter is required and must be an object")
	}
	return rancherClient.UpdateKubeconfig(ctx, name, kubeconfig)
}

func patchKubeconfig(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	patch, ok := args["patch"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("patch parameter is required and must be an object")
	}
	return rancherClient.PatchKubeconfig(ctx, name, patch)
}
