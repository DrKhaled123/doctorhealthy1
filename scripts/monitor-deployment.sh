#!/bin/bash

# ================================================================
# Post-Deployment Monitoring Script
# Monitors application health and logs after deployment
# ================================================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
DOMAIN="${DOMAIN:-my.doctorhealthy1.com}"
CHECK_INTERVAL=${CHECK_INTERVAL:-30}
ALERT_THRESHOLD=${ALERT_THRESHOLD:-3}

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      POST-DEPLOYMENT MONITORING                    ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""

CONSECUTIVE_FAILURES=0
TOTAL_CHECKS=0
SUCCESSFUL_CHECKS=0
FAILED_CHECKS=0

# ================================================================
# Health Check Function
# ================================================================

check_health() {
    local url=$1
    local name=$2
    
    START_TIME=$(date +%s%N)
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -m 10 "$url" 2>/dev/null || echo "000")
    END_TIME=$(date +%s%N)
    RESPONSE_TIME=$(( (END_TIME - START_TIME) / 1000000 ))
    
    ((TOTAL_CHECKS++))
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} $(date +%H:%M:%S) | $name | HTTP $HTTP_CODE | ${RESPONSE_TIME}ms"
        ((SUCCESSFUL_CHECKS++))
        CONSECUTIVE_FAILURES=0
        return 0
    else
        echo -e "${RED}✗${NC} $(date +%H:%M:%S) | $name | HTTP $HTTP_CODE | ${RESPONSE_TIME}ms"
        ((FAILED_CHECKS++))
        ((CONSECUTIVE_FAILURES++))
        return 1
    fi
}

# ================================================================
# Performance Metrics Function
# ================================================================

check_performance() {
    local url="https://$DOMAIN/health"
    
    # Run 5 requests and calculate average
    TIMES=()
    for i in {1..5}; do
        START=$(date +%s%N)
        curl -s -o /dev/null -m 10 "$url" 2>/dev/null || true
        END=$(date +%s%N)
        ELAPSED=$(( (END - START) / 1000000 ))
        TIMES+=($ELAPSED)
    done
    
    # Calculate average
    SUM=0
    for time in "${TIMES[@]}"; do
        SUM=$((SUM + time))
    done
    AVG=$((SUM / ${#TIMES[@]}))
    
    # Find min and max
    MIN=${TIMES[0]}
    MAX=${TIMES[0]}
    for time in "${TIMES[@]}"; do
        [ $time -lt $MIN ] && MIN=$time
        [ $time -gt $MAX ] && MAX=$time
    done
    
    echo -e "${BLUE}Performance: Avg=${AVG}ms, Min=${MIN}ms, Max=${MAX}ms${NC}"
    
    if [ $AVG -lt 500 ]; then
        echo -e "${GREEN}  ✓ Performance: Excellent${NC}"
    elif [ $AVG -lt 1000 ]; then
        echo -e "${YELLOW}  ⚠ Performance: Acceptable${NC}"
    else
        echo -e "${RED}  ✗ Performance: Slow${NC}"
    fi
}

# ================================================================
# API Endpoint Tests
# ================================================================

test_endpoints() {
    echo ""
    echo -e "${BLUE}Testing API Endpoints:${NC}"
    
    # Health endpoint
    check_health "https://$DOMAIN/health" "Health"
    
    # API root
    check_health "https://$DOMAIN/" "API Root"
    
    # Recipes (may require auth)
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/api/recipes")
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "401" ]; then
        echo -e "${GREEN}✓${NC} $(date +%H:%M:%S) | Recipes API | HTTP $HTTP_CODE"
    else
        echo -e "${RED}✗${NC} $(date +%H:%M:%S) | Recipes API | HTTP $HTTP_CODE"
    fi
    
    # Workouts
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN/api/workouts")
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "401" ]; then
        echo -e "${GREEN}✓${NC} $(date +%H:%M:%S) | Workouts API | HTTP $HTTP_CODE"
    else
        echo -e "${RED}✗${NC} $(date +%H:%M:%S) | Workouts API | HTTP $HTTP_CODE"
    fi
    
    echo ""
}

# ================================================================
# Alert Function
# ================================================================

send_alert() {
    local message=$1
    
    echo ""
    echo -e "${RED}╔════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                  ALERT                             ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════╝${NC}"
    echo -e "${RED}$message${NC}"
    echo ""
    
    # Log to file
    echo "[$(date)] ALERT: $message" >> deployment-monitor.log
    
    # Could integrate with notification services here
    # Example: curl webhook, email, Slack, etc.
}

# ================================================================
# Main Monitoring Loop
# ================================================================

echo "Starting continuous monitoring..."
echo "Press Ctrl+C to stop"
echo ""
echo "Configuration:"
echo "  Domain: $DOMAIN"
echo "  Check interval: ${CHECK_INTERVAL}s"
echo "  Alert threshold: $ALERT_THRESHOLD consecutive failures"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Initial comprehensive test
test_endpoints
check_performance

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Continuous monitoring (every ${CHECK_INTERVAL}s):"
echo ""

# Continuous monitoring
while true; do
    check_health "https://$DOMAIN/health" "Health"
    
    # Check if alert threshold reached
    if [ $CONSECUTIVE_FAILURES -ge $ALERT_THRESHOLD ]; then
        send_alert "Application is DOWN! $CONSECUTIVE_FAILURES consecutive failures detected."
        echo "Waiting 60s before continuing..."
        sleep 60
    fi
    
    # Show summary every 10 checks
    if [ $((TOTAL_CHECKS % 10)) -eq 0 ]; then
        echo ""
        echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${BLUE}Summary (last $TOTAL_CHECKS checks):${NC}"
        echo -e "  ${GREEN}Successful: $SUCCESSFUL_CHECKS${NC}"
        echo -e "  ${RED}Failed: $FAILED_CHECKS${NC}"
        UPTIME=$(echo "scale=2; $SUCCESSFUL_CHECKS * 100 / $TOTAL_CHECKS" | bc)
        echo -e "  Uptime: ${UPTIME}%"
        echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        # Run comprehensive test every 10 checks
        test_endpoints
        check_performance
        echo ""
    fi
    
    sleep $CHECK_INTERVAL
done
