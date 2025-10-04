#!/bin/bash

# Git Auto-Commit and Push Script
# Provides convenient commands for automated git workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not in a git repository"
        print_error "Please navigate to your project directory and try again"
        exit 1
    fi
}

# Function to get current branch
get_current_branch() {
    git branch --show-current
}

# Function to auto-commit and push
auto_commit_push() {
    local message="$1"

    if [ -z "$message" ]; then
        print_error "Commit message is required"
        echo "Usage: git-auto commit 'Your commit message'"
        exit 1
    fi

    check_git_repo

    # Add all changes
    print_status "Adding all changes..."
    git add .

    # Check if there are changes to commit
    if git diff --cached --quiet; then
        print_warning "No changes to commit"
        return 0
    fi

    # Commit with message
    print_status "Committing with message: '$message'"
    git commit -m "$message"

    # Push to current branch
    local branch=$(get_current_branch)
    print_status "Pushing to origin/$branch..."
    git push origin "$branch"

    print_success "Changes committed and pushed successfully!"
}

# Function to quick commit with timestamp
quick_commit() {
    local message="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    if [ -z "$message" ]; then
        message="Auto-commit at $timestamp"
    else
        message="$message - $timestamp"
    fi

    auto_commit_push "$message"
}

# Function to sync with remote (pull and push)
sync_remote() {
    check_git_repo

    local branch=$(get_current_branch)

    print_status "Fetching latest changes..."
    git fetch origin

    print_status "Pulling changes from origin/$branch..."
    git pull origin "$branch"

    print_status "Pushing local changes to origin/$branch..."
    git push origin "$branch"

    print_success "Repository synchronized!"
}

# Function to show status with colors
show_status() {
    check_git_repo

    print_status "Current branch: $(get_current_branch)"
    print_status "Git status:"
    git status --short
}

# Function to show recent commits
show_recent() {
    check_git_repo

    local count="${1:-5}"
    print_status "Recent $count commits:"
    git log --oneline -$count
}

# Function to create backup branch
backup_branch() {
    check_git_repo

    local timestamp=$(date '+%Y%m%d_%H%M%S')
    local backup_branch="backup_$timestamp"

    print_status "Creating backup branch: $backup_branch"
    git checkout -b "$backup_branch"

    print_success "Backup branch created: $backup_branch"
}

# Help function
show_help() {
    echo "Git Auto-Commit and Push Script"
    echo ""
    echo "Usage: git-auto <command> [arguments]"
    echo ""
    echo "Commands:"
    echo "  commit <message>    - Add all changes, commit with message, and push"
    echo "  quick [message]     - Quick commit with timestamp"
    echo "  sync               - Sync with remote (pull and push)"
    echo "  status             - Show git status"
    echo "  recent [count]     - Show recent commits (default: 5)"
    echo "  backup             - Create a backup branch"
    echo "  help               - Show this help message"
    echo ""
    echo "Examples:"
    echo "  git-auto commit 'Add new feature'"
    echo "  git-auto quick 'WIP'"
    echo "  git-auto sync"
    echo "  git-auto recent 10"
}

# Main script logic
main() {
    local command="$1"

    case "$command" in
        "commit")
            auto_commit_push "$2"
            ;;
        "quick")
            quick_commit "$2"
            ;;
        "sync")
            sync_remote
            ;;
        "status")
            show_status
            ;;
        "recent")
            show_recent "$2"
            ;;
        "backup")
            backup_branch
            ;;
        "help"|"-h"|"--help"|"")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
