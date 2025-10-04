#!/bin/bash

# ================================================================
# Rollback Script
# Quickly revert to previous deployment if issues detected
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
COOLIFY_API_TOKEN="${COOLIFY_API_TOKEN:-4|jdTX2lUb2q6IOrwNGkHyQBCO74JJeeRHZVvFNwgI6b376a50}"
COOLIFY_URL="http://128.140.111.171:8000"
APP_UUID="hcw0gc8wcwk440gw4c88408o"
DOMAIN="my.doctorhealthy1.com"

echo -e "${RED}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${RED}║         DEPLOYMENT ROLLBACK PROCEDURE             ║${NC}"
echo -e "${RED}╚════════════════════════════════════════════════════╝${NC}"
echo ""

# ================================================================
# Helper Functions
# ================================================================

call_api() {
    local endpoint=$1
    local method=${2:-GET}
    local data=${3:-}
    
    if [ -n "$data" ]; then
        curl -s -X "$method" \
            -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$COOLIFY_URL/api/$endpoint"
    else
        curl -s -X "$method" \
            -H "Authorization: Bearer $COOLIFY_API_TOKEN" \
            "$COOLIFY_URL/api/$endpoint"
    fi
}

# ================================================================
# Step 1: Confirm Rollback
# ================================================================

echo -e "${YELLOW}⚠ WARNING: You are about to rollback the deployment${NC}"
echo ""
echo "This will:"
echo "  - Stop the current version"
echo "  - Redeploy the previous working version"
echo "  - May cause brief downtime"
echo ""

read -p "Are you sure you want to proceed? (yes/no): " -r CONFIRM
echo ""

if [ "$CONFIRM" != "yes" ]; then
    echo "Rollback cancelled"
    exit 0
fi

# ================================================================
# Step 2: Capture Current State
# ================================================================

echo -e "${BLUE}Step 1: Capturing Current State${NC}"
echo ""

CURRENT_COMMIT=$(git rev-parse HEAD)
CURRENT_BRANCH=$(git branch --show-current)

echo "Current commit: $CURRENT_COMMIT"
echo "Current branch: $CURRENT_BRANCH"

# Save rollback info
cat > rollback-info.json <<EOF
{
  "rollback_timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "failed_commit": "$CURRENT_COMMIT",
  "failed_branch": "$CURRENT_BRANCH",
  "rolled_back_by": "$(git config user.name)"
}
EOF

echo "Rollback info saved to rollback-info.json"
echo ""

# ================================================================
# Step 3: Get Deployment History
# ================================================================

echo -e "${BLUE}Step 2: Fetching Deployment History${NC}"
echo ""

DEPLOYMENTS=$(call_api "applications/$APP_UUID/deployments")

echo "Recent deployments:"
echo "$DEPLOYMENTS" | jq -r '.[] | select(.status == "success") | "\(.id) - \(.created_at) - \(.commit_sha // "unknown")"' | head -5

echo ""

# ================================================================
# Step 4: Git Rollback Options
# ================================================================

echo -e "${BLUE}Step 3: Rollback Options${NC}"
echo ""
echo "Choose rollback method:"
echo "  1) Rollback to previous commit (git revert)"
echo "  2) Rollback to specific commit"
echo "  3) Rollback using Coolify UI (manual)"
echo ""

read -p "Select option (1-3): " -r OPTION
echo ""

case $OPTION in
    1)
        echo "Rolling back to previous commit..."
        
        # Get previous commit
        PREVIOUS_COMMIT=$(git rev-parse HEAD~1)
        echo "Previous commit: $PREVIOUS_COMMIT"
        
        # Create revert commit
        git revert --no-edit HEAD
        
        echo -e "${GREEN}✓ Reverted to previous commit${NC}"
        ROLLBACK_COMMIT=$(git rev-parse HEAD)
        ;;
        
    2)
        echo "Available recent commits:"
        git log --oneline -10
        echo ""
        
        read -p "Enter commit SHA to rollback to: " -r TARGET_COMMIT
        
        if [ -z "$TARGET_COMMIT" ]; then
            echo -e "${RED}✗ No commit specified${NC}"
            exit 1
        fi
        
        # Validate commit exists
        if ! git cat-file -e "$TARGET_COMMIT" 2>/dev/null; then
            echo -e "${RED}✗ Invalid commit: $TARGET_COMMIT${NC}"
            exit 1
        fi
        
        echo "Rolling back to $TARGET_COMMIT..."
        git revert --no-edit $CURRENT_COMMIT..$TARGET_COMMIT
        
        echo -e "${GREEN}✓ Reverted to commit $TARGET_COMMIT${NC}"
        ROLLBACK_COMMIT=$(git rev-parse HEAD)
        ;;
        
    3)
        echo -e "${YELLOW}Manual rollback selected${NC}"
        echo ""
        echo "To rollback manually:"
        echo "  1. Go to: $COOLIFY_URL/applications/$APP_UUID"
        echo "  2. Navigate to 'Deployments' tab"
        echo "  3. Find a previous successful deployment"
        echo "  4. Click 'Redeploy' button"
        echo ""
        echo "Press Enter when rollback is complete..."
        read
        exit 0
        ;;
        
    *)
        echo -e "${RED}✗ Invalid option${NC}"
        exit 1
        ;;
