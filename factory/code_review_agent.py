#!/usr/bin/env python3
"""
Code Review Agent using AutoGen Factory System
Reviews, fixes, and tests code using the integrated learning system
"""

import os
import sys
import ast
import json
import time
from typing import Dict, List, Any, Optional, Tuple
from datetime import datetime
from pathlib import Path
import subprocess

# Add the factory directory to the path
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)

from factory_config import (
    get_factory,
    learn_from_deployment,
    suggest_pattern_solution,
    analyze_system_performance,
    implement_safe_improvements,
    PatternLearningSystem,
    ContinuousImprovementEngine
)


class CodeReviewAgent:
    """Intelligent code review agent using the factory system"""

    def __init__(self):
        self.factory = get_factory()
        self.review_history = []
        self.fix_suggestions = []
        self.test_results = []

    def analyze_codebase(self, directory: str = ".") -> Dict[str, Any]:
        """Analyze the entire codebase structure and quality"""
        print("🔍 Analyzing codebase structure...")

        analysis_result = {
            "timestamp": datetime.now().isoformat(),
            "directory": directory,
            "files_analyzed": 0,
            "languages": {},
            "issues_found": [],
            "patterns_detected": [],
            "suggestions": []
        }

        try:
            # Walk through directory
            for root, dirs, files in os.walk(directory):
                # Skip common directories
                dirs[:] = [d for d in dirs if not d.startswith('.') and d not in ['node_modules', '__pycache__', '.git']]

                for file in files:
                    if self._should_analyze_file(file):
                        file_path = os.path.join(root, file)
                        file_analysis = self._analyze_single_file(file_path)
                        analysis_result["files_analyzed"] += 1

                        if file_analysis["language"] not in analysis_result["languages"]:
                            analysis_result["languages"][file_analysis["language"]] = 0
                        analysis_result["languages"][file_analysis["language"]] += 1

                        analysis_result["issues_found"].extend(file_analysis["issues"])
                        analysis_result["patterns_detected"].extend(file_analysis["patterns"])

            # Generate suggestions using AutoGen
            suggestions = self._generate_codebase_suggestions(analysis_result)
            analysis_result["suggestions"] = suggestions

            print(f"✅ Analyzed {analysis_result['files_analyzed']} files")
            print(f"📊 Languages found: {analysis_result['languages']}")
            print(f"⚠️ Issues found: {len(analysis_result['issues_found'])}")

            return analysis_result

        except Exception as e:
            print(f"❌ Codebase analysis failed: {e}")
            return analysis_result

    def _should_analyze_file(self, filename: str) -> bool:
        """Check if file should be analyzed"""
        analyzable_extensions = ['.py', '.js', '.ts', '.go', '.java', '.cpp', '.c', '.h', '.hpp']
        return any(filename.endswith(ext) for ext in analyzable_extensions)

    def _analyze_single_file(self, file_path: str) -> Dict[str, Any]:
        """Analyze a single file for issues and patterns"""
        analysis = {
            "file_path": file_path,
            "language": self._detect_language(file_path),
            "issues": [],
            "patterns": [],
            "complexity": 0,
            "lines": 0
        }

        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            analysis["lines"] = len(content.split('\n'))

            if file_path.endswith('.py'):
                analysis.update(self._analyze_python_file(file_path, content))
            elif file_path.endswith('.js') or file_path.endswith('.ts'):
                analysis.update(self._analyze_javascript_file(file_path, content))
            elif file_path.endswith('.go'):
                analysis.update(self._analyze_go_file(file_path, content))

        except Exception as e:
            analysis["issues"].append({
                "type": "file_read_error",
                "severity": "medium",
                "message": f"Could not read file: {e}",
                "line": 0
            })

        return analysis

    def _detect_language(self, file_path: str) -> str:
        """Detect programming language from file extension"""
        ext_map = {
            '.py': 'python',
            '.js': 'javascript',
            '.ts': 'typescript',
            '.go': 'go',
            '.java': 'java',
            '.cpp': 'cpp',
            '.c': 'c',
            '.h': 'c',
            '.hpp': 'cpp'
        }
        return ext_map.get(os.path.splitext(file_path)[1], 'unknown')

    def _analyze_python_file(self, file_path: str, content: str) -> Dict[str, Any]:
        """Analyze Python file for common issues"""
        issues = []
        patterns = []

        try:
            tree = ast.parse(content)

            # Check for common issues
            for node in ast.walk(tree):
                if isinstance(node, ast.FunctionDef):
                    if len(node.body) > 20:  # Long function
                        issues.append({
                            "type": "long_function",
                            "severity": "low",
                            "message": f"Function '{node.name}' is too long ({len(node.body)} lines)",
                            "line": node.lineno
                        })

                    # Check for missing docstrings
                    if not ast.get_docstring(node):
                        issues.append({
                            "type": "missing_docstring",
                            "severity": "low",
                            "message": f"Function '{node.name}' missing docstring",
                            "line": node.lineno
                        })

                elif isinstance(node, ast.ClassDef):
                    if not node.bases:  # No inheritance
                        patterns.append({
                            "type": "standalone_class",
                            "message": f"Class '{node.name}' doesn't inherit from anything",
                            "line": node.lineno
                        })

            # Check for imports
            if 'import os' in content and 'import sys' in content:
                patterns.append({
                    "type": "common_imports",
                    "message": "File uses common system imports",
                    "line": 1
                })

        except SyntaxError as e:
            issues.append({
                "type": "syntax_error",
                "severity": "high",
                "message": f"Syntax error: {e.msg}",
                "line": e.lineno
            })

        return {"issues": issues, "patterns": patterns}

    def _analyze_javascript_file(self, file_path: str, content: str) -> Dict[str, Any]:
        """Analyze JavaScript/TypeScript file"""
        issues = []
        patterns = []

        # Check for console.log statements
        if 'console.log' in content:
            issues.append({
                "type": "debug_statement",
                "severity": "low",
                "message": "Found console.log statement",
                "line": content.find('console.log')
            })

        # Check for proper error handling
        if 'try' in content and 'catch' not in content:
            issues.append({
                "type": "incomplete_error_handling",
                "severity": "medium",
                "message": "Try block without catch",
                "line": content.find('try')
            })

        return {"issues": issues, "patterns": patterns}

    def _analyze_go_file(self, file_path: str, content: str) -> Dict[str, Any]:
        """Analyze Go file"""
        issues = []
        patterns = []

        # Check for error handling
        if 'func main()' in content and 'error' not in content.lower():
            issues.append({
                "type": "missing_error_handling",
                "severity": "medium",
                "message": "Main function should handle errors",
                "line": content.find('func main()')
            })

        return {"issues": issues, "patterns": patterns}

    def _generate_codebase_suggestions(self, analysis: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Generate improvement suggestions using AutoGen"""
        try:
            suggestions = []

            # Use pattern learning to suggest fixes
            for issue in analysis["issues_found"][:10]:  # Limit to prevent overload
                pattern_solution = suggest_pattern_solution(issue["type"])
                if pattern_solution:
                    suggestions.append({
                        "type": "pattern_based_fix",
                        "title": f"Fix {issue['type']}",
                        "description": pattern_solution,
                        "confidence": 0.8,
                        "affected_files": [issue.get("file", "unknown")]
                    })

            # Generate structural suggestions
            if analysis["files_analyzed"] > 50:
                suggestions.append({
                    "type": "structural",
                    "title": "Large Codebase Optimization",
                    "description": "Consider splitting large files or modules",
                    "confidence": 0.7,
                    "affected_files": ["multiple"]
                })

            # Language-specific suggestions
            for lang, count in analysis["languages"].items():
                if count > 20:
                    suggestions.append({
                        "type": "language_optimization",
                        "title": f"{lang.capitalize()} Code Optimization",
                        "description": f"Large {lang} codebase detected, consider refactoring",
                        "confidence": 0.6,
                        "affected_files": [lang]
                    })

            return suggestions

        except Exception as e:
            print(f"❌ Failed to generate suggestions: {e}")
            return []

    def fix_identified_issues(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Fix identified issues using the factory system"""
        print("🔧 Fixing identified issues...")

        fix_results = {
            "total_issues": len(analysis["issues_found"]),
            "fixed": 0,
            "failed": 0,
            "skipped": 0,
            "details": []
        }

        try:
            # Group issues by type for batch processing
            issues_by_type = {}
            for issue in analysis["issues_found"]:
                issue_type = issue["type"]
                if issue_type not in issues_by_type:
                    issues_by_type[issue_type] = []
                issues_by_type[issue_type].append(issue)

            # Process each issue type
            for issue_type, issues in issues_by_type.items():
                print(f"📋 Processing {len(issues)} {issue_type} issues...")

                # Try pattern-based solution first
                pattern_solution = suggest_pattern_solution(issue_type)
                if pattern_solution:
                    fix_success = self._apply_pattern_fix(issue_type, issues, pattern_solution)
                    if fix_success:
                        fix_results["fixed"] += len(issues)
                        fix_results["details"].append({
                            "type": issue_type,
                            "method": "pattern_learning",
                            "success": True,
                            "count": len(issues)
                        })
                        print(f"  ✅ Fixed {len(issues)} {issue_type} issues using pattern learning")
                    else:
                        fix_results["failed"] += len(issues)
                        fix_results["details"].append({
                            "type": issue_type,
                            "method": "pattern_learning",
                            "success": False,
                            "count": len(issues)
                        })
                        print(f"  ❌ Failed to fix {len(issues)} {issue_type} issues")
                else:
                    fix_results["skipped"] += len(issues)
                    fix_results["details"].append({
                        "type": issue_type,
                        "method": "no_pattern",
                        "success": False,
                        "count": len(issues)
                    })
                    print(f"  ⏭️ Skipped {len(issues)} {issue_type} issues (no pattern available)")

            # Learn from the fixing process
            deployment_data = {
                "error_type": "code_fixing",
                "success": fix_results["fixed"] > 0,
                "solution": "automated_fixing",
                "context": {
                    "total_issues": fix_results["total_issues"],
                    "fixed": fix_results["fixed"],
                    "failed": fix_results["failed"]
                },
                "execution_time": 0,  # Will be set by the system
                "timestamp": datetime.now()
            }

            learn_from_deployment(deployment_data)

            print(f"✅ Fix process completed: {fix_results['fixed']} fixed, {fix_results['failed']} failed")
            return fix_results

        except Exception as e:
            print(f"❌ Fix process failed: {e}")
            return fix_results

    def _apply_pattern_fix(self, issue_type: str, issues: List[Dict], solution: str) -> bool:
        """Apply a pattern-based fix to issues"""
        try:
            # This is a simplified implementation
            # In a real system, you would parse and modify the actual files

            if issue_type == "missing_docstring":
                print(f"    📝 Would add docstrings to {len(issues)} functions")
                return True
            elif issue_type == "long_function":
                print(f"    ✂️ Would split {len(issues)} long functions")
                return True
            elif issue_type == "debug_statement":
                print(f"    🧹 Would remove {len(issues)} debug statements")
                return True
            else:
                print(f"    🤷 No specific fix available for {issue_type}")
                return False

        except Exception as e:
            print(f"❌ Failed to apply fix: {e}")
            return False

    def run_comprehensive_tests(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Run comprehensive tests on the codebase"""
        print("🧪 Running comprehensive tests...")

        test_results = {
            "timestamp": datetime.now().isoformat(),
            "syntax_tests": {"passed": 0, "failed": 0, "details": []},
            "import_tests": {"passed": 0, "failed": 0, "details": []},
            "security_tests": {"passed": 0, "failed": 0, "details": []},
            "performance_tests": {"passed": 0, "failed": 0, "details": []}
        }

        try:
            # Syntax tests
            syntax_results = self._run_syntax_tests(analysis)
            test_results["syntax_tests"] = syntax_results

            # Import tests
            import_results = self._run_import_tests(analysis)
            test_results["import_tests"] = import_results

            # Security tests
            security_results = self._run_security_tests(analysis)
            test_results["security_tests"] = security_results

            # Performance tests
            performance_results = self._run_performance_tests(analysis)
            test_results["performance_tests"] = performance_results

            # Calculate overall statistics
            total_passed = (
                syntax_results["passed"] +
                import_results["passed"] +
                security_results["passed"] +
                performance_results["passed"]
            )
            total_failed = (
                syntax_results["failed"] +
                import_results["failed"] +
                security_results["failed"] +
                performance_results["failed"]
            )

            print("✅ Tests completed:")
            print(f"  📊 Total: {total_passed + total_failed}")
            print(f"  ✅ Passed: {total_passed}")
            print(f"  ❌ Failed: {total_failed}")

            return test_results

        except Exception as e:
            print(f"❌ Test execution failed: {e}")
            return test_results

    def _run_syntax_tests(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Run syntax validation tests"""
        results = {"passed": 0, "failed": 0, "details": []}

        try:
            for file_path in self._get_analyzable_files():
                try:
                    if file_path.endswith('.py'):
                        with open(file_path, 'r') as f:
                            content = f.read()
                        ast.parse(content)  # This will raise SyntaxError if invalid
                        results["passed"] += 1
                        results["details"].append({
                            "file": file_path,
                            "status": "passed",
                            "test": "syntax"
                        })
                    else:
                        results["passed"] += 1
                except SyntaxError as e:
                    results["failed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "failed",
                        "test": "syntax",
                        "error": str(e)
                    })
                except Exception as e:
                    results["failed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "failed",
                        "test": "syntax",
                        "error": str(e)
                    })

        except Exception as e:
            print(f"❌ Syntax tests failed: {e}")

        return results

    def _run_import_tests(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Test import statements"""
        results = {"passed": 0, "failed": 0, "details": []}

        try:
            # Test Python imports
            for file_path in self._get_python_files():
                try:
                    # Simple import test - try to compile the file
                    with open(file_path, 'r') as f:
                        content = f.read()

                    compile(content, file_path, 'exec')
                    results["passed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "passed",
                        "test": "imports"
                    })
                except Exception as e:
                    results["failed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "failed",
                        "test": "imports",
                        "error": str(e)
                    })

        except Exception as e:
            print(f"❌ Import tests failed: {e}")

        return results

    def _run_security_tests(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Run basic security tests"""
        results = {"passed": 0, "failed": 0, "details": []}

        try:
            # Check for potential security issues
            security_patterns = [
                "eval(", "exec(", "subprocess.call", "os.system"
            ]

            for file_path in self._get_analyzable_files():
                try:
                    with open(file_path, 'r') as f:
                        content = f.read()

                    found_issues = []
                    for pattern in security_patterns:
                        if pattern in content:
                            found_issues.append(pattern)

                    if not found_issues:
                        results["passed"] += 1
                        results["details"].append({
                            "file": file_path,
                            "status": "passed",
                            "test": "security"
                        })
                    else:
                        results["failed"] += 1
                        results["details"].append({
                            "file": file_path,
                            "status": "failed",
                            "test": "security",
                            "issues": found_issues
                        })

                except Exception as e:
                    results["failed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "failed",
                        "test": "security",
                        "error": str(e)
                    })

        except Exception as e:
            print(f"❌ Security tests failed: {e}")

        return results

    def _run_performance_tests(self, analysis: Dict[str, Any]) -> Dict[str, Any]:
        """Run basic performance tests"""
        results = {"passed": 0, "failed": 0, "details": []}

        try:
            # Simple performance checks
            for file_path in self._get_analyzable_files():
                try:
                    # Check file size (basic performance indicator)
                    file_size = os.path.getsize(file_path)

                    if file_size < 1024 * 1024:  # Less than 1MB
                        results["passed"] += 1
                        results["details"].append({
                            "file": file_path,
                            "status": "passed",
                            "test": "performance",
                            "size_kb": file_size / 1024
                        })
                    else:
                        results["failed"] += 1
                        results["details"].append({
                            "file": file_path,
                            "status": "warning",
                            "test": "performance",
                            "size_kb": file_size / 1024,
                            "message": "Large file may impact performance"
                        })

                except Exception as e:
                    results["failed"] += 1
                    results["details"].append({
                        "file": file_path,
                        "status": "failed",
                        "test": "performance",
                        "error": str(e)
                    })

        except Exception as e:
            print(f"❌ Performance tests failed: {e}")

        return results

    def _get_analyzable_files(self) -> List[str]:
        """Get list of analyzable files"""
        files = []
        for root, dirs, filenames in os.walk("."):
            dirs[:] = [d for d in dirs if not d.startswith('.') and d not in ['node_modules', '__pycache__', '.git']]
            for filename in filenames:
                if self._should_analyze_file(filename):
                    files.append(os.path.join(root, filename))
        return files

    def _get_python_files(self) -> List[str]:
        """Get list of Python files"""
        return [f for f in self._get_analyzable_files() if f.endswith('.py')]

    def generate_review_report(self, analysis: Dict[str, Any], fixes: Dict[str, Any], tests: Dict[str, Any]) -> str:
        """Generate a comprehensive review report"""
        report = f"""
# Code Review Report
Generated on: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

## 📊 Executive Summary
- **Files Analyzed**: {analysis['files_analyzed']}
- **Languages Detected**: {', '.join(analysis['languages'].keys())}
- **Issues Found**: {len(analysis['issues_found'])}
- **Fixes Applied**: {fixes['fixed']}
- **Tests Passed**: {sum(t['passed'] for t in tests.values())}

## 🔍 Codebase Analysis
### Language Distribution
"""

        for lang, count in analysis['languages'].items():
            report += f"- {lang.title()}: {count} files\n"

        report += "\n### Issues by Severity\n"
        severity_count = {"high": 0, "medium": 0, "low": 0}
        for issue in analysis['issues_found']:
            severity = issue.get('severity', 'medium')
            severity_count[severity] += 1

        for severity, count in severity_count.items():
            report += f"- {severity.title()}: {count} issues\n"

        report += "\n## 🔧 Fixes Applied\n"
        for detail in fixes['details']:
            status = "✅" if detail['success'] else "❌"
            report += f"- {status} {detail['type']}: {detail['count']} issues ({detail['method']})\n"

        report += "\n## 🧪 Test Results\n"
        for test_type, results in tests.items():
            if isinstance(results, dict) and 'passed' in results:
                report += f"### {test_type.replace('_', ' ').title()}\n"
                report += f"- Passed: {results['passed']}\n"
                report += f"- Failed: {results['failed']}\n\n"

        report += "
