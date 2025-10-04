#!/bin/bash

# ================================================================
# Automated Coolify Deployment Script
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

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      COOLIFY DEPLOYMENT AUTOMATION                ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
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

log_step() {
    echo -e "${BLUE}→ $1${NC}"
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# ================================================================
# Step 1: Pre-deployment Validation
# ================================================================

echo -e "${BLUE}Step 1: Pre-deployment Validation${NC}"
echo ""

log_step "Running pre-deployment checks..."
if bash scripts/pre-deploy-check.sh; then
    log_success "Pre-deployment checks passed"
else
    log_error "Pre-deployment checks failed"
    exit 1
fi

echo ""

# ================================================================
# Step 2: Test API Connection
# ================================================================

echo -e "${BLUE}Step 2: API Connection${NC}"
echo ""

log_step "Testing Coolify API connection..."
HEALTH_CHECK=$(call_api "health")

if echo "$HEALTH_CHECK" | grep -q "ok\|healthy\|success"; then
    log_success "Coolify API connection successful"
else
    log_error "Cannot connect to Coolify API"
    exit 1
fi

echo ""

# ================================================================
# Step 3: Get Application Status
# ================================================================

echo -e "${BLUE}Step 3: Application Status${NC}"
echo ""

log_step "Fetching current application status..."
APP_STATUS=$(call_api "applications/$APP_UUID")

CURRENT_STATE=$(echo "$APP_STATUS" | jq -r '.status // "unknown"')
log_success "Current status: $CURRENT_STATE"

echo ""

# ================================================================
# Step 4: Create Deployment Backup Point
# ================================================================

echo -e "${BLUE}Step 4: Backup Point${NC}"
echo ""

log_step "Creating deployment backup information..."

DEPLOYMENT_INFO=$(cat <<EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "commit": "$(git rev-parse HEAD)",
  "branch": "$(git branch --show-current)",
  "deployed_by": "$(git config user.name)",
  "domain": "$DOMAIN"
}
EOF
)

echo "$DEPLOYMENT_INFO" > deployment-info.json
log_success "Backup info saved to deployment-info.json"

echo ""

# ================================================================
# Step 5: Update Environment Variables
# ================================================================

echo -e "${BLUE}Step 5: Environment Variables${NC}"
echo ""

log_step "Updating environment variables..."

if [ -f "coolify-env-vars.txt" ]; then
    # Read env vars from file
    ENV_VARS=$(cat coolify-env-vars.txt | jq -R -s 'split("\n") | map(select(length > 0)) | map(split("=")) | map({(.[0]): .[1]}) | add')
    
    # Update via API (example - adjust based on Coolify API)
    log_success "Environment variables prepared"
    log_warning "Manual verification recommended in Coolify UI"
else
    log_warning "coolify-env-vars.txt not found, skipping env update"
fi

echo ""

# ================================================================
# Step 6: Trigger Deployment
# ================================================================

echo -e "${BLUE}Step 6: Trigger Deployment${NC}"
echo ""

log_step "Triggering deployment..."

# Option 1: Using Coolify API
DEPLOY_RESPONSE=$(call_api "applications/$APP_UUID/deploy" "POST" '{"force": false}')

DEPLOYMENT_ID=$(echo "$DEPLOY_RESPONSE" | jq -r '.deployment_id // .id // "unknown"')

if [ "$DEPLOYMENT_ID" != "unknown" ]; then
    log_success "Deployment triggered: $DEPLOYMENT_ID"
else
    log_warning "Deployment triggered but no ID returned"
    log_warning "Check Coolify UI for deployment status"
fi

echo ""

# ================================================================
# Step 7: Monitor Deployment
# ================================================================

echo -e "${BLUE}Step 7: Monitor Deployment${NC}"
echo ""

log_step "Monitoring deployment progress..."
echo "This may take several minutes..."
echo ""

# Wait for deployment to start
sleep 10

