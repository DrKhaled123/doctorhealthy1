# 🚀 Quick Reference Card - Testing & Deployment

**Quick access to all testing and deployment commands**

---

## 📦 Testing Commands

### Show All Available Targets
```bash
make -f Makefile.testing help
```

### Run Tests

| Command | Description | Time |
|---------|-------------|------|
| `make -f Makefile.testing test-unit` | Run unit tests with coverage | ~30s |
| `make -f Makefile.testing test-integration` | Run E2E integration tests | ~2min |
| `make -f Makefile.testing test-security` | Run security scans | ~1min |
| `make -f Makefile.testing test-performance` | Run benchmarks & load tests | ~2min |
| `make -f Makefile.testing smoke-tests` | Run smoke tests | ~1min |
| `make -f Makefile.testing test-all` | Run ALL tests | ~5min |

### QA-GPT Container

```bash
# Build
docker build -f Dockerfile.qa-gpt -t qa-gpt-tester .

# Run
docker run --rm -v $(pwd)/coverage:/app/coverage qa-gpt-tester

# View results
cat coverage/qa-gpt-report.json | jq .
```

---

## 🚀 Deployment Commands

### Full Automated Deployment
```bash
make -f Makefile.testing deploy
```
This runs pre-checks, deploys to Coolify, and verifies health.

### Step-by-Step Deployment

| Step | Command | Description |
|------|---------|-------------|
| 1 | `make -f Makefile.testing pre-deploy` | Validate before deploying |
| 2 | `bash scripts/deploy-to-coolify.sh` | Deploy to Coolify |
| 3 | `make -f Makefile.testing monitor` | Monitor health |
| 4 | `make -f Makefile.testing rollback` | Rollback if needed |

---

## 📊 Monitoring Commands

### Start Monitoring
```bash
# Default (30s interval)
make -f Makefile.testing monitor

# Custom interval
CHECK_INTERVAL=60 bash scripts/monitor-deployment.sh
```

### Check Logs
```bash
# Monitor log
tail -f deployment-monitor.log

# Application logs (Coolify)
# Go to: http://128.140.111.171:8000/applications/hcw0gc8wcwk440gw4c88408o/logs
```

---

## 🔄 Rollback Commands

### Quick Rollback
```bash
make -f Makefile.testing rollback
```

### Manual Rollback (Coolify UI)
1. Go to: http://128.140.111.171:8000
2. Navigate to Applications → Your App
3. Click "Deployments" tab
4. Find previous successful deployment
5. Click "Redeploy"

---

## 🧹 Cleanup Commands

```bash
# Clean all test artifacts
make -f Makefile.testing clean

# Clean and run fresh tests
make -f Makefile.testing clean test-all
```

---

## 🔐 Environment Variables

### Set Coolify API Token
```bash
export COOLIFY_API_TOKEN="your-token-here"
```

### Check Environment
```bash
# View all configured environment variables
cat coolify-env-vars.txt
```

---

## 📁 Important Files

| File | Purpose |
|------|---------|
| `Makefile.testing` | Testing orchestration (30+ targets) |
| `scripts/e2e-tests.sh` | E2E test suite (10 suites) |
| `scripts/smoke-tests.sh` | Manual smoke tests |
| `scripts/check-secrets.sh` | Secret scanner |
| `scripts/load-test.sh` | Load testing |
| `scripts/pre-deploy-check.sh` | Pre-deployment validation |
| `scripts/deploy-to-coolify.sh` | Automated deployment |
| `scripts/monitor-deployment.sh` | Health monitoring |
| `scripts/rollback-deployment.sh` | Rollback procedure |
| `.github/workflows/comprehensive-testing.yml` | CI/CD pipeline |
| `Dockerfile.qa-gpt` | QA-GPT container |

---

## 📊 Reports & Artifacts

### Test Reports Location
```bash
# Coverage reports
coverage/
├── unit.html           # Unit test coverage
├── unit-coverage.out   # Coverage profile
├── integration-results.txt
├── security-results.txt
└── performance-results.txt

# Load test results
coverage/load-tests/
├── health-plot.html
├── recipes-plot.html
├── mixed-plot.html
└── spike-plot.html

# QA-GPT report
coverage/qa-gpt-report.json
```

---

## 🔗 Quick Links

- **Coolify Dashboard:** http://128.140.111.171:8000
- **Application:** https://my.doctorhealthy1.com
- **Health Check:** https://my.doctorhealthy1.com/health
- **API Root:** https://my.doctorhealthy1.com/

---

## 🆘 Emergency Commands

### If Deployment Fails
```bash
# Quick rollback
make -f Makefile.testing rollback
```

### If Application is Down
```bash
# Check health
curl -I https://my.doctorhealthy1.com/health

# Check if server is reachable
ping 128.140.111.171

# Check DNS
dig my.doctorhealthy1.com
```

### If Tests Fail
```bash
# Run specific test
go test -v -run TestName ./...

# Check test coverage
go test -cover ./...

# View detailed errors
make -f Makefile.testing test-unit | tee test-output.log
```

---

## 📖 Documentation

- [Comprehensive Testing Guide](TESTING-FRAMEWORK-IMPLEMENTATION-COMPLETE.md)
- [Deployment Automation Guide](DEPLOYMENT-AUTOMATION-GUIDE.md)
- [Coolify Manual Steps](COOLIFY-MANUAL-STEPS.md)
- [Coolify Domain Guide](COOLIFY-DOMAIN-GUIDE.md)

---

## ✅ Pre-Flight Checklist

Before deploying:
- [ ] All tests passing: `make -f Makefile.testing test-all`
- [ ] No secrets in code: `bash scripts/check-secrets.sh`
- [ ] Code formatted: `gofmt -l . | wc -l` (should be 0)
- [ ] Changes committed: `git status`
- [ ] Pre-deploy checks pass: `make -f Makefile.testing pre-deploy`

After deploying:
- [ ] Health check passes: `curl https://my.doctorhealthy1.com/health`
- [ ] API responding: `curl https://my.doctorhealthy1.com/`
- [ ] Monitor for 30 mins: `make -f Makefile.testing monitor`
- [ ] Document deployment: Update `deployment-info.json`

---

## 🎯 Common Workflows

### **Development Workflow**
```bash
# 1. Make changes
# 2. Run tests
make -f Makefile.testing test-unit
# 3. Commit
git commit -m "..."
# 4. Push (triggers CI/CD)
git push
```

### **Deployment Workflow**
```bash
# 1. Validate
make -f Makefile.testing pre-deploy
# 2. Deploy
make -f Makefile.testing deploy
# 3. Monitor
make -f Makefile.testing monitor
```

### **Troubleshooting Workflow**
```bash
# 1. Check health
curl https://my.doctorhealthy1.com/health
# 2. Check logs
tail -f deployment-monitor.log
# 3. If broken, rollback
make -f Makefile.testing rollback
```

---

**Print this card and keep it nearby for quick reference!**

---

**Last Updated:** October 3, 2025  
**Version:** 1.0.0
