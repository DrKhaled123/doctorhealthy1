## DoctorHealthy API

Production-ready Go (Echo) API for health, workouts, recipes, and PDFs.

### Local dev
```
make run
```

### Docker
```
docker build -t doctorhealthy-api:latest .
docker run -p 8081:8081 --name doctorhealthy --rm \
  -e PORT=8081 -e DB_PATH=/data/app.db \
  -e JWT_SECRET=change_me \
  -e CORS_ORIGINS=http://localhost:3000 \
  -v $(pwd)/.data:/data doctorhealthy-api:latest
```

### Coolify deploy
1. New Application → From Git
2. Build: Dockerfile, internal port 8081, health GET /health
3. Storage: mount /data
4. Env:
   - PORT=8081
   - DB_PATH=/data/app.db
   - JWT_SECRET=change_me
   - CORS_ORIGINS=https://www.doctorhealthy1.com
   - EMAIL_FOR_TLS=you@example.com
5. Domain: api.doctorhealthy1.com

### Required endpoints
- /health
- Recipes, Enhanced, Ultimate handlers (Echo)

### PDF
Returns application/pdf for diet/workout/lifestyle/recipes exports.
# 🏥 Comprehensive Health Management System

A complete health and nutrition management platform built with Go and Echo framework, providing personalized nutrition plans, workout routines, health management, and recipe generation with API key security.

## ✨ Features

### 🔑 API Key Management
- **Cryptographically Secure**: API keys with customizable length and prefix
- **Usage Tracking**: Track API key usage, statistics, and rate limiting
- **Permission System**: Granular permissions for different access levels

### 👤 User Health Profiles
- **Complete Health Data**: Age, weight, height, activity level, goals, diseases
- **Calorie Calculations**: BMI, BMR, TDEE using Mifflin-St Jeor equation
- **Goal-Based Planning**: Weight loss, muscle building, maintenance, etc.

### 🍽️ Nutrition Management
- **Personalized Meal Plans**: Weekly meal plans based on user data
- **8 Diet Types**: Keto, Mediterranean, Low-carb, DASH, Vegan, etc.
- **Allergy & Preference Filtering**: Exclude disliked foods and allergens
- **Macro Calculations**: Protein, carbs, fats based on goals and plan type

### 💪 Workout Planning
- **Gym & Home Workouts**: Exercise routines for different environments
- **Injury-Aware**: Filter exercises based on user injuries
- **Goal-Specific**: Muscle building, strength, weight loss, endurance
- **Exercise Alternatives**: Multiple options for each exercise

### 🏥 Health Management
- **Disease-Specific Advice**: Treatment plans for diabetes, hypertension, etc.
- **Supplement Recommendations**: Dosage calculations based on weight
- **Medication Awareness**: Drug interaction and safety information
- **Lifestyle Modifications**: Evidence-based lifestyle recommendations

### 🍳 Recipe Generation
- **Multi-Cuisine Support**: 15+ cuisines (Mediterranean, Asian, Middle Eastern, etc.)
- **Dietary Restrictions**: Filter by allergies, preferences, and health conditions
- **Difficulty Levels**: Easy, medium, hard recipes with cooking tips
- **Nutritional Information**: Calories, prep time, serving sizes

### 🌍 Multilingual Support
- **English & Arabic**: Full support with medical disclaimers
- **Cultural Sensitivity**: Cuisine preferences and dietary considerations

### 🔐 Security & Performance
- **OWASP Compliance**: Following security best practices
- **High Performance**: Built with Echo v4 framework
- **SQLite Database**: Lightweight, fast, and reliable
- **Rate Limiting**: Built-in abuse prevention
- **Medical Disclaimers**: Proper medical safety warnings

## 🏗️ Architecture

```
├── main.go                     # Application entry point
├── internal/
│   ├── config/                # Configuration management
│   ├── database/              # Database setup and migrations
│   ├── handlers/              # HTTP handlers (controllers)
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Data models and DTOs
│   └── services/              # Business logic layer
├── Makefile                   # Build automation
└── .env.example              # Environment configuration template
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- SQLite3

### Installation

1. **Clone and setup**:
```bash
git clone <repository>
cd api-key-generator
make setup-dev
```

2. **Install dependencies**:
```bash
make deps
```

3. **Run the application**:
```bash
make run
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8081/api/v1
```

### Endpoints

#### 👤 User Management

**Create User Profile**
```http
POST /users
Content-Type: application/json

{
  "name": "Ahmed Mohamed",
  "email": "ahmed@example.com",
  "age": 30,
  "weight": 75,
  "height": 175,
  "gender": "male",
  "activity_level": "moderate",
  "metabolic_rate": "medium",
  "goal": "build_muscle",
  "food_dislikes": ["mushrooms"],
  "allergies": ["nuts"],
  "diseases": ["diabetes"],
  "medications": ["metformin"],
  "preferred_cuisine": "mediterranean",
  "language": "en"
}
```

**Get User Profile**
```http
GET /users/{id}
```

**Calculate Daily Calories**
```http
GET /users/{id}/calories
```

#### 🍽️ Nutrition Management

**Generate Nutrition Plan**
```http
POST /nutrition/generate-plan
Content-Type: application/json

