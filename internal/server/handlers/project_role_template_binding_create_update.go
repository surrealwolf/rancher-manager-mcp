package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterProjectRoleTemplateBindingCreateUpdateTools registers project role template binding create and update tools
func RegisterProjectRoleTemplateBindingCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_project_role_template_binding", "Create a new project role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Project role template binding object to create",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace for the binding",
			},
		},
		"required": []string{"binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createProjectRoleTemplateBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_project_role_template_binding", "Update/replace a project role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the project role template binding",
			},
			"binding": map[string]interface{}{
				"type":        "object",
				"description": "Updated project role template binding object",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the binding",
			},
		},
		"required": []string{"name", "binding"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateProjectRoleTemplateBinding(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_project_role_template_binding", "Partially update a project role template binding", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the project role template binding",
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
		return patchProjectRoleTemplateBinding(ctx, args, rancherClient)
	})
}

func createProjectRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	binding, ok := args["binding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("binding parameter is required and must be an object")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.CreateProjectRoleTemplateBinding(ctx, binding, namespace)
}

func updateProjectRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.UpdateProjectRoleTemplateBinding(ctx, name, binding, namespace)
}

func patchProjectRoleTemplateBinding(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchProjectRoleTemplateBinding(ctx, name, patch, namespace)
}
