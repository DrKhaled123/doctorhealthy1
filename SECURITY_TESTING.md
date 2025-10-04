# Security Testing Documentation

This document outlines the comprehensive security test suite for the Pure Nutrition API project. The security tests are designed to validate the robustness of the application against common security vulnerabilities and attack vectors.

## Overview

The security test suite is organized into several categories, each targeting specific security concerns:

- **API Security Tests** - Core API endpoint security validation
- **Input Validation Tests** - Protection against malicious input
- **Rate Limiting Tests** - Protection against abuse and DoS attacks
- **Authentication & Authorization Tests** - Access control validation

## Test Files Structure

```
internal/security/
├── api_security_test.go      # Core API security tests
├── input_validation_test.go  # Input sanitization and validation tests
└── rate_limit_test.go       # Rate limiting and quota tests
```

## Security Test Categories

### 1. API Security Tests (`api_security_test.go`)

#### API Key Authentication Security
- **Purpose**: Validates API key authentication mechanisms
- **Test Cases**:
  - Valid API key acceptance
  - Empty API key rejection
  - Invalid API key rejection
  - Malformed API key handling
  - SQL injection attempts in API keys
  - XSS attempts in API keys

#### JWT Token Security
- **Purpose**: Validates JSON Web Token handling
- **Test Cases**:
  - Valid JWT token processing
  - Invalid JWT token rejection
  - Malicious JWT token handling
  - Missing authorization header handling
  - SQL injection attempts in JWT tokens

#### Permission-Based Access Control
- **Purpose**: Validates scope-based authorization
- **Test Cases**:
  - Full access key permissions
  - Limited key permissions
  - Forbidden permission attempts
  - Admin permission restrictions

#### Security Headers
- **Purpose**: Validates security headers implementation
- **Test Cases**:
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `X-XSS-Protection: 1; mode=block`
  - `Strict-Transport-Security`
  - `Referrer-Policy`
  - `Content-Security-Policy`

#### API Key Expiration Security
- **Purpose**: Validates expired API key handling
- **Test Cases**:
  - Expired key rejection
  - Active key acceptance
  - Expiration boundary testing

#### Concurrent Access Security
- **Purpose**: Validates thread safety
- **Test Cases**:
  - Concurrent API key validation
  - Race condition prevention
  - Data integrity validation

### 2. Input Validation Tests (`input_validation_test.go`)

#### Sanitization Functions
- **Purpose**: Validates input sanitization utilities
- **Test Cases**:
  - `SanitizeForLog()` function
  - `SanitizeForHTML()` function
  - `SanitizeForJSON()` function
  - `SanitizeInput()` function

#### XSS Prevention
- **Purpose**: Validates Cross-Site Scripting protection
- **Test Cases**:
  - Script tag injection attempts
  - Event handler injection
  - JavaScript protocol attempts
  - SVG-based XSS attempts
  - Unicode-based XSS bypasses

#### SQL Injection Prevention
- **Purpose**: Validates SQL injection protection
- **Test Cases**:
  - Classic SQL injection patterns
  - Union-based injections
  - Boolean-based blind injections
  - Time-based blind injections
  - Comment-based injections

#### API Key Creation Input Validation
- **Purpose**: Validates API key creation security
- **Test Cases**:
  - XSS in name fields
  - SQL injection in name fields
  - Empty field validation
  - Length validation
  - Invalid permission handling

#### HTTP Header Injection
- **Purpose**: Validates header injection protection
- **Test Cases**:
  - CRLF injection attempts
  - Content-Type manipulation
  - Cookie injection attempts
  - URL-encoded injection attempts

#### JSON Injection
- **Purpose**: Validates JSON parsing security
- **Test Cases**:
  - Prototype pollution attempts
  - Constructor pollution attempts
  - Large number handling
  - Special character handling

#### Path Traversal Prevention
- **Purpose**: Validates path traversal protection
- **Test Cases**:
  - Directory traversal attempts
  - URL-encoded traversal
  - Unicode traversal attempts
  - Mixed encoding attempts

#### NoSQL Injection Prevention
- **Purpose**: Validates NoSQL injection protection
- **Test Cases**:
  - MongoDB operator injection
  - Query manipulation attempts
  - Logical operator bypasses

### 3. Rate Limiting Tests (`rate_limit_test.go`)

#### Rate Limiting Middleware
- **Purpose**: Validates request rate limiting
- **Test Cases**:
  - Basic rate limit enforcement
  - Rate limit reset after time window
  - Per-IP rate limiting
  - Burst request handling

#### User-Based Rate Limiting
- **Purpose**: Validates per-user rate limits
- **Test Cases**:
  - Authenticated user rate limits
  - Unauthenticated request handling
  - Per-user isolation
  - Rate limit bypass prevention

#### Monthly Quota Middleware
- **Purpose**: Validates monthly usage quotas
- **Test Cases**:
  - Free plan quota enforcement (3/month)
  - Pro plan quota enforcement (50/month)
  - Lifetime plan quota handling (unlimited)
  - Shared plan bonus validation (11/month)

#### Concurrent Rate Limiting
- **Purpose**: Validates thread-safe rate limiting
- **Test Cases**:
  - Concurrent request handling
  - Race condition prevention
  - Accurate rate counting

## Security Vulnerabilities Covered

### OWASP Top 10 Coverage

