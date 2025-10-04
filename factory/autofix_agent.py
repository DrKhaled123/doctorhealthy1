import ast
import autopep8
import black
from typing import List, Tuple, Dict, Any
import re
import json
import subprocess
import sys
import os

class AutofixAgent:
    """Automatic code fixing and formatting agent"""

    def __init__(self, config: Dict = None):
        self.config = config or {}
        self.fix_strategies = [
            self.fix_syntax_errors,
            self.fix_import_errors,
            self.fix_type_errors,
            self.apply_formatting,
            self.fix_security_issues,
            self.fix_common_patterns
        ]

        # Setup formatters
        self.black_mode = black.FileMode()
        self.autopep8_options = {
            'aggressive': 1,
            'max_line_length': 88
        }

    async def fix_code(self, code: str, errors: List[Dict] = None) -> Dict:
        """Apply automatic fixes based on errors and code analysis"""
        result = {
            "original_code": code,
            "fixed_code": code,
            "fixes_applied": [],
            "errors_remaining": [],
            "success": True
        }

        try:
            # Run all fix strategies
            for strategy in self.fix_strategies:
                try:
                    strategy_result = await strategy(result["fixed_code"], errors or [])
                    result["fixed_code"] = strategy_result["code"]
                    result["fixes_applied"].extend(strategy_result["fixes"])
                except Exception as e:
                    print(f"Fix strategy failed: {e}")
                    result["errors_remaining"].append(str(e))

            # Final validation
            validation_result = self.validate_fixes(result["fixed_code"])
            result["errors_remaining"].extend(validation_result["errors"])
            result["success"] = len(result["errors_remaining"]) == 0

        except Exception as e:
            result["success"] = False
            result["errors_remaining"].append(f"General fix error: {str(e)}")

        return result

    async def fix_syntax_errors(self, code: str, errors: List[Dict]) -> Dict:
        """Fix common syntax errors"""
        fixes = []

        try:
            # Try to parse the code to find syntax errors
            ast.parse(code)
            return {"code": code, "fixes": fixes}
        except SyntaxError as e:
            fixes.append(f"Fixed syntax error at line {e.lineno}: {e.msg}")

            # Get the problematic lines
            lines = code.split('\n')

            # Fix common syntax issues
            if "invalid syntax" in str(e):
                # Fix missing colons
                for i, line in enumerate(lines):
                    stripped = line.strip()
                    if (stripped.startswith(('if', 'elif', 'else', 'for', 'while',
                                          'def', 'class', 'try', 'except', 'finally')) and
                        not stripped.endswith(':')):
                        lines[i] = line.rstrip() + ':'
                        fixes.append(f"Added missing colon at line {i+1}")

            # Fix indentation issues
            if "unexpected indent" in str(e) or "indentation" in str(e):
                lines = self.fix_indentation(lines)
                fixes.append("Fixed indentation issues")

            # Fix missing parentheses/brackets
            if "missing" in str(e).lower():
                lines = self.fix_missing_brackets(lines, e)
                fixes.append("Fixed missing brackets/parentheses")

        except Exception as e:
            fixes.append(f"Could not fix syntax error: {str(e)}")

        return {"code": '\n'.join(lines), "fixes": fixes}

    async def fix_import_errors(self, code: str, errors: List[Dict]) -> Dict:
        """Fix import-related errors"""
        fixes = []
        lines = code.split('\n')

        # Check for common import issues
        for i, line in enumerate(lines):
            stripped = line.strip()

            # Fix incorrect import statements
            if stripped.startswith('import '):
                parts = stripped.split()
                if len(parts) >= 2:
                    module = parts[1]
                    # Check if module needs to be installed
                    if not self.module_exists(module):
                        fixes.append(f"Module '{module}' may need to be installed")
                    # Fix relative imports
                    if module.startswith('.'):
                        if not module.startswith('..'):
                            lines[i] = line.replace('import .', 'from . import ')
                            fixes.append(f"Fixed relative import at line {i+1}")

            # Fix from ... import ... statements
            elif stripped.startswith('from '):
                if 'import' in stripped:
                    parts = stripped.split()
                    if len(parts) >= 4 and parts[2] == 'import':
                        module = parts[1]
                        if not self.module_exists(module):
                            fixes.append(f"Module '{module}' may need to be installed")

        return {"code": '\n'.join(lines), "fixes": fixes}

    async def fix_type_errors(self, code: str, errors: List[Dict]) -> Dict:
        """Fix type-related errors"""
        fixes = []
        lines = code.split('\n')

        # Add type hints where missing
        for i, line in enumerate(lines):
            stripped = line.strip()

            # Fix function definitions without type hints
            if (stripped.startswith('def ') and
                '(' in stripped and
                not '->' in stripped and
                not 'self' in stripped.split('(')[1].split(',')[0].strip()):

                # Try to infer return type
                func_body = self.get_function_body(lines, i)
                if self.has_return_statement(func_body):
                    lines[i] = line.replace('):', ') -> Any:').replace('def ', 'def ')
                    fixes.append(f"Added return type hint to function at line {i+1}")

        return {"code": '\n'.join(lines), "fixes": fixes}

    async def apply_formatting(self, code: str, errors: List[Dict]) -> Dict:
        """Apply code formatting"""
        fixes = []

        try:
            # Apply autopep8
            formatted_code = autopep8.fix_code(code, options=self.autopep8_options)
            fixes.append("Applied PEP8 formatting")

            # Apply black formatting
            try:
                formatted_code = black.format_str(formatted_code, mode=self.black_mode)
                fixes.append("Applied Black formatting")
            except Exception as e:
                # Black might fail on some code, that's okay'
                pass

        except Exception as e:
            fixes.append(f"Formatting failed: {str(e)}")
            formatted_code = code

        return {"code": formatted_code, "fixes": fixes}

    async def fix_security_issues(self, code: str, errors: List[Dict]) -> Dict:
        """Fix common security issues"""
        fixes = []
        lines = code.split('\n')

        for i, line in enumerate(lines):
            stripped = line.strip()

            # Fix hardcoded passwords
            if re.search(r'password\s*=\s*["\'][^"\']+["\']', stripped, re.IGNORECASE):
                lines[i] = re.sub(
                    r'(["\'])([^"\']+)(["\'])',
                    r'\1***REDACTED***\3',
                    stripped
                )
                fixes.append(f"Redacted hardcoded password at line {i+1}")

            # Fix SQL injection vulnerabilities
            if re.search(r'["\']\+.*\+["\']', stripped):
                if 'SELECT' in stripped.upper() or 'INSERT' in stripped.upper():
                    lines[i] = stripped.replace('+"', " %s ").replace('+ "', " %s ")
                    fixes.append(f"Fixed potential SQL injection at line {i+1}")

            # Fix use of eval
            if 'eval(' in stripped:
                lines[i] = "# SECURITY: Removed eval() - use ast.literal_eval() instead"
                fixes.append(f"Removed unsafe eval() at line {i+1}")

        return {"code": '\n'.join(lines), "fixes": fixes}

    async def fix_common_patterns(self, code: str, errors: List[Dict]) -> Dict:
        """Fix common code patterns"""
        fixes = []
        lines = code.split('\n')

        for i, line in enumerate(lines):
            stripped = line.strip()

            # Fix bare except clauses
            if stripped == 'except:' or stripped.startswith('except Exception:'):
                lines[i] = line.replace('except:', 'except Exception as e:').replace(
                    'except Exception:', 'except Exception as e:'
                )
                fixes.append(f"Fixed bare except clause at line {i+1}")

            # Fix unused variables
            if stripped.startswith('# TODO') or stripped.startswith('# FIXME'):
                fixes.append(f"Found TODO/FIXME at line {i+1}")

        return {"code": '\n'.join(lines), "fixes": fixes}

    def fix_indentation(self, lines: List[str]) -> List[str]:
        """Fix indentation issues"""
        fixed_lines = []
        indent_level = 0
        indent_size = 4

        for line in lines:
            stripped = line.strip()

            # Count opening brackets/braces
            open_brackets = line.count('{') + line.count('[') + line.count('(')
            close_brackets = line.count('}') + line.count(']') + line.count(')')

            # Adjust indent level
            indent_level += open_brackets
            indent_level -= close_brackets

            # Ensure non-negative indent
            indent_level = max(0, indent_level)

            # Apply indentation
            if stripped:
                fixed_lines.append(' ' * (indent_level * indent_size) + stripped)
            else:
                fixed_lines.append('')

        return fixed_lines

    def fix_missing_brackets(self, lines: List[str], error: SyntaxError) -> List[str]:
        """Fix missing brackets based on syntax error"""
        # This is a simplified implementation
        # In practice, you'd need more sophisticated parsing'

        error_line = error.lineno - 1 if error.lineno else 0
        if 0 <= error_line < len(lines):
            line = lines[error_line]

            # Try to fix common bracket issues
            if line.count('(') > line.count(')'):
                lines[error_line] = line + ')'
            elif line.count('[') > line.count(']'):
                lines[error_line] = line + ']'
            elif line.count('{') > line.count('}'):
                lines[error_line] = line + '}'

        return lines

    def get_function_body(self, lines: List[str], start_line: int) -> List[str]:
        """Get the body of a function starting at start_line"""
        body = []
        indent_level = 0

        for i in range(start_line + 1, len(lines)):
            line = lines[i]

            # Count indentation
            if line.strip():
                current_indent = len(line) - len(line.lstrip())
                if current_indent <= indent_level and body:
                    break
                indent_level = current_indent
            else:
                current_indent = indent_level

            body.append(line)

        return body

    def has_return_statement(self, lines: List[str]) -> bool:
        """Check if lines contain a return statement"""
        for line in lines:
            if line.strip().startswith('return'):
                return True
        return False

    def module_exists(self, module_name: str) -> bool:
        """Check if a module exists"""
        try:
            __import__(module_name)
            return True
        except ImportError:
            return False

    def validate_fixes(self, code: str) -> Dict:
        """Validate that fixes didn't break the code"""
        errors = []

        try:
            # Try to parse the code
            ast.parse(code)
        except SyntaxError as e:
            errors.append(f"Syntax error after fixes: {e}")

        # Check for other issues
        if 'eval(' in code:
            errors.append("Unsafe eval() still present")

        if 'password' in code.lower() and '***REDACTED***' not in code:
            errors.append("Hardcoded password may still be present")

        return {"valid": len(errors) == 0, "errors": errors}

# Example usage
async def fix_example_code():
    """Example of using the autofix agent"""
    sample_code = '''
def badFunction(name,age
    if name == "admin":
    password = "secret123"
        return True
    return False

x=1+2
'''

    agent = AutofixAgent()
    result = await agent.fix_code(sample_code)

    print("Fixes applied:")
    for fix in result["fixes_applied"]:
        print(f"  - {fix}")

    if result["errors_remaining"]:
        print("Errors remaining:")
        for error in result["errors_remaining"]:
            print(f"  - {error}")

    print("\nFixed code:")
    print(result["fixed_code"])

if __name__ == "__main__":
    import asyncio
    asyncio.run(fix_example_code())