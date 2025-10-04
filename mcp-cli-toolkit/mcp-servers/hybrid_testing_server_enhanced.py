#!/usr/bin/env python3
"""
Enhanced Hybrid Testing Framework with MCP
Combines Model Context Protocol with traditional CLI tools for comprehensive testing
Enhanced with safety features, error handling, and production-ready capabilities
"""

import json
import subprocess
import sys
import os
import asyncio
import logging
import time
import tempfile
from typing import Any, Dict, List, Optional
from pathlib import Path
import traceback

# Setup enhanced logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/tmp/hybrid_testing_server.log'),
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)

class SafeSubprocessRunner:
    """Safe subprocess runner with timeout and error handling"""

    def __init__(self, timeout: int = 300):
        self.timeout = timeout
        self.max_output_size = 10 * 1024 * 1024  # 10MB limit

    def run_command(self, cmd: List[str], cwd: str = None) -> Dict[str, Any]:
        """Run command safely with comprehensive error handling"""
        try:
            logger.info(f"Running command: {' '.join(cmd)}")

            # Validate command
            if not cmd or not isinstance(cmd, list):
                return {
                    "exit_code": -1,
                    "stdout": "",
                    "stderr": "Invalid command format",
                    "error": "Command must be a non-empty list"
                }

            # Security check - prevent dangerous commands
            dangerous_commands = ['rm', 'del', 'format', 'fdisk', 'mkfs', 'dd']
            if any(dangerous in cmd[0] for dangerous in dangerous_commands):
                return {
                    "exit_code": -1,
                    "stdout": "",
                    "stderr": "Dangerous command blocked for safety",
                    "error": "Command blocked for security reasons"
                }

            process = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                cwd=cwd,
                text=True,
                bufsize=8192
            )

            try:
                stdout, stderr = process.communicate(timeout=self.timeout)

                # Truncate output if too large
                if len(stdout) > self.max_output_size:
                    stdout = stdout[:self.max_output_size] + "\n... (output truncated)"
                if len(stderr) > self.max_output_size:
                    stderr = stderr[:self.max_output_size] + "\n... (error truncated)"

                return {
                    "exit_code": process.returncode,
                    "stdout": stdout,
                    "stderr": stderr,
                    "command": " ".join(cmd),
                    "execution_time": time.time()
                }

            except subprocess.TimeoutExpired:
                process.kill()
                return {
                    "exit_code": -1,
                    "stdout": "",
                    "stderr": f"Command timed out after {self.timeout} seconds",
                    "error": "Command execution timeout",
                    "command": " ".join(cmd)
                }

        except FileNotFoundError:
            return {
                "exit_code": -1,
                "stdout": "",
                "stderr": f"Command not found: {cmd[0] if cmd else 'unknown'}",
                "error": "Command not found in PATH"
            }
        except Exception as e:
            logger.error(f"Error running command: {str(e)}")
            return {
                "exit_code": -1,
                "stdout": "",
                "stderr": f"Command execution failed: {str(e)}",
                "error": f"Unexpected error: {str(e)}",
                "command": " ".join(cmd) if cmd else "unknown"
            }

