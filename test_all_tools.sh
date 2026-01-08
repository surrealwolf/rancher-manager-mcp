#!/bin/bash

# Comprehensive test script for all MCP tools
# This script tests all MCP tools with test objects and uses curl to verify API calls

set +e  # Don't exit on errors, we want to test all endpoints

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if environment variables are set
if [ -z "$RANCHER_URL" ] || [ -z "$RANCHER_TOKEN" ]; then
    echo -e "${RED}Error: RANCHER_URL and RANCHER_TOKEN must be set${NC}"
    echo "Usage: source .env && ./test_all_tools.sh"
    exit 1
fi

# Normalize RANCHER_URL (remove trailing slash)
RANCHER_URL=$(echo "$RANCHER_URL" | sed 's|/$||')

# Check for insecure skip verify option
CURL_INSECURE_OPT=""
if [ "$RANCHER_INSECURE_SKIP_VERIFY" = "true" ] || [ "$RANCHER_INSECURE_SKIP_VERIFY" = "1" ]; then
    CURL_INSECURE_OPT="-k"
    echo -e "${YELLOW}Warning: SSL certificate verification is disabled${NC}"
fi

# Always follow redirects
CURL_FOLLOW_REDIRECTS="-L"

# Clean up function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up test objects...${NC}"
    # Add cleanup logic here if needed
}