# Monitor deployment status (adjust based on actual Coolify API)
MAX_ATTEMPTS=60
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    DEPLOY_STATUS=$(call_api "applications/$APP_UUID/deployments/$DEPLOYMENT_ID" 2>/dev/null || echo '{"status":"pending"}')
    STATUS=$(echo "$DEPLOY_STATUS" | jq -r '.status // "unknown"')
    
    echo -ne "\r${YELLOW}Status: $STATUS (attempt $((ATTEMPT+1))/$MAX_ATTEMPTS)${NC}     "
    
    if [ "$STATUS" = "success" ] || [ "$STATUS" = "completed" ]; then
        echo ""
        log_success "Deployment completed successfully"
        break
    elif [ "$STATUS" = "failed" ] || [ "$STATUS" = "error" ]; then
        echo ""
        log_error "Deployment failed"
        echo ""
        echo "Check logs at: $COOLIFY_URL/applications/$APP_UUID/deployments/$DEPLOYMENT_ID"
        exit 1
    fi
    
    sleep 10
    ((ATTEMPT++))
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo ""
    log_warning "Deployment monitoring timeout"
    log_warning "Check Coolify UI for final status: $COOLIFY_URL"
fi

echo ""

# ================================================================
# Step 8: Post-Deployment Health Check
# ================================================================

echo -e "${BLUE}Step 8: Post-Deployment Health Check${NC}"
echo ""

log_step "Waiting for application to be ready..."
sleep 30

log_step "Testing application health..."

HEALTH_URL="https://$DOMAIN/health"
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL" || echo "000")

if [ "$HEALTH_STATUS" = "200" ]; then
    log_success "Application is healthy (HTTP $HEALTH_STATUS)"
else
    log_error "Health check failed (HTTP $HEALTH_STATUS)"
    log_warning "Application may still be starting up"
fi

echo ""

# ================================================================
# Step 9: Smoke Tests
# ================================================================

echo -e "${BLUE}Step 9: Smoke Tests${NC}"
echo ""

log_step "Running production smoke tests..."

# Test API root
API_ROOT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/" || echo "000")
if [ "$API_ROOT_STATUS" = "200" ]; then
    log_success "API root accessible"
else
    log_warning "API root returned HTTP $API_ROOT_STATUS"
fi

# Test API endpoint
RECIPES_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/api/recipes" || echo "000")
if [ "$RECIPES_STATUS" = "200" ] || [ "$RECIPES_STATUS" = "401" ]; then
    log_success "API endpoints responding"
else
    log_warning "API test returned HTTP $RECIPES_STATUS"
fi

echo ""

# ================================================================
# Step 10: Deployment Summary
# ================================================================

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           DEPLOYMENT SUMMARY                       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""

cat <<EOF
Domain:          https://$DOMAIN
Application:     $APP_UUID
Deployment ID:   $DEPLOYMENT_ID
Commit:          $(git rev-parse --short HEAD)
Branch:          $(git branch --show-current)
Deployed:        $(date)
Health Status:   $HEALTH_STATUS

Next Steps:
  1. Verify application in browser: https://$DOMAIN
  2. Check Coolify logs: $COOLIFY_URL/applications/$APP_UUID
  3. Monitor for errors: tail -f logs/app.log
  4. Run full smoke tests: bash scripts/smoke-tests.sh

Rollback (if needed):
  1. Go to: $COOLIFY_URL/applications/$APP_UUID
  2. Navigate to Deployments tab
  3. Click "Redeploy" on previous successful deployment
EOF

echo ""

if [ "$HEALTH_STATUS" = "200" ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║    ✓ DEPLOYMENT SUCCESSFUL                         ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════╝${NC}"
    exit 0
else
    echo -e "${YELLOW}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║    ⚠ DEPLOYMENT COMPLETED WITH WARNINGS            ║${NC}"
    echo -e "${YELLOW}╚════════════════════════════════════════════════════╝${NC}"
    exit 0
fi
