# Rancher Manager MCP Server

Model Context Protocol (MCP) server for Rancher Manager API. Control and monitor your Rancher Kubernetes clusters through an AI-powered interface.

## Features

- **75 MCP Tools**: Comprehensive coverage of all Rancher Manager operations
- **Full CRUD Support**: Create, Read, Update, Patch, and Delete operations for all resources
- **Dual Transport Support**: Works with both stdio (for CLI tools) and HTTP (for web services)
- **Rancher API Integration**: Full integration with Rancher Manager Kubernetes API
- **Status Operations**: Get status information from all resource types
- **Namespace Support**: Handles both namespaced and cluster-scoped resources
- **SSL Verification Control**: Configurable SSL certificate verification
- **Environment-based Configuration**: Secure credential management via `.env` files
- **Comprehensive Testing**: Test scripts verify all tools with actual Rancher objects

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/surrealwolf/rancher-manager-mcp.git
cd rancher-manager-mcp

# Build the server
go build -o bin/rancher-mcp ./cmd
```

### Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your Rancher credentials:
```bash
RANCHER_URL=https://your-rancher-server
RANCHER_TOKEN=token-XXXXX:YYYYY
# Optional: Skip SSL verification (not recommended for production)
# RANCHER_INSECURE_SKIP_VERIFY=false
```

### Usage

**Stdio Transport (Default - for MCP clients):**
```bash
source .env
./bin/rancher-mcp --transport stdio
```

**HTTP Transport:**
```bash
source .env
./bin/rancher-mcp --transport http --http-addr :8080
```

## Available Tools

This MCP server provides **75 tools** covering all Rancher Manager operations:

### Cluster Management (9 tools)
* `list_clusters` - List all Rancher clusters
* `get_cluster` - Get details of a specific cluster
* `create_cluster` - Create a new cluster
* `update_cluster` - Update/replace a cluster
* `patch_cluster` - Partially update a cluster
* `delete_cluster` - Delete a cluster
* `get_cluster_status` - Get cluster status

### User Management (8 tools)
* `list_users` - List all Rancher users
* `get_user` - Get details of a specific user
* `create_user` - Create a new user
* `update_user` - Update/replace a user
* `patch_user` - Partially update a user
* `delete_user` - Delete a user
* `get_user_status` - Get user status

### Project Management (9 tools)
* `list_projects` - List all Rancher projects
* `get_project` - Get details of a specific project
* `create_project` - Create a new project
* `update_project` - Update/replace a project
* `patch_project` - Partially update a project
* `delete_project` - Delete a project
* `get_project_status` - Get project status

### Role Templates (9 tools)
* `list_role_templates` - List all role templates
* `get_role_template` - Get details of a role template
* `create_role_template` - Create a new role template
* `update_role_template` - Update/replace a role template
* `patch_role_template` - Partially update a role template
* `delete_role_template` - Delete a role template
* `get_role_template_status` - Get role template status

### Global Roles (9 tools)
* `list_global_roles` - List all global roles
* `get_global_role` - Get details of a global role
* `create_global_role` - Create a new global role
* `update_global_role` - Update/replace a global role
* `patch_global_role` - Partially update a global role
* `delete_global_role` - Delete a global role
* `get_global_role_status` - Get global role status

### Global Role Bindings (9 tools)
* `list_global_role_bindings` - List all global role bindings
* `get_global_role_binding` - Get details of a global role binding
* `create_global_role_binding` - Create a new global role binding
* `update_global_role_binding` - Update/replace a global role binding
* `patch_global_role_binding` - Partially update a global role binding
* `delete_global_role_binding` - Delete a global role binding
* `get_global_role_binding_status` - Get global role binding status

### Cluster Role Template Bindings (9 tools)
* `list_cluster_role_template_bindings` - List all cluster role template bindings
* `get_cluster_role_template_binding` - Get details of a cluster role template binding
* `create_cluster_role_template_binding` - Create a new cluster role template binding
* `update_cluster_role_template_binding` - Update/replace a cluster role template binding
* `patch_cluster_role_template_binding` - Partially update a cluster role template binding
* `delete_cluster_role_template_binding` - Delete a cluster role template binding
* `get_cluster_role_template_binding_status` - Get cluster role template binding status

### Project Role Template Bindings (9 tools)
* `list_project_role_template_bindings` - List all project role template bindings
* `get_project_role_template_binding` - Get details of a project role template binding
* `create_project_role_template_binding` - Create a new project role template binding
* `update_project_role_template_binding` - Update/replace a project role template binding
* `patch_project_role_template_binding` - Partially update a project role template binding
* `delete_project_role_template_binding` - Delete a project role template binding
* `get_project_role_template_binding_status` - Get project role template binding status

### Tokens (6 tools)
* `list_tokens` - List all API tokens
* `get_token` - Get details of a token
* `create_token` - Create a new API token
* `update_token` - Update/replace a token
* `patch_token` - Partially update a token
* `delete_token` - Delete an API token

### Kubeconfigs (6 tools)
* `list_kubeconfigs` - List all kubeconfigs
* `get_kubeconfig` - Get details of a kubeconfig
* `create_kubeconfig` - Create a new kubeconfig
* `update_kubeconfig` - Update/replace a kubeconfig
* `patch_kubeconfig` - Partially update a kubeconfig
* `delete_kubeconfig` - Delete a kubeconfig

### Audit Policies (9 tools)
* `list_audit_policies` - List all audit policies
* `get_audit_policy` - Get details of an audit policy
* `create_audit_policy` - Create a new audit policy
* `update_audit_policy` - Update/replace an audit policy
* `patch_audit_policy` - Partially update an audit policy
* `delete_audit_policy` - Delete an audit policy
* `get_audit_policy_status` - Get audit policy status

**Total: 75 tools** covering all Rancher Manager operations with full CRUD support.

See [docs/TOOLS_REFERENCE.md](docs/TOOLS_REFERENCE.md) for complete tool documentation.

## Cursor IDE Integration

To use this MCP server with Cursor IDE:

1. Open Cursor Settings (Cmd/Ctrl + ,)
2. Navigate to "Features" → "Model Context Protocol"
3. Add a new MCP server with the following configuration:

```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "/absolute/path/to/bin/rancher-mcp",
      "args": ["--transport", "stdio"],
      "env": {
        "RANCHER_URL": "https://your-rancher-server",
        "RANCHER_TOKEN": "your-token-here"
      }
    }
  }
}
```

Or use environment file:
```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "bash",
      "args": ["-c", "source /path/to/.env && /path/to/bin/rancher-mcp --transport stdio"],
      "env": {}
    }
  }
}
```

See [docs/CURSOR_SETUP.md](docs/CURSOR_SETUP.md) for detailed setup instructions.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RANCHER_URL` | Rancher Manager API URL | Required |
| `RANCHER_TOKEN` | Rancher API token (format: `token-XXXXX:YYYYY`) | Required |
| `RANCHER_INSECURE_SKIP_VERIFY` | Skip SSL certificate verification | `false` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

