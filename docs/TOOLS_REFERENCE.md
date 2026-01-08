# Rancher Manager MCP Tools Reference

Complete reference for all available MCP tools in the Rancher Manager MCP server.

## Cluster Management

### list_clusters

List all Rancher clusters in the system.

**Parameters**: None

**Returns**: JSON object containing cluster list with metadata

**Example**:
```json
{
  "name": "list_clusters",
  "arguments": {}
}
```

**Response**:
```json
{
  "apiVersion": "management.cattle.io/v3",
  "items": [
    {
      "apiVersion": "management.cattle.io/v3",
      "kind": "Cluster",
      "metadata": {
        "name": "c-abc123",
        "displayName": "production-cluster"
      },
      "status": {
        "conditions": [...]
      }
    }
  ]
}
```

### get_cluster

Get detailed information about a specific cluster.

**Parameters**:
- `name` (string, required) - The name or ID of the cluster

**Example**:
```json
{
  "name": "get_cluster",
  "arguments": {
    "name": "c-abc123"
  }
}
```

**Response**: Full cluster object with all details including:
- Metadata (name, labels, annotations)
- Spec (configuration)
- Status (conditions, node count, version, etc.)

## User Management

### list_users

List all Rancher users.

**Parameters**: None

**Returns**: JSON object containing user list

**Example**:
```json
{
  "name": "list_users",
  "arguments": {}
}
```

**Response**:
```json
{
  "apiVersion": "management.cattle.io/v3",
  "items": [
    {
      "apiVersion": "management.cattle.io/v3",
      "kind": "User",
      "metadata": {
        "name": "user-xyz",
        "displayName": "John Doe"
      },
      "username": "jdoe",
      "enabled": true
    }
  ]
}
```

### get_user

Get detailed information about a specific user.

**Parameters**:
- `name` (string, required) - The name or ID of the user

**Example**:
```json
{
  "name": "get_user",
  "arguments": {
    "name": "user-xyz"
  }
}
```

**Response**: Full user object with:
- User metadata
- Principal IDs
- Enabled status
- Password requirements

## Project Management

### list_projects

List all Rancher projects across all namespaces or filter by namespace.

**Parameters**:
- `namespace` (string, optional) - Optional namespace to filter projects

**Returns**: JSON object containing project list

**Example**:
```json
{
  "name": "list_projects",
  "arguments": {}
}
```

**Response**:
```json
{
  "apiVersion": "management.cattle.io/v3",
  "items": [
    {
      "apiVersion": "management.cattle.io/v3",
      "kind": "Project",
      "metadata": {
        "name": "p-abc123",
        "namespace": "default",
        "displayName": "production"
      }
    }
  ]
}
```

### get_project

Get detailed information about a specific project.

**Parameters**:
- `name` (string, required) - The name or ID of the project
- `namespace` (string, optional) - The namespace containing the project

**Example**:
```json
{
  "name": "get_project",
  "arguments": {
    "name": "p-abc123",
    "namespace": "default"
  }
}
```

**Response**: Full project object with:
- Project metadata
- Resource quotas
- Cluster assignments
- Member access

## Audit Policy Management (auditlogCattleIo_v1)

### list_audit_policies

List all audit policies.

**Parameters**: None

**Returns**: JSON object containing audit policy list

### get_audit_policy

Get detailed information about a specific audit policy.

**Parameters**:
- `name` (string, required) - The name of the audit policy

## Kubeconfig Management (extCattleIo_v1)

### list_kubeconfigs

List all kubeconfigs.

**Parameters**: None

**Returns**: JSON object containing kubeconfig list

### get_kubeconfig

Get detailed information about a specific kubeconfig.

**Parameters**:
- `name` (string, required) - The name of the kubeconfig

## Token Management (extCattleIo_v1)

### list_tokens

List all API tokens.

**Parameters**: None

**Returns**: JSON object containing token list

### get_token

Get detailed information about a specific API token.

**Parameters**:
- `name` (string, required) - The name of the token

## Global Role Management (managementCattleIo_v3)

