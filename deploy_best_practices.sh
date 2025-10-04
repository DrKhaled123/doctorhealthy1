#!/bin/bash
#
# Comprehensive Deployment Best Practices Script
# Implements expert recommendations for error-free deployments
#

set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/deployment_$(date +%Y%m%d_%H%M%S).log"
ERROR_TRACKER="${SCRIPT_DIR}/deployment_error_tracker.py"
MONITORING_SYSTEM="${SCRIPT_DIR}/deployment_monitoring.py"
PREVENTION_SYSTEM="${SCRIPT_DIR}/deployment_error_prevention.py"

# Logging function
log() {
    local level=$1
    shift
    local message="$*"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case $level in
        "INFO")
            echo -e "${GREEN}[INFO]${NC} $message" | tee -a "$LOG_FILE"
            ;;
        "WARN")
            echo -e "${YELLOW}[WARN]${NC} $message" | tee -a "$LOG_FILE"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message" | tee -a "$LOG_FILE"
            ;;
        "DEBUG")
            echo -e "${BLUE}[DEBUG]${NC} $message" | tee -a "$LOG_FILE"
            ;;
    esac
}

# Error handling function
handle_error() {
    local exit_code=$1
    local line_number=$2
    log "ERROR" "Script failed at line $line_number with exit code $exit_code"
    
    # Record error in tracking system
    if [[ -f "$ERROR_TRACKER" ]]; then
        python3 "$ERROR_TRACKER" --add-error \
            --platform "bash_script" \
            --error-type "script_failure" \
            --root-cause "Script failed at line $line_number" \
            --severity "medium" || true
    fi
    
    exit $exit_code
}

trap 'handle_error $? $LINENO' ERR

# Pre-deployment validation
validate_environment() {
    log "INFO" "🔍 Validating deployment environment..."
    
    # Check required tools
    local required_tools=("curl" "docker" "go" "python3")
    local missing_tools=()
    
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log "ERROR" "Missing required tools: ${missing_tools[*]}"
        return 1
    fi
    
    # Check disk space (require at least 1GB)
    local available_space_kb=$(df . | awk 'NR==2 {print $4}')
    local required_space_kb=1048576  # 1GB in KB
    
    if [[ $available_space_kb -lt $required_space_kb ]]; then
        log "ERROR" "Insufficient disk space. Available: ${available_space_kb}KB, Required: ${required_space_kb}KB"
        return 1
    fi
    
    # Validate Go environment
    if ! go version &> /dev/null; then
        log "ERROR" "Go is not properly configured"
        return 1
    fi
    
    # Check if application is already running
    if pgrep -f "api-key-generator" > /dev/null; then
        log "INFO" "Application is currently running"
    else
        log "WARN" "Application is not running - this may be expected during deployment"
    fi
    
    log "INFO" "✅ Environment validation completed successfully"
    return 0
}

# Nginx configuration validation
validate_nginx_config() {
    log "INFO" "🌐 Validating Nginx configuration..."
    
    # Generate sample Nginx configuration for reference
    local nginx_config_sample="${SCRIPT_DIR}/nginx_sample.conf"
    
    cat > "$nginx_config_sample" << 'EOF'
# Sample Nginx Configuration for Go Application
# Location: /etc/nginx/sites-available/doctorhealthy1

server {
    listen 80;
    server_name my.doctorhealthy1.com;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;

    # Main proxy configuration
    location / {
        proxy_pass http://localhost:8085;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 300s;
        
        # Buffer settings
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }

    # Health check endpoint (bypass rate limiting)
    location /health {
        limit_req off;
        proxy_pass http://localhost:8085;
        proxy_set_header Host $host;
    }

    # Static files (if any)
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Security: Block access to sensitive files
    location ~ /\. {
        deny all;
    }

    # File upload limits
    client_max_body_size 64M;
}

# Redirect HTTP to HTTPS (enable after SSL setup)
# server {
#     listen 80;
#     server_name my.doctorhealthy1.com;
#     return 301 https://$server_name$request_uri;
# }

# HTTPS configuration (enable after SSL setup)
# server {
#     listen 443 ssl http2;
#     server_name my.doctorhealthy1.com;
#
#     ssl_certificate /path/to/certificate.crt;
#     ssl_certificate_key /path/to/private.key;
#     ssl_protocols TLSv1.2 TLSv1.3;
#     ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
#     ssl_prefer_server_ciphers off;
#
#     # Include the same location blocks as HTTP version above
# }
EOF

    log "INFO" "📄 Generated Nginx configuration sample: $nginx_config_sample"
    
    # Validate existing Nginx config if available
    if command -v nginx &> /dev/null; then
        if nginx -t &> /dev/null; then
            log "INFO" "✅ Nginx configuration is valid"
        else
            log "WARN" "⚠️  Nginx configuration has issues. Check with: nginx -t"
        fi
    else
        log "INFO" "ℹ️  Nginx not installed - using alternative proxy"
    fi
    
    return 0
}

