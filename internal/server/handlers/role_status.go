package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterRoleStatusTools registers role status tools
func RegisterRoleStatusTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	// GlobalRole status
	mcpServer.RegisterToolWithSchema("get_global_role_status", "Get status of a specific global role", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getGlobalRoleStatus(ctx, args, rancherClient)
	})

	// GlobalRoleBinding status
	mcpServer.RegisterToolWithSchema("get_global_role_binding_status", "Get status of a specific global role binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the global role binding",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getGlobalRoleBindingStatus(ctx, args, rancherClient)
	})

	// RoleTemplate status
	mcpServer.RegisterToolWithSchema("get_role_template_status", "Get status of a specific role template", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the role template",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getRoleTemplateStatus(ctx, args, rancherClient)
	})
}

func getGlobalRoleStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetGlobalRoleStatus(ctx, name)
}

func getGlobalRoleBindingStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetGlobalRoleBindingStatus(ctx, name)
}

func getRoleTemplateStatus(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetRoleTemplateStatus(ctx, name)
}
