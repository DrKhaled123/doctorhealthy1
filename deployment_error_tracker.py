#!/usr/bin/env python3
"""
Error Tracking & Memory System for Future Error Prevention
Implements systematic error tracking with best practices from deployment guide
"""

import json
import os
from datetime import datetime
from dataclasses import dataclass, asdict
from typing import List, Dict, Optional
from pathlib import Path

@dataclass
class ErrorRecord:
    """Structured error record for tracking deployment issues"""
    date: str
    platform: str  # Coolify, Hetzner, Docker, etc.
    error_type: str  # 404, 502, deployment failure, etc.
    root_cause: str
    solution_applied: Dict[str, any]
    prevention_measures: Dict[str, any]
    related_resources: List[str]
    severity: str  # critical, high, medium, low
    resolved: bool = False
    follow_up_date: Optional[str] = None

class DeploymentErrorTracker:
    """Comprehensive error tracking and prevention system"""
    
    def __init__(self, tracker_file: str = "deployment_error_tracker.json"):
        self.tracker_file = tracker_file
        self.errors: List[ErrorRecord] = []
        self.load_existing_errors()
        
    def load_existing_errors(self):
        """Load existing error records from file"""
        if os.path.exists(self.tracker_file):
            try:
                with open(self.tracker_file, 'r') as f:
                    data = json.load(f)
                
                self.errors = []
                for error_data in data.get('errors', []):
                    error_record = ErrorRecord(**error_data)
                    self.errors.append(error_record)
                    
                print(f"📚 Loaded {len(self.errors)} existing error records")
            except Exception as e:
                print(f"⚠️  Could not load existing error tracker: {e}")
                self.errors = []
    
    def add_error(self, 
                  platform: str,
                  error_type: str, 
                  root_cause: str,
                  severity: str = "medium",
                  steps_taken: List[str] = None,
                  config_changes: Dict[str, str] = None,
                  commands_used: List[str] = None,
                  documentation_updates: List[str] = None,
                  monitoring_added: List[str] = None,
                  process_changes: List[str] = None,
                  related_resources: List[str] = None) -> ErrorRecord:
        """Add a new error record with comprehensive tracking"""
        
        error_record = ErrorRecord(
            date=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            platform=platform,
            error_type=error_type,
            root_cause=root_cause,
            severity=severity,
            solution_applied={
                "steps_taken": steps_taken or [],
                "configuration_changes": config_changes or {},
                "commands_used": commands_used or []
            },
            prevention_measures={
                "documentation_updates": documentation_updates or [],
                "monitoring_added": monitoring_added or [],
                "process_changes": process_changes or []
            },
            related_resources=related_resources or [],
            resolved=False
        )
        
        self.errors.append(error_record)
        self.save_errors()
        
        print(f"📝 Added new error record: {error_type} on {platform}")
        
        return error_record
    
    def mark_resolved(self, error_index: int, follow_up_date: str = None):
        """Mark an error as resolved with optional follow-up"""
        if 0 <= error_index < len(self.errors):
            self.errors[error_index].resolved = True
            if follow_up_date:
                self.errors[error_index].follow_up_date = follow_up_date
            
            self.save_errors()
            print(f"✅ Marked error #{error_index} as resolved")
    
    def get_unresolved_errors(self) -> List[ErrorRecord]:
        """Get list of unresolved errors"""
        return [error for error in self.errors if not error.resolved]
    
    def get_errors_by_platform(self, platform: str) -> List[ErrorRecord]:
        """Get errors filtered by platform"""
        return [error for error in self.errors if error.platform.lower() == platform.lower()]
    
    def get_errors_by_type(self, error_type: str) -> List[ErrorRecord]:
        """Get errors filtered by type"""
        return [error for error in self.errors if error_type.lower() in error.error_type.lower()]
    
    def save_errors(self):
        """Save all error records to file"""
        data = {
            "metadata": {
                "last_updated": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                "total_errors": len(self.errors),
                "resolved_errors": len([e for e in self.errors if e.resolved]),
                "platforms": list(set([e.platform for e in self.errors]))
            },
            "errors": [asdict(error) for error in self.errors]
        }
        
        with open(self.tracker_file, 'w') as f:
            json.dump(data, f, indent=2)
    
    def generate_monthly_report(self) -> Dict:
        """Generate monthly error analysis report"""
        current_month = datetime.now().strftime("%Y-%m")
        
        monthly_errors = [
            error for error in self.errors 
            if error.date.startswith(current_month)
        ]
        
        # Platform breakdown
        platform_counts = {}
        for error in monthly_errors:
            platform_counts[error.platform] = platform_counts.get(error.platform, 0) + 1
        
        # Error type breakdown
        type_counts = {}
        for error in monthly_errors:
            type_counts[error.error_type] = type_counts.get(error.error_type, 0) + 1
        
        # Severity breakdown
        severity_counts = {}
        for error in monthly_errors:
            severity_counts[error.severity] = severity_counts.get(error.severity, 0) + 1
        
        report = {
            "month": current_month,
            "total_errors": len(monthly_errors),
            "resolved_errors": len([e for e in monthly_errors if e.resolved]),
            "platform_breakdown": platform_counts,
            "error_type_breakdown": type_counts,
            "severity_breakdown": severity_counts,
            "top_issues": self._get_top_issues(monthly_errors)
        }
        
        return report
    
    def _get_top_issues(self, errors: List[ErrorRecord], limit: int = 5) -> List[Dict]:
        """Get top recurring issues"""
        issue_patterns = {}
        
        for error in errors:
            key = f"{error.platform}:{error.error_type}"
            if key not in issue_patterns:
                issue_patterns[key] = {
                    "platform": error.platform,
                    "error_type": error.error_type,
                    "count": 0,
                    "latest_date": error.date,
                    "resolved_count": 0
                }
            
            issue_patterns[key]["count"] += 1
            if error.resolved:
                issue_patterns[key]["resolved_count"] += 1
            
            if error.date > issue_patterns[key]["latest_date"]:
                issue_patterns[key]["latest_date"] = error.date
        
        # Sort by count descending
        top_issues = sorted(
            issue_patterns.values(), 
            key=lambda x: x["count"], 
            reverse=True
        )
        
        return top_issues[:limit]
    
    def suggest_prevention_measures(self, error_type: str, platform: str) -> List[str]:
        """Suggest prevention measures based on historical data"""
        similar_errors = [
            error for error in self.errors 
            if (error.error_type.lower() == error_type.lower() or 
                error.platform.lower() == platform.lower()) and 
            error.resolved
        ]
        
        suggestions = []
        
        # Extract successful prevention measures
        for error in similar_errors:
            prevention = error.prevention_measures
            
            for category, measures in prevention.items():
                if measures:
                    suggestions.extend([
                        f"[{category}] {measure}" for measure in measures
                    ])
        
        # Remove duplicates and return top suggestions
        unique_suggestions = list(set(suggestions))
        return unique_suggestions[:10]  # Top 10 suggestions
    
    def print_summary(self):
        """Print a summary of the error tracking system"""
        total_errors = len(self.errors)
        resolved_errors = len([e for e in self.errors if e.resolved])
        unresolved_errors = total_errors - resolved_errors
        
        print("📊 Error Tracking System Summary")
        print("=" * 50)
        print(f"Total Errors Tracked: {total_errors}")
        print(f"Resolved: {resolved_errors}")
        print(f"Unresolved: {unresolved_errors}")
        
        if total_errors > 0:
            resolution_rate = (resolved_errors / total_errors) * 100
            print(f"Resolution Rate: {resolution_rate:.1f}%")
        
        # Show platform breakdown
        platforms = {}
        for error in self.errors:
            platforms[error.platform] = platforms.get(error.platform, 0) + 1
        
        if platforms:
            print(f"\nPlatform Breakdown:")
            for platform, count in sorted(platforms.items(), key=lambda x: x[1], reverse=True):
                print(f"  {platform}: {count} errors")
        
        # Show recent unresolved errors
        recent_unresolved = [
            error for error in self.errors 
            if not error.resolved
        ][-5:]  # Last 5 unresolved
        
        if recent_unresolved:
            print(f"\nRecent Unresolved Issues:")
            for i, error in enumerate(recent_unresolved, 1):
                print(f"  {i}. [{error.platform}] {error.error_type} - {error.date}")

