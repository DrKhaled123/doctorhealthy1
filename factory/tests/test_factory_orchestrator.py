import pytest
import asyncio
import json
import tempfile
import os
from unittest.mock import Mock, patch, AsyncMock
from factory.factory_orchestrator import (
    FactoryOrchestrator,
    PipelineStage,
    PipelineResult,
    run_factory
)


class TestFactoryOrchestrator:
    """Test cases for FactoryOrchestrator"""

    @pytest.fixture
    def temp_config(self):
        """Create temporary config file"""
        config_data = {
            "redis": {"host": "localhost", "port": 6379, "db": 0},
            "debug": False,
            "memory": {"max_memories": 1000}
        }
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(config_data, f)
            temp_path = f.name
        yield temp_path
        os.unlink(temp_path)

    @pytest.fixture
    def orchestrator(self, temp_config):
        """Create FactoryOrchestrator instance"""
        return FactoryOrchestrator(temp_config)

    def test_initialization(self, orchestrator):
        """Test orchestrator initialization"""
        assert orchestrator.config is not None
        assert orchestrator.redis_client is None
        assert orchestrator.memory_store is None
        assert orchestrator.current_stage is None

    def test_config_loading(self, temp_config):
        """Test configuration loading"""
        orchestrator = FactoryOrchestrator(temp_config)
        assert orchestrator.config["redis"]["host"] == "localhost"
        assert orchestrator.config["debug"] is False

    def test_config_loading_missing_file(self):
        """Test configuration loading with missing file"""
        orchestrator = FactoryOrchestrator("nonexistent.json")
        assert orchestrator.config["redis"]["host"] == "localhost"
        assert orchestrator.config["debug"] is False

    @pytest.mark.asyncio
    async def test_initialize_success(self, orchestrator):
        """Test successful initialization"""
        with patch('redis.asyncio.Redis') as mock_redis:
            mock_redis_instance = AsyncMock()
            mock_redis.return_value = mock_redis_instance

            await orchestrator.initialize()

            assert orchestrator.redis_client == mock_redis_instance
            assert orchestrator.memory_store is not None
            mock_redis.assert_called_once()

    @pytest.mark.asyncio
    async def test_initialize_redis_failure(self, orchestrator):
        """Test initialization with Redis failure"""
        with patch('redis.asyncio.Redis') as mock_redis:
            mock_redis.side_effect = Exception("Redis connection failed")

            with pytest.raises(Exception, match="Initialization failed"):
                await orchestrator.initialize()

    @pytest.mark.asyncio
    async def test_coding_stage_success(self, orchestrator):
        """Test successful coding stage"""
        await orchestrator.initialize()

        result = await orchestrator._coding_stage("create login system")

        assert result.success is True
        assert result.stage == PipelineStage.CODING
        assert "code" in result.data
        assert "file" in result.data

    @pytest.mark.asyncio
    async def test_coding_stage_error(self, orchestrator):
        """Test coding stage with error"""
        await orchestrator.initialize()

        # Mock autopep8 to raise an exception
        with patch('factory.factory_orchestrator.autopep8') as mock_autopep8:
            mock_autopep8.fix_code.side_effect = Exception("Formatting failed")

            result = await orchestrator._coding_stage("create login system")

            assert result.success is False
            assert result.stage == PipelineStage.CODING
            assert "error" in result.data

    @pytest.mark.asyncio
    async def test_review_stage_success(self, orchestrator):
        """Test successful review stage"""
        await orchestrator.initialize()

        # First set some code in Redis
        await orchestrator.redis_client.set("current_code", "print('hello world')")

        result = await orchestrator._review_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.REVIEW
        assert "issues" in result.data
        assert "code_length" in result.data

    @pytest.mark.asyncio
    async def test_review_stage_no_code(self, orchestrator):
        """Test review stage with no code"""
        await orchestrator.initialize()

        result = await orchestrator._review_stage("test spec")

        assert result.success is False
        assert result.stage == PipelineStage.REVIEW
        assert "No code found" in result.error

    @pytest.mark.asyncio
    async def test_autofix_stage_success(self, orchestrator):
        """Test successful autofix stage"""
        await orchestrator.initialize()

        # Set some code in Redis
        test_code = "print('hello world')"
        await orchestrator.redis_client.set("current_code", test_code)

        result = await orchestrator._autofix_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.AUTOFIX
        assert "original_length" in result.data
        assert "fixed_length" in result.data

    @pytest.mark.asyncio
    async def test_first_test_stage(self, orchestrator):
        """Test first test stage"""
        await orchestrator.initialize()

        result = await orchestrator._first_test_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.FIRST_TEST
        assert "tests_run" in result.data
        assert "passed" in result.data
        assert "failed" in result.data
        assert "coverage" in result.data

    @pytest.mark.asyncio
    async def test_error_solving_stage_no_errors(self, orchestrator):
        """Test error solving stage with no errors"""
        await orchestrator.initialize()

        result = await orchestrator._error_solving_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.ERROR_SOLVING
        assert result.data["no_errors"] is True

    @pytest.mark.asyncio
    async def test_error_solving_stage_with_errors(self, orchestrator):
        """Test error solving stage with errors"""
        await orchestrator.initialize()

        # Set test logs with failures
        test_logs = {"failed": 2, "passed": 3}
        await orchestrator.redis_client.set("test_logs", json.dumps(test_logs))

        result = await orchestrator._error_solving_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.ERROR_SOLVING
        assert "errors_found" in result.data
        assert "errors_fixed" in result.data

    @pytest.mark.asyncio
    async def test_validation_stage_success(self, orchestrator):
        """Test successful validation stage"""
        await orchestrator.initialize()

        # Set validation to true
        await orchestrator.redis_client.set("validated", "true")

        result = await orchestrator._validation_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.VALIDATION
        assert result.data["valid"] is True

    @pytest.mark.asyncio
    async def test_validation_stage_failure(self, orchestrator):
        """Test failed validation stage"""
        await orchestrator.initialize()

        # Set validation to false
        await orchestrator.redis_client.set("validated", "false")

        result = await orchestrator._validation_stage("test spec")

        assert result.success is False
        assert result.stage == PipelineStage.VALIDATION
        assert "Validation failed" in result.error

    @pytest.mark.asyncio
    async def test_monitoring_stage(self, orchestrator):
        """Test monitoring stage"""
        await orchestrator.initialize()

        result = await orchestrator._monitoring_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.MONITORING
        assert "uptime" in result.data
        assert "response_time" in result.data
        assert "error_rate" in result.data

    @pytest.mark.asyncio
    async def test_deployment_stage(self, orchestrator):
        """Test deployment stage"""
        await orchestrator.initialize()

        result = await orchestrator._deployment_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.DEPLOYMENT
        assert "deployed" in result.data
        assert "url" in result.data
        assert "version" in result.data

    @pytest.mark.asyncio
    async def test_reporting_stage(self, orchestrator):
        """Test reporting stage"""
        await orchestrator.initialize()

        result = await orchestrator._reporting_stage("test spec")

        assert result.success is True
        assert result.stage == PipelineStage.REPORTING
        assert "specification" in result.data
        assert "stages_completed" in result.data
        assert "success" in result.data
        assert "timestamp" in result.data

        # Check if report file was created
        assert os.path.exists("factory_report.json")

        # Clean up
        if os.path.exists("factory_report.json"):
            os.remove("factory_report.json")

    @pytest.mark.asyncio
    async def test_full_pipeline_success(self, orchestrator):
        """Test complete pipeline execution"""
        await orchestrator.initialize()

        result = await orchestrator.run_pipeline("create a simple calculator")

        assert result.success is True
        assert result.stage == PipelineStage.REPORTING
        assert result.data["status"] == "completed"

    @pytest.mark.asyncio
    async def test_pipeline_stage_failure(self, orchestrator):
        """Test pipeline with stage failure"""
        await orchestrator.initialize()

        # Mock a stage to fail
        original_coding = orchestrator._coding_stage
        async def failing_coding(spec):
            return PipelineResult(PipelineStage.CODING, False, {}, "Mock failure")

        orchestrator._coding_stage = failing_coding

        result = await orchestrator.run_pipeline("test spec")

        assert result.success is False
        assert result.stage == PipelineStage.CODING
        assert "Mock failure" in result.error

        # Restore original method
        orchestrator._coding_stage = original_coding


