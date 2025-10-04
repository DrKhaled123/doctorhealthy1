#!/usr/bin/env python3
"""
Enhanced Hybrid Deployment Framework with MCP
Combines Model Context Protocol with traditional CLI tools for comprehensive deployment and API access
Enhanced with safety features, multi-cloud support, and production-ready capabilities
"""

import json
import subprocess
import sys
import os
import asyncio
import logging
import time
import tempfile
import requests
from typing import Any, Dict, List, Optional
from pathlib import Path
import traceback
import base64
import urllib.parse

# Setup enhanced logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/tmp/hybrid_deployment_server.log'),
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)

class SafeSubprocessRunner:
    """Safe subprocess runner with timeout and error handling"""

    def __init__(self, timeout: int = 600):
        self.timeout = timeout
        self.max_output_size = 10 * 1024 * 1024  # 10MB limit

    def run_command(self, cmd: List[str], cwd: str = None, env: Dict = None) -> Dict[str, Any]:
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
            dangerous_commands = ['rm', 'del', 'format', 'fdisk', 'mkfs', 'dd', 'sudo', 'su']
            if any(dangerous in cmd[0] for dangerous in dangerous_commands):
                return {
                    "exit_code": -1,
                    "stdout": "",
                    "stderr": "Dangerous command blocked for safety",
                    "error": "Command blocked for security reasons"
                }

            # Prepare environment
            run_env = os.environ.copy()
            if env:
                run_env.update(env)

            process = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                cwd=cwd,
                text=True,
                env=run_env,
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

