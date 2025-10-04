#!/bin/bash

# MCP-CLI Toolkit Installation Script
# Universal installer for any device, any agent, any time

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
TOOLKIT_NAME="MCP-CLI Toolkit"
TOOLKIT_VERSION="1.0.0"
INSTALL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$HOME/.local/bin"
LIB_DIR="$HOME/.local/lib/mcp-cli-toolkit"
CONFIG_DIR="$HOME/.config/mcp-cli-toolkit"

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

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        "Darwin")
            echo "macos"
            ;;
        "Linux")
            echo "linux"
            ;;
        "CYGWIN"*|"MINGW"*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()
    local os_type=$(detect_os)

    # Core dependencies
    local core_deps=("python3" "pip3" "curl" "jq")

    # OS-specific dependencies
    case $os_type in
        "linux")
            core_deps+=("git")
            ;;
        "macos")
            core_deps+=("git" "brew")
            ;;
        "windows")
            core_deps+=("git")
            ;;
    esac

    for dep in "${core_deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_deps+=("$dep")
        fi
    done

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_warning "Missing dependencies: ${missing_deps[*]}"
        log_info "Please install missing dependencies and run this script again."
        return 1
    fi

    log_success "All dependencies satisfied"
    return 0
}

# Create necessary directories
create_directories() {
    log_info "Creating directories..."

    local dirs=(
        "$BIN_DIR"
        "$LIB_DIR"
        "$CONFIG_DIR"
        "$CONFIG_DIR/templates"
        "$CONFIG_DIR/profiles"
    )

    for dir in "${dirs[@]}"; do
        if [[ ! -d "$dir" ]]; then
            if mkdir -p "$dir"; then
                log_success "Created directory: $dir"
            else
                log_error "Failed to create directory: $dir"
                return 1
            fi
        fi
    done
}

# Install Python dependencies
install_python_deps() {
    log_info "Installing Python dependencies..."

    # Create virtual environment if it doesn't exist
    local venv_path="$LIB_DIR/venv"
    if [[ ! -d "$venv_path" ]]; then
        python3 -m venv "$venv_path"
        log_success "Created Python virtual environment"
    fi

    # Activate virtual environment and install dependencies
    source "$venv_path/bin/activate"
    pip install --upgrade pip

    # Install required packages
    local requirements=(
        "requests"
        "pyyaml"
        "colorama"
        "tqdm"
        "python-dotenv"
    )

    for package in "${requirements[@]}"; do
        if pip install "$package"; then
            log_success "Installed Python package: $package"
        else
            log_warning "Failed to install: $package"
        fi
    done

    log_success "Python dependencies installed"
}

# Create symbolic links
create_symlinks() {
    log_info "Creating symbolic links..."

    # Create main executable
    local main_script="$INSTALL_DIR/bin/mcp-cli-toolkit"
    if [[ ! -f "$main_script" ]]; then
        log_error "Main script not found: $main_script"
        return 1
    fi

    # Make script executable
    chmod +x "$main_script"

    # Create symlink in PATH
    local symlink_target="$BIN_DIR/mcp-cli-toolkit"
    if [[ -L "$symlink_target" ]]; then
        rm "$symlink_target"
    fi

    if ln -s "$main_script" "$symlink_target"; then
        log_success "Created symlink: $symlink_target"
    else
        log_error "Failed to create symlink"
        return 1
    fi
}

# Configure environment
configure_environment() {
    log_info "Configuring environment..."

    # Create environment file
    local env_file="$CONFIG_DIR/.env"
    if [[ ! -f "$env_file" ]]; then
        cat > "$env_file" << 'EOF'
# MCP-CLI Toolkit Environment Configuration
export MCP_CLI_TOOLKIT_HOME="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
export MCP_CLI_VENV_PATH="$MCP_CLI_TOOLKIT_HOME/lib/venv"
export MCP_CLI_CONFIG_PATH="$HOME/.config/mcp-cli-toolkit"
export MCP_CLI_LOG_LEVEL="INFO"
export MCP_CLI_DEFAULT_PROFILE="default"

# Add to PATH if not already present
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    export PATH="$HOME/.local/bin:$PATH"
fi
EOF
        log_success "Created environment configuration: $env_file"
    fi

    # Create default profile
    local default_profile="$CONFIG_DIR/profiles/default.json"
    if [[ ! -f "$default_profile" ]]; then
        cat > "$default_profile" << 'EOF'
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
        log_success "Created default profile: $default_profile"
    fi
}

