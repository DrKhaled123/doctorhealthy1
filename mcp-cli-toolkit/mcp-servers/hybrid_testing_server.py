#!/usr/bin/env python3
"""
Hybrid Testing Framework with MCP
Combines Model Context Protocol with traditional CLI tools for comprehensive testing
"""

import json
import subprocess
import sys
import os
import asyncio
import logging
from typing import Any, Dict, List, Optional

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class HybridTestingServer:
    """MCP server for hybrid testing framework"""
    
    def __init__(self):
        self.name = "hybrid-testing"
        self.tools = self._register_tools()
    
    def _register_tools(self) -> List[Dict[str, Any]]:
        """Register available testing tools"""
        return [
            {
                "name": "run_unit_tests",
                "description": "Execute unit tests using appropriate CLI tool",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "test_type": {
                            "enum": ["go", "python", "java", "js", "ruby"],
                            "description": "Language/framework"
                        },
                        "test_path": {
                            "type": "string", 
                            "description": "Path to test file or directory"
                        },
                        "verbose": {
                            "type": "boolean",
                            "description": "Enable verbose output",
                            "default": False
                        }
                    },
                    "required": ["test_type"]
                }
            },
            {
                "name": "run_integration_tests",
                "description": "Execute integration tests",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["postman", "newman", "curl", "custom"],
                            "description": "Testing tool"
                        },
                        "config": {
                            "type": "string",
                            "description": "Path to configuration file or URL"
                        },
                        "environment": {
                            "type": "string",
                            "description": "Environment variables or config",
                            "default": ""
                        }
                    },
                    "required": ["tool", "config"]
                }
            },
            {
                "name": "run_performance_tests",
                "description": "Execute performance tests",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["ab", "wrk", "hey", "k6"],
                            "description": "Performance tool"
                        },
                        "target_url": {
                            "type": "string",
                            "description": "Target URL for testing"
                        },
                        "concurrent_users": {
                            "type": "integer",
                            "description": "Number of concurrent users",
                            "default": 10
                        },
                        "duration": {
                            "type": "string",
                            "description": "Test duration (e.g., '30s', '2m')",
                            "default": "30s"
                        }
                    },
                    "required": ["tool", "target_url"]
                }
            },
            {
                "name": "run_security_tests",
                "description": "Execute security tests",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["nmap", "curl_security", "custom_scan"],
                            "description": "Security testing tool"
                        },
                        "target": {
                            "type": "string",
                            "description": "Target URL or IP for security testing"
                        },
                        "scan_type": {
                            "enum": ["quick", "full", "custom"],
                            "description": "Type of security scan",
                            "default": "quick"
                        }
                    },
                    "required": ["tool", "target"]
                }
            },
            {
                "name": "validate_api_endpoints",
                "description": "Validate API endpoints and their responses",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "base_url": {
                            "type": "string",
                            "description": "Base URL of the API"
                        },
                        "endpoints": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "path": {"type": "string"},
                                    "method": {"type": "string"},
                                    "expected_status": {"type": "integer"}
                                }
                            },
                            "description": "List of endpoints to validate"
                        }
                    },
                    "required": ["base_url", "endpoints"]
                }
            }
        ]
    
    def list_tools(self) -> List[Dict[str, Any]]:
        """Return list of available tools"""
        return self.tools
    
    async def call_tool(self, name: str, arguments: Dict[str, Any]) -> Dict[str, Any]:
        """Execute the specified tool with given arguments"""
        try:
            if name == "run_unit_tests":
                return await self._run_unit_tests(arguments)
            elif name == "run_integration_tests":
                return await self._run_integration_tests(arguments)
            elif name == "run_performance_tests":
                return await self._run_performance_tests(arguments)
            elif name == "run_security_tests":
                return await self._run_security_tests(arguments)
            elif name == "validate_api_endpoints":
                return await self._validate_api_endpoints(arguments)
            else:
                return {"error": f"Unknown tool: {name}"}
        
        except Exception as e:
            logger.error(f"Error executing tool {name}: {str(e)}")
            return {"error": f"Tool execution failed: {str(e)}"}
    
    async def _run_unit_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute unit tests based on the test type"""
        test_type = args["test_type"]
        test_path = args.get("test_path", ".")
        verbose = args.get("verbose", False)
        
        # Map test types to CLI commands
        commands = {
            "go": ["go", "test", "-v" if verbose else "", "./..."],
            "python": ["python", "-m", "pytest", test_path, "-v" if verbose else ""],
            "java": ["mvn", "test"],
            "js": ["npm", "test"],
            "ruby": ["rspec", test_path]
        }
        
        cmd = commands.get(test_type)
        if not cmd:
            return {"error": f"Unsupported test type: {test_type}"}
        
        # Remove empty strings from command
        cmd = [arg for arg in cmd if arg]
        
        logger.info(f"Running unit tests: {' '.join(cmd)}")
        result = await self._run_command(cmd)
        
        # Parse test results for better reporting
        if test_type == "go" and result.get("exit_code") == 0:
            result["summary"] = "Go tests passed successfully"
        elif test_type == "python" and "pytest" in str(result.get("stdout", "")):
            result["summary"] = self._parse_pytest_output(result.get("stdout", ""))
        
        return result
    
    async def _run_integration_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute integration tests"""
        tool = args["tool"]
        config = args["config"]
        environment = args.get("environment", "")
        
        if tool == "postman" or tool == "newman":
            cmd = ["newman", "run", config]
            if environment:
                cmd.extend(["-e", environment])
        
        elif tool == "curl":
            # For simple curl-based integration tests
            cmd = ["curl", "-s", "-f", config]
        
        elif tool == "custom":
            # Custom integration test script
            if os.path.exists(config):
                cmd = ["bash", config]
            else:
                return {"error": f"Custom test script not found: {config}"}
        else:
            return {"error": f"Unsupported integration tool: {tool}"}
        
        logger.info(f"Running integration tests: {' '.join(cmd)}")
        result = await self._run_command(cmd)
        
        # Enhance result with test summary
        if result.get("exit_code") == 0:
            result["summary"] = f"Integration tests passed using {tool}"
        else:
            result["summary"] = f"Integration tests failed using {tool}"
        
        return result
    
    async def _run_performance_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute performance tests"""
        tool = args["tool"]
        target_url = args["target_url"]
        concurrent_users = args.get("concurrent_users", 10)
        duration = args.get("duration", "30s")
        
        if tool == "ab":
            # Apache Bench
            requests = concurrent_users * 100  # Total requests
            cmd = ["ab", "-n", str(requests), "-c", str(concurrent_users), target_url]
        
        elif tool == "wrk":
            cmd = ["wrk", "-t", str(concurrent_users), "-c", str(concurrent_users), 
                   "-d", duration, target_url]
        
        elif tool == "hey":
            cmd = ["hey", "-z", duration, "-c", str(concurrent_users), target_url]
        
        elif tool == "k6":
            # Create a simple k6 script
            k6_script = f"""