1. **A01:2021 – Broken Access Control**
   - ✅ API key validation
   - ✅ Permission-based access control
   - ✅ JWT token validation
   - ✅ Rate limiting and quotas

2. **A02:2021 – Cryptographic Failures**
   - ✅ Secure API key generation
   - ✅ JWT token validation
   - ✅ Secure headers implementation

3. **A03:2021 – Injection**
   - ✅ SQL injection prevention
   - ✅ NoSQL injection prevention
   - ✅ Command injection prevention
   - ✅ LDAP injection prevention

4. **A04:2021 – Insecure Design**
   - ✅ Rate limiting implementation
   - ✅ Quota system design
   - ✅ Permission model validation

5. **A05:2021 – Security Misconfiguration**
   - ✅ Security headers validation
   - ✅ Error handling testing
   - ✅ Configuration validation

6. **A06:2021 – Vulnerable and Outdated Components**
   - ✅ Input validation libraries
   - ✅ JWT library security

7. **A07:2021 – Identification and Authentication Failures**
   - ✅ API key authentication
   - ✅ JWT token validation
   - ✅ Session management

8. **A08:2021 – Software and Data Integrity Failures**
   - ✅ Input validation
   - ✅ Data sanitization
   - ✅ API response integrity

9. **A09:2021 – Security Logging and Monitoring Failures**
   - ✅ Log sanitization
   - ✅ Error logging validation

10. **A10:2021 – Server-Side Request Forgery (SSRF)**
    - ✅ Input validation
    - ✅ URL validation

### Additional Security Measures

- **HTTP Header Injection Prevention**
- **Path Traversal Protection**
- **XML External Entity (XXE) Prevention**
- **Prototype Pollution Prevention**
- **Unicode Security Validation**
- **Control Character Filtering**

## Running Security Tests

### Prerequisites

```bash
go mod tidy
```

### Running All Security Tests

```bash
# Run all security tests
go test ./internal/security/... -v

# Run with coverage
go test ./internal/security/... -v -cover

# Run specific test file
go test ./internal/security/api_security_test.go -v
go test ./internal/security/input_validation_test.go -v
go test ./internal/security/rate_limit_test.go -v
```

### Running Specific Test Categories

```bash
# API Security Tests
go test ./internal/security/ -run TestAPIKey -v
go test ./internal/security/ -run TestJWT -v
go test ./internal/security/ -run TestPermission -v

# Input Validation Tests
go test ./internal/security/ -run TestXSS -v
go test ./internal/security/ -run TestSQL -v
go test ./internal/security/ -run TestSanitization -v

# Rate Limiting Tests
go test ./internal/security/ -run TestRateLimit -v
go test ./internal/security/ -run TestQuota -v
```

## Test Configuration

### Environment Variables

```bash
# JWT Secret for testing
export JWT_SECRET="test-secret-key-for-security-tests"

# Bootstrap token for API key creation
export BOOTSTRAP_TOKEN="test-bootstrap-token"

# Rate limiting configuration
export RATE_LIMIT_REQUESTS=100
export RATE_LIMIT_WINDOW=60s
```

### Test Database

Tests use in-memory SQLite databases to ensure isolation and fast execution.

## Security Test Metrics

### Coverage Goals

- **API Security**: 100% of authentication/authorization paths
- **Input Validation**: 95%+ of input processing functions
- **Rate Limiting**: 100% of rate limiting middleware
- **Overall Security**: 90%+ security-related code coverage

### Performance Benchmarks

- Security tests should complete within 30 seconds
- Individual test categories should complete within 10 seconds
- Memory usage should remain under 100MB during testing

## Security Test Maintenance

### Regular Updates

1. **Monthly Review**: Review and update test cases based on new threats
2. **Dependency Updates**: Update security testing dependencies
3. **Vulnerability Research**: Incorporate new attack vectors
4. **Performance Monitoring**: Ensure tests remain fast and reliable

### Adding New Security Tests

1. Identify the security concern or vulnerability
2. Create test cases covering both positive and negative scenarios
3. Add appropriate assertions and error checking
4. Update this documentation
5. Ensure tests are deterministic and isolated

## Integration with CI/CD

### GitHub Actions Integration

```yaml
name: Security Tests
on: [push, pull_request]
jobs:
  security-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.20'
      - name: Run Security Tests
        run: |
          go test ./internal/security/... -v -cover
          go test ./internal/security/... -bench=. -benchmem
```

### Pre-commit Hooks

```bash
# Install pre-commit hook
cat > .pre-commit-config.yaml << EOF
repos:
  - repo: local
    hooks:
      - id: security-tests
        name: Security Tests
        entry: go test ./internal/security/... -v
        language: system
        pass_filenames: false
EOF
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Ensure SQLite is available
   - Check file permissions for test database

2. **Rate Limiting Test Flakiness**
   - Tests may be timing-sensitive
   - Adjust time windows if needed

3. **Compilation Errors**
   - Verify all imports are available
   - Check Go version compatibility

### Debug Mode

```bash
# Run tests with debug output
go test ./internal/security/... -v -debug

# Run with race detection
go test ./internal/security/... -race -v
```

## Security Contact

For security-related issues or questions about the security test suite:

- Create an issue in the project repository
- Follow responsible disclosure practices
- Include detailed reproduction steps

---

*This document is maintained alongside the security test suite and should be updated whenever new security tests are added or existing tests are modified.*