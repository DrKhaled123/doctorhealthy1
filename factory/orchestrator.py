import asyncio
import json
import time
import os
from typing import Dict, Any, Optional
from dataclasses import dataclass
from enum import Enum

# Import our custom modules
from factory_orchestrator import FactoryOrchestrator, PipelineStage
from browser_testing_agent import BrowserTestingAgent
from autofix_agent import AutofixAgent
from memo_ai_memory import MemoAIMemory

class AgentStatus(Enum):
    IDLE = "idle"
    WORKING = "working"
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class AgentResult:
    agent_name: str
    status: AgentStatus
    data: Dict[str, Any]
    execution_time: float
    error: Optional[str] = None

class LightAgent:
    """Lightweight agent for autonomous task execution"""

    def __init__(self, name: str, backend_config: Dict = None):
        self.name = name
        self.backend_config = backend_config or {}
        self.status = AgentStatus.IDLE
        self.memory = None
        self.tools = {}

        # Initialize components
        self._initialize_components()

    def _initialize_components(self):
        """Initialize agent components"""
        try:
            # Initialize memory system
            self.memory = MemoAIMemory()

            # Register available tools
            self.tools = {
                "memory": self.memory,
                "browser_test": BrowserTestingAgent(),
                "autofix": AutofixAgent()
            }

            print(f"Agent {self.name} initialized successfully")

        except Exception as e:
            print(f"Failed to initialize agent {self.name}: {e}")
            raise

    async def execute_task(self, task_spec: str, context: Dict = None) -> AgentResult:
        """Execute a task autonomously"""
        start_time = time.time()
        self.status = AgentStatus.WORKING

        try:
            # Store task in memory
            await self.memory.remember("task_started", {
                "task": task_spec,
                "context": context or {},
                "agent": self.name
            })

            # Analyze task and determine approach
            approach = await self._analyze_task(task_spec)

            # Execute based on approach
            if approach["type"] == "code_generation":
                result = await self._execute_code_generation(approach, context)
            elif approach["type"] == "testing":
                result = await self._execute_testing(approach, context)
            elif approach["type"] == "fixing":
                result = await self._execute_fixing(approach, context)
            else:
                result = await self._execute_generic_task(approach, context)

            # Store result in memory
            await self.memory.remember("task_completed", {
                "task": task_spec,
                "result": result,
                "execution_time": time.time() - start_time
            })

            self.status = AgentStatus.COMPLETED
            return AgentResult(
                agent_name=self.name,
                status=AgentStatus.COMPLETED,
                data=result,
                execution_time=time.time() - start_time
            )

        except Exception as e:
            error_msg = str(e)
            self.status = AgentStatus.FAILED

            # Store failure in memory
            await self.memory.remember("task_failed", {
                "task": task_spec,
                "error": error_msg,
                "execution_time": time.time() - start_time
            })

            return AgentResult(
                agent_name=self.name,
                status=AgentStatus.FAILED,
                data={},
                execution_time=time.time() - start_time,
                error=error_msg
            )

    async def _analyze_task(self, task_spec: str) -> Dict:
        """Analyze task to determine execution approach"""
        # Simple rule-based analysis
        analysis = {
            "type": "generic",
            "complexity": "medium",
            "tools_needed": []
        }

        task_lower = task_spec.lower()

        if any(keyword in task_lower for keyword in ["code", "implement", "write", "create"]):
            analysis["type"] = "code_generation"
            analysis["tools_needed"] = ["memory", "autofix"]
        elif any(keyword in task_lower for keyword in ["test", "verify", "check"]):
            analysis["type"] = "testing"
            analysis["tools_needed"] = ["browser_test", "memory"]
        elif any(keyword in task_lower for keyword in ["fix", "error", "bug", "debug"]):
            analysis["type"] = "fixing"
            analysis["tools_needed"] = ["autofix", "memory"]
        elif any(keyword in task_lower for keyword in ["deploy", "monitor", "validate"]):
            analysis["type"] = "deployment"
            analysis["tools_needed"] = ["memory"]

        return analysis

    async def _execute_code_generation(self, approach: Dict, context: Dict) -> Dict:
        """Execute code generation task"""
        # Mock code generation
        generated_code = f"""
# Generated code based on: {context.get('spec', 'unknown specification')}

def main():
    print("Hello, World!")
    return True

if __name__ == "__main__":
    main()
"""

        # Apply autofix if available
        if "autofix" in self.tools:
            fix_result = await self.tools["autofix"].fix_code(generated_code)
            if fix_result["success"]:
                generated_code = fix_result["fixed_code"]

        return {
            "generated_code": generated_code,
            "language": "python",
            "complexity": approach["complexity"],
            "fixes_applied": len(approach["tools_needed"])
        }

    async def _execute_testing(self, approach: Dict, context: Dict) -> Dict:
        """Execute testing task"""
        # Use browser testing agent
        if "browser_test" in self.tools:
            browser_agent = self.tools["browser_test"]

            # Generate test cases
            test_cases = await browser_agent.generate_test_cases(
                context.get("app_description", "web application")
            )

            # Run tests
            test_url = context.get("test_url", "http://localhost:3000")
            results = await browser_agent.test_web_application(test_url, test_cases)

            return {
                "test_results": results,
                "tests_run": len(test_cases),
                "passed": results.get("passed", 0),
                "failed": results.get("failed", 0)
            }

        return {"error": "Browser testing not available"}

    async def _execute_fixing(self, approach: Dict, context: Dict) -> Dict:
        """Execute code fixing task"""
        # Use autofix agent
        if "autofix" in self.tools:
            code_to_fix = context.get("code", "")
            errors = context.get("errors", [])

            fix_result = await self.tools["autofix"].fix_code(code_to_fix, errors)

            return {
                "fix_result": fix_result,
                "fixes_applied": len(fix_result["fixes_applied"]),
                "success": fix_result["success"]
            }

        return {"error": "Autofix not available"}

    async def _execute_generic_task(self, approach: Dict, context: Dict) -> Dict:
        """Execute generic task"""
        return {
            "message": f"Generic task executed: {approach}",
            "context_used": bool(context),
            "tools_available": list(self.tools.keys())
        }

    async def get_status(self) -> Dict:
        """Get current agent status"""
        return {
            "name": self.name,
            "status": self.status.value,
            "tools": list(self.tools.keys()),
            "memory_stats": await self.memory.get_memory_stats() if self.memory else {}
        }

