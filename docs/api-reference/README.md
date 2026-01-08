# Rancher Manager API Reference

This directory contains the local API reference documentation for Rancher Manager.

## Source

The API reference is sourced from: https://ranchermanager.docs.rancher.com/api/api-reference

## API Groups

### auditlogCattleIo_v1
- **AuditPolicy**: Audit policy resources

### extCattleIo_v1
- **Kubeconfig**: Kubernetes configuration resources
- **Token**: API token resources

### managementCattleIo_v3
- **ClusterRoleTemplateBinding**: Cluster role template bindings
- **GlobalRoleBinding**: Global role bindings
- **GlobalRole**: Global roles
- **ProjectRoleTemplateBinding**: Project role template bindings
- **Project**: Projects
- **RoleTemplate**: Role templates
- **User**: Users

## Endpoints

All endpoints use the base path: `/apis/{group}/{version}/{resource}`

### List Operations
- `GET /apis/{group}/{version}/{resource}` - List all resources
- `GET /apis/{group}/{version}/namespaces/{namespace}/{resource}` - List namespaced resources

### Get Operations
- `GET /apis/{group}/{version}/{resource}/{name}` - Get a specific resource
- `GET /apis/{group}/{version}/namespaces/{namespace}/{resource}/{name}` - Get a namespaced resource

## Authentication

All requests require Bearer token authentication:
```
Authorization: Bearer <token>
```

## Reference

For the full interactive API reference, visit: https://ranchermanager.docs.rancher.com/api/api-reference
