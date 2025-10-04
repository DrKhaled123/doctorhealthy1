# 🎯 COOLIFY MANUAL CONFIGURATION STEPS

## ✅ What's Done
- ✅ API Connection: Successfully connected to Coolify
- ✅ DNS Resolution: my.doctorhealthy1.com → 128.140.111.171
- ✅ JWT Secret: Generated secure 64-character secret
- ✅ Environment Variables: Saved to `coolify-env-vars.txt`

## 📋 What You Need to Do (3 Steps - 10 Minutes)

### 🔗 Dashboard Access
Open: http://128.140.111.171:8000

---

## STEP 1: Environment Variables (5 minutes)

**Navigation:**
1. Click on **Applications** (left sidebar)
2. Find and click **doctorhealthy1-api**
3. Click **Environment Variables** tab

**Quick Add Method:**
1. Click **"Bulk Add"** or **"Bulk Edit"** button
2. Copy ALL the content from `coolify-env-vars.txt` (shown below)
3. Paste into the text area
4. Click **Save**

**Environment Variables to Add:**
```bash
JWT_SECRET=sE6oOg8fJkAuXZlC21TDQphFVK613/hh5HQQhMsHKMmR5QilVXk2/jDGnPaY+9II
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

✅ **Verify:** You should see 12 environment variables listed

---

## STEP 2: Domain & SSL (3 minutes)

**Navigation:**
1. Stay in **doctorhealthy1-api** application
2. Click **Domains** tab

**Configure Domain (Coolify v4):**

**Option A - If you see a domain input field:**
1. Look for the **"FQDN"** or **"Domain"** input field
2. Clear any existing domain (like "localhost")
3. Enter: `my.doctorhealthy1.com`
4. Click **Save**
5. SSL will be automatically provisioned by Coolify (using Let's Encrypt)

**Option B - If you see "Add Domain" button:**
1. Click **"+ Add Domain"** or **"Add"** button
2. Enter domain: `my.doctorhealthy1.com`
3. Click **Save** or **Add**
4. SSL will be automatically enabled

**Option C - Configuration tab:**
If Domains tab is not visible:
1. Go to **Configuration** or **General** tab
2. Look for **"Domain"** or **"URLs"** section
3. Enter: `my.doctorhealthy1.com`
4. Click **Save**

**Important Notes:**
- In Coolify v4, SSL/HTTPS is usually **automatic** when you add a domain
- You don't need to toggle anything - Coolify handles Let's Encrypt automatically
- The domain should start showing as `https://` once provisioned
- If you see a padlock icon 🔒 next to the domain, SSL is active

**Wait for SSL (2-5 minutes):**
- The certificate provisioning happens in the background
- Monitor the logs or domain status
- Once complete, the domain will be accessible via HTTPS

✅ **Verify:** Domain shows `https://my.doctorhealthy1.com` and is accessible

---

## STEP 3: Persistent Storage (2 minutes)

**Navigation:**
1. Stay in **doctorhealthy1-api** application
2. Click **Storages** or **Volumes** tab

**Add Storage:**
1. Click **"Add Persistent Storage"** or **"+ Add Volume"**
2. Fill in:
   - **Name:** `app-data`
   - **Mount Path:** `/app/data`
   - **Host Path:** (leave empty or auto-generated)
   - **Size:** `2GB` (or more if needed)
3. Click **Save** or **Create**

**Container Restart:**
- Coolify will automatically restart the container to mount the volume
- Wait 30-60 seconds for restart

✅ **Verify:** Storage shows "Mounted" status

---

## 🚀 STEP 4: Deploy (1 minute)

**After completing steps 1-3:**

1. Go back to **doctorhealthy1-api** main page
2. Look for the **"Deploy"** or **"Redeploy"** button (usually top-right)
3. Click **Deploy**
4. Monitor the deployment logs (should take 3-5 minutes)

**Watch for these log messages:**
- ✅ "Building image..."
- ✅ "Starting container..."
- ✅ "Server started on port 8081"
- ✅ "Health check passed"

---

## 🧪 VERIFICATION TESTS

### After deployment completes, run these tests:

### Test 1: SSL Certificate
```bash
curl -I https://my.doctorhealthy1.com/health
```
**Expected:** `HTTP/2 200` (no SSL errors)

### Test 2: Health Check
```bash
curl https://my.doctorhealthy1.com/health
```
**Expected:** `{"status":"healthy",...}`

