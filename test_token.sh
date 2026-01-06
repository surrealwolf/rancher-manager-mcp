#!/bin/bash
# Test script to verify Rancher API token
# 
# SECURITY: Never commit tokens to git! Use environment variables.
# Usage: RANCHER_URL=https://your-server RANCHER_TOKEN=your-token ./test_token.sh [--insecure]

TOKEN="${RANCHER_TOKEN:-}"

# Check for --insecure flag
INSECURE=false
if [ "$1" = "--insecure" ] || [ "${RANCHER_INSECURE_SKIP_VERIFY}" = "true" ]; then
  INSECURE=true
fi

if [ -z "$TOKEN" ]; then
  echo "Error: RANCHER_TOKEN environment variable is required"
  echo "Usage: RANCHER_URL=https://your-server RANCHER_TOKEN=your-token ./test_token.sh [--insecure]"
  exit 1
fi

# You'll need to set your Rancher URL
RANCHER_URL="${RANCHER_URL:-https://your-rancher-server}"

echo "Testing Rancher API token..."
echo "URL: $RANCHER_URL"
echo "Token: ${TOKEN:0:20}..."
if [ "$INSECURE" = "true" ]; then
  echo "⚠️  SSL verification disabled"
fi

# Test the token by making a request to the users endpoint
# Use -L to follow redirects
# Ensure URL doesn't have double slashes
RANCHER_URL=$(echo "$RANCHER_URL" | sed 's|/$||')

# Build curl command with optional -k flag for insecure
CURL_CMD="curl -L -s -w \"\n%{http_code}\""
if [ "$INSECURE" = "true" ]; then
  CURL_CMD="curl -k -L -s -w \"\n%{http_code}\""
fi

response=$(eval "$CURL_CMD" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  "$RANCHER_URL/apis/management.cattle.io/v3/users")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
  echo "✓ Token is valid!"
  echo "Response preview:"
  echo "$body" | head -c 200
  echo "..."
else
  echo "✗ Token verification failed (HTTP $http_code)"
  echo "Response:"
  echo "$body"
  exit 1
fi
