#!/usr/bin/env python3
"""
Hybrid Deployment Framework with MCP
Combines Model Context Protocol with traditional CLI tools for deployment and API access
"""

import asyncio
import subprocess
import json
import sys
import os
from datetime import datetime

class HybridDeploymentServer:
    def __init__(self):
        self.name = "hybrid-deployment"
        self.version = "1.0.0"
        
    async def list_tools(self):
        return [
            {
                "name": "deploy_to_coolify",
                "description": "Deploy application to Coolify platform",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "app_uuid": {"type": "string", "description": "Application UUID"},
                        "host": {"type": "string", "description": "Coolify host"},
                        "token": {"type": "string", "description": "API token"},
                        "force_rebuild": {"type": "boolean", "description": "Force rebuild"}
                    }
                }
            },
            {
                "name": "setup_ssh_tunnel",
                "description": "Setup SSH tunnel for Coolify API access",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "host": {"type": "string", "description": "Remote host"},
                        "port": {"type": "string", "description": "Local port"},
                        "remote_port": {"type": "string", "description": "Remote port"},
                        "key_path": {"type": "string", "description": "SSH key path"}
                    }
                }
            },
            {
                "name": "access_coolify_api",
                "description": "Access Coolify API endpoints",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "endpoint": {"type": "string", "description": "API endpoint"},
                        "method": {"enum": ["GET", "POST", "PUT", "DELETE"], "description": "HTTP method"},
                        "token": {"type": "string", "description": "API token"},
                        "data": {"type": "object", "description": "Request data"}
                    }
                }
            },
            {
                "name": "monitor_deployment",
                "description": "Monitor deployment status",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "deployment_uuid": {"type": "string", "description": "Deployment UUID"},
                        "token": {"type": "string", "description": "API token"},
                        "timeout": {"type": "integer", "description": "Timeout in seconds"}
                    }
                }
            },
            {
                "name": "validate_application",
                "description": "Validate deployed application health",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "app_url": {"type": "string", "description": "Application URL"},
                        "endpoints": {"type": "array", "description": "Endpoints to test"},
                        "expected_status": {"type": "integer", "description": "Expected HTTP status"}
                    }
                }
            },
            {
                "name": "configure_domain",
                "description": "Configure domain for application",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "app_uuid": {"type": "string", "description": "Application UUID"},
                        "domain": {"type": "string", "description": "Domain name"},
                        "token": {"type": "string", "description": "API token"}
                    }
                }
            }
        ]
    
    async def call_tool(self, name, arguments):
        try:
            if name == "deploy_to_coolify":
                return await self._deploy_to_coolify(arguments)
            elif name == "setup_ssh_tunnel":
                return await self._setup_ssh_tunnel(arguments)
            elif name == "access_coolify_api":
                return await self._access_coolify_api(arguments)
            elif name == "monitor_deployment":
                return await self._monitor_deployment(arguments)
            elif name == "validate_application":
                return await self._validate_application(arguments)
            elif name == "configure_domain":
                return await self._configure_domain(arguments)
            else:
                return {"error": f"Unknown tool: {name}"}
        except Exception as e:
            return {"error": str(e)}
    
    async def _deploy_to_coolify(self, args):
        app_uuid = args.get("app_uuid")
        host = args.get("host", "localhost:8000")
        token = args.get("token")
        force_rebuild = args.get("force_rebuild", True)
        
        if not app_uuid or not token:
            return {"error": "app_uuid and token are required"}
        
        # Trigger deployment
        deploy_data = {"uuid": app_uuid}
        if force_rebuild:
            deploy_data["force_rebuild"] = True
        
        cmd = [
            "curl", "-s", "-X", "POST",
            f"http://{host}/api/v1/deploy",
            "-H", f"Authorization: Bearer {token}",
            "-H", "Content-Type: application/json",
            "-d", json.dumps(deploy_data)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            try:
                response = json.loads(result.stdout)
                return {
                    "success": True,
                    "deployment_uuid": response.get("deployment_uuid"),
                    "message": "Deployment triggered successfully",
                    "response": response
                }
            except json.JSONDecodeError:
                return {
                    "success": False,
                    "error": "Invalid JSON response from Coolify API",
                    "stdout": result.stdout,
                    "stderr": result.stderr
                }
        else:
            return {
                "success": False,
                "error": f"Deployment failed with exit code {result.returncode}",
                "stdout": result.stdout,
                "stderr": result.stderr
            }
    
    async def _setup_ssh_tunnel(self, args):
        host = args.get("host")
        port = args.get("port", "8000")
        remote_port = args.get("remote_port", "8000")
        key_path = args.get("key_path", "~/.ssh/coolify_doctorhealthy1")
        
        if not host:
            return {"error": "host is required"}
        
        # Check if tunnel already exists
        check_cmd = ["netstat", "-tlnp"]
        check_result = subprocess.run(check_cmd, capture_output=True, text=True)
        
        if f":{port}" in check_result.stdout:
            return {
                "success": True,
                "message": f"SSH tunnel already active on port {port}",
                "existing": True
            }
        
        # Start SSH tunnel
        cmd = [
            "ssh", "-i", key_path, "-N", "-L",
            f"{port}:localhost:{remote_port}",
            f"root@{host}"
        ]
        
        # Start tunnel in background
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        
        # Wait a moment and check if tunnel is established
        await asyncio.sleep(3)
        
        check_result = subprocess.run(["netstat", "-tlnp"], capture_output=True, text=True)
        
        if f":{port}" in check_result.stdout:
            return {
                "success": True,
                "message": f"SSH tunnel established on port {port}",
                "pid": process.pid
            }
        else:
            return {
                "success": False,
                "error": "Failed to establish SSH tunnel",
                "stderr": process.stderr.read().decode() if process.stderr else ""
            }
    
    async def _access_coolify_api(self, args):
        endpoint = args.get("endpoint")
        method = args.get("method", "GET")
        token = args.get("token")
        data = args.get("data", {})
        
        if not endpoint or not token:
            return {"error": "endpoint and token are required"}
        
        # Build curl command
        cmd = ["curl", "-s", "-X", method]
        
        # Add headers
        cmd.extend(["-H", f"Authorization: Bearer {token}"])
        cmd.extend(["-H", "Content-Type: application/json"])
        cmd.extend(["-H", "Accept: application/json"])
        
        # Add data if present
        if data and method in ["POST", "PUT", "PATCH"]:
            cmd.extend(["-d", json.dumps(data)])
        
        # Add endpoint
        if not endpoint.startswith("http"):
            endpoint = f"http://localhost:8000{endpoint}"
        cmd.append(endpoint)
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        try:
            response_data = json.loads(result.stdout) if result.stdout else {}
        except json.JSONDecodeError:
            response_data = {"raw_response": result.stdout}
        
        return {
            "success": result.returncode == 0,
            "status_code": result.returncode,
            "response": response_data,
            "stderr": result.stderr if result.stderr else None,
            "command": " ".join(cmd[:-1] + ["<endpoint>"])  # Hide endpoint for security
        }
    
    async def _monitor_deployment(self, args):
        deployment_uuid = args.get("deployment_uuid")
        token = args.get("token")
        timeout = args.get("timeout", 600)  # 10 minutes default
        
        if not deployment_uuid or not token:
            return {"error": "deployment_uuid and token are required"}
        
        start_time = datetime.now()
        status_history = []
        
        while (datetime.now() - start_time).seconds < timeout:
            # Check deployment status
            cmd = [
                "curl", "-s", "-X", "GET",
                f"http://localhost:8000/api/v1/deployments/{deployment_uuid}",
                "-H", f"Authorization: Bearer {token}"
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                try:
                    response = json.loads(result.stdout)
                    status = response.get("status", "unknown")
                    
                    status_entry = {
                        "timestamp": datetime.now().isoformat(),
                        "status": status,
                        "elapsed": (datetime.now() - start_time).seconds
                    }
                    status_history.append(status_entry)
                    
                    if status == "finished":
                        return {
                            "success": True,
                            "final_status": "completed",
                            "duration": (datetime.now() - start_time).seconds,
                            "status_history": status_history
                        }
                    elif status == "failed":
                        return {
                            "success": False,
                            "final_status": "failed",
                            "duration": (datetime.now() - start_time).seconds,
                            "status_history": status_history,
                            "error": response.get("error", "Deployment failed")
                        }
                    
                except json.JSONDecodeError:
                    pass
            
            await asyncio.sleep(10)  # Check every 10 seconds
        
        return {
            "success": False,
            "error": "Deployment monitoring timeout",
            "duration": timeout,
            "status_history": status_history
        }
    
    async def _validate_application(self, args):
        app_url = args.get("app_url")
        endpoints = args.get("endpoints", ["/health"])
        expected_status = args.get("expected_status", 200)
        
        if not app_url:
            return {"error": "app_url is required"}
        
        results = []
        
        for endpoint in endpoints:
            full_url = app_url.rstrip("/") + endpoint
            
            cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null", full_url]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            try:
                status_code = int(result.stdout.strip())
                success = status_code == expected_status
            except (ValueError, AttributeError):
                status_code = None
                success = False
            
            results.append({
                "endpoint": endpoint,
                "url": full_url,
                "status_code": status_code,
                "expected": expected_status,
                "success": success
            })
        
        overall_success = all(r["success"] for r in results)
        
        return {
            "success": overall_success,
            "results": results,
            "summary": {
                "total_endpoints": len(endpoints),
                "successful": sum(1 for r in results if r["success"]),
                "failed": sum(1 for r in results if not r["success"])
            }
        }
    
    async def _configure_domain(self, args):
        app_uuid = args.get("app_uuid")
        domain = args.get("domain")
        token = args.get("token")
        
        if not all([app_uuid, domain, token]):
            return {"error": "app_uuid, domain, and token are required"}
        
        # Configure domain in Coolify
        domain_data = {
            "uuid": app_uuid,
            "fqdn": domain,
            "redirect": False
        }
        
        cmd = [
            "curl", "-s", "-X", "POST",
            f"http://localhost:8000/api/v1/applications/{app_uuid}/domains",
            "-H", f"Authorization: Bearer {token}",
            "-H", "Content-Type: application/json",
            "-d", json.dumps(domain_data)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            try:
                response = json.loads(result.stdout)
                return {
                    "success": True,
                    "message": f"Domain {domain} configured successfully",
                    "response": response
                }
            except json.JSONDecodeError:
                return {
                    "success": True,  # Assume success if no JSON error
                    "message": f"Domain {domain} configured (non-JSON response)",
                    "raw_response": result.stdout
                }
        else:
            return {
                "success": False,
                "error": f"Domain configuration failed with exit code {result.returncode}",
                "stdout": result.stdout,
                "stderr": result.stderr
            }

# CLI interface
async def main():
    if len(sys.argv) < 2:
        print("Usage: python hybrid_deployment_server.py <command> [args...]")
        print("Commands: deploy, tunnel, api, monitor, validate, domain")
        return
    
    server = HybridDeploymentServer()
    command = sys.argv[1]
    
    if command == "deploy":
        app_uuid = sys.argv[2] if len(sys.argv) > 2 else None
        token = sys.argv[3] if len(sys.argv) > 3 else None
        result = await server.call_tool("deploy_to_coolify", {"app_uuid": app_uuid, "token": token})
    
    elif command == "tunnel":
        host = sys.argv[2] if len(sys.argv) > 2 else None
        result = await server.call_tool("setup_ssh_tunnel", {"host": host})
    
    elif command == "validate":
        app_url = sys.argv[2] if len(sys.argv) > 2 else None
        result = await server.call_tool("validate_application", {"app_url": app_url})
    
    else:
        print(f"Unknown command: {command}")
        return
    
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    asyncio.run(main())