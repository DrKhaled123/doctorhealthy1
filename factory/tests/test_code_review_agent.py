import pytest
import asyncio
from unittest.mock import Mock, patch, AsyncMock
import tempfile
import json
from factory.code_review_agent import CodeReviewAgent


class TestCodeReviewAgent:
    """Test cases for CodeReviewAgent"""

    @pytest.fixture
    def temp_config(self):
        """Create temporary config file"""
        config_data = {
            "review": {
                "max_issues": 50,
                "severity_levels": ["low", "medium", "high", "critical"],
                "auto_fix": True,
                "style_guide": "pep8"
            },
            "ai": {
                "model": "gpt-3.5-turbo",
                "temperature": 0.3,
                "max_tokens": 1000
            }
        }
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(config_data, f)
            temp_path = f.name
        yield temp_path

    @pytest.fixture
    def agent(self, temp_config):
        """Create CodeReviewAgent instance"""
        return CodeReviewAgent(temp_config)

    def test_initialization(self, agent):
        """Test agent initialization"""
        assert agent.config is not None
        assert agent.review_history == []
        assert agent.issues_found == 0

    @pytest.mark.asyncio
    async def test_review_code_basic(self, agent):
        """Test basic code review"""
        test_code = '''
def badFunction():
    x=1+2
    print(x)
    return x
'''

        result = await agent.review_code(test_code, "test.py")

        assert result["success"] is True
        assert "issues" in result
        assert "suggestions" in result
        assert "overall_score" in result
        assert len(agent.review_history) == 1

    @pytest.mark.asyncio
    async def test_review_code_perfect(self, agent):
        """Test review of perfect code"""
        perfect_code = '''"""Perfect code example."""

def calculate_fibonacci(n: int) -> int:
    """
    Calculate the nth Fibonacci number.

    Args:
        n: The position in the Fibonacci sequence.

    Returns:
        The nth Fibonacci number.

    Raises:
        ValueError: If n is negative.
    """
    if n < 0:
        raise ValueError("n must be non-negative")

    if n == 0:
        return 0
    elif n == 1:
        return 1
    else:
        return calculate_fibonacci(n - 1) + calculate_fibonacci(n - 2)


if __name__ == "__main__":
    # Test the function
    for i in range(10):
        print(f"Fibonacci({i}) = {calculate_fibonacci(i)}")
'''

        result = await agent.review_code(perfect_code, "fibonacci.py")

        assert result["success"] is True
        assert result["overall_score"] >= 90  # Should be very high score
        assert len(result["issues"]) == 0  # No issues expected

    @pytest.mark.asyncio
    async def test_review_code_with_syntax_error(self, agent):
        """Test review of code with syntax errors"""
        syntax_error_code = '''
def broken_function(
    print("This has a syntax error"
'''

        result = await agent.review_code(syntax_error_code, "broken.py")

        assert result["success"] is True
        assert len(result["issues"]) > 0
        assert any("syntax" in issue.lower() for issue in result["issues"])

    @pytest.mark.asyncio
    async def test_review_code_empty(self, agent):
        """Test review of empty code"""
        result = await agent.review_code("", "empty.py")

        assert result["success"] is True
        assert result["overall_score"] == 0
        assert len(result["issues"]) > 0

    @pytest.mark.asyncio
    async def test_review_code_whitespace_only(self, agent):
        """Test review of whitespace-only code"""
        result = await agent.review_code("   \n\t  \n  ", "whitespace.py")

        assert result["success"] is True
        assert result["overall_score"] == 0
        assert len(result["issues"]) > 0

    @pytest.mark.asyncio
    async def test_generate_suggestions(self, agent):
        """Test suggestion generation"""
        issues = [
            "Missing docstring",
            "Line too long",
            "Unused import"
        ]

        suggestions = await agent.generate_suggestions(issues)

        assert len(suggestions) == len(issues)
        assert all(isinstance(suggestion, dict) for suggestion in suggestions)
        assert all("fix" in suggestion for suggestion in suggestions)

    @pytest.mark.asyncio
    async def test_apply_auto_fixes(self, agent):
        """Test automatic fix application"""
        test_code = '''
def badFunction():
    x=1+2
    print(x)
    return x
'''

        fixed_code = await agent.apply_auto_fixes(test_code)

        assert isinstance(fixed_code, str)
        assert len(fixed_code) > 0
        # Should be different from original (fixed)
        assert fixed_code != test_code

    @pytest.mark.asyncio
    async def test_calculate_code_metrics(self, agent):
        """Test code metrics calculation"""
        test_code = '''
def example_function(param1, param2):
    """Example function."""
    if param1 > 0:
        for i in range(param2):
            print(f"Value: {i}")
    return param1 + param2

class ExampleClass:
    """Example class."""

    def __init__(self, value):
        self.value = value

    def get_value(self):
        return self.value
'''

        metrics = await agent.calculate_code_metrics(test_code)

        assert "lines_of_code" in metrics
        assert "cyclomatic_complexity" in metrics
        assert "function_count" in metrics
        assert "class_count" in metrics
        assert "maintainability_index" in metrics

        assert metrics["lines_of_code"] > 0
        assert metrics["function_count"] >= 1
        assert metrics["class_count"] >= 1

    @pytest.mark.asyncio
    async def test_check_pep8_compliance(self, agent):
        """Test PEP8 compliance checking"""
        pep8_code = '''"""PEP8 compliant code."""

def calculate_area(length: int, width: int) -> int:
    """
    Calculate rectangle area.

    Args:
        length: The length of the rectangle.
        width: The width of the rectangle.

    Returns:
        The area of the rectangle.
    """
    return length * width
'''

        compliance = await agent.check_pep8_compliance(pep8_code)

        assert compliance["is_compliant"] is True
        assert compliance["violation_count"] == 0

    @pytest.mark.asyncio
    async def test_check_pep8_non_compliant(self, agent):
        """Test PEP8 compliance checking for non-compliant code"""
        non_pep8_code = '''
def badFunction(): # missing spaces around parameter
    x=1+2# missing spaces
    print(x)
    return x
'''

        compliance = await agent.check_pep8_compliance(non_pep8_code)

        assert compliance["is_compliant"] is False
        assert compliance["violation_count"] > 0
        assert len(compliance["violations"]) > 0

    @pytest.mark.asyncio
    async def test_analyze_security_issues(self, agent):
        """Test security issue analysis"""
        insecure_code = '''
import os
import subprocess

def insecure_function(user_input):
    # Dangerous: Using eval
    result = eval(user_input)

    # Dangerous: Command injection
    os.system(f"echo {user_input}")

    # Dangerous: SQL injection vulnerability
    query = f"SELECT * FROM users WHERE name = '{user_input}'"

    return result
'''

        security_issues = await agent.analyze_security_issues(insecure_code)

        assert len(security_issues) > 0
        assert any("eval" in issue.lower() for issue in security_issues)
        assert any("injection" in issue.lower() for issue in security_issues)

    @pytest.mark.asyncio
    async def test_analyze_performance_issues(self, agent):
        """Test performance issue analysis"""
        inefficient_code = '''
def inefficient_function():
    # Inefficient: Multiple loops
    for i in range(1000):
        for j in range(1000):
            print(f"{i}, {j}")

    # Inefficient: Repeated calculations
    numbers = [1, 2, 3, 4, 5]
    for num in numbers:
        result = num * 2 + 1
        print(result)
'''

        performance_issues = await agent.analyze_performance_issues(inefficient_code)

        assert len(performance_issues) > 0
        assert any("nested" in issue.lower() or "loop" in issue.lower() for issue in performance_issues)

    @pytest.mark.asyncio
    async def test_generate_comprehensive_report(self, agent):
        """Test comprehensive report generation"""
        test_code = '''
def example():
    x=1+2
    print(x)
    return x
'''

        # First do a review
        review_result = await agent.review_code(test_code, "example.py")

        # Generate report
        report = await agent.generate_comprehensive_report("example.py", test_code)

        assert "file_name" in report
        assert "review_summary" in report
        assert "metrics" in report
        assert "security_analysis" in report
        assert "performance_analysis" in report
        assert "overall_assessment" in report

        assert report["file_name"] == "example.py"

    def test_review_history_tracking(self, agent):
        """Test review history tracking"""
        # Initially empty
        assert len(agent.review_history) == 0
        assert agent.issues_found == 0

        # Simulate some reviews
        agent.review_history = [
            {"file": "test1.py", "issues": 5},
            {"file": "test2.py", "issues": 3},
            {"file": "test3.py", "issues": 0}
        ]
        agent.issues_found = 8

        assert len(agent.review_history) == 3
        assert agent.issues_found == 8

    @pytest.mark.asyncio
    async def test_batch_review(self, agent):
        """Test batch code review"""
        code_files = {
            "file1.py": 'def func1(): pass',
            "file2.py": 'def func2(): pass',
            "file3.py": 'def func3(): pass'
        }

        results = await agent.batch_review(code_files)

        assert len(results) == 3
        assert all(result["success"] is True for result in results.values())
        assert len(agent.review_history) == 3

    @pytest.mark.asyncio
    async def test_review_with_custom_rules(self, agent):
        """Test review with custom rules"""
        custom_rules = [
            "No print statements allowed",
            "All functions must have docstrings",
            "Maximum line length is 100 characters"
        ]

        test_code = '''
def bad_function():
    print("This should trigger custom rules")
    return "test"
'''

        result = await agent.review_code(test_code, "test.py", custom_rules)

        assert result["success"] is True
        assert "issues" in result
        # Should detect at least some of the custom rule violations
        assert len(result["issues"]) > 0

    def test_config_validation(self, agent):
        """Test configuration validation"""
        # Test valid config structure
        assert "review" in agent.config
        assert "ai" in agent.config
        assert "max_issues" in agent.config["review"]
        assert "model" in agent.config["ai"]

        # Test default values
        assert agent.config["review"]["auto_fix"] is True
        assert agent.config["ai"]["temperature"] == 0.3

    @pytest.mark.asyncio
    async def test_error_handling_invalid_code(self, agent):
        """Test error handling for invalid code"""
        # Test with None
        result = await agent.review_code(None, "test.py")
        assert result["success"] is False
        assert "error" in result

        # Test with non-string
        result = await agent.review_code(123, "test.py")
        assert result["success"] is False
        assert "error" in result

    @pytest.mark.asyncio
    async def test_memory_usage_tracking(self, agent):
        """Test memory usage tracking during reviews"""
        initial_history_length = len(agent.review_history)

        # Perform multiple reviews
        for i in range(5):
            await agent.review_code(f'def func{i}(): pass', f'test{i}.py')

        # Check that history is tracked
        assert len(agent.review_history) == initial_history_length + 5
        assert agent.issues_found >= 0  # Should be non-negative

    @pytest.mark.asyncio
    async def test_concurrent_reviews(self, agent):
        """Test handling of concurrent reviews"""
        import asyncio

        async def review_code_async(code, filename):
            return await agent.review_code(code, filename)

        # Create multiple review tasks
        tasks = [
            review_code_async(f'def func{i}(): pass', f'test{i}.py')
            for i in range(10)
        ]

        # Execute concurrently
        results = await asyncio.gather(*tasks)

        # All should succeed
        assert len(results) == 10
        assert all(result["success"] is True for result in results)
        assert len(agent.review_history) == 10
