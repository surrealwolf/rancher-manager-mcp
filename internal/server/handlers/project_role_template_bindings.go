package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterProjectRoleTemplateBindingTools registers all project role template binding management tools
func RegisterProjectRoleTemplateBindingTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_project_role_template_bindings", "List all project role template bindings", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace to filter project role template bindings",
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listProjectRoleTemplateBindings(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_project_role_template_binding", "Get details of a specific project role template binding", map[string]interface{}{
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
		return getProjectRoleTemplateBinding(ctx, args, rancherClient)
	})
}

func listProjectRoleTemplateBindings(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.ListNamespacedProjectRoleTemplateBindings(ctx, namespace)
}

func getProjectRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.GetProjectRoleTemplateBinding(ctx, name, namespace)
}