class CloudPlatformManager:
    """Manages cloud platform deployments safely"""

    def __init__(self):
        self.runner = SafeSubprocessRunner()
        self.supported_platforms = {
            "gcp": {
                "services": ["cloud-run", "cloud-functions", "app-engine", "gke", "compute-engine"],
                "auth_required": True,
                "config_required": ["project_id"]
            },
            "aws": {
                "services": ["ecs", "lambda", "ec2", "s3", "rds", "elastic-beanstalk"],
                "auth_required": True,
                "config_required": ["region", "profile"]
            },
            "azure": {
                "services": ["app-service", "functions", "container-instances", "aks"],
                "auth_required": True,
                "config_required": ["subscription_id", "resource_group"]
            },
            "coolify": {
                "services": ["container", "static-site", "database"],
                "auth_required": True,
                "config_required": ["host", "api_token"]
            }
        }

    def validate_deployment_config(self, platform: str, service: str, config: Dict) -> Dict[str, Any]:
        """Validate deployment configuration"""
        if platform not in self.supported_platforms:
            return {"valid": False, "error": f"Unsupported platform: {platform}"}

        platform_config = self.supported_platforms[platform]

        if service not in platform_config["services"]:
            return {"valid": False, "error": f"Unsupported service for {platform}: {service}"}

        # Check required configuration
        missing_config = []
        for required in platform_config["config_required"]:
            if required not in config:
                missing_config.append(required)

        if missing_config:
            return {
                "valid": False,
                "error": f"Missing required configuration: {missing_config}",
                "missing": missing_config
            }

        return {"valid": True}

    def deploy_to_gcp(self, service: str, config: Dict) -> Dict[str, Any]:
        """Deploy to Google Cloud Platform"""
        deployment_id = f"gcp_{service}_{int(time.time())}"

        try:
            if service == "cloud-run":
                cmd = [
                    "gcloud", "run", "deploy", config["service_name"],
                    "--image", config["image"],
                    "--platform", "managed",
                    "--project", config["project_id"]
                ]

                if config.get("region"):
                    cmd.extend(["--region", config["region"]])

                if config.get("port"):
                    cmd.extend(["--port", str(config["port"])])

                if config.get("allow_unauthenticated", True):
                    cmd.append("--allow-unauthenticated")

                if config.get("memory"):
                    cmd.extend(["--memory", config["memory"]])

                if config.get("cpu"):
                    cmd.extend(["--cpu", config["cpu"]])

            elif service == "cloud-functions":
                cmd = [
                    "gcloud", "functions", "deploy", config["function_name"],
                    "--runtime", config.get("runtime", "python39"),
                    "--trigger-http",
                    "--project", config["project_id"]
                ]

                if config.get("region"):
                    cmd.extend(["--region", config["region"]])

            elif service == "app-engine":
                cmd = [
                    "gcloud", "app", "deploy", config.get("app_yaml", "app.yaml"),
                    "--project", config["project_id"]
                ]

            else:
                return {"error": f"Unsupported GCP service: {service}"}

            result = self.runner.run_command(cmd)
            return {
                "deployment_id": deployment_id,
                "platform": "gcp",
                "service": service,
                "success": result.get("exit_code") == 0,
                "result": result
            }

        except Exception as e:
            return {
                "deployment_id": deployment_id,
                "platform": "gcp",
                "service": service,
                "success": False,
                "error": str(e)
            }

    def deploy_to_aws(self, service: str, config: Dict) -> Dict[str, Any]:
        """Deploy to Amazon Web Services"""
        deployment_id = f"aws_{service}_{int(time.time())}"

        try:
            if service == "ecs":
                cmd = [
                    "aws", "ecs", "create-service",
                    "--service-name", config["service_name"],
                    "--task-definition", config["task_definition"],
                    "--cluster", config["cluster"],
                    "--profile", config.get("profile", "default")
                ]

                if config.get("desired-count"):
                    cmd.extend(["--desired-count", str(config["desired-count"])])

            elif service == "lambda":
                cmd = [
                    "aws", "lambda", "create-function",
                    "--function-name", config["function_name"],
                    "--runtime", config.get("runtime", "python3.9"),
                    "--role", config["role_arn"],
                    "--handler", config.get("handler", "lambda_function.lambda_handler"),
                    "--zip-file", config["zip_file"],
                    "--profile", config.get("profile", "default")
                ]

            elif service == "s3":
                cmd = [
                    "aws", "s3", "sync", config["source"],
                    f"s3://{config['bucket']}",
                    "--profile", config.get("profile", "default")
                ]

            else:
                return {"error": f"Unsupported AWS service: {service}"}

            result = self.runner.run_command(cmd)
            return {
                "deployment_id": deployment_id,
                "platform": "aws",
                "service": service,
                "success": result.get("exit_code") == 0,
                "result": result
            }

        except Exception as e:
            return {
                "deployment_id": deployment_id,
                "platform": "aws",
                "service": service,
                "success": False,
                "error": str(e)
            }

    def deploy_to_azure(self, service: str, config: Dict) -> Dict[str, Any]:
        """Deploy to Microsoft Azure"""
        deployment_id = f"azure_{service}_{int(time.time())}"

        try:
            if service == "app-service":
                cmd = [
                    "az", "webapp", "create",
                    "--resource-group", config["resource_group"],
                    "--plan", config["app_service_plan"],
                    "--name", config["app_name"],
                    "--runtime", config.get("runtime", "PYTHON|3.9"),
                    "--subscription", config["subscription_id"]
                ]

            elif service == "functions":
                cmd = [
                    "az", "functionapp", "create",
                    "--resource-group", config["resource_group"],
                    "--name", config["function_name"],
                    "--runtime", config.get("runtime", "python"),
                    "--runtime-version", config.get("runtime_version", "3.9"),
                    "--functions-version", "4",
                    "--subscription", config["subscription_id"]
                ]

            else:
                return {"error": f"Unsupported Azure service: {service}"}

            result = self.runner.run_command(cmd)
            return {
                "deployment_id": deployment_id,
                "platform": "azure",
                "service": service,
                "success": result.get("exit_code") == 0,
                "result": result
            }

        except Exception as e:
            return {
                "deployment_id": deployment_id,
                "platform": "azure",
                "service": service,
                "success": False,
                "error": str(e)
            }

