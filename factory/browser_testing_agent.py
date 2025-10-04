from playwright.async_api import async_playwright
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import asyncio
import json
import time
from typing import Dict, List, Optional
import os

class BrowserTestingAgent:
    """Handles all browser-based testing with Playwright and Selenium fallback"""

    def __init__(self):
        self.playwright = None
        self.browser = None
        self.page = None
        self.selenium_driver = None

    async def setup_playwright(self):
        """Setup Playwright for modern async testing"""
        try:
            self.playwright = await async_playwright().start()
            self.browser = await self.playwright.chromium.launch(
                headless=True,
                args=['--no-sandbox', '--disable-dev-shm-usage']
            )
            self.page = await self.browser.new_page()
            print("Playwright setup successful")
        except Exception as e:
            print(f"Playwright setup failed: {e}")
            raise

    def setup_selenium_fallback(self):
        """Setup Selenium for legacy compatibility"""
        try:
            chrome_options = Options()
            chrome_options.add_argument('--headless')
            chrome_options.add_argument('--no-sandbox')
            chrome_options.add_argument('--disable-dev-shm-usage')
            self.selenium_driver = webdriver.Chrome(options=chrome_options)
            print("Selenium setup successful")
        except Exception as e:
            print(f"Selenium setup failed: {e}")
            raise

    async def test_web_application(self, url: str, test_cases: List[Dict]) -> Dict:
        """Run comprehensive web application tests"""
        results = {
            "url": url,
            "tests": [],
            "passed": 0,
            "failed": 0,
            "total_time": 0
        }

        start_time = time.time()

        try:
            # Try Playwright first
            await self.setup_playwright()

            for test_case in test_cases:
                test_result = await self._run_playwright_test(test_case)
                results["tests"].append(test_result)

                if test_result["passed"]:
                    results["passed"] += 1
                else:
                    results["failed"] += 1

            await self.cleanup_playwright()

        except Exception as playwright_error:
            print(f"Playwright testing failed: {playwright_error}")
            print("Falling back to Selenium...")

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
                print(f"Selenium testing also failed: {selenium_error}")
                results["error"] = f"Both Playwright and Selenium failed: {playwright_error}, {selenium_error}"

        results["total_time"] = time.time() - start_time
        return results

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
        """Generate test cases based on application description"""
        # Mock test case generation
        test_cases = [
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

        return test_cases

# Example usage
async def run_browser_tests():
    """Example function to run browser tests"""
    agent = BrowserTestingAgent()

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

    try:
        results = await agent.test_web_application("http://localhost:3000", test_cases)

        # Print results
        print(f"Tests completed: {results['passed']} passed, {results['failed']} failed")
        print(f"Total time: {results['total_time']:.2f}s")

        for test in results["tests"]:
            status = "PASS" if test["passed"] else "FAIL"
            print(f"  {test['name']}: {status} ({test['duration']:.2f}s)")

        return results

    except Exception as e:
        print(f"Browser testing failed: {e}")
        return {"error": str(e)}

if __name__ == "__main__":
    # Run example tests
    asyncio.run(run_browser_tests())