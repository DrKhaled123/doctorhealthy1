#!/bin/bash

# ================================================================
# Secret Scanning Script
# Checks for accidentally committed secrets and sensitive data
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Secret Scanning${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

ISSUES_FOUND=0

# Patterns to search for
PATTERNS=(
    "password\s*=\s*['\"]"
    "api[_-]?key\s*=\s*['\"]"
    "secret\s*=\s*['\"]"
    "token\s*=\s*['\"]"
    "private[_-]?key"
    "aws[_-]?access"
    "BEGIN\s+(RSA|DSA|EC)\s+PRIVATE\s+KEY"
    "sk_live_"
    "pk_live_"
    "xox[a-zA-Z]-"
)

# Files to exclude
EXCLUDE_PATTERNS=(
    ".git/*"
    "node_modules/*"
    "vendor/*"
    "*.log"
    "coverage/*"
    "*.test.go"
    "test-*"
    "*_test.go"
    "*.md"
)

echo -e "${YELLOW}Scanning for secrets in source code...${NC}"
echo ""

for pattern in "${PATTERNS[@]}"; do
    echo -en "${YELLOW}Checking for: $pattern...${NC} "
    
    # Build exclude arguments
    EXCLUDE_ARGS=""
    for exclude in "${EXCLUDE_PATTERNS[@]}"; do
        EXCLUDE_ARGS="$EXCLUDE_ARGS --exclude='$exclude'"
    done
    
    # Search for pattern
    results=$(grep -riE "$pattern" . \
        --exclude-dir=.git \
        --exclude-dir=node_modules \
        --exclude-dir=vendor \
        --exclude-dir=coverage \
        --exclude="*.log" \
        --exclude="*.md" \
        --exclude="*_test.go" \
        --exclude="test-*" \
        2>/dev/null || true)
    
    if [ -z "$results" ]; then
        echo -e "${GREEN}✓ Clean${NC}"
    else
        echo -e "${RED}✗ Found matches:${NC}"
        echo "$results" | while read -r line; do
            echo -e "${RED}  → $line${NC}"
        done
        ((ISSUES_FOUND++))
    fi
done

echo ""

# Check for .env files in git
echo -en "${YELLOW}Checking for .env files in git...${NC} "
ENV_FILES=$(git ls-files | grep -E '\.env$|\.env\.production$' || true)
if [ -z "$ENV_FILES" ]; then
    echo -e "${GREEN}✓ No .env files tracked${NC}"
else
    echo -e "${RED}✗ Found .env files in git:${NC}"
    echo "$ENV_FILES" | while read -r file; do
        echo -e "${RED}  → $file${NC}"
    done
    ((ISSUES_FOUND++))
fi

echo ""

# Check for hardcoded IPs
echo -en "${YELLOW}Checking for hardcoded IP addresses...${NC} "
IPS=$(grep -riE '\b([0-9]{1,3}\.){3}[0-9]{1,3}\b' \
    --include="*.go" \
    --exclude-dir=.git \
    --exclude="*_test.go" \
    . | grep -v "127.0.0.1\|localhost\|0.0.0.0" || true)
if [ -z "$IPS" ]; then
    echo -e "${GREEN}✓ No suspicious IPs${NC}"
else
    echo -e "${YELLOW}⚠ Found IP addresses:${NC}"
    echo "$IPS" | head -n 5
    if [ $(echo "$IPS" | wc -l) -gt 5 ]; then
        echo -e "${YELLOW}  ... and $(( $(echo "$IPS" | wc -l) - 5 )) more${NC}"
    fi
fi

echo ""

# Summary
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✓ No secrets found${NC}"
    exit 0
else
    echo -e "${RED}✗ Found $ISSUES_FOUND potential security issues${NC}"
    echo -e "${YELLOW}Please review and remove any exposed secrets${NC}"
    exit 1
fi
