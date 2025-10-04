# Coolify Configuration Guide - Step by Step

## 🎯 Goal
Configure your DoctorHealthy1 application with:
- ✅ SSL/HTTPS certificates
- ✅ Environment variables
- ✅ Persistent storage

## 📝 Your Configuration Details

**Coolify Server:**
- Host: `128.140.111.171`
- Port: `8000`
- Dashboard URL: http://128.140.111.171:8000

**API Access Token:**
```
4|jdTX2lUb2q6IOrwNGkHyQBCO74JJeeRHZVvFNwgI6b376a50
```
(Alternative: `jdTX2lUb2q6IOrwNGkHyQBCO74JJeeRHZVvFNwgI6b376a50`)

**Domains:**
- API Domain: `api.doctorhealthy1.com`
- Web App Domain: `my.doctorhealthy1.com`

**Application UUID:** `hcw0gc8wcwk440gw4c88408o`

---

## 🚀 Quick Start

### Option 1: Automated Script (Recommended)
```bash
cd /Users/khaledahmedmohamed/Desktop/pure\ nutrition\ cursor\ or\ kiro
chmod +x configure-coolify.sh
./configure-coolify.sh
```

### Option 2: Manual Configuration (Follow sections below)

---

## 📋 Part 1: Domain & SSL Configuration

### Step 1.1: Access Coolify Dashboard
1. Open browser
2. Navigate to: http://128.140.111.171:8000
3. Log in with your credentials

### Step 1.2: Find Your Application
1. Click on **"Applications"** in left sidebar
2. Find: **"DoctorHealthy1 API"** or search by UUID: `hcw0gc8wcwk440gw4c88408o`
3. Click on the application to open it

### Step 1.3: Configure API Domain
1. In application view, click **"Domains"** tab
2. Click **"Add Domain"** button
3. Fill in the form:
   ```
   Domain: api.doctorhealthy1.com
   ```