## API Reference

This server uses the Rancher Manager Kubernetes API:
- **Base Path**: `/apis/management.cattle.io/v3/`
- **Authentication**: Bearer token in `Authorization` header
- **Documentation**: https://ranchermanager.docs.rancher.com/api/api-reference

## Docker Support

### Build Docker Image
```bash
docker build -t rancher-mcp:latest .
```

### Run with Docker (stdio mode)
```bash
docker run --rm -i \
  -e RANCHER_URL=https://your-rancher-server \
  -e RANCHER_TOKEN=token-XXXXX:YYYYY \
  rancher-mcp:latest \
  --transport stdio
```

### Run with Docker (HTTP mode)
```bash
docker run --rm -d \
  -p 8080:8080 \
  -e RANCHER_URL=https://your-rancher-server \
  -e RANCHER_TOKEN=token-XXXXX:YYYYY \
  rancher-mcp:latest \
  --transport http --http-addr :8080
```

### Docker Compose Example
```yaml
version: '3.8'
services:
  rancher-mcp:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RANCHER_URL=https://your-rancher-server
      - RANCHER_TOKEN=token-XXXXX:YYYYY
      - RANCHER_INSECURE_SKIP_VERIFY=false
    command: ["--transport", "http", "--http-addr", ":8080"]
```

## Development

### Using Make

The project includes a `Makefile` with common development tasks:

```bash
# Show all available commands
make help

# Build all binaries
make build-all

# Run tests
make test

# Run integration tests (requires RANCHER_URL and RANCHER_TOKEN)
make test-integration

# Run linters
make lint

# Build Docker image
make docker-build

# Clean build artifacts
make clean
```

### Build
```bash
go build -o bin/rancher-mcp ./cmd
```

Or using Make:
```bash
make build
```

### Testing

The project includes comprehensive Go tests for the Rancher client. Tests are integration tests that require a live Rancher instance.

**Prerequisites:**
Set environment variables for your Rancher instance:
```bash
export RANCHER_URL=https://your-rancher-server
export RANCHER_TOKEN=token-XXXXX:YYYYY
export RANCHER_INSECURE_SKIP_VERIFY=false  # Optional: set to true for self-signed certs
```

**Run all tests:**
```bash
go test ./internal/client/... -v
```

**Run specific test suites:**
```bash
# Test list/get operations with existing resources
go test ./internal/client/... -v -run TestList

# Test CRUD operations (creates test objects, tests operations, then cleans up)
go test ./internal/client/... -v -run TestCRUD
```

**Test coverage:**
```bash
go test ./internal/client/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Note:** Tests will be skipped if `RANCHER_URL` and `RANCHER_TOKEN` are not set, making them safe to run in CI/CD pipelines without credentials.

## Project Structure

```
rancher-manager-mcp/
├── cmd/
│   ├── main.go              # Main MCP server entry point
│   └── verify-token/        # Token verification tool
├── internal/
│   ├── client/
│   │   └── rancher_client.go # Rancher API client
│   ├── mcp/
│   │   ├── server.go        # MCP protocol implementation
│   │   └── types.go         # MCP protocol types
│   └── server/
│       ├── server.go        # Server wrapper and tool registration
│       └── handlers/        # Tool handlers for all resources
│           ├── clusters.go
│           ├── users.go
│           ├── projects.go
│           ├── roles*.go
│           └── ...
├── docs/                    # Documentation
│   ├── TOOLS_REFERENCE.md   # Complete tool documentation
│   ├── CURSOR_SETUP.md      # Cursor IDE setup guide
│   └── api-reference/       # API reference documentation
├── bin/                     # Build output (git-ignored)
├── .env.example            # Environment configuration template
└── README.md               # This file
```

## Skills & Capabilities

This MCP server enables:

1. **Cluster Management** - List and inspect Rancher clusters
2. **User Management** - Manage Rancher users and permissions
3. **Project Management** - Organize resources into projects
4. **Infrastructure Monitoring** - Query cluster and resource status
5. **API Integration** - Direct access to Rancher Manager Kubernetes API

## License

MIT License - See LICENSE file for details

## Support

For issues and questions:
* Check the [Rancher API Documentation](https://ranchermanager.docs.rancher.com/api/api-reference)
* Review [docs/TOOLS_REFERENCE.md](docs/TOOLS_REFERENCE.md) for complete tool documentation
* Review implementation examples in `internal/`

---

**Built for Rancher Manager** - Extend your Kubernetes management capabilities with AI. 🤖✨
