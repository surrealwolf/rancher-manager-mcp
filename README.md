# Rancher Manager MCP Server

A Model Context Protocol (MCP) server for interacting with Rancher Manager API. Supports both stdio and HTTP transports.

## Features

- **Dual Transport Support**: Works with both stdio (for CLI tools) and HTTP (for web services)
- **Rancher API Integration**: Full integration with Rancher Manager Kubernetes API
- **Tool-based Interface**: Exposes Rancher operations as MCP tools

## Installation

```bash
go build -o rancher-mcp ./cmd
```

## Usage

### Stdio Transport (Default)

```bash
./rancher-mcp \
  --transport stdio \
  --rancher-url https://your-rancher-server \
  --rancher-token token-lk4pv:your-token
```

### HTTP Transport

```bash
./rancher-mcp \
  --transport http \
  --http-addr :8080 \
  --rancher-url https://your-rancher-server \
  --rancher-token token-lk4pv:your-token
```

### Verify Token

```bash
# Using the Go tool
go run ./cmd/verify-token/main.go \
  --rancher-url https://your-rancher-server \
  --rancher-token YOUR_TOKEN_HERE

# Or using the shell script
export RANCHER_URL=https://your-rancher-server
./test_token.sh
```

See [VERIFY_TOKEN.md](VERIFY_TOKEN.md) for more details.

## Available Tools

- `list_clusters` - List all Rancher clusters
- `get_cluster` - Get details of a specific cluster (requires `name` parameter)
- `list_users` - List all Rancher users
- `get_user` - Get details of a specific user (requires `name` parameter)
- `list_projects` - List all Rancher projects
- `get_project` - Get details of a specific project (requires `name` parameter, optional `namespace`)

## API Reference

This server uses the Rancher Manager Kubernetes API. See the [official API documentation](https://ranchermanager.docs.rancher.com/api/api-reference) for more details.

## Development

```bash
# Run tests
go test ./...

# Build
go build -o rancher-mcp ./cmd

# Run with debug logging
./rancher-mcp --log-level debug --transport stdio --rancher-url <url> --rancher-token <token>
```

## License

MIT
