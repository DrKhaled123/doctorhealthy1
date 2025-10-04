#!/bin/bash

# Universal MCP-CLI Toolkit Activation Script
# Can be used by any AI agent on any device

set -euo pipefail

# Find the toolkit root (directory containing this script)
TOOLKIT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$HOME/.config/mcp-cli-toolkit"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

print_header() {
    echo -e "${BOLD}${BLUE}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "        $1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${NC}"
}

# Detect system information
detect_system() {
    local system_info=""

    system_info="${system_info}OS: $(uname -s) ($(uname -m))"
    system_info="${system_info}, Shell: $SHELL"
    system_info="${system_info}, User: $(whoami)"
    system_info="${system_info}, Date: $(date)"

    echo "$system_info"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    local missing_prereqs=()

    # Required commands
    local required_commands=("python3" "curl" "jq")
    for cmd in "${required_commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing_prereqs+=("$cmd")
        fi
    done

    if [[ ${#missing_prereqs[@]} -gt 0 ]]; then
        log_warning "Missing prerequisites: ${missing_prereqs[*]}"
        log_info "Some features may not work correctly."
        log_info "Consider installing missing dependencies."
        return 1
    fi

    log_success "All prerequisites satisfied"
    return 0
}

# Setup configuration
setup_configuration() {
    log_info "Setting up configuration..."

    # Create config directory
    mkdir -p "$CONFIG_DIR/profiles"

    # Create environment file if it doesn't exist
    if [[ ! -f "$CONFIG_DIR/.env" ]]; then
        cat > "$CONFIG_DIR/.env" << EOF
# MCP-CLI Toolkit Environment Configuration
export MCP_CLI_TOOLKIT_HOME="$TOOLKIT_ROOT"
export MCP_CLI_VENV_PATH="$TOOLKIT_ROOT/lib/venv"
export MCP_CLI_CONFIG_PATH="$CONFIG_DIR"
export MCP_CLI_LOG_LEVEL="INFO"
export MCP_CLI_DEFAULT_PROFILE="default"

# Add toolkit bin to PATH
if [[ ":$PATH:" != *":$TOOLKIT_ROOT/bin:"* ]]; then
    export PATH="$TOOLKIT_ROOT/bin:$PATH"
fi
EOF
        log_success "Created environment configuration"
    fi

    # Create default profile if it doesn't exist
    if [[ ! -f "$CONFIG_DIR/profiles/default.json" ]]; then
        cat > "$CONFIG_DIR/profiles/default.json" << 'EOF'
{
  "name": "default",
  "description": "Default MCP-CLI Toolkit profile",
  "version": "1.0.0",
  "created": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')",
  "mcp_servers": {
    "hybrid_deployment": {
      "enabled": true,
      "path": "mcp-servers/hybrid_deployment_server.py",
      "description": "Hybrid deployment server with cloud platform support"
    },
    "hybrid_testing": {
      "enabled": true,
      "path": "mcp-servers/hybrid_testing_server.py",
      "description": "Hybrid testing server with comprehensive testing tools"
    }
  },
  "cli_scripts": {
    "deploy_best_practices": {
      "enabled": true,
      "path": "cli-scripts/deploy_best_practices.sh",
      "description": "Comprehensive deployment best practices"
    },
    "deploy": {
      "enabled": true,
      "path": "cli-scripts/deploy.sh",
      "description": "Coolify deployment script"
    }
  },
  "settings": {
    "log_level": "INFO",
    "parallel_execution": true,
    "timeout": 300,
    "retry_attempts": 3
  }
}
EOF
        log_success "Created default profile"
    fi
}

# Setup Python virtual environment
setup_python_environment() {
    log_info "Setting up Python environment..."

    local venv_path="$TOOLKIT_ROOT/lib/venv"

    # Create virtual environment if it doesn't exist
    if [[ ! -d "$venv_path" ]]; then
        log_info "Creating Python virtual environment..."
        python3 -m venv "$venv_path"
        log_success "Created Python virtual environment"
    fi

    # Install/upgrade pip
    log_info "Upgrading pip..."
    "$venv_path/bin/pip" install --upgrade pip

    # Install required packages
    log_info "Installing Python dependencies..."
    local requirements=(
        "requests"
        "pyyaml"
        "colorama"
        "tqdm"
        "python-dotenv"
    )

    for package in "${requirements[@]}"; do
        if "$venv_path/bin/pip" install "$package" &> /dev/null; then
            log_success "Installed: $package"
        else
            log_warning "Failed to install: $package"
        fi
    done
}

# Make scripts executable
make_scripts_executable() {
    log_info "Making scripts executable..."

    # Find and make all shell scripts executable
    find "$TOOLKIT_ROOT" -name "*.sh" -type f -exec chmod +x {} \; 2>/dev/null || true

    # Make Python scripts executable
    find "$TOOLKIT_ROOT" -name "*.py" -type f -exec chmod +x {} \; 2>/dev/null || true

    # Make main binaries executable
    chmod +x "$TOOLKIT_ROOT/bin/mcp-cli-toolkit" 2>/dev/null || true
    chmod +x "$TOOLKIT_ROOT/bin/mcp-config" 2>/dev/null || true

    log_success "Scripts made executable"
}

# Test installation
test_installation() {
    log_info "Testing installation..."

    local tests_passed=0
    local total_tests=0

    # Test 1: Main executable
    ((total_tests++))
    if command -v mcp-cli-toolkit &> /dev/null || [[ -x "$TOOLKIT_ROOT/bin/mcp-cli-toolkit" ]]; then
        ((tests_passed++))
        log_success "Main executable test passed"
    else
        log_warning "Main executable test failed"
    fi

    # Test 2: Configuration
    ((total_tests++))
    if [[ -f "$CONFIG_DIR/.env" ]] && [[ -f "$CONFIG_DIR/profiles/default.json" ]]; then
        ((tests_passed++))
        log_success "Configuration test passed"
    else
        log_warning "Configuration test failed"
    fi

    # Test 3: MCP servers
    ((total_tests++))
    if [[ -f "$TOOLKIT_ROOT/mcp-servers/hybrid_deployment_server.py" ]]; then
        ((tests_passed++))
        log_success "MCP servers test passed"
    else
        log_warning "MCP servers test failed"
    fi

    # Test 4: CLI scripts
    ((total_tests++))
    if [[ -f "$TOOLKIT_ROOT/cli-scripts/deploy_best_practices.sh" ]]; then
        ((tests_passed++))
        log_success "CLI scripts test passed"
    else
        log_warning "CLI scripts test failed"
    fi

    # Test 5: Python environment
    ((total_tests++))
    if [[ -d "$TOOLKIT_ROOT/lib/venv" ]]; then
        ((tests_passed++))
        log_success "Python environment test passed"
    else
        log_warning "Python environment test failed"
    fi

    local success_rate=$((tests_passed * 100 / total_tests))
    log_info "Test results: $tests_passed/$total_tests passed ($success_rate%)"

    if [[ $success_rate -ge 80 ]]; then
        log_success "Installation test completed successfully"
        return 0
    else
        log_warning "Installation test completed with issues"
        return 1
    fi
}

# Show final status
show_final_status() {
    print_header "🎉 MCP-CLI Toolkit Activated Successfully!"

    echo ""
    echo "📋 System Information:"
    detect_system
    echo ""
    echo "📂 Installation Details:"
    echo "   Root: $TOOLKIT_ROOT"
    echo "   Config: $CONFIG_DIR"
    echo "   Profile: default"
    echo "   Version: 1.0.0"
    echo ""
    echo "🤖 Available MCP Servers:"
    echo "   • hybrid_deployment_server.py - Cloud deployment automation"
    echo "   • hybrid_testing_server.py - Comprehensive testing framework"
    echo ""
    echo "🛠️  Available CLI Scripts:"
    echo "   • deploy_best_practices.sh - Deployment validation & automation"
    echo "   • deploy.sh - Coolify platform deployment"
    echo "   • deployment_monitoring.py - Real-time monitoring"
    echo "   • deployment_error_tracker.py - Error tracking & analysis"
    echo ""
    echo "💡 Quick Start Commands:"
    echo "   mcp-cli-toolkit status              # Show toolkit status"
    echo "   mcp-cli-toolkit list                # List available tools"
    echo "   mcp-cli-toolkit mcp hybrid_deployment  # Start deployment server"
    echo "   mcp-cli-toolkit cli deploy_best_practices.sh validate"
    echo ""
    echo "🔧 Management Commands:"
    echo "   mcp-config profile list             # List profiles"
    echo "   mcp-config profile create <name>    # Create new profile"
    echo "   mcp-config env show                 # Show environment"
    echo "   mcp-config env set <key> <value>   # Set environment variable"
    echo ""
    echo "📚 Documentation:"
    echo "   cat $TOOLKIT_ROOT/README.md         # Full documentation"
    echo "   mcp-cli-toolkit help                # Command help"
    echo ""
    echo "🔗 Need Help?"
    echo "   • Check the README: $TOOLKIT_ROOT/README.md"
    echo "   • Run: mcp-cli-toolkit help"
    echo "   • Run: mcp-config --help"
    echo ""
}

# Main activation function
main() {
    print_header "🚀 Activating MCP-CLI Toolkit"

    log_info "Starting activation process..."
    log_info "Toolkit root: $TOOLKIT_ROOT"

    # Run activation steps
    check_prerequisites
    setup_configuration
    setup_python_environment
    make_scripts_executable
    test_installation

    # Source the environment to make commands available
    if [[ -f "$CONFIG_DIR/.env" ]]; then
        source "$CONFIG_DIR/.env"
        log_success "Environment configuration loaded"
    fi

    show_final_status

    log_success "MCP-CLI Toolkit activation completed!"
    echo ""
    echo "🎯 Next Steps:"
    echo "1. Run: mcp-cli-toolkit status"
    echo "2. Run: mcp-cli-toolkit list"
    echo "3. Start using any of the available tools!"
    echo ""
}

# Run main function
main "$@"