### Test 3: Frontend
Open in browser: https://my.doctorhealthy1.com/

**Expected:**
- ✅ Green padlock 🔒 in address bar
- ✅ Website loads completely
- ✅ No "Not Secure" warnings

### Test 4: Browser Console
1. Open website
2. Press F12 (Developer Tools)
3. Click **Console** tab

**Expected:**
- ✅ No CORS errors
- ✅ No "Failed to fetch" errors
- ✅ API calls succeed

---

## ⏱️ Time Estimate

| Step | Time |
|------|------|
| Environment Variables | 5 min |
| Domain & SSL | 3 min |
| Persistent Storage | 2 min |
| Deploy | 5 min |
| **Total** | **15 minutes** |

---

## 🐛 Troubleshooting

### ❌ Can't Find SSL Toggle in Coolify v4
**What's Happening:** Coolify v4 enables SSL automatically

**Solution:**
1. Just add your domain: `my.doctorhealthy1.com`
2. Click Save
3. Coolify will automatically request and install Let's Encrypt SSL certificate
4. Wait 2-5 minutes for certificate provisioning
5. Check if domain becomes accessible via `https://`

**Alternative - Check Settings:**
1. Go to your application
2. Look for **Settings** or **Configuration** tab
3. Check if there's a **"Force HTTPS"** or **"Redirect to HTTPS"** option
4. Enable it if available

### ❌ SSL Certificate Fails
**Problem:** "Failed to obtain certificate" or "DNS verification failed"

**Solution:**
1. Verify DNS: `nslookup my.doctorhealthy1.com`
2. Should return: `128.140.111.171`
3. Wait 15 minutes for DNS propagation
4. In Coolify, try redeploying the application
5. Check application logs for SSL errors
6. Ensure port 80 and 443 are accessible (Coolify needs these for Let's Encrypt verification)

### ❌ Container Won't Start
**Problem:** Deployment fails or container stops immediately

**Solution:**
1. Check logs in Coolify dashboard
2. Look for error about JWT_SECRET
3. Verify all environment variables are set correctly
4. Ensure PORT=8081 is set

### ❌ "Server Not Found"
**Problem:** Browser can't reach the site

**Solution:**
1. Check deployment status (should be "Running")
2. Check logs for "Server started on port 8081"
3. Verify domain in Coolify matches `my.doctorhealthy1.com`
4. Test direct IP: http://128.140.111.171:8000

### ❌ Data Lost After Restart
**Problem:** Database is empty after redeployment

**Solution:**
1. Check Storages tab shows "Mounted"
2. Verify mount path is `/app/data`
3. In container shell, run: `ls -la /app/data`
4. Should show `app.db` file

---

## 📸 Quick Visual Guide

### Environment Variables Screen
Look for:
- [ ] "Bulk Add" or "Bulk Edit" button
- [ ] Text area to paste variables
- [ ] "Save" button at bottom

### Domains Screen
Look for:
- [ ] "FQDN" or "Domain" input field
- [ ] "Generate Let's Encrypt SSL" toggle switch
- [ ] Green indicator when SSL is active

### Storages Screen
Look for:
- [ ] "+ Add Persistent Storage" button
- [ ] Fields for Name, Mount Path
- [ ] "Mounted" status indicator

---

## ✅ Completion Checklist

After completing all steps, verify:

- [ ] 12 environment variables are set (including JWT_SECRET)
- [ ] Domain shows: my.doctorhealthy1.com
- [ ] SSL certificate is active (green padlock)
- [ ] Persistent storage shows "Mounted" at /app/data
- [ ] Container is running (deployment successful)
- [ ] Health check returns HTTP 200
- [ ] Website loads with HTTPS
- [ ] No CORS errors in browser console
- [ ] Meal/workout generation works
- [ ] Data persists after restart

---

## 📞 Need Help?

**Coolify Dashboard:** http://128.140.111.171:8000

**Application Logs:** 
http://128.140.111.171:8000/application/hcw0gc8wcwk440gw4c88408o/logs

**Environment Variables File:** `coolify-env-vars.txt`

**Full Documentation:** `COOLIFY-SETUP-GUIDE.md`

---

**🎯 You're Almost There!**

The hard work is done - configuration is prepared, DNS is working, and the secure JWT secret is generated. Just follow these 3 manual steps in the Coolify dashboard and you'll be live in 15 minutes! 🚀
