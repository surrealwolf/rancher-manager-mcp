package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterProjectCreateUpdateTools registers project create and update tools
func RegisterProjectCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_project", "Create a new project", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project": map[string]interface{}{
				"type":        "object",
				"description": "Project object to create",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace for the project",
			},
		},
		"required": []string{"project"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createProject(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_project", "Update/replace a project", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the project",
			},
			"project": map[string]interface{}{
				"type":        "object",
				"description": "Updated project object",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the project",
			},
		},
		"required": []string{"name", "project"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateProject(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_project", "Partially update a project", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name or ID of the project",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
			"namespace": map[string]interface{}{
				"type":        "string",
				"description": "Optional namespace of the project",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchProject(ctx, args, rancherClient)
	})
}

func createProject(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	project, ok := args["project"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("project parameter is required and must be an object")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.CreateProject(ctx, project, namespace)
}

func updateProject(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	project, ok := args["project"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("project parameter is required and must be an object")
	}
	namespace, _ := args["namespace"].(string)
	return rancherClient.UpdateProject(ctx, name, project, namespace)
}

func patchProject(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchProject(ctx, name, patch, namespace)
}
