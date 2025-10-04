# Quick Start Guide - Doctor Healthy Application

**After VIP Database Migration - September 30, 2025**

---

## 🚀 Quick Start

### **1. Start the Application**
```bash
cd "/Users/khaledahmedmohamed/Desktop/pure nutrition cursor or kiro"
./main
```

### **2. Verify Health**
```bash
curl http://localhost:8085/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-09-30T...",
  "checks": {
    "database": "healthy",
    "filesystem": "healthy"
  }
}
```

### **3. Test API Endpoint**
```bash
# Get VIP data summary (requires rate limiting adjustment or API key)
curl http://localhost:8085/api/v1/vip/summary
```

---

## 📁 VIP Database Files

All VIP database files are in the **project root directory**:

```
pure nutrition cursor or kiro/
├── vip-complaints.js          # 1.3 MB - Health complaints & recommendations
├── vip-drugs-nutrition.js     # 321 KB - Medications & supplements
├── vip-workouts.js            # 1.2 MB - Workout plans (all levels)
├── vip-injuries.js            # 713 KB - Injury protocols
├── vip-workouts-techniques.js # 104 KB - Training techniques
└── vip-type-plans.js          # 1.7 MB - Diet & meal plans
```

**Total:** 5.3 MB of comprehensive health data

---

## 🔧 Build & Run

### **Rebuild After Changes**
```bash
# Update dependencies
go mod tidy

# Build
go build -o main .

# Run
./main
```

### **Background Run**
```bash
# Start in background
./main &

# Stop
pkill -f "./main"
```

### **Docker Build**
```bash
# Build image
docker build -t doctorhealthy:latest .

# Run container
docker run -d -p 8085:8080 \
  --name doctorhealthy \
  -v $(pwd)/data:/app/data \
  doctorhealthy:latest
```

---

## 🌐 Access Points

### **Local Development**
- **API Server:** http://localhost:8085
- **Health Check:** http://localhost:8085/health
- **Ready Check:** http://localhost:8085/ready

### **Production** (After Deployment)
- **Domain:** https://my.doctorhealthy1.com
- **API:** https://my.doctorhealthy1.com/api/v1
- **Health:** https://my.doctorhealthy1.com/health

---

## 📚 Key API Endpoints

### **Public Endpoints**
```bash
GET  /health                    # Health check
GET  /ready                     # Readiness probe
GET  /                          # Frontend home
```

### **API Key Management** (Requires Auth)
```bash
POST   /api/v1/apikeys          # Create API key
GET    /api/v1/apikeys          # List API keys
GET    /api/v1/apikeys/:id      # Get API key details
PUT    /api/v1/apikeys/:id      # Update API key
DELETE /api/v1/apikeys/:id      # Delete API key
```

### **VIP Data Access** (Requires API Key)
```bash
GET  /api/v1/vip/summary        # Get VIP data summary
GET  /api/v1/medications        # Get medications database
GET  /api/v1/workouts           # Get workout plans
GET  /api/v1/supplements        # Get supplements
GET  /api/v1/injuries           # Get injury protocols
```

### **Health System** (Requires API Key)
```bash
POST /api/v1/enhanced/diet-plan      # Generate custom diet plan
POST /api/v1/enhanced/workout-plan   # Generate workout plan
POST /api/v1/enhanced/lifestyle-plan # Generate lifestyle plan
```

---

## 🔑 Authentication

### **API Key Format**
```
dh_<64_hex_characters>
```

### **Create API Key**
```bash
curl -X POST http://localhost:8085/api/v1/apikeys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Key",
    "permissions": ["read", "write"],
    "expiry_days": 365
  }'
```

### **Use API Key**
```bash
curl http://localhost:8085/api/v1/vip/summary \
  -H "X-API-Key: dh_your_key_here"
```

---

## ⚙️ Configuration

### **Environment Variables**
Create `.env` file in project root:

```bash
# Server
PORT=8085
HOST=0.0.0.0

# Database
DB_PATH=./data/apikeys.db

# Security
JWT_SECRET=your-super-secret-jwt-key-change-this
API_KEY_PREFIX=dh_
API_KEY_LENGTH=32
API_KEY_EXPIRY=8760h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:8080,https://my.doctorhealthy1.com

# Rate Limiting
SECURITY_RATE_LIMIT_REQUESTS=100
SECURITY_RATE_LIMIT_WINDOW=1m
```

