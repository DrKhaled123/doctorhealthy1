#!/bin/bash

# ================================================================
# Pre-Deployment Validation Script
# Comprehensive checks before deployment
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      PRE-DEPLOYMENT VALIDATION                     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""

CHECKS_PASSED=0
CHECKS_FAILED=0
CRITICAL_FAILURES=()

check() {
    local check_name=$1
    local is_critical=${2:-false}
    
    echo -en "${YELLOW}→ $check_name...${NC} "
}

pass() {
    echo -e "${GREEN}✓ PASS${NC}"
    ((CHECKS_PASSED++))
}

fail() {
    local message=$1
    local is_critical=$2
    
    echo -e "${RED}✗ FAIL${NC}"
    if [ "$message" != "" ]; then
        echo -e "${RED}  Error: $message${NC}"
    fi
    ((CHECKS_FAILED++))
    
    if [ "$is_critical" = "true" ]; then
        CRITICAL_FAILURES+=("$message")
    fi
}

# ================================================================
# 1. Code Quality Checks
# ================================================================

echo -e "${BLUE}1. Code Quality Checks${NC}"
echo ""

check "Go code formatting" true
if [ "$(gofmt -l . | wc -l)" -eq 0 ]; then
    pass
else
    fail "Code is not formatted. Run: gofmt -s -w ." true
fi

check "Go vet" true
if go vet ./... 2>&1 | grep -q "no issues found\|^$"; then
    pass
else
    fail "go vet found issues" true
fi

check "Static analysis" false
if command -v staticcheck &> /dev/null; then
    if staticcheck ./... &> /dev/null; then
        pass
    else
        fail "staticcheck found issues" false
    fi
else
    echo -e "${YELLOW}⚠ SKIP (staticcheck not installed)${NC}"
fi

echo ""

# ================================================================
# 2. Test Coverage
# ================================================================

echo -e "${BLUE}2. Test Coverage${NC}"
echo ""

check "Unit tests" true
if go test ./... &> /dev/null; then
    pass
else
    fail "Unit tests failing" true
fi

check "Test coverage > 50%" false
COVERAGE=$(go test -cover ./... 2>&1 | grep -o 'coverage: [0-9.]*%' | sed 's/coverage: //;s/%//' | awk '{sum+=$1; count++} END {if(count>0) print sum/count; else print 0}')
if (( $(echo "$COVERAGE > 50" | bc -l) )); then
    echo -e "${GREEN}✓ PASS (${COVERAGE}%)${NC}"
    ((CHECKS_PASSED++))
else
    echo -e "${YELLOW}⚠ WARNING (${COVERAGE}%)${NC}"
fi

echo ""

# ================================================================
# 3. Security Checks
# ================================================================

echo -e "${BLUE}3. Security Checks${NC}"
echo ""

check "No secrets in code" true
if bash scripts/check-secrets.sh &> /dev/null; then
    pass
else
    fail "Secrets found in code" true
fi

check "No .env files in git" true
if git ls-files | grep -q '\.env$'; then
    fail ".env file is tracked in git" true
else
    pass
fi

check "Dependencies audit" false
if command -v nancy &> /dev/null; then
    if go list -json -m all | nancy sleuth &> /dev/null; then
        pass
    else
        fail "Vulnerable dependencies found" false
    fi
else
    echo -e "${YELLOW}⚠ SKIP (nancy not installed)${NC}"
fi

echo ""

# ================================================================
# 4. Configuration Validation
# ================================================================

echo -e "${BLUE}4. Configuration Validation${NC}"
echo ""

check ".env.example exists" false
if [ -f ".env.example" ]; then
    pass
else
    fail ".env.example not found" false
fi

check "Required env vars documented" false
REQUIRED_VARS=("JWT_SECRET" "PORT" "DB_PATH" "ENV")
MISSING_VARS=()
for var in "${REQUIRED_VARS[@]}"; do
    if ! grep -q "^$var=" .env.example 2>/dev/null; then
        MISSING_VARS+=("$var")
    fi
done

if [ ${#MISSING_VARS[@]} -eq 0 ]; then
    pass
else
    fail "Missing vars in .env.example: ${MISSING_VARS[*]}" false
fi

check "Dockerfile exists" true
if [ -f "Dockerfile" ]; then
    pass
else
    fail "Dockerfile not found" true
fi

echo ""

# ================================================================
# 5. Build Validation
# ================================================================

echo -e "${BLUE}5. Build Validation${NC}"
echo ""

check "Application builds" true
if CGO_ENABLED=1 go build -o /tmp/test-build . &> /dev/null; then
    pass
    rm -f /tmp/test-build
else
    fail "Build failed" true
fi

check "Docker image builds" false
if command -v docker &> /dev/null; then
    if docker build -t test-build:latest . &> /dev/null; then
        pass
        docker rmi test-build:latest &> /dev/null || true
    else
        fail "Docker build failed" false
    fi
else
    echo -e "${YELLOW}⚠ SKIP (docker not available)${NC}"
fi

echo ""

# ================================================================
# 6. Documentation
# ================================================================

echo -e "${BLUE}6. Documentation${NC}"
echo ""

check "README exists" false
if [ -f "README.md" ]; then
    pass
else
    fail "README.md not found" false
fi

check "Deployment docs exist" false
if [ -f "COOLIFY-MANUAL-STEPS.md" ] || [ -f "DEPLOYMENT-CHECKLIST.md" ]; then
    pass
else
    fail "Deployment documentation missing" false
fi

echo ""

# ================================================================
# 7. Git Status
# ================================================================

echo -e "${BLUE}7. Git Status${NC}"
echo ""

check "No uncommitted changes" false
if [ -z "$(git status --porcelain)" ]; then
    pass
else
    fail "Uncommitted changes present" false
fi

check "On main/master branch" false
BRANCH=$(git branch --show-current)
if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
    pass
else
    echo -e "${YELLOW}⚠ WARNING (on branch: $BRANCH)${NC}"
fi

echo ""

# ================================================================
# Summary
# ================================================================

TOTAL_CHECKS=$((CHECKS_PASSED + CHECKS_FAILED))

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Validation Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Total Checks: $TOTAL_CHECKS"
echo -e "${GREEN}Passed: $CHECKS_PASSED${NC}"
echo -e "${RED}Failed: $CHECKS_FAILED${NC}"
echo ""

if [ ${#CRITICAL_FAILURES[@]} -gt 0 ]; then
    echo -e "${RED}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║    CRITICAL FAILURES - DEPLOYMENT BLOCKED         ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Critical issues:"
    for failure in "${CRITICAL_FAILURES[@]}"; do
        echo -e "${RED}  • $failure${NC}"
    done
    echo ""
    exit 1
elif [ $CHECKS_FAILED -gt 0 ]; then
    echo -e "${YELLOW}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║    WARNINGS - REVIEW BEFORE DEPLOYMENT            ║${NC}"
    echo -e "${YELLOW}╚════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Some non-critical checks failed. Review before deploying."
    echo ""
    read -p "Continue with deployment? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Deployment cancelled"
        exit 1
    fi
else
    echo -e "${GREEN}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║    ✓ ALL CHECKS PASSED - READY TO DEPLOY          ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════╝${NC}"
    echo ""
fi

echo -e "${GREEN}✓ Pre-deployment validation complete${NC}"
exit 0
