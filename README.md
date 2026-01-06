# Rancher Manager MCP Server

Model Context Protocol (MCP) server for Rancher Manager API. Control and monitor your Rancher Kubernetes clusters through an AI-powered interface.

## Features

- **Dual Transport Support**: Works with both stdio (for CLI tools) and HTTP (for web services)
- **Rancher API Integration**: Full integration with Rancher Manager Kubernetes API
- **Tool-based Interface**: Exposes Rancher operations as MCP tools
- **SSL Verification Control**: Configurable SSL certificate verification
- **Environment-based Configuration**: Secure credential management via `.env` files

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

### Cluster Management (2 tools)
* `list_clusters` - List all Rancher clusters
* `get_cluster` - Get details of a specific cluster (requires `name` parameter)

### User Management (2 tools)
* `list_users` - List all Rancher users
* `get_user` - Get details of a specific user (requires `name` parameter)

### Project Management (2 tools)
* `list_projects` - List all Rancher projects
* `get_project` - Get details of a specific project (requires `name` parameter, optional `namespace`)

**Total: 6 tools** covering essential Rancher Manager operations.

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

## Development

### Build
```bash
go build -o bin/rancher-mcp ./cmd
```

### Test Token
```bash
source .env
./test_token.sh
```

### Test MCP Server
```bash
source .env
./test_mcp.sh
```

### Run Tests
```bash
go test ./...
```

## Project Structure

```
rancher-manager-mcp/
├── cmd/
│   ├── main.go              # Main MCP server entry point
│   └── verify-token/        # Token verification tool
├── internal/
│   ├── mcp/
│   │   ├── server.go        # MCP protocol implementation
│   │   └── types.go         # MCP protocol types
│   └── server/
│       ├── server.go        # Server wrapper and tool registration
│       └── rancher_client.go # Rancher API client
├── docs/                    # Documentation
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
