package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterBindingStatusTools registers binding status tools
func RegisterBindingStatusTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	// ClusterRoleTemplateBinding status
	mcpServer.RegisterToolWithSchema("get_cluster_role_template_binding_status", "Get status of a specific cluster role template binding", map[string]interface{}{
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
		return getClusterRoleTemplateBindingStatus(ctx, args, rancherClient)
	})

	// ProjectRoleTemplateBinding status
	mcpServer.RegisterToolWithSchema("get_project_role_template_binding_status", "Get status of a specific project role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the project role template binding",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the project role template binding",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getProjectRoleTemplateBindingStatus(ctx, args, rancherClient)
	})
}

func getClusterRoleTemplateBindingStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetClusterRoleTemplateBindingStatus(ctx, name, namespace)
}

func getProjectRoleTemplateBindingStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetProjectRoleTemplateBindingStatus(ctx, name, namespace)
}
