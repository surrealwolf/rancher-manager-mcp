#!/bin/bash

# Comprehensive test script that creates test objects and tests all MCP tools
# This script creates test objects, tests all operations, and cleans up

set +e  # Don't exit on errors, we want to test all operations

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if environment variables are set
if [ -z "$RANCHER_URL" ] || [ -z "$RANCHER_TOKEN" ]; then
    echo -e "${RED}Error: RANCHER_URL and RANCHER_TOKEN must be set${NC}"
    echo "Usage: source .env && ./test_all_tools_with_objects.sh"
    exit 1
fi

# Normalize RANCHER_URL (remove trailing slash)
RANCHER_URL=$(echo "$RANCHER_URL" | sed 's|/$||')

# Check for insecure skip verify option
CURL_INSECURE_OPT=""
if [ "$RANCHER_INSECURE_SKIP_VERIFY" = "true" ] || [ "$RANCHER_INSECURE_SKIP_VERIFY" = "1" ]; then
    CURL_INSECURE_OPT="-k"
fi

# Always follow redirects
CURL_FOLLOW_REDIRECTS="-L"

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_SKIPPED=0

# Cleanup tracking
CLEANUP_ITEMS=()

# Cleanup function
cleanup() {
    echo -e "\n${CYAN}========================================${NC}"
    echo -e "${CYAN}Cleaning up test objects...${NC}"
    echo -e "${CYAN}========================================${NC}"
    
    local cleanup_count=0
    for item in "${CLEANUP_ITEMS[@]}"; do
        if [ -n "$item" ]; then
            cleanup_count=$((cleanup_count + 1))
        fi
    done
    
    echo "Items to clean up: $cleanup_count"
    # Cleanup is handled at the end of each test section
}

trap cleanup EXIT

# Helper function to make API calls
api_call() {
    local method=$1
    local url=$2
    local data=$3
    
    if [ "$method" = "GET" ]; then
        curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X GET \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            "$RANCHER_URL$url" 2>&1
    elif [ "$method" = "POST" ]; then
        curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X POST \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            -d "$data" \
            "$RANCHER_URL$url" 2>&1
    elif [ "$method" = "PUT" ]; then
        curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X PUT \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            -d "$data" \
            "$RANCHER_URL$url" 2>&1
    elif [ "$method" = "PATCH" ]; then
        curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X PATCH \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/merge-patch+json" \
            -H "Accept: application/json" \
            -d "$data" \
            "$RANCHER_URL$url" 2>&1
    elif [ "$method" = "DELETE" ]; then
        curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X DELETE \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Accept: application/json" \
            "$RANCHER_URL$url" 2>&1
    fi
}

# Test function
test_operation() {
    local operation_name=$1
    local method=$2
    local url=$3
    local test_data=$4
    local expected_status=${5:-200}
    
    echo -e "\n${BLUE}Testing: $operation_name${NC}"
    
    local response=$(api_call "$method" "$url" "$test_data")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -eq "$expected_status" ] || [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ $operation_name passed (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
        return 0
    else
        echo -e "${RED}✗ $operation_name failed (HTTP $http_code)${NC}"
        if [ -n "$body" ] && [ "$body" != "null" ]; then
            echo "Response: $(echo "$body" | jq -c . 2>/dev/null || echo "$body" | head -3)"
        fi
        ((TESTS_FAILED++))
        return 1
    fi
}

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}Comprehensive MCP Tools Test Suite${NC}"
echo -e "${CYAN}Testing with actual Rancher objects${NC}"
echo -e "${CYAN}========================================${NC}"
echo "Rancher URL: $RANCHER_URL"
echo ""

# Generate unique test identifiers
TIMESTAMP=$(date +%s)
TEST_PREFIX="mcp-test-${TIMESTAMP}"

echo -e "${YELLOW}Test prefix: $TEST_PREFIX${NC}"
echo ""