class LightOrchestrator:
    """Lightweight orchestrator for managing multiple agents"""

    def __init__(self, backend_config: Dict = None):
        self.backend_config = backend_config or {}
        self.agents = {}
        self.redis_client = None
        self.memory = None

        # Initialize backend
        self._initialize_backend()

    def _initialize_backend(self):
        """Initialize backend services"""
        try:
            # Initialize Redis
            import redis.asyncio as redis
            self.redis_client = redis.Redis(
                host=self.backend_config.get("host", "localhost"),
                port=self.backend_config.get("port", 6379),
                db=self.backend_config.get("db", 0),
                decode_responses=True
            )

            # Initialize shared memory
            self.memory = MemoAIMemory(redis_client=self.redis_client)

            print("Orchestrator backend initialized")

        except Exception as e:
            print(f"Backend initialization failed: {e}")

    def register_agent(self, name: str, agent_class=LightAgent) -> LightAgent:
        """Register a new agent"""
        agent = agent_class(name, self.backend_config)
        self.agents[name] = agent
        return agent

    async def start_workflow(self, workflow_spec: str) -> Dict:
        """Start a workflow with multiple agents"""
        results = {}

        try:
            # Parse workflow specification
            workflow_steps = self._parse_workflow_spec(workflow_spec)

            for step in workflow_steps:
                agent_name = step["agent"]
                task = step["task"]

                if agent_name in self.agents:
                    agent = self.agents[agent_name]
                    result = await agent.execute_task(task, step.get("context", {}))

                    results[agent_name] = {
                        "status": result.status.value,
                        "data": result.data,
                        "execution_time": result.execution_time
                    }

                    # Stop if any agent fails
                    if result.status == AgentStatus.FAILED:
                        break
                else:
                    results[agent_name] = {
                        "status": "failed",
                        "error": f"Agent {agent_name} not found"
                    }

            return results

        except Exception as e:
            return {"error": str(e)}

    def _parse_workflow_spec(self, spec: str) -> List[Dict]:
        """Parse workflow specification"""
        # Simple mock parser - in reality this would be more sophisticated
        return [
            {
                "agent": "coder",
                "task": spec,
                "context": {"priority": "high"}
            },
            {
                "agent": "tester",
                "task": "test the generated code",
                "context": {"test_type": "automated"}
            }
        ]

    async def get_orchestrator_status(self) -> Dict:
        """Get orchestrator status"""
        agent_statuses = {}

        for name, agent in self.agents.items():
            agent_statuses[name] = await agent.get_status()

        return {
            "total_agents": len(self.agents),
            "agents": agent_statuses,
            "backend_healthy": self.redis_client is not None
        }

# Factory function to create pre-configured orchestrator
def create_light_orchestrator() -> LightOrchestrator:
    """Create a pre-configured light orchestrator"""
    orchestrator = LightOrchestrator()

    # Register default agents
    coder = orchestrator.register_agent("coder")
    tester = orchestrator.register_agent("tester")
    fixer = orchestrator.register_agent("fixer")

    return orchestrator

# Example usage
async def run_light_agent_example():
    """Example of running the light agent system"""
    print("Starting Light Agent System...")

    # Create orchestrator
    orchestrator = create_light_orchestrator()

    # Check status
    status = await orchestrator.get_orchestrator_status()
    print(f"Orchestrator status: {status}")

    # Run a workflow
    workflow_spec = "Implement user authentication system"
    results = await orchestrator.start_workflow(workflow_spec)

    print("Workflow results:")
    for agent, result in results.items():
        print(f"  {agent}: {result['status']} ({result['execution_time']:.2f}s)")

    return results

if __name__ == "__main__":
    # Run example
    asyncio.run(run_light_agent_example())