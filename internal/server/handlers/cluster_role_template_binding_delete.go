package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterClusterRoleTemplateBindingDeleteTools registers cluster role template binding delete tools
func RegisterClusterRoleTemplateBindingDeleteTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("delete_cluster_role_template_binding", "Delete a cluster role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the cluster role template binding to delete",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the binding",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return deleteClusterRoleTemplateBinding(ctx, args, rancherClient)
	})
}

func deleteClusterRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.DeleteClusterRoleTemplateBinding(ctx, name, namespace)
}