# Get existing resources for reference
echo -e "${CYAN}=== Getting existing resources for reference ===${NC}"
EXISTING_CLUSTER=$(api_call "GET" "/apis/management.cattle.io/v3/clusters" "" | sed '$d' | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
EXISTING_USER=$(api_call "GET" "/apis/management.cattle.io/v3/users" "" | sed '$d' | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
EXISTING_GLOBAL_ROLE=$(api_call "GET" "/apis/management.cattle.io/v3/globalroles" "" | sed '$d' | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
EXISTING_ROLE_TEMPLATE=$(api_call "GET" "/apis/management.cattle.io/v3/roletemplates" "" | sed '$d' | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")

echo "Existing cluster: ${EXISTING_CLUSTER:-none}"
echo "Existing user: ${EXISTING_USER:-none}"
echo "Existing global role: ${EXISTING_GLOBAL_ROLE:-none}"
echo "Existing role template: ${EXISTING_ROLE_TEMPLATE:-none}"
echo ""

# Test 1: List Operations (Read-only, safe to test)
echo -e "${CYAN}=== Testing List Operations ===${NC}"

test_operation "list_clusters" "GET" "/apis/management.cattle.io/v3/clusters"
test_operation "list_users" "GET" "/apis/management.cattle.io/v3/users"
test_operation "list_projects" "GET" "/apis/management.cattle.io/v3/projects"
test_operation "list_role_templates" "GET" "/apis/management.cattle.io/v3/roletemplates"
test_operation "list_global_roles" "GET" "/apis/management.cattle.io/v3/globalroles"
test_operation "list_global_role_bindings" "GET" "/apis/management.cattle.io/v3/globalrolebindings"
test_operation "list_cluster_role_template_bindings" "GET" "/apis/management.cattle.io/v3/clusterroletemplatebindings"
test_operation "list_project_role_template_bindings" "GET" "/apis/management.cattle.io/v3/projectroletemplatebindings"
test_operation "list_tokens" "GET" "/apis/ext.cattle.io/v1/tokens"
test_operation "list_kubeconfigs" "GET" "/apis/ext.cattle.io/v1/kubeconfigs"
test_operation "list_audit_policies" "GET" "/apis/auditlog.cattle.io/v1/auditpolicies"

# Test 2: Get Operations (Read-only, safe to test if resources exist)
echo -e "\n${CYAN}=== Testing Get Operations ===${NC}"

if [ -n "$EXISTING_CLUSTER" ]; then
    test_operation "get_cluster" "GET" "/apis/management.cattle.io/v3/clusters/$EXISTING_CLUSTER"
    # Test status extraction
    STATUS_RESPONSE=$(api_call "GET" "/apis/management.cattle.io/v3/clusters/$EXISTING_CLUSTER" "")
    if echo "$STATUS_RESPONSE" | sed '$d' | jq -e '.status' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ get_cluster_status (status field exists)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠ Status field not found in cluster${NC}"
        ((TESTS_SKIPPED++))
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_cluster tests: No clusters found${NC}"
    ((TESTS_SKIPPED++))
fi

if [ -n "$EXISTING_USER" ]; then
    test_operation "get_user" "GET" "/apis/management.cattle.io/v3/users/$EXISTING_USER"
    # Test status extraction
    STATUS_RESPONSE=$(api_call "GET" "/apis/management.cattle.io/v3/users/$EXISTING_USER" "")
    if echo "$STATUS_RESPONSE" | sed '$d' | jq -e '.status' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ get_user_status (status field exists)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠ Status field not found in user${NC}"
        ((TESTS_SKIPPED++))
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_user tests: No users found${NC}"
    ((TESTS_SKIPPED++))
fi

if [ -n "$EXISTING_GLOBAL_ROLE" ]; then
    test_operation "get_global_role" "GET" "/apis/management.cattle.io/v3/globalroles/$EXISTING_GLOBAL_ROLE"
    # Test status extraction
    STATUS_RESPONSE=$(api_call "GET" "/apis/management.cattle.io/v3/globalroles/$EXISTING_GLOBAL_ROLE" "")
    if echo "$STATUS_RESPONSE" | sed '$d' | jq -e '.status' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ get_global_role_status (status field exists)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠ Status field not found in global role${NC}"
        ((TESTS_SKIPPED++))
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_global_role tests: No global roles found${NC}"
    ((TESTS_SKIPPED++))
fi

if [ -n "$EXISTING_ROLE_TEMPLATE" ]; then
    test_operation "get_role_template" "GET" "/apis/management.cattle.io/v3/roletemplates/$EXISTING_ROLE_TEMPLATE"
    # Test status extraction
    STATUS_RESPONSE=$(api_call "GET" "/apis/management.cattle.io/v3/roletemplates/$EXISTING_ROLE_TEMPLATE" "")
    if echo "$STATUS_RESPONSE" | sed '$d' | jq -e '.status' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ get_role_template_status (status field exists)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠ Status field not found in role template${NC}"
        ((TESTS_SKIPPED++))
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_role_template tests: No role templates found${NC}"
    ((TESTS_SKIPPED++))
fi

# Test 3: Create Operations (Creating test objects)
echo -e "\n${CYAN}=== Testing Create Operations ===${NC}"
echo -e "${YELLOW}Note: Some resources may not support creation via API or require specific permissions${NC}"

# Test Global Role Creation
TEST_GLOBAL_ROLE_NAME="${TEST_PREFIX}-global-role"
TEST_GLOBAL_ROLE_DATA=$(cat <<EOF
{
  "apiVersion": "management.cattle.io/v3",
  "kind": "GlobalRole",
  "metadata": {
    "name": "$TEST_GLOBAL_ROLE_NAME"
  },
  "displayName": "MCP Test Global Role",
  "rules": []
}
EOF
)

if test_operation "create_global_role" "POST" "/apis/management.cattle.io/v3/globalroles" "$TEST_GLOBAL_ROLE_DATA" 201; then
    CLEANUP_ITEMS+=("globalrole:$TEST_GLOBAL_ROLE_NAME")
fi

# Test Role Template Creation
TEST_ROLE_TEMPLATE_NAME="${TEST_PREFIX}-role-template"
TEST_ROLE_TEMPLATE_DATA=$(cat <<EOF
{
  "apiVersion": "management.cattle.io/v3",
  "kind": "RoleTemplate",
  "metadata": {
    "name": "$TEST_ROLE_TEMPLATE_NAME"
  },
  "displayName": "MCP Test Role Template",
  "rules": []
}
EOF
)

if test_operation "create_role_template" "POST" "/apis/management.cattle.io/v3/roletemplates" "$TEST_ROLE_TEMPLATE_DATA" 201; then
    CLEANUP_ITEMS+=("roletemplate:$TEST_ROLE_TEMPLATE_NAME")
fi

# Test 4: Update Operations (Updating test objects we created)
echo -e "\n${CYAN}=== Testing Update Operations ===${NC}"

if [ -n "$TEST_GLOBAL_ROLE_NAME" ]; then
    # Get the existing resource first to get resourceVersion
    EXISTING_RESOURCE=$(api_call "GET" "/apis/management.cattle.io/v3/globalroles/$TEST_GLOBAL_ROLE_NAME" "" | sed '$d')
    RESOURCE_VERSION=$(echo "$EXISTING_RESOURCE" | jq -r '.metadata.resourceVersion // empty' 2>/dev/null)
    
    if [ -n "$RESOURCE_VERSION" ]; then
        UPDATED_GLOBAL_ROLE_DATA=$(cat <<EOF
{
  "apiVersion": "management.cattle.io/v3",
  "kind": "GlobalRole",
  "metadata": {
    "name": "$TEST_GLOBAL_ROLE_NAME",
    "resourceVersion": "$RESOURCE_VERSION"
  },
  "displayName": "MCP Test Global Role - Updated",
  "rules": []
}
EOF
)
        test_operation "update_global_role" "PUT" "/apis/management.cattle.io/v3/globalroles/$TEST_GLOBAL_ROLE_NAME" "$UPDATED_GLOBAL_ROLE_DATA"
    else
        echo -e "${YELLOW}⚠ Skipping update_global_role: Could not get resourceVersion${NC}"
        ((TESTS_SKIPPED++))
    fi
fi

# Test 5: Patch Operations (Partial updates)
echo -e "\n${CYAN}=== Testing Patch Operations ===${NC}"

if [ -n "$TEST_GLOBAL_ROLE_NAME" ]; then
    PATCH_DATA='{"displayName": "MCP Test Global Role - Patched"}'
    test_operation "patch_global_role" "PATCH" "/apis/management.cattle.io/v3/globalroles/$TEST_GLOBAL_ROLE_NAME" "$PATCH_DATA"
fi

# Test 6: Status Operations (on test objects)
echo -e "\n${CYAN}=== Testing Status Operations ===${NC}"

if [ -n "$TEST_GLOBAL_ROLE_NAME" ]; then
    # Status is part of main resource, test getting it
    STATUS_RESPONSE=$(api_call "GET" "/apis/management.cattle.io/v3/globalroles/$TEST_GLOBAL_ROLE_NAME" "")
    if echo "$STATUS_RESPONSE" | sed '$d' | jq -e '.status' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ get_global_role_status (status field exists)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠ Status field not found in global role${NC}"
        ((TESTS_SKIPPED++))
    fi
fi

# Test 7: Delete Operations (Cleanup test objects)
echo -e "\n${CYAN}=== Testing Delete Operations ===${NC}"

for item in "${CLEANUP_ITEMS[@]}"; do
    if [ -z "$item" ]; then
        continue
    fi
    
    resource_type=$(echo "$item" | cut -d':' -f1)
    resource_name=$(echo "$item" | cut -d':' -f2)
    
    if [ -z "$resource_type" ] || [ -z "$resource_name" ]; then
        continue
    fi
    
    case "$resource_type" in
        "globalrole")
            test_operation "delete_global_role" "DELETE" "/apis/management.cattle.io/v3/globalroles/$resource_name" "" 200
            ;;
        "roletemplate")
            test_operation "delete_role_template" "DELETE" "/apis/management.cattle.io/v3/roletemplates/$resource_name" "" 200
            ;;
        *)
            echo -e "${YELLOW}⚠ Unknown resource type for cleanup: $resource_type${NC}"
            ;;
    esac
done

# Final Summary
echo -e "\n${CYAN}========================================${NC}"
echo -e "${CYAN}Test Summary${NC}"
echo -e "${CYAN}========================================${NC}"
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo -e "${YELLOW}Tests Skipped: $TESTS_SKIPPED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All executed tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi
