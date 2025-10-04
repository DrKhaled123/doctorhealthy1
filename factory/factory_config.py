"""
AutoGen + Redis + RQ Factory Configuration
Memory Learning and Deployment Aid System
"""

import os
import redis
import json
import logging
from typing import Dict, List, Any, Optional, Tuple
from datetime import datetime, timedelta
from dataclasses import dataclass, asdict
from pydantic import BaseModel, Field
from collections import defaultdict, deque
import autogen
from rq import Worker, Queue, Connection
from loguru import logger
import uuid
import time
import statistics
from enum import Enum

# Configuration Models
class RedisConfig(BaseModel):
    host: str = Field(default="localhost")
    port: int = Field(default=6379)
    db: int = Field(default=0)
    password: Optional[str] = Field(default=None)
    decode_responses: bool = Field(default=True)

class RQConfig(BaseModel):
    queue_name: str = Field(default="factory_tasks")
    worker_name: str = Field(default="factory_worker")
    connection_pool_size: int = Field(default=10)
    job_timeout: int = Field(default=3600)  # 1 hour

class AutoGenConfig(BaseModel):
    model: str = Field(default="gpt-4")
    temperature: float = Field(default=0.7)
    max_tokens: int = Field(default=2000)
    api_key: Optional[str] = Field(default=None)

class MemoryConfig(BaseModel):
    max_memories: int = Field(default=1000)
    memory_ttl: int = Field(default=86400 * 30)  # 30 days
    learning_rate: float = Field(default=0.1)
    similarity_threshold: float = Field(default=0.8)

class FactoryConfig(BaseModel):
    redis: RedisConfig = Field(default_factory=RedisConfig)
    rq: RQConfig = Field(default_factory=RQConfig)
    autogen: AutoGenConfig = Field(default_factory=AutoGenConfig)
    memory: MemoryConfig = Field(default_factory=MemoryConfig)
    debug: bool = Field(default=False)

