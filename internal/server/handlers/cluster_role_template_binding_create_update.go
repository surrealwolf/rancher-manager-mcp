package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterRoleTemplateBindingCreateUpdateTools registers cluster role template binding create and update tools
func RegisterClusterRoleTemplateBindingCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_cluster_role_template_binding", "Create a new cluster role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Cluster role template binding object to create",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace for the binding",
			},
		},
		"required": []string{"binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createClusterRoleTemplateBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_cluster_role_template_binding", "Update/replace a cluster role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the cluster role template binding",
			},
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Updated cluster role template binding object",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the binding",
			},
		},
		"required": []string{"name", "binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateClusterRoleTemplateBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_cluster_role_template_binding", "Partially update a cluster role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the cluster role template binding",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the binding",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchClusterRoleTemplateBinding(ctx, args, rancherClient)
	})
}

func createClusterRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	binding, ok := args["binding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("binding parameter is required and must be an object")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.CreateClusterRoleTemplateBinding(ctx, binding, namespace)
}

func updateClusterRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	binding, ok := args["binding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("binding parameter is required and must be an object")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.UpdateClusterRoleTemplateBinding(ctx, name, binding, namespace)
}

func patchClusterRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	namespace, _ := args["namespace"].(string)
	return rancherClient.PatchClusterRoleTemplateBinding(ctx, name, patch, namespace)
}
