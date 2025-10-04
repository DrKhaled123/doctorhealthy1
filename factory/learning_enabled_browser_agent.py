#!/usr/bin/env python3
"""
Learning-Enabled Browser Testing Agent
An enhanced browser testing agent that participates in inter-agent learning
"""

import asyncio
import time
import json
import uuid
from typing import Dict, List, Any, Optional
from datetime import datetime

from playwright.async_api import async_playwright
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

# Import the learning system
from inter_agent_learning_system import (
    LearningEnabledAgent,
    AgentCapability,
    LearningType
)

class LearningEnabledBrowserTestingAgent(LearningEnabledAgent):
    """Browser testing agent with inter-agent learning capabilities"""

    def __init__(self, agent_id: str = None):
        agent_id = agent_id or f"browser_agent_{int(time.time())}"
        super().__init__(
            agent_id=agent_id,
            agent_type="BrowserTestingAgent",
            capabilities=[
                AgentCapability.TESTING,
                AgentCapability.MONITORING,
                AgentCapability.DEBUGGING
            ]
        )

        # Browser testing components
        self.playwright = None
        self.browser = None
        self.page = None
        self.selenium_driver = None

        # Learning tracking
        self.current_test_session = None
        self.test_history = []
        self.performance_metrics = {}

    async def initialize(self):
        """Initialize the learning-enabled browser agent"""
        await self.initialize_learning()
        logger.info(f"🎭 Learning-Enabled Browser Agent {self.agent_id} initialized")

    async def run_test_session(self, url: str, test_cases: List[Dict],
                             session_context: Dict[str, Any] = None) -> Dict[str, Any]:
        """Run a test session with learning integration"""
        session_id = str(uuid.uuid4())
        start_time = time.time()

        self.current_test_session = {
            'session_id': session_id,
            'url': url,
            'test_cases_count': len(test_cases),
            'start_time': start_time,
            'context': session_context or {}
        }

        try:
            # Run the actual tests
            results = await self.test_web_application(url, test_cases)

            # Calculate session metrics
            session_duration = time.time() - start_time
            success_rate = results['passed'] / max(results['passed'] + results['failed'], 1)

            # Share learning experience
            await self._share_test_experience(results, session_context)

            # Update performance metrics
            self.performance_metrics[session_id] = {
                'duration': session_duration,
                'success_rate': success_rate,
                'tests_run': len(test_cases),
                'timestamp': datetime.now().isoformat()
            }

            # Store session results
            self.test_history.append({
                'session_id': session_id,
                'results': results,
                'duration': session_duration,
                'timestamp': datetime.now().isoformat()
            })

            return {
                'session_id': session_id,
                'results': results,
                'duration': session_duration,
                'learning_shared': True
            }

        except Exception as e:
            # Share failure experience
            await self._share_failure_experience(str(e), session_context)

            return {
                'session_id': session_id,
                'error': str(e),
                'duration': time.time() - start_time,
                'learning_shared': True
            }

    async def _share_test_experience(self, results: Dict[str, Any],
                                   context: Dict[str, Any]):
        """Share successful testing experience with other agents"""
        try:
            success = results.get('passed', 0) > results.get('failed', 0)
            confidence = min(1.0, results.get('passed', 0) / max(results.get('passed', 0) + results.get('failed', 0), 1))

            lessons_learned = []

            if success:
                lessons_learned.extend([
                    'Browser testing session completed successfully',
                    'Test stability improved with proper wait conditions',
                    'Page load optimization enhances test reliability'
                ])

                if results.get('total_time', 0) < 30:  # Fast execution
                    lessons_learned.append('Fast test execution indicates good page performance')

            else:
                lessons_learned.extend([
                    'Test failures often indicate timing or element issues',
                    'Retry mechanisms can help with flaky tests',
                    'Proper error handling is crucial for test reliability'
                ])

            await self.share_experience(
                capability=AgentCapability.TESTING,
                context={
                    'task_type': 'browser_testing',
                    'test_count': len(results.get('tests', [])),
                    'url': self.current_test_session.get('url', ''),
                    **context
                },
                outcome={
                    'success': success,
                    'success_patterns': ['proper_wait_conditions', 'error_handling'],
                    'failure_patterns': ['timing_issues', 'element_not_found'],
                    'performance_metrics': {
                        'total_time': results.get('total_time', 0),
                        'success_rate': confidence
                    }
                },
                success=success,
                confidence=confidence,
                lessons_learned=lessons_learned
            )

        except Exception as e:
            logger.error(f"❌ Failed to share test experience: {e}")

    async def _share_failure_experience(self, error: str, context: Dict[str, Any]):
        """Share failure experience for learning"""
        try:
            await self.share_experience(
                capability=AgentCapability.DEBUGGING,
                context={
                    'task_type': 'error_analysis',
                    'error_type': 'browser_testing_failure',
                    **context
                },
                outcome={
                    'success': False,
                    'error_patterns': [error[:100]],  # Truncate long errors
                    'debugging_insights': ['Need better error handling', 'Consider retry mechanisms']
                },
                success=False,
                confidence=0.3,  # Low confidence for failures
                lessons_learned=[
                    'Browser testing failures often have multiple causes',
                    'Network issues can cause intermittent failures',
                    'Element selectors need regular maintenance'
                ]
            )

        except Exception as e:
            logger.error(f"❌ Failed to share failure experience: {e}")

    async def learn_from_other_agents(self):
        """Learn from experiences of other agents"""
        try:
            # Get recommendations for testing capability
            recommendations = await self.get_learning_recommendations(AgentCapability.TESTING)

            if recommendations:
                print(f"📚 Learning recommendations for {self.agent_id}:")
                for rec in recommendations:
                    print(f"  • {rec['type']}: {rec['reason']}")

                    # Request knowledge transfer if recommended
                    if rec['type'] == 'learn_from_expert':
                        transfer_id = await self.request_knowledge(
                            AgentCapability.TESTING,
                            rec['recommended_agent']
                        )
                        if transfer_id:
                            print(f"    🔄 Knowledge transfer initiated: {transfer_id}")

            # Get collaboration opportunities
            opportunities = await self.get_learning_recommendations()
            if opportunities:
                print(f"🤝 Collaboration opportunities for {self.agent_id}:")
                for opp in opportunities[:3]:  # Top 3
                    print(f"  • {opp['agent_type']}: {opp['reason']}")

        except Exception as e:
            logger.error(f"❌ Failed to learn from other agents: {e}")

    async def test_web_application(self, url: str, test_cases: List[Dict]) -> Dict:
        """Enhanced web application testing with learning integration"""
        results = {
            "url": url,
            "tests": [],
            "passed": 0,
            "failed": 0,
            "total_time": 0,
            "learning_insights": []
        }

        start_time = time.time()

        try:
            # Try Playwright first (modern async testing)
            await self.setup_playwright()

            for test_case in test_cases:
                test_result = await self._run_playwright_test(test_case)
                results["tests"].append(test_result)

                if test_result["passed"]:
                    results["passed"] += 1
                else:
                    results["failed"] += 1

                # Add learning insights
                if test_result["passed"]:
                    results["learning_insights"].append(
                        f"Test '{test_case.get('name', 'Unknown')}' passed - good test design"
                    )
                else:
                    results["learning_insights"].append(
                        f"Test '{test_case.get('name', 'Unknown')}' failed: {test_result.get('error', 'Unknown error')}"
                    )

            await self.cleanup_playwright()

        except Exception as playwright_error:
            logger.warning(f"⚠️ Playwright testing failed: {playwright_error}")

            # Fallback to Selenium
            try:
                self.setup_selenium_fallback()

                for test_case in test_cases:
                    test_result = self._run_selenium_test(test_case)
                    results["tests"].append(test_result)

                    if test_result["passed"]:
                        results["passed"] += 1
                    else:
                        results["failed"] += 1

                self.cleanup_selenium()

            except Exception as selenium_error:
                logger.error(f"❌ Both Playwright and Selenium failed: {playwright_error}, {selenium_error}")
                results["error"] = f"Both testing frameworks failed: {playwright_error}, {selenium_error}"

        results["total_time"] = time.time() - start_time
        return results

    async def setup_playwright(self):
        """Setup Playwright for modern async testing"""
        try:
            self.playwright = await async_playwright().start()
            self.browser = await self.playwright.chromium.launch(
                headless=True,
                args=['--no-sandbox', '--disable-dev-shm-usage']
            )
            self.page = await self.browser.new_page()
            print("✅ Playwright setup successful")
        except Exception as e:
            print(f"❌ Playwright setup failed: {e}")
            raise

    def setup_selenium_fallback(self):
        """Setup Selenium for legacy compatibility"""
        try:
            chrome_options = Options()
            chrome_options.add_argument('--headless')
            chrome_options.add_argument('--no-sandbox')
            chrome_options.add_argument('--disable-dev-shm-usage')
            self.selenium_driver = webdriver.Chrome(options=chrome_options)
            print("✅ Selenium setup successful")
        except Exception as e:
            print(f"❌ Selenium setup failed: {e}")
            raise

    async def _run_playwright_test(self, test_case: Dict) -> Dict:
        """Run a single test case with Playwright"""
        test_result = {
            "name": test_case.get("name", "Unnamed test"),
            "passed": False,
            "error": None,
            "duration": 0,
            "screenshots": []
        }

        start_time = time.time()

        try:
            if "url" in test_case:
                await self.page.goto(test_case["url"])
            elif hasattr(self, 'base_url'):
                await self.page.goto(self.base_url)

            # Wait for page to load
            await self.page.wait_for_load_state("networkidle")

            # Execute test actions
            if "actions" in test_case:
                for action in test_case["actions"]:
                    await self._execute_playwright_action(action)

            # Check assertions
            if "assertions" in test_case:
                for assertion in test_case["assertions"]:
                    if not await self._check_playwright_assertion(assertion):
                        test_result["error"] = f"Assertion failed: {assertion}"
                        return test_result

            # Take screenshot if requested
            if test_case.get("screenshot", False):
                screenshot_path = f"screenshot_{test_case['name']}_{int(time.time())}.png"
                await self.page.screenshot(path=screenshot_path)
                test_result["screenshots"].append(screenshot_path)

            test_result["passed"] = True

        except Exception as e:
            test_result["error"] = str(e)

        finally:
            test_result["duration"] = time.time() - start_time

        return test_result

    def _run_selenium_test(self, test_case: Dict) -> Dict:
        """Run a single test case with Selenium"""
        test_result = {
            "name": test_case.get("name", "Unnamed test"),
            "passed": False,
            "error": None,
            "duration": 0,
            "screenshots": []
        }

        start_time = time.time()

        try:
            if "url" in test_case:
                self.selenium_driver.get(test_case["url"])
            elif hasattr(self, 'base_url'):
                self.selenium_driver.get(self.base_url)

            # Wait for page to load
            WebDriverWait(self.selenium_driver, 10).until(
                lambda driver: driver.execute_script("return document.readyState") == "complete"
            )

            # Execute test actions
            if "actions" in test_case:
                for action in test_case["actions"]:
                    self._execute_selenium_action(action)

            # Check assertions
            if "assertions" in test_case:
                for assertion in test_case["assertions"]:
                    if not self._check_selenium_assertion(assertion):
                        test_result["error"] = f"Assertion failed: {assertion}"
                        return test_result

            # Take screenshot if requested
            if test_case.get("screenshot", False):
                screenshot_path = f"selenium_screenshot_{test_case['name']}_{int(time.time())}.png"
                self.selenium_driver.save_screenshot(screenshot_path)
                test_result["screenshots"].append(screenshot_path)

            test_result["passed"] = True

        except Exception as e:
            test_result["error"] = str(e)

        finally:
            test_result["duration"] = time.time() - start_time

        return test_result

    async def _execute_playwright_action(self, action: Dict):
        """Execute a Playwright action"""
        action_type = action.get("type")

        if action_type == "click":
            await self.page.click(action["selector"])
        elif action_type == "type":
            await self.page.fill(action["selector"], action["text"])
        elif action_type == "wait":
            await asyncio.sleep(action.get("seconds", 1))
        elif action_type == "scroll":
            await self.page.evaluate("window.scrollTo(0, document.body.scrollHeight)")

    def _execute_selenium_action(self, action: Dict):
        """Execute a Selenium action"""
        action_type = action.get("type")

        if action_type == "click":
            element = WebDriverWait(self.selenium_driver, 10).until(
                EC.element_to_be_clickable((By.CSS_SELECTOR, action["selector"]))
            )
            element.click()
        elif action_type == "type":
            element = WebDriverWait(self.selenium_driver, 10).until(
                EC.presence_of_element_located((By.CSS_SELECTOR, action["selector"]))
            )
            element.clear()
            element.send_keys(action["text"])
        elif action_type == "wait":
            time.sleep(action.get("seconds", 1))
        elif action_type == "scroll":
            self.selenium_driver.execute_script("window.scrollTo(0, document.body.scrollHeight)")

    async def _check_playwright_assertion(self, assertion: Dict) -> bool:
        """Check a Playwright assertion"""
        assertion_type = assertion.get("type")

        if assertion_type == "element_exists":
            try:
                await self.page.wait_for_selector(assertion["selector"], timeout=5000)
                return True
            except:
                return False
        elif assertion_type == "text_contains":
            content = await self.page.content()
            return assertion["text"] in content
        elif assertion_type == "url_contains":
            current_url = self.page.url
            return assertion["text"] in current_url

        return False

    def _check_selenium_assertion(self, assertion: Dict) -> bool:
        """Check a Selenium assertion"""
        assertion_type = assertion.get("type")

        if assertion_type == "element_exists":
            try:
                WebDriverWait(self.selenium_driver, 5).until(
                    EC.presence_of_element_located((By.CSS_SELECTOR, assertion["selector"]))
                )
                return True
            except:
                return False
        elif assertion_type == "text_contains":
            return assertion["text"] in self.selenium_driver.page_source
        elif assertion_type == "url_contains":
            return assertion["text"] in self.selenium_driver.current_url

        return False

    async def cleanup_playwright(self):
        """Clean up Playwright resources"""
        if self.page:
            await self.page.close()
        if self.browser:
            await self.browser.close()
        if self.playwright:
            await self.playwright.stop()

    def cleanup_selenium(self):
        """Clean up Selenium resources"""
        if self.selenium_driver:
            self.selenium_driver.quit()

    async def generate_test_cases(self, app_description: str) -> List[Dict]:
        """Generate test cases based on application description with learning insights"""
        # Get learning insights from other agents
        recommendations = await self.get_learning_recommendations(AgentCapability.TESTING)

        base_test_cases = [
            {
                "name": "Load homepage",
                "actions": [
                    {"type": "wait", "seconds": 2}
                ],
                "assertions": [
                    {"type": "element_exists", "selector": "body"}
                ],
                "screenshot": True
            },
            {
                "name": "Check title",
                "actions": [
                    {"type": "wait", "seconds": 1}
                ],
                "assertions": [
                    {"type": "text_contains", "text": "Welcome"}
                ]
            }
        ]

        # Enhance test cases based on learning recommendations
        if recommendations:
            for rec in recommendations:
                if "pattern_based" in rec['type']:
                    # Add pattern-based test case
                    base_test_cases.append({
                        "name": f"Pattern-based test: {rec['reason'][:30]}",
                        "actions": [
                            {"type": "wait", "seconds": 1}
                        ],
                        "assertions": [
                            {"type": "element_exists", "selector": "body"}
                        ]
                    })

        return base_test_cases

