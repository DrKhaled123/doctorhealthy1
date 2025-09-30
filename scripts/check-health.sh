#!/bin/bash

# check-health.sh - Verify deployment health

set -e

DOMAIN="${1:-my.doctorhealthy1.com}"
MAX_RETRIES=30
RETRY_INTERVAL=10

echo "🏥 Checking health of ${DOMAIN}..."

for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $i/$MAX_RETRIES..."

    # Try to hit the health endpoint
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" https://${DOMAIN}/health || echo "000")

    if [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Health check passed! (HTTP $HTTP_CODE)"

        # Get detailed health info
        echo ""
        echo "Health details:"
        curl -s https://${DOMAIN}/health | jq '.' || curl -s https://${DOMAIN}/health
        exit 0
    elif [ "$HTTP_CODE" = "503" ]; then
        echo "⚠️  Service unavailable (HTTP $HTTP_CODE) - checking details..."
        curl -s https://${DOMAIN}/health | jq '.' || curl -s https://${DOMAIN}/health
    else
        echo "❌ Health check failed (HTTP $HTTP_CODE)"
    fi

    if [ $i -lt $MAX_RETRIES ]; then
        echo "Waiting ${RETRY_INTERVAL}s before retry..."
        sleep $RETRY_INTERVAL
    fi
done

echo ""
echo "❌ Health check failed after $MAX_RETRIES attempts"
echo "Please check Coolify logs for more details"
exit 1