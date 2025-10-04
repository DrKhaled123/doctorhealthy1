#!/usr/bin/env python3
"""
Complete Factory Orchestrator System Demo
Tests all components: orchestrator, memory, browser testing, and autofix
"""

import asyncio
import json
import time
from factory_orchestrator import FactoryOrchestrator
from memo_ai_memory import MemoAIMemory
from browser_testing_agent import BrowserTestingAgent
from autofix_agent import AutofixAgent

async def test_memory_system():
    """Test the memory system"""
    print("🧠 Testing Memory System...")
    memory = MemoAIMemory()

    # Store some experiences
    await memory.store_success_pattern(
        "user_authentication",
        {"method": "jwt", "secure": True},
        "successful_implementation"
    )

    await memory.store_failure_pattern(
        "database_connection",
        "connection_timeout",
        {"host": "localhost", "timeout": 30}
    )

    # Recall similar experiences
    similar = await memory.recall("user login system", k=3)
    print(f"✅ Found {len(similar)} similar experiences")

    # Get stats
    stats = await memory.get_memory_stats()
    print(f"✅ Memory stats: {stats}")

    return True

async def test_autofix_agent():
    """Test the autofix agent"""
    print("🔧 Testing Autofix Agent...")

    # Sample problematic code
    problematic_code = '''
def badFunction(name,age
    if name == "admin":
    password = "secret123"
        return True
    return False

x=1+2
'''

    agent = AutofixAgent()
    result = await agent.fix_code(problematic_code)

    print(f"✅ Fixes applied: {len(result['fixes_applied'])}")
    print(f"✅ Success: {result['success']}")

    if result['errors_remaining']:
        print(f"⚠️  Errors remaining: {len(result['errors_remaining'])}")

    return result['success']

async def test_browser_agent():
    """Test the browser testing agent"""
    print("🌐 Testing Browser Agent...")

    agent = BrowserTestingAgent()

    # Generate test cases
    test_cases = await agent.generate_test_cases("web application")
    print(f"✅ Generated {len(test_cases)} test cases")

    # Note: Actual browser testing would require a running web server
    print("✅ Browser agent ready (requires web server for full testing)")

    return True

async def test_factory_orchestrator():
    """Test the main factory orchestrator"""
    print("🏭 Testing Factory Orchestrator...")

    orchestrator = FactoryOrchestrator()

    # Test initialization
    await orchestrator.initialize()
    print("✅ Orchestrator initialized")

    # Test a simple pipeline run
    spec = "Create a simple calculator function"
    result = await orchestrator.run_pipeline(spec)

    print(f"✅ Pipeline completed: {result.success}")
    print(f"✅ Final stage: {result.stage.value}")

    return result.success

async def run_complete_demo():
    """Run the complete system demo"""
    print("🚀 Factory Orchestrator Complete System Demo")
    print("=============================================")

    start_time = time.time()

    tests = [
        ("Memory System", test_memory_system),
        ("Autofix Agent", test_autofix_agent),
        ("Browser Agent", test_browser_agent),
        ("Factory Orchestrator", test_factory_orchestrator),
    ]

    results = {}

    for test_name, test_func in tests:
        try:
            print(f"\n🔄 Running {test_name}...")
            success = await test_func()
            results[test_name] = success

            if success:
                print(f"✅ {test_name}: PASSED")
            else:
                print(f"❌ {test_name}: FAILED")

        except Exception as e:
            print(f"❌ {test_name}: ERROR - {e}")
            results[test_name] = False

    # Summary
    end_time = time.time()
    duration = end_time - start_time

    print("\n=============================================")
    print("📊 Demo Summary:")
    print(f"⏱️  Total time: {duration:.2f}s")

    passed = sum(1 for result in results.values() if result)
    total = len(results)

    print(f"✅ Passed: {passed}/{total}")

    for test_name, success in results.items():
        status = "✅ PASS" if success else "❌ FAIL"
        print(f"  {status}: {test_name}")

    if passed == total:
        print("\n🎉 All tests passed! Factory system is ready.")
        print("\n📖 Next steps:")
        print("  1. Start Redis: brew services start redis")
        print("  2. Run factory: ./start_factory.sh")
        print("  3. Use orchestrator: python3 factory_orchestrator.py 'your spec'")
    else:
        print(f"\n⚠️  {total - passed} test(s) failed. Check the errors above.")

    return passed == total

if __name__ == "__main__":
    success = asyncio.run(run_complete_demo())
    exit(0 if success else 1)