class EnhancedHybridTestingServer:
    """Enhanced MCP server for hybrid testing framework with safety features"""

    def __init__(self):
        self.name = "enhanced-hybrid-testing"
        self.tools = self._register_tools()
        self.runner = SafeSubprocessRunner()
        self.test_results = {}
        self.supported_languages = {
            "java": {"test_cmd": ["mvn", "test"], "build_tool": "maven"},
            "python": {"test_cmd": ["python", "-m", "pytest"], "build_tool": "pytest"},
            "javascript": {"test_cmd": ["npm", "test"], "build_tool": "npm"},
            "typescript": {"test_cmd": ["npm", "test"], "build_tool": "npm"},
            "go": {"test_cmd": ["go", "test"], "build_tool": "go"},
            "ruby": {"test_cmd": ["rspec"], "build_tool": "rspec"},
            "csharp": {"test_cmd": ["dotnet", "test"], "build_tool": "dotnet"},
            "php": {"test_cmd": ["phpunit"], "build_tool": "phpunit"}
        }

    def _register_tools(self) -> List[Dict[str, Any]]:
        """Register available testing tools with enhanced schemas"""
        return [
            {
                "name": "run_unit_tests",
                "description": "Execute unit tests using appropriate CLI tool with safety checks",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "test_type": {
                            "enum": list(self.supported_languages.keys()),
                            "description": "Programming language/framework"
                        },
                        "test_path": {
                            "type": "string",
                            "description": "Path to test file or directory",
                            "default": "."
                        },
                        "verbose": {
                            "type": "boolean",
                            "description": "Enable verbose output",
                            "default": False
                        },
                        "coverage": {
                            "type": "boolean",
                            "description": "Generate coverage report",
                            "default": False
                        },
                        "parallel": {
                            "type": "boolean",
                            "description": "Run tests in parallel",
                            "default": True
                        }
                    },
                    "required": ["test_type"]
                }
            },
            {
                "name": "run_integration_tests",
                "description": "Execute integration tests with enhanced tool support",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["postman", "newman", "rest-assured", "soapui", "cypress", "playwright"],
                            "description": "Testing tool to use"
                        },
                        "config": {
                            "type": "string",
                            "description": "Path to configuration file or URL"
                        },
                        "environment": {
                            "type": "string",
                            "description": "Environment configuration",
                            "default": ""
                        },
                        "timeout": {
                            "type": "integer",
                            "description": "Test timeout in seconds",
                            "default": 300
                        }
                    },
                    "required": ["tool", "config"]
                }
            },
            {
                "name": "run_performance_tests",
                "description": "Execute performance tests with multiple tools",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["jmeter", "k6", "artillery", "locust", "ab", "wrk"],
                            "description": "Performance testing tool"
                        },
                        "target_url": {
                            "type": "string",
                            "description": "Target URL for testing"
                        },
                        "script": {
                            "type": "string",
                            "description": "Path to test script (for script-based tools)"
                        },
                        "duration": {
                            "type": "string",
                            "description": "Test duration (e.g., '30s', '2m')",
                            "default": "30s"
                        },
                        "virtual_users": {
                            "type": "integer",
                            "description": "Number of virtual users",
                            "default": 10
                        }
                    },
                    "required": ["tool"]
                }
            },
            {
                "name": "run_security_tests",
                "description": "Execute security tests with safety controls",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "tool": {
                            "enum": ["owasp-zap", "sqlmap", "nikto", "nuclei", "semgrep"],
                            "description": "Security testing tool"
                        },
                        "target": {
                            "type": "string",
                            "description": "Target URL or file for security testing"
                        },
                        "scan_type": {
                            "enum": ["quick", "full", "custom"],
                            "description": "Type of security scan",
                            "default": "quick"
                        },
                        "severity": {
                            "enum": ["low", "medium", "high", "critical"],
                            "description": "Minimum severity level",
                            "default": "medium"
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
                                    "method": {"type": "string", "default": "GET"},
                                    "expected_status": {"type": "integer", "default": 200},
                                    "headers": {"type": "object"},
                                    "body": {"type": "object"}
                                }
                            },
                            "description": "List of endpoints to validate"
                        },
                        "timeout": {
                            "type": "integer",
                            "description": "Request timeout in seconds",
                            "default": 30
                        }
                    },
                    "required": ["base_url", "endpoints"]
                }
            },
            {
                "name": "get_test_report",
                "description": "Get comprehensive test report and analysis",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "test_run_id": {
                            "type": "string",
                            "description": "Test run identifier"
                        },
                        "format": {
                            "enum": ["json", "html", "xml", "summary"],
                            "description": "Report format",
                            "default": "summary"
                        },
                        "include_coverage": {
                            "type": "boolean",
                            "description": "Include coverage data",
                            "default": True
                        }
                    },
                    "required": ["test_run_id"]
                }
            }
        ]

    def list_tools(self) -> List[Dict[str, Any]]:
        """Return list of available tools"""
        return self.tools

    async def call_tool(self, name: str, arguments: Dict[str, Any]) -> Dict[str, Any]:
        """Execute the specified tool with given arguments and comprehensive error handling"""
        try:
            logger.info(f"Calling tool: {name} with args: {arguments}")

            # Validate arguments
            if not isinstance(arguments, dict):
                return {"error": "Arguments must be a dictionary"}

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
            elif name == "get_test_report":
                return await self._get_test_report(arguments)
            else:
                return {"error": f"Unknown tool: {name}"}

        except Exception as e:
            logger.error(f"Error executing tool {name}: {str(e)}")
            logger.error(traceback.format_exc())
            return {
                "error": f"Tool execution failed: {str(e)}",
                "traceback": traceback.format_exc()
            }

    async def _run_unit_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute unit tests with enhanced safety and reporting"""
        test_type = args["test_type"]
        test_path = args.get("test_path", ".")
        verbose = args.get("verbose", False)
        coverage = args.get("coverage", False)
        parallel = args.get("parallel", True)

        if test_type not in self.supported_languages:
            return {"error": f"Unsupported test type: {test_type}"}

        # Generate test run ID
        test_run_id = f"unit_{test_type}_{int(time.time())}"

        # Build command based on language
        base_cmd = self.supported_languages[test_type]["test_cmd"].copy()

        # Add language-specific options
        if test_type == "python":
            if verbose:
                base_cmd.append("-v")
            if coverage:
                base_cmd.extend(["--cov", test_path, "--cov-report", "html"])
            if parallel:
                base_cmd.append("-n")
                base_cmd.append("auto")

        elif test_type == "java":
            if verbose:
                base_cmd.append("-X")
            if coverage:
                base_cmd.extend(["-Dcoverage.enabled=true"])

        elif test_type == "javascript" or test_type == "typescript":
            if coverage:
                base_cmd.extend(["--coverage"])

        # Add test path
        if test_path != ".":
            base_cmd.append(test_path)

        logger.info(f"Running unit tests: {' '.join(base_cmd)}")

        # Execute tests
        result = self.runner.run_command(base_cmd)

        # Parse results
        test_summary = self._parse_test_results(test_type, result.get("stdout", ""))

        # Store results
        self.test_results[test_run_id] = {
            "type": "unit",
            "language": test_type,
            "result": result,
            "summary": test_summary,
            "timestamp": time.time()
        }

        # Enhanced result
        enhanced_result = {
            "test_run_id": test_run_id,
            "test_type": "unit",
            "language": test_type,
            "success": result.get("exit_code") == 0,
            "summary": test_summary,
            "raw_output": result,
            "execution_time": result.get("execution_time")
        }

        return enhanced_result

    async def _run_integration_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute integration tests with enhanced tool support"""
        tool = args["tool"]
        config = args["config"]
        environment = args.get("environment", "")
        timeout = args.get("timeout", 300)

        # Update runner timeout
        old_timeout = self.runner.timeout
        self.runner.timeout = timeout

        try:
            test_run_id = f"integration_{tool}_{int(time.time())}"

            if tool == "postman" or tool == "newman":
                cmd = ["newman", "run", config]
                if environment:
                    cmd.extend(["-e", environment])

            elif tool == "rest-assured":
                cmd = ["mvn", "verify", "-Dtest=" + config]

            elif tool == "soapui":
                cmd = ["testrunner.sh", "-s", config]

            elif tool == "cypress":
                cmd = ["npx", "cypress", "run", "--spec", config]

            elif tool == "playwright":
                cmd = ["npx", "playwright", "test", config]

            else:
                return {"error": f"Unsupported integration tool: {tool}"}

            logger.info(f"Running integration tests: {' '.join(cmd)}")
            result = self.runner.run_command(cmd)

            # Parse integration test results
            summary = self._parse_integration_results(tool, result.get("stdout", ""))

            # Store results
            self.test_results[test_run_id] = {
                "type": "integration",
                "tool": tool,
                "result": result,
                "summary": summary,
                "timestamp": time.time()
            }

            return {
                "test_run_id": test_run_id,
                "test_type": "integration",
                "tool": tool,
                "success": result.get("exit_code") == 0,
                "summary": summary,
                "raw_output": result
            }

        finally:
            self.runner.timeout = old_timeout

    async def _run_performance_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute performance tests with multiple tools"""
        tool = args["tool"]
        target_url = args.get("target_url")
        script = args.get("script")
        duration = args.get("duration", "30s")
        virtual_users = args.get("virtual_users", 10)

        test_run_id = f"performance_{tool}_{int(time.time())}"

        if tool == "jmeter":
            if not script:
                return {"error": "Script path required for JMeter"}
            cmd = ["jmeter", "-n", "-t", script, "-Jusers=" + str(virtual_users)]

        elif tool == "k6":
            if script:
                cmd = ["k6", "run", "--vus", str(virtual_users), "--duration", duration, script]
            elif target_url:
                # Generate simple k6 script
                k6_script = f"""
