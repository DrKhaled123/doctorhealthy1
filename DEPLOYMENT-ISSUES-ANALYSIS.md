# Deployment Issues Analysis & Solutions

## 🚨 Critical Issues Identified

Based on your symptoms:
- ❌ "Server not found"
- ❌ "Webpage not secure"
- ❌ "Boxes with no functions"
- ❌ "No generation for meals or workouts"
- ❌ "Clicks not working"

---

## 🔴 Issue #1: HTTPS/SSL Certificate Missing (CRITICAL)

### Problem
**Line 92 in `internal/handlers/pdf_handler.go`:**
```go
Secure:   false,  // Set to true in production with HTTPS
```

Your cookies are not secure, and the browser warns "webpage not secure" because:
1. No SSL/TLS certificate configured in Coolify
2. Application running on HTTP instead of HTTPS
3. Mixed content warnings blocking API calls

### Impact
- 🔴 Browser security warnings
- 🔴 Cookies may not work (especially in Chrome/Safari)
- 🔴 API calls from frontend may be blocked
- 🔴 User trust issues

### Solution

#### Step 1: Configure SSL in Coolify
```bash
# In your Coolify dashboard:
1. Go to Application Settings → hcw0gc8wcwk440gw4c88408o
2. Navigate to "Domains" section
3. Add domain: api.doctorhealthy1.com
4. Enable "Automatic SSL" (Let's Encrypt)
5. Wait for certificate provisioning (~2 minutes)
```

#### Step 2: Update Cookie Security
```go
// In internal/handlers/pdf_handler.go line 92
Secure:   true,  // ENABLE for HTTPS
HttpOnly: true,
SameSite: http.SameSiteLaxMode,
```

#### Step 3: Force HTTPS Redirect
Add to `main.go` after line 95:
```go
// Force HTTPS in production
e.Pre(echomiddleware.HTTPSRedirect())
```

---

## 🔴 Issue #2: CORS Configuration Problem (CRITICAL)

### Problem
**Line 95 in `main.go`:**
```go
e.Use(echomiddleware.CORS())  // Using default CORS - TOO PERMISSIVE
```

**Line 59 in `internal/middleware/middleware.go`:**
```go
AllowOrigins: []string{"*"},  // Allows ALL origins - SECURITY RISK
```

This causes:
- ❌ Frontend-backend communication fails
- ❌ Preflight OPTIONS requests fail
- ❌ Cookies not sent with requests
- ❌ "Boxes with no functions" (API calls blocked)

### Impact
- 🔴 API calls from browser fail
- 🔴 Login doesn't work
- 🔴 No meal/workout generation (API blocked)
- 🔴 Clicks do nothing (AJAX fails)

### Solution

#### Update `main.go` (Replace line 95):
```go
// REPLACE:
e.Use(echomiddleware.CORS())

// WITH:
e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
	AllowOrigins: []string{
		"https://api.doctorhealthy1.com",
		"https://doctorhealthy1.com",
		"https://www.doctorhealthy1.com",
		"http://localhost:3000",  // For local development
	},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	},
	AllowHeaders: []string{
		echo.HeaderOrigin,
		echo.HeaderContentType,
		echo.HeaderAccept,
		echo.HeaderAuthorization,
		"X-API-Key",
		"X-Requested-With",
	},
	AllowCredentials: true,
	MaxAge:           300,
}))
```

---

## 🟠 Issue #3: Port Configuration Mismatch

### Problem
**Dockerfile exposes port 8080:**
```dockerfile
EXPOSE 8080
```

**But deploy.sh shows Internal Port: 8081:**
```bash
Internal Port: 8081
```

**And default config uses 8080:**
```go
Port: getEnv("PORT", "8080"),
```

### Impact
- 🟠 Server starts but Coolify can't reach it
- 🟠 Health checks fail
- 🟠 "Server not found" intermittently

### Solution

#### Option A: Update Dockerfile (Recommended)
```dockerfile
# Change line in Dockerfile:
EXPOSE 8081

# Update HEALTHCHECK:
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1
```

#### Option B: Set PORT Environment Variable in Coolify
```bash
# In Coolify dashboard → Environment Variables:
PORT=8081
```

---

## 🟠 Issue #4: Missing Environment Variables in Production