esac

echo ""

# ================================================================
# Step 5: Push Rollback Commit
# ================================================================

echo -e "${BLUE}Step 4: Pushing Rollback${NC}"
echo ""

read -p "Push rollback commit to remote? (y/n): " -r PUSH_CONFIRM
echo ""

if [[ $PUSH_CONFIRM =~ ^[Yy]$ ]]; then
    echo "Pushing to $CURRENT_BRANCH..."
    git push origin "$CURRENT_BRANCH"
    echo -e "${GREEN}✓ Rollback commit pushed${NC}"
else
    echo -e "${YELLOW}⚠ Rollback commit not pushed. Manual push required.${NC}"
fi

echo ""

# ================================================================
# Step 6: Trigger Deployment
# ================================================================

echo -e "${BLUE}Step 5: Triggering Deployment${NC}"
echo ""

read -p "Trigger automatic deployment on Coolify? (y/n): " -r DEPLOY_CONFIRM
echo ""

if [[ $DEPLOY_CONFIRM =~ ^[Yy]$ ]]; then
    echo "Triggering deployment..."
    
    DEPLOY_RESPONSE=$(call_api "applications/$APP_UUID/deploy" "POST" '{"force": false}')
    DEPLOYMENT_ID=$(echo "$DEPLOY_RESPONSE" | jq -r '.deployment_id // .id // "unknown"')
    
    if [ "$DEPLOYMENT_ID" != "unknown" ]; then
        echo -e "${GREEN}✓ Deployment triggered: $DEPLOYMENT_ID${NC}"
        echo ""
        echo "Monitor at: $COOLIFY_URL/applications/$APP_UUID/deployments/$DEPLOYMENT_ID"
    else
        echo -e "${YELLOW}⚠ Could not trigger automatic deployment${NC}"
        echo "Please trigger manually in Coolify UI"
    fi
else
    echo -e "${YELLOW}⚠ Automatic deployment skipped${NC}"
    echo "Trigger deployment manually in Coolify UI"
fi

echo ""

# ================================================================
# Step 7: Wait and Verify
# ================================================================

if [[ $DEPLOY_CONFIRM =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Step 6: Waiting for Deployment${NC}"
    echo ""
    
    echo "Waiting 60s for deployment to start..."
    sleep 60
    
    echo "Testing application health..."
    
    MAX_ATTEMPTS=10
    ATTEMPT=0
    
    while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/health" || echo "000")
        
        if [ "$HTTP_CODE" = "200" ]; then
            echo -e "${GREEN}✓ Application is healthy (HTTP $HTTP_CODE)${NC}"
            break
        else
            echo "Attempt $((ATTEMPT+1))/$MAX_ATTEMPTS: HTTP $HTTP_CODE"
            sleep 10
            ((ATTEMPT++))
        fi
    done
    
    if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
        echo -e "${RED}✗ Health check timeout${NC}"
        echo "Check Coolify logs: $COOLIFY_URL/applications/$APP_UUID"
    fi
    
    echo ""
fi

# ================================================================
# Step 8: Summary
# ================================================================

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           ROLLBACK SUMMARY                         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""

cat <<EOF
Rolled back from: $CURRENT_COMMIT
Current state:    $(git rev-parse --short HEAD)
Branch:           $CURRENT_BRANCH
Timestamp:        $(date)

Next Steps:
  1. Verify application: https://$DOMAIN
  2. Check logs: $COOLIFY_URL/applications/$APP_UUID
  3. Investigate root cause of failure
  4. Fix issues before next deployment
  5. Update deployment-info.json and rollback-info.json

Documentation:
  - Rollback info: rollback-info.json
  - Deployment info: deployment-info.json
EOF

echo ""

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║    ✓ ROLLBACK SUCCESSFUL                           ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════╝${NC}"
else
    echo -e "${YELLOW}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║    ⚠ ROLLBACK COMPLETED - VERIFY MANUALLY         ║${NC}"
    echo -e "${YELLOW}╚════════════════════════════════════════════════════╝${NC}"
fi

echo ""
