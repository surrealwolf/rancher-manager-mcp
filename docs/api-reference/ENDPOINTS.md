# Rancher Manager API Endpoints Reference

This document lists all available list and get endpoints in the Rancher Manager API.

## Source

API Reference: https://ranchermanager.docs.rancher.com/api/api-reference

## API Groups

### auditlogCattleIo_v1

#### AuditPolicy
- **List**: `GET /apis/auditlog.cattle.io/v1/auditpolicies`
- **Get**: `GET /apis/auditlog.cattle.io/v1/auditpolicies/{name}`

### extCattleIo_v1

#### Kubeconfig
- **List**: `GET /apis/ext.cattle.io/v1/kubeconfigs`
- **Get**: `GET /apis/ext.cattle.io/v1/kubeconfigs/{name}`

#### Token
- **List**: `GET /apis/ext.cattle.io/v1/tokens`
- **Get**: `GET /apis/ext.cattle.io/v1/tokens/{name}`

### managementCattleIo_v3

#### Cluster
- **List**: `GET /apis/management.cattle.io/v3/clusters`
- **Get**: `GET /apis/management.cattle.io/v3/clusters/{name}`

#### ClusterRoleTemplateBinding
- **List (all namespaces)**: `GET /apis/management.cattle.io/v3/clusterroletemplatebindings`
- **List (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/clusterroletemplatebindings`
- **Get (all namespaces)**: `GET /apis/management.cattle.io/v3/clusterroletemplatebindings/{name}`
- **Get (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/clusterroletemplatebindings/{name}`

#### GlobalRole
- **List**: `GET /apis/management.cattle.io/v3/globalroles`
- **Get**: `GET /apis/management.cattle.io/v3/globalroles/{name}`

#### GlobalRoleBinding
- **List**: `GET /apis/management.cattle.io/v3/globalrolebindings`
- **Get**: `GET /apis/management.cattle.io/v3/globalrolebindings/{name}`

#### Project
- **List (all namespaces)**: `GET /apis/management.cattle.io/v3/projects`
- **List (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/projects`
- **Get (all namespaces)**: `GET /apis/management.cattle.io/v3/projects/{name}`
- **Get (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/projects/{name}`

#### ProjectRoleTemplateBinding
- **List (all namespaces)**: `GET /apis/management.cattle.io/v3/projectroletemplatebindings`
- **List (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/projectroletemplatebindings`
- **Get (all namespaces)**: `GET /apis/management.cattle.io/v3/projectroletemplatebindings/{name}`
- **Get (namespaced)**: `GET /apis/management.cattle.io/v3/namespaces/{namespace}/projectroletemplatebindings/{name}`

#### RoleTemplate
- **List**: `GET /apis/management.cattle.io/v3/roletemplates`
- **Get**: `GET /apis/management.cattle.io/v3/roletemplates/{name}`

#### User
- **List**: `GET /apis/management.cattle.io/v3/users`
- **Get**: `GET /apis/management.cattle.io/v3/users/{name}`

## Authentication

All requests require Bearer token authentication:
```
Authorization: Bearer <token>
```

## MCP Tools

All these endpoints are available as MCP tools. See `docs/TOOLS_REFERENCE.md` for complete tool documentation.