# Memory and Learning Models
class ErrorMemory(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    error_type: str
    error_message: str
    context: Dict[str, Any]
    solution: Optional[str] = None
    timestamp: datetime = Field(default_factory=datetime.now)
    frequency: int = 1
    learned: bool = False

class DeploymentAid(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    task_type: str
    parameters: Dict[str, Any]
    result: Optional[Dict[str, Any]] = None
    success: bool = False
    timestamp: datetime = Field(default_factory=datetime.now)
    execution_time: Optional[float] = None

# Enums for type safety
class RiskLevel(Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"

class OptimizationType(Enum):
    PERFORMANCE = "performance"
    MEMORY = "memory"
    RELIABILITY = "reliability"
    SECURITY = "security"

# Pattern Learning System
class PatternLearningSystem:
    """Safe and robust pattern learning system for deployment outcomes"""

    def __init__(self):
        self.error_database: Dict[str, List[Dict[str, Any]]] = defaultdict(list)
        self.success_patterns: Dict[str, List[Dict[str, Any]]] = defaultdict(list)
        self.performance_metrics: Dict[str, deque] = defaultdict(lambda: deque(maxlen=100))
        self.max_patterns_per_error = 50  # Limit to prevent memory bloat
        self.min_success_rate = 0.7  # Minimum success rate to consider a pattern reliable

    def learn_from_deployment(self, deployment_data: Dict[str, Any]) -> bool:
        """
        Safely learn from deployment outcomes with validation

        Args:
            deployment_data: Dictionary containing deployment information
                - error_type: Type of error encountered
                - solution: Solution that was applied
                - success: Whether the deployment succeeded
                - context: Additional context information
                - execution_time: How long the deployment took
                - timestamp: When the deployment occurred

        Returns:
            bool: True if learning was successful, False otherwise
        """
        try:
            # Validate required fields
            required_fields = ['error_type', 'success']
            for field in required_fields:
                if field not in deployment_data:
                    logger.warning(f"Missing required field '{field}' in deployment data")
                    return False

            error_type = deployment_data['error_type']
            solution = deployment_data.get('solution')
            success = deployment_data['success']
            timestamp = deployment_data.get('timestamp', datetime.now())

            # Validate data types
            if not isinstance(success, bool):
                logger.warning(f"Invalid success type: {type(success)}, expected bool")
                return False

            if not isinstance(error_type, str) or not error_type.strip():
                logger.warning(f"Invalid error_type: {error_type}")
                return False

            # Create learning entry
            learning_entry = {
                'solution': solution,
                'success': success,
                'timestamp': timestamp,
                'context': deployment_data.get('context', {}),
                'execution_time': deployment_data.get('execution_time')
            }

            # Store in appropriate database
            if success:
                self._store_success_pattern(error_type, learning_entry)
            else:
                self._store_error_pattern(error_type, learning_entry)

            # Update performance metrics
            if 'execution_time' in deployment_data:
                self.performance_metrics[f"{error_type}_success" if success else f"{error_type}_failure"].append(
                    deployment_data['execution_time']
                )

            logger.info(f"Successfully learned from deployment: {error_type} (success: {success})")
            return True

        except Exception as e:
            logger.error(f"Failed to learn from deployment: {e}")
            return False

    def _store_success_pattern(self, error_type: str, entry: Dict[str, Any]):
        """Store successful deployment pattern"""
        if len(self.success_patterns[error_type]) >= self.max_patterns_per_error:
            # Remove oldest entry to maintain size limit
            self.success_patterns[error_type].pop(0)

        self.success_patterns[error_type].append(entry)

    def _store_error_pattern(self, error_type: str, entry: Dict[str, Any]):
        """Store error pattern for future learning"""
        if len(self.error_database[error_type]) >= self.max_patterns_per_error:
            # Remove oldest entry to maintain size limit
            self.error_database[error_type].pop(0)

        self.error_database[error_type].append(entry)

    def suggest_solution(self, error_type: str, context: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """
        Suggest best solution based on historical data with safety checks

        Args:
            error_type: The type of error to find a solution for
            context: Additional context to help select the best solution

        Returns:
            Optional[str]: Best solution if available and reliable, None otherwise
        """
        try:
            if error_type not in self.error_database and error_type not in self.success_patterns:
                logger.info(f"No historical data available for error type: {error_type}")
                return None

            # Get all relevant patterns
            error_patterns = self.error_database.get(error_type, [])
            success_patterns = self.success_patterns.get(error_type, [])

            # Combine and filter patterns
            all_patterns = error_patterns + success_patterns

            if not all_patterns:
                return None

            # Calculate success rates for each solution
            solution_stats = self._calculate_solution_stats(all_patterns)

            # Find best solution based on success rate and recency
            best_solution = self._select_best_solution(solution_stats, context)

            if best_solution and best_solution['success_rate'] >= self.min_success_rate:
                logger.info(f"Suggested solution for {error_type} with {best_solution['success_rate']:.2f} success rate")
                return best_solution['solution']

            logger.info(f"No reliable solution found for {error_type} (best success rate: {best_solution['success_rate']:.2f} if best_solution else 0)")
            return None

        except Exception as e:
            logger.error(f"Failed to suggest solution: {e}")
            return None

    def _calculate_solution_stats(self, patterns: List[Dict[str, Any]]) -> Dict[str, Dict[str, Any]]:
        """Calculate statistics for each solution"""
        solution_stats = defaultdict(lambda: {
            'count': 0,
            'success_count': 0,
            'total_execution_time': 0,
            'recent_count': 0,
            'timestamps': []
        })

        cutoff_time = datetime.now() - timedelta(days=30)  # Consider last 30 days

        for pattern in patterns:
            solution = pattern.get('solution', 'unknown')
            if not solution or solution == 'unknown':
                continue

            stats = solution_stats[solution]
            stats['count'] += 1
            stats['timestamps'].append(pattern['timestamp'])

            if pattern['success']:
                stats['success_count'] += 1

            if pattern.get('execution_time'):
                stats['total_execution_time'] += pattern['execution_time']

            # Count recent occurrences
            if pattern['timestamp'] > cutoff_time:
                stats['recent_count'] += 1

        # Calculate final statistics
        for solution, stats in solution_stats.items():
            if stats['count'] > 0:
                stats['success_rate'] = stats['success_count'] / stats['count']
                stats['avg_execution_time'] = stats['total_execution_time'] / stats['count']
                stats['recency_score'] = min(stats['recent_count'] / max(stats['count'], 1), 1.0)
            else:
                stats['success_rate'] = 0.0
                stats['avg_execution_time'] = 0.0
                stats['recency_score'] = 0.0

        return dict(solution_stats)

    def _select_best_solution(self, solution_stats: Dict[str, Dict[str, Any]], context: Optional[Dict[str, Any]] = None) -> Optional[Dict[str, Any]]:
        """Select the best solution based on multiple factors"""
        if not solution_stats:
            return None

        scored_solutions = []

        for solution, stats in solution_stats.items():
            # Calculate composite score
            success_weight = 0.6
            recency_weight = 0.3
            performance_weight = 0.1

            # Normalize execution time (lower is better)
            avg_time = stats.get('avg_execution_time', 0)
            time_score = max(0, 1 - (avg_time / 3600))  # Normalize to 1 hour max

            composite_score = (
                stats['success_rate'] * success_weight +
                stats['recency_score'] * recency_weight +
                time_score * performance_weight
            )

            scored_solutions.append({
                'solution': solution,
                'score': composite_score,
                'success_rate': stats['success_rate'],
                'recency_score': stats['recency_score'],
                'avg_execution_time': avg_time,
                'count': stats['count']
            })

        # Sort by composite score
        scored_solutions.sort(key=lambda x: x['score'], reverse=True)

        return scored_solutions[0] if scored_solutions else None

    def get_learning_stats(self) -> Dict[str, Any]:
        """Get comprehensive learning statistics"""
        try:
            total_errors = sum(len(patterns) for patterns in self.error_database.values())
            total_successes = sum(len(patterns) for patterns in self.success_patterns.values())

            # Calculate average success rates
            avg_success_rate = 0.0
            if total_errors + total_successes > 0:
                total_attempts = sum(stats['count'] for solution_stats in [
                    self._calculate_solution_stats(patterns)
                    for patterns in list(self.error_database.values()) + list(self.success_patterns.values())
                ] for stats in solution_stats.values())

                if total_attempts > 0:
                    total_successes_count = sum(stats['success_count'] for solution_stats in [
                        self._calculate_solution_stats(patterns)
                        for patterns in list(self.error_database.values()) + list(self.success_patterns.values())
                    ] for stats in solution_stats.values())

                    avg_success_rate = total_successes_count / total_attempts

            return {
                'total_error_patterns': total_errors,
                'total_success_patterns': total_successes,
                'unique_error_types': len(self.error_database),
                'average_success_rate': avg_success_rate,
                'performance_metrics_count': {k: len(v) for k, v in self.performance_metrics.items()},
                'last_updated': datetime.now().isoformat()
            }

        except Exception as e:
            logger.error(f"Failed to get learning stats: {e}")
            return {}

# Continuous Improvement Engine
class ContinuousImprovementEngine:
    """Safe continuous improvement engine with risk assessment"""

    def __init__(self):
        self.improvement_log: List[Dict[str, Any]] = []
        self.optimization_suggestions: List[Dict[str, Any]] = []
        self.applied_improvements: List[Dict[str, Any]] = []
        self.performance_history: deque = deque(maxlen=1000)
        self.risk_thresholds = {
            RiskLevel.LOW: 0.8,
            RiskLevel.MEDIUM: 0.6,
            RiskLevel.HIGH: 0.4,
            RiskLevel.CRITICAL: 0.2
        }

    def collect_metrics(self) -> Dict[str, Any]:
        """Safely collect system performance metrics"""
        try:
            metrics = {
                'timestamp': datetime.now(),
                'memory_usage': {},
                'response_times': {},
                'error_rates': {},
                'throughput': {}
            }

            # Memory metrics (safe defaults if collection fails)
            try:
                import psutil
                process = psutil.Process()
                metrics['memory_usage'] = {
                    'rss_mb': process.memory_info().rss / 1024 / 1024,
                    'vms_mb': process.memory_info().vms / 1024 / 1024,
                    'cpu_percent': process.cpu_percent()
                }
            except ImportError:
                metrics['memory_usage'] = {'note': 'psutil not available'}
            except Exception as e:
                logger.warning(f"Failed to collect memory metrics: {e}")
                metrics['memory_usage'] = {'error': str(e)}

            # Response time metrics from history
            if self.performance_history:
                response_times = [entry.get('response_time', 0) for entry in self.performance_history]
                if response_times:
                    metrics['response_times'] = {
                        'average': statistics.mean(response_times),
                        'median': statistics.median(response_times),
                        'min': min(response_times),
                        'max': max(response_times),
                        'count': len(response_times)
                    }

            # Error rate calculation
            recent_entries = list(self.performance_history)[-100:] if len(self.performance_history) > 100 else self.performance_history
            if recent_entries:
                error_count = sum(1 for entry in recent_entries if not entry.get('success', True))
                metrics['error_rates'] = {
                    'recent_error_rate': error_count / len(recent_entries),
                    'total_entries': len(recent_entries)
                }

            return metrics

        except Exception as e:
            logger.error(f"Failed to collect metrics: {e}")
            return {'error': str(e), 'timestamp': datetime.now()}

    def generate_optimization_suggestions(self, metrics: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate safe optimization suggestions based on metrics"""
        suggestions = []

        try:
            # Memory optimization suggestions
            memory_usage = metrics.get('memory_usage', {})
            if isinstance(memory_usage, dict) and 'rss_mb' in memory_usage:
                memory_mb = memory_usage['rss_mb']

                if memory_mb > 500:  # High memory usage
                    suggestions.append({
                        'type': OptimizationType.MEMORY,
                        'title': 'High Memory Usage Detected',
                        'description': f'Process using {memory_mb:.1f}MB of memory',
                        'risk_level': RiskLevel.MEDIUM,
                        'suggested_action': 'Consider optimizing memory usage or increasing available memory',
                        'confidence': 0.8,
                        'metrics': memory_usage
                    })

            # Response time optimization
            response_times = metrics.get('response_times', {})
            if isinstance(response_times, dict) and 'average' in response_times:
                avg_response = response_times['average']

                if avg_response > 30:  # Slow response time
                    suggestions.append({
                        'type': OptimizationType.PERFORMANCE,
                        'title': 'Slow Response Times',
                        'description': f'Average response time: {avg_response:.2f}s',
                        'risk_level': RiskLevel.LOW,
                        'suggested_action': 'Consider caching frequently accessed data or optimizing algorithms',
                        'confidence': 0.7,
                        'metrics': response_times
                    })

            # Error rate optimization
            error_rates = metrics.get('error_rates', {})
            if isinstance(error_rates, dict) and 'recent_error_rate' in error_rates:
                error_rate = error_rates['recent_error_rate']

                if error_rate > 0.1:  # High error rate
                    suggestions.append({
                        'type': OptimizationType.RELIABILITY,
                        'title': 'High Error Rate',
                        'description': f'Error rate: {error_rate:.2%}',
                        'risk_level': RiskLevel.HIGH,
                        'suggested_action': 'Review error patterns and implement additional error handling',
                        'confidence': 0.9,
                        'metrics': error_rates
                    })

            # Learning pattern optimization
            if len(self.optimization_suggestions) > 10:
                suggestions.append({
                    'type': OptimizationType.MEMORY,
                    'title': 'Suggestion Cleanup',
                    'description': 'Accumulated many optimization suggestions',
                    'risk_level': RiskLevel.LOW,
                    'suggested_action': 'Clean up old optimization suggestions to free memory',
                    'confidence': 0.6,
                    'metrics': {'suggestion_count': len(self.optimization_suggestions)}
                })

            self.optimization_suggestions.extend(suggestions)
            return suggestions

        except Exception as e:
            logger.error(f"Failed to generate optimization suggestions: {e}")
            return []

    def analyze_performance(self) -> List[Dict[str, Any]]:
        """Analyze system performance and suggest improvements"""
        try:
            metrics = self.collect_metrics()
            suggestions = self.generate_optimization_suggestions(metrics)

            # Log analysis results
            self.improvement_log.append({
                'timestamp': datetime.now(),
                'metrics': metrics,
                'suggestions_generated': len(suggestions),
                'analysis_type': 'performance'
            })

            logger.info(f"Performance analysis completed: {len(suggestions)} suggestions generated")
            return suggestions

        except Exception as e:
            logger.error(f"Failed to analyze performance: {e}")
            return []

    def implement_improvements(self) -> Dict[str, Any]:
        """
        Safely implement low-risk optimizations

        Returns:
            Dict containing results of improvement implementation
        """
        results = {
            'implemented': 0,
            'skipped': 0,
            'failed': 0,
            'details': []
        }

        try:
            # Only implement low-risk suggestions
            low_risk_suggestions = [
                s for s in self.optimization_suggestions
                if s.get('risk_level') == RiskLevel.LOW
            ]

            for suggestion in low_risk_suggestions[:5]:  # Limit to 5 per cycle
                try:
                    success = self.apply_improvement(suggestion)

                    if success:
                        results['implemented'] += 1
                        self.log_improvement(suggestion)
                        self.optimization_suggestions.remove(suggestion)
                    else:
                        results['failed'] += 1

                    results['details'].append({
                        'title': suggestion['title'],
                        'success': success,
                        'risk_level': suggestion['risk_level'].value
                    })

                except Exception as e:
                    logger.error(f"Failed to apply improvement {suggestion.get('title', 'unknown')}: {e}")
                    results['failed'] += 1
                    results['details'].append({
                        'title': suggestion.get('title', 'unknown'),
                        'success': False,
                        'error': str(e),
                        'risk_level': suggestion.get('risk_level', 'unknown').value
                    })

            # Remove processed suggestions
            for suggestion in low_risk_suggestions[:results['implemented']]:
                if suggestion in self.optimization_suggestions:
                    self.optimization_suggestions.remove(suggestion)

            logger.info(f"Improvement implementation completed: {results['implemented']} implemented, {results['failed']} failed")
            return results

        except Exception as e:
            logger.error(f"Failed to implement improvements: {e}")
            results['failed'] = len(low_risk_suggestions)
            return results

    def apply_improvement(self, suggestion: Dict[str, Any]) -> bool:
        """Apply a specific improvement with safety checks"""
        try:
            improvement_type = suggestion.get('type')
            title = suggestion.get('title', 'Unknown improvement')

            # Validate improvement type
            if not isinstance(improvement_type, OptimizationType):
                logger.warning(f"Invalid improvement type: {improvement_type}")
                return False

            # Apply improvement based on type
            if improvement_type == OptimizationType.MEMORY:
                return self._apply_memory_optimization(suggestion)
            elif improvement_type == OptimizationType.PERFORMANCE:
                return self._apply_performance_optimization(suggestion)
            elif improvement_type == OptimizationType.RELIABILITY:
                return self._apply_reliability_optimization(suggestion)
            elif improvement_type == OptimizationType.SECURITY:
                return self._apply_security_optimization(suggestion)
            else:
                logger.warning(f"Unknown improvement type: {improvement_type}")
                return False

        except Exception as e:
            logger.error(f"Failed to apply improvement {suggestion.get('title', 'unknown')}: {e}")
            return False

    def _apply_memory_optimization(self, suggestion: Dict[str, Any]) -> bool:
        """Apply memory optimization improvements"""
        try:
            # Clear old performance history if it's getting large
            if len(self.performance_history) > 800:
                # Keep only the most recent 500 entries
                kept_entries = list(self.performance_history)[-500:]
                self.performance_history.clear()
                self.performance_history.extend(kept_entries)
                return True

            # Clear old improvement log if it's getting large
            if len(self.improvement_log) > 100:
                self.improvement_log = self.improvement_log[-50:]  # Keep last 50 entries
                return True

            return True

        except Exception as e:
            logger.error(f"Failed to apply memory optimization: {e}")
            return False

    def _apply_performance_optimization(self, suggestion: Dict[str, Any]) -> bool:
        """Apply performance optimization improvements"""
        try:
            # This is a placeholder for actual performance optimizations
            # In a real system, you might adjust caching strategies, connection pooling, etc.
            logger.info(f"Applied performance optimization: {suggestion.get('title', 'unknown')}")
            return True

        except Exception as e:
            logger.error(f"Failed to apply performance optimization: {e}")
            return False

    def _apply_reliability_optimization(self, suggestion: Dict[str, Any]) -> bool:
        """Apply reliability optimization improvements"""
        try:
            # This is a placeholder for actual reliability optimizations
            # In a real system, you might adjust retry logic, circuit breakers, etc.
            logger.info(f"Applied reliability optimization: {suggestion.get('title', 'unknown')}")
            return True

        except Exception as e:
            logger.error(f"Failed to apply reliability optimization: {e}")
            return False

    def _apply_security_optimization(self, suggestion: Dict[str, Any]) -> bool:
        """Apply security optimization improvements"""
        try:
            # This is a placeholder for actual security optimizations
            # In a real system, you might adjust authentication, encryption, etc.
            logger.info(f"Applied security optimization: {suggestion.get('title', 'unknown')}")
            return True

        except Exception as e:
            logger.error(f"Failed to apply security optimization: {e}")
            return False

    def log_improvement(self, suggestion: Dict[str, Any]):
        """Log an applied improvement"""
        try:
            log_entry = {
                'timestamp': datetime.now(),
                'title': suggestion.get('title'),
                'type': suggestion.get('type', 'unknown').value,
                'risk_level': suggestion.get('risk_level', 'unknown').value,
                'confidence': suggestion.get('confidence', 0.0),
                'description': suggestion.get('description', ''),
                'status': 'applied'
            }

            self.applied_improvements.append(log_entry)

            # Keep only last 100 applied improvements
            if len(self.applied_improvements) > 100:
                self.applied_improvements = self.applied_improvements[-100:]

        except Exception as e:
            logger.error(f"Failed to log improvement: {e}")

    def get_improvement_stats(self) -> Dict[str, Any]:
        """Get comprehensive improvement statistics"""
        try:
            return {
                'total_suggestions': len(self.optimization_suggestions),
                'applied_improvements': len(self.applied_improvements),
                'improvement_log_size': len(self.improvement_log),
                'performance_history_size': len(self.performance_history),
                'risk_level_distribution': self._get_risk_distribution(),
                'optimization_type_distribution': self._get_optimization_type_distribution(),
                'last_analysis': datetime.now().isoformat()
            }

        except Exception as e:
            logger.error(f"Failed to get improvement stats: {e}")
            return {}

    def _get_risk_distribution(self) -> Dict[str, int]:
        """Get distribution of suggestions by risk level"""
        distribution = defaultdict(int)
        for suggestion in self.optimization_suggestions:
            risk_level = suggestion.get('risk_level', RiskLevel.MEDIUM)
            distribution[risk_level.value] += 1
        return dict(distribution)

    def _get_optimization_type_distribution(self) -> Dict[str, int]:
        """Get distribution of suggestions by optimization type"""
        distribution = defaultdict(int)
        for suggestion in self.optimization_suggestions:
            opt_type = suggestion.get('type', OptimizationType.PERFORMANCE)
            distribution[opt_type.value] += 1
        return dict(distribution)

# Factory Core Class
class AutoGenFactory:
    def __init__(self, config: FactoryConfig):
        self.config = config
        self.redis_client = None
        self.queue = None
        self.worker = None
        self.llm_config = None
        self.memory_store = {}

        # Initialize learning systems
        self.pattern_learning = PatternLearningSystem()
        self.improvement_engine = ContinuousImprovementEngine()

        self.setup_connections()
        self.setup_autogen()
        logger.info("AutoGen Factory initialized successfully")

    def setup_connections(self):
        """Setup Redis and RQ connections"""
        try:
            # Redis connection
            self.redis_client = redis.Redis(
                host=self.config.redis.host,
                port=self.config.redis.port,
                db=self.config.redis.db,
                password=self.config.redis.password,
                decode_responses=self.config.redis.decode_responses
            )

            # Test connection
            self.redis_client.ping()
            logger.info("Redis connection established")

            # RQ setup
            with Connection(self.redis_client):
                self.queue = Queue(self.config.rq.queue_name)

            logger.info("RQ queue setup completed")

        except Exception as e:
            logger.error(f"Failed to setup connections: {e}")
            raise

    def setup_autogen(self):
        """Setup AutoGen configuration"""
        api_key = self.config.autogen.api_key or os.getenv("OPENAI_API_KEY")
        if not api_key:
            raise ValueError("OpenAI API key not found")

        self.llm_config = {
            "config_list": [{
                "model": self.config.autogen.model,
                "api_key": api_key,
                "temperature": self.config.autogen.temperature,
                "max_tokens": self.config.autogen.max_tokens
            }],
            "timeout": 120,
        }

        logger.info("AutoGen configuration setup completed")

    def store_error_memory(self, error_memory: ErrorMemory):
        """Store error memory for learning"""
        try:
            key = f"error_memory:{error_memory.id}"
            data = json.dumps(asdict(error_memory), default=str)
            self.redis_client.setex(key, self.config.memory.memory_ttl, data)

            # Update frequency tracking
            freq_key = f"error_freq:{error_memory.error_type}:{error_memory.error_message}"
            self.redis_client.incr(freq_key)
            self.redis_client.expire(freq_key, self.config.memory.memory_ttl)

            logger.info(f"Stored error memory: {error_memory.id}")
            return True
        except Exception as e:
            logger.error(f"Failed to store error memory: {e}")
            return False

    def get_similar_errors(self, error_type: str, error_message: str, limit: int = 5) -> List[ErrorMemory]:
        """Retrieve similar error memories for learning"""
        try:
            # Get all error memories
            keys = self.redis_client.keys("error_memory:*")
            similar_memories = []

            for key in keys[:100]:  # Limit search space
                data = self.redis_client.get(key)
                if data:
                    memory_dict = json.loads(data)
                    memory = ErrorMemory(**memory_dict)

                    # Simple similarity check
                    if (memory.error_type == error_type and
                        self._calculate_similarity(memory.error_message, error_message) > self.config.memory.similarity_threshold):
                        similar_memories.append(memory)

            # Sort by frequency and return top matches
            similar_memories.sort(key=lambda x: x.frequency, reverse=True)
            return similar_memories[:limit]

        except Exception as e:
            logger.error(f"Failed to retrieve similar errors: {e}")
            return []

    def _calculate_similarity(self, text1: str, text2: str) -> float:
        """Simple text similarity calculation"""
        words1 = set(text1.lower().split())
        words2 = set(text2.lower().split())

        if not words1 or not words2:
            return 0.0

        intersection = words1.intersection(words2)
        union = words1.union(words2)

        return len(intersection) / len(union)

    def learn_from_errors(self, error_type: str, error_message: str, context: Dict[str, Any]) -> Optional[str]:
        """Learn from similar errors and suggest solutions"""
        try:
            similar_errors = self.get_similar_errors(error_type, error_message)

            if not similar_errors:
                # Store new error for future learning
                new_memory = ErrorMemory(
                    error_type=error_type,
                    error_message=error_message,
                    context=context
                )
                self.store_error_memory(new_memory)
                return None

            # Use AutoGen to analyze and learn
            user_proxy = autogen.UserProxyAgent(
                "user_proxy",
                code_execution_config=False,
                human_input_mode="NEVER"
            )

            assistant = autogen.AssistantAgent(
                "assistant",
                llm_config=self.llm_config,
                system_message=f"""You are an expert system that learns from past errors.
                Analyze the current error and similar past errors to suggest a solution.

                Current Error: {error_type}: {error_message}
                Context: {context}

                Similar Past Errors:
                {json.dumps([asdict(mem) for mem in similar_errors], indent=2, default=str)}

                Provide a specific solution based on patterns you observe."""
            )

            user_proxy.initiate_chat(assistant, message="Analyze the error and suggest a solution.")

            # Extract the last message from assistant
            solution = assistant.last_message()["content"]

            # Store the learning
            learned_memory = ErrorMemory(
                error_type=error_type,
                error_message=error_message,
                context=context,
                solution=solution,
                learned=True
            )
            self.store_error_memory(learned_memory)

            return solution

        except Exception as e:
            logger.error(f"Failed to learn from errors: {e}")
            return None

    def queue_deployment_task(self, task_type: str, parameters: Dict[str, Any]) -> str:
        """Queue a deployment aid task"""
        try:
            job_id = f"deploy_{uuid.uuid4().hex[:8]}"

            with Connection(self.redis_client):
                job = self.queue.enqueue(
                    "factory.tasks.execute_deployment_task",
                    task_type,
                    parameters,
                    job_id=job_id,
                    timeout=self.config.rq.job_timeout
                )

            logger.info(f"Queued deployment task: {job_id}")
            return job_id

        except Exception as e:
            logger.error(f"Failed to queue deployment task: {e}")
            raise

    def get_deployment_aid(self, task_type: str, parameters: Dict[str, Any]) -> Dict[str, Any]:
        """Get deployment aid for a specific task"""
        try:
            # Check for cached results first
            cache_key = f"deploy_aid:{task_type}:{hash(str(parameters))}"
            cached = self.redis_client.get(cache_key)

            if cached:
                return json.loads(cached)

            # Use AutoGen to provide deployment aid
            user_proxy = autogen.UserProxyAgent(
                "user_proxy",
                code_execution_config=False,
                human_input_mode="NEVER"
            )

            assistant = autogen.AssistantAgent(
                "deployment_assistant",
                llm_config=self.llm_config,
                system_message=f"""You are a deployment expert. Provide detailed guidance for: {task_type}

                Parameters: {parameters}

                Include:
                1. Step-by-step instructions
                2. Required tools and dependencies
                3. Potential issues and solutions
                4. Best practices
                5. Verification steps"""
            )

            user_proxy.initiate_chat(assistant, message=f"Provide deployment aid for {task_type}")

            # Extract and structure the response
            aid_content = assistant.last_message()["content"]

            result = {
                "task_type": task_type,
                "parameters": parameters,
                "guidance": aid_content,
                "timestamp": datetime.now().isoformat(),
                "cached": False
            }

            # Cache the result
            self.redis_client.setex(cache_key, 3600, json.dumps(result))  # Cache for 1 hour

            return result

        except Exception as e:
            logger.error(f"Failed to get deployment aid: {e}")
            raise

    def start_worker(self):
        """Start the RQ worker"""
        try:
            with Connection(self.redis_client):
                worker = Worker([self.queue], name=self.config.rq.worker_name)
                logger.info(f"Starting RQ worker: {self.config.rq.worker_name}")
                worker.work()

        except Exception as e:
            logger.error(f"Failed to start worker: {e}")
            raise

# Global factory instance
factory_instance = None

def get_factory(config_path: str = "factory_config.json") -> AutoGenFactory:
    """Get or create factory instance"""
    global factory_instance

    if factory_instance is None:
        # Load config from file or use defaults
        if os.path.exists(config_path):
            with open(config_path, 'r') as f:
                config_data = json.load(f)
            config = FactoryConfig(**config_data)
        else:
            config = FactoryConfig()

        factory_instance = AutoGenFactory(config)

    return factory_instance

# Convenience functions
def learn_from_error(error_type: str, error_message: str, context: Dict[str, Any]) -> Optional[str]:
    """Convenience function to learn from an error"""
    factory = get_factory()
    return factory.learn_from_errors(error_type, error_message, context)

def get_deployment_aid(task_type: str, **parameters) -> Dict[str, Any]:
    """Convenience function to get deployment aid"""
    factory = get_factory()
    return factory.get_deployment_aid(task_type, parameters)

def queue_task(task_type: str, **parameters) -> str:
    """Convenience function to queue a task"""
    factory = get_factory()
    return factory.queue_deployment_task(task_type, parameters)

# Pattern Learning System convenience functions
def learn_from_deployment(deployment_data: Dict[str, Any]) -> bool:
    """Learn from deployment outcome using pattern learning system"""
    factory = get_factory()
    return factory.pattern_learning.learn_from_deployment(deployment_data)

def suggest_pattern_solution(error_type: str, context: Optional[Dict[str, Any]] = None) -> Optional[str]:
    """Suggest solution based on learned patterns"""
    factory = get_factory()
    return factory.pattern_learning.suggest_solution(error_type, context)

def get_learning_stats() -> Dict[str, Any]:
    """Get pattern learning statistics"""
    factory = get_factory()
    return factory.pattern_learning.get_learning_stats()

# Continuous Improvement Engine convenience functions
def analyze_system_performance() -> List[Dict[str, Any]]:
    """Analyze system performance and generate improvement suggestions"""
    factory = get_factory()
    return factory.improvement_engine.analyze_performance()

def implement_safe_improvements() -> Dict[str, Any]:
    """Implement low-risk system improvements"""
    factory = get_factory()
    return factory.improvement_engine.implement_improvements()

def get_improvement_stats() -> Dict[str, Any]:
    """Get continuous improvement statistics"""
    factory = get_factory()
    return factory.improvement_engine.get_improvement_stats()

def collect_system_metrics() -> Dict[str, Any]:
    """Collect current system performance metrics"""
    factory = get_factory()
    return factory.improvement_engine.collect_metrics()
