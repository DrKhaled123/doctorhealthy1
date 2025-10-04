#!/bin/bash

# ================================================================
# Smoke Test Suite
# Manual validation checklist for critical functionality
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="${BASE_URL:-http://localhost:8081}"

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           SMOKE TEST SUITE                         ║${NC}"
echo -e "${BLUE}║      Manual Validation Checklist                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Testing: $BASE_URL${NC}"
echo ""

PASSED=0
FAILED=0

smoke_test() {
    local test_name=$1
    local endpoint=$2
    local expected=$3
    
    echo -en "${YELLOW}→ $test_name...${NC} "
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" 2>&1)
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL (got $status, expected $expected)${NC}"
        ((FAILED++))
        return 1
    fi
}

# ================================================================
# Critical Path Tests
# ================================================================

echo -e "${BLUE}Critical Path Tests:${NC}"
echo ""

smoke_test "Server is running" "/health" "200"
smoke_test "Frontend loads" "/" "200"
smoke_test "API root responds" "/api" "200"
smoke_test "Recipes endpoint" "/api/recipes" "200"
smoke_test "Workouts endpoint" "/api/workouts" "200"
smoke_test "Diseases endpoint" "/api/diseases" "200"

echo ""

# ================================================================
# Authentication Tests
# ================================================================

echo -e "${BLUE}Authentication Tests:${NC}"
echo ""

# Test without API key
echo -en "${YELLOW}→ Protected endpoint without key...${NC} "
status=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/protected")
if [ "$status" = "401" ]; then
    echo -e "${GREEN}✓ PASS (correctly rejected)${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL (got $status)${NC}"
    ((FAILED++))
fi

echo ""

# ================================================================
# Error Handling Tests
# ================================================================

echo -e "${BLUE}Error Handling Tests:${NC}"
echo ""

smoke_test "404 for invalid endpoint" "/api/nonexistent" "404"
smoke_test "Invalid ID handling" "/api/recipes/99999" "404"

echo ""

# ================================================================
# Performance Tests
# ================================================================

echo -e "${BLUE}Performance Tests:${NC}"
echo ""

echo -en "${YELLOW}→ Response time check...${NC} "
start=$(date +%s%N)
curl -s "$BASE_URL/health" > /dev/null
end=$(date +%s%N)
duration=$((($end - $start) / 1000000))

if [ $duration -lt 1000 ]; then
    echo -e "${GREEN}✓ PASS (${duration}ms)${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ SLOW (${duration}ms)${NC}"
    ((FAILED++))
fi

echo ""

# ================================================================
# Browser Checks (if applicable)
# ================================================================

echo -e "${BLUE}Browser Checks:${NC}"
echo ""

echo -e "${YELLOW}Manual checks needed:${NC}"
echo "  1. Open $BASE_URL in browser"
echo "  2. Check console for errors (F12)"
echo "  3. Verify all pages load"
echo "  4. Test navigation"
echo "  5. Check mobile responsiveness"
echo "  6. Verify forms submit correctly"
echo ""

read -p "Have you completed manual browser checks? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}✓ Manual checks completed${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ Manual checks skipped${NC}"
fi

echo ""

# ================================================================
# Summary
# ================================================================

TOTAL=$((PASSED + FAILED))

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Smoke Test Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Total Tests: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ ALL SMOKE TESTS PASSED${NC}"
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    exit 1
fi
