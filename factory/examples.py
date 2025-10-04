#!/usr/bin/env python3
"""
AutoGen Factory Usage Examples
Demonstrates how to use the memory learning and deployment aid system
"""

import os
import sys
import time
from typing import Dict, Any

# Add the factory directory to the path
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)

from factory_config import (
    get_factory,
    learn_from_error,
    get_deployment_aid,
    queue_task,
    ErrorMemory
)


def example_error_learning():
    """Example: Learning from deployment errors"""
    print("=== Error Learning Example ===")

    # Simulate a deployment error
    error_type = "docker_build_failed"
    error_message = "failed to build docker image: context not found"
    context = {
        "dockerfile": "Dockerfile",
        "build_context": ".",
        "base_image": "node:18-alpine",
        "error_stage": "build"
    }

    print(f"Learning from error: {error_type}")
    print(f"Error message: {error_message}")

    # Learn from the error
    solution = learn_from_error(error_type, error_message, context)

    if solution:
        print(f"AI suggested solution: {solution}")
    else:
        print("No solution found - stored error for future learning")

    print()


def example_deployment_aid():
    """Example: Getting deployment aid"""
    print("=== Deployment Aid Example ===")

    # Get deployment aid for Docker setup
    task_type = "docker_setup"
    parameters = {
        "application_type": "node.js",
        "port": 3000,
        "environment": "production"
    }

    print(f"Getting deployment aid for: {task_type}")

    try:
        aid_result = get_deployment_aid(task_type, **parameters)
        print("Deployment guidance:")
        print(aid_result.get("guidance", "No guidance available"))
    except Exception as e:
        print(f"Error getting deployment aid: {e}")

    print()


def example_task_queuing():
    """Example: Queueing deployment tasks"""
    print("=== Task Queuing Example ===")

    # Queue various deployment tasks
    tasks = [
        {
            "task_type": "deployment_validation",
            "required_files": ["Dockerfile", "docker-compose.yml", ".env"]
        },
        {
            "task_type": "security_scan",
            "scan_type": "dependency"
        },
        {
            "task_type": "docker_build",
            "image_name": "myapp:latest",
            "dockerfile": "Dockerfile"
        }
    ]

    job_ids = []

    for task in tasks:
        task_type = task.pop("task_type")
        print(f"Queueing task: {task_type}")

        try:
            job_id = queue_task(task_type, **task)
            job_ids.append(job_id)
            print(f"Queued successfully - Job ID: {job_id}")
        except Exception as e:
            print(f"Failed to queue task: {e}")

    print(f"Total jobs queued: {len(job_ids)}")
    print()


def example_memory_storage():
    """Example: Manual memory storage"""
    print("=== Memory Storage Example ===")

    factory = get_factory()

    # Store a successful deployment memory
    success_memory = ErrorMemory(
        error_type="deployment_success",
        error_message="Docker deployment successful",
        context={
            "deployment_type": "docker",
            "image_tag": "v1.2.3",
            "environment": "production",
            "duration": "45 seconds"
        },
        learned=True
    )

    if factory.store_error_memory(success_memory):
        print("Successfully stored deployment success memory")
    else:
        print("Failed to store memory")

    print()


def example_batch_operations():
    """Example: Batch error learning"""
    print("=== Batch Operations Example ===")

    factory = get_factory()

    # Simulate multiple related errors
    errors = [
        {
            "type": "database_connection_failed",
            "message": "connection timeout after 30 seconds",
            "context": {"database": "postgres", "host": "localhost", "port": 5432}
        },
        {
            "type": "database_connection_failed",
            "message": "authentication failed for user 'app'",
            "context": {"database": "postgres", "host": "localhost", "port": 5432}
        },
        {
            "type": "database_connection_failed",
            "message": "connection refused on port 5432",
            "context": {"database": "postgres", "host": "localhost", "port": 5432}
        }
    ]

    for error in errors:
        print(f"Learning from: {error['type']}")
        solution = factory.learn_from_errors(error['type'], error['message'], error['context'])

        if solution:
            print(f"Solution: {solution[:100]}...")
        else:
            print("No solution found")

    print()


def example_redis_monitoring():
    """Example: Monitoring Redis memory usage"""
    print("=== Redis Monitoring Example ===")

    factory = get_factory()

    try:
        # Get Redis info
        info = factory.redis_client.info()

        print("Redis Memory Usage:")
        print(f"  Used Memory: {info.get('used_memory_human', 'Unknown')}")
        print(f"  Peak Memory: {info.get('used_memory_peak_human', 'Unknown')}")
        print(f"  Memory Fragmentation: {info.get('mem_fragmentation_ratio', 'Unknown')}")

        # Get error memory count
        error_keys = factory.redis_client.keys("error_memory:*")
        print(f"  Stored Error Memories: {len(error_keys)}")

        # Get frequency tracking
        freq_keys = factory.redis_client.keys("error_freq:*")
        print(f"  Error Frequency Entries: {len(freq_keys)}")

    except Exception as e:
        print(f"Error monitoring Redis: {e}")

    print()


def main():
    """Run all examples"""
    print("AutoGen Factory Examples")
    print("=" * 50)
    print()

    # Set a dummy API key for examples (replace with real key)
    os.environ["OPENAI_API_KEY"] = "your-api-key-here"

    try:
        # Run examples
        example_error_learning()
        example_deployment_aid()
        example_task_queuing()
        example_memory_storage()
        example_batch_operations()
        example_redis_monitoring()

        print("All examples completed!")

    except Exception as e:
        print(f"Example execution failed: {e}")
        import traceback
        traceback.print_exc()


if __name__ == "__main__":
    main()
