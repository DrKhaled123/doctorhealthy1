#!/bin/bash

# ================================================================
# Load Testing Script
# Tests application performance under various load conditions
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="${BASE_URL:-http://localhost:8081}"
DURATION="${DURATION:-30}"
WORKERS="${WORKERS:-10}"

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           LOAD TESTING SUITE                       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Configuration:${NC}"
echo "  Base URL: $BASE_URL"
echo "  Duration: ${DURATION}s"
echo "  Concurrent Workers: $WORKERS"
echo ""

# Check if vegeta is installed
if ! command -v vegeta &> /dev/null; then
    echo -e "${YELLOW}Installing vegeta load testing tool...${NC}"
    go install github.com/tsenart/vegeta/v12@latest
fi

mkdir -p coverage/load-tests

# ================================================================
# Test 1: Health Endpoint Load Test
# ================================================================

echo -e "${BLUE}Test 1: Health Endpoint Load${NC}"
echo ""

echo "GET $BASE_URL/health" | vegeta attack \
    -duration=${DURATION}s \
    -rate=100 \
    -workers=$WORKERS \
    | tee coverage/load-tests/health-results.bin \
    | vegeta report

echo ""
echo -e "${GREEN}✓ Health endpoint test complete${NC}"
echo "  Report: coverage/load-tests/health-results.bin"
echo ""

# ================================================================
# Test 2: Recipe API Load Test
# ================================================================

echo -e "${BLUE}Test 2: Recipe API Load${NC}"
echo ""

echo "GET $BASE_URL/api/recipes" | vegeta attack \
    -duration=${DURATION}s \
    -rate=50 \
    -workers=$WORKERS \
    | tee coverage/load-tests/recipes-results.bin \
    | vegeta report

echo ""
echo -e "${GREEN}✓ Recipe API test complete${NC}"
echo ""

# ================================================================
# Test 3: Mixed Workload
# ================================================================

echo -e "${BLUE}Test 3: Mixed Workload${NC}"
echo ""

cat > coverage/load-tests/targets.txt <<EOF
GET $BASE_URL/health
GET $BASE_URL/api/recipes
GET $BASE_URL/api/workouts
GET $BASE_URL/api/diseases
GET $BASE_URL/
EOF

vegeta attack \
    -targets=coverage/load-tests/targets.txt \
    -duration=${DURATION}s \
    -rate=100 \
    -workers=$WORKERS \
    | tee coverage/load-tests/mixed-results.bin \
    | vegeta report

echo ""
echo -e "${GREEN}✓ Mixed workload test complete${NC}"
echo ""

# ================================================================
# Test 4: Spike Test
# ================================================================

echo -e "${BLUE}Test 4: Spike Test${NC}"
echo ""

echo "GET $BASE_URL/health" | vegeta attack \
    -duration=10s \
    -rate=500 \
    -workers=50 \
    | tee coverage/load-tests/spike-results.bin \
    | vegeta report

echo ""
echo -e "${GREEN}✓ Spike test complete${NC}"
echo ""

# ================================================================
# Generate HTML Reports
# ================================================================

echo -e "${BLUE}Generating HTML reports...${NC}"
echo ""

vegeta plot coverage/load-tests/health-results.bin > coverage/load-tests/health-plot.html
vegeta plot coverage/load-tests/recipes-results.bin > coverage/load-tests/recipes-plot.html
vegeta plot coverage/load-tests/mixed-results.bin > coverage/load-tests/mixed-plot.html
vegeta plot coverage/load-tests/spike-results.bin > coverage/load-tests/spike-plot.html

echo -e "${GREEN}✓ HTML reports generated${NC}"
echo ""

# ================================================================
# Summary
# ================================================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Load Test Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Reports available in: coverage/load-tests/"
echo "  - health-plot.html (Health endpoint)"
echo "  - recipes-plot.html (Recipe API)"
echo "  - mixed-plot.html (Mixed workload)"
echo "  - spike-plot.html (Spike test)"
echo ""
echo -e "${GREEN}✓ All load tests completed${NC}"