class TestPipelineResult:
    """Test cases for PipelineResult"""

    def test_pipeline_result_creation(self):
        """Test PipelineResult creation"""
        result = PipelineResult(
            stage=PipelineStage.CODING,
            success=True,
            data={"test": "data"},
            error=""
        )

        assert result.stage == PipelineStage.CODING
        assert result.success is True
        assert result.data == {"test": "data"}
        assert result.error == ""

    def test_pipeline_result_with_error(self):
        """Test PipelineResult with error"""
        result = PipelineResult(
            stage=PipelineStage.REVIEW,
            success=False,
            data={},
            error="Test error"
        )

        assert result.success is False
        assert result.error == "Test error"


class TestPipelineStage:
    """Test cases for PipelineStage enum"""

    def test_pipeline_stage_values(self):
        """Test PipelineStage enum values"""
        assert PipelineStage.CODING.value == "coding"
        assert PipelineStage.REVIEW.value == "review"
        assert PipelineStage.AUTOFIX.value == "autofix"
        assert PipelineStage.FIRST_TEST.value == "first_test"
        assert PipelineStage.ERROR_SOLVING.value == "error_solving"
        assert PipelineStage.SECOND_TEST.value == "second_test"
        assert PipelineStage.VALIDATION.value == "validation"
        assert PipelineStage.MONITORING.value == "monitoring"
        assert PipelineStage.DEPLOYMENT.value == "deployment"
        assert PipelineStage.REPORTING.value == "reporting"


@pytest.mark.integration
class TestFactoryIntegration:
    """Integration tests for the factory system"""

    @pytest.mark.asyncio
    async def test_run_factory_function(self):
        """Test the run_factory convenience function"""
        result = await run_factory("create a simple function")

        # Should complete successfully (though with mock data)
        assert result is not None
        assert isinstance(result, PipelineResult)

    def test_memory_store_functionality(self):
        """Test memory store operations"""
        from factory.factory_orchestrator import FactoryOrchestrator

        orchestrator = FactoryOrchestrator()
        assert orchestrator.memory_store is None  # Not initialized yet

    @pytest.mark.slow
    @pytest.mark.asyncio
    async def test_full_factory_workflow(self):
        """Test complete factory workflow with realistic timing"""
        import time

        start_time = time.time()
        result = await run_factory("implement user authentication system")
        end_time = time.time()

        # Should complete within reasonable time (less than 30 seconds for mock)
        assert end_time - start_time < 30
        assert result is not None