{
  "user_id": "user-id-here",
  "plan_type": "mediterranean",
  "duration": 1
}
```

**Get Available Plan Types**
```http
GET /nutrition/plan-types
```

#### 💪 Workout Management

**Generate Workout Plan**
```http
POST /workouts/generate-plan
Content-Type: application/json

{
  "user_id": "user-id-here",
  "goal": "build_muscle",
  "workout_type": "gym",
  "injuries": [],
  "complaints": [],
  "duration": 1
}
```

**Get Workout Goals**
```http
GET /workouts/goals
```

**Get Available Injuries**
```http
GET /workouts/injuries
```

#### 🏥 Health Management

**Generate Health Plan**
```http
POST /health/generate-plan
Content-Type: application/json

{
  "user_id": "user-id-here",
  "diseases": ["diabetes", "hypertension"],
  "medications": ["metformin"],
  "complaints": ["fatigue"]
}
```

**Get Available Diseases**
```http
GET /health/diseases
```

**Get Health Complaints**
```http
GET /health/complaints
```

#### 🍳 Recipe Management

**Generate Recipe**
```http
POST /recipes/generate
Content-Type: application/json

{
  "user_id": "user-id-here",
  "cuisine": "mediterranean",
  "meal_type": "lunch",
  "difficulty": "easy",
  "max_calories": 500
}
```

**Get Available Cuisines**
```http
GET /recipes/cuisines
```

**Get User Recipes**
```http
GET /users/{id}/recipes?limit=10
```

#### 🔑 API Key Management

**Create API Key**
```http
POST /api-keys
Content-Type: application/json

{
  "name": "Health App Key",
  "permissions": ["nutrition:generate", "workout:generate"],
  "expiry_days": 365
}
```

**Validate API Key**
```http
POST /validate
X-API-Key: ak_your_api_key_here
```

#### 🔍 Utility Endpoints

**Get Available Goals**
```http
GET /goals
```

**Get Medical Disclaimer**
```http
GET /disclaimer?lang=en
```

**Health Check**
```http
GET /health
```

## 🔐 Security Features

### API Key Generation
- **Cryptographically Secure**: Uses `crypto/rand` for secure random generation
- **Customizable Length**: Default 32 bytes (64 hex characters)
- **Prefix Support**: Configurable prefix (default: `ak_`)
- **Expiration**: Configurable expiration dates

### Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`
- `Content-Security-Policy: default-src 'self'`

### Rate Limiting
- Configurable rate limits per API key
- Global rate limiting (default: 100 requests/minute)
- Memory-based rate limiter for high performance

## 🎯 Permission System

Available permissions:
- `read` - Read access to resources
- `write` - Write access to resources
- `delete` - Delete access to resources
- `admin` - Administrative access
- `users:read` - Read user data
- `users:write` - Modify user data
- `meals:read` - Read meal data
- `meals:write` - Modify meal data
- `workouts:read` - Read workout data
- `workouts:write` - Modify workout data

## 📊 Usage Tracking

The system automatically tracks:
- Total API key usage count
- Last used timestamp
- Request endpoints and methods
- HTTP status codes
- IP addresses and user agents
- Rate limit usage

## ⚙️ Configuration

Environment variables (see `.env.example`):

```bash
# Server
PORT=8080
ENV=development

# Database
DB_PATH=data/apikeys.db

# API Keys
API_KEY_LENGTH=32
API_KEY_EXPIRY=365d
API_KEY_PREFIX=ak_

# Security
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

## 🛠️ Development

### Available Make Commands

```bash
make build      # Build the application
make run        # Build and run
make test       # Run tests
make deps       # Download dependencies
make setup-dev  # Setup development environment
make db-reset   # Reset database
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
make docs       # Show API documentation
make test-api   # Test API endpoints
make help       # Show all commands
```

### Testing API Endpoints

```bash
# Test health endpoint
curl http://localhost:8080/health

# Get available permissions
curl http://localhost:8080/api/v1/permissions

# Create an API key
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Key",
    "permissions": ["read", "write"],
    "expiry_days": 30
  }'

# Validate API key
curl -X POST http://localhost:8080/api/v1/validate \
  -H "X-API-Key: ak_your_generated_key_here"
```

## 🏆 Performance

- **High Throughput**: Built with Echo v4 for maximum performance
- **Low Latency**: Optimized database queries with proper indexing
- **Memory Efficient**: SQLite with WAL mode for concurrent reads
- **Fast API Key Generation**: Cryptographically secure in microseconds

## 🔒 Security Best Practices

✅ **Never hardcode secrets**  
✅ **Always validate inputs**  
✅ **Use prepared statements**  
✅ **Implement proper error handling**  
✅ **Rate limiting enabled**  
✅ **Security headers configured**  
✅ **Structured logging**  
✅ **Resource cleanup with defer**  

## 📈 Monitoring

The application provides:
- Health check endpoint (`/health`)
- Structured JSON logging
- Request correlation IDs
- Performance metrics
- Usage statistics

## 🤝 Contributing

1. Follow the established patterns in the codebase
2. Run `make fmt` and `make lint` before committing
3. Add tests for new functionality
4. Update documentation as needed

## 📄 License

MIT License - Feel free to use and modify as needed.

---

**Built with ❤️ using Go, Echo, and security best practices**