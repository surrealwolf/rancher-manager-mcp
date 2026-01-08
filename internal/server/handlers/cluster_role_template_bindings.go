package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterRoleTemplateBindingTools registers all cluster role template binding management tools
func RegisterClusterRoleTemplateBindingTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_cluster_role_template_bindings", "List all cluster role template bindings", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace to filter cluster role template bindings",
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listClusterRoleTemplateBindings(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_cluster_role_template_binding", "Get details of a specific cluster role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the cluster role template binding",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the cluster role template binding",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getClusterRoleTemplateBinding(ctx, args, rancherClient)
	})
}

func listClusterRoleTemplateBindings(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.ListNamespacedClusterRoleTemplateBindings(ctx, namespace)
}

func getClusterRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetClusterRoleTemplateBinding(ctx, name, namespace)
}