def add_common_deployment_errors():
    """Pre-populate with common deployment errors and solutions"""
    tracker = DeploymentErrorTracker()
    
    # Add some common deployment errors with solutions
    common_errors = [
        {
            "platform": "Nginx",
            "error_type": "502 Bad Gateway", 
            "root_cause": "Proxy misconfiguration or backend service down",
            "severity": "high",
            "steps_taken": [
                "Check nginx configuration syntax with 'nginx -t'",
                "Verify backend service is running",
                "Check proxy_pass directive points to correct port",
                "Restart nginx service"
            ],
            "config_changes": {
                "nginx.conf": "Added proper proxy headers and timeout settings"
            },
            "commands_used": [
                "nginx -t",
                "systemctl restart nginx",
                "netstat -tuln | grep :8080"
            ]
        },
        {
            "platform": "Docker",
            "error_type": "Container Build Failure",
            "root_cause": "Missing dependencies or incorrect Dockerfile syntax",
            "severity": "medium",
            "steps_taken": [
                "Review Dockerfile for syntax errors",
                "Check base image availability",
                "Verify COPY paths are correct",
                "Test build with --no-cache flag"
            ],
            "config_changes": {
                "Dockerfile": "Fixed COPY paths and added proper USER directive"
            },
            "commands_used": [
                "docker build --no-cache -t app .",
                "docker run --rm app npm test"
            ]
        },
        {
            "platform": "Coolify",
            "error_type": "Deployment Timeout",
            "root_cause": "Build process taking too long or resource limits exceeded",
            "severity": "medium",
            "steps_taken": [
                "Check build logs for bottlenecks",
                "Optimize build process",
                "Increase resource limits",
                "Split build into stages"
            ],
            "monitoring_added": [
                "Build time monitoring",
                "Resource usage alerts"
            ]
        }
    ]
    
    for error_data in common_errors:
        tracker.add_error(**error_data)
    
    return tracker

