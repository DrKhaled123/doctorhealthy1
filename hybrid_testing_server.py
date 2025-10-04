#!/usr/bin/env python3
"""
Hybrid Testing Framework with MCP
Combines Model Context Protocol with traditional CLI tools for optimal testing
"""

import asyncio
import subprocess
import json
import sys
import os
from pathlib import Path

# MCP Server implementation
class HybridTestingServer:
    def __init__(self):
        self.name = "hybrid-testing"
        self.version = "1.0.0"
        
    async def list_tools(self):
        return [
            {
                "name": "run_unit_tests",
                "description": "Execute unit tests using appropriate CLI tool",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "test_type": {
                            "enum": ["go", "python", "js", "java"], 
                            "description": "Language/framework"
                        },
                        "test_path": {
                            "type": "string", 
                            "description": "Path to test file or directory"
                        }
                    }
                }
            },
            {
                "name": "run_integration_tests",
                "description": "Execute integration tests",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["postman", "curl", "go-test"], 
                            "description": "Testing tool"
                        },
                        "config": {
                            "type": "string", 
                            "description": "Path to configuration file or endpoint URL"
                        }
                    }
                }
            },
            {
                "name": "run_security_tests",
                "description": "Execute security tests",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "test_type": {
                            "enum": ["sanitization", "rate_limit", "auth"], 
                            "description": "Security test type"
                        },
                        "endpoint": {
                            "type": "string", 
                            "description": "API endpoint to test"
                        }
                    }
                }
            },
            {
                "name": "validate_build",
                "description": "Validate build and compilation",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "build_type": {
                            "enum": ["go", "docker"], 
                            "description": "Build system"
                        },
                        "target": {
                            "type": "string", 
                            "description": "Build target"
                        }
                    }
                }
            }
        ]
    
    async def call_tool(self, name, arguments):
        try:
            if name == "run_unit_tests":
                return await self._run_unit_tests(arguments)
            elif name == "run_integration_tests":
                return await self._run_integration_tests(arguments)
            elif name == "run_security_tests":
                return await self._run_security_tests(arguments)
            elif name == "validate_build":
                return await self._validate_build(arguments)
            else:
                return {"error": f"Unknown tool: {name}"}
        except Exception as e:
            return {"error": str(e)}
    
    async def _run_unit_tests(self, args):
        test_type = args.get("test_type", "go")
        test_path = args.get("test_path", "./...")
        
        commands = {
            "go": ["go", "test", "-v", test_path],
            "python": ["pytest", "-v", test_path],
            "js": ["npm", "test"],
            "java": ["mvn", "test"]
        }
        
        cmd = commands.get(test_type)
        if not cmd:
            return {"error": f"Unsupported test type: {test_type}"}
        
        result = subprocess.run(cmd, capture_output=True, text=True, cwd=".")
        return {
            "success": result.returncode == 0,
            "stdout": result.stdout,
            "stderr": result.stderr,
            "exit_code": result.returncode,
            "command": " ".join(cmd)
        }
    
    async def _run_integration_tests(self, args):
        tool = args.get("tool", "curl")
        config = args.get("config", "")
        
        if tool == "go-test":
            cmd = ["go", "test", "-v", "-tags=integration", "./..."]
        elif tool == "curl":
            if not config.startswith("http"):
                return {"error": "Config must be a valid URL for curl tests"}
            cmd = ["curl", "-f", "-s", "-S", config]
        elif tool == "postman":
            if not config.endswith(".json"):
                return {"error": "Config must be a Postman collection JSON file"}
            cmd = ["newman", "run", config]
        else:
            return {"error": f"Unsupported integration tool: {tool}"}
        
        result = subprocess.run(cmd, capture_output=True, text=True, cwd=".")
        return {
            "success": result.returncode == 0,
            "stdout": result.stdout,
            "stderr": result.stderr,
            "exit_code": result.returncode,
            "command": " ".join(cmd)
        }
    
    async def _run_security_tests(self, args):
        test_type = args.get("test_type", "sanitization")
        endpoint = args.get("endpoint", "http://localhost:8080")
        
        results = []
        
        if test_type == "sanitization":
            # Test XSS protection
            xss_payloads = [
                "<script>alert('xss')</script>",
                "javascript:alert('xss')",
                "<img src=x onerror=alert('xss')>",
                "' OR '1'='1",
                "../../../etc/passwd"
            ]
            
            for payload in xss_payloads:
                cmd = ["curl", "-s", "-X", "POST", 
                       f"{endpoint}/api/test", 
                       "-H", "Content-Type: application/json",
                       "-d", json.dumps({"input": payload})]
                
                result = subprocess.run(cmd, capture_output=True, text=True)
                results.append({
                    "payload": payload,
                    "blocked": "error" in result.stdout.lower() or result.returncode != 0,
                    "response": result.stdout[:200]  # Truncate response
                })
        
        elif test_type == "rate_limit":
            # Test rate limiting
            cmd = ["bash", "-c", f"for i in {{1..20}}; do curl -s {endpoint}/health > /dev/null; done"]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            # Test if rate limiting kicks in
            rapid_cmd = ["curl", "-s", endpoint + "/health"]
            for i in range(10):
                rapid_result = subprocess.run(rapid_cmd, capture_output=True, text=True)
                if "429" in rapid_result.stdout or rapid_result.returncode != 0:
                    results.append({"rate_limit_active": True, "requests_before_limit": i})
                    break
            else:
                results.append({"rate_limit_active": False})
        
        elif test_type == "auth":
            # Test authentication
            cmd = ["curl", "-s", "-X", "POST", 
                   f"{endpoint}/api/v1/enhanced/diet/generate",
                   "-H", "Content-Type: application/json",
                   "-d", "{}"]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            auth_required = "unauthorized" in result.stdout.lower() or "api key" in result.stdout.lower()
            results.append({
                "auth_required": auth_required,
                "response": result.stdout[:200]
            })
        
        return {
            "success": True,
            "test_type": test_type,
            "results": results
        }
    
    async def _validate_build(self, args):
        build_type = args.get("build_type", "go")
        target = args.get("target", ".")
        
        if build_type == "go":
            # Check Go build
            cmd = ["go", "build", "-o", "/tmp/test_build", target]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                # Clean up test binary
                subprocess.run(["rm", "-f", "/tmp/test_build"], capture_output=True)
            
            return {
                "success": result.returncode == 0,
                "stdout": result.stdout,
                "stderr": result.stderr,
                "exit_code": result.returncode,
                "command": " ".join(cmd)
            }
        
        elif build_type == "docker":
            # Check Docker build
            cmd = ["docker", "build", "-t", "test-image", target]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                # Clean up test image
                subprocess.run(["docker", "rmi", "test-image"], capture_output=True)
            
            return {
                "success": result.returncode == 0,
                "stdout": result.stdout,
                "stderr": result.stderr,
                "exit_code": result.returncode,
                "command": " ".join(cmd)
            }
        
        return {"error": f"Unsupported build type: {build_type}"}