# Setup shell integration
setup_shell_integration() {
    log_info "Setting up shell integration..."

    local shell_config=""
    local shell_name=""

    # Detect user's shell
    case "$SHELL" in
        *bash*)
            shell_config="$HOME/.bashrc"
            shell_name="bash"
            ;;
        *zsh*)
            shell_config="$HOME/.zshrc"
            shell_name="zsh"
            ;;
        *fish*)
            shell_config="$HOME/.config/fish/config.fish"
            shell_name="fish"
            ;;
        *)
            log_warning "Unsupported shell: $SHELL"
            return 0
            ;;
    esac

    if [[ -n "$shell_config" ]]; then
        # Check if already integrated
        if ! grep -q "mcp-cli-toolkit" "$shell_config" 2>/dev/null; then
            echo "" >> "$shell_config"
            echo "# MCP-CLI Toolkit Integration" >> "$shell_config"
            echo "source \"$CONFIG_DIR/.env\"" >> "$shell_config"
            echo "" >> "$shell_config"

            log_success "Added integration to $shell_name configuration"
            log_info "Please restart your shell or run: source $shell_config"
        else
            log_info "Shell integration already exists"
        fi
    fi
}

# Create activation script
create_activation_script() {
    log_info "Creating activation script..."

    local activate_script="$INSTALL_DIR/activate.sh"
    cat > "$activate_script" << 'EOF'
#!/bin/bash
# MCP-CLI Toolkit Activation Script
# Source this script to activate the toolkit in your current session

TOOLKIT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$HOME/.config/mcp-cli-toolkit"

# Source environment configuration
if [[ -f "$CONFIG_DIR/.env" ]]; then
    source "$CONFIG_DIR/.env"
    echo "✅ MCP-CLI Toolkit activated"
    echo "📂 Toolkit root: $TOOLKIT_ROOT"
    echo "🔧 Configuration: $CONFIG_DIR"
    echo ""
    echo "Available commands:"
    echo "  mcp-cli-toolkit --help    Show help information"
    echo "  mcp-cli-toolkit status    Show toolkit status"
    echo "  mcp-cli-toolkit list      List available tools"
    echo ""
else
    echo "❌ MCP-CLI Toolkit configuration not found"
    echo "Please run the installation script first"
    return 1
fi
EOF

    chmod +x "$activate_script"
    log_success "Created activation script: $activate_script"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."

    local issues=()

    # Check if main executable is accessible
    if ! command -v mcp-cli-toolkit &> /dev/null; then
        issues+=("Main executable not in PATH")
    fi

    # Check if configuration exists
    if [[ ! -f "$CONFIG_DIR/.env" ]]; then
        issues+=("Environment configuration missing")
    fi

    # Check if default profile exists
    if [[ ! -f "$CONFIG_DIR/profiles/default.json" ]]; then
        issues+=("Default profile missing")
    fi

    # Check if MCP servers are present
    if [[ ! -f "$INSTALL_DIR/mcp-servers/hybrid_deployment_server.py" ]]; then
        issues+=("MCP deployment server missing")
    fi

    if [[ ! -f "$INSTALL_DIR/mcp-servers/hybrid_testing_server.py" ]]; then
        issues+=("MCP testing server missing")
    fi

    if [[ ${#issues[@]} -eq 0 ]]; then
        log_success "Installation verification passed"
        return 0
    else
        log_warning "Installation issues found:"
        for issue in "${issues[@]}"; do
            log_warning "  - $issue"
        done
        return 1
    fi
}

# Main installation function
main() {
    print_header "🚀 $TOOLKIT_NAME Installation v$TOOLKIT_VERSION"

    local os_type=$(detect_os)
    log_info "Detected operating system: $os_type"
    log_info "Installation directory: $INSTALL_DIR"

    # Run installation steps
    if ! check_dependencies; then
        log_error "Dependency check failed"
        exit 1
    fi

    create_directories
    install_python_deps
    create_symlinks
    configure_environment
    setup_shell_integration
    create_activation_script

    if verify_installation; then
        print_header "🎉 Installation Completed Successfully!"

        echo ""
        echo "📋 Next Steps:"
        echo "1. Restart your shell or run: source ~/.bashrc (or ~/.zshrc)"
        echo "2. Activate the toolkit: source $INSTALL_DIR/activate.sh"
        echo "3. Run: mcp-cli-toolkit --help"
        echo ""
        echo "📂 Installation Details:"
        echo "   Root: $INSTALL_DIR"
        echo "   Bin:  $BIN_DIR"
        echo "   Config: $CONFIG_DIR"
        echo "   Virtual Env: $LIB_DIR/venv"
        echo ""
        echo "🔧 Quick Commands:"
        echo "   mcp-cli-toolkit status    - Show toolkit status"
        echo "   mcp-cli-toolkit list      - List available tools"
        echo "   mcp-cli-toolkit deploy    - Run deployment tools"
        echo "   mcp-cli-toolkit test      - Run testing tools"
        echo ""

        log_success "MCP-CLI Toolkit is ready to use!"
        exit 0
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

# Run main function
main "$@"