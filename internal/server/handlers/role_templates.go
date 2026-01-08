package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/mcp"
	"github.com/rancher/rancher-manager-mcp/internal/client"
)

// RegisterRoleTemplateTools registers all role template management tools
func RegisterRoleTemplateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("list_role_templates", "List all role templates", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return listRoleTemplates(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("get_role_template", "Get details of a specific role template", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the role template",
			},
		},
		"required": []string{"name"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return getRoleTemplate(ctx, args, rancherClient)
	})
}

func listRoleTemplates(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	return rancherClient.ListRoleTemplates(ctx)
}

func getRoleTemplate(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	return rancherClient.GetRoleTemplate(ctx, name)
}
