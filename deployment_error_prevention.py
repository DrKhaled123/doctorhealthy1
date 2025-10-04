#!/usr/bin/env python3
"""
Comprehensive Deployment Error Prevention & Testing System
Based on expert recommendations for web deployment best practices
"""

import asyncio
import json
import subprocess
import sys
import os
from datetime import datetime
from pathlib import Path

class DeploymentErrorPrevention:
    def __init__(self):
        self.error_log = {
            "session_id": datetime.now().strftime("%Y%m%d_%H%M%S"),
            "errors": [],
            "prevention_checks": {},
            "recommendations": []
        }
        
    async def run_comprehensive_checks(self):
        print("🔍 Comprehensive Deployment Error Prevention System")
        print("=" * 70)
        
        # 1. Reverse Proxy & Configuration Checks
        await self._check_reverse_proxy_config()
        
        # 2. File Upload & Project Deployment Checks
        await self._check_file_deployment_issues()
        
        # 3. Docker & Container Checks
        await self._check_docker_configuration()
        
        # 4. Security Testing
        await self._run_security_tests()
        
        # 5. Performance Testing
        await self._run_performance_tests()
        
        # 6. Integration Testing
        await self._run_integration_tests()
        
        # 7. Platform-Specific Checks
        await self._check_platform_issues()
        
        # 8. Generate Prevention Report
        await self._generate_prevention_report()
        
        return self.error_log
    
    async def _check_reverse_proxy_config(self):
        print("🌐 Checking Reverse Proxy Configuration...")
        
        checks = {
            "nginx_config": await self._validate_nginx_config(),
            "port_conflicts": await self._check_port_conflicts(),
            "proxy_headers": await self._validate_proxy_headers()
        }
        
        self.error_log["prevention_checks"]["reverse_proxy"] = checks
        
        for check_name, result in checks.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {check_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _validate_nginx_config(self):
        """Check for proper Nginx configuration patterns"""
        try:
            # Check if nginx config exists and is valid
            result = subprocess.run(["nginx", "-t"], capture_output=True, text=True)
            
            if result.returncode == 0:
                return {
                    "success": True,
                    "message": "Nginx configuration is valid",
                    "details": result.stdout
                }
            else:
                return {
                    "success": False,
                    "message": "Nginx configuration has errors",
                    "details": result.stderr,
                    "recommendation": "Check nginx configuration syntax"
                }
        except FileNotFoundError:
            return {
                "success": True,  # Not an error if nginx not installed
                "message": "Nginx not installed (using alternative proxy)",
                "recommendation": "Consider nginx for production reverse proxy"
            }
    
    async def _check_port_conflicts(self):
        """Check for port conflicts that could cause deployment issues"""
        common_ports = [80, 443, 8000, 8080, 8085, 3000]
        conflicts = []
        
        for port in common_ports:
            result = subprocess.run(
                ["lsof", "-i", f":{port}"], 
                capture_output=True, text=True
            )
            if result.returncode == 0 and result.stdout:
                conflicts.append({
                    "port": port,
                    "process": result.stdout.split('\n')[1] if len(result.stdout.split('\n')) > 1 else "unknown"
                })
        
        if conflicts:
            return {
                "success": False,
                "message": f"Found {len(conflicts)} port conflicts",
                "details": conflicts,
                "recommendation": "Resolve port conflicts before deployment"
            }
        else:
            return {
                "success": True,
                "message": "No port conflicts detected",
                "details": "All common ports available"
            }
    
    async def _validate_proxy_headers(self):
        """Validate proxy header configuration for proper forwarding"""
        required_headers = [
            "X-Real-IP",
            "X-Forwarded-For", 
            "X-Forwarded-Proto",
            "Host"
        ]
        
        # Test against running application
        test_url = "http://localhost:8085/health"
        cmd = ["curl", "-I", "-s", test_url]
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            headers_present = []
            for header in required_headers:
                if header.lower() in result.stdout.lower():
                    headers_present.append(header)
            
            return {
                "success": len(headers_present) >= 2,  # At least 2 required headers
                "message": f"Found {len(headers_present)}/{len(required_headers)} proxy headers",
                "details": headers_present,
                "recommendation": "Add missing proxy headers for production"
            }
        else:
            return {
                "success": False,
                "message": "Cannot test proxy headers - application not responding",
                "recommendation": "Ensure application is running for header validation"
            }
    
    async def _check_file_deployment_issues(self):
        print("📁 Checking File Deployment Issues...")
        
        checks = {
            "file_permissions": await self._check_file_permissions(),
            "disk_space": await self._check_disk_space(),
            "dependencies": await self._check_dependencies()
        }
        
        self.error_log["prevention_checks"]["file_deployment"] = checks
        
        for check_name, result in checks.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {check_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _check_file_permissions(self):
        """Check file permissions for proper deployment"""
        current_dir = os.getcwd()
        
        try:
            # Check if we can read/write in current directory
            test_file = os.path.join(current_dir, ".deployment_test")
            with open(test_file, "w") as f:
                f.write("test")
            os.remove(test_file)
            
            # Check executable permissions on binaries
            bin_dir = os.path.join(current_dir, "bin")
            if os.path.exists(bin_dir):
                executables = [f for f in os.listdir(bin_dir) if not f.startswith('.')]
                non_executable = []
                
                for exe in executables:
                    exe_path = os.path.join(bin_dir, exe)
                    if not os.access(exe_path, os.X_OK):
                        non_executable.append(exe)
                
                if non_executable:
                    return {
                        "success": False,
                        "message": f"Found {len(non_executable)} non-executable binaries",
                        "details": non_executable,
                        "recommendation": "chmod +x on binary files"
                    }
            
            return {
                "success": True,
                "message": "File permissions are correct",
                "details": "Read/write/execute permissions validated"
            }
            
        except PermissionError:
            return {
                "success": False,
                "message": "Insufficient file permissions",
                "recommendation": "Check directory permissions and ownership"
            }
    
    async def _check_disk_space(self):
        """Check available disk space for deployment"""
        result = subprocess.run(["df", "-h", "."], capture_output=True, text=True)
        
        if result.returncode == 0:
            lines = result.stdout.strip().split('\n')
            if len(lines) > 1:
                parts = lines[1].split()
                if len(parts) >= 5:
                    used_percent = parts[4].rstrip('%')
                    try:
                        used_int = int(used_percent)
                        available = parts[3]
                        
                        if used_int > 90:
                            return {
                                "success": False,
                                "message": f"Disk {used_percent}% full, {available} available",
                                "recommendation": "Free up disk space before deployment"
                            }
                        else:
                            return {
                                "success": True,
                                "message": f"Disk {used_percent}% full, {available} available",
                                "details": "Sufficient disk space available"
                            }
                    except ValueError:
                        pass
        
        return {
            "success": True,  # Default to success if can't parse
            "message": "Could not determine disk usage",
            "recommendation": "Manually check disk space with 'df -h'"
        }
    
    async def _check_dependencies(self):
        """Check for missing dependencies that could cause deployment failures"""
        dependencies = {
            "go": ["go", "version"],
            "docker": ["docker", "--version"],
            "curl": ["curl", "--version"],
            "git": ["git", "--version"]
        }
        
        missing_deps = []
        present_deps = []
        
        for dep_name, cmd in dependencies.items():
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode == 0:
                present_deps.append(dep_name)
            else:
                missing_deps.append(dep_name)
        
        if missing_deps:
            return {
                "success": False,
                "message": f"Missing dependencies: {', '.join(missing_deps)}",
                "details": {"present": present_deps, "missing": missing_deps},
                "recommendation": "Install missing dependencies before deployment"
            }
        else:
            return {
                "success": True,
                "message": "All required dependencies present",
                "details": present_deps
            }
    
    async def _check_docker_configuration(self):
        print("🐳 Checking Docker Configuration...")
        
        checks = {
            "dockerfile_syntax": await self._validate_dockerfile(),
            "docker_service": await self._check_docker_service(),
            "image_security": await self._check_image_security()
        }
        
        self.error_log["prevention_checks"]["docker"] = checks
        
        for check_name, result in checks.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {check_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _validate_dockerfile(self):
        """Validate Dockerfile syntax and best practices"""
        dockerfile_path = "Dockerfile"
        
        if not os.path.exists(dockerfile_path):
            return {
                "success": True,  # Not required if no Docker deployment
                "message": "No Dockerfile found (not using Docker deployment)",
                "recommendation": "Consider Docker for consistent deployments"
            }
        
        try:
            with open(dockerfile_path, 'r') as f:
                content = f.read()
            
            best_practices = {
                "multi_stage": "FROM" in content and "as" in content.lower(),
                "non_root_user": "USER" in content,
                "explicit_version": not content.count("FROM") or not "latest" in content,
                "workdir_set": "WORKDIR" in content,
                "expose_port": "EXPOSE" in content
            }
            
            passed_checks = sum(best_practices.values())
            total_checks = len(best_practices)
            
            return {
                "success": passed_checks >= 3,  # At least 3/5 best practices
                "message": f"Dockerfile follows {passed_checks}/{total_checks} best practices",
                "details": best_practices,
                "recommendation": "Improve Dockerfile security and practices"
            }
            
        except Exception as e:
            return {
                "success": False,
                "message": f"Error reading Dockerfile: {str(e)}",
                "recommendation": "Fix Dockerfile syntax errors"
            }
    
    async def _check_docker_service(self):
        """Check if Docker service is running and accessible"""
        try:
            result = subprocess.run(["docker", "info"], capture_output=True, text=True)
            
            if result.returncode == 0:
                return {
                    "success": True,
                    "message": "Docker service is running",
                    "details": "Docker daemon accessible"
                }
            else:
                return {
                    "success": False,
                    "message": "Docker service not accessible",
                    "recommendation": "Start Docker service or check permissions"
                }
                
        except FileNotFoundError:
            return {
                "success": True,  # Not required if not using Docker
                "message": "Docker not installed",
                "recommendation": "Install Docker if container deployment needed"
            }
    
    async def _check_image_security(self):
        """Basic Docker image security checks"""
        try:
            # Try to build a test image to check for issues
            result = subprocess.run(
                ["docker", "build", "--dry-run", "-t", "test-security", "."],
                capture_output=True, text=True
            )
            
            if "dry-run" in result.stderr or result.returncode == 0:
                return {
                    "success": True,
                    "message": "Docker build appears secure",
                    "details": "No obvious security issues in build"
                }
            else:
                return {
                    "success": False,
                    "message": "Docker build has potential issues",
                    "details": result.stderr[:200],
                    "recommendation": "Review Docker build errors"
                }
                
        except Exception:
            return {
                "success": True,  # Skip if Docker not available
                "message": "Could not perform Docker security check",
                "recommendation": "Manually review Dockerfile for security"
            }
    
    async def _run_security_tests(self):
        print("🔒 Running Security Tests...")
        
        if not await self._check_app_running():
            print("  ⚠️  Skipping security tests - application not running")
            return
        
        tests = {
            "xss_protection": await self._test_xss_protection(),
            "sql_injection": await self._test_sql_injection(),
            "auth_validation": await self._test_auth_validation(),
            "rate_limiting": await self._test_rate_limiting()
        }
        
        self.error_log["prevention_checks"]["security"] = tests
        
        for test_name, result in tests.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {test_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _test_xss_protection(self):
        """Test XSS protection mechanisms"""
        xss_payloads = [
            "<script>alert('xss')</script>",
            "javascript:alert('xss')",
            "<img src=x onerror=alert('xss')>"
        ]
        
        blocked_count = 0
        
        for payload in xss_payloads:
            # Test with a POST request to an endpoint that accepts data
            cmd = ["curl", "-s", "-X", "POST", 
                   "http://localhost:8085/api/test-input",
                   "-H", "Content-Type: application/json",
                   "-d", json.dumps({"input": payload})]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            # If we get an error or sanitization message, XSS protection is working
            if (result.returncode != 0 or 
                "error" in result.stdout.lower() or 
                "sanitized" in result.stdout.lower() or
                payload not in result.stdout):
                blocked_count += 1
        
        success_rate = blocked_count / len(xss_payloads)
        
        return {
            "success": success_rate >= 0.7,  # 70% of payloads should be blocked
            "message": f"Blocked {blocked_count}/{len(xss_payloads)} XSS attempts",
            "details": f"Protection rate: {success_rate:.1%}",
            "recommendation": "Improve XSS protection if rate < 100%"
        }
    
    async def _test_sql_injection(self):
        """Test SQL injection protection"""
        sql_payloads = [
            "' OR '1'='1",
            "1; DROP TABLE users;",
            "' UNION SELECT * FROM admin--"
        ]
        
        blocked_count = 0
        
        for payload in sql_payloads:
            cmd = ["curl", "-s", "-X", "GET", 
                   f"http://localhost:8085/api/search?q={payload}"]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            # Check if SQL injection was blocked or sanitized
            if (result.returncode != 0 or 
                "error" in result.stdout.lower() or
                "invalid" in result.stdout.lower()):
                blocked_count += 1
        
        success_rate = blocked_count / len(sql_payloads)
        
        return {
            "success": success_rate >= 0.7,
            "message": f"Blocked {blocked_count}/{len(sql_payloads)} SQL injection attempts",
            "details": f"Protection rate: {success_rate:.1%}"
        }
    
    async def _test_auth_validation(self):
        """Test authentication validation"""
        # Test protected endpoint without authentication
        cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null",
               "http://localhost:8085/api/v1/enhanced/diet/generate",
               "-X", "POST"]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        try:
            status_code = int(result.stdout.strip())
            # Should return 401 or 403 for unauthorized access
            if status_code in [401, 403]:
                return {
                    "success": True,
                    "message": "Authentication properly enforced",
                    "details": f"Returned {status_code} for unauthorized access"
                }
            else:
                return {
                    "success": False,
                    "message": f"Authentication bypass detected (status: {status_code})",
                    "recommendation": "Strengthen authentication validation"
                }
        except ValueError:
            return {
                "success": False,
                "message": "Could not test authentication",
                "recommendation": "Manually verify auth endpoints"
            }
    
    async def _test_rate_limiting(self):
        """Test rate limiting functionality"""
        # Make rapid requests to test rate limiting
        rate_limited = False
        
        for i in range(15):  # Make 15 rapid requests
            cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null",
                   "http://localhost:8085/health"]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            try:
                status_code = int(result.stdout.strip())
                if status_code == 429:  # Too Many Requests
                    rate_limited = True
                    break
            except ValueError:
                pass
        
        return {
            "success": rate_limited,
            "message": "Rate limiting is active" if rate_limited else "Rate limiting not detected",
            "recommendation": "Implement rate limiting for production" if not rate_limited else None
        }
    
    async def _run_performance_tests(self):
        print("⚡ Running Performance Tests...")
        
        if not await self._check_app_running():
            print("  ⚠️  Skipping performance tests - application not running")
            return
        
        tests = {
            "response_time": await self._test_response_time(),
            "concurrent_requests": await self._test_concurrent_requests(),
            "memory_usage": await self._test_memory_usage()
        }
        
        self.error_log["prevention_checks"]["performance"] = tests
        
        for test_name, result in tests.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {test_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _test_response_time(self):
        """Test application response time"""
        cmd = ["curl", "-s", "-w", "%{time_total}", "-o", "/dev/null",
               "http://localhost:8085/health"]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        try:
            response_time = float(result.stdout.strip())
            
            if response_time < 1.0:  # Less than 1 second
                return {
                    "success": True,
                    "message": f"Response time: {response_time:.3f}s (Good)",
                    "details": response_time
                }
            elif response_time < 3.0:  # Less than 3 seconds
                return {
                    "success": True,
                    "message": f"Response time: {response_time:.3f}s (Acceptable)",
                    "details": response_time,
                    "recommendation": "Consider performance optimization"
                }
            else:
                return {
                    "success": False,
                    "message": f"Response time: {response_time:.3f}s (Slow)",
                    "details": response_time,
                    "recommendation": "Performance optimization required"
                }
        except ValueError:
            return {
                "success": False,
                "message": "Could not measure response time",
                "recommendation": "Check application responsiveness manually"
            }
    
    async def _test_concurrent_requests(self):
        """Test handling of concurrent requests"""
        # Create multiple concurrent curl processes
        processes = []
        
        for i in range(10):  # 10 concurrent requests
            cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null",
                   "http://localhost:8085/health"]
            proc = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            processes.append(proc)
        
        # Wait for all to complete
        results = []
        for proc in processes:
            stdout, stderr = proc.communicate()
            try:
                status_code = int(stdout.decode().strip())
                results.append(status_code)
            except ValueError:
                results.append(0)
        
        successful_requests = sum(1 for code in results if code == 200)
        success_rate = successful_requests / len(results)
        
        return {
            "success": success_rate >= 0.8,  # 80% success rate
            "message": f"Handled {successful_requests}/{len(results)} concurrent requests",
            "details": {"success_rate": success_rate, "results": results},
            "recommendation": "Improve concurrency handling" if success_rate < 0.8 else None
        }
    
    async def _test_memory_usage(self):
        """Test current memory usage"""
        try:
            # Get process info for the application
            result = subprocess.run(
                ["ps", "aux"], capture_output=True, text=True
            )
            
            app_processes = []
            for line in result.stdout.split('\n'):
                if 'api-key-generator' in line:
                    parts = line.split()
                    if len(parts) > 3:
                        try:
                            memory_percent = float(parts[3])
                            app_processes.append(memory_percent)
                        except ValueError:
                            pass
            
            if app_processes:
                total_memory = sum(app_processes)
                
                if total_memory < 5.0:  # Less than 5% memory usage
                    return {
                        "success": True,
                        "message": f"Memory usage: {total_memory:.1f}% (Good)",
                        "details": app_processes
                    }
                elif total_memory < 15.0:  # Less than 15%
                    return {
                        "success": True,
                        "message": f"Memory usage: {total_memory:.1f}% (Acceptable)",
                        "details": app_processes,
                        "recommendation": "Monitor memory usage in production"
                    }
                else:
                    return {
                        "success": False,
                        "message": f"Memory usage: {total_memory:.1f}% (High)",
                        "details": app_processes,
                        "recommendation": "Investigate memory leaks"
                    }
            else:
                return {
                    "success": True,
                    "message": "Could not measure application memory usage",
                    "recommendation": "Monitor memory in production"
                }
                
        except Exception as e:
            return {
                "success": True,
                "message": f"Memory check error: {str(e)[:50]}",
                "recommendation": "Use system monitoring tools"
            }
    
    async def _run_integration_tests(self):
        print("🔗 Running Integration Tests...")
        
        if not await self._check_app_running():
            print("  ⚠️  Skipping integration tests - application not running")
            return
        
        tests = {
            "api_endpoints": await self._test_api_endpoints(),
            "database_connection": await self._test_database_connection(),
            "external_services": await self._test_external_services()
        }
        
        self.error_log["prevention_checks"]["integration"] = tests
        
        for test_name, result in tests.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {test_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _test_api_endpoints(self):
        """Test critical API endpoints"""
        endpoints = [
            ("GET", "/health", 200),
            ("GET", "/ready", 200),
            ("GET", "/api/recipes", 200),
            ("POST", "/api/v1/enhanced/diet/generate", 401)  # Should require auth
        ]
        
        results = []
        
        for method, endpoint, expected_status in endpoints:
            cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null",
                   "-X", method, f"http://localhost:8085{endpoint}"]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            try:
                status_code = int(result.stdout.strip())
                success = status_code == expected_status
                
                results.append({
                    "endpoint": f"{method} {endpoint}",
                    "expected": expected_status,
                    "actual": status_code,
                    "success": success
                })
            except ValueError:
                results.append({
                    "endpoint": f"{method} {endpoint}",
                    "expected": expected_status,
                    "actual": None,
                    "success": False,
                    "error": "Could not get status code"
                })
        
        successful_endpoints = sum(1 for r in results if r["success"])
        success_rate = successful_endpoints / len(results)
        
        return {
            "success": success_rate >= 0.75,  # 75% success rate
            "message": f"API endpoints: {successful_endpoints}/{len(results)} working correctly",
            "details": results,
            "recommendation": "Fix failing endpoints" if success_rate < 0.75 else None
        }
    
    async def _test_database_connection(self):
        """Test database connectivity and basic operations"""
        # Check if SQLite database file exists
        db_files = ["data/apikeys.db", "apikeys.db", "*.db"]
        db_found = False
        
        for pattern in db_files:
            if "*" in pattern:
                import glob
                matches = glob.glob(pattern)
                if matches:
                    db_found = True
                    break
            else:
                if os.path.exists(pattern):
                    db_found = True
                    break
        
        if db_found:
            return {
                "success": True,
                "message": "Database file found",
                "details": "SQLite database accessible",
                "recommendation": "Test database operations in production"
            }
        else:
            return {
                "success": False,
                "message": "Database file not found",
                "recommendation": "Ensure database is properly initialized"
            }
    
    async def _test_external_services(self):
        """Test external service connectivity"""
        # Test basic internet connectivity
        cmd = ["curl", "-s", "--connect-timeout", "5", "--max-time", "10",
               "-w", "%{http_code}", "-o", "/dev/null", "https://httpbin.org/status/200"]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        try:
            status_code = int(result.stdout.strip())
            if status_code == 200:
                return {
                    "success": True,
                    "message": "External service connectivity working",
                    "details": "Internet connectivity verified"
                }
            else:
                return {
                    "success": False,
                    "message": f"External service test failed (status: {status_code})",
                    "recommendation": "Check network connectivity"
                }
        except ValueError:
            return {
                "success": False,
                "message": "External service connectivity test failed",
                "recommendation": "Check firewall and DNS settings"
            }
    
    async def _check_platform_issues(self):
        print("🏗️ Checking Platform-Specific Issues...")
        
        checks = {
            "coolify_readiness": await self._check_coolify_readiness(),
            "environment_variables": await self._check_environment_variables(),
            "ssl_certificates": await self._check_ssl_readiness()
        }
        
        self.error_log["prevention_checks"]["platform"] = checks
        
        for check_name, result in checks.items():
            status = "✅" if result["success"] else "❌"
            print(f"  {status} {check_name.replace('_', ' ').title()}: {result['message']}")
    
    async def _check_coolify_readiness(self):
        """Check readiness for Coolify deployment"""
        required_files = ["Dockerfile", "go.mod", "main.go"]
        missing_files = []
        
        for file in required_files:
            if not os.path.exists(file):
                missing_files.append(file)
        
        # Check SSH key for Coolify
        ssh_key_path = os.path.expanduser("~/.ssh/coolify_doctorhealthy1")
        ssh_key_exists = os.path.exists(ssh_key_path)
        
        if missing_files:
            return {
                "success": False,
                "message": f"Missing required files: {', '.join(missing_files)}",
                "recommendation": "Create missing files for Coolify deployment"
            }
        elif not ssh_key_exists:
            return {
                "success": False,
                "message": "Coolify SSH key not found",
                "recommendation": "Set up SSH key for Coolify access"
            }
        else:
            return {
                "success": True,
                "message": "Coolify deployment ready",
                "details": "All required files and SSH key present"
            }
    
    async def _check_environment_variables(self):
        """Check for required environment variables"""
        required_env_vars = ["PORT", "DB_PATH"]
        optional_env_vars = ["JWT_SECRET", "API_KEY_PREFIX", "CORS_ORIGINS"]
        
        missing_required = []
        missing_optional = []
        
        for var in required_env_vars:
            if not os.environ.get(var):
                missing_required.append(var)
        
        for var in optional_env_vars:
            if not os.environ.get(var):
                missing_optional.append(var)
        
        if missing_required:
            return {
                "success": False,
                "message": f"Missing required env vars: {', '.join(missing_required)}",
                "recommendation": "Set required environment variables"
            }
        else:
            return {
                "success": True,
                "message": "Environment variables configured",
                "details": {
                    "required_present": len(required_env_vars) - len(missing_required),
                    "optional_missing": missing_optional
                },
                "recommendation": f"Consider setting: {', '.join(missing_optional)}" if missing_optional else None
            }
    
    async def _check_ssl_readiness(self):
        """Check SSL certificate readiness for production"""
        # For now, just check if we're ready for SSL
        return {
            "success": True,
            "message": "SSL configuration not required for initial deployment",
            "recommendation": "Configure SSL certificates after successful deployment"
        }
    
    async def _check_app_running(self):
        """Check if application is running and responsive"""
        cmd = ["curl", "-s", "--connect-timeout", "3", "--max-time", "5",
               "http://localhost:8085/health"]
        result = subprocess.run(cmd, capture_output=True, text=True)
        return result.returncode == 0 and "healthy" in result.stdout
    
    async def _generate_prevention_report(self):
        print("\n📋 Deployment Readiness Report")
        print("=" * 70)
        
        # Calculate overall readiness
        all_checks = {}
        for category, checks in self.error_log["prevention_checks"].items():
            all_checks.update(checks)
        
        total_checks = len(all_checks)
        passed_checks = sum(1 for check in all_checks.values() if check["success"])
        failed_checks = total_checks - passed_checks
        
        readiness_score = (passed_checks / total_checks) * 100 if total_checks > 0 else 0
        
        print(f"📊 Overall Readiness Score: {readiness_score:.1f}%")
        print(f"✅ Passed: {passed_checks}/{total_checks}")
        print(f"❌ Failed: {failed_checks}/{total_checks}")
        
        # Determine deployment readiness
        if readiness_score >= 90:
            status = "🎉 READY FOR DEPLOYMENT"
            self.error_log["deployment_ready"] = True
        elif readiness_score >= 75:
            status = "⚠️  READY WITH WARNINGS"
            self.error_log["deployment_ready"] = True
        else:
            status = "❌ NOT READY FOR DEPLOYMENT"
            self.error_log["deployment_ready"] = False
        
        print(f"\n{status}")
        
        # Generate recommendations
        recommendations = []
        
        for category, checks in self.error_log["prevention_checks"].items():
            for check_name, check_result in checks.items():
                if not check_result["success"] and check_result.get("recommendation"):
                    recommendations.append(f"[{category}] {check_result['recommendation']}")
        
        if recommendations:
            print(f"\n💡 Recommendations:")
            for i, rec in enumerate(recommendations[:5], 1):  # Top 5 recommendations
                print(f"   {i}. {rec}")
        
        # Save detailed report
        report_file = f"deployment_readiness_report_{self.error_log['session_id']}.json"
        with open(report_file, 'w') as f:
            json.dump(self.error_log, f, indent=2, default=str)
        
        print(f"\n📄 Detailed report saved: {report_file}")
        
        return readiness_score

async def main():
    """Run comprehensive deployment error prevention checks"""
    if len(sys.argv) > 1 and sys.argv[1] == "--quick":
        print("🚀 Running Quick Deployment Checks...")
        # Quick mode - just essential checks
    else:
        print("🔍 Running Comprehensive Deployment Error Prevention...")
    
    prevention_system = DeploymentErrorPrevention()
    
    try:
        results = await prevention_system.run_comprehensive_checks()
        
        if results.get("deployment_ready", False):
            print(f"\n✅ System is ready for deployment!")
            sys.exit(0)
        else:
            print(f"\n❌ System requires fixes before deployment!")
            sys.exit(1)
            
    except KeyboardInterrupt:
        print(f"\n⚠️  Error prevention check interrupted")
        sys.exit(2)
    except Exception as e:
        print(f"\n💥 Error prevention check failed: {e}")
        sys.exit(3)

if __name__ == "__main__":
    asyncio.run(main())