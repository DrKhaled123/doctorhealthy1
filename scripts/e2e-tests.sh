#!/bin/bash

# ================================================================
# E2E (End-to-End) Test Suite
# Tests complete user workflows and API functionality
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="${BASE_URL:-http://localhost:8082}"
TEST_RESULTS=()
FAILED_TESTS=0
PASSED_TESTS=0

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  End-to-End Test Suite${NC}"
echo -e "${BLUE}  Testing: $BASE_URL${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test helper functions
test_endpoint() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local expected_status=$4
    local data=${5:-""}
    
    echo -en "${YELLOW}Testing: $test_name...${NC} "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $status_code)"
        TEST_RESULTS+=("PASS: $test_name")
        ((PASSED_TESTS++))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (Expected: $expected_status, Got: $status_code)"
        echo -e "${RED}Response: $body${NC}"
        TEST_RESULTS+=("FAIL: $test_name - Expected $expected_status, got $status_code")
        ((FAILED_TESTS++))
        return 1
    fi
}

# ================================================================
# Test Suite 1: Health & Status Endpoints
# ================================================================

echo -e "${BLUE}Test Suite 1: Health & Status${NC}"
echo ""

test_endpoint "Health Check" "GET" "/health" "200"
test_endpoint "API Status" "GET" "/" "200"
test_endpoint "Frontend Loads" "GET" "/" "200"

echo ""

# ================================================================
# Test Suite 2: API Key Management
# ================================================================

echo -e "${BLUE}Test Suite 2: API Key Management${NC}"
echo ""

# Generate test API key
API_KEY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/keys/generate" \
    -H "Content-Type: application/json" \
    -d '{"name":"test-key","scopes":["read","write"]}')

TEST_API_KEY=$(echo "$API_KEY_RESPONSE" | grep -o '"key":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$TEST_API_KEY" ]; then
    echo -e "${GREEN}✓ API Key Generated: ${TEST_API_KEY:0:20}...${NC}"
    TEST_RESULTS+=("PASS: API Key Generation")
    ((PASSED_TESTS++))
else
    echo -e "${RED}✗ Failed to generate API key${NC}"
    TEST_RESULTS+=("FAIL: API Key Generation")
    ((FAILED_TESTS++))
fi

# Test API key validation
test_endpoint "Valid API Key" "GET" "/api/test" "200" "" \
    -H "X-API-Key: $TEST_API_KEY"

test_endpoint "Invalid API Key" "GET" "/api/test" "401" "" \
    -H "X-API-Key: invalid-key-123"

test_endpoint "Missing API Key" "GET" "/api/test" "401" ""

echo ""

# ================================================================
# Test Suite 3: Recipe Endpoints
# ================================================================

echo -e "${BLUE}Test Suite 3: Recipe Management${NC}"
echo ""

test_endpoint "Get All Recipes" "GET" "/api/recipes" "200"
test_endpoint "Get Recipe by ID" "GET" "/api/recipes/1" "200"
test_endpoint "Get Non-existent Recipe" "GET" "/api/recipes/99999" "404"
test_endpoint "Search Recipes" "GET" "/api/recipes/search?q=healthy" "200"
test_endpoint "Filter by Cuisine" "GET" "/api/recipes?cuisine=mediterranean" "200"

echo ""

# ================================================================
# Test Suite 4: Nutrition Endpoints
# ================================================================

echo -e "${BLUE}Test Suite 4: Nutrition Data${NC}"
echo ""

test_endpoint "Get Nutrition Info" "GET" "/api/nutrition/1" "200"
test_endpoint "Calculate Calories" "POST" "/api/nutrition/calculate" "200" \
    '{"ingredients":["chicken","rice","vegetables"],"servings":2}'

echo ""

# ================================================================
# Test Suite 5: Workout Endpoints
# ================================================================

echo -e "${BLUE}Test Suite 5: Workout Management${NC}"
echo ""

test_endpoint "Get All Workouts" "GET" "/api/workouts" "200"
test_endpoint "Get Workout by ID" "GET" "/api/workouts/1" "200"
test_endpoint "Search Workouts" "GET" "/api/workouts/search?q=strength" "200"

echo ""

# ================================================================
# Test Suite 6: Health Conditions
# ================================================================

echo -e "${BLUE}Test Suite 6: Health Conditions${NC}"
echo ""

test_endpoint "Get Diseases" "GET" "/api/diseases" "200"
test_endpoint "Get Disease by ID" "GET" "/api/diseases/1" "200"
test_endpoint "Get Injuries" "GET" "/api/injuries" "200"
test_endpoint "Get Complaints" "GET" "/api/complaints" "200"

echo ""

# ================================================================
# Test Suite 7: Plan Generation
# ================================================================

echo -e "${BLUE}Test Suite 7: Plan Generation${NC}"
echo ""

test_endpoint "Generate Meal Plan" "POST" "/api/plans/meal" "200" \
    '{"goal":"weight_loss","dietType":"balanced","calories":2000}'

test_endpoint "Generate Workout Plan" "POST" "/api/plans/workout" "200" \
    '{"goal":"strength","level":"intermediate","days":5}'

echo ""

# ================================================================
# Test Suite 8: Error Handling
# ================================================================

echo -e "${BLUE}Test Suite 8: Error Handling${NC}"
echo ""

test_endpoint "Invalid Endpoint" "GET" "/api/nonexistent" "404"
test_endpoint "Invalid Method" "DELETE" "/health" "405"
test_endpoint "Malformed JSON" "POST" "/api/nutrition/calculate" "400" \
    '{invalid json}'

echo ""

# ================================================================
# Test Suite 9: Security Headers
# ================================================================

echo -e "${BLUE}Test Suite 9: Security Headers${NC}"
echo ""

echo -en "${YELLOW}Checking security headers...${NC} "
HEADERS=$(curl -s -I "$BASE_URL/")

check_header() {
    local header=$1
    if echo "$HEADERS" | grep -qi "$header"; then
        echo -e "${GREEN}✓ $header present${NC}"
        ((PASSED_TESTS++))
        return 0
    else
        echo -e "${RED}✗ $header missing${NC}"
        ((FAILED_TESTS++))
        return 1
    fi
}

check_header "X-Content-Type-Options"
check_header "X-Frame-Options"
check_header "X-XSS-Protection"

echo ""

# ================================================================
# Test Suite 10: Performance
# ================================================================

echo -e "${BLUE}Test Suite 10: Performance Tests${NC}"
echo ""

echo -en "${YELLOW}Testing response time...${NC} "
START_TIME=$(date +%s%N)
curl -s "$BASE_URL/health" > /dev/null
END_TIME=$(date +%s%N)
DURATION=$((($END_TIME - $START_TIME) / 1000000))

if [ $DURATION -lt 500 ]; then
    echo -e "${GREEN}✓ Response time: ${DURATION}ms${NC}"
    ((PASSED_TESTS++))
else
    echo -e "${YELLOW}⚠ Response time: ${DURATION}ms (> 500ms)${NC}"
    ((FAILED_TESTS++))
fi

echo ""

# ================================================================
# Test Summary
# ================================================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Test Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Total Tests: $((PASSED_TESTS + FAILED_TESTS))"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║    ALL E2E TESTS PASSED ✓                          ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════╝${NC}"
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║    SOME TESTS FAILED ✗                             ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Failed Tests:"
    for result in "${TEST_RESULTS[@]}"; do
        if [[ $result == FAIL:* ]]; then
            echo -e "${RED}  - ${result#FAIL: }${NC}"
        fi
    done
    exit 1
fi