trap cleanup EXIT

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Test function
test_tool() {
    local tool_name=$1
    local description=$2
    local curl_url=$3
    local curl_method=${4:-GET}
    local test_data=$5
    
    echo -e "\n${BLUE}Testing: $tool_name${NC}"
    echo "Description: $description"
    
    # Test with curl first
    echo -e "${YELLOW}Curl test:${NC} $curl_method $curl_url"
    if [ "$curl_method" = "GET" ]; then
        curl_response=$(curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X GET \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            "$RANCHER_URL$curl_url" 2>&1)
    elif [ "$curl_method" = "POST" ]; then
        curl_response=$(curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X POST \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            -d "$test_data" \
            "$RANCHER_URL$curl_url" 2>&1)
    elif [ "$curl_method" = "PUT" ]; then
        curl_response=$(curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X PUT \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/json" \
            -H "Accept: application/json" \
            -d "$test_data" \
            "$RANCHER_URL$curl_url" 2>&1)
    elif [ "$curl_method" = "PATCH" ]; then
        curl_response=$(curl $CURL_INSECURE_OPT $CURL_FOLLOW_REDIRECTS -s -w "\n%{http_code}" -X PATCH \
            -H "Authorization: Bearer $RANCHER_TOKEN" \
            -H "Content-Type: application/merge-patch+json" \
            -H "Accept: application/json" \
            -d "$test_data" \
            "$RANCHER_URL$curl_url" 2>&1)
    fi
    
    http_code=$(echo "$curl_response" | tail -n1)
    curl_body=$(echo "$curl_response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ Curl test passed (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
        return 0
    else
        echo -e "${RED}✗ Curl test failed (HTTP $http_code)${NC}"
        echo "Response: $curl_body"
        ((TESTS_FAILED++))
        return 1
    fi
}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Comprehensive MCP Tools Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo "Rancher URL: $RANCHER_URL"
echo ""

# Test 1: List Clusters
test_tool "list_clusters" "List all clusters" "/apis/management.cattle.io/v3/clusters" "GET"

# Test 2: Get Cluster (use first cluster or local if available)
echo -e "\n${BLUE}Testing: get_cluster${NC}"
cluster_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/clusters")
cluster_name=$(echo "$cluster_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$cluster_name" ]; then
    test_tool "get_cluster" "Get cluster details" "/apis/management.cattle.io/v3/clusters/$cluster_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_cluster test: No clusters found${NC}"
fi

# Test 3: Get Cluster Status
if [ -n "$cluster_name" ]; then
    test_tool "get_cluster_status" "Get cluster status" "/apis/management.cattle.io/v3/clusters/$cluster_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_cluster_status test: No clusters found${NC}"
fi

# Test 4: List Users
test_tool "list_users" "List all users" "/apis/management.cattle.io/v3/users" "GET"

# Test 5: Get User
echo -e "\n${BLUE}Testing: get_user${NC}"
user_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/users")
user_name=$(echo "$user_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$user_name" ]; then
    test_tool "get_user" "Get user details" "/apis/management.cattle.io/v3/users/$user_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_user test: No users found${NC}"
fi

# Test 6: Get User Status
if [ -n "$user_name" ]; then
    test_tool "get_user_status" "Get user status" "/apis/management.cattle.io/v3/users/$user_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_user_status test: No users found${NC}"
fi

# Test 7: List Projects
test_tool "list_projects" "List all projects" "/apis/management.cattle.io/v3/projects" "GET"

# Test 8: Get Project
echo -e "\n${BLUE}Testing: get_project${NC}"
project_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/projects")
project_name=$(echo "$project_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
project_namespace=$(echo "$project_list" | jq -r '.items[0].metadata.namespace // .data[0].metadata.namespace // empty' 2>/dev/null || echo "")
if [ -n "$project_name" ]; then
    if [ -n "$project_namespace" ]; then
        test_tool "get_project" "Get project details (namespaced)" "/apis/management.cattle.io/v3/namespaces/$project_namespace/projects/$project_name" "GET"
    else
        test_tool "get_project" "Get project details" "/apis/management.cattle.io/v3/projects/$project_name" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_project test: No projects found${NC}"
fi

# Test 9: Get Project Status
if [ -n "$project_name" ]; then
    if [ -n "$project_namespace" ]; then
        test_tool "get_project_status" "Get project status (namespaced)" "/apis/management.cattle.io/v3/namespaces/$project_namespace/projects/$project_name/status" "GET"
    else
        test_tool "get_project_status" "Get project status" "/apis/management.cattle.io/v3/projects/$project_name/status" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_project_status test: No projects found${NC}"
fi

# Test 10: List Role Templates
test_tool "list_role_templates" "List all role templates" "/apis/management.cattle.io/v3/roletemplates" "GET"

# Test 11: Get Role Template
echo -e "\n${BLUE}Testing: get_role_template${NC}"
role_template_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/roletemplates")
role_template_name=$(echo "$role_template_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$role_template_name" ]; then
    test_tool "get_role_template" "Get role template details" "/apis/management.cattle.io/v3/roletemplates/$role_template_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_role_template test: No role templates found${NC}"
fi

# Test 12: Get Role Template Status
if [ -n "$role_template_name" ]; then
    test_tool "get_role_template_status" "Get role template status" "/apis/management.cattle.io/v3/roletemplates/$role_template_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_role_template_status test: No role templates found${NC}"
fi

# Test 13: List Global Roles
test_tool "list_global_roles" "List all global roles" "/apis/management.cattle.io/v3/globalroles" "GET"

# Test 14: Get Global Role
echo -e "\n${BLUE}Testing: get_global_role${NC}"
global_role_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/globalroles")
global_role_name=$(echo "$global_role_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$global_role_name" ]; then
    test_tool "get_global_role" "Get global role details" "/apis/management.cattle.io/v3/globalroles/$global_role_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_global_role test: No global roles found${NC}"
fi

# Test 15: Get Global Role Status
if [ -n "$global_role_name" ]; then
    test_tool "get_global_role_status" "Get global role status" "/apis/management.cattle.io/v3/globalroles/$global_role_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_global_role_status test: No global roles found${NC}"
fi

# Test 16: List Global Role Bindings
test_tool "list_global_role_bindings" "List all global role bindings" "/apis/management.cattle.io/v3/globalrolebindings" "GET"

# Test 17: Get Global Role Binding
echo -e "\n${BLUE}Testing: get_global_role_binding${NC}"
global_role_binding_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/globalrolebindings")
global_role_binding_name=$(echo "$global_role_binding_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$global_role_binding_name" ]; then
    test_tool "get_global_role_binding" "Get global role binding details" "/apis/management.cattle.io/v3/globalrolebindings/$global_role_binding_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_global_role_binding test: No global role bindings found${NC}"
fi

# Test 18: Get Global Role Binding Status
if [ -n "$global_role_binding_name" ]; then
    test_tool "get_global_role_binding_status" "Get global role binding status" "/apis/management.cattle.io/v3/globalrolebindings/$global_role_binding_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_global_role_binding_status test: No global role bindings found${NC}"
fi

# Test 19: List Cluster Role Template Bindings
test_tool "list_cluster_role_template_bindings" "List all cluster role template bindings" "/apis/management.cattle.io/v3/clusterroletemplatebindings" "GET"

# Test 20: Get Cluster Role Template Binding
echo -e "\n${BLUE}Testing: get_cluster_role_template_binding${NC}"
cluster_rtb_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/clusterroletemplatebindings")
cluster_rtb_name=$(echo "$cluster_rtb_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
cluster_rtb_namespace=$(echo "$cluster_rtb_list" | jq -r '.items[0].metadata.namespace // .data[0].metadata.namespace // empty' 2>/dev/null || echo "")
if [ -n "$cluster_rtb_name" ]; then
    if [ -n "$cluster_rtb_namespace" ]; then
        test_tool "get_cluster_role_template_binding" "Get cluster role template binding (namespaced)" "/apis/management.cattle.io/v3/namespaces/$cluster_rtb_namespace/clusterroletemplatebindings/$cluster_rtb_name" "GET"
    else
        test_tool "get_cluster_role_template_binding" "Get cluster role template binding" "/apis/management.cattle.io/v3/clusterroletemplatebindings/$cluster_rtb_name" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_cluster_role_template_binding test: No bindings found${NC}"
fi

# Test 21: Get Cluster Role Template Binding Status
if [ -n "$cluster_rtb_name" ]; then
    if [ -n "$cluster_rtb_namespace" ]; then
        test_tool "get_cluster_role_template_binding_status" "Get cluster role template binding status (namespaced)" "/apis/management.cattle.io/v3/namespaces/$cluster_rtb_namespace/clusterroletemplatebindings/$cluster_rtb_name/status" "GET"
    else
        test_tool "get_cluster_role_template_binding_status" "Get cluster role template binding status" "/apis/management.cattle.io/v3/clusterroletemplatebindings/$cluster_rtb_name/status" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_cluster_role_template_binding_status test: No bindings found${NC}"
fi

# Test 22: List Project Role Template Bindings
test_tool "list_project_role_template_bindings" "List all project role template bindings" "/apis/management.cattle.io/v3/projectroletemplatebindings" "GET"

# Test 23: Get Project Role Template Binding
echo -e "\n${BLUE}Testing: get_project_role_template_binding${NC}"
project_rtb_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/management.cattle.io/v3/projectroletemplatebindings")
project_rtb_name=$(echo "$project_rtb_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
project_rtb_namespace=$(echo "$project_rtb_list" | jq -r '.items[0].metadata.namespace // .data[0].metadata.namespace // empty' 2>/dev/null || echo "")
if [ -n "$project_rtb_name" ]; then
    if [ -n "$project_rtb_namespace" ]; then
        test_tool "get_project_role_template_binding" "Get project role template binding (namespaced)" "/apis/management.cattle.io/v3/namespaces/$project_rtb_namespace/projectroletemplatebindings/$project_rtb_name" "GET"
    else
        test_tool "get_project_role_template_binding" "Get project role template binding" "/apis/management.cattle.io/v3/projectroletemplatebindings/$project_rtb_name" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_project_role_template_binding test: No bindings found${NC}"
fi

# Test 24: Get Project Role Template Binding Status
if [ -n "$project_rtb_name" ]; then
    if [ -n "$project_rtb_namespace" ]; then
        test_tool "get_project_role_template_binding_status" "Get project role template binding status (namespaced)" "/apis/management.cattle.io/v3/namespaces/$project_rtb_namespace/projectroletemplatebindings/$project_rtb_name/status" "GET"
    else
        test_tool "get_project_role_template_binding_status" "Get project role template binding status" "/apis/management.cattle.io/v3/projectroletemplatebindings/$project_rtb_name/status" "GET"
    fi
else
    echo -e "${YELLOW}⚠ Skipping get_project_role_template_binding_status test: No bindings found${NC}"
fi

# Test 25: List Tokens
test_tool "list_tokens" "List all tokens" "/apis/ext.cattle.io/v1/tokens" "GET"

# Test 26: Get Token
echo -e "\n${BLUE}Testing: get_token${NC}"
token_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/ext.cattle.io/v1/tokens")
token_name=$(echo "$token_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$token_name" ]; then
    test_tool "get_token" "Get token details" "/apis/ext.cattle.io/v1/tokens/$token_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_token test: No tokens found${NC}"
fi

# Test 27: List Kubeconfigs
test_tool "list_kubeconfigs" "List all kubeconfigs" "/apis/ext.cattle.io/v1/kubeconfigs" "GET"

# Test 28: Get Kubeconfig
echo -e "\n${BLUE}Testing: get_kubeconfig${NC}"
kubeconfig_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/ext.cattle.io/v1/kubeconfigs")
kubeconfig_name=$(echo "$kubeconfig_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$kubeconfig_name" ]; then
    test_tool "get_kubeconfig" "Get kubeconfig details" "/apis/ext.cattle.io/v1/kubeconfigs/$kubeconfig_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_kubeconfig test: No kubeconfigs found${NC}"
fi

# Test 29: List Audit Policies
test_tool "list_audit_policies" "List all audit policies" "/apis/auditlog.cattle.io/v1/auditpolicies" "GET"

# Test 30: Get Audit Policy
echo -e "\n${BLUE}Testing: get_audit_policy${NC}"
audit_policy_list=$(curl $CURL_INSECURE_OPT -s -X GET \
    -H "Authorization: Bearer $RANCHER_TOKEN" \
    -H "Accept: application/json" \
    "$RANCHER_URL/apis/auditlog.cattle.io/v1/auditpolicies")
audit_policy_name=$(echo "$audit_policy_list" | jq -r '.items[0].metadata.name // .data[0].metadata.name // empty' 2>/dev/null || echo "")
if [ -n "$audit_policy_name" ]; then
    test_tool "get_audit_policy" "Get audit policy details" "/apis/auditlog.cattle.io/v1/auditpolicies/$audit_policy_name" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_audit_policy test: No audit policies found${NC}"
fi

# Test 31: Get Audit Policy Status
if [ -n "$audit_policy_name" ]; then
    test_tool "get_audit_policy_status" "Get audit policy status" "/apis/auditlog.cattle.io/v1/auditpolicies/$audit_policy_name/status" "GET"
else
    echo -e "${YELLOW}⚠ Skipping get_audit_policy_status test: No audit policies found${NC}"
fi

# Summary
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi
