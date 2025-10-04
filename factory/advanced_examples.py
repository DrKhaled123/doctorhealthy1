#!/usr/bin/env python3
"""
Advanced AutoGen Factory Examples
Demonstrates PatternLearningSystem and ContinuousImprovementEngine integration
"""

import os
import sys
import time
import random
from typing import Dict, Any
from datetime import datetime, timedelta

# Add the factory directory to the path
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)

from factory_config import (
    get_factory,
    PatternLearningSystem,
    ContinuousImprovementEngine,
    RiskLevel,
    OptimizationType,
    learn_from_deployment,
    suggest_pattern_solution,
    get_learning_stats,
    analyze_system_performance,
    implement_safe_improvements,
    get_improvement_stats,
    collect_system_metrics
)


def demonstrate_pattern_learning():
    """Demonstrate the PatternLearningSystem capabilities"""
    print("🔄 Pattern Learning System Demonstration")
    print("=" * 50)

    # Create a standalone pattern learning system for demonstration
    learning_system = PatternLearningSystem()

    # Simulate various deployment scenarios
    scenarios = [
        {
            "error_type": "docker_build_failed",
            "success": False,
            "solution": "Check Dockerfile syntax",
            "context": {"dockerfile": "Dockerfile", "stage": "build"},
            "execution_time": 45.2
        },
        {
            "error_type": "docker_build_failed",
            "success": True,
            "solution": "Check Dockerfile syntax",
            "context": {"dockerfile": "Dockerfile", "stage": "build"},
            "execution_time": 52.1
        },
        {
            "error_type": "docker_build_failed",
            "success": True,
            "solution": "Check Dockerfile syntax",
            "context": {"dockerfile": "Dockerfile", "stage": "build"},
            "execution_time": 38.7
        },
        {
            "error_type": "database_connection_failed",
            "success": False,
            "solution": "Verify connection string",
            "context": {"database": "postgres", "port": 5432},
            "execution_time": 15.3
        },
        {
            "error_type": "database_connection_failed",
            "success": True,
            "solution": "Verify connection string",
            "context": {"database": "postgres", "port": 5432},
            "execution_time": 12.8
        },
        {
            "error_type": "test_execution_failed",
            "success": False,
            "solution": "Install test dependencies",
            "context": {"test_framework": "pytest", "coverage": True},
            "execution_time": 67.4
        },
        {
            "error_type": "test_execution_failed",
            "success": True,
            "solution": "Install test dependencies",
            "context": {"test_framework": "pytest", "coverage": True},
            "execution_time": 58.9
        }
    ]

    print("📚 Training the learning system with deployment scenarios...")

    for i, scenario in enumerate(scenarios, 1):
        success = learning_system.learn_from_deployment(scenario)
        status = "✅" if success else "❌"
        print(f"  {status} Scenario {i}: {scenario['error_type']} ({scenario['success']})")

        # Small delay to show progression
        time.sleep(0.1)

    print("\n🎯 Testing pattern recognition and solution suggestions...")
    # Test solution suggestions
    test_cases = [
        "docker_build_failed",
        "database_connection_failed",
        "test_execution_failed",
        "unknown_error_type"  # This should return None
    ]

    for error_type in test_cases:
        solution = learning_system.suggest_solution(error_type)
        if solution:
            print(f"  💡 {error_type} → Suggested: {solution}")
        else:
            print(f"  ❓ {error_type} → No reliable solution available")

    # Display learning statistics
    stats = learning_system.get_learning_stats()
    print("\n📊 Learning Statistics:")
    print(f"  • Total error patterns: {stats.get('total_error_patterns', 0)}")
    print(f"  • Total success patterns: {stats.get('total_success_patterns', 0)}")
    print(f"  • Unique error types: {stats.get('unique_error_types', 0)}")
    print(f"  • Average success rate: {stats.get('average_success_rate', 0):.2%}")

    print("\n✅ Pattern learning demonstration completed!\n")


def demonstrate_continuous_improvement():
    """Demonstrate the ContinuousImprovementEngine capabilities"""
    print("🔄 Continuous Improvement Engine Demonstration")
    print("=" * 50)

    # Create improvement engine
    improvement_engine = ContinuousImprovementEngine()

    print("📈 Simulating system performance data...")

    # Simulate performance history
    for i in range(15):
        # Generate realistic performance data
        success_rate = random.uniform(0.7, 0.95)
        response_time = random.uniform(10, 100)
        memory_usage = random.uniform(50, 300)

        entry = {
            "task_type": random.choice(["docker_build", "test_run", "security_scan"]),
            "success": random.random() > 0.2,  # 80% success rate
            "response_time": response_time,
            "memory_usage": memory_usage,
            "timestamp": datetime.now().isoformat()
        }

        improvement_engine.performance_history.append(entry)

        print(f"  📊 Entry {i+1}: {entry['task_type']} - {'✅' if entry['success'] else '❌'} "
              f"({response_time:.1f}ms, {memory_usage:.1f}MB)")

    print("\n🔍 Analyzing system performance...")
    # Analyze performance
    suggestions = improvement_engine.analyze_performance()

    print(f"🎯 Generated {len(suggestions)} optimization suggestions:")

    for i, suggestion in enumerate(suggestions, 1):
        risk_emoji = {"low": "🟢", "medium": "🟡", "high": "🔴", "critical": "🚨"}
        risk_icon = risk_emoji.get(suggestion.get('risk_level', 'medium'), '⚪')

        print(f"  {risk_icon} {i}. {suggestion['title']}")
        print(f"     Type: {suggestion.get('type', 'unknown').value}")
        print(f"     Confidence: {suggestion.get('confidence', 0):.1%}")
        print(f"     Action: {suggestion.get('suggested_action', 'N/A')[:60]}...")
        print()

    print("🚀 Implementing safe improvements...")

    # Implement improvements
    results = improvement_engine.implement_improvements()

    print("📋 Implementation Results:")
    print(f"  • Implemented: {results['implemented']}")
    print(f"  • Failed: {results['failed']}")

    for detail in results['details']:
        status = "✅" if detail['success'] else "❌"
        print(f"    {status} {detail['title']} (Risk: {detail['risk_level']})")

    # Display improvement statistics
    stats = improvement_engine.get_improvement_stats()
    print("\n📊 Improvement Statistics:")
    print(f"  • Total suggestions: {stats.get('total_suggestions', 0)}")
    print(f"  • Applied improvements: {stats.get('applied_improvements', 0)}")
    print(f"  • Performance history size: {stats.get('performance_history_size', 0)}")

    risk_dist = stats.get('risk_level_distribution', {})
    if risk_dist:
        print("  • Risk distribution:")
        for risk_level, count in risk_dist.items():
            print(f"    - {risk_level}: {count}")

    print("\n✅ Continuous improvement demonstration completed!\n")


def demonstrate_integrated_system():
    """Demonstrate the integrated factory system"""
    print("🏭 Integrated Factory System Demonstration")
    print("=" * 50)

    try:
        # Initialize the factory
        factory = get_factory()

        print("🚀 Running integrated deployment simulation...")

        # Simulate a series of deployment tasks
        deployment_tasks = [
            {
                "task_type": "deployment_validation",
                "required_files": ["Dockerfile", "docker-compose.yml"],
                "required_env_vars": ["DATABASE_URL", "API_KEY"]
            },
            {
                "task_type": "docker_build",
                "image_name": "demo-app:latest",
                "dockerfile": "Dockerfile"
            },
            {
                "task_type": "security_scan",
                "scan_type": "dependency"
            },
            {
                "task_type": "test_run",
                "test_command": "go test ./...",
                "coverage": True
            }
        ]

        for i, task_config in enumerate(deployment_tasks, 1):
            task_type = task_config.pop("task_type")
            print(f"\n📋 Task {i}: {task_type}")

            # Queue the task
            job_id = factory.queue_deployment_task(task_type, task_config)
            print(f"  🎫 Queued job: {job_id}")

            # Simulate task execution for demo (in real scenario, this would be handled by RQ worker)
            time.sleep(0.5)

            # Simulate learning from the deployment
            deployment_data = {
                "error_type": task_type,
                "success": random.choice([True, True, True, False]),  # 75% success rate
                "solution": f"Standard {task_type} procedure",
                "context": task_config,
                "execution_time": random.uniform(20, 120),
                "timestamp": datetime.now()
            }

            # Learn from the deployment
            learning_success = learn_from_deployment(deployment_data)
            status = "✅" if learning_success else "❌"
            print(f"  {status} Learning recorded: {deployment_data['success']}")

        print("\n🎓 Learning from deployment patterns...")
        # Test pattern-based solution suggestions
        test_error_types = ["deployment_validation", "docker_build", "security_scan"]

        for error_type in test_error_types:
            solution = suggest_pattern_solution(error_type)
            if solution:
                print(f"  💡 Pattern solution for {error_type}: {solution}")
            else:
                print(f"  ❓ No pattern solution for {error_type}")

        print("\n📊 System Statistics:")
        # Get comprehensive statistics
        learning_stats = get_learning_stats()
        improvement_stats = get_improvement_stats()

        print("  🎯 Learning System:")
        print(f"    - Error patterns: {learning_stats.get('total_error_patterns', 0)}")
        print(f"    - Success patterns: {learning_stats.get('total_success_patterns', 0)}")
        print(f"    - Success rate: {learning_stats.get('average_success_rate', 0):.2%}")

        print("  🔧 Improvement Engine:")
        print(f"    - Suggestions: {improvement_stats.get('total_suggestions', 0)}")
        print(f"    - Applied: {improvement_stats.get('applied_improvements', 0)}")

        print("\n✅ Integrated system demonstration completed!\n")

    except Exception as e:
        print(f"❌ Integrated demonstration failed: {e}")
        import traceback
        traceback.print_exc()


def demonstrate_error_recovery():
    """Demonstrate error recovery and learning capabilities"""
    print("🔧 Error Recovery and Learning Demonstration")
    print("=" * 50)

    try:
        factory = get_factory()

        print("💥 Simulating deployment errors and recovery...")

        # Simulate a problematic deployment scenario
        error_scenarios = [
            {
                "error_type": "docker_build_context_missing",
                "error_message": "docker build failed: context not found",
                "context": {"build_context": ".", "dockerfile": "Dockerfile"},
                "should_succeed": False
            },
            {
                "error_type": "docker_build_context_missing",
                "error_message": "docker build failed: context not found",
                "context": {"build_context": ".", "dockerfile": "Dockerfile"},
                "should_succeed": True  # Fixed the issue
            },
            {
                "error_type": "database_port_blocked",
                "error_message": "connection refused on port 5432",
                "context": {"database": "postgres", "port": 5432},
                "should_succeed": False
            },
            {
                "error_type": "database_port_blocked",
                "error_message": "connection refused on port 5432",
                "context": {"database": "postgres", "port": 5433},  # Different port
                "should_succeed": True  # Fixed by changing port
            }
        ]

        for i, scenario in enumerate(error_scenarios, 1):
            print(f"\n🔍 Scenario {i}: {scenario['error_type']}")

            # Learn from the error
            solution = factory.learn_from_errors(
                scenario['error_type'],
                scenario['error_message'],
                scenario['context']
            )

            if solution:
                print(f"  🤖 AI Solution: {solution[:80]}...")
            else:
                print("  🤷 No AI solution available")

            # Learn from deployment outcome
            deployment_data = {
                "error_type": scenario['error_type'],
                "success": scenario['should_succeed'],
                "solution": solution or "manual_fix",
                "context": scenario['context'],
                "execution_time": random.uniform(30, 90),
                "timestamp": datetime.now()
            }

            learning_success = learn_from_deployment(deployment_data)
            status = "✅" if learning_success else "❌"
            print(f"  {status} Deployment learning: {'Success' if scenario['should_succeed'] else 'Failed'}")

        print("
🎯 Testing learned pattern recognition..."
        # Test if the system learned from the patterns
        for scenario in error_scenarios:
            if scenario['should_succeed']:
                pattern_solution = suggest_pattern_solution(scenario['error_type'])
                if pattern_solution:
                    print(f"  💡 Learned pattern for {scenario['error_type']}: {pattern_solution}")
                else:
                    print(f"  ❓ No learned pattern for {scenario['error_type']}")

        print("\n✅ Error recovery demonstration completed!\n")

    except Exception as e:
        print(f"❌ Error recovery demonstration failed: {e}")
        import traceback
        traceback.print_exc()


def demonstrate_performance_monitoring():
    """Demonstrate performance monitoring and optimization"""
    print("📊 Performance Monitoring Demonstration")
    print("=" * 50)

    try:
        print("📈 Collecting system metrics...")

        # Collect current metrics
        metrics = collect_system_metrics()

        print("📋 Current System Metrics:")
        for category, values in metrics.items():
            if isinstance(values, dict):
                print(f"  • {category}:")
                for metric, value in values.items():
                    if isinstance(value, float):
                        print(f"    - {metric}: {value:.2f}")
                    else:
                        print(f"    - {metric}: {value}")
            else:
                print(f"  • {category}: {values}")

        print("
🔍 Analyzing performance and generating suggestions..."
        # Analyze performance
        suggestions = analyze_system_performance()

        print(f"🎯 Generated {len(suggestions)} optimization suggestions:")

        for i, suggestion in enumerate(suggestions, 1):
            print(f"\n  {i}. {suggestion['title']}")
            print(f"     Type: {suggestion.get('type', 'unknown').value}")
            print(f"     Risk Level: {suggestion.get('risk_level', 'unknown').value}")
            print(f"     Confidence: {suggestion.get('confidence', 0):.1%}")
            print(f"     Description: {suggestion.get('description', 'N/A')}")

        print("
🚀 Implementing safe optimizations..."
        # Implement improvements
        improvement_results = implement_safe_improvements()

        print("📋 Implementation Results:")
        print(f"  • Successfully implemented: {improvement_results['implemented']}")
        print(f"  • Failed to implement: {improvement_results['failed']}")

        if improvement_results['details']:
            print("  • Details:")
            for detail in improvement_results['details']:
                status = "✅" if detail['success'] else "❌"
                print(f"    {status} {detail['title']}")

        print("
📊 Final System Statistics:"
        # Display final statistics
        final_learning_stats = get_learning_stats()
        final_improvement_stats = get_improvement_stats()

        print("  🎓 Learning System:")
        for key, value in final_learning_stats.items():
            if isinstance(value, float):
                print(f"    • {key}: {value:.3f}")
            else:
                print(f"    • {key}: {value}")

        print("  🔧 Improvement Engine:")
        for key, value in final_improvement_stats.items():
            if key == 'risk_level_distribution':
                print(f"    • {key}:")
                for risk, count in value.items():
                    print(f"      - {risk}: {count}")
            elif key == 'optimization_type_distribution':
                print(f"    • {key}:")
                for opt_type, count in value.items():
                    print(f"      - {opt_type}: {count}")
            elif isinstance(value, float):
                print(f"    • {key}: {value:.3f}")
            else:
                print(f"    • {key}: {value}")

        print("\n✅ Performance monitoring demonstration completed!\n")

    except Exception as e:
        print(f"❌ Performance monitoring demonstration failed: {e}")
        import traceback
        traceback.print_exc()


def main():
    """Run all advanced demonstrations"""
    print("🚀 Advanced AutoGen Factory Demonstrations")
    print("=" * 60)
    print()

    # Set a dummy API key for examples (replace with real key)
    os.environ["OPENAI_API_KEY"] = "your-api-key-here"

    try:
        # Run all demonstrations
        demonstrate_pattern_learning()
        print("\n" + "="*60 + "\n")

        demonstrate_continuous_improvement()
        print("\n" + "="*60 + "\n")

        demonstrate_integrated_system()
        print("\n" + "="*60 + "\n")

        demonstrate_error_recovery()
        print("\n" + "="*60 + "\n")

        demonstrate_performance_monitoring()

        print("\n🎉 All advanced demonstrations completed successfully!")
        print("\nKey Features Demonstrated:")
        print("✅ Pattern Learning System - Learns from deployment outcomes")
        print("✅ Continuous Improvement Engine - Analyzes and optimizes performance")
        print("✅ Integrated Factory System - AutoGen + Redis + RQ + Learning")
        print("✅ Error Recovery - Intelligent error analysis and solution suggestions")
        print("✅ Performance Monitoring - Real-time metrics and optimization")

    except Exception as e:
        print(f"❌ Demonstration execution failed: {e}")
        import traceback
        traceback.print_exc()


if __name__ == "__main__":
    main()
