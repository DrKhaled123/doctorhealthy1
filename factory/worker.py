#!/usr/bin/env python3
"""
AutoGen Factory RQ Worker
Background worker for processing deployment and error learning tasks
"""

import os
import sys
import redis
from rq import Worker, Queue, Connection
from loguru import logger
import argparse
import signal
import time
from typing import List

# Add the factory directory to the path
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)

from factory_config import FactoryConfig
from tasks import execute_deployment_task


class GracefulWorker:
    """RQ Worker with graceful shutdown handling"""

    def __init__(self, config_path: str = "factory_config.json"):
        self.config = self._load_config(config_path)
        self.redis_client = None
        self.worker = None
        self.queue = None
        self.shutdown_requested = False

        # Setup signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, self._signal_handler)
        signal.signal(signal.SIGTERM, self._signal_handler)

    def _load_config(self, config_path: str) -> FactoryConfig:
        """Load configuration from file or use defaults"""
        if os.path.exists(config_path):
            import json
            with open(config_path, 'r') as f:
                config_data = json.load(f)
            return FactoryConfig(**config_data)
        else:
            return FactoryConfig()

    def _signal_handler(self, signum, frame):
        """Handle shutdown signals"""
        logger.info(f"Received signal {signum}. Requesting graceful shutdown...")
        self.shutdown_requested = True

        if self.worker:
            self.worker.request_force_stop()

    def setup_connections(self):
        """Setup Redis connection"""
        try:
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

            # Setup queue
            with Connection(self.redis_client):
                self.queue = Queue(self.config.rq.queue_name)

        except Exception as e:
            logger.error(f"Failed to setup connections: {e}")
            raise

    def start_worker(self, queue_names: List[str] = None):
        """Start the RQ worker"""
        try:
            if not queue_names:
                queue_names = [self.config.rq.queue_name]

            with Connection(self.redis_client):
                # Create worker with specified queues
                queues = [Queue(name) for name in queue_names]

                self.worker = Worker(
                    queues,
                    name=self.config.rq.worker_name,
                    connection=self.redis_client
                )

                logger.info(f"Starting RQ worker: {self.config.rq.worker_name}")
                logger.info(f"Listening to queues: {', '.join(queue_names)}")

                # Start working
                self.worker.work(
                    with_scheduler=False,
                    logging_level=self.config.debug and "DEBUG" or "INFO"
                )

        except Exception as e:
            logger.error(f"Failed to start worker: {e}")
            raise

    def run(self, queue_names: List[str] = None):
        """Main worker run loop"""
        logger.info("AutoGen Factory Worker starting up...")

        try:
            self.setup_connections()

            # Keep running until shutdown is requested
            while not self.shutdown_requested:
                try:
                    self.start_worker(queue_names)
                except KeyboardInterrupt:
                    break
                except Exception as e:
                    logger.error(f"Worker error: {e}")
                    if not self.shutdown_requested:
                        logger.info("Restarting worker in 5 seconds...")
                        time.sleep(5)

        except Exception as e:
            logger.error(f"Fatal worker error: {e}")
        finally:
            logger.info("AutoGen Factory Worker shutting down gracefully")


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description="AutoGen Factory RQ Worker")
    parser.add_argument(
        "--config",
        default="factory_config.json",
        help="Path to configuration file"
    )
    parser.add_argument(
        "--queues",
        nargs="+",
        help="Queue names to listen to (default: factory_tasks)"
    )
    parser.add_argument(
        "--log-level",
        default="INFO",
        choices=["DEBUG", "INFO", "WARNING", "ERROR"],
        help="Logging level"
    )
    parser.add_argument(
        "--worker-name",
        help="Custom worker name"
    )

    args = parser.parse_args()

    # Setup logging
    logger.remove()
    logger.add(
        sys.stdout,
        level=args.log_level,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>"
    )

    try:
        worker = GracefulWorker(args.config)

        # Override config with command line arguments
        if args.worker_name:
            worker.config.rq.worker_name = args.worker_name

        # Run worker
        worker.run(args.queues)

    except Exception as e:
        logger.error(f"Worker failed to start: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