import http from 'k6/http';
import {{ check }} from 'k6';

export let options = {{
  vus: {virtual_users},
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
                with tempfile.NamedTemporaryFile(mode='w', suffix='.js', delete=False) as f:
                    f.write(k6_script)
                    script_path = f.name

                cmd = ["k6", "run", "--vus", str(virtual_users), "--duration", duration, script_path]
            else:
                return {"error": "Either script or target_url required for K6"}

        elif tool == "artillery":
            if script:
                cmd = ["artillery", "run", "--config", script]
            else:
                return {"error": "Script required for Artillery"}

        elif tool == "locust":
            if script:
                cmd = ["locust", "-f", script, "--headless", "-u", str(virtual_users), "-t", duration]
            else:
                return {"error": "Script required for Locust"}

        elif tool == "ab":  # Apache Bench
            if not target_url:
                return {"error": "Target URL required for Apache Bench"}
            requests = virtual_users * 100
            cmd = ["ab", "-n", str(requests), "-c", str(virtual_users), target_url]

        elif tool == "wrk":
            if not target_url:
                return {"error": "Target URL required for WRK"}
            cmd = ["wrk", "-t", str(virtual_users), "-c", str(virtual_users), "-d", duration, target_url]

        else:
            return {"error": f"Unsupported performance tool: {tool}"}

        logger.info(f"Running performance tests: {' '.join(cmd)}")
        result = self.runner.run_command(cmd)

        # Parse performance results
        summary = self._parse_performance_results(tool, result.get("stdout", ""))

        # Store results
        self.test_results[test_run_id] = {
            "type": "performance",
            "tool": tool,
            "result": result,
            "summary": summary,
            "timestamp": time.time()
        }

        return {
            "test_run_id": test_run_id,
            "test_type": "performance",
            "tool": tool,
            "success": result.get("exit_code") == 0,
            "summary": summary,
            "raw_output": result
        }

    async def _run_security_tests(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Execute security tests with safety controls"""
        tool = args["tool"]
        target = args["target"]
        scan_type = args.get("scan_type", "quick")
        severity = args.get("severity", "medium")

        test_run_id = f"security_{tool}_{int(time.time())}"

        # Security tools with safe configurations
        if tool == "owasp-zap":
            if scan_type == "quick":
                cmd = ["zap.sh", "-cmd", "-autorun", "/zap/wrk/quick.yaml"]
            else:
                cmd = ["zap.sh", "-cmd", "-autorun", f"/zap/wrk/{scan_type}.yaml"]

        elif tool == "sqlmap":
            cmd = ["sqlmap", "-u", target, "--batch", "--risk=1", "--level=1"]

        elif tool == "nikto":
            cmd = ["nikto", "-h", target, "-timeout", "30"]

        elif tool == "nuclei":
            cmd = ["nuclei", "-u", target, "-severity", severity]

        elif tool == "semgrep":
            cmd = ["semgrep", "--config=auto", target]

        else:
            return {"error": f"Unsupported security tool: {tool}"}

        logger.info(f"Running security tests: {' '.join(cmd)}")
        result = self.runner.run_command(cmd)

        # Parse security results
        summary = self._parse_security_results(tool, result.get("stdout", ""))

        # Store results
        self.test_results[test_run_id] = {
            "type": "security",
            "tool": tool,
            "result": result,
            "summary": summary,
            "timestamp": time.time()
        }

        return {
            "test_run_id": test_run_id,
            "test_type": "security",
            "tool": tool,
            "success": result.get("exit_code") == 0,
            "summary": summary,
            "raw_output": result
        }

    async def _validate_api_endpoints(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Validate multiple API endpoints with comprehensive reporting"""
        base_url = args["base_url"]
        endpoints = args["endpoints"]
        timeout = args.get("timeout", 30)

        test_run_id = f"api_validation_{int(time.time())}"
        results = []
        overall_success = True

        for endpoint in endpoints:
            path = endpoint["path"]
            method = endpoint.get("method", "GET")
            expected_status = endpoint.get("expected_status", 200)
            headers = endpoint.get("headers", {})
            body = endpoint.get("body", {})

            url = f"{base_url.rstrip('/')}/{path.lstrip('/')}"

            # Build curl command
            cmd = ["curl", "-s", "-w",
                   "response_time:%{time_total}\\nhttp_code:%{http_code}\\n",
                   "--max-time", str(timeout),
                   "-X", method]

            # Add headers
            for key, value in headers.items():
                cmd.extend(["-H", f"{key}: {value}"])

            # Add body for POST/PUT/PATCH
            if body and method.upper() in ["POST", "PUT", "PATCH"]:
                cmd.extend(["-d", json.dumps(body)])

            cmd.append(url)

            result = self.runner.run_command(cmd)

            # Parse response
            try:
                output = result.get("stdout", "")
                response_time = None
                actual_status = None

                if "response_time:" in output:
                    response_time = float(output.split("response_time:")[1].split()[0]) * 1000  # Convert to ms

                if "http_code:" in output:
                    actual_status = int(output.split("http_code:")[1].split()[0])

                success = actual_status == expected_status

                if not success:
                    overall_success = False

                endpoint_result = {
                    "endpoint": path,
                    "method": method,
                    "url": url,
                    "expected_status": expected_status,
                    "actual_status": actual_status,
                    "response_time_ms": response_time,
                    "success": success,
                    "raw_output": result
                }

            except Exception as e:
                overall_success = False
                endpoint_result = {
                    "endpoint": path,
                    "method": method,
                    "url": url,
                    "expected_status": expected_status,
                    "actual_status": "error",
                    "response_time_ms": None,
                    "success": False,
                    "error": str(e),
                    "raw_output": result
                }

            results.append(endpoint_result)

        # Calculate summary
        successful_endpoints = sum(1 for r in results if r["success"])
        avg_response_time = sum(r["response_time_ms"] for r in results if r["response_time_ms"]) / len([r for r in results if r["response_time_ms"]]) if any(r["response_time_ms"] for r in results) else 0

        summary = {
            "total_endpoints": len(endpoints),
            "successful_endpoints": successful_endpoints,
            "failed_endpoints": len(endpoints) - successful_endpoints,
            "success_rate": round((successful_endpoints / len(endpoints)) * 100, 2) if endpoints else 0,
            "average_response_time_ms": round(avg_response_time, 2) if avg_response_time else 0,
            "overall_success": overall_success
        }

        # Store results
        self.test_results[test_run_id] = {
            "type": "api_validation",
            "result": {"results": results, "summary": summary},
            "timestamp": time.time()
        }

        return {
            "test_run_id": test_run_id,
            "test_type": "api_validation",
            "success": overall_success,
            "summary": summary,
            "results": results
        }

    async def _get_test_report(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Get comprehensive test report and analysis"""
        test_run_id = args["test_run_id"]
        format_type = args.get("format", "summary")
        include_coverage = args.get("include_coverage", True)

        if test_run_id not in self.test_results:
            return {"error": f"Test run not found: {test_run_id}"}

        test_data = self.test_results[test_run_id]

        if format_type == "json":
            return {
                "test_run_id": test_run_id,
                "report": test_data,
                "format": "json"
            }

        elif format_type == "summary":
            # Generate human-readable summary
            summary = {
                "test_run_id": test_run_id,
                "test_type": test_data.get("type", "unknown"),
                "timestamp": test_data.get("timestamp", 0),
                "success": test_data.get("result", {}).get("exit_code") == 0,
                "summary": test_data.get("summary", {}),
                "execution_time": test_data.get("result", {}).get("execution_time", 0)
            }

            return {
                "test_run_id": test_run_id,
                "report": summary,
                "format": "summary"
            }

        elif format_type == "html":
            # Generate HTML report
            html_report = f"""
            <html>
            <head><title>Test Report - {test_run_id}</title></head>
            <body>
                <h1>Test Report: {test_run_id}</h1>
                <p><strong>Test Type:</strong> {test_data.get("type", "unknown")}</p>
                <p><strong>Timestamp:</strong> {test_data.get("timestamp", 0)}</p>
                <p><strong>Success:</strong> {test_data.get("result", {}).get("exit_code") == 0}</p>
                <pre>{json.dumps(test_data, indent=2)}</pre>
            </body>
            </html>
            """

            return {
                "test_run_id": test_run_id,
                "report": html_report,
                "format": "html"
            }

    def _parse_test_results(self, test_type: str, output: str) -> Dict[str, Any]:
        """Parse test results based on test type"""
        summary = {
            "total_tests": 0,
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "duration": 0
        }

        try:
            if test_type == "python":
                # Parse pytest output
                lines = output.split('\n')
                for line in reversed(lines):
                    if "passed" in line.lower() and "failed" in line.lower():
                        parts = line.split()
                        for i, part in enumerate(parts):
                            if part == "passed" and i > 0:
                                try:
                                    summary["passed"] = int(parts[i-1])
                                except ValueError:
                                    pass
                            elif part == "failed" and i > 0:
                                try:
                                    summary["failed"] = int(parts[i-1])
                                except ValueError:
                                    pass

            elif test_type == "java":
                # Parse JUnit/Maven output
                if "Tests run:" in output:
                    for line in output.split('\n'):
                        if "Tests run:" in line:
                            parts = line.split()
                            try:
                                summary["total_tests"] = int(parts[2])
                                summary["failed"] = int(parts[4].rstrip(','))
                                summary["passed"] = summary["total_tests"] - summary["failed"]
                            except (ValueError, IndexError):
                                pass

            elif test_type == "javascript":
                # Parse Jest/Mocha output
                if "Tests:" in output:
                    for line in output.split('\n'):
                        if "Tests:" in line:
                            parts = line.split()
                            try:
                                summary["total_tests"] = int(parts[1])
                                summary["passed"] = int(parts[3].rstrip(','))
                                summary["failed"] = summary["total_tests"] - summary["passed"]
                            except (ValueError, IndexError):
                                pass

        except Exception as e:
            logger.warning(f"Error parsing test results: {str(e)}")

        return summary

    def _parse_integration_results(self, tool: str, output: str) -> Dict[str, Any]:
        """Parse integration test results"""
        summary = {
            "total_requests": 0,
            "passed": 0,
            "failed": 0,
            "response_time_avg": 0
        }

        try:
            if tool == "postman" or tool == "newman":
                # Parse Newman output
                if "requests:" in output.lower():
                    for line in output.split('\n'):
                        if "requests:" in line.lower():
                            parts = line.split()
                            try:
                                summary["total_requests"] = int(parts[1])
                            except (ValueError, IndexError):
                                pass

                if "passed" in output.lower() and "failed" in output.lower():
                    for line in reversed(output.split('\n')):
                        if "passed" in line.lower() and "failed" in line.lower():
                            parts = line.split()
                            for i, part in enumerate(parts):
                                if part == "passed" and i > 0:
                                    try:
                                        summary["passed"] = int(parts[i-1])
                                    except ValueError:
                                        pass
                                elif part == "failed" and i > 0:
                                    try:
                                        summary["failed"] = int(parts[i-1])
                                    except ValueError:
                                        pass

        except Exception as e:
            logger.warning(f"Error parsing integration results: {str(e)}")

        return summary

    def _parse_performance_results(self, tool: str, output: str) -> Dict[str, Any]:
        """Parse performance test results"""
        summary = {"tool": tool}

        try:
            if tool == "jmeter":
                # Parse JMeter summary
                if "summary =" in output:
                    for line in output.split('\n'):
                        if "summary =" in line:
                            parts = line.split()
                            try:
                                summary["requests_per_sec"] = float(parts[6])
                                summary["avg_response_time"] = float(parts[8])
                                summary["error_rate"] = float(parts[14].rstrip('%')) / 100
                            except (ValueError, IndexError):
                                pass

            elif tool == "k6":
                # Parse K6 output
                if "http_req_duration" in output:
                    for line in output.split('\n'):
                        if "http_req_duration" in line:
                            parts = line.split()
                            try:
                                summary["avg_response_time"] = float(parts[5])
                            except (ValueError, IndexError):
                                pass

            elif tool == "ab":
                # Parse Apache Bench output
                if "Requests per second:" in output:
                    for line in output.split('\n'):
                        if "Requests per second:" in line:
                            parts = line.split()
                            try:
                                summary["requests_per_sec"] = float(parts[3])
                            except (ValueError, IndexError):
                                pass

        except Exception as e:
            logger.warning(f"Error parsing performance results: {str(e)}")

        return summary

    def _parse_security_results(self, tool: str, output: str) -> Dict[str, Any]:
        """Parse security test results"""
        summary = {
            "tool": tool,
            "vulnerabilities_found": 0,
            "high_severity": 0,
            "medium_severity": 0,
            "low_severity": 0
        }

        try:
            # Count vulnerabilities by severity
            high_count = output.lower().count("high") - output.lower().count("highlight")
            medium_count = output.lower().count("medium")
            low_count = output.lower().count("low")

            summary["high_severity"] = high_count
            summary["medium_severity"] = medium_count
            summary["low_severity"] = low_count
            summary["vulnerabilities_found"] = high_count + medium_count + low_count

        except Exception as e:
            logger.warning(f"Error parsing security results: {str(e)}")

        return summary

def main():
    """Main function to run the enhanced MCP server"""
    server = EnhancedHybridTestingServer()

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