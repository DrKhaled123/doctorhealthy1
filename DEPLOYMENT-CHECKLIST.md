# 🚀 Deployment Checklist - Quick Reference

## ⚡ Critical Actions (Do These First!)

### 🔴 Step 1: Configure SSL/HTTPS in Coolify
**Why:** Fixes "webpage not secure" warning

**How:**
```
1. Open Coolify Dashboard: http://128.140.111.171:8000
2. Find your app: hcw0gc8wcwk440gw4c88408o
3. Go to: Settings → Domains
4. Add domain: api.doctorhealthy1.com
5. Toggle: "Enable Automatic SSL" ✅
6. Wait: 2-3 minutes for certificate
7. Verify: Green lock icon appears
```

**Status:** [ ] Not Started / [ ] In Progress / [ ] ✅ Complete

---

### 🔴 Step 2: Set Environment Variables
**Why:** Fixes server configuration and startup

**How:**
```
In Coolify → Application → Environment Variables

Add these (copy-paste):

JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-change-this-now
ENV=production
PORT=8081
LOG_LEVEL=warn
DB_PATH=/app/data/app.db
RATE_LIMIT=100
API_KEY_PREFIX=dh_
API_KEY_LENGTH=32
```

**Status:** [ ] Not Started / [ ] In Progress / [ ] ✅ Complete

---

### 🟠 Step 3: Add Persistent Storage
**Why:** Prevents data loss on container restart

**How:**
```
In Coolify → Application → Storage

1. Click: "Add Volume"
2. Name: app-data
3. Mount Path: /app/data
4. Size: 2GB
5. Click: "Save"
```

**Status:** [ ] Not Started / [ ] In Progress / [ ] ✅ Complete

---

### 🟢 Step 4: Deploy Application
**Why:** Apply all fixes to production

**How:**
```bash
# On your local machine:
cd /Users/khaledahmedmohamed/Desktop/pure\ nutrition\ cursor\ or\ kiro

# Option A: Automated (Recommended)
./fix-deployment.sh

# Option B: Manual
./deploy.sh
```

**Status:** [ ] Not Started / [ ] In Progress / [ ] ✅ Complete

---

## ✅ Verification Steps (After Deployment)

### Test 1: SSL Certificate Working
```bash
# Open in browser (should show green lock):
https://api.doctorhealthy1.com/health

# Or test with curl:
curl -I https://api.doctorhealthy1.com/health
# Look for: HTTP/2 200
```
**Result:** [ ] ❌ Failed / [ ] ✅ Passed

---

### Test 2: Health Endpoint Responding
```bash
curl https://api.doctorhealthy1.com/health
# Expected output:
# {"status":"healthy","database":"connected",...}
```
**Result:** [ ] ❌ Failed / [ ] ✅ Passed

---

### Test 3: Frontend Loading
```bash
# Open in browser:
https://api.doctorhealthy1.com/

# Should see: Full website with navigation
# Not: "Not Found" or blank page
```
**Result:** [ ] ❌ Failed / [ ] ✅ Passed

---

### Test 4: API Endpoints Working
```bash
# Test recipes endpoint:
curl https://api.doctorhealthy1.com/api/v1/recipes \
     -H "X-API-Key: test-key"

# Should return: JSON data (even if "unauthorized" - that's OK)
# Should NOT return: CORS error or "Not Found"
```
**Result:** [ ] ❌ Failed / [ ] ✅ Passed

---

### Test 5: CORS Configuration
```bash
# Open browser console (F12) on your frontend
# Go to Network tab
# Try generating a meal/workout plan
# Check: No "CORS policy" errors in console
```
**Result:** [ ] ❌ Failed / [ ] ✅ Passed

---

## 🐛 Troubleshooting Guide

### Issue: "Still seeing 'Not Secure' warning"
**Likely Cause:** SSL certificate not yet provisioned  
**Solution:** Wait 5 more minutes, refresh browser  
**Check:** Coolify logs for certificate status

---

### Issue: "Server not found"
**Likely Cause:** PORT environment variable not set  
**Solution:**
1. Go to Coolify → Environment Variables
2. Verify: PORT=8081
3. Redeploy application

---

### Issue: "Boxes with no functions"
**Likely Cause:** CORS blocking API calls or frontend not loaded  
**Solution:**
1. Open browser console (F12)
2. Look for CORS errors
3. Verify: Frontend files loaded (view page source)
4. Check: main.go has correct CORS config

---

### Issue: "Clicks not working"
**Likely Cause:** JavaScript errors from CORS/SSL issues  
**Solution:**
1. Open browser console (F12)
2. Look for JavaScript errors
3. Check: Mixed content warnings
4. Verify: SSL certificate installed

---

### Issue: "No meal/workout generation"
**Likely Cause:** API calls being blocked  
**Solution:**
1. Test API directly: `curl https://api.../api/v1/health`
2. Check browser Network tab
3. Verify: CORS errors resolved
4. Check: API key is valid

---

## 📊 Deployment Status Overview

### Before Fixes:
- [ ] ❌ SSL/HTTPS configured
- [ ] ❌ CORS properly configured
- [ ] ❌ Correct port (8081)
- [ ] ❌ Environment variables set
- [ ] ❌ Persistent storage
- [ ] ❌ Frontend verification

### After Fixes (Code):
- [x] ✅ CORS fixed in main.go
- [x] ✅ Port fixed in Dockerfile
- [x] ✅ Frontend check added
- [x] ✅ Health check improved
- [x] ✅ .env.example updated
- [x] ✅ Documentation created

### After Deployment (You Do):
- [ ] ⏳ SSL/HTTPS configured
- [ ] ⏳ Environment variables set
- [ ] ⏳ Persistent storage added
- [ ] ⏳ Application deployed
- [ ] ⏳ Verification tests passed

---

## 📞 Quick Help

### Where to Find Information:
- **Detailed Analysis:** `DEPLOYMENT-ISSUES-ANALYSIS.md` (67 pages)
- **Quick Summary:** `DEPLOYMENT-FIX-SUMMARY.md` (this is the action list)
- **Environment Setup:** `.env.example`
- **Fix Script:** `./fix-deployment.sh`
- **Deploy Script:** `./deploy.sh`

### Coolify Dashboard:
- **URL:** http://128.140.111.171:8000
- **App ID:** hcw0gc8wcwk440gw4c88408o
- **SSH:** Already configured in deploy.sh

---

## ✨ Success Criteria

Your deployment is successful when ALL of these are true:

- ✅ Browser shows green lock icon (HTTPS)
- ✅ No "not secure" warnings
- ✅ All UI elements visible ("boxes" have content)
- ✅ Meal generation works
- ✅ Workout generation works
- ✅ All buttons/clicks work
- ✅ No CORS errors in browser console
- ✅ API calls succeed
- ✅ Server stays up (doesn't restart)

---

## 🎯 Estimated Time

- **Setup in Coolify:** 10-15 minutes
- **Deployment:** 5-10 minutes
- **Verification:** 5 minutes
- **Total:** ~30 minutes

---

**Last Updated:** 2025-10-01  
**Status:** Ready to Deploy  
**Next Action:** Start with Step 1 (SSL Configuration)

---

**Quick Start:**
```bash
# 1. Do Steps 1-3 in Coolify Dashboard
# 2. Run this command:
./fix-deployment.sh
# 3. Check tests above
# 4. Done! ✨
```