import http from 'k6/http';
import {{ check }} from 'k6';

export let options = {{
  vus: {concurrent_users},
  duration: '{duration}',
}};

export default function() {{
  let response = http.get('{target_url}');
  check(response, {{
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  }});
}}
"""
            script_path = "/tmp/k6_test.js"
            with open(script_path, "w") as f:
                f.write(k6_script)
            cmd = ["k6", "run", script_path]
        else:
            return {"error": f"Unsupported performance tool: {tool}"}
        
        logger.info(f"Running performance tests: {' '.join(cmd)}")
        result = await self._run_command(cmd)
        
        # Parse performance metrics
        result["performance_summary"] = self._parse_performance_output(tool, result.get("stdout", ""))
        
        return result
    
    async def _run_security_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute security tests"""
        tool = args["tool"]
        target = args["target"]
        scan_type = args.get("scan_type", "quick")
        
        if tool == "nmap":
            if scan_type == "quick":
                cmd = ["nmap", "-sS", "-O", target]
            elif scan_type == "full":
                cmd = ["nmap", "-A", "-T4", target]
            else:
                cmd = ["nmap", "-sS", target]
        
        elif tool == "curl_security":
            # Basic security headers check
            cmd = ["curl", "-I", "-s", target]
        
        elif tool == "custom_scan":
            # Custom security scan - checking common vulnerabilities
            return await self._custom_security_scan(target)
        
        else:
            return {"error": f"Unsupported security tool: {tool}"}
        
        logger.info(f"Running security tests: {' '.join(cmd)}")
        result = await self._run_command(cmd)
        
        # Analyze security results
        if tool == "curl_security":
            result["security_analysis"] = self._analyze_security_headers(result.get("stdout", ""))
        
        return result
    
    async def _validate_api_endpoints(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Validate multiple API endpoints"""
        base_url = args["base_url"]
        endpoints = args["endpoints"]
        
        results = []
        overall_success = True
        
        for endpoint in endpoints:
            path = endpoint["path"]
            method = endpoint.get("method", "GET")
            expected_status = endpoint.get("expected_status", 200)
            
            url = f"{base_url.rstrip('/')}/{path.lstrip('/')}"
            
            if method == "GET":
                cmd = ["curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", url]
            else:
                cmd = ["curl", "-s", "-X", method, "-o", "/dev/null", "-w", "%{http_code}", url]
            
            result = await self._run_command(cmd)
            
            try:
                actual_status = int(result.get("stdout", "0"))
                success = actual_status == expected_status
                
                if not success:
                    overall_success = False
                
                results.append({
                    "endpoint": path,
                    "method": method,
                    "expected_status": expected_status,
                    "actual_status": actual_status,
                    "success": success
                })
            
            except ValueError:
                overall_success = False
                results.append({
                    "endpoint": path,
                    "method": method,
                    "expected_status": expected_status,
                    "actual_status": "error",
                    "success": False,
                    "error": "Failed to get response"
                })
        
        return {
            "overall_success": overall_success,
            "total_endpoints": len(endpoints),
            "successful_endpoints": sum(1 for r in results if r["success"]),
            "results": results
        }
    
    async def _run_command(self, cmd: List[str]) -> Dict[str, Any]:
        """Run a command asynchronously and return result"""
        try:
            process = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await process.communicate()
            
            return {
                "exit_code": process.returncode,
                "stdout": stdout.decode('utf-8', errors='replace'),
                "stderr": stderr.decode('utf-8', errors='replace'),
                "command": " ".join(cmd)
            }
        
        except Exception as e:
            return {
                "exit_code": -1,
                "stdout": "",
                "stderr": str(e),
                "command": " ".join(cmd),
                "error": f"Command execution failed: {str(e)}"
            }
    
    def _parse_pytest_output(self, output: str) -> str:
        """Parse pytest output for summary"""
        lines = output.split('\n')
        for line in reversed(lines):
            if "passed" in line or "failed" in line:
                return line.strip()
        return "Pytest execution completed"
    
    def _parse_performance_output(self, tool: str, output: str) -> Dict[str, Any]:
        """Parse performance test output"""
        summary = {"tool": tool}
        
        if tool == "ab" and "Requests per second" in output:
            for line in output.split('\n'):
                if "Requests per second" in line:
                    summary["rps"] = line.split(':')[1].strip()
                elif "Time per request" in line and "mean" in line:
                    summary["avg_response_time"] = line.split(':')[1].strip()
        
        elif tool == "wrk":
            lines = output.split('\n')
            for line in lines:
                if "Requests/sec" in line:
                    summary["rps"] = line.split(':')[1].strip()
                elif "Latency" in line:
                    summary["avg_latency"] = line.split()[1]
        
        return summary
    
    def _analyze_security_headers(self, headers: str) -> Dict[str, Any]:
        """Analyze security headers from curl response"""
        analysis = {
            "security_headers_present": [],
            "missing_security_headers": [],
            "recommendations": []
        }
        
        security_headers = [
            "X-Content-Type-Options",
            "X-Frame-Options",
            "X-XSS-Protection",
            "Strict-Transport-Security",
            "Content-Security-Policy"
        ]
        
        headers_lower = headers.lower()
        
        for header in security_headers:
            if header.lower() in headers_lower:
                analysis["security_headers_present"].append(header)
            else:
                analysis["missing_security_headers"].append(header)
                analysis["recommendations"].append(f"Add {header} header for better security")
        
        return analysis
    
    async def _custom_security_scan(self, target: str) -> Dict[str, Any]:
        """Perform custom security scan"""
        results = {
            "target": target,
            "vulnerabilities": [],
            "security_score": 100
        }
        
        # Test for common security issues
        tests = [
            {"path": "/admin", "description": "Admin panel exposure"},
            {"path": "/.env", "description": "Environment file exposure"},
            {"path": "/config", "description": "Configuration exposure"},
            {"path": "/api/v1/health", "description": "Health endpoint availability"}
        ]
        
        for test in tests:
            url = f"{target.rstrip('/')}{test['path']}"
            cmd = ["curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", url]
            result = await self._run_command(cmd)
            
            try:
                status_code = int(result.get("stdout", "0"))
                if status_code == 200 and test["path"] in ["/admin", "/.env", "/config"]:
                    results["vulnerabilities"].append({
                        "type": "Information Disclosure",
                        "description": test["description"],
                        "url": url,
                        "severity": "Medium"
                    })
                    results["security_score"] -= 20
            except ValueError:
                pass
        
        return results

def main():
    """Main function to run the MCP server"""
    server = HybridTestingServer()
    
    # Simple stdio server implementation
    while True:
        try:
            line = sys.stdin.readline()
            if not line:
                break
            
            request = json.loads(line.strip())
            
            if request.get("method") == "tools/list":
                response = {
                    "jsonrpc": "2.0",
                    "id": request.get("id"),
                    "result": {"tools": server.list_tools()}
                }
            
            elif request.get("method") == "tools/call":
                params = request.get("params", {})
                tool_name = params.get("name")
                arguments = params.get("arguments", {})
                
                # Run the tool
                result = asyncio.run(server.call_tool(tool_name, arguments))
                
                response = {
                    "jsonrpc": "2.0",
                    "id": request.get("id"),
                    "result": {"content": [{"type": "text", "text": json.dumps(result, indent=2)}]}
                }
            
            else:
                response = {
                    "jsonrpc": "2.0",
                    "id": request.get("id"),
                    "error": {"code": -32601, "message": "Method not found"}
                }
            
            print(json.dumps(response))
            sys.stdout.flush()
        
        except Exception as e:
            logger.error(f"Error processing request: {str(e)}")
            error_response = {
                "jsonrpc": "2.0",
                "id": request.get("id", None),
                "error": {"code": -32603, "message": f"Internal error: {str(e)}"}
            }
            print(json.dumps(error_response))
            sys.stdout.flush()

if __name__ == "__main__":
    main()