### Problem
Configuration expects these but they may not be set in Coolify:

```go
// Required
JWT_SECRET=<min-32-chars>  // If missing, server won't start

// Missing = Default values used
PORT=8080                   // Should be 8081
HOST=0.0.0.0               // Correct
DB_PATH=./data/app.db      // May not have write permissions
ENV=production             // Important for logging
LOG_LEVEL=info             // Default OK
CORS_ORIGINS=              // Not used (hardcoded *)
```

### Impact
- 🟠 Server starts with wrong configuration
- 🟠 Database writes may fail
- 🟠 Debug logs in production (performance hit)
- 🟠 Security issues

### Solution

#### Set Environment Variables in Coolify:
```bash
# Required
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-abc123

# Production settings
ENV=production
PORT=8081
HOST=0.0.0.0
DB_PATH=/app/data/app.db
LOG_LEVEL=warn

# Security
API_KEY_PREFIX=dh_
API_KEY_LENGTH=32
RATE_LIMIT=100

# Optional (if you have frontend domain)
ALLOWED_ORIGIN=https://doctorhealthy1.com
```

---

## 🟠 Issue #5: Database Persistence Not Configured

### Problem
**Dockerfile creates data directory:**
```dockerfile
RUN mkdir -p /app/data && chown -R appuser:appuser /app/data
```

But Coolify doesn't have a persistent volume mounted. On container restart:
- ❌ All data lost
- ❌ API keys deleted
- ❌ User data gone
- ❌ Appears as "server not found" (data inconsistency)

### Solution

#### Configure Persistent Volume in Coolify:
```bash
1. Go to Application → Storage
2. Add Volume:
   - Name: app-data
   - Mount Path: /app/data
   - Size: 2GB
3. Redeploy application
```

---

## 🟡 Issue #6: Frontend Not Found (Static Files Missing)

### Problem
Frontend files are embedded at build time:
```go
//go:embed frontend/*
var embeddedFrontend embed.FS
```

But if frontend directory is missing during Docker build:
- ❌ No HTML files served
- ❌ "Boxes with no functions" (no UI loaded)
- ❌ Clicks don't work (no JavaScript loaded)

### Solution

#### Verify frontend directory exists in Docker build:
```dockerfile
# Add to Dockerfile before build step:
# List frontend files for debugging
RUN ls -la frontend/ || echo "WARNING: frontend directory not found!"

# Ensure frontend is copied
COPY frontend/ ./frontend/
```

#### Update Dockerfile (add after line 27):
```dockerfile
# Copy source code
COPY . .

# ADD THIS: Verify frontend exists
RUN ls -la frontend/ && echo "✅ Frontend files found" || (echo "❌ Frontend missing!" && exit 1)

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .
```

---

## 🟡 Issue #7: Health Check Configuration Issues

### Problem
**Dockerfile health check:**
```dockerfile
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
```

Issues:
1. Wrong port (8080 vs 8081)
2. Short timeout (3s may not be enough)
3. Only 3 retries (server marked dead quickly)

### Impact
- 🟡 Coolify marks server as "unhealthy"
- 🟡 Container restarts frequently
- 🟡 "Server not found" during restarts

### Solution

#### Update Dockerfile:
```dockerfile
# Increase resilience
HEALTHCHECK --interval=15s --timeout=10s --start-period=30s --retries=5 \
    CMD curl -f http://localhost:8081/health || exit 1
```

---

## 🟡 Issue #8: No Readiness Check Delay

### Problem
Server starts immediately but:
- Database may not be ready
- VIP data files not loaded yet
- Services not initialized

Users see "boxes with no functions" because API returns errors during startup.

### Solution

#### Add startup delay in main.go (after line 234):
```go
log.Printf("🚀 Health Management System server started on port %s", cfg.Server.Port)

// ADD THIS: Wait for services to be ready
log.Println("⏳ Warming up services...")
time.Sleep(2 * time.Second)
log.Println("✅ All services ready")
```

---

## 📋 Complete Fix Checklist

### Immediate Actions (Critical - Do First)

- [ ] **1. Configure SSL/HTTPS in Coolify**
  - Add domain in Coolify
  - Enable automatic SSL
  - Wait for certificate

