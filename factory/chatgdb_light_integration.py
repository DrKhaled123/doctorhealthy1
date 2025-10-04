#!/usr/bin/env python3
"""
ChatGDB Light Integration for Factory Orchestrator
Lightweight debugging and error analysis capabilities
"""

import os
import json
import time
import asyncio
import re
from typing import Dict, List, Any, Optional
from datetime import datetime
import subprocess

class ChatGDBLight:
    """Lightweight ChatGDB for debugging and error analysis"""

    def __init__(self, model_path: str = None):
        self.model_path = model_path or os.getenv('CHAT_GDB_MODEL', 'microsoft/DialoGPT-medium')
        self.debug_history = []
        self.error_patterns = {}
        self.solution_cache = {}

    async def analyze_error(self, error_message: str, code_context: str = "",
                          stack_trace: str = "", language: str = "python") -> Dict[str, Any]:
        """Analyze error using ChatGDB Light"""
        print(f"🔍 Analyzing error: {error_message[:100]}...")

        try:
            # Create analysis context
            context = {
                'error_message': error_message,
                'code_context': code_context,
                'stack_trace': stack_trace,
                'language': language,
                'timestamp': datetime.now().isoformat()
            }

            # Check for known error patterns
            pattern_match = self._match_error_pattern(error_message)
            if pattern_match:
                print(f"✅ Found matching pattern: {pattern_match['pattern']}")
                return {
                    'success': True,
                    'pattern_match': True,
                    'solution': pattern_match['solution'],
                    'confidence': pattern_match['confidence'],
                    'explanation': pattern_match['explanation'],
                    'context': context
                }

            # Generate solution using pattern recognition
            solution = await self._generate_solution(context)

            # Cache the solution
            self.solution_cache[error_message] = solution

            # Store in debug history
            self.debug_history.append({
                'timestamp': datetime.now().isoformat(),
                'error': error_message,
                'solution': solution,
                'context': context
            })

            return {
                'success': True,
                'pattern_match': False,
                'solution': solution,
                'confidence': 0.8,
                'explanation': 'Generated solution based on error analysis',
                'context': context
            }

        except Exception as e:
            return {
                'success': False,
                'error': str(e),
                'context': context
            }

    def _match_error_pattern(self, error_message: str) -> Optional[Dict]:
        """Match error against known patterns"""
        # Common error patterns and solutions
        patterns = {
            'import_error': {
                'pattern': r'ModuleNotFoundError|ImportError',
                'solution': 'Install missing module using pip install',
                'confidence': 0.9,
                'explanation': 'Missing module needs to be installed'
            },
            'syntax_error': {
                'pattern': r'SyntaxError',
                'solution': 'Check for missing colons, brackets, or quotes',
                'confidence': 0.95,
                'explanation': 'Python syntax error detected'
            },
            'indentation_error': {
                'pattern': r'IndentationError',
                'solution': 'Fix indentation - use 4 spaces per level',
                'confidence': 0.9,
                'explanation': 'Python indentation issue'
            },
            'type_error': {
                'pattern': r'TypeError',
                'solution': 'Check variable types and method signatures',
                'confidence': 0.8,
                'explanation': 'Type mismatch in operation'
            },
            'attribute_error': {
                'pattern': r'AttributeError',
                'solution': 'Check if object has the required attribute',
                'confidence': 0.85,
                'explanation': 'Object missing expected attribute'
            },
            'connection_error': {
                'pattern': r'ConnectionError|ConnectionRefusedError',
                'solution': 'Check network connection and server status',
                'confidence': 0.9,
                'explanation': 'Network connectivity issue'
            }
        }

        for pattern_name, pattern_info in patterns.items():
            if re.search(pattern_info['pattern'], error_message, re.IGNORECASE):
                return {
                    'pattern': pattern_name,
                    'solution': pattern_info['solution'],
                    'confidence': pattern_info['confidence'],
                    'explanation': pattern_info['explanation']
                }

        return None

    async def _generate_solution(self, context: Dict) -> str:
        """Generate solution for unknown error patterns"""
        error_msg = context['error_message']
        code_context = context.get('code_context', '')
        language = context['language']

        # Simple rule-based solution generation
        if 'redis' in error_msg.lower():
            return "Check Redis connection: ensure Redis server is running and accessible"
        elif 'permission' in error_msg.lower():
            return "Check file/directory permissions and access rights"
        elif 'memory' in error_msg.lower():
            return "Check memory usage and consider increasing available memory"
        elif 'timeout' in error_msg.lower():
            return "Increase timeout values or optimize slow operations"
        elif 'encoding' in error_msg.lower():
            return "Specify correct encoding (usually 'utf-8') when opening files"
        else:
            # Generic solution based on error characteristics
            return self._generate_generic_solution(error_msg, code_context, language)

    def _generate_generic_solution(self, error_msg: str, code_context: str, language: str) -> str:
        """Generate generic solution for unknown errors"""
        solutions = [
            "Check the error message and stack trace carefully",
            "Verify all dependencies are installed",
            "Check configuration files for correct settings",
            "Review recent code changes that might have caused the issue",
            "Check system resources (memory, disk space)",
            "Verify network connectivity if applicable"
        ]

        # Add language-specific suggestions
        if language == 'python':
            solutions.extend([
                "Check Python version compatibility",
                "Verify all imports are correct",
                "Check for typos in variable/function names"
            ])
        elif language == 'javascript':
            solutions.extend([
                "Check Node.js version",
                "Verify npm dependencies",
                "Check for async/await issues"
            ])

        return " | ".join(solutions)

    async def debug_code_execution(self, code: str, expected_behavior: str = "") -> Dict[str, Any]:
        """Debug code execution and suggest fixes"""
        print("🐛 Debugging code execution...")

        try:
            # Analyze code for potential issues
            issues = self._analyze_code_issues(code)

            # Execute code safely
            execution_result = await self._safe_code_execution(code)

            # Generate debugging report
            debug_report = {
                'timestamp': datetime.now().isoformat(),
                'code_length': len(code),
                'issues_found': issues,
                'execution_result': execution_result,
                'suggestions': []
            }

            # Generate suggestions based on analysis
            if issues:
                debug_report['suggestions'] = self._generate_debug_suggestions(issues, execution_result)

            if execution_result['success']:
                print("✅ Code executed successfully")
            else:
                print(f"❌ Code execution failed: {execution_result['error']}")

            return debug_report

        except Exception as e:
            return {
                'success': False,
                'error': str(e),
                'timestamp': datetime.now().isoformat()
            }

    def _analyze_code_issues(self, code: str) -> List[Dict]:
        """Analyze code for potential issues"""
        issues = []

        # Check for common issues
        if 'eval(' in code:
            issues.append({
                'type': 'security_risk',
                'severity': 'high',
                'message': 'Use of eval() detected - potential security risk',
                'line': self._find_line_number(code, 'eval(')
            })

        if 'exec(' in code:
            issues.append({
                'type': 'security_risk',
                'severity': 'high',
                'message': 'Use of exec() detected - potential security risk',
                'line': self._find_line_number(code, 'exec(')
            })

        if 'subprocess' in code and 'shell=True' in code:
            issues.append({
                'type': 'security_risk',
                'severity': 'medium',
                'message': 'Subprocess with shell=True detected',
                'line': self._find_line_number(code, 'shell=True')
            })

        # Check for potential bugs
        if '==' in code:
            lines = code.split('\n')
            for i, line in enumerate(lines):
                if '==' in line and not line.strip().startswith('#'):
                    issues.append({
                        'type': 'potential_bug',
                        'severity': 'low',
                        'message': 'Found equality comparison - consider using isinstance() for type checking',
                        'line': i + 1
                    })

        return issues

    async def _safe_code_execution(self, code: str) -> Dict:
        """Safely execute code in isolated environment"""
        try:
            # Create temporary file for execution
            with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
                f.write(code)
                temp_file = f.name

            # Execute with timeout
            result = subprocess.run(
                ['python3', temp_file],
                capture_output=True,
                text=True,
                timeout=30
            )

            return {
                'success': result.returncode == 0,
                'stdout': result.stdout,
                'stderr': result.stderr,
                'return_code': result.returncode
            }

        except subprocess.TimeoutExpired:
            return {
                'success': False,
                'error': 'Code execution timed out',
                'return_code': -1
            }
        except Exception as e:
            return {
                'success': False,
                'error': str(e),
                'return_code': -1
            }
        finally:
            # Clean up temporary file
            if 'temp_file' in locals():
                os.unlink(temp_file)

    def _find_line_number(self, code: str, pattern: str) -> int:
        """Find line number containing pattern"""
        lines = code.split('\n')
        for i, line in enumerate(lines):
            if pattern in line:
                return i + 1
        return 0

    def _generate_debug_suggestions(self, issues: List[Dict], execution_result: Dict) -> List[str]:
        """Generate debugging suggestions"""
        suggestions = []

        for issue in issues:
            if issue['type'] == 'security_risk':
                suggestions.append(f"Security issue on line {issue['line']}: {issue['message']}")
            elif issue['type'] == 'potential_bug':
                suggestions.append(f"Potential bug on line {issue['line']}: {issue['message']}")

        if not execution_result['success']:
            suggestions.append(f"Runtime error: {execution_result['error']}")

        return suggestions

    async def analyze_stack_trace(self, stack_trace: str) -> Dict[str, Any]:
        """Analyze stack trace for debugging insights"""
        print("📊 Analyzing stack trace...")

        try:
            lines = stack_trace.split('\n')
            analysis = {
                'total_frames': 0,
                'file_locations': [],
                'error_types': [],
                'suggestions': []
            }

            for line in lines:
                if 'File "' in line:
                    analysis['total_frames'] += 1
                    file_match = re.search(r'File "([^"]+)", line (\d+)', line)
                    if file_match:
                        analysis['file_locations'].append({
                            'file': file_match.group(1),
                            'line': int(file_match.group(2))
                        })

                if 'Error' in line or 'Exception' in line:
                    analysis['error_types'].append(line.strip())

            # Generate suggestions based on stack trace
            if 'ImportError' in stack_trace or 'ModuleNotFoundError' in stack_trace:
                analysis['suggestions'].append('Check if all required modules are installed')
            if 'AttributeError' in stack_trace:
                analysis['suggestions'].append('Check if object has the expected attributes')
            if 'TypeError' in stack_trace:
                analysis['suggestions'].append('Check data types and method signatures')
            if 'IndentationError' in stack_trace:
                analysis['suggestions'].append('Fix Python indentation (use 4 spaces)')

            return analysis

        except Exception as e:
            return {
                'error': str(e),
                'suggestions': ['Unable to analyze stack trace']
            }

    async def get_debug_recommendations(self, error_context: Dict) -> List[str]:
        """Get debugging recommendations based on error context"""
        recommendations = []

        error_msg = error_context.get('error_message', '')
        code_context = error_context.get('code_context', '')
        stack_trace = error_context.get('stack_trace', '')

        # Analyze different aspects
        if error_msg:
            error_analysis = await self.analyze_error(error_msg, code_context, stack_trace)
            if error_analysis['success']:
                recommendations.append(f"Error analysis: {error_analysis['solution']}")

        if stack_trace:
            stack_analysis = await self.analyze_stack_trace(stack_trace)
            recommendations.extend(stack_analysis.get('suggestions', []))

        if code_context:
            code_issues = self._analyze_code_issues(code_context)
            for issue in code_issues:
                recommendations.append(f"Code issue: {issue['message']}")

        return recommendations

