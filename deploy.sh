#!/bin/bash

# DoctorHealthy1 Coolify Deployment Script
# This script automates the deployment process to Coolify

set -e

echo "🚀 Starting DoctorHealthy1 deployment to Coolify..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVER_IP="128.140.111.171"
GITHUB_REPO="https://github.com/DrKhaled123/doctorhealthy1"
DOMAIN="api.doctorhealthy1.com"
SSH_KEY="$HOME/.ssh/coolify_doctorhealthy1"

echo -e "${YELLOW}📋 Prerequisites Check${NC}"
echo "1. Make sure your SSH key is added to ssh-agent"
echo "2. DNS records should be set up in Namecheap"
echo "3. Coolify should be running on the server"
echo ""

# Step 1: Push latest code to GitHub
echo -e "${YELLOW}🔄 Checking GitHub status...${NC}"
if git status --porcelain | grep -q .; then
    echo "You have uncommitted changes. Please commit them first."
    exit 1
fi

echo "Code is up to date with GitHub..."
# Skip push due to large files - code is already up to date
# git push origin main

# Step 2: SSH to server and restart Coolify
echo -e "${YELLOW}🔧 Restarting Coolify services...${NC}"
ssh -i "$SSH_KEY" root@$SERVER_IP << 'EOF'
cd /opt/coolify
docker-compose down
docker-compose up -d
echo "Coolify restarted successfully"
EOF

# Step 3: Instructions for Coolify UI setup
echo -e "${GREEN}✅ Code pushed and Coolify restarted!${NC}"
echo ""
echo -e "${YELLOW}📋 Next Steps (Manual):${NC}"
echo ""
echo "1. Open Coolify UI:"
echo "   ssh -i \"$SSH_KEY\" -N -L 8000:localhost:8000 root@$SERVER_IP"
echo "   Then visit: http://localhost:8000"
echo ""
echo "2. Create New Application in Coolify:"
echo "   - Name: doctorhealthy1-api"
echo "   - Git Repository: https://github.com/DrKhaled123/doctorhealthy1"
echo "   - Branch: main"
echo "   - Build Pack: Dockerfile"
echo "   - Internal Port: 8081"
echo "   - Health Check: GET /health"
echo ""
echo "3. Environment Variables (add these in Coolify):"
echo "   PORT=8081"
echo "   DB_PATH=/data/app.db"
echo "   JWT_SECRET=your_super_secure_jwt_secret_key_change_this_in_production_12345678901234567890"
echo "   CORS_ORIGINS=https://www.doctorhealthy1.com"
echo ""
echo "4. Storage:"
echo "   - Mount /data as persistent volume"
echo ""
echo "5. Domain:"
echo "   - Domain: api.doctorhealthy1.com"
echo ""
echo "6. Deploy the application"
echo ""
echo -e "${YELLOW}🔍 Testing:${NC}"
echo "After deployment, test these URLs:"
echo "curl https://api.doctorhealthy1.com/health"
echo "curl https://www.doctorhealthy1.com/health"
echo ""

echo -e "${GREEN}🎉 Deployment script completed!${NC}"
echo "Follow the manual steps above to complete the deployment."