# Example usage
async def demonstrate_learning_browser_agent():
    """Demonstrate the learning-enabled browser testing agent"""
    print("🎭 Learning-Enabled Browser Testing Agent Demo")
    print("=" * 50)

    try:
        # Create and initialize the agent
        agent = LearningEnabledBrowserTestingAgent()
        await agent.initialize()

        # Example test cases
        test_cases = [
            {
                "name": "Basic page load",
                "url": "http://localhost:3000",
                "actions": [
                    {"type": "wait", "seconds": 2}
                ],
                "assertions": [
                    {"type": "element_exists", "selector": "h1"},
                    {"type": "text_contains", "text": "Welcome"}
                ],
                "screenshot": True
            }
        ]

        # Run test session with learning
        print("\n🧪 Running test session with learning...")
        results = await agent.run_test_session(
            "http://localhost:3000",
            test_cases,
            {"test_context": "demo_session", "priority": "high"}
        )

        print("✅ Test session completed:")
        print(f"  • Session ID: {results['session_id']}")
        print(f"  • Duration: {results['duration']:.2f}s")
        print(f"  • Learning shared: {results['learning_shared']}")

        # Learn from other agents
        print("\n📚 Learning from other agents...")
        await agent.learn_from_other_agents()

        # Generate enhanced test cases
        print("\n🔍 Generating enhanced test cases...")
        enhanced_tests = await agent.generate_test_cases("E-commerce website")
        print(f"✅ Generated {len(enhanced_tests)} enhanced test cases")

        for test in enhanced_tests:
            print(f"  • {test['name']}")

        print("\n🎉 Learning-Enabled Browser Agent demonstration completed!")
        return True

    except Exception as e:
        print(f"❌ Demonstration failed: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    # Run demonstration
    asyncio.run(demonstrate_learning_browser_agent())
