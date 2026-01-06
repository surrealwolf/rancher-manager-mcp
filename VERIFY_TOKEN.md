# Token Verification

To verify the Rancher API token, you can use either the Go tool or the shell script.

## Using the Go Tool

```bash
go run ./cmd/verify-token/main.go \
  --rancher-url https://your-rancher-server \
  --rancher-token YOUR_TOKEN_HERE
```

Or build it first:

```bash
go build -o bin/verify-token ./cmd/verify-token
./bin/verify-token \
  --rancher-url https://your-rancher-server \
  --rancher-token YOUR_TOKEN_HERE
```

## Using the Shell Script

```bash
export RANCHER_URL=https://your-rancher-server
./test_token.sh
```

## Token Details

- **Token**: `YOUR_TOKEN_HERE`
- **Format**: Bearer token for Rancher Manager API
- **API Endpoint**: `/apis/management.cattle.io/v3/*`

The token is verified by making a request to the `/apis/management.cattle.io/v3/users` endpoint. A successful response (HTTP 200) indicates the token is valid.
