#!/bin/bash

# DoctorHealthy1 - Quick Deployment Fix Script
# This script applies critical fixes for production deployment

set -e

echo "🔧 DoctorHealthy1 Deployment Fix Script"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "main.go not found. Please run this script from the project root."
    exit 1
fi

print_status "Project root directory confirmed"

# Check for frontend directory
echo ""
echo "1️⃣  Checking frontend directory..."
if [ ! -d "frontend" ]; then
    print_error "Frontend directory not found!"
    echo "   This will cause 'boxes with no functions' issue."
    echo "   Please ensure frontend directory exists before deploying."
    exit 1
else
    FILE_COUNT=$(ls -1 frontend | wc -l)
    print_status "Frontend directory found with $FILE_COUNT files"
fi

# Check for required files
echo ""
echo "2️⃣  Checking required files..."
REQUIRED_FILES=("Dockerfile" "deploy.sh" "go.mod" "go.sum")
for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        print_error "$file not found"
        exit 1
    fi
done
print_status "All required files present"

# Verify fixes are applied
echo ""
echo "3️⃣  Verifying fixes..."

# Check CORS configuration
if grep -q "AllowOrigins.*https://api.doctorhealthy1.com" main.go; then
    print_status "CORS configuration updated"
else
    print_warning "CORS may not be properly configured in main.go"
fi

# Check Dockerfile port
if grep -q "EXPOSE 8081" Dockerfile; then
    print_status "Dockerfile port configured to 8081"
else
    print_warning "Dockerfile still uses port 8080 (should be 8081)"
fi

# Check for JWT_SECRET in .env
echo ""
echo "4️⃣  Checking environment configuration..."
if [ -f ".env" ]; then
    if grep -q "JWT_SECRET=" .env; then
        JWT_LEN=$(grep "JWT_SECRET=" .env | cut -d'=' -f2 | wc -c)
        if [ $JWT_LEN -ge 32 ]; then
            print_status "JWT_SECRET configured (${JWT_LEN} characters)"
        else
            print_warning "JWT_SECRET is too short (needs 32+ characters)"
        fi
    else
        print_warning ".env file exists but JWT_SECRET not set"
    fi
else
    print_warning ".env file not found (will use environment variables)"
fi

# Build test
echo ""
echo "5️⃣  Testing build..."
if go build -o /tmp/doctorhealthy1-test .; then
    print_status "Build successful"
    rm -f /tmp/doctorhealthy1-test
else
    print_error "Build failed"
    exit 1
fi

# Summary
echo ""
echo "======================================"
echo "📋 Pre-Deployment Checklist"
echo "======================================"
echo ""
echo "Before deploying to Coolify, ensure:"
echo ""
echo "  [ ] SSL/HTTPS configured in Coolify"
echo "      - Add domain: api.doctorhealthy1.com"
echo "      - Enable automatic SSL (Let's Encrypt)"
echo ""
echo "  [ ] Environment variables set in Coolify:"
echo "      - JWT_SECRET=<32+chars>"
echo "      - ENV=production"
echo "      - PORT=8081"
echo "      - LOG_LEVEL=warn"
echo ""
echo "  [ ] Persistent volume configured:"
echo "      - Mount path: /app/data"
echo "      - Size: 2GB minimum"
echo ""
echo "  [ ] Review DEPLOYMENT-ISSUES-ANALYSIS.md"
echo ""
echo "======================================"
echo ""

read -p "Ready to deploy? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Starting deployment..."
    echo ""
    
    # Make deploy.sh executable if not already
    chmod +x deploy.sh
    
    # Run deployment
    ./deploy.sh
else
    echo ""
    print_warning "Deployment cancelled"
    echo "Review the checklist above and run this script again when ready."
    exit 0
fi
