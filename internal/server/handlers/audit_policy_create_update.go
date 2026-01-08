package handlers

import (
	"context"
	"fmt"

	"github.com/rancher/rancher-manager-mcp/internal/client"
	"github.com/rancher/rancher-manager-mcp/internal/mcp"
)

// RegisterAuditPolicyCreateUpdateTools registers audit policy create and update tools
func RegisterAuditPolicyCreateUpdateTools(mcpServer *mcp.Server, rancherClient *client.RancherClient) {
	mcpServer.RegisterToolWithSchema("create_audit_policy", "Create a new audit policy", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"policy": map[string]interface{}{
				"type":        "object",
				"description": "Audit policy object to create",
			},
		},
		"required": []string{"policy"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return createAuditPolicy(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("update_audit_policy", "Update/replace an audit policy", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the audit policy",
			},
			"policy": map[string]interface{}{
				"type":        "object",
				"description": "Updated audit policy object",
			},
		},
		"required": []string{"name", "policy"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return updateAuditPolicy(ctx, args, rancherClient)
	})

	mcpServer.RegisterToolWithSchema("patch_audit_policy", "Partially update an audit policy", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the audit policy",
			},
			"patch": map[string]interface{}{
				"type":        "object",
				"description": "Partial update patch object",
			},
		},
		"required": []string{"name", "patch"},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return patchAuditPolicy(ctx, args, rancherClient)
	})
}

func createAuditPolicy(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	policy, ok := args["policy"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("policy parameter is required and must be an object")
	}
	return rancherClient.CreateAuditPolicy(ctx, policy)
}

func updateAuditPolicy(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
	if rancherClient == nil {
		return nil, fmt.Errorf("Rancher client not configured")
	}
	name, ok := args["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter is required")
	}
	policy, ok := args["policy"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("policy parameter is required and must be an object")
	}
	return rancherClient.UpdateAuditPolicy(ctx, name, policy)
}

func patchAuditPolicy(ctx context.Context, args map[string]interface{}, rancherClient *client.RancherClient) (interface{}, error) {
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
	return rancherClient.PatchAuditPolicy(ctx, name, patch)
}
