#!/usr/bin/env python3
"""
Comprehensive Monitoring & Alerting System for Deployment Health
Implements continuous monitoring with proactive alerts and performance tracking
"""

import asyncio
import json
import subprocess
import time
import os
import sys
from datetime import datetime, timedelta
from dataclasses import dataclass, asdict
from typing import Dict, List, Optional, Callable
import threading
import queue
from pathlib import Path

@dataclass
class Alert:
    """Alert structure for monitoring system"""
    timestamp: str
    severity: str  # critical, high, medium, low
    category: str  # performance, security, availability, resource
    message: str
    metric_name: str
    current_value: float
    threshold_value: float
    suggested_action: str
    resolved: bool = False

@dataclass
class HealthMetric:
    """Health metric data structure"""
    name: str
    current_value: float
    threshold: float
    unit: str
    status: str  # healthy, warning, critical
    last_updated: str
    trend: str  # improving, stable, degrading

class DeploymentMonitor:
    """Comprehensive deployment health monitoring system"""
    
    def __init__(self, config_file: str = "monitoring_config.json"):
        self.config_file = config_file
        self.monitoring_active = False
        self.alerts_queue = queue.Queue()
        self.metrics_history = {}
        self.current_metrics = {}
        self.alert_callbacks = {}
        self.load_config()
        
    def load_config(self):
        """Load monitoring configuration"""
        default_config = {
            "monitoring_interval": 30,  # seconds
            "retention_days": 7,
            "thresholds": {
                "response_time_ms": 2000,
                "cpu_usage_percent": 80,
                "memory_usage_percent": 85,
                "disk_usage_percent": 90,
                "error_rate_percent": 5,
                "request_rate_per_minute": 1000,
                "concurrent_connections": 100
            },
            "endpoints_to_monitor": [
                "http://localhost:8085/health",
                "http://localhost:8085/ready",
                "http://localhost:8085/metrics"
            ],
            "alert_rules": {
                "consecutive_failures": 3,
                "response_time_samples": 5,
                "escalation_minutes": 15
            }
        }
        
        if os.path.exists(self.config_file):
            try:
                with open(self.config_file, 'r') as f:
                    loaded_config = json.load(f)
                    # Merge with defaults
                    self.config = {**default_config, **loaded_config}
            except Exception as e:
                print(f"⚠️  Could not load config: {e}, using defaults")
                self.config = default_config
        else:
            self.config = default_config
            self.save_config()
    
    def save_config(self):
        """Save current configuration"""
        with open(self.config_file, 'w') as f:
            json.dump(self.config, f, indent=2)
    
    async def start_monitoring(self):
        """Start continuous monitoring"""
        print("🔍 Starting Deployment Health Monitoring...")
        print("=" * 50)
        
        self.monitoring_active = True
        
        # Start monitoring tasks
        tasks = [
            asyncio.create_task(self._monitor_application_health()),
            asyncio.create_task(self._monitor_system_resources()),
            asyncio.create_task(self._monitor_performance_metrics()),
            asyncio.create_task(self._process_alerts()),
            asyncio.create_task(self._periodic_health_check())
        ]
        
        try:
            await asyncio.gather(*tasks)
        except KeyboardInterrupt:
            print("\n⚠️  Monitoring interrupted by user")
            self.monitoring_active = False
        except Exception as e:
            print(f"\n💥 Monitoring error: {e}")
            self.monitoring_active = False
    
    async def _monitor_application_health(self):
        """Monitor application health endpoints"""
        while self.monitoring_active:
            for endpoint in self.config["endpoints_to_monitor"]:
                try:
                    # Measure response time
                    start_time = time.time()
                    
                    cmd = ["curl", "-s", "-w", "%{http_code},%{time_total}", "-o", "/dev/null", endpoint]
                    result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
                    
                    response_time_ms = (time.time() - start_time) * 1000
                    
                    if result.returncode == 0:
                        try:
                            status_code, curl_time = result.stdout.strip().split(',')
                            status_code = int(status_code)
                            curl_time_ms = float(curl_time) * 1000
                            
                            # Use more accurate curl timing if available
                            if curl_time_ms > 0:
                                response_time_ms = curl_time_ms
                            
                            # Create health metric
                            self._record_metric(
                                name=f"response_time_{endpoint.split('/')[-1]}",
                                value=response_time_ms,
                                unit="ms",
                                threshold=self.config["thresholds"]["response_time_ms"]
                            )
                            
                            self._record_metric(
                                name=f"status_code_{endpoint.split('/')[-1]}",
                                value=status_code,
                                unit="code",
                                threshold=400  # HTTP error threshold
                            )
                            
                            # Check for alerts
                            if status_code >= 400:
                                await self._create_alert(
                                    severity="high" if status_code >= 500 else "medium",
                                    category="availability",
                                    message=f"HTTP error {status_code} from {endpoint}",
                                    metric_name=f"status_code_{endpoint.split('/')[-1]}",
                                    current_value=status_code,
                                    threshold_value=400,
                                    suggested_action="Check application logs and service status"
                                )
                            
                            if response_time_ms > self.config["thresholds"]["response_time_ms"]:
                                await self._create_alert(
                                    severity="medium",
                                    category="performance",
                                    message=f"Slow response from {endpoint}: {response_time_ms:.1f}ms",
                                    metric_name=f"response_time_{endpoint.split('/')[-1]}",
                                    current_value=response_time_ms,
                                    threshold_value=self.config["thresholds"]["response_time_ms"],
                                    suggested_action="Check application performance and resource usage"
                                )
                        
                        except ValueError:
                            # Could not parse curl output
                            await self._create_alert(
                                severity="medium",
                                category="availability",
                                message=f"Could not parse response from {endpoint}",
                                metric_name=f"response_parse_{endpoint.split('/')[-1]}",
                                current_value=0,
                                threshold_value=1,
                                suggested_action="Check endpoint response format"
                            )
                    
                    else:
                        # Connection failed
                        await self._create_alert(
                            severity="critical",
                            category="availability",
                            message=f"Cannot connect to {endpoint}",
                            metric_name=f"connectivity_{endpoint.split('/')[-1]}",
                            current_value=0,
                            threshold_value=1,
                            suggested_action="Check if application is running and accessible"
                        )
                
                except asyncio.TimeoutError:
                    await self._create_alert(
                        severity="high",
                        category="performance",
                        message=f"Timeout connecting to {endpoint}",
                        metric_name=f"timeout_{endpoint.split('/')[-1]}",
                        current_value=10000,  # 10 seconds in ms
                        threshold_value=self.config["thresholds"]["response_time_ms"],
                        suggested_action="Check network connectivity and application responsiveness"
                    )
                
                except Exception as e:
                    await self._create_alert(
                        severity="medium",
                        category="monitoring",
                        message=f"Monitoring error for {endpoint}: {str(e)[:100]}",
                        metric_name="monitoring_error",
                        current_value=1,
                        threshold_value=0,
                        suggested_action="Check monitoring system configuration"
                    )
            
            await asyncio.sleep(self.config["monitoring_interval"])
    
    async def _monitor_system_resources(self):
        """Monitor system resource usage"""
        while self.monitoring_active:
            try:
                # CPU Usage
                cpu_usage = await self._get_cpu_usage()
                if cpu_usage is not None:
                    self._record_metric("cpu_usage", cpu_usage, "%", self.config["thresholds"]["cpu_usage_percent"])
                    
                    if cpu_usage > self.config["thresholds"]["cpu_usage_percent"]:
                        await self._create_alert(
                            severity="high",
                            category="resource",
                            message=f"High CPU usage: {cpu_usage:.1f}%",
                            metric_name="cpu_usage",
                            current_value=cpu_usage,
                            threshold_value=self.config["thresholds"]["cpu_usage_percent"],
                            suggested_action="Investigate high CPU processes and optimize performance"
                        )
                
                # Memory Usage
                memory_usage = await self._get_memory_usage()
                if memory_usage is not None:
                    self._record_metric("memory_usage", memory_usage, "%", self.config["thresholds"]["memory_usage_percent"])
                    
                    if memory_usage > self.config["thresholds"]["memory_usage_percent"]:
                        await self._create_alert(
                            severity="high",
                            category="resource",
                            message=f"High memory usage: {memory_usage:.1f}%",
                            metric_name="memory_usage",
                            current_value=memory_usage,
                            threshold_value=self.config["thresholds"]["memory_usage_percent"],
                            suggested_action="Check for memory leaks and optimize memory usage"
                        )
                
                # Disk Usage
                disk_usage = await self._get_disk_usage()
                if disk_usage is not None:
                    self._record_metric("disk_usage", disk_usage, "%", self.config["thresholds"]["disk_usage_percent"])
                    
                    if disk_usage > self.config["thresholds"]["disk_usage_percent"]:
                        await self._create_alert(
                            severity="critical",
                            category="resource",
                            message=f"High disk usage: {disk_usage:.1f}%",
                            metric_name="disk_usage",
                            current_value=disk_usage,
                            threshold_value=self.config["thresholds"]["disk_usage_percent"],
                            suggested_action="Free up disk space immediately to prevent service disruption"
                        )
                
                # Process Count for Application
                process_count = await self._get_process_count()
                if process_count is not None:
                    self._record_metric("process_count", process_count, "processes", 10)  # Alert if >10 processes
            
            except Exception as e:
                print(f"⚠️  Error monitoring system resources: {e}")
            
            await asyncio.sleep(self.config["monitoring_interval"] * 2)  # Check less frequently
    
    async def _monitor_performance_metrics(self):
        """Monitor application performance metrics"""
        request_counts = {}
        
        while self.monitoring_active:
            try:
                # Check application logs for error patterns (if available)
                error_rate = await self._calculate_error_rate()
                if error_rate is not None:
                    self._record_metric("error_rate", error_rate, "%", self.config["thresholds"]["error_rate_percent"])
                    
                    if error_rate > self.config["thresholds"]["error_rate_percent"]:
                        await self._create_alert(
                            severity="high",
                            category="performance",
                            message=f"High error rate: {error_rate:.1f}%",
                            metric_name="error_rate",
                            current_value=error_rate,
                            threshold_value=self.config["thresholds"]["error_rate_percent"],
                            suggested_action="Check application logs for errors and fix issues"
                        )
                
                # Monitor concurrent connections (approximate)
                connection_count = await self._get_connection_count()
                if connection_count is not None:
                    self._record_metric("connections", connection_count, "conn", self.config["thresholds"]["concurrent_connections"])
                    
                    if connection_count > self.config["thresholds"]["concurrent_connections"]:
                        await self._create_alert(
                            severity="medium",
                            category="performance",
                            message=f"High connection count: {connection_count}",
                            metric_name="connections",
                            current_value=connection_count,
                            threshold_value=self.config["thresholds"]["concurrent_connections"],
                            suggested_action="Monitor traffic patterns and consider scaling"
                        )
            
            except Exception as e:
                print(f"⚠️  Error monitoring performance metrics: {e}")
            
            await asyncio.sleep(self.config["monitoring_interval"])
    
    async def _process_alerts(self):
        """Process and handle alerts"""
        while self.monitoring_active:
            try:
                # Check for alerts in queue (non-blocking)
                try:
                    alert = self.alerts_queue.get_nowait()
                    await self._handle_alert(alert)
                except queue.Empty:
                    pass
            except Exception as e:
                print(f"⚠️  Error processing alerts: {e}")
            
            await asyncio.sleep(1)  # Check alerts frequently
    
    async def _periodic_health_check(self):
        """Periodic comprehensive health check"""
        while self.monitoring_active:
            try:
                print(f"\n🩺 Health Check - {datetime.now().strftime('%H:%M:%S')}")
                
                # Show current metrics
                healthy_metrics = 0
                warning_metrics = 0
                critical_metrics = 0
                
                for metric_name, metric in self.current_metrics.items():
                    if metric.status == "healthy":
                        healthy_metrics += 1
                    elif metric.status == "warning":
                        warning_metrics += 1
                    else:
                        critical_metrics += 1
                
                total_metrics = len(self.current_metrics)
                
                if total_metrics > 0:
                    health_score = (healthy_metrics / total_metrics) * 100
                    
                    status_emoji = "✅" if health_score >= 90 else "⚠️" if health_score >= 70 else "❌"
                    
                    print(f"{status_emoji} Health Score: {health_score:.1f}% ({healthy_metrics}/{total_metrics} healthy)")
                    
                    if warning_metrics > 0:
                        print(f"⚠️  {warning_metrics} metrics in warning state")
                    if critical_metrics > 0:
                        print(f"❌ {critical_metrics} metrics in critical state")
                
                # Show recent alerts
                recent_alerts = [alert for alert in list(self.alerts_queue.queue) if not alert.resolved]
                if recent_alerts:
                    print(f"🚨 {len(recent_alerts)} active alerts")
            
            except Exception as e:
                print(f"⚠️  Error in health check: {e}")
            
            await asyncio.sleep(60)  # Health check every minute
    
    def _record_metric(self, name: str, value: float, unit: str, threshold: float):
        """Record a metric value"""
        # Determine status
        if value <= threshold * 0.8:  # 80% of threshold
            status = "healthy"
        elif value <= threshold:
            status = "warning"  
        else:
            status = "critical"
        
        # Calculate trend (simplified)
        trend = "stable"
        if name in self.current_metrics:
            old_value = self.current_metrics[name].current_value
            if value > old_value * 1.1:  # 10% increase
                trend = "degrading"
            elif value < old_value * 0.9:  # 10% decrease
                trend = "improving"
        
        metric = HealthMetric(
            name=name,
            current_value=value,
            threshold=threshold,
            unit=unit,
            status=status,
            last_updated=datetime.now().strftime("%H:%M:%S"),
            trend=trend
        )
        
        self.current_metrics[name] = metric
        
        # Store in history
        if name not in self.metrics_history:
            self.metrics_history[name] = []
        
        self.metrics_history[name].append({
            "timestamp": datetime.now().isoformat(),
            "value": value,
            "status": status
        })
        
        # Limit history size
        if len(self.metrics_history[name]) > 1000:
            self.metrics_history[name] = self.metrics_history[name][-500:]
    
    async def _create_alert(self, severity: str, category: str, message: str, 
                           metric_name: str, current_value: float, threshold_value: float, 
                           suggested_action: str):
        """Create and queue an alert"""
        alert = Alert(
            timestamp=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            severity=severity,
            category=category,
            message=message,
            metric_name=metric_name,
            current_value=current_value,
            threshold_value=threshold_value,
            suggested_action=suggested_action
        )
        
        self.alerts_queue.put(alert)
    
    async def _handle_alert(self, alert: Alert):
        """Handle an alert"""
        severity_emoji = {
            "critical": "🔴",
            "high": "🟠", 
            "medium": "🟡",
            "low": "🟢"
        }
        
        emoji = severity_emoji.get(alert.severity, "⚪")
        
        print(f"\n{emoji} ALERT [{alert.severity.upper()}] - {alert.category}")
        print(f"📍 {alert.message}")
        print(f"📊 {alert.metric_name}: {alert.current_value} (threshold: {alert.threshold_value})")
        print(f"💡 Action: {alert.suggested_action}")
        print(f"🕐 {alert.timestamp}")
        
        # Log to file
        self._log_alert(alert)
    
    def _log_alert(self, alert: Alert):
        """Log alert to file"""
        log_file = f"alerts_{datetime.now().strftime('%Y%m%d')}.log"
        
        try:
            with open(log_file, 'a') as f:
                log_entry = {
                    "timestamp": alert.timestamp,
                    "severity": alert.severity,
                    "category": alert.category,
                    "message": alert.message,
                    "metric": alert.metric_name,
                    "value": alert.current_value,
                    "threshold": alert.threshold_value,
                    "action": alert.suggested_action
                }
                f.write(json.dumps(log_entry) + "\n")
        except Exception as e:
            print(f"⚠️  Could not log alert: {e}")
    
    async def _get_cpu_usage(self) -> Optional[float]:
        """Get CPU usage percentage"""
        try:
            # Use top command to get CPU usage
            result = subprocess.run(
                ["top", "-l", "1", "-n", "0"], 
                capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                lines = result.stdout.split('\n')
                for line in lines:
                    if 'CPU usage' in line:
                        # Parse macOS top output: "CPU usage: 10.5% user, 5.2% sys, 84.3% idle"
                        parts = line.split()
                        for i, part in enumerate(parts):
                            if part == 'idle':
                                idle_percent = float(parts[i-1].rstrip('%'))
                                return 100 - idle_percent
            return None
        except Exception:
            return None
    
    async def _get_memory_usage(self) -> Optional[float]:
        """Get memory usage percentage"""
        try:
            result = subprocess.run(
                ["vm_stat"], capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                # Parse vm_stat output for macOS
                lines = result.stdout.split('\n')
                page_size = 4096  # Default page size
                
                free_pages = 0
                inactive_pages = 0
                total_pages = 0
                
                for line in lines:
                    if 'page size of' in line:
                        page_size = int(line.split()[-2])
                    elif 'Pages free:' in line:
                        free_pages = int(line.split()[-1].rstrip('.'))
                    elif 'Pages inactive:' in line:
                        inactive_pages = int(line.split()[-1].rstrip('.'))
                
                # Get total memory
                total_result = subprocess.run(
                    ["sysctl", "-n", "hw.memsize"], 
                    capture_output=True, text=True, timeout=5
                )
                
                if total_result.returncode == 0:
                    total_memory = int(total_result.stdout.strip())
                    total_pages = total_memory // page_size
                    
                    available_pages = free_pages + inactive_pages
                    used_pages = total_pages - available_pages
                    
                    if total_pages > 0:
                        return (used_pages / total_pages) * 100
            
            return None
        except Exception:
            return None
    
    async def _get_disk_usage(self) -> Optional[float]:
        """Get disk usage percentage"""
        try:
            result = subprocess.run(
                ["df", "-h", "."], capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                lines = result.stdout.strip().split('\n')
                if len(lines) > 1:
                    parts = lines[1].split()
                    if len(parts) >= 5:
                        usage_str = parts[4].rstrip('%')
                        return float(usage_str)
            
            return None
        except Exception:
            return None
    
    async def _get_process_count(self) -> Optional[int]:
        """Get count of application processes"""
        try:
            result = subprocess.run(
                ["pgrep", "-c", "api-key-generator"], 
                capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                return int(result.stdout.strip())
            else:
                return 0  # No processes found
        except Exception:
            return None
    
    async def _calculate_error_rate(self) -> Optional[float]:
        """Calculate error rate from recent requests (simplified)"""
        # This is a simplified implementation
        # In a real system, you'd analyze application logs
        try:
            # Check health endpoint for errors
            cmd = ["curl", "-s", "-w", "%{http_code}", "-o", "/dev/null", "http://localhost:8085/health"]
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=5)
            
            if result.returncode == 0:
                status_code = int(result.stdout.strip())
                return 0.0 if status_code < 400 else 100.0  # Simple: 0% or 100%
            else:
                return 100.0  # Connection failed = 100% error rate
        except Exception:
            return None
    
    async def _get_connection_count(self) -> Optional[int]:
        """Get current connection count"""
        try:
            result = subprocess.run(
                ["netstat", "-an", "|", "grep", ":8085", "|", "grep", "ESTABLISHED", "|", "wc", "-l"],
                shell=True, capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                return int(result.stdout.strip())
            else:
                return 0
        except Exception:
            return None

async def main():
    """Main monitoring function"""
    print("🚀 Starting Comprehensive Deployment Monitoring System")
    print("=" * 60)
    
    monitor = DeploymentMonitor()
    
    try:
        await monitor.start_monitoring()
    except KeyboardInterrupt:
        print("\n👋 Monitoring stopped by user")
    except Exception as e:
        print(f"\n💥 Monitoring failed: {e}")

if __name__ == "__main__":
    asyncio.run(main())