### list_global_roles

List all global roles.

**Parameters**: None

**Returns**: JSON object containing global role list

### get_global_role

Get detailed information about a specific global role.

**Parameters**:
- `name` (string, required) - The name of the global role

## Global Role Binding Management (managementCattleIo_v3)

### list_global_role_bindings

List all global role bindings.

**Parameters**: None

**Returns**: JSON object containing global role binding list

### get_global_role_binding

Get detailed information about a specific global role binding.

**Parameters**:
- `name` (string, required) - The name of the global role binding

## Role Template Management (managementCattleIo_v3)

### list_role_templates

List all role templates.

**Parameters**: None

**Returns**: JSON object containing role template list

### get_role_template

Get detailed information about a specific role template.

**Parameters**:
- `name` (string, required) - The name of the role template

## Cluster Role Template Binding Management (managementCattleIo_v3)

### list_cluster_role_template_bindings

List all cluster role template bindings, optionally filtered by namespace.

**Parameters**:
- `namespace` (string, optional) - Optional namespace to filter cluster role template bindings

**Returns**: JSON object containing cluster role template binding list

### get_cluster_role_template_binding

Get detailed information about a specific cluster role template binding.

**Parameters**:
- `name` (string, required) - The name of the cluster role template binding
- `namespace` (string, optional) - Optional namespace of the cluster role template binding

## Project Role Template Binding Management (managementCattleIo_v3)

### list_project_role_template_bindings

List all project role template bindings, optionally filtered by namespace.

**Parameters**:
- `namespace` (string, optional) - Optional namespace to filter project role template bindings

**Returns**: JSON object containing project role template binding list

### get_project_role_template_binding

Get detailed information about a specific project role template binding.

**Parameters**:
- `name` (string, required) - The name of the project role template binding
- `namespace` (string, optional) - Optional namespace of the project role template binding

## Error Handling

All tools return errors in the following format:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Error: <error message>"
      }
    ],
    "isError": true
  }
}
```

Common error scenarios:
- **Authentication failed**: Invalid or expired token
- **Resource not found**: Cluster/user/project doesn't exist
- **Permission denied**: Token lacks required permissions
- **Network error**: Cannot reach Rancher server

## API Endpoints

All tools use the Rancher Manager Kubernetes API:

### managementCattleIo_v3
- **Clusters**: `/apis/management.cattle.io/v3/clusters`
- **Users**: `/apis/management.cattle.io/v3/users`
- **Projects**: `/apis/management.cattle.io/v3/projects`
- **GlobalRoles**: `/apis/management.cattle.io/v3/globalroles`
- **GlobalRoleBindings**: `/apis/management.cattle.io/v3/globalrolebindings`
- **RoleTemplates**: `/apis/management.cattle.io/v3/roletemplates`
- **ClusterRoleTemplateBindings**: `/apis/management.cattle.io/v3/clusterroletemplatebindings`
- **ProjectRoleTemplateBindings**: `/apis/management.cattle.io/v3/projectroletemplatebindings`

### extCattleIo_v1
- **Kubeconfigs**: `/apis/ext.cattle.io/v1/kubeconfigs`
- **Tokens**: `/apis/ext.cattle.io/v1/tokens`

### auditlogCattleIo_v1
- **AuditPolicies**: `/apis/auditlog.cattle.io/v1/auditpolicies`

For detailed API documentation, see:
- https://ranchermanager.docs.rancher.com/api/api-reference
- `docs/api-reference/ENDPOINTS.md` (local reference)

## Usage Examples

### List all clusters and their status
```json
{
  "method": "tools/call",
  "params": {
    "name": "list_clusters",
    "arguments": {}
  }
}
```

### Get specific cluster details
```json
{
  "method": "tools/call",
  "params": {
    "name": "get_cluster",
    "arguments": {
      "name": "c-abc123"
    }
  }
}
```

### List all users
```json
{
  "method": "tools/call",
  "params": {
    "name": "list_users",
    "arguments": {}
  }
}
```
