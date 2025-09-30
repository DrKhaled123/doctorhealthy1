#!/bin/bash

# force-rebuild.sh - Force Coolify to rebuild Docker image

set -e

echo "🔄 Forcing Coolify rebuild..."

# Method 1: Update Dockerfile with timestamp
TIMESTAMP=$(date +%s)
sed -i '' "s/# Build version:.*/# Build version: ${TIMESTAMP}/" Dockerfile

# Method 2: Create/Update .dockerignore to force rebuild
echo "# Last rebuild: ${TIMESTAMP}" >> .dockerignore

# Method 3: Tag commit with rebuild flag
git add Dockerfile .dockerignore
git commit -m "force rebuild: ${TIMESTAMP} [rebuild]"
git push

echo "✅ Changes pushed. Coolify should detect changes and rebuild."
echo ""
echo "If Coolify still doesn't rebuild, use the manual rebuild option:"
echo "1. Go to Coolify dashboard"
echo "2. Navigate to your application"
echo "3. Click 'Force Rebuild' or 'Redeploy'"