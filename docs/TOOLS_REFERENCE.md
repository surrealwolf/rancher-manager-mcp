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

List all Rancher projects across all namespaces.

**Parameters**: None

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

- **Clusters**: `/apis/management.cattle.io/v3/clusters`
- **Users**: `/apis/management.cattle.io/v3/users`
- **Projects**: `/apis/management.cattle.io/v3/projects`

For detailed API documentation, see:
https://ranchermanager.docs.rancher.com/api/api-reference

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