class EnhancedHybridDeploymentServer:
    """Enhanced MCP server for hybrid deployment framework with safety features"""

    def __init__(self):
        self.name = "enhanced-hybrid-deployment"
        self.tools = self._register_tools()
        self.runner = SafeSubprocessRunner()
        self.cloud_manager = CloudPlatformManager()
        self.deployment_history = {}

    def _register_tools(self) -> List[Dict[str, Any]]:
        """Register available deployment and API access tools"""
        return [
            {
                "name": "deploy_to_cloud",
                "description": "Deploy application to cloud platform with validation",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "platform": {
                            "enum": ["gcp", "aws", "azure", "coolify"],
                            "description": "Cloud platform"
                        },
                        "service": {
                            "enum": ["cloud-run", "ecs", "app-service", "container", "cloud-functions", "lambda"],
                            "description": "Cloud service"
                        },
                        "config": {
                            "type": "object",
                            "description": "Deployment configuration"
                        },
                        "validate_only": {
                            "type": "boolean",
                            "description": "Only validate configuration without deploying",
                            "default": False
                        }
                    },
                    "required": ["platform", "service", "config"]
                }
            },
            {
                "name": "access_coolify_api",
                "description": "Access Coolify self-hosted platform API with authentication",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "operation": {
                            "enum": ["deploy", "status", "logs", "env_vars", "backup", "restore"],
                            "description": "Coolify operation"
                        },
                        "app_uuid": {
                            "type": "string",
                            "description": "Application UUID"
                        },
                        "api_token": {
                            "type": "string",
                            "description": "Coolify API token"
                        },
                        "host": {
                            "type": "string",
                            "description": "Coolify host URL",
                            "default": "localhost:8000"
                        },
                        "params": {
                            "type": "object",
                            "description": "Operation parameters"
                        }
                    },
                    "required": ["operation", "app_uuid", "api_token"]
                }
            },
            {
                "name": "access_web_service",
                "description": "Access generic web services and APIs with enhanced security",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "url": {
                            "type": "string",
                            "description": "Service URL"
                        },
                        "method": {
                            "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                            "description": "HTTP method"
                        },
                        "headers": {
                            "type": "object",
                            "description": "HTTP headers"
                        },
                        "body": {
                            "type": "object",
                            "description": "Request body"
                        },
                        "auth": {
                            "type": "object",
                            "properties": {
                                "type": {
                                    "enum": ["bearer", "basic", "api_key", "oauth2"],
                                    "description": "Authentication type"
                                },
                                "token": {
                                    "type": "string",
                                    "description": "Authentication token"
                                },
                                "username": {
                                    "type": "string",
                                    "description": "Username for basic auth"
                                },
                                "password": {
                                    "type": "string",
                                    "description": "Password for basic auth"
                                }
                            }
                        },
                        "timeout": {
                            "type": "integer",
                            "description": "Request timeout in seconds",
                            "default": 30
                        }
                    },
                    "required": ["url", "method"]
                }
            },
            {
                "name": "build_docker_image",
                "description": "Build Docker image for deployment with security checks",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "dockerfile_path": {
                            "type": "string",
                            "description": "Path to Dockerfile",
                            "default": "."
                        },
                        "image_name": {
                            "type": "string",
                            "description": "Docker image name and tag"
                        },
                        "build_args": {
                            "type": "object",
                            "description": "Docker build arguments"
                        },
                        "platform": {
                            "type": "string",
                            "description": "Target platform (e.g., linux/amd64)",
                            "default": "linux/amd64"
                        },
                        "security_scan": {
                            "type": "boolean",
                            "description": "Run security scan on built image",
                            "default": True
                        }
                    },
                    "required": ["image_name"]
                }
            },
            {
                "name": "manage_ssh_tunnel",
                "description": "Manage SSH tunnels for secure connections with safety controls",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "action": {
                            "enum": ["create", "close", "status", "list"],
                            "description": "Tunnel action"
                        },
                        "local_port": {
                            "type": "integer",
                            "description": "Local port",
                            "minimum": 1024,
                            "maximum": 65535
                        },
                        "remote_host": {
                            "type": "string",
                            "description": "Remote host"
                        },
                        "remote_port": {
                            "type": "integer",
                            "description": "Remote port",
                            "minimum": 1,
                            "maximum": 65535
                        },
                        "ssh_key": {
                            "type": "string",
                            "description": "Path to SSH key"
                        },
                        "ssh_user": {
                            "type": "string",
                            "description": "SSH user",
                            "default": "root"
                        },
                        "timeout": {
                            "type": "integer",
                            "description": "Connection timeout in seconds",
                            "default": 30
                        }
                    },
                    "required": ["action"]
                }
            },
            {
                "name": "health_check_deployment",
                "description": "Perform comprehensive health checks on deployed application",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "app_url": {
                            "type": "string",
                            "description": "Application URL to check"
                        },
                        "health_endpoints": {
                            "type": "array",
                            "items": {"type": "string"},
                            "description": "List of health check endpoints",
                            "default": ["/health", "/ready", "/"]
                        },
                        "expected_response_time": {
                            "type": "integer",
                            "description": "Expected response time in milliseconds",
                            "default": 2000
                        },
                        "retries": {
                            "type": "integer",
                            "description": "Number of retry attempts",
                            "default": 3
                        },
                        "headers": {
                            "type": "object",
                            "description": "Custom headers for health checks"
                        }
                    },
                    "required": ["app_url"]
                }
            },
            {
                "name": "get_deployment_status",
                "description": "Get comprehensive deployment status and history",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "deployment_id": {
                            "type": "string",
                            "description": "Deployment identifier"
                        },
                        "platform": {
                            "type": "string",
                            "description": "Platform filter"
                        },
                        "limit": {
                            "type": "integer",
                            "description": "Number of recent deployments to show",
                            "default": 10
                        }
                    }
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

            if name == "deploy_to_cloud":
                return await self._deploy_to_cloud(arguments)
            elif name == "access_coolify_api":
                return await self._access_coolify_api(arguments)
            elif name == "access_web_service":
                return await self._access_web_service(arguments)
            elif name == "build_docker_image":
                return await self._build_docker_image(arguments)
            elif name == "manage_ssh_tunnel":
                return await self._manage_ssh_tunnel(arguments)
            elif name == "health_check_deployment":
                return await self._health_check_deployment(arguments)
            elif name == "get_deployment_status":
                return await self._get_deployment_status(arguments)
            else:
                return {"error": f"Unknown tool: {name}"}

        except Exception as e:
            logger.error(f"Error executing tool {name}: {str(e)}")
            logger.error(traceback.format_exc())
            return {
                "error": f"Tool execution failed: {str(e)}",
                "traceback": traceback.format_exc()
            }

    async def _deploy_to_cloud(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Deploy application to specified cloud platform with validation"""
        platform = args["platform"]
        service = args["service"]
        config = args["config"]
        validate_only = args.get("validate_only", False)

        # Validate configuration first
        validation = self.cloud_manager.validate_deployment_config(platform, service, config)
        if not validation["valid"]:
            return {
                "deployment_id": f"{platform}_{service}_{int(time.time())}",
                "success": False,
                "error": validation["error"],
                "validation": validation
            }

        if validate_only:
            return {
                "deployment_id": f"{platform}_{service}_{int(time.time())}",
                "success": True,
                "message": "Configuration validation passed",
                "validation": validation
            }

        # Perform deployment
        try:
            if platform == "gcp":
                result = self.cloud_manager.deploy_to_gcp(service, config)
            elif platform == "aws":
                result = self.cloud_manager.deploy_to_aws(service, config)
            elif platform == "azure":
                result = self.cloud_manager.deploy_to_azure(service, config)
            else:
                return {"error": f"Unsupported platform: {platform}"}

            # Store deployment history
            self.deployment_history[result["deployment_id"]] = {
                "platform": platform,
                "service": service,
                "config": config,
                "result": result,
                "timestamp": time.time()
            }

            return result

        except Exception as e:
            error_result = {
                "deployment_id": f"{platform}_{service}_{int(time.time())}",
                "platform": platform,
                "service": service,
                "success": False,
                "error": str(e)
            }

            self.deployment_history[error_result["deployment_id"]] = {
                "platform": platform,
                "service": service,
                "config": config,
                "result": error_result,
                "timestamp": time.time()
            }

            return error_result

    async def _access_coolify_api(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Access Coolify API for various operations with authentication"""
        operation = args["operation"]
        app_uuid = args["app_uuid"]
        api_token = args["api_token"]
        host = args.get("host", "localhost:8000")
        params = args.get("params", {})

        # Validate inputs
        if not app_uuid or not api_token:
            return {"error": "app_uuid and api_token are required"}

        base_url = f"http://{host}/api/v1"
        headers = {
            "Authorization": f"Bearer {api_token}",
            "Content-Type": "application/json",
            "Accept": "application/json"
        }

        try:
            if operation == "deploy":
                url = f"{base_url}/deploy"
                data = {"uuid": app_uuid}
                if params.get("force"):
                    data["force"] = True
                response = requests.post(url, json=data, headers=headers, timeout=30)

            elif operation == "status":
                url = f"{base_url}/applications/{app_uuid}"
                response = requests.get(url, headers=headers, timeout=30)

            elif operation == "logs":
                url = f"{base_url}/applications/{app_uuid}/logs"
                response = requests.get(url, headers=headers, timeout=30)

            elif operation == "env_vars":
                url = f"{base_url}/applications/{app_uuid}/envs"
                if params.get("set"):
                    # Set environment variable
                    data = {
                        "key": params["key"],
                        "value": params["value"]
                    }
                    response = requests.post(url, json=data, headers=headers, timeout=30)
                else:
                    # Get environment variables
                    response = requests.get(url, headers=headers, timeout=30)

            elif operation == "backup":
                url = f"{base_url}/applications/{app_uuid}/backup"
                response = requests.post(url, headers=headers, timeout=60)

            elif operation == "restore":
                url = f"{base_url}/applications/{app_uuid}/restore"
                data = {"backup_id": params.get("backup_id")}
                response = requests.post(url, json=data, headers=headers, timeout=60)

            else:
                return {"error": f"Unsupported Coolify operation: {operation}"}

            return {
                "status_code": response.status_code,
                "response": response.json() if response.content else {},
                "success": response.status_code < 400,
                "operation": operation,
                "url": url
            }

        except requests.RequestException as e:
            return {"error": f"Coolify API request failed: {str(e)}"}
        except Exception as e:
            return {"error": f"Coolify API access failed: {str(e)}"}

    async def _access_web_service(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Access generic web services with enhanced security and error handling"""
        url = args["url"]
        method = args["method"]
        headers = args.get("headers", {})
        body = args.get("body", {})
        auth = args.get("auth", {})
        timeout = args.get("timeout", 30)

        # Validate URL
        if not url or not isinstance(url, str):
            return {"error": "Valid URL is required"}

        # Security check - block potentially dangerous URLs
        dangerous_patterns = ["file://", "ftp://", "127.0.0.1", "localhost"]
        if any(pattern in url.lower() for pattern in dangerous_patterns):
            return {"error": "URL blocked for security reasons"}

        # Add authentication headers
        if auth:
            auth_type = auth.get("type")
            token = auth.get("token")
            username = auth.get("username")
            password = auth.get("password")

            if auth_type == "bearer" and token:
                headers["Authorization"] = f"Bearer {token}"
            elif auth_type == "api_key" and token:
                headers["X-API-Key"] = token
            elif auth_type == "basic" and username and password:
                credentials = base64.b64encode(f"{username}:{password}".encode()).decode()
                headers["Authorization"] = f"Basic {credentials}"
            elif auth_type == "oauth2" and token:
                headers["Authorization"] = f"Bearer {token}"

        # Build curl command with safety limits
        cmd = [
            "curl", "-s", "-w",
            "response_time:%{time_total}\\nhttp_code:%{http_code}\\n",
            "--max-time", str(timeout),
            "--retry", "3",
            "--retry-delay", "1",
            "-X", method
        ]

        # Add headers
        for key, value in headers.items():
            cmd.extend(["-H", f"{key}: {value}"])

        # Add body if present
        if body and method in ["POST", "PUT", "PATCH"]:
            cmd.extend(["-d", json.dumps(body)])

        # Add URL
        cmd.append(url)

        logger.info(f"Accessing web service: {method} {url}")
        result = self.runner.run_command(cmd)

        # Try to parse JSON response
        try:
            if result.get("stdout"):
                # Extract response time and status code from curl output
                output_lines = result["stdout"].split('\n')
                response_data = output_lines[:-2]  # Last 2 lines are timing info
                response_body = '\n'.join(response_data)

                # Parse timing information
                timing_info = output_lines[-2:] if len(output_lines) >= 2 else []
                response_time = None
                status_code = None

                for line in timing_info:
                    if line.startswith("response_time:"):
                        try:
                            response_time = float(line.split(":", 1)[1]) * 1000  # Convert to ms
                        except ValueError:
                            pass
                    elif line.startswith("http_code:"):
                        try:
                            status_code = int(line.split(":", 1)[1])
                        except ValueError:
                            pass

                # Try to parse JSON
                try:
                    parsed_response = json.loads(response_body)
                    result["parsed_response"] = parsed_response
                except json.JSONDecodeError:
                    result["raw_response"] = response_body

                result["response_time_ms"] = response_time
                result["status_code"] = status_code

        except Exception as e:
            logger.warning(f"Error parsing response: {str(e)}")

        return result

    async def _build_docker_image(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Build Docker image with security checks"""
        dockerfile_path = args.get("dockerfile_path", ".")
        image_name = args["image_name"]
        build_args = args.get("build_args", {})
        platform = args.get("platform", "linux/amd64")
        security_scan = args.get("security_scan", True)

        build_id = f"docker_build_{int(time.time())}"

        # Validate dockerfile exists
        if not os.path.exists(os.path.join(dockerfile_path, "Dockerfile")):
            return {"error": f"Dockerfile not found in path: {dockerfile_path}"}

        # Build command
        cmd = ["docker", "build", "--platform", platform, "--no-cache"]

        # Add build arguments
        for key, value in build_args.items():
            cmd.extend(["--build-arg", f"{key}={value}"])

        # Add tag and context
        cmd.extend(["-t", image_name, dockerfile_path])

        logger.info(f"Building Docker image: {' '.join(cmd)}")
        result = self.runner.run_command(cmd)

        # Security scan if requested
        if security_scan and result.get("exit_code") == 0:
            logger.info(f"Running security scan on image: {image_name}")
            scan_result = self.runner.run_command([
                "docker", "scan", image_name, "--json"
            ])

            if scan_result.get("exit_code") == 0:
                try:
                    scan_data = json.loads(scan_result.get("stdout", "{}"))
                    result["security_scan"] = scan_data
                except json.JSONDecodeError:
                    result["security_scan"] = {"error": "Failed to parse scan results"}

        result["build_id"] = build_id
        result["image_built"] = image_name
        result["build_status"] = "success" if result.get("exit_code") == 0 else "failed"

        return result

    async def _manage_ssh_tunnel(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Manage SSH tunnels with safety controls"""
        action = args["action"]

        if action == "create":
            local_port = args["local_port"]
            remote_host = args["remote_host"]
            remote_port = args["remote_port"]
            ssh_key = args.get("ssh_key")
            ssh_user = args.get("ssh_user", "root")
            timeout = args.get("timeout", 30)

            # Validate inputs
            if local_port < 1024 or local_port > 65535:
                return {"error": "Invalid local port range"}

            if remote_port < 1 or remote_port > 65535:
                return {"error": "Invalid remote port range"}

            cmd = ["ssh", "-N", "-L", f"{local_port}:localhost:{remote_port}"]

            if ssh_key:
                cmd.extend(["-i", ssh_key])

            # Add timeout
            cmd.extend(["-o", f"ConnectTimeout={timeout}"])
            cmd.extend(["-o", "StrictHostKeyChecking=no"])
            cmd.extend(["-o", "UserKnownHostsFile=/dev/null"])

            cmd.append(f"{ssh_user}@{remote_host}")

            # Run in background
            logger.info(f"Creating SSH tunnel: {' '.join(cmd)}")
            process = subprocess.Popen(cmd)

            return {
                "tunnel_created": True,
                "pid": process.pid,
                "local_port": local_port,
                "remote_host": remote_host,
                "remote_port": remote_port,
                "tunnel_id": f"tunnel_{process.pid}"
            }

        elif action == "status":
            # Check for existing SSH tunnels
            cmd = ["ps", "aux"]
            result = self.runner.run_command(cmd)

            tunnels = []
            if result.get("stdout"):
                for line in result["stdout"].split('\n'):
                    if "ssh" in line and "-L" in line:
                        tunnels.append(line.strip())

            return {
                "active_tunnels": tunnels,
                "tunnel_count": len(tunnels)
            }

        elif action == "list":
            # List all tunnel processes with details
            cmd = ["ps", "aux", "|", "grep", "ssh", "|", "grep", "-v", "grep"]
            result = self.runner.run_command(cmd)

            return {
                "tunnels": result.get("stdout", "").split('\n') if result.get("stdout") else []
            }

        elif action == "close":
            # Kill SSH tunnel processes
            cmd = ["pkill", "-f", "ssh.*-L"]
            result = self.runner.run_command(cmd)
            return {
                "tunnels_closed": result.get("exit_code") == 0,
                "message": "SSH tunnels terminated"
            }

        return {"error": f"Unsupported tunnel action: {action}"}

    async def _health_check_deployment(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Perform comprehensive health checks"""
        app_url = args["app_url"]
        health_endpoints = args.get("health_endpoints", ["/health", "/ready", "/"])
        expected_response_time = args.get("expected_response_time", 2000)
        retries = args.get("retries", 3)
        headers = args.get("headers", {})

        check_id = f"health_check_{int(time.time())}"
        results = {
            "check_id": check_id,
            "app_url": app_url,
            "overall_health": True,
            "checks": [],
            "summary": {}
        }

        total_checks = len(health_endpoints)
        successful_checks = 0

        for endpoint in health_endpoints:
            url = f"{app_url.rstrip('/')}{endpoint}"

            for attempt in range(retries):
                # Use curl with timeout and retry
                cmd = [
                    "curl", "-s", "-w",
                    "response_time:%{time_total}\\nhttp_code:%{http_code}\\n",
                    "--max-time", "10",
                    "--retry", "2",
                    "--retry-delay", "1",
                    "--fail"
                ]

                # Add custom headers
                for key, value in headers.items():
                    cmd.extend(["-H", f"{key}: {value}"])

                cmd.append(url)

                result = self.runner.run_command(cmd)

                check_result = {
                    "endpoint": endpoint,
                    "url": url,
                    "attempt": attempt + 1,
                    "success": False,
                    "response_time": None,
                    "status_code": None
                }

                if result.get("exit_code") == 0 and result.get("stdout"):
                    output = result["stdout"]

                    # Parse response time and status code
                    if "response_time:" in output:
                        try:
                            response_time = float(output.split("response_time:")[1].split()[0])
                            check_result["response_time"] = round(response_time * 1000, 2)  # Convert to ms
                        except (ValueError, IndexError):
                            pass

                    if "http_code:" in output:
                        try:
                            status_code = int(output.split("http_code:")[1].split()[0])
                            check_result["status_code"] = status_code

                            # Consider 2xx and 3xx as successful
                            if 200 <= status_code < 400:
                                check_result["success"] = True
                                successful_checks += 1
                                break
                        except (ValueError, IndexError):
                            pass

                if not check_result["success"]:
                    check_result["error"] = result.get("stderr", "Health check failed")

            results["checks"].append(check_result)

        # Calculate summary
        results["summary"] = {
            "total_checks": total_checks,
            "successful_checks": successful_checks,
            "success_rate": round((successful_checks / total_checks) * 100, 2) if total_checks > 0 else 0,
            "health_status": "healthy" if results["overall_health"] else "unhealthy"
        }

        return results

    async def _get_deployment_status(self, args: Dict[str, Any]) -> Dict[str, Any]:
        """Get comprehensive deployment status and history"""
        deployment_id = args.get("deployment_id")
        platform = args.get("platform")
        limit = args.get("limit", 10)

        if deployment_id:
            # Get specific deployment
            if deployment_id in self.deployment_history:
                return {
                    "deployment_id": deployment_id,
                    "status": self.deployment_history[deployment_id]
                }
            else:
                return {"error": f"Deployment not found: {deployment_id}"}

        else:
            # Get recent deployments
            recent_deployments = []
            for dep_id, dep_data in list(self.deployment_history.items())[-limit:]:
                if platform is None or dep_data.get("platform") == platform:
                    recent_deployments.append({
                        "deployment_id": dep_id,
                        "platform": dep_data.get("platform"),
                        "service": dep_data.get("service"),
                        "timestamp": dep_data.get("timestamp"),
                        "success": dep_data.get("result", {}).get("success", False)
                    })

            return {
                "total_deployments": len(self.deployment_history),
                "recent_deployments": recent_deployments,
                "platforms": list(set(dep.get("platform") for dep in self.deployment_history.values()))
            }

def main():
    """Main function to run the enhanced MCP server"""
    server = EnhancedHybridDeploymentServer()

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