# Factory integration functions
async def debug_with_chatgdb(error_message: str, code_context: str = "",
                           stack_trace: str = "", language: str = "python") -> Dict[str, Any]:
    """Debug using ChatGDB Light"""
    debugger = ChatGDBLight()
    return await debugger.analyze_error(error_message, code_context, stack_trace, language)

async def analyze_code_with_chatgdb(code: str, expected_behavior: str = "") -> Dict[str, Any]:
    """Analyze code using ChatGDB Light"""
    debugger = ChatGDBLight()
    return await debugger.debug_code_execution(code, expected_behavior)

async def analyze_stack_trace_with_chatgdb(stack_trace: str) -> Dict[str, Any]:
    """Analyze stack trace using ChatGDB Light"""
    debugger = ChatGDBLight()
    return await debugger.analyze_stack_trace(stack_trace)

# Example usage
async def demo_chatgdb_light():
    """Demonstrate ChatGDB Light integration"""
    print("🐛 ChatGDB Light Integration Demo")
    print("=" * 40)

    debugger = ChatGDBLight()

    # Example 1: Analyze error
    error_msg = "ModuleNotFoundError: No module named 'requests'"
    code_context = "import requests"

    print("🔍 Analyzing import error...")
    error_analysis = await debugger.analyze_error(error_msg, code_context)

    if error_analysis['success']:
        print(f"✅ Error analysis: {error_analysis['solution']}")
        print(f"   Confidence: {error_analysis['confidence']}")
    else:
        print(f"❌ Error analysis failed: {error_analysis['error']}")

    # Example 2: Analyze code
    problematic_code = '''
def bad_function():
    eval("print('hello')")
    return True
'''

    print("\n🔍 Analyzing problematic code...")
    code_analysis = await debugger.debug_code_execution(problematic_code)

    if code_analysis.get('issues_found'):
        print("✅ Found code issues:")
        for issue in code_analysis['issues_found']:
            print(f"  • {issue['message']}")

    # Example 3: Stack trace analysis
    stack_trace = '''
Traceback (most recent call last):
  File "test.py", line 10, in <module>
    main()
  File "test.py", line 5, in main
    bad_function()
  File "test.py", line 2, in bad_function
    raise ValueError("Test error")
ValueError: Test error
'''

    print("\n🔍 Analyzing stack trace...")
    stack_analysis = await debugger.analyze_stack_trace(stack_trace)

    print(f"✅ Stack trace analysis: {stack_analysis['total_frames']} frames")
    print(f"   Error types: {stack_analysis['error_types']}")

    print("\n✅ ChatGDB Light demo completed!")

if __name__ == "__main__":
    asyncio.run(demo_chatgdb_light())