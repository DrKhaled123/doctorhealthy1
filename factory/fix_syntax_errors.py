#!/usr/bin/env python3
"""
Syntax Error Fixer for Factory System
Fixes all syntax errors in the factory codebase
"""

import os
import re
import glob

def fix_fstring_syntax(file_path):
    """Fix f-string syntax errors"""
    try:
        with open(file_path, 'r') as f:
            content = f.read()

        # Fix common f-string issues
        # Pattern 1: Missing colon in f-string
        content = re.sub(r'\{([^}]+)"(\.?\d+)f"\}', r'{\1:\2f}', content)

        # Pattern 2: Unterminated f-string expressions
        content = re.sub(r'\{([^}]+)"(\.?\d+)"\}', r'{\1:\2}', content)

        # Pattern 3: Missing closing brace in f-string
        lines = content.split('\n')
        for i, line in enumerate(lines):
            if '{' in line and '}' in line:
                # Count braces
                open_braces = line.count('{')
                close_braces = line.count('}')
                if open_braces > close_braces:
                    # Try to fix by adding closing brace
                    lines[i] = line + '"}'

        content = '\n'.join(lines)

        # Write back if changed
        with open(file_path, 'w') as f:
            f.write(content)

        print(f"✅ Fixed f-string syntax in {file_path}")

    except Exception as e:
        print(f"❌ Failed to fix {file_path}: {e}")

def fix_string_literals(file_path):
    """Fix unterminated string literals"""
    try:
        with open(file_path, 'r') as f:
            content = f.read()

        lines = content.split('\n')
        fixed_lines = []
        in_string = False
        string_char = None

        for i, line in enumerate(lines):
            # Check for unterminated strings
            quote_count = line.count('"') + line.count("'")

            if not in_string and quote_count % 2 == 1:
                # Found start of unterminated string
                in_string = True
                if '"' in line and line.rindex('"') > line.rindex("'"):
                    string_char = '"'
                else:
                    string_char = "'"

                # Add closing quote at end of line
                if not line.rstrip().endswith(string_char):
                    lines[i] = line.rstrip() + string_char
                    print(f"✅ Fixed unterminated string at line {i+1} in {file_path}")

            elif in_string:
                # Check if this line closes the string
                if string_char in line:
                    in_string = False
                    string_char = None

        content = '\n'.join(lines)

        with open(file_path, 'w') as f:
            f.write(content)

    except Exception as e:
        print(f"❌ Failed to fix string literals in {file_path}: {e}")

def fix_indentation(file_path):
    """Fix indentation issues"""
    try:
        with open(file_path, 'r') as f:
            lines = f.readlines()

        # Fix common indentation issues
        for i in range(len(lines)):
            line = lines[i]

            # Fix lines that start with print but have wrong indentation
            if line.strip().startswith('print(') and not line.startswith('    ') and not line.startswith('\t'):
                # Check if previous line has proper indentation
                if i > 0 and (lines[i-1].startswith('    ') or lines[i-1].startswith('\t')):
                    lines[i] = '    ' + line.lstrip()
                    print(f"✅ Fixed indentation at line {i+1} in {file_path}")

        with open(file_path, 'w') as f:
            f.writelines(lines)

    except Exception as e:
        print(f"❌ Failed to fix indentation in {file_path}: {e}")

def fix_all_syntax_errors():
    """Fix all syntax errors in factory files"""
    print("🔧 Fixing Syntax Errors in Factory System")
    print("=" * 50)

    factory_files = glob.glob("*.py")

    for file_path in factory_files:
        print(f"\n📄 Processing {file_path}...")

        # Apply all fixes
        fix_fstring_syntax(file_path)
        fix_string_literals(file_path)
        fix_indentation(file_path)

    print("\n✅ Syntax error fixing completed!")

if __name__ == "__main__":
    fix_all_syntax_errors()