# Docker configuration validation
validate_docker_config() {
    log "INFO" "🐳 Validating Docker configuration..."
    
    # Check if Docker is available
    if ! command -v docker &> /dev/null; then
        log "WARN" "Docker not installed - skipping Docker validation"
        return 0
    fi
    
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        log "WARN" "Docker daemon is not running"
        return 0
    fi
    
    # Validate Dockerfile if exists
    if [[ -f "Dockerfile" ]]; then
        # Basic Dockerfile validation
        local dockerfile_issues=()
        
        if ! grep -q "^FROM" Dockerfile; then
            dockerfile_issues+=("Missing FROM instruction")
        fi
        
        if ! grep -q "^WORKDIR" Dockerfile; then
            dockerfile_issues+=("Missing WORKDIR instruction")
        fi
        
        if ! grep -q "^USER" Dockerfile; then
            dockerfile_issues+=("Missing USER instruction (security risk)")
        fi
        
        if grep -q ":latest" Dockerfile; then
            dockerfile_issues+=("Using :latest tag (not recommended)")
        fi
        
        if [[ ${#dockerfile_issues[@]} -gt 0 ]]; then
            log "WARN" "Dockerfile issues found:"
            for issue in "${dockerfile_issues[@]}"; do
                log "WARN" "  - $issue"
            done
        else
            log "INFO" "✅ Dockerfile follows best practices"
        fi
        
        # Test Docker build (dry run if possible)
        log "INFO" "Testing Docker build..."
        if docker build --dry-run -t test-build . &> /dev/null; then
            log "INFO" "✅ Docker build validation passed"
        else
            log "WARN" "⚠️  Docker build may have issues"
        fi
    else
        log "INFO" "ℹ️  No Dockerfile found - not using Docker deployment"
    fi
    
    return 0
}

# Security validation
validate_security() {
    log "INFO" "🔒 Validating security configuration..."
    
    # Check for sensitive files
    local sensitive_patterns=("*.key" "*.pem" "*secret*" "*password*" "*.env")
    local exposed_files=()
    
    for pattern in "${sensitive_patterns[@]}"; do
        if find . -name "$pattern" -type f 2>/dev/null | grep -q .; then
            exposed_files+=("$pattern")
        fi
    done
    
    if [[ ${#exposed_files[@]} -gt 0 ]]; then
        log "WARN" "Potentially sensitive files found: ${exposed_files[*]}"
        log "WARN" "Ensure these files are properly secured and not in version control"
    fi
    
    # Check .gitignore for security patterns
    if [[ -f ".gitignore" ]]; then
        local gitignore_patterns=("*.key" "*.pem" ".env" "secrets/" "config.json")
        local missing_patterns=()
        
        for pattern in "${gitignore_patterns[@]}"; do
            if ! grep -q "$pattern" .gitignore; then
                missing_patterns+=("$pattern")
            fi
        done
        
        if [[ ${#missing_patterns[@]} -gt 0 ]]; then
            log "WARN" "Consider adding to .gitignore: ${missing_patterns[*]}"
        fi
    fi
    
    # Validate file permissions
    local executable_files=()
    while IFS= read -r -d '' file; do
        executable_files+=("$file")
    done < <(find . -name "*.sh" -o -name "*.py" -print0 2>/dev/null)
    
    for file in "${executable_files[@]}"; do
        if [[ ! -x "$file" ]]; then
            log "WARN" "Script file not executable: $file"
            log "INFO" "Run: chmod +x $file"
        fi
    done
    
    log "INFO" "✅ Security validation completed"
    return 0
}

# Performance optimization
optimize_performance() {
    log "INFO" "⚡ Applying performance optimizations..."
    
    # Go build optimizations
    if [[ -f "main.go" ]]; then
        log "INFO" "Building Go application with optimizations..."
        
        # Set Go environment for optimal builds
        export CGO_ENABLED=0
        export GOOS=linux
        export GOARCH=amd64
        
        # Build with optimizations
        go build -ldflags="-w -s" -o main main.go
        
        # Verify build
        if [[ -x "main" ]]; then
            log "INFO" "✅ Optimized build completed successfully"
        else
            log "ERROR" "Build failed"
            return 1
        fi
    fi
    
    # Database optimizations (SQLite)
    local db_files=($(find . -name "*.db" 2>/dev/null))
    for db_file in "${db_files[@]}"; do
        if [[ -f "$db_file" ]]; then
            # Check database file size
            local db_size=$(du -h "$db_file" | cut -f1)
            log "INFO" "Database size: $db_file = $db_size"
            
            # Basic SQLite optimization (VACUUM if size > 100MB)
            local db_size_bytes=$(du -b "$db_file" | cut -f1)
            if [[ $db_size_bytes -gt 104857600 ]]; then  # 100MB
                log "INFO" "Optimizing large database: $db_file"
                sqlite3 "$db_file" "VACUUM; ANALYZE;" 2>/dev/null || true
            fi
        fi
    done
    
    log "INFO" "✅ Performance optimization completed"
    return 0
}

# Run comprehensive error prevention
run_error_prevention() {
    log "INFO" "🔍 Running comprehensive error prevention checks..."
    
    if [[ -f "$PREVENTION_SYSTEM" ]]; then
        if python3 "$PREVENTION_SYSTEM"; then
            log "INFO" "✅ Error prevention checks passed"
            return 0
        else
            log "ERROR" "Error prevention checks failed"
            return 1
        fi
    else
        log "WARN" "Error prevention system not found: $PREVENTION_SYSTEM"
        return 0
    fi
}

# Start monitoring
start_monitoring() {
    log "INFO" "📊 Starting deployment monitoring..."
    
    if [[ -f "$MONITORING_SYSTEM" ]]; then
        # Start monitoring in background
        nohup python3 "$MONITORING_SYSTEM" > "monitoring_$(date +%Y%m%d_%H%M%S).log" 2>&1 &
        local monitor_pid=$!
        
        log "INFO" "✅ Monitoring started (PID: $monitor_pid)"
        echo "$monitor_pid" > "${SCRIPT_DIR}/monitoring.pid"
        
        # Give monitoring system time to start
        sleep 5
        
        # Verify monitoring is running
        if kill -0 "$monitor_pid" 2>/dev/null; then
            log "INFO" "✅ Monitoring system is running"
        else
            log "WARN" "⚠️  Monitoring system may have failed to start"
        fi
    else
        log "WARN" "Monitoring system not found: $MONITORING_SYSTEM"
    fi
}

# Stop monitoring
stop_monitoring() {
    log "INFO" "🛑 Stopping deployment monitoring..."
    
    local pid_file="${SCRIPT_DIR}/monitoring.pid"
    if [[ -f "$pid_file" ]]; then
        local monitor_pid=$(cat "$pid_file")
        if kill -0 "$monitor_pid" 2>/dev/null; then
            kill "$monitor_pid"
            log "INFO" "✅ Monitoring stopped (PID: $monitor_pid)"
        else
            log "INFO" "ℹ️  Monitoring process was not running"
        fi
        rm -f "$pid_file"
    else
        log "INFO" "ℹ️  No monitoring PID file found"
    fi
}

# Pre-deployment checklist
run_pre_deployment_checklist() {
    log "INFO" "📋 Running pre-deployment checklist..."
    
    local checks=(
        "validate_environment"
        "validate_nginx_config" 
        "validate_docker_config"
        "validate_security"
        "optimize_performance"
        "run_error_prevention"
    )
    
    local passed=0
    local total=${#checks[@]}
    
    for check in "${checks[@]}"; do
        log "INFO" "Running: $check"
        if $check; then
            ((passed++))
        else
            log "ERROR" "Check failed: $check"
        fi
    done
    
    local success_rate=$((passed * 100 / total))
    
    log "INFO" "📊 Pre-deployment checklist results: $passed/$total checks passed ($success_rate%)"
    
    if [[ $success_rate -ge 80 ]]; then
        log "INFO" "✅ System ready for deployment (score: $success_rate%)"
        return 0
    else
        log "ERROR" "❌ System not ready for deployment (score: $success_rate%)"
        return 1
    fi
}

# Deploy to production
deploy_to_production() {
    log "INFO" "🚀 Starting production deployment..."
    
    # Ensure application is built
    if [[ ! -f "main" ]]; then
        log "INFO" "Building application..."
        go build -o main main.go
    fi
    
    # Stop existing application gracefully
    if pgrep -f "api-key-generator" > /dev/null; then
        log "INFO" "Stopping existing application..."
        pkill -f "api-key-generator" || true
        sleep 3
    fi
    
    # Start application
    log "INFO" "Starting application..."
    nohup ./main > "app_$(date +%Y%m%d_%H%M%S).log" 2>&1 &
    local app_pid=$!
    
    echo "$app_pid" > "${SCRIPT_DIR}/app.pid"
    log "INFO" "✅ Application started (PID: $app_pid)"
    
    # Wait for application to start
    sleep 5
    
    # Verify application is running
    if kill -0 "$app_pid" 2>/dev/null; then
        log "INFO" "✅ Application is running"
    else
        log "ERROR" "❌ Application failed to start"
        return 1
    fi
    
    # Test application health
    local max_attempts=10
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -s --connect-timeout 5 --max-time 10 "http://localhost:8085/health" > /dev/null; then
            log "INFO" "✅ Application health check passed"
            break
        else
            log "WARN" "Health check attempt $attempt/$max_attempts failed, retrying..."
            sleep 2
            ((attempt++))
        fi
    done
    
    if [[ $attempt -gt $max_attempts ]]; then
        log "ERROR" "❌ Application health checks failed"
        return 1
    fi
    
    log "INFO" "✅ Production deployment completed successfully"
    return 0
}

# Show deployment status
show_deployment_status() {
    log "INFO" "📊 Current deployment status:"
    
    # Application status
    if pgrep -f "api-key-generator" > /dev/null; then
        local app_pid=$(pgrep -f "api-key-generator")
        log "INFO" "✅ Application running (PID: $app_pid)"
        
        # Test health endpoint
        if curl -s --connect-timeout 3 "http://localhost:8085/health" > /dev/null; then
            log "INFO" "✅ Application responding to health checks"
        else
            log "WARN" "⚠️  Application not responding to health checks"
        fi
    else
        log "WARN" "⚠️  Application not running"
    fi
    
    # Monitoring status
    local pid_file="${SCRIPT_DIR}/monitoring.pid"
    if [[ -f "$pid_file" ]]; then
        local monitor_pid=$(cat "$pid_file")
        if kill -0 "$monitor_pid" 2>/dev/null; then
            log "INFO" "✅ Monitoring system running (PID: $monitor_pid)"
        else
            log "WARN" "⚠️  Monitoring system not running"
        fi
    else
        log "INFO" "ℹ️  No monitoring system detected"
    fi
    
    # Port usage
    local ports=("8085" "80" "443")
    for port in "${ports[@]}"; do
        if lsof -i ":$port" &> /dev/null; then
            local process=$(lsof -i ":$port" | awk 'NR==2 {print $1}')
            log "INFO" "✅ Port $port in use by: $process"
        else
            log "INFO" "ℹ️  Port $port available"
        fi
    done
    
    # Recent logs
    if [[ -f "$LOG_FILE" ]]; then
        log "INFO" "📄 Recent deployment log: $LOG_FILE"
    fi
}

# Main function
main() {
    log "INFO" "🚀 Starting Comprehensive Deployment Best Practices Script"
    log "INFO" "📄 Log file: $LOG_FILE"
    log "INFO" "=" "=================================================================="
    
    case "${1:-}" in
        "validate")
            run_pre_deployment_checklist
            ;;
        "deploy")
            if run_pre_deployment_checklist; then
                start_monitoring
                deploy_to_production
            else
                log "ERROR" "Pre-deployment validation failed. Aborting deployment."
                exit 1
            fi
            ;;
        "monitor")
            start_monitoring
            ;;
        "stop-monitor")
            stop_monitoring
            ;;
        "status")
            show_deployment_status
            ;;
        "full")
            run_pre_deployment_checklist
            start_monitoring
            deploy_to_production
            show_deployment_status
            ;;
        *)
            echo "Usage: $0 {validate|deploy|monitor|stop-monitor|status|full}"
            echo ""
            echo "Commands:"
            echo "  validate      - Run pre-deployment validation checks"
            echo "  deploy        - Run validation and deploy to production" 
            echo "  monitor       - Start monitoring system"
            echo "  stop-monitor  - Stop monitoring system"
            echo "  status        - Show current deployment status"
            echo "  full          - Complete deployment workflow"
            echo ""
            exit 1
            ;;
    esac
    
    log "INFO" "✅ Script completed successfully"
}

# Run main function with all arguments
main "$@"