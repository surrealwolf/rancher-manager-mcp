package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterCreateUpdateTools registers cluster create and update tools
func RegisterClusterCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_cluster", "Create a new cluster", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"cluster": map[string]interface{}{
				"type":        "object",
				"description": "Cluster object to create",
			},
		},
		"required": []string{"cluster"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createCluster(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_cluster", "Update/replace a cluster", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the cluster",
			},
			"cluster": map[string]interface{}{
				"type":        "object",
				"description": "Updated cluster object",
			},
		},
		"required": []string{"name", "cluster"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateCluster(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_cluster", "Partially update a cluster", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the cluster",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchCluster(ctx, args, rancherClient)
	})
}

func createCluster(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	cluster, ok := args["cluster"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cluster parameter is required and must be an object")
	}
	return rancherClient.CreateCluster(ctx, cluster)
}

func updateCluster(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	cluster, ok := args["cluster"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cluster parameter is required and must be an object")
	}
	return rancherClient.UpdateCluster(ctx, name, cluster)
}

func patchCluster(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchCluster(ctx, name, patch)
}
