package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterRoleTemplateCreateUpdateTools registers role template create and update tools
func RegisterRoleTemplateCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_role_template", "Create a new role template", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"template": map[string]interface{}{
				"type":        "object",
				"description": "Role template object to create",
			},
		},
		"required": []string{"template"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createRoleTemplate(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_role_template", "Update/replace a role template", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the role template",
			},
			"template": map[string]interface{}{
				"type":        "object",
				"description": "Updated role template object",
			},
		},
		"required": []string{"name", "template"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateRoleTemplate(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_role_template", "Partially update a role template", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the role template",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchRoleTemplate(ctx, args, rancherClient)
	})
}

func createRoleTemplate(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	template, ok := args["template"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("template parameter is required and must be an object")
	}
	return rancherClient.CreateRoleTemplate(ctx, template)
}

func updateRoleTemplate(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	template, ok := args["template"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("template parameter is required and must be an object")
	}
	return rancherClient.UpdateRoleTemplate(ctx, name, template)
}

func patchRoleTemplate(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchRoleTemplate(ctx, name, patch)
}