- [ ] **2. Fix CORS Configuration**
  - Update `main.go` with proper CORS config
  - Add your actual domain
  - Enable credentials

- [ ] **3. Fix Port Mismatch**
  - Set `PORT=8081` in Coolify environment variables
  - Update Dockerfile EXPOSE and HEALTHCHECK

- [ ] **4. Set Required Environment Variables**
  ```bash
  JWT_SECRET=your-secret-key-min-32-chars
  ENV=production
  PORT=8081
  LOG_LEVEL=warn
  ```

### Important Actions (Do Next)

- [ ] **5. Configure Persistent Storage**
  - Add volume in Coolify for `/app/data`
  
- [ ] **6. Verify Frontend Files**
  - Check `frontend/` directory exists
  - Update Dockerfile to verify

- [ ] **7. Update Health Check**
  - Fix port in Dockerfile
  - Increase timeout and retries

### Nice to Have

- [ ] **8. Add Monitoring**
  - Enable Coolify logs
  - Set up alerts

- [ ] **9. Test Deployment**
  - Run `./deploy.sh`
  - Check health endpoint
  - Test API calls

---

## 🔧 Quick Fix Script

Save as `fix-deployment.sh`:

```bash
#!/bin/bash

echo "🔧 Applying deployment fixes..."

# 1. Update CORS in main.go
echo "📝 Updating CORS configuration..."
# (Manual step - see Issue #2)

# 2. Update Dockerfile
echo "📝 Updating Dockerfile..."
sed -i.bak 's/EXPOSE 8080/EXPOSE 8081/' Dockerfile
sed -i.bak 's/localhost:8080/localhost:8081/g' Dockerfile

# 3. Verify frontend exists
if [ ! -d "frontend" ]; then
    echo "❌ ERROR: frontend directory not found!"
    exit 1
fi
echo "✅ Frontend directory found"

# 4. Build and deploy
echo "🚀 Building and deploying..."
./deploy.sh

echo "✅ Fixes applied! Check Coolify dashboard."
echo ""
echo "🔔 REMEMBER TO:"
echo "   1. Set environment variables in Coolify"
echo "   2. Configure SSL/HTTPS"
echo "   3. Add persistent volume"
```

---

## 🧪 Testing After Fixes

### Test 1: Health Check
```bash
curl https://api.doctorhealthy1.com/health
# Expected: {"status": "healthy", ...}
```

### Test 2: HTTPS Certificate
```bash
curl -I https://api.doctorhealthy1.com/health
# Expected: No SSL errors
```

### Test 3: CORS
```bash
curl -H "Origin: https://doctorhealthy1.com" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     https://api.doctorhealthy1.com/api/v1/health
# Expected: Access-Control-Allow-Origin header present
```

### Test 4: Frontend Loading
```bash
curl https://api.doctorhealthy1.com/
# Expected: HTML content with <title>Doctor Healthy</title>
```

### Test 5: API Endpoint
```bash
curl https://api.doctorhealthy1.com/api/v1/recipes \
     -H "X-API-Key: your-key"
# Expected: JSON response with recipes
```

---

## 🎯 Root Causes Summary

| Issue | Symptom | Root Cause |
|-------|---------|------------|
| "Server not found" | Intermittent | Port mismatch (8080/8081), health checks failing |
| "Not secure" | Browser warning | No SSL certificate configured |
| "Boxes with no functions" | UI broken | CORS blocking API calls |
| "No meal/workout generation" | API fails | CORS + Frontend not loaded |
| "Clicks not working" | JavaScript errors | CORS + Mixed content blocking |

---

## 🚀 Next Steps

1. **Start with HTTPS/SSL** - This fixes the "not secure" warning
2. **Fix CORS** - This fixes API communication
3. **Fix ports** - This fixes "server not found"
4. **Set environment variables** - This ensures correct configuration
5. **Add persistent storage** - This prevents data loss
6. **Test thoroughly** - Follow testing checklist above

---

## 📞 Need Help?

If issues persist after these fixes:

1. Check Coolify logs: Application → Logs
2. Check application logs: `docker logs <container-id>`
3. Test locally: `docker build -t test . && docker run -p 8081:8081 test`
4. Review this checklist again

---

**Generated:** 2025-10-01  
**Status:** Ready for implementation
