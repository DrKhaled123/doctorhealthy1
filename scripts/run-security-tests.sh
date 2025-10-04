#!/bin/bash

# Security Test Runner Script
# This script runs comprehensive security tests for the Pure Nutrition API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SECURITY_TEST_DIR="$PROJECT_ROOT/internal/security"
COVERAGE_DIR="$PROJECT_ROOT/coverage"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Pure Nutrition API Security Tests    ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print section headers
print_section() {
    echo -e "${YELLOW}$1${NC}"
    echo "----------------------------------------"
}

# Function to run tests with coverage
run_test_with_coverage() {
    local test_file=$1
    local test_name=$2
    local coverage_file="$COVERAGE_DIR/${test_name}_${TIMESTAMP}.out"
    
    echo -e "Running ${GREEN}$test_name${NC} tests..."
    
    if go test "$test_file" -v -coverprofile="$coverage_file" -covermode=atomic; then
        echo -e "${GREEN}✓${NC} $test_name tests passed"
        
        # Generate coverage report
        local coverage_percent=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
        echo -e "  Coverage: ${GREEN}$coverage_percent${NC}"
    else
        echo -e "${RED}✗${NC} $test_name tests failed"
        return 1
    fi
    echo ""
}

# Function to run all tests in a directory
run_all_security_tests() {
    local coverage_file="$COVERAGE_DIR/security_all_${TIMESTAMP}.out"
    
    echo -e "Running ${GREEN}All Security Tests${NC}..."
    
    if go test "$SECURITY_TEST_DIR/..." -v -coverprofile="$coverage_file" -covermode=atomic; then
        echo -e "${GREEN}✓${NC} All security tests passed"
        
        # Generate coverage report
        local coverage_percent=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
        echo -e "  Overall Security Coverage: ${GREEN}$coverage_percent${NC}"
        
        # Generate HTML coverage report
        local html_report="$COVERAGE_DIR/security_coverage_${TIMESTAMP}.html"
        go tool cover -html="$coverage_file" -o "$html_report"
        echo -e "  HTML Coverage Report: ${BLUE}$html_report${NC}"
    else
        echo -e "${RED}✗${NC} Some security tests failed"
        return 1
    fi
    echo ""
}

# Function to run specific test categories
run_test_category() {
    local pattern=$1
    local category_name=$2
    
    echo -e "Running ${GREEN}$category_name${NC} tests..."
    
    if go test "$SECURITY_TEST_DIR/..." -run "$pattern" -v; then
        echo -e "${GREEN}✓${NC} $category_name tests passed"
    else
        echo -e "${RED}✗${NC} $category_name tests failed"
        return 1
    fi
    echo ""
}

# Function to run benchmark tests
run_benchmark_tests() {
    echo -e "Running ${GREEN}Security Benchmarks${NC}..."
    
    local bench_file="$COVERAGE_DIR/security_bench_${TIMESTAMP}.txt"
    if go test "$SECURITY_TEST_DIR/..." -bench=. -benchmem > "$bench_file" 2>&1; then
        echo -e "${GREEN}✓${NC} Security benchmarks completed"
        echo -e "  Benchmark Results: ${BLUE}$bench_file${NC}"
        
        # Show summary of benchmark results
        echo "Benchmark Summary:"
        grep "Benchmark" "$bench_file" | head -5
    else
        echo -e "${YELLOW}⚠${NC} Benchmark tests had issues (check $bench_file)"
    fi
    echo ""
}

# Function to run race condition tests
run_race_tests() {
    echo -e "Running ${GREEN}Race Condition Tests${NC}..."
    
    if go test "$SECURITY_TEST_DIR/..." -race -v; then
        echo -e "${GREEN}✓${NC} Race condition tests passed"
    else
        echo -e "${RED}✗${NC} Race condition tests failed"
        return 1
    fi
    echo ""
}

# Function to validate test environment
validate_environment() {
    print_section "Environment Validation"
    
    # Check Go version
    local go_version=$(go version | awk '{print $3}')
    echo -e "Go Version: ${GREEN}$go_version${NC}"
    
    # Set test environment variables
    export JWT_SECRET="test-secret-key-for-security-tests"
    export BOOTSTRAP_TOKEN="test-bootstrap-token"
    export RATE_LIMIT_REQUESTS="100"
    export RATE_LIMIT_WINDOW="60s"
    
    echo -e "${GREEN}✓${NC} Environment validation passed"
    echo ""
}

# Function to clean up old coverage files
cleanup_old_coverage() {
    echo "Cleaning up old coverage files..."
    find "$COVERAGE_DIR" -name "*.out" -mtime +7 -delete 2>/dev/null || true
    find "$COVERAGE_DIR" -name "*.html" -mtime +7 -delete 2>/dev/null || true
    echo -e "${GREEN}✓${NC} Cleanup completed"
    echo ""
}

# Main execution logic
main() {
    local test_mode=${1:-"all"}
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    # Validate environment first
    validate_environment
    
    # Clean up old files
    cleanup_old_coverage
    
    case "$test_mode" in
        "all")
            print_section "Running All Security Tests"
            run_all_security_tests
            ;;
        "individual")
            print_section "Running Individual Test Files"
            run_test_with_coverage "$SECURITY_TEST_DIR/api_security_test.go" "API Security"
            run_test_with_coverage "$SECURITY_TEST_DIR/input_validation_test.go" "Input Validation"
            run_test_with_coverage "$SECURITY_TEST_DIR/rate_limit_test.go" "Rate Limiting"
            ;;
        "categories")
            print_section "Running Test Categories"
            run_test_category "TestAPIKey" "API Key Security"
            run_test_category "TestJWT" "JWT Security"
            run_test_category "TestXSS" "XSS Prevention"
            run_test_category "TestSQL" "SQL Injection Prevention"
            run_test_category "TestRateLimit" "Rate Limiting"
            run_test_category "TestQuota" "Quota Management"
            ;;
        "benchmark")
            print_section "Running Benchmark Tests"
            run_benchmark_tests
            ;;
        "race")
            print_section "Running Race Condition Tests"
            run_race_tests
            ;;
        "quick")
            print_section "Running Quick Security Tests"
            echo "Running essential security tests only..."
            go test "$SECURITY_TEST_DIR/..." -v -short
            ;;
        *)
            echo "Usage: $0 [all|individual|categories|benchmark|race|quick]"
            echo ""
            echo "Options:"
            echo "  all         - Run all security tests with coverage (default)"
            echo "  individual  - Run each test file separately with coverage"
            echo "  categories  - Run tests by security category"
            echo "  benchmark   - Run performance benchmarks"
            echo "  race        - Run race condition detection tests"
            echo "  quick       - Run quick essential tests only"
            exit 1
            ;;
    esac
    
    print_section "Security Test Summary"
    echo -e "${GREEN}✓${NC} Security test execution completed"
    echo -e "Coverage reports saved to: ${BLUE}$COVERAGE_DIR${NC}"
    echo -e "Documentation: ${BLUE}$PROJECT_ROOT/SECURITY_TESTING.md${NC}"
}

# Error handling
trap 'echo -e "\n${RED}✗${NC} Security tests interrupted"; exit 1' INT TERM

# Run main function with all arguments
main "$@"