# CLI interface for direct usage
async def main():
    if len(sys.argv) < 2:
        print("Usage: python hybrid_testing_server.py <command> [args...]")
        print("Commands: unit-tests, integration-tests, security-tests, validate-build")
        return
    
    server = HybridTestingServer()
    command = sys.argv[1]
    
    if command == "unit-tests":
        test_type = sys.argv[2] if len(sys.argv) > 2 else "go"
        test_path = sys.argv[3] if len(sys.argv) > 3 else "./..."
        result = await server.call_tool("run_unit_tests", {"test_type": test_type, "test_path": test_path})
    
    elif command == "integration-tests":
        tool = sys.argv[2] if len(sys.argv) > 2 else "curl"
        config = sys.argv[3] if len(sys.argv) > 3 else "http://localhost:8080/health"
        result = await server.call_tool("run_integration_tests", {"tool": tool, "config": config})
    
    elif command == "security-tests":
        test_type = sys.argv[2] if len(sys.argv) > 2 else "sanitization"
        endpoint = sys.argv[3] if len(sys.argv) > 3 else "http://localhost:8080"
        result = await server.call_tool("run_security_tests", {"test_type": test_type, "endpoint": endpoint})
    
    elif command == "validate-build":
        build_type = sys.argv[2] if len(sys.argv) > 2 else "go"
        target = sys.argv[3] if len(sys.argv) > 3 else "."
        result = await server.call_tool("validate_build", {"build_type": build_type, "target": target})
    
    else:
        print(f"Unknown command: {command}")
        return
    
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    asyncio.run(main())