# Quick Reference Guide

Quick reference for using the Rancher Manager MCP server.

## Installation

```bash
git clone https://github.com/surrealwolf/rancher-manager-mcp.git
cd rancher-manager-mcp
go build -o bin/rancher-mcp ./cmd
cp .env.example .env
# Edit .env with your credentials
```

## Configuration

### Environment Variables

```bash
RANCHER_URL=https://your-rancher-server
RANCHER_TOKEN=token-XXXXX:YYYYY
RANCHER_INSECURE_SKIP_VERIFY=false  # Optional
```

### Command Line

```bash
# Stdio transport (for MCP clients)
./bin/rancher-mcp --transport stdio

# HTTP transport
./bin/rancher-mcp --transport http --http-addr :8080

# With flags
./bin/rancher-mcp \
  --rancher-url https://server \
  --rancher-token token \
  --insecure-skip-verify
```

## Available Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `list_clusters` | List all clusters | None |
| `get_cluster` | Get cluster details | `name` (required) |
| `list_users` | List all users | None |
| `get_user` | Get user details | `name` (required) |
| `list_projects` | List all projects | None |
| `get_project` | Get project details | `name` (required), `namespace` (optional) |

## Testing

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

### Verify Token Tool
```bash
go run ./cmd/verify-token/main.go \
  --rancher-url https://server \
  --rancher-token token
```

## Cursor IDE Setup

1. Open Cursor Settings → Features → Model Context Protocol
2. Add MCP server:
```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "/path/to/bin/rancher-mcp",
      "args": ["--transport", "stdio"],
      "env": {
        "RANCHER_URL": "https://server",
        "RANCHER_TOKEN": "token"
      }
    }
  }
}
```

See [CURSOR_SETUP.md](CURSOR_SETUP.md) for detailed instructions.

## API Endpoints

- Clusters: `/apis/management.cattle.io/v3/clusters`
- Users: `/apis/management.cattle.io/v3/users`
- Projects: `/apis/management.cattle.io/v3/projects`

## Troubleshooting

**Token not working?**
```bash
./test_token.sh
```

**SSL certificate errors?**
```bash
# Add to .env
RANCHER_INSECURE_SKIP_VERIFY=true
```

**Server not starting?**
```bash
# Check logs
./bin/rancher-mcp --log-level debug
```

## Resources

- [Full Tools Reference](TOOLS_REFERENCE.md)
- [Cursor Setup Guide](CURSOR_SETUP.md)
- [Rancher API Docs](https://ranchermanager.docs.rancher.com/api/api-reference)