4. ✅ Toggle ON: **"Generate Automatic HTTPS"** (Let's Encrypt)
5. Click **"Save"**

### Step 1.4: Configure Web App Domain (if separate)
If your web app is a separate application:
1. Find your frontend/webapp application in Coolify
2. Repeat steps 1.2-1.3 with domain: `my.doctorhealthy1.com`

### Step 1.5: Wait for SSL Certificate
1. Status will show "Provisioning" → "Active" (2-5 minutes)
2. Green indicator means SSL is active
3. If it fails, check DNS configuration (see Step 1.6)

### Step 1.6: Verify DNS Configuration
**IMPORTANT:** Before SSL can work, DNS must point to Coolify server.

Check DNS resolution:
```bash
nslookup api.doctorhealthy1.com
# Should return: 128.140.111.171

nslookup my.doctorhealthy1.com
# Should return: 128.140.111.171
```

If DNS doesn't resolve:
1. Go to your domain registrar (e.g., GoDaddy, Namecheap, Cloudflare)
2. Add/Update A records:
   - Host: `api` → Value: `128.140.111.171`
   - Host: `my` → Value: `128.140.111.171`
3. Wait 5-15 minutes for DNS propagation
4. Retry SSL configuration in Coolify

### Step 1.7: Verify SSL is Working
```bash
# Test 1: Check certificate
curl -I https://api.doctorhealthy1.com/health

# Expected: HTTP/2 200 (no SSL errors)

# Test 2: Browser test
# Open: https://api.doctorhealthy1.com/health
# Should see: Green lock icon 🔒
```

✅ **Part 1 Complete!** Your domains now have SSL certificates.

---

## 🔧 Part 2: Environment Variables Configuration

### Step 2.1: Navigate to Environment Variables
1. In your application view (DoctorHealthy1 API)
2. Click **"Environment Variables"** tab
3. You'll see a list of existing variables (if any)

### Step 2.2: Add Required Variables

**Method A: Add One by One**
1. Click **"Add Variable"** button
2. For each variable below, fill in:
   - Key: (variable name)
   - Value: (variable value)
   - Click **"Add"**

**Method B: Bulk Add (Faster)**
1. Click **"Bulk Add Variables"** (if available)
2. Copy-paste all variables at once (see Step 2.3)
3. Click **"Save"**

### Step 2.3: Variables to Add

Copy this entire block for bulk add:

```bash
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-CHANGE-THIS-NOW
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
ALLOWED_ORIGIN=https://api.doctorhealthy1.com
```

**⚠️ IMPORTANT:** Change `JWT_SECRET` to a secure random string!

### Step 2.4: Generate Secure JWT_SECRET

**Option A: Use openssl (Mac/Linux)**
```bash
openssl rand -base64 48 | tr -d '\n'
```

**Option B: Use online generator**
Visit: https://randomkeygen.com/
Use "CodeIgniter Encryption Keys" (256-bit)

**Option C: Use Python**
```bash
python3 -c "import secrets; print(secrets.token_urlsafe(48))"
```

Replace the JWT_SECRET value with your generated secret.

### Step 2.5: Save and Verify
1. Click **"Save"** or **"Update"**
2. Application will restart automatically
3. Check logs to verify: "Server started on port 8081"

### Step 2.6: Verify Environment Variables
In Coolify dashboard:
1. Go to: Application → **"Shell"** or **"Terminal"**
2. Run: `env | grep JWT_SECRET`
3. Should display your JWT_SECRET (partially masked)

✅ **Part 2 Complete!** Environment variables are configured.

---

## 💾 Part 3: Persistent Storage Configuration

### Step 3.1: Navigate to Storage
1. In your application view (DoctorHealthy1 API)
2. Click **"Storages"** or **"Volumes"** tab
3. You'll see a list of mounted volumes (if any)

### Step 3.2: Add Persistent Volume
1. Click **"Add Persistent Storage"** or **"Add Volume"**
2. Fill in the form:

```
Name: app-data
Mount Path: /app/data
Volume Type: Local Volume (default)
Host Path: (leave empty - Coolify auto-manages)
Size: 2GB (or adjust as needed)
```

3. Click **"Save"** or **"Create"**

### Step 3.3: Container Restart
- Coolify will automatically restart the container
- New volume will be mounted at `/app/data`
- Your SQLite database will persist across deployments

### Step 3.4: Verify Storage
**Method A: Via Coolify Shell**
1. Go to: Application → **"Shell"** or **"Terminal"**
2. Run:
```bash
ls -la /app/data
# Should show app.db file after first run

df -h /app/data
# Should show mounted volume
```

**Method B: Via Logs**
1. Check application logs
2. Look for: "Database initialized at /app/data/app.db"

### Step 3.5: Test Persistence
1. Create some data (API key, etc.)
2. Trigger a redeploy: Application → **"Deploy"** → **"Force Deploy"**
3. After deployment, verify data still exists
4. If data persists → ✅ Storage configured correctly!

✅ **Part 3 Complete!** Persistent storage is configured.

---

## 🚀 Part 4: Deploy & Verify

### Step 4.1: Trigger Deployment
1. Go to your application in Coolify
2. Click **"Deploy"** button (top right)
3. Or: Click **"Force Redeploy"** to rebuild from scratch

### Step 4.2: Monitor Deployment
1. Click **"Logs"** tab
2. Watch for:
   - ✅ "Building image..."
   - ✅ "Frontend files found"
   - ✅ "Database connection pool configured"
   - ✅ "Server started on port 8081"
   - ✅ "Health Management System server started"

### Step 4.3: Check Health Status
Wait 2-3 minutes, then test:

```bash
# Test 1: Health endpoint
curl https://api.doctorhealthy1.com/health

# Expected output:
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-10-01T..."
}

# Test 2: Ready endpoint
curl https://api.doctorhealthy1.com/ready

# Test 3: Frontend
curl -I https://api.doctorhealthy1.com/

# Expected: HTTP/2 200 with HTML content
```

### Step 4.4: Browser Tests
1. Open: https://api.doctorhealthy1.com/
2. Should see: Full website with navigation
3. Check browser console (F12):
   - ❌ No CORS errors
   - ❌ No SSL warnings
   - ✅ API calls succeed

### Step 4.5: Test Features
1. Try generating a meal plan
2. Try generating a workout plan
3. Click various buttons
4. All features should work!

✅ **Part 4 Complete!** Application is deployed and working.

---

## 🔍 Troubleshooting

### Issue: SSL Certificate Failed to Provision

**Symptoms:** Domain shows "SSL Error" or "Not Secure"

**Solutions:**
1. **Check DNS:**
   ```bash
   nslookup api.doctorhealthy1.com
   # Must return: 128.140.111.171
   ```
   If not, update DNS A record and wait 15 minutes

2. **Check Port 80/443:**
   Ensure ports 80 and 443 are open on firewall

3. **Retry SSL:**
   - Coolify → Domains → Click "Retry SSL"
   - Or delete domain and re-add

4. **Check Let's Encrypt Rate Limits:**
   - Max 5 failures per hour
   - Wait 1 hour if rate limited

---

### Issue: Environment Variables Not Applied

**Symptoms:** Server starts but uses wrong port or missing JWT_SECRET

**Solutions:**
1. **Verify Variables Exist:**
   - Coolify → Environment Variables
   - Check all 11+ variables are listed

2. **Check Variable Format:**
   - No quotes needed: `JWT_SECRET=value` ✅
   - Not: `JWT_SECRET="value"` ❌

3. **Restart Application:**
   - Coolify → Deploy → Force Redeploy

4. **Check Logs:**
   ```
   Application → Logs
   Look for: "Server Port: 8081"
   ```

---

### Issue: Persistent Storage Not Working

**Symptoms:** Data lost after deployment

**Solutions:**
1. **Verify Mount Path:**
   - Must be exactly: `/app/data`
   - Not: `/data` or `app/data`

2. **Check Volume Status:**
   - Coolify → Storages → Should show "Mounted"

3. **Permissions:**
   Container runs as `appuser` (UID 1000)
   ```bash
   # In container shell:
   ls -la /app/data
   # Owner should be: appuser:appuser
   ```

4. **Recreate Volume:**
   - Delete existing volume
   - Add new volume
   - Redeploy

---

### Issue: CORS Errors in Browser

**Symptoms:** "CORS policy blocked" in browser console

**Solutions:**
1. **Verify CORS Config:**
   Our fix in main.go should have resolved this

2. **Check Allowed Origins:**
   ```bash
   # In environment variables, verify:
   ALLOWED_ORIGIN=https://api.doctorhealthy1.com
   ```

3. **Test from Command Line:**
   ```bash
   curl -H "Origin: https://my.doctorhealthy1.com" \
        -H "Access-Control-Request-Method: POST" \
        -X OPTIONS \
        https://api.doctorhealthy1.com/api/v1/recipes
   
   # Should return: Access-Control-Allow-Origin header
   ```

---

### Issue: Port Conflict

**Symptoms:** "Port already in use" in logs

**Solutions:**
1. **Verify PORT=8081:**
   Check environment variable

2. **Check Dockerfile:**
   Should have: `EXPOSE 8081`

3. **Coolify Port Mapping:**
   - Internal: 8081
   - External: Managed by Coolify (usually 80/443)

---

## 📊 Configuration Checklist

Use this to verify everything is configured:

### Domain & SSL
- [ ] API domain added: `api.doctorhealthy1.com`
- [ ] Web app domain added: `my.doctorhealthy1.com`
- [ ] DNS resolves correctly (nslookup)
- [ ] SSL certificate provisioned (green indicator)
- [ ] HTTPS works in browser (green lock)

### Environment Variables
- [ ] JWT_SECRET (32+ characters)
- [ ] ENV=production
- [ ] PORT=8081
- [ ] LOG_LEVEL=warn
- [ ] DB_PATH=/app/data/app.db
- [ ] All 11+ variables configured
- [ ] Variables visible in Coolify dashboard

### Persistent Storage
- [ ] Volume added: app-data
- [ ] Mount path: /app/data
- [ ] Status: Mounted
- [ ] Database persists across deployments
- [ ] Volume shows in Coolify dashboard

### Deployment
- [ ] Application deployed successfully
- [ ] No errors in logs
- [ ] Health check returns 200
- [ ] Frontend loads correctly
- [ ] API endpoints work
- [ ] No CORS errors

### Verification
- [ ] curl https://api.doctorhealthy1.com/health works
- [ ] Browser shows green lock icon
- [ ] Meal generation works
- [ ] Workout generation works
- [ ] All UI elements functional

---

## 🎯 Quick Command Reference

```bash
# Test DNS
nslookup api.doctorhealthy1.com
nslookup my.doctorhealthy1.com

# Test SSL
curl -I https://api.doctorhealthy1.com/health

# Test Health
curl https://api.doctorhealthy1.com/health | jq

# Test Frontend
curl -I https://api.doctorhealthy1.com/

# Generate JWT Secret
openssl rand -base64 48 | tr -d '\n'

# Check Coolify API
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://128.140.111.171:8000/api/v1/applications/hcw0gc8wcwk440gw4c88408o
```

---

## 📞 Support

If you encounter issues:

1. **Check Coolify Logs:**
   Dashboard → Application → Logs

2. **Review Documentation:**
   - DEPLOYMENT-CHECKLIST.md
   - DEPLOYMENT-ISSUES-ANALYSIS.md
   - DEPLOYMENT-FIX-SUMMARY.md

3. **Verify Configuration:**
   Use checklist above

4. **Common Issues:**
   See Troubleshooting section

---

**Last Updated:** 2025-10-01  
**Status:** Ready to Configure  
**Estimated Time:** 30-45 minutes