### **Adjust Rate Limiting for Testing**
In `.env`:
```bash
SECURITY_RATE_LIMIT_REQUESTS=1000  # Increase for testing
SECURITY_RATE_LIMIT_WINDOW=1m
```

---

## 🐛 Troubleshooting

### **Issue: "rate limit exceeded"**
**Solution:**
1. Increase rate limit in `.env`:
   ```bash
   SECURITY_RATE_LIMIT_REQUESTS=1000
   ```
2. Restart application

### **Issue: "database locked"**
**Solution:**
1. Stop all running instances: `pkill -f "./main"`
2. Check for stuck processes: `ps aux | grep main`
3. Restart application

### **Issue: "bind: address already in use"**
**Solution:**
1. Find process using port: `lsof -i :8085`
2. Kill process: `kill -9 <PID>`
3. Or change port in `.env`: `PORT=8086`

### **Issue: VIP data not loading**
**Solution:**
1. Verify VIP files exist:
   ```bash
   ls -lah vip-*.js
   ```
2. Check file permissions:
   ```bash
   chmod 644 vip-*.js
   ```
3. Check logs for errors

---

## 📊 Monitoring

### **Check Application Status**
```bash
# Health check
curl http://localhost:8085/health

# Process check
ps aux | grep main

# Port check
lsof -i :8085
```

### **View Logs**
```bash
# Run in foreground to see logs
./main

# Or redirect to file
./main > app.log 2>&1 &

# View logs
tail -f app.log
```

### **Database Check**
```bash
# Check database file
ls -lh ./data/apikeys.db

# SQLite command line
sqlite3 ./data/apikeys.db "SELECT COUNT(*) FROM api_keys;"
```

---

## 🧪 Testing

### **Health Check Test**
```bash
# Should return 200 OK
curl -w "\nStatus: %{http_code}\n" http://localhost:8085/health
```

### **API Key Creation Test**
```bash
# Create test API key
curl -X POST http://localhost:8085/api/v1/apikeys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test API Key",
    "permissions": ["read"],
    "expiry_days": 30
  }' | jq '.'
```

### **VIP Data Test**
```bash
# Test with API key
API_KEY="dh_your_key_here"
curl http://localhost:8085/api/v1/vip/summary \
  -H "X-API-Key: $API_KEY" | jq '.'
```

---

## 📦 Dependencies

### **Current Stable Versions**
```go
Go 1.22                                    // Language
github.com/labstack/echo/v4 v4.12.0       // HTTP router
github.com/mattn/go-sqlite3 v1.14.22      // SQLite driver
github.com/google/uuid v1.6.0             // UUID generation
github.com/golang-jwt/jwt v3.2.2          // JWT auth
github.com/go-playground/validator v10.22.0 // Validation
```

### **Update Dependencies**
```bash
go mod tidy
go mod verify
```

---

## 📝 Quick Commands

### **Most Common Commands**
```bash
# Start
./main

# Start in background
./main &

# Stop
pkill -f "./main"

# Rebuild
go build -o main .

# Health check
curl http://localhost:8085/health

# View VIP files
ls -lh vip-*.js

# Check logs
tail -f app.log
```

---

## 🎯 Success Criteria

### **Application is Working if:**
- ✅ Health endpoint returns `"status": "healthy"`
- ✅ No "file not found" errors in logs
- ✅ Server starts on port 8085
- ✅ Database connects successfully
- ✅ VIP files are present (6 files, ~5.3 MB total)

### **VIP Database is Working if:**
- ✅ All 6 VIP `.js` files present in root directory
- ✅ No warnings about missing COMPLETE files
- ✅ File sizes match: complaints(1.3MB), workouts(1.2MB), type-plans(1.7MB)
- ✅ API endpoints return VIP data

---

## 📞 Support

### **Documentation**
- **Migration Guide:** `VIP-DATABASE-MIGRATION-COMPLETE.md`
- **Status Report:** `CURRENT-STATUS-REPORT.md`
- **This Guide:** `QUICK-START-GUIDE.md`
- **Deployment:** `PRODUCTION_DEPLOYMENT_GUIDE.md`

### **Resources**
- **Repository:** github.com/DrKhaled123/doctorhealthy1
- **Branch:** main
- **Go Version:** 1.22 (stable)
- **Echo Docs:** https://echo.labstack.com

---

**Last Updated:** September 30, 2025  
**Version:** 1.0.0 (VIP Database Edition)  
**Status:** ✅ Production Ready

---

*For detailed information, see CURRENT-STATUS-REPORT.md*