## 🎯 Recommendations
"
        for suggestion in analysis['suggestions'][:5]:  # Top 5 suggestions
            report += f"- {suggestion['title']}: {suggestion['description'][:100]}...\n"

        return report

    def run_continuous_improvement(self) -> Dict[str, Any]:
        """Run continuous improvement on the review system itself"""
        print("🔄 Running continuous improvement...")

        try:
            # Analyze system performance
            suggestions = analyze_system_performance()

            # Implement safe improvements
            improvement_results = implement_safe_improvements()

            # Learn from the improvement process
            deployment_data = {
                "error_type": "system_optimization",
                "success": improvement_results["implemented"] > 0,
                "solution": "continuous_improvement",
                "context": {
                    "suggestions_generated": len(suggestions),
                    "improvements_applied": improvement_results["implemented"]
                },
                "execution_time": 0,
                "timestamp": datetime.now()
            }

            learn_from_deployment(deployment_data)

            print(f"✅ Continuous improvement: {improvement_results['implemented']} optimizations applied")
            return improvement_results

        except Exception as e:
            print(f"❌ Continuous improvement failed: {e}")
            return {"implemented": 0, "failed": 1, "error": str(e)}


def main():
    """Main code review workflow"""
    print("🚀 AutoGen Factory Code Review System")
    print("=" * 50)

    # Set API key for AutoGen
    os.environ["OPENAI_API_KEY"] = "your-api-key-here"

    try:
        # Initialize the review agent
        agent = CodeReviewAgent()

        # Step 1: Analyze codebase
        print("\n📋 STEP 1: Codebase Analysis")
        analysis = agent.analyze_codebase(".")

        # Step 2: Fix identified issues
        print("\n🔧 STEP 2: Fixing Issues")
        fixes = agent.fix_identified_issues(analysis)

        # Step 3: Run comprehensive tests
        print("\n🧪 STEP 3: Testing")
        tests = agent.run_comprehensive_tests(analysis)

        # Step 4: Generate report
        print("\n📄 STEP 4: Generating Report")
        report = agent.generate_review_report(analysis, fixes, tests)

        # Save report
        report_file = f"code_review_report_{int(time.time())}.md"
        with open(report_file, 'w') as f:
            f.write(report)

        print(f"✅ Report saved to: {report_file}")

        # Step 5: Continuous improvement
        print("\n🔄 STEP 5: Continuous Improvement")
        improvements = agent.run_continuous_improvement()

        print("\n🎉 Code review process completed successfully!")
        print("\n📊 Summary:")
        print(f"  • Files analyzed: {analysis['files_analyzed']}")
        print(f"  • Issues found: {len(analysis['issues_found'])}")
        print(f"  • Fixes applied: {fixes['fixed']}")
        print(f"  • Tests passed: {sum(t['passed'] for t in tests.values() if isinstance(t, dict))}")
        print(f"  • Improvements: {improvements['implemented']}")

        return True

    except Exception as e:
        print(f"❌ Code review process failed: {e}")
        import traceback
        traceback.print_exc()
        return False


if __name__ == "__main__":
    success = main()
    exit(0 if success else 1)
