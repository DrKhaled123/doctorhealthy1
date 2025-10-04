#!/bin/bash

# Setup script for Git Automation
# This script sets up convenient aliases and configurations for automated git workflow

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Detect shell configuration file
detect_shell() {
    if [ -n "$ZSH_VERSION" ]; then
        echo "$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        echo "$HOME/.bashrc"
    else
        echo "$HOME/.profile"
    fi
}

# Add alias to shell configuration
add_alias() {
    local shell_config=$(detect_shell)
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local git_auto_path="$script_dir/git-auto.sh"

    # Create alias command
    local alias_command="alias git-auto='. $git_auto_path'"

    # Check if alias already exists
    if grep -q "alias git-auto" "$shell_config" 2>/dev/null; then
        print_warning "git-auto alias already exists in $shell_config"
        return 0
    fi

    # Add alias to shell configuration
    echo "" >> "$shell_config"
    echo "# Git automation alias" >> "$shell_config"
    echo "$alias_command" >> "$shell_config"

    print_success "Added git-auto alias to $shell_config"
    print_info "Please restart your terminal or run: source $shell_config"
}

# Create local bin directory if it doesn't exist
create_local_bin() {
    local local_bin="$HOME/.local/bin"

    if [ ! -d "$local_bin" ]; then
        mkdir -p "$local_bin"
        print_info "Created $local_bin directory"
    fi

    # Add to PATH if not already there
    local shell_config=$(detect_shell)
    if ! grep -q "$local_bin" "$shell_config" 2>/dev/null; then
        echo "" >> "$shell_config"
        echo "# Add local bin to PATH" >> "$shell_config"
        echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$shell_config"
        print_info "Added $local_bin to PATH in $shell_config"
    fi
}

# Create symlink in local bin
create_symlink() {
    local local_bin="$HOME/.local/bin"
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local git_auto_path="$script_dir/git-auto.sh"
    local symlink_path="$local_bin/git-auto"

    if [ -L "$symlink_path" ]; then
        print_warning "Symlink already exists at $symlink_path"
        return 0
    fi

    ln -sf "$git_auto_path" "$symlink_path"
    print_success "Created symlink at $symlink_path"
}

# Main setup function
main() {
    print_info "Setting up Git automation..."

    # Add alias to shell configuration
    add_alias

    # Create local bin and symlink
    create_local_bin
    create_symlink

    print_success "Git automation setup complete!"
    echo ""
    print_info "Available commands:"
    echo "  git-auto commit 'message'  - Auto-commit and push"
    echo "  git-auto quick [message]   - Quick commit with timestamp"
    echo "  git-auto sync             - Sync with remote"
    echo "  git-auto status           - Show git status"
    echo "  git-auto recent [count]   - Show recent commits"
    echo "  git-auto backup           - Create backup branch"
    echo ""
    print_info "Examples:"
    echo "  git-auto commit 'Add new feature'"
    echo "  git-auto quick 'WIP'"
    echo "  git-auto sync"
    echo ""
    print_warning "Please restart your terminal or source your shell config:"
    echo "  source $(detect_shell)"
}

# Run main function
main "$@"
