import pytest
import asyncio
from unittest.mock import Mock, patch, AsyncMock
import tempfile
import json
from factory.browser_testing_agent import BrowserTestingAgent


class TestBrowserTestingAgent:
    """Test cases for BrowserTestingAgent"""

    @pytest.fixture
    def temp_config(self):
        """Create temporary config file"""
        config_data = {
            "browser": {
                "headless": True,
                "timeout": 30000,
                "viewport": {"width": 1280, "height": 720}
            },
            "test": {
                "base_url": "http://localhost:3000",
                "wait_time": 5000
            }
        }
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(config_data, f)
            temp_path = f.name
        yield temp_path

    @pytest.fixture
    def agent(self, temp_config):
        """Create BrowserTestingAgent instance"""
        return BrowserTestingAgent(temp_config)

    def test_initialization(self, agent):
        """Test agent initialization"""
        assert agent.config is not None
        assert agent.browser is None
        assert agent.page is None

    @pytest.mark.asyncio
    async def test_initialize_browser_success(self, agent):
        """Test successful browser initialization"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_pw_instance.chromium.launch.return_value = AsyncMock()
            mock_pw_instance.chromium.launch.return_value.new_page.return_value = AsyncMock()

            await agent.initialize_browser()

            assert agent.browser is not None
            assert agent.page is not None

    @pytest.mark.asyncio
    async def test_initialize_browser_failure(self, agent):
        """Test browser initialization failure"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_pw_instance.chromium.launch.side_effect = Exception("Browser launch failed")

            with pytest.raises(Exception):
                await agent.initialize_browser()

    @pytest.mark.asyncio
    async def test_navigate_to_page_success(self, agent):
        """Test successful page navigation"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful navigation
            mock_page.goto.return_value = asyncio.Future()
            mock_page.goto.return_value.set_result(None)

            result = await agent.navigate_to_page("http://example.com")

            assert result is True
            mock_page.goto.assert_called_once_with("http://example.com")

    @pytest.mark.asyncio
    async def test_navigate_to_page_failure(self, agent):
        """Test page navigation failure"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock failed navigation
            mock_page.goto.return_value = asyncio.Future()
            mock_page.goto.return_value.set_exception(Exception("Navigation failed"))

            result = await agent.navigate_to_page("http://invalid-url")

            assert result is False

    @pytest.mark.asyncio
    async def test_wait_for_element_success(self, agent):
        """Test successful element waiting"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful element wait
            mock_element = AsyncMock()
            mock_page.wait_for_selector.return_value = asyncio.Future()
            mock_page.wait_for_selector.return_value.set_result(mock_element)

            result = await agent.wait_for_element("#test-element")

            assert result is True
            mock_page.wait_for_selector.assert_called_once_with("#test-element")

    @pytest.mark.asyncio
    async def test_wait_for_element_timeout(self, agent):
        """Test element waiting timeout"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock timeout
            mock_page.wait_for_selector.return_value = asyncio.Future()
            mock_page.wait_for_selector.return_value.set_exception(
                Exception("Timeout waiting for element")
            )

            result = await agent.wait_for_element("#missing-element")

            assert result is False

    @pytest.mark.asyncio
    async def test_click_element_success(self, agent):
        """Test successful element clicking"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_element = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful click
            mock_element.click.return_value = asyncio.Future()
            mock_element.click.return_value.set_result(None)
            mock_page.query_selector.return_value = mock_element

            result = await agent.click_element("#clickable")

            assert result is True
            mock_element.click.assert_called_once()

    @pytest.mark.asyncio
    async def test_click_element_not_found(self, agent):
        """Test clicking non-existent element"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock element not found
            mock_page.query_selector.return_value = None

            result = await agent.click_element("#nonexistent")

            assert result is False

    @pytest.mark.asyncio
    async def test_type_text_success(self, agent):
        """Test successful text typing"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_element = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful typing
            mock_element.type.return_value = asyncio.Future()
            mock_element.type.return_value.set_result(None)
            mock_page.query_selector.return_value = mock_element

            result = await agent.type_text("#input", "test text")

            assert result is True
            mock_element.type.assert_called_once_with("test text")

    @pytest.mark.asyncio
    async def test_get_element_text_success(self, agent):
        """Test successful text retrieval"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_element = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful text retrieval
            mock_element.text_content.return_value = asyncio.Future()
            mock_element.text_content.return_value.set_result("Retrieved text")
            mock_page.query_selector.return_value = mock_element

            result = await agent.get_element_text("#text-element")

            assert result == "Retrieved text"
            mock_element.text_content.assert_called_once()

    @pytest.mark.asyncio
    async def test_get_element_text_not_found(self, agent):
        """Test text retrieval for non-existent element"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock element not found
            mock_page.query_selector.return_value = None

            result = await agent.get_element_text("#nonexistent")

            assert result is None

    @pytest.mark.asyncio
    async def test_take_screenshot_success(self, agent):
        """Test successful screenshot capture"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful screenshot
            mock_page.screenshot.return_value = asyncio.Future()
            mock_page.screenshot.return_value.set_result(b"screenshot_data")

            result = await agent.take_screenshot("test_screenshot.png")

            assert result is True
            mock_page.screenshot.assert_called_once()

    @pytest.mark.asyncio
    async def test_run_basic_test_success(self, agent):
        """Test successful basic test execution"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_element = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful test steps
            mock_page.goto.return_value = asyncio.Future()
            mock_page.goto.return_value.set_result(None)

            mock_page.wait_for_selector.return_value = asyncio.Future()
            mock_page.wait_for_selector.return_value.set_result(mock_element)

            mock_element.click.return_value = asyncio.Future()
            mock_element.click.return_value.set_result(None)
            mock_page.query_selector.return_value = mock_element

            test_steps = [
                {"action": "navigate", "url": "http://example.com"},
                {"action": "wait", "selector": "#button"},
                {"action": "click", "selector": "#button"}
            ]

            result = await agent.run_basic_test(test_steps)

            assert result["success"] is True
            assert result["steps_completed"] == 3
            assert len(result["errors"]) == 0

    @pytest.mark.asyncio
    async def test_run_basic_test_with_failures(self, agent):
        """Test basic test execution with failures"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock failed navigation
            mock_page.goto.return_value = asyncio.Future()
            mock_page.goto.return_value.set_exception(Exception("Navigation failed"))

            test_steps = [
                {"action": "navigate", "url": "http://invalid-url"},
                {"action": "wait", "selector": "#button"}
            ]

            result = await agent.run_basic_test(test_steps)

            assert result["success"] is False
            assert result["steps_completed"] == 0
            assert len(result["errors"]) > 0
            assert "Navigation failed" in result["errors"][0]

    @pytest.mark.asyncio
    async def test_cleanup(self, agent):
        """Test browser cleanup"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock cleanup methods
            mock_page.close.return_value = asyncio.Future()
            mock_page.close.return_value.set_result(None)
            mock_browser.close.return_value = asyncio.Future()
            mock_browser.close.return_value.set_result(None)

            await agent.cleanup()

            assert agent.browser is None
            assert agent.page is None
            mock_page.close.assert_called_once()
            mock_browser.close.assert_called_once()

    @pytest.mark.asyncio
    async def test_execute_test_suite(self, agent):
        """Test test suite execution"""
        with patch('playwright.async_api.Playwright') as mock_pw:
            # Setup mocks
            mock_pw_instance = AsyncMock()
            mock_pw.return_value = mock_pw_instance
            mock_browser = AsyncMock()
            mock_page = AsyncMock()
            mock_pw_instance.chromium.launch.return_value = mock_browser
            mock_browser.new_page.return_value = mock_page

            await agent.initialize_browser()

            # Mock successful test execution
            mock_page.goto.return_value = asyncio.Future()
            mock_page.goto.return_value.set_result(None)

            test_suite = {
                "name": "Login Test Suite",
                "tests": [
                    {
                        "name": "Valid Login",
                        "steps": [
                            {"action": "navigate", "url": "http://example.com"},
                            {"action": "wait", "selector": "#username"},
                            {"action": "type", "selector": "#username", "text": "testuser"},
                            {"action": "click", "selector": "#login-btn"}
                        ]
                    }
                ]
            }

            result = await agent.execute_test_suite(test_suite)

            assert result["success"] is True
            assert result["suite_name"] == "Login Test Suite"
            assert len(result["test_results"]) == 1
            assert result["test_results"][0]["success"] is True

    def test_config_validation(self, agent):
        """Test configuration validation"""
        # Test valid config
        assert agent.config["browser"]["headless"] is True
        assert agent.config["test"]["base_url"] == "http://localhost:3000"

        # Test missing keys have defaults
        assert "timeout" in agent.config["browser"]
        assert "wait_time" in agent.config["test"]
