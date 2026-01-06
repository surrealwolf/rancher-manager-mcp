# Cursor IDE Setup Guide

This guide explains how to configure the Rancher Manager MCP server in Cursor IDE.

## Prerequisites

1. Build the MCP server:
```bash
go build -o bin/rancher-mcp ./cmd
```

2. Create and configure your `.env` file:
```bash
cp .env.example .env
# Edit .env with your Rancher credentials
```

## Configuration Methods

### Method 1: Direct Command (Recommended)

1. Open Cursor Settings:
   - **macOS**: `Cmd + ,` or `Cursor` → `Settings`
   - **Windows/Linux**: `Ctrl + ,` or `File` → `Preferences` → `Settings`

2. Navigate to: **Features** → **Model Context Protocol**

3. Click **"Add MCP Server"** or edit `settings.json` directly

4. Add the following configuration:

```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "/absolute/path/to/rancher-manager-mcp/bin/rancher-mcp",
      "args": ["--transport", "stdio"],
      "env": {
        "RANCHER_URL": "https://your-rancher-server",
        "RANCHER_TOKEN": "token-XXXXX:YYYYY",
        "RANCHER_INSECURE_SKIP_VERIFY": "false"
      }
    }
  }
}
```

**Important**: Replace `/absolute/path/to/rancher-manager-mcp` with the actual path to your repository.

### Method 2: Using Environment File

If you prefer to use the `.env` file:

```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "bash",
      "args": [
        "-c",
        "source /absolute/path/to/rancher-manager-mcp/.env && /absolute/path/to/rancher-manager-mcp/bin/rancher-mcp --transport stdio"
      ],
      "env": {}
    }
  }
}
```

### Method 3: Shell Script Wrapper

Create a wrapper script for easier management:

1. Create `bin/rancher-mcp-wrapper.sh`:
```bash
#!/bin/bash
cd "$(dirname "$0")/.."
source .env
exec ./bin/rancher-mcp --transport stdio
```

2. Make it executable:
```bash
chmod +x bin/rancher-mcp-wrapper.sh
```

3. Configure in Cursor:
```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "/absolute/path/to/rancher-manager-mcp/bin/rancher-mcp-wrapper.sh",
      "args": [],
      "env": {}
    }
  }
}
```

## Verification

After adding the configuration:

1. Restart Cursor IDE
2. Open the MCP panel (usually in the sidebar or via command palette)
3. You should see "rancher-manager" listed as an available MCP server
4. The server should show as "connected" with available tools listed

## Troubleshooting

### Server Not Connecting

1. **Check the path**: Ensure the path to `rancher-mcp` is absolute and correct
2. **Check permissions**: Make sure the binary is executable:
   ```bash
   chmod +x bin/rancher-mcp
   ```
3. **Check logs**: Look for error messages in Cursor's developer console
   - Open: `Help` → `Toggle Developer Tools` → `Console`

### Authentication Errors

1. **Verify token**: Test your token with:
   ```bash
   ./test_token.sh
   ```
2. **Check environment variables**: Ensure `RANCHER_URL` and `RANCHER_TOKEN` are set correctly
3. **SSL issues**: If using self-signed certificates, set `RANCHER_INSECURE_SKIP_VERIFY=true`

### Tools Not Appearing

1. **Check server status**: Verify the MCP server is running and connected
2. **Check tool registration**: Review `internal/server/server.go` to see registered tools
3. **Restart Cursor**: Sometimes a restart is needed after configuration changes

## Example Configuration (Full)

Here's a complete example configuration with all options:

```json
{
  "mcpServers": {
    "rancher-manager": {
      "command": "/home/user/git/rancher-manager-mcp/bin/rancher-mcp",
      "args": [
        "--transport", "stdio",
        "--log-level", "info"
      ],
      "env": {
        "RANCHER_URL": "https://rancher.example.com",
        "RANCHER_TOKEN": "token-abc123:secret456",
        "RANCHER_INSECURE_SKIP_VERIFY": "false",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

## Next Steps

Once configured, you can:
- Ask Cursor to list your Rancher clusters
- Query cluster status and details
- Manage users and projects
- Get help with Rancher operations

Try asking: "List all my Rancher clusters" or "Show me the details of cluster X"