if __name__ == "__main__":
    # Example usage and setup
    print("🚀 Initializing Deployment Error Tracking System")
    
    # Load or create tracker
    tracker = DeploymentErrorTracker()
    
    # If no errors exist, add common ones
    if len(tracker.errors) == 0:
        print("📚 Adding common deployment errors to knowledge base...")
        tracker = add_common_deployment_errors()
    
    # Print summary
    tracker.print_summary()
    
    # Generate monthly report
    monthly_report = tracker.generate_monthly_report()
    
    print(f"\n📈 Monthly Report for {monthly_report['month']}:")
    print(f"  Total errors this month: {monthly_report['total_errors']}")
    print(f"  Resolved this month: {monthly_report['resolved_errors']}")
    
    if monthly_report['top_issues']:
        print(f"  Top issues:")
        for issue in monthly_report['top_issues'][:3]:
            print(f"    - {issue['platform']}: {issue['error_type']} ({issue['count']} times)")
    
    # Example: Get prevention suggestions
    suggestions = tracker.suggest_prevention_measures("502 Bad Gateway", "Nginx")
    if suggestions:
        print(f"\n💡 Prevention suggestions for Nginx 502 errors:")
        for suggestion in suggestions[:3]:
            print(f"  - {suggestion}")
    
    print(f"\n✅ Error tracking system ready!")
    print(f"📄 Data saved to: {tracker.tracker_file}")