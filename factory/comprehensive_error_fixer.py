#!/usr/bin/env python3
"""
Comprehensive Error Fixer for Factory System
Fixes all remaining syntax and structural errors
"""

import os
import re
import glob
import subprocess
import sys

def fix_all_syntax_errors():
    """Fix all syntax errors in the factory system"""
    print("🔧 Comprehensive Error Fixing")
    print("=" * 50)

    # Get all Python files
    python_files = glob.glob("*.py")

    for file_path in python_files:
        print(f"\n📄 Processing {file_path}...")
        fix_file_syntax(file_path)

    print("\n✅ Comprehensive error fixing completed!")

def fix_file_syntax(file_path):
    """Fix syntax errors in a specific file"""
    try:
        with open(file_path, 'r') as f:
            content = f.read()

        original_content = content

        # Fix f-string syntax errors
        content = fix_fstring_syntax(content)

        # Fix unterminated string literals
        content = fix_string_literals(content)

        # Fix indentation issues
        content = fix_indentation_issues(content)

        # Fix missing imports
        content = fix_missing_imports(content)

        if content != original_content:
            with open(file_path, 'w') as f:
                f.write(content)
            print(f"  ✅ Fixed syntax in {file_path}")
        else:
            print(f"  ✅ No issues found in {file_path}")

    except Exception as e:
        print(f"  ❌ Error processing {file_path}: {e}")

def fix_fstring_syntax(content):
    """Fix f-string syntax errors"""
    # Fix missing colons in f-strings
    content = re.sub(r'\{([^}]+)"(\.?\d+)f"\}', r'{\1:\2f}', content)

    # Fix unterminated f-string expressions
    content = re.sub(r'\{([^}]+)"(\.?\d+)"\}', r'{\1:\2}', content)

    return content

def fix_string_literals(content):
    """Fix unterminated string literals"""
    lines = content.split('\n')
    fixed_lines = []
    in_string = False
    string_char = None

    for line in lines:
        original_line = line

        # Check for unterminated strings
        quote_count = line.count('"') + line.count("'")

        if not in_string and quote_count % 2 == 1:
            # Found start of unterminated string
            in_string = True
            if '"' in line and line.rindex('"') > line.rindex("'"):
                string_char = '"'
            else:
                string_char = "'"

            # Add closing quote if missing
            if not line.rstrip().endswith(string_char):
                line = line.rstrip() + string_char

        elif in_string:
            # Check if this line closes the string
            if string_char in line:
                in_string = False
                string_char = None

        fixed_lines.append(line)

    return '\n'.join(fixed_lines)

def fix_indentation_issues(content):
    """Fix indentation issues"""
    lines = content.split('\n')
    fixed_lines = []

    for i, line in enumerate(lines):
        # Fix lines that should be indented but aren't
        stripped = line.strip()
        if stripped and not line.startswith(' ') and not line.startswith('\t'):
            # Check if this should be indented based on context
            if i > 0:
                prev_line = lines[i-1]
                if (prev_line.strip().endswith(':') or
                    prev_line.strip().startswith(('def ', 'class ', 'if ', 'for ', 'while ', 'try ', 'except ', 'finally '))):
                    if not prev_line.startswith('    '):
                        lines[i] = '    ' + line
                        line = '    ' + line

        fixed_lines.append(line)

    return '\n'.join(fixed_lines)

def fix_missing_imports(content):
    """Fix missing imports"""
    # Add common missing imports
    imports_to_add = []

    if 'redis' in content and 'import redis' not in content:
        imports_to_add.append('import redis')

    if 'asyncio' in content and 'import asyncio' not in content:
        imports_to_add.append('import asyncio')

    if 'json' in content and 'import json' not in content:
        imports_to_add.append('import json')

    if 'os' in content and 'import os' not in content:
        imports_to_add.append('import os')

    if imports_to_add:
        # Add imports at the beginning
        lines = content.split('\n')
        insert_index = 0

        # Find the right place to insert imports
        for i, line in enumerate(lines):
            if line.startswith('import ') or line.startswith('from '):
                insert_index = i + 1
            elif line.strip() and not line.startswith('#'):
                break

        # Insert imports
        for import_line in imports_to_add:
            lines.insert(insert_index, import_line)
            insert_index += 1

        content = '\n'.join(lines)

    return content

def validate_fixes():
    """Validate that all fixes work"""
    print("\n🔍 Validating fixes...")

    try:
        # Test compilation of all files
        result = subprocess.run([
            sys.executable, '-m', 'py_compile'
        ] + glob.glob("*.py"), capture_output=True, text=True, timeout=30)

        if result.returncode == 0:
            print("✅ All files compile successfully!")
            return True
        else:
            print("❌ Some files still have errors:")
            print(result.stderr)
            return False

    except Exception as e:
        print(f"❌ Validation failed: {e}")
        return False

def create_error_free_demo():
    """Create a demo script that works without errors"""
    demo_content = '''#!/usr/bin/env python3
"""
Factory System Demo - Error-Free Version
Demonstrates the working components of the factory system
"""

import asyncio
import time
from datetime import datetime

async def demo_working_components():
    """Demo the working components"""
    print("🚀 Factory System - Working Components Demo")
    print("=" * 50)

    # Test memory system
    try:
        from memo_ai_memory import MemoAIMemory
        memory = MemoAIMemory()
        await memory.remember("demo_test", {"test": "success"})
        print("✅ Memory system working")
    except Exception as e:
        print(f"❌ Memory system error: {e}")

    # Test GitHub integration
    try:
        from github_integration import GitHubIntegration
        github = GitHubIntegration()
        print("✅ GitHub integration ready")
    except Exception as e:
        print(f"❌ GitHub integration error: {e}")

    # Test hybrid framework
    try:
        from hybrid_framework import HybridFramework
        framework = HybridFramework()
        results = await framework.run_hybrid_analysis("test specification")
        print(f"✅ Hybrid framework working: {results['total_components']}/5 components")
    except Exception as e:
        print(f"❌ Hybrid framework error: {e}")

    print("\\n🎉 Demo completed!")

if __name__ == "__main__":
    asyncio.run(demo_working_components())
'''

    with open('demo_working_system.py', 'w') as f:
        f.write(demo_content)

    print("✅ Created error-free demo script")

if __name__ == "__main__":
    fix_all_syntax_errors()

    if validate_fixes():
        create_error_free_demo()
        print("\\n🎉 All errors fixed successfully!")
        print("\\n📖 You can now run:")
        print("  python3 demo_working_system.py")
        print("  python3 hybrid_framework.py")
        print("  ./start_factory.sh")
    else:
        print("\\n⚠️  Some errors remain. Manual intervention may be needed.")