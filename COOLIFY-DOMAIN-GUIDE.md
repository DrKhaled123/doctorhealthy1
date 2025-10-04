# 🔍 Finding Domain & SSL Settings in Coolify v4

## Where to Find Domain Configuration

Coolify v4 has made SSL automatic! Here's where to find the domain settings:

---

## 📍 Location 1: Domains Tab (Most Common)

**Navigation:**
```
Dashboard → Applications → doctorhealthy1-api → Domains
```

**What You'll See:**
- A list of current domains (might show "localhost" or be empty)
- An input field or **"+ Add Domain"** button
- Existing domains with status indicators

**What to Do:**
1. If you see an input field:
   - Type: `my.doctorhealthy1.com`
   - Press Enter or click Save
   
2. If you see "+ Add Domain" button:
   - Click it
   - Enter: `my.doctorhealthy1.com`
   - Click Add/Save

**Result:** SSL will be provisioned automatically! ✅

---

## 📍 Location 2: Configuration Tab

**Navigation:**
```
Dashboard → Applications → doctorhealthy1-api → Configuration
```

**Look for:**
- **"FQDN"** field (Fully Qualified Domain Name)
- **"Domain"** or **"URL"** field
- **"Exposed Ports"** section

**What to Do:**
1. Find the domain/FQDN field
2. Enter: `my.doctorhealthy1.com`
3. Scroll down and click **Save**

---

## 📍 Location 3: General Settings

**Navigation:**
```
Dashboard → Applications → doctorhealthy1-api → General
```

**Look for:**
- **"Domain"** section
- **"Network"** section
- **"Expose"** settings

---

## 🔐 About SSL in Coolify v4

### ⚡ SSL is Automatic!

In Coolify v4, you **don't need to toggle SSL on/off**. Here's what happens:

1. **You add a domain** (e.g., `my.doctorhealthy1.com`)
2. **Coolify detects** it's a public domain
3. **Automatically requests** Let's Encrypt certificate
4. **Installs certificate** within 2-5 minutes
5. **Domain becomes accessible** via HTTPS

### 🔍 How to Check SSL Status

**Method 1: Check the Domain List**
- Look for a 🔒 (padlock) icon next to your domain
- Green indicator means SSL is active
- If domain shows `https://` prefix, SSL is working

**Method 2: Check Application Logs**
```
Dashboard → Applications → doctorhealthy1-api → Logs
```
Look for messages like:
- "Certificate obtained successfully"
- "SSL certificate installed"
- "HTTPS enabled"

**Method 3: Test in Browser**
Open: https://my.doctorhealthy1.com
- Green padlock 🔒 in address bar = SSL working
- "Not Secure" warning = SSL not yet provisioned

**Method 4: Use curl**
```bash
curl -I https://my.doctorhealthy1.com/health
```
If SSL is working, you'll see:
```
HTTP/2 200
```

If SSL is not ready:
```
curl: (60) SSL certificate problem
```

---

## 🎯 Step-by-Step Visual Guide

### Step 1: Find Your Application
```
1. Click "Applications" in left sidebar
2. Find "doctorhealthy1-api" in the list
3. Click on it
```

### Step 2: Look for Domain Settings

**Check these tabs (in order):**

**Tab 1: "Domains"** ⭐ (Most likely here)
```
┌─────────────────────────────────┐
│ Domains                         │
├─────────────────────────────────┤
│                                 │
│ [+ Add Domain]                  │
│                                 │
│ OR                              │
│                                 │
│ Domain: [________________]      │
│         [Save]                  │
└─────────────────────────────────┘
```

**Tab 2: "Configuration"**
```
┌─────────────────────────────────┐
│ Configuration                   │
├─────────────────────────────────┤
│ FQDN: [________________]        │
│                                 │
│ Port: 8081                      │
│                                 │
│ [Save Configuration]            │
└─────────────────────────────────┘
```

**Tab 3: "General" or "Settings"**
```
┌─────────────────────────────────┐
│ General Settings                │
├─────────────────────────────────┤
│ Application Name: ...           │
│ Domain: [________________]      │
│ Port: 8081                      │
│ [Save]                          │
└─────────────────────────────────┘
```

### Step 3: Add Your Domain

**Just type:**
```
my.doctorhealthy1.com
```

**Then click Save**

That's it! 🎉

---

## ⏱️ SSL Provisioning Timeline

```
0:00 - You add domain and click Save
0:30 - Coolify validates DNS (checks if domain points to server)
1:00 - Coolify requests certificate from Let's Encrypt
2:00 - Let's Encrypt verifies domain ownership (HTTP challenge)
3:00 - Certificate issued and installed
5:00 - HTTPS fully active ✅
```

**Total time: 2-5 minutes** (usually 3 minutes)

---

## 🚨 What If There's No Domain Field?

### Option A: Port Mapping
If you can't find a domain field, Coolify might be using port mapping:

1. Go to **Configuration** or **Network** tab
2. Look for **"Port Mapping"** or **"Ports"**
3. You should see something like: `80:8081` or `443:8081`
4. Add a **"Proxy"** or **"Load Balancer"** configuration
5. Set domain there

### Option B: Use Coolify Proxy
1. Go to **Proxy** or **Traefik** settings (in main dashboard)
2. Add a new route:
   - Host: `my.doctorhealthy1.com`
   - Target: Your application
   - Port: 8081
3. Save

### Option C: Check Documentation
```
In Coolify Dashboard:
Top-right → Help or Documentation → Domain Setup
```

---

## ✅ What Success Looks Like

After adding the domain, you should see:

**In Coolify:**
```
✅ my.doctorhealthy1.com (Active)
🔒 HTTPS Enabled
📊 Status: Running
```

**In Browser:**
```
Address bar: 🔒 https://my.doctorhealthy1.com
(Green padlock, "Connection is secure")
```

**In Terminal:**
```bash
$ curl -I https://my.doctorhealthy1.com/health
HTTP/2 200 
```

---

## 📞 Still Can't Find It?

### Screenshot the Interface

If you're still stuck:
1. Take a screenshot of your application page in Coolify
2. Look at all the available tabs
3. Check each tab for domain-related fields

### Common Tab Names to Check:
- ✓ Domains
- ✓ Configuration  
- ✓ General
- ✓ Settings
- ✓ Network
- ✓ Routing
- ✓ Proxy
- ✓ Advanced

### Default Behavior

Even without explicit SSL configuration, Coolify should:
- Detect public domains automatically
- Enable HTTPS by default
- Use Let's Encrypt for certificates
- Redirect HTTP → HTTPS automatically

**So just add the domain and wait!** 🚀

---

## 🎯 Quick Checklist

- [ ] Found domain configuration (Domains/Configuration/General tab)
- [ ] Added domain: `my.doctorhealthy1.com`
- [ ] Clicked Save
- [ ] Waited 3-5 minutes
- [ ] Checked domain shows HTTPS/🔒
- [ ] Tested in browser: https://my.doctorhealthy1.com
- [ ] Verified: Green padlock appears
- [ ] Confirmed: No "Not Secure" warning

---

**Remember:** In Coolify v4, SSL is automatic! You don't need to enable it manually. Just add your domain and Coolify handles the rest! 🎉
