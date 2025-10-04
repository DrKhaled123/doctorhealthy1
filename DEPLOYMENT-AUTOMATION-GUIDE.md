# 🚀 Deployment Automation Guide

Complete guide for automated deployment, monitoring, and rollback procedures.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Deployment Workflow](#deployment-workflow)
4. [Scripts Reference](#scripts-reference)
5. [Monitoring](#monitoring)
6. [Rollback Procedures](#rollback-procedures)
7. [Troubleshooting](#troubleshooting)

---

## 🎯 Overview

This project includes a comprehensive deployment automation system with:

- **Pre-deployment validation** - Automated checks before deployment
- **Automated deployment** - One-command deployment to Coolify
- **Health monitoring** - Continuous application health checks
- **Quick rollback** - Fast recovery from failed deployments

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Deployment Pipeline                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Pre-Deployment Checks                                   │
│     ├─ Code Quality (gofmt, go vet, staticcheck)           │
│     ├─ Tests (unit tests, coverage >50%)                   │
│     ├─ Security (secret scanner, dependencies)             │
│     ├─ Build Validation (binary + Docker)                  │
│     └─ Git Status                                           │
│                                                              │
│  2. Deployment Execution                                    │
│     ├─ Update Environment Variables                        │
│     ├─ Trigger Coolify Deployment                          │
│     ├─ Monitor Deployment Progress                         │
│     └─ Post-Deployment Health Check                        │
│                                                              │
│  3. Continuous Monitoring                                   │
│     ├─ Health Checks (every 30s)                           │
│     ├─ Performance Metrics                                 │
│     ├─ API Endpoint Tests                                  │
│     └─ Alert on Failures                                   │
│                                                              │
│  4. Rollback (if needed)                                    │
│     ├─ Capture Current State                               │
│     ├─ Git Revert                                          │
│     ├─ Redeploy Previous Version                           │
│     └─ Verify Recovery                                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Prerequisites

### Required Tools

```bash
# Check if you have all required tools
which go        # Go 1.22+
which git       # Git
which curl      # cURL
which docker    # Docker (optional but recommended)
which jq        # jq (for JSON parsing)
```

### Install Missing Tools

```bash
# macOS
brew install go git curl docker jq

# Linux (Ubuntu/Debian)
sudo apt-get install golang git curl docker.io jq

# Linux (CentOS/RHEL)
sudo yum install golang git curl docker jq
```

### Environment Configuration

Create/verify `.env.example`:

```bash
# Required environment variables
JWT_SECRET=your-secret-here
PORT=8081
DB_PATH=/app/data/app.db
ENV=production
LOG_LEVEL=warn
ALLOWED_ORIGIN=https://my.doctorhealthy1.com
```

### Coolify Configuration

Set your Coolify API token as environment variable:

```bash
export COOLIFY_API_TOKEN="4|jdTX2lUb2q6IOrwNGkHyQBCO74JJeeRHZVvFNwgI6b376a50"
```

Or it will use the default token in the scripts.

---

## 🚀 Deployment Workflow

### Quick Start (Recommended)

```bash
# Complete deployment in one command
make -f Makefile.testing deploy
```

This will:
1. ✓ Run all pre-deployment checks
2. ✓ Validate code quality and tests
3. ✓ Scan for security issues
4. ✓ Deploy to Coolify
5. ✓ Verify deployment health

### Step-by-Step Deployment

#### Step 1: Pre-Deployment Validation

```bash
# Run validation checks
make -f Makefile.testing pre-deploy

# Or run script directly
bash scripts/pre-deploy-check.sh
```

**What it checks:**
- ✓ Code formatting (gofmt)
- ✓ Go vet analysis
- ✓ Static analysis (staticcheck)
- ✓ Unit tests passing
- ✓ Test coverage >50%
- ✓ No secrets in code
- ✓ No .env in git
- ✓ Dependencies audit
- ✓ Configuration files
- ✓ Build success
- ✓ Docker image builds
- ✓ Documentation exists
- ✓ Git status

**Output Example:**
```
╔════════════════════════════════════════════════════╗
║      PRE-DEPLOYMENT VALIDATION                     ║
╚════════════════════════════════════════════════════╝

1. Code Quality Checks

→ Go code formatting... ✓ PASS
→ Go vet... ✓ PASS
→ Static analysis... ✓ PASS

...

╔════════════════════════════════════════════════════╗
║    ✓ ALL CHECKS PASSED - READY TO DEPLOY          ║
╚════════════════════════════════════════════════════╝
```

#### Step 2: Deploy

```bash
# Deploy to Coolify
bash scripts/deploy-to-coolify.sh
```

**Deployment Process:**
1. Runs pre-deployment checks
2. Tests Coolify API connection
3. Gets current application status
4. Creates backup point (deployment-info.json)
5. Updates environment variables
6. Triggers deployment
7. Monitors progress
8. Runs post-deployment health checks
9. Runs smoke tests

**Output Example:**
```
╔════════════════════════════════════════════════════╗
║      COOLIFY DEPLOYMENT AUTOMATION                ║
╚════════════════════════════════════════════════════╝

Step 1: Pre-deployment Validation
✓ Pre-deployment checks passed

Step 2: API Connection
✓ Coolify API connection successful

...

╔════════════════════════════════════════════════════╗
║           DEPLOYMENT SUMMARY                       ║
╚════════════════════════════════════════════════════╝

Domain:          https://my.doctorhealthy1.com
Deployment ID:   abc123
Commit:          a1b2c3d
Branch:          main
Deployed:        2025-10-03 14:30:00
Health Status:   200
```

#### Step 3: Monitor

```bash
# Start continuous monitoring
make -f Makefile.testing monitor

# Or run script directly
bash scripts/monitor-deployment.sh
```

**Monitoring Features:**
- Health checks every 30 seconds
- Performance metrics (avg, min, max response time)
- API endpoint tests
- Automatic alerts on failures
- Summary statistics

**Output Example:**
```
╔════════════════════════════════════════════════════╗
║      POST-DEPLOYMENT MONITORING                    ║
╚════════════════════════════════════════════════════╝

Configuration:
  Domain: my.doctorhealthy1.com
  Check interval: 30s
  Alert threshold: 3 consecutive failures

Continuous monitoring:

✓ 14:30:15 | Health | HTTP 200 | 45ms
✓ 14:30:45 | Health | HTTP 200 | 52ms
✓ 14:31:15 | Health | HTTP 200 | 48ms

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Summary (last 10 checks):
  Successful: 10
  Failed: 0
  Uptime: 100.00%
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## 📚 Scripts Reference

### `scripts/pre-deploy-check.sh`

**Purpose:** Comprehensive pre-deployment validation

**Usage:**
```bash
bash scripts/pre-deploy-check.sh
```

**Exit Codes:**
- `0` - All checks passed
- `1` - Critical failures detected (deployment blocked)

**Features:**
- 7 validation categories
- Critical vs non-critical checks
- Color-coded output
- Interactive confirmation for warnings

---

### `scripts/deploy-to-coolify.sh`

**Purpose:** Automated deployment to Coolify platform

**Usage:**
```bash
# With default configuration
bash scripts/deploy-to-coolify.sh

# With custom API token
COOLIFY_API_TOKEN="your-token" bash scripts/deploy-to-coolify.sh
```

**Environment Variables:**
- `COOLIFY_API_TOKEN` - API authentication token
- `COOLIFY_URL` - Coolify instance URL (default: http://128.140.111.171:8000)
- `APP_UUID` - Application UUID (default: hcw0gc8wcwk440gw4c88408o)
- `DOMAIN` - Application domain (default: my.doctorhealthy1.com)

**Generated Files:**
- `deployment-info.json` - Deployment metadata (timestamp, commit, branch)

**Exit Codes:**
- `0` - Deployment successful
- `1` - Deployment failed

---

### `scripts/monitor-deployment.sh`

**Purpose:** Continuous health monitoring

**Usage:**
```bash
# Default monitoring (30s interval)
bash scripts/monitor-deployment.sh

# Custom interval and threshold
CHECK_INTERVAL=60 ALERT_THRESHOLD=5 bash scripts/monitor-deployment.sh
```

**Environment Variables:**
- `DOMAIN` - Application domain (default: my.doctorhealthy1.com)
- `CHECK_INTERVAL` - Seconds between checks (default: 30)
- `ALERT_THRESHOLD` - Consecutive failures before alert (default: 3)

**Generated Files:**
- `deployment-monitor.log` - Alert log file

**Controls:**
- Press `Ctrl+C` to stop monitoring

---

### `scripts/rollback-deployment.sh`

**Purpose:** Quick rollback to previous working version

**Usage:**
```bash
bash scripts/rollback-deployment.sh
```

**Rollback Options:**
1. **Automatic revert** - Reverts to previous commit (git revert HEAD)
2. **Specific commit** - Reverts to a specific commit you choose
3. **Manual rollback** - Use Coolify UI manually

**Generated Files:**
- `rollback-info.json` - Rollback metadata (timestamp, failed commit, rolled back by)

**Interactive Prompts:**
- Confirms rollback decision
- Asks for push confirmation
- Asks for deployment trigger

---

## 📊 Monitoring

### Real-Time Monitoring

```bash
# Start monitoring
make -f Makefile.testing monitor
```

### Check Logs

```bash
# View deployment monitor log
tail -f deployment-monitor.log

# View application logs (on Coolify)
# Go to: http://128.140.111.171:8000/applications/hcw0gc8wcwk440gw4c88408o/logs
```

### Performance Metrics

Monitoring includes:
- **Response Time**: Average, min, max
- **Uptime**: Percentage of successful checks
- **Endpoint Health**: Multiple API endpoints tested
- **Performance Rating**: Excellent (<500ms), Acceptable (<1000ms), Slow (>1000ms)

### Alerts

Automatic alerts trigger when:
- 3 consecutive health check failures (default threshold)
- Logged to `deployment-monitor.log`
- Can be integrated with notification services (Slack, email, etc.)

---

## 🔄 Rollback Procedures

### When to Rollback

Rollback if you observe:
- ❌ Application not responding (HTTP 500/502/503)
- ❌ High error rate in logs
- ❌ Database migration failure
- ❌ Critical functionality broken
- ❌ Security vulnerability introduced

### Quick Rollback

```bash
# Start rollback procedure
make -f Makefile.testing rollback

# Or run script directly
bash scripts/rollback-deployment.sh
```

### Rollback Process

1. **Confirm Decision**
   - Script will ask for confirmation
   - Type `yes` to proceed

2. **Choose Rollback Method**
   - Option 1: Automatic (revert to previous commit)
   - Option 2: Specific commit
   - Option 3: Manual via Coolify UI

3. **Push Changes**
   - Confirm git push of rollback commit

4. **Trigger Deployment**
   - Confirm automatic deployment trigger

5. **Verify Recovery**
   - Script waits 60s
   - Tests health endpoint
   - Reports success/failure

### Manual Rollback (Coolify UI)

If automated rollback fails:

1. Go to Coolify: http://128.140.111.171:8000
2. Navigate to Applications → Your App
3. Click "Deployments" tab
4. Find previous successful deployment
5. Click "Redeploy" button
6. Wait for deployment to complete
7. Verify application health

---

## 🔧 Troubleshooting

### Common Issues

#### Pre-Deployment Checks Fail

**Issue:** "Code is not formatted"
```bash
# Fix: Format code
gofmt -s -w .
```

**Issue:** "Unit tests failing"
```bash
# Fix: Run tests to see failures
go test -v ./...

# Debug specific test
go test -v -run TestName ./...
```

**Issue:** "Secrets found in code"
```bash
# Fix: Check what was found
bash scripts/check-secrets.sh

# Remove secrets from code
# Move to environment variables
```

#### Deployment Fails

**Issue:** "Cannot connect to Coolify API"
```bash
# Check: Is Coolify running?
curl http://128.140.111.171:8000/api/health

# Check: Is API token correct?
echo $COOLIFY_API_TOKEN
```

**Issue:** "Deployment timeout"
```bash
# Check Coolify logs in UI
# Look for build errors or resource issues

# Try manual deployment
# Go to Coolify UI → Applications → Deploy
```

#### Health Checks Fail

**Issue:** "HTTP 502 Bad Gateway"
```bash
# Application may be starting up
# Wait 2-3 minutes and check again

# Check application logs in Coolify
```

**Issue:** "HTTP 500 Internal Server Error"
```bash
# Check application logs
# Look for:
# - Database connection errors
# - Missing environment variables
# - Panic/crash logs

# Verify environment variables are set correctly
```

#### Monitoring Issues

**Issue:** "Connection timeout"
```bash
# Check: Is domain resolving?
dig my.doctorhealthy1.com

# Check: Is application running?
curl -I https://my.doctorhealthy1.com/health

# Check: SSL certificate valid?
curl -vI https://my.doctorhealthy1.com/health 2>&1 | grep -i ssl
```

### Debug Mode

Run scripts with debug output:

```bash
# Enable bash debug mode
bash -x scripts/deploy-to-coolify.sh

# Enable verbose curl
curl -v https://my.doctorhealthy1.com/health
```

### Get Help

1. **Check Logs:**
   - Coolify UI: http://128.140.111.171:8000/applications/hcw0gc8wcwk440gw4c88408o/logs
   - Monitor log: `cat deployment-monitor.log`
   - Git log: `git log --oneline -10`

2. **Verify Configuration:**
   ```bash
   # Check environment variables
   cat coolify-env-vars.txt
   
   # Check deployment info
   cat deployment-info.json
   
   # Check git status
   git status
   git log --oneline -5
   ```

3. **Test Manually:**
   ```bash
   # Test API directly
   curl https://my.doctorhealthy1.com/health
   curl https://my.doctorhealthy1.com/api/recipes
   
   # Check DNS
   nslookup my.doctorhealthy1.com
   
   # Check SSL
   openssl s_client -connect my.doctorhealthy1.com:443 -servername my.doctorhealthy1.com
   ```

---

## 🎯 Best Practices

### Before Deployment

1. ✅ Always run pre-deployment checks
2. ✅ Commit all changes to git
3. ✅ Push to remote repository
4. ✅ Review recent changes
5. ✅ Have rollback plan ready

### During Deployment

1. ✅ Monitor deployment progress
2. ✅ Watch for errors in logs
3. ✅ Don't interrupt deployment
4. ✅ Wait for health checks

### After Deployment

1. ✅ Run smoke tests
2. ✅ Monitor for 30 minutes
3. ✅ Check error rates
4. ✅ Verify critical functionality
5. ✅ Document any issues

### Rollback Strategy

1. ✅ Rollback immediately if critical issues
2. ✅ Document reason for rollback
3. ✅ Fix issues before redeployment
4. ✅ Test fix locally first
5. ✅ Re-run pre-deployment checks

---

## 📖 Additional Resources

- [Coolify Manual Steps](COOLIFY-MANUAL-STEPS.md)
- [Coolify Domain Guide](COOLIFY-DOMAIN-GUIDE.md)
- [Testing Documentation](Makefile.testing)
- [Deployment Checklist](DEPLOYMENT-CHECKLIST.md)

---

## 🆘 Emergency Contacts

**If deployment fails and you need immediate rollback:**

```bash
# Quick rollback command
bash scripts/rollback-deployment.sh
```

**If rollback fails:**

1. Access Coolify UI: http://128.140.111.171:8000
2. Go to Applications → Your App → Deployments
3. Find last working deployment
4. Click "Redeploy"

**If everything is down:**

1. Check server status
2. Verify DNS is resolving
3. Check Coolify is running
4. Contact server administrator

---

**Last Updated:** October 3, 2025  
**Version:** 1.0.0  
**Maintainer:** Development Team
