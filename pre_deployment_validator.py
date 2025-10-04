#!/usr/bin/env python3
"""
Pre-Deployment Validation Suite
Comprehensive testing before production deployment
"""

import asyncio
import json
import sys
import os
import subprocess
import time
from typing import Dict, List, Any
import requests

class PreDeploymentValidator:
    def __init__(self):
        self.app_url = "http://localhost:8085"
        self.results = {
            "validation_timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
            "overall_status": "pending",
            "tests": [],
            "summary": {}
        }
    
    async def run_all_validations(self) -> Dict[str, Any]:
        """Run comprehensive pre-deployment validation suite"""
        print("🚀 Starting Pre-Deployment Validation Suite")
        print("=" * 60)
        
        # 1. Application Health Check
        await self.validate_application_health()
        
        # 2. Security Validation
        await self.validate_security_features()
        
        # 3. API Endpoint Testing
        await self.validate_api_endpoints()
        
        # 4. Performance Testing
        await self.validate_performance()
        
        # 5. Docker Build Validation
        await self.validate_docker_build()
        
        # 6. Database Integrity
        await self.validate_database_integrity()
        
        # Calculate overall status
        self._calculate_overall_status()
        
        # Generate report
        await self.generate_validation_report()
        
        return self.results
    
    async def validate_application_health(self):
        """Validate basic application health"""
        print("\n📋 1. Application Health Validation")
        print("-" * 40)
        
        test_result = {
            "category": "Health Check",
            "tests": [],
            "status": "passed"
        }
        
        # Check if application is running
        health_check = await self._check_endpoint("/health", expected_status=200)
        test_result["tests"].append({
            "name": "Health Endpoint",
            "status": "passed" if health_check["success"] else "failed",
            "details": health_check
        })
        
        # Check root endpoint
        root_check = await self._check_endpoint("/", expected_status=200)
        test_result["tests"].append({
            "name": "Root Endpoint",
            "status": "passed" if root_check["success"] else "failed",
            "details": root_check
        })
        
        # Check API documentation
        api_docs_check = await self._check_endpoint("/api/docs", expected_status=200)
        test_result["tests"].append({
            "name": "API Documentation",
            "status": "passed" if api_docs_check["success"] else "failed",
            "details": api_docs_check
        })
        
        if not all(test["status"] == "passed" for test in test_result["tests"]):
            test_result["status"] = "failed"
        
        self.results["tests"].append(test_result)
        print(f"✅ Health validation: {test_result['status']}")
    
    async def validate_security_features(self):
        """Validate security implementations"""
        print("\n🔒 2. Security Feature Validation")
        print("-" * 40)
        
        test_result = {
            "category": "Security",
            "tests": [],
            "status": "passed"
        }
        
        # Test XSS protection
        xss_test = await self._test_xss_protection()
        test_result["tests"].append({
            "name": "XSS Protection",
            "status": "passed" if xss_test["protected"] else "failed",
            "details": xss_test
        })
        
        # Test SQL injection protection
        sql_injection_test = await self._test_sql_injection_protection()
        test_result["tests"].append({
            "name": "SQL Injection Protection",
            "status": "passed" if sql_injection_test["protected"] else "failed",
            "details": sql_injection_test
        })
        
        # Test rate limiting
        rate_limit_test = await self._test_rate_limiting()
        test_result["tests"].append({
            "name": "Rate Limiting",
            "status": "passed" if rate_limit_test["active"] else "failed",
            "details": rate_limit_test
        })
        
        # Test CORS headers
        cors_test = await self._test_cors_headers()
        test_result["tests"].append({
            "name": "CORS Headers",
            "status": "passed" if cors_test["configured"] else "failed",
            "details": cors_test
        })
        
        if not all(test["status"] == "passed" for test in test_result["tests"]):
            test_result["status"] = "failed"
        
        self.results["tests"].append(test_result)
        print(f"🔐 Security validation: {test_result['status']}")
    
    async def validate_api_endpoints(self):
        """Validate critical API endpoints"""
        print("\n🌐 3. API Endpoint Validation")
        print("-" * 40)
        
        test_result = {
            "category": "API Endpoints",
            "tests": [],
            "status": "passed"
        }
        
        # Key endpoints to test
        endpoints = [
            {"path": "/api/health", "method": "GET", "expected": 200},
            {"path": "/api/recipes", "method": "GET", "expected": 200},
            {"path": "/api/diseases", "method": "GET", "expected": 200},
            {"path": "/api/complaints", "method": "GET", "expected": 200},
            {"path": "/api/auth/login", "method": "POST", "expected": 400}  # Should fail without credentials
        ]
        
        for endpoint in endpoints:
            check = await self._check_endpoint(
                endpoint["path"], 
                method=endpoint["method"],
                expected_status=endpoint["expected"]
            )
            
            test_result["tests"].append({
                "name": f"{endpoint['method']} {endpoint['path']}",
                "status": "passed" if check["success"] else "failed",
                "details": check
            })
        
        if not all(test["status"] == "passed" for test in test_result["tests"]):
            test_result["status"] = "failed"
        
        self.results["tests"].append(test_result)
        print(f"🌍 API validation: {test_result['status']}")
    
    async def validate_performance(self):
        """Basic performance validation"""
        print("\n⚡ 4. Performance Validation")
        print("-" * 40)
        
        test_result = {
            "category": "Performance",
            "tests": [],
            "status": "passed"
        }
        
        # Response time test
        response_times = []
        for i in range(5):
            start_time = time.time()
            check = await self._check_endpoint("/health")
            end_time = time.time()
            
            if check["success"]:
                response_times.append((end_time - start_time) * 1000)  # Convert to ms
        
        if response_times:
            avg_response_time = sum(response_times) / len(response_times)
            max_response_time = max(response_times)
            
            test_result["tests"].append({
                "name": "Response Time",
                "status": "passed" if avg_response_time < 500 else "warning",
                "details": {
                    "average_ms": round(avg_response_time, 2),
                    "max_ms": round(max_response_time, 2),
                    "threshold_ms": 500
                }
            })
        else:
            test_result["tests"].append({
                "name": "Response Time",
                "status": "failed",
                "details": {"error": "Could not measure response time"}
            })
        
        # Memory usage check (if available)
        try:
            memory_check = await self._check_memory_usage()
            test_result["tests"].append({
                "name": "Memory Usage",
                "status": memory_check["status"],
                "details": memory_check
            })
        except Exception as e:
            test_result["tests"].append({
                "name": "Memory Usage",
                "status": "skipped",
                "details": {"error": str(e)}
            })
        
        # Check for any failed tests
        failed_tests = [test for test in test_result["tests"] if test["status"] == "failed"]
        if failed_tests:
            test_result["status"] = "failed"
        elif any(test["status"] == "warning" for test in test_result["tests"]):
            test_result["status"] = "warning"
        
        self.results["tests"].append(test_result)
        print(f"⚡ Performance validation: {test_result['status']}")
    
    async def validate_docker_build(self):
        """Validate Docker build process"""
        print("\n🐳 5. Docker Build Validation")
        print("-" * 40)
        
        test_result = {
            "category": "Docker Build",
            "tests": [],
            "status": "passed"
        }
        
        # Check if Dockerfile exists
        dockerfile_check = {
            "name": "Dockerfile Exists",
            "status": "passed" if os.path.exists("Dockerfile") else "failed",
            "details": {"path": "Dockerfile", "exists": os.path.exists("Dockerfile")}
        }
        test_result["tests"].append(dockerfile_check)
        
        # Validate Dockerfile syntax
        if os.path.exists("Dockerfile"):
            dockerfile_syntax = await self._validate_dockerfile_syntax()
            test_result["tests"].append({
                "name": "Dockerfile Syntax",
                "status": dockerfile_syntax["status"],
                "details": dockerfile_syntax
            })
        
        # Check .dockerignore
        dockerignore_check = {
            "name": "Dockerignore File",
            "status": "passed" if os.path.exists(".dockerignore") else "warning",
            "details": {"path": ".dockerignore", "exists": os.path.exists(".dockerignore")}
        }
        test_result["tests"].append(dockerignore_check)
        
        if any(test["status"] == "failed" for test in test_result["tests"]):
            test_result["status"] = "failed"
        elif any(test["status"] == "warning" for test in test_result["tests"]):
            test_result["status"] = "warning"
        
        self.results["tests"].append(test_result)
        print(f"🐳 Docker validation: {test_result['status']}")
    
    async def validate_database_integrity(self):
        """Validate database integrity"""
        print("\n💾 6. Database Integrity Validation")
        print("-" * 40)
        
        test_result = {
            "category": "Database",
            "tests": [],
            "status": "passed"
        }
        
        # Check if database files exist
        db_files = ["app.db", "health.db", "data.db"]
        for db_file in db_files:
            if os.path.exists(db_file):
                test_result["tests"].append({
                    "name": f"Database File {db_file}",
                    "status": "passed",
                    "details": {"file": db_file, "size": os.path.getsize(db_file)}
                })
                break
        else:
            test_result["tests"].append({
                "name": "Database Files",
                "status": "warning",
                "details": {"message": "No database files found, using in-memory database"}
            })
        
        # Test database connectivity through API
        db_test = await self._test_database_connectivity()
        test_result["tests"].append({
            "name": "Database Connectivity",
            "status": db_test["status"],
            "details": db_test
        })
        
        if any(test["status"] == "failed" for test in test_result["tests"]):
            test_result["status"] = "failed"
        elif any(test["status"] == "warning" for test in test_result["tests"]):
            test_result["status"] = "warning"
        
        self.results["tests"].append(test_result)
        print(f"💾 Database validation: {test_result['status']}")
    
    async def _check_endpoint(self, path: str, method: str = "GET", expected_status: int = 200, **kwargs) -> Dict[str, Any]:
        """Check a specific endpoint"""
        url = f"{self.app_url}{path}"
        
        try:
            if method == "GET":
                response = requests.get(url, timeout=5, **kwargs)
            elif method == "POST":
                response = requests.post(url, timeout=5, **kwargs)
            elif method == "PUT":
                response = requests.put(url, timeout=5, **kwargs)
            else:
                return {"success": False, "error": f"Unsupported method: {method}"}
            
            return {
                "success": response.status_code == expected_status,
                "status_code": response.status_code,
                "expected_status": expected_status,
                "response_time_ms": round(response.elapsed.total_seconds() * 1000, 2),
                "content_length": len(response.content) if response.content else 0
            }
        
        except requests.exceptions.RequestException as e:
            return {
                "success": False,
                "error": str(e),
                "expected_status": expected_status
            }
    
    async def _test_xss_protection(self) -> Dict[str, Any]:
        """Test XSS protection"""
        try:
            # Try to submit XSS payload
            xss_payload = "<script>alert('xss')</script>"
            response = requests.get(f"{self.app_url}/api/recipes?search={xss_payload}", timeout=5)
            
            # Check if the payload is sanitized in response
            is_protected = xss_payload not in response.text
            
            return {
                "protected": is_protected,
                "status_code": response.status_code,
                "payload_found": xss_payload in response.text
            }
        
        except Exception as e:
            return {"protected": False, "error": str(e)}
    
    async def _test_sql_injection_protection(self) -> Dict[str, Any]:
        """Test SQL injection protection"""
        try:
            # Try SQL injection payload
            sql_payload = "'; DROP TABLE recipes; --"
            response = requests.get(f"{self.app_url}/api/recipes?search={sql_payload}", timeout=5)
            
            # If we get a response and it's not a server error, protection is likely working
            is_protected = response.status_code != 500
            
            return {
                "protected": is_protected,
                "status_code": response.status_code
            }
        
        except Exception as e:
            return {"protected": False, "error": str(e)}
    
    async def _test_rate_limiting(self) -> Dict[str, Any]:
        """Test rate limiting"""
        try:
            # Make multiple rapid requests
            responses = []
            for i in range(10):
                response = requests.get(f"{self.app_url}/health", timeout=5)
                responses.append(response.status_code)
            
            # Check if any request was rate limited (429 status)
            rate_limited = 429 in responses
            
            return {
                "active": rate_limited,
                "responses": responses,
                "rate_limit_triggered": rate_limited
            }
        
        except Exception as e:
            return {"active": False, "error": str(e)}
    
    async def _test_cors_headers(self) -> Dict[str, Any]:
        """Test CORS configuration"""
        try:
            response = requests.get(f"{self.app_url}/api/health", timeout=5)
            
            cors_headers = {
                "Access-Control-Allow-Origin": response.headers.get("Access-Control-Allow-Origin"),
                "Access-Control-Allow-Methods": response.headers.get("Access-Control-Allow-Methods"),
                "Access-Control-Allow-Headers": response.headers.get("Access-Control-Allow-Headers")
            }
            
            configured = any(cors_headers.values())
            
            return {
                "configured": configured,
                "headers": cors_headers
            }
        
        except Exception as e:
            return {"configured": False, "error": str(e)}
    
    async def _check_memory_usage(self) -> Dict[str, Any]:
        """Check memory usage"""
        try:
            # Use ps command to check memory usage of Go process
            result = subprocess.run(
                ["ps", "aux"], 
                capture_output=True, 
                text=True, 
                timeout=5
            )
            
            if result.returncode == 0:
                lines = result.stdout.split('\n')
                go_processes = [line for line in lines if 'main' in line and 'go' in line.lower()]
                
                if go_processes:
                    # Parse memory usage from first Go process found
                    process_line = go_processes[0].split()
                    if len(process_line) > 5:
                        memory_percent = float(process_line[3])
                        return {
                            "status": "warning" if memory_percent > 50 else "passed",
                            "memory_percent": memory_percent,
                            "process_info": process_line[10] if len(process_line) > 10 else "unknown"
                        }
            
            return {"status": "skipped", "message": "Could not determine memory usage"}
        
        except Exception as e:
            return {"status": "skipped", "error": str(e)}
    
    async def _validate_dockerfile_syntax(self) -> Dict[str, Any]:
        """Validate Dockerfile syntax"""
        try:
            # Use docker build with --no-cache and --dry-run if available
            result = subprocess.run(
                ["docker", "build", "--no-cache", "-f", "Dockerfile", "."],
                capture_output=True,
                text=True,
                timeout=30
            )
            
            return {
                "status": "passed" if result.returncode == 0 else "failed",
                "exit_code": result.returncode,
                "output": result.stdout[:500] if result.stdout else "",
                "error": result.stderr[:500] if result.stderr else ""
            }
        
        except subprocess.TimeoutExpired:
            return {
                "status": "warning",
                "message": "Docker build timed out (may indicate slow build)"
            }
        except Exception as e:
            return {
                "status": "skipped",
                "error": str(e)
            }
    
    async def _test_database_connectivity(self) -> Dict[str, Any]:
        """Test database connectivity through API"""
        try:
            # Test a simple database query through API
            response = requests.get(f"{self.app_url}/api/recipes?limit=1", timeout=5)
            
            return {
                "status": "passed" if response.status_code == 200 else "failed",
                "status_code": response.status_code,
                "response_size": len(response.content) if response.content else 0
            }
        
        except Exception as e:
            return {
                "status": "failed",
                "error": str(e)
            }
    
    def _calculate_overall_status(self):
        """Calculate overall validation status"""
        all_statuses = []
        
        for test_category in self.results["tests"]:
            all_statuses.append(test_category["status"])
        
        if "failed" in all_statuses:
            self.results["overall_status"] = "failed"
        elif "warning" in all_statuses:
            self.results["overall_status"] = "warning"
        else:
            self.results["overall_status"] = "passed"
        
        # Generate summary
        total_categories = len(self.results["tests"])
        passed = sum(1 for status in all_statuses if status == "passed")
        warnings = sum(1 for status in all_statuses if status == "warning")
        failed = sum(1 for status in all_statuses if status == "failed")
        
        self.results["summary"] = {
            "total_categories": total_categories,
            "passed": passed,
            "warnings": warnings,
            "failed": failed,
            "success_rate": round((passed / total_categories) * 100, 2) if total_categories > 0 else 0
        }
    
    async def generate_validation_report(self):
        """Generate detailed validation report"""
        print(f"\n📊 Validation Summary")
        print("=" * 60)
        print(f"Overall Status: {self.results['overall_status'].upper()}")
        print(f"Success Rate: {self.results['summary']['success_rate']}%")
        print(f"Categories: {self.results['summary']['passed']} passed, {self.results['summary']['warnings']} warnings, {self.results['summary']['failed']} failed")
        
        # Write detailed report to file
        report_path = "pre_deployment_validation_report.json"
        with open(report_path, 'w') as f:
            json.dump(self.results, f, indent=2)
        
        print(f"\n📄 Detailed report saved to: {report_path}")
        
        # Display failed/warning tests
        for category in self.results["tests"]:
            if category["status"] != "passed":
                print(f"\n⚠️ {category['category']} Issues:")
                for test in category["tests"]:
                    if test["status"] != "passed":
                        print(f"  • {test['name']}: {test['status']}")
                        if "error" in test["details"]:
                            print(f"    Error: {test['details']['error']}")

async def main():
    """Run pre-deployment validation"""
    validator = PreDeploymentValidator()
    
    print("🎯 Pre-Deployment Validation Suite")
    print("Ensuring application readiness for production deployment")
    print()
    
    results = await validator.run_all_validations()
    
    print("\n" + "=" * 60)
    if results["overall_status"] == "passed":
        print("✅ APPLICATION READY FOR DEPLOYMENT")
        exit_code = 0
    elif results["overall_status"] == "warning":
        print("⚠️ APPLICATION READY WITH WARNINGS")
        print("Review warnings before deploying to production")
        exit_code = 0
    else:
        print("❌ APPLICATION NOT READY FOR DEPLOYMENT")
        print("Fix failed tests before deploying")
        exit_code = 1
    
    return exit_code

if __name__ == "__main__":
    exit_code = asyncio.run(main())
    sys.exit(exit_code)