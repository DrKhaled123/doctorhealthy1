# Quick Reference Card - Coolify Configuration

## 🎯 Your Configuration at a Glance

### 🌐 Access Information
```
Coolify Dashboard: http://128.140.111.171:8000
Application UUID: hcw0gc8wcwk440gw4c88408o
API Token: 4|jdTX2lUb2q6IOrwNGkHyQBCO74JJeeRHZVvFNwgI6b376a50
```

### 🔗 Your Domains
```
Primary Domain: my.doctorhealthy1.com
Server IP:      128.140.111.171
```

---

## ⚡ Super Quick Setup

### Option 1: Automated (2 minutes)
```bash
cd /Users/khaledahmedmohamed/Desktop/pure\ nutrition\ cursor\ or\ kiro
./configure-coolify.sh
```

### Option 2: Manual (15 minutes)
Follow: `COOLIFY-SETUP-GUIDE.md`

---

## 🔐 Part 1: SSL/HTTPS (5 min)

**Where:** Coolify → Applications → DoctorHealthy1 → Domains

**Steps:**
1. Add domain: `my.doctorhealthy1.com`
2. Toggle ON: "Generate Automatic HTTPS" ✅
3. Save
4. Wait 2-5 minutes for Let's Encrypt

**Verify:**
```bash
curl -I https://my.doctorhealthy1.com/health
# Should return: HTTP/2 200
```

**If DNS doesn't resolve:**
- Go to your domain registrar
- Add A record: `my` → `128.140.111.171`
- Wait 15 minutes

---

## 🔧 Part 2: Environment Variables (5 min)

**Where:** Coolify → Applications → DoctorHealthy1 → Environment Variables

**Quick Add (copy-paste this):**
```bash
JWT_SECRET=CHANGE-THIS-TO-RANDOM-64-CHARS-USE-openssl-rand-base64-48
ENV=production
PORT=8081
HOST=0.0.0.0
LOG_LEVEL=warn
DB_PATH=/app/data/app.db
RATE_LIMIT=100
API_KEY_PREFIX=dh_
API_KEY_LENGTH=32
SECURITY_RATE_LIMIT_REQUESTS=100
SECURITY_RATE_LIMIT_WINDOW=1m
ALLOWED_ORIGIN=https://my.doctorhealthy1.com
```

**Generate Secure JWT_SECRET:**
```bash
openssl rand -base64 48 | tr -d '\n'
```
Then replace the JWT_SECRET value above.

---

## 💾 Part 3: Persistent Storage (2 min)

**Where:** Coolify → Applications → DoctorHealthy1 → Storages

**Steps:**
1. Click "Add Persistent Storage"
2. Fill in:
   - Name: `app-data`
   - Mount Path: `/app/data`
   - Size: `2GB`
3. Save

**Verify:**
```bash
# In Coolify → Shell:
ls -la /app/data
```

---

## 🚀 Part 4: Deploy (5 min)

**Where:** Coolify → Applications → DoctorHealthy1

**Steps:**
1. Click "Deploy" button (top right)
2. Wait 3-5 minutes
3. Monitor logs

**Or deploy from local:**
```bash
./deploy.sh
```

---

## ✅ Verification Tests

Run these after deployment:

```bash
# Test 1: SSL Certificate
curl -I https://my.doctorhealthy1.com/health
# Expected: HTTP/2 200

# Test 2: Health Check
curl https://my.doctorhealthy1.com/health
# Expected: {"status":"healthy",...}

# Test 3: Frontend
open https://my.doctorhealthy1.com/
# Expected: Website loads with green lock 🔒

# Test 4: DNS
nslookup my.doctorhealthy1.com
# Expected: 128.140.111.171
```

---

## 🐛 Quick Troubleshooting

### ❌ "Not Secure" Warning
```bash
# Check DNS first
nslookup my.doctorhealthy1.com

# If DNS OK, check SSL in Coolify
# Dashboard → Domains → Should show green indicator

# If failed, click "Retry SSL"
```

### ❌ "Server Not Found"
```bash
# Check PORT environment variable
# Should be: PORT=8081

# Check logs for:
# "Server started on port 8081" ✅
```

### ❌ "Boxes with No Functions"
```bash
# Check browser console (F12)
# Look for CORS errors

# Verify CORS in main.go was updated
# Should allow: my.doctorhealthy1.com
```

### ❌ Data Lost After Deployment
```bash
# Check persistent storage
# Coolify → Storages → Should show "Mounted"

# Verify mount path: /app/data
```

---

## 📊 Configuration Checklist

**Before Deployment:**
- [ ] DNS configured (A records)
- [ ] SSL enabled in Coolify
- [ ] All environment variables set
- [ ] JWT_SECRET is secure (32+ chars)
- [ ] Persistent storage added
- [ ] Code changes committed

**After Deployment:**
- [ ] HTTPS works (green lock)
- [ ] Health check responds
- [ ] Frontend loads
- [ ] API calls work
- [ ] No CORS errors
- [ ] Data persists

---

## 🎯 Expected Results

| Feature | Status |
|---------|--------|
| SSL/HTTPS | ✅ Green lock icon |
| Health Check | ✅ Returns 200 |
| Frontend | ✅ Full website loads |
| Meal Generation | ✅ Working |
| Workout Generation | ✅ Working |
| Data Persistence | ✅ Survives restart |

---

## 📞 Quick Help

**Documentation:**
- Full Guide: `COOLIFY-SETUP-GUIDE.md`
- Deployment: `DEPLOYMENT-FIX-SUMMARY.md`
- Issues: `DEPLOYMENT-ISSUES-ANALYSIS.md`
- Checklist: `DEPLOYMENT-CHECKLIST.md`

**Scripts:**
- Configure: `./configure-coolify.sh`
- Deploy: `./deploy.sh`
- Fix Check: `./fix-deployment.sh`

**Coolify Dashboard:**
http://128.140.111.171:8000

**Primary Domain:**
https://my.doctorhealthy1.com

---

## 🚀 One-Liner Setup

```bash
# Full automated setup:
cd /Users/khaledahmedmohamed/Desktop/pure\ nutrition\ cursor\ or\ kiro && \
./configure-coolify-api.sh && \
./deploy.sh
```

---

**Total Time: ~30 minutes**  
**Difficulty: Easy (mostly automated)**  
**Support: See full guides above**

---

**Print this page for quick reference! 📄**
