#!/bin/bash

# MCP-CLI Toolkit Packaging Script
# Creates portable packages for distribution across devices

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
TOOLKIT_NAME="mcp-cli-toolkit"
VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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

# Detect package format
detect_package_format() {
    if command -v tar &> /dev/null; then
        echo "tar"
    elif command -v zip &> /dev/null; then
        echo "zip"
    else
        echo "none"
    fi
}

# Create tar package
create_tar_package() {
    local package_name="$TOOLKIT_NAME-$VERSION.tar.gz"
    local package_path="$SCRIPT_DIR/dist/$package_name"

    log_info "Creating tar package: $package_name"

    # Create dist directory
    mkdir -p "$SCRIPT_DIR/dist"

    # Create tar package
    cd "$SCRIPT_DIR/.."
    tar -czf "$package_path" \
        --exclude=".git" \
        --exclude="__pycache__" \
        --exclude="*.pyc" \
        --exclude=".DS_Store" \
        --exclude="dist" \
        --exclude="*.log" \
        "$(basename "$SCRIPT_DIR")/"

    log_success "Created tar package: $package_path"
    echo "$package_path"
}

# Create zip package
create_zip_package() {
    local package_name="$TOOLKIT_NAME-$VERSION.zip"
    local package_path="$SCRIPT_DIR/dist/$package_name"

    log_info "Creating zip package: $package_name"

    # Create dist directory
    mkdir -p "$SCRIPT_DIR/dist"

    # Create zip package
    cd "$SCRIPT_DIR/.."
    zip -r "$package_path" \
        -x "*.git*" \
        -x "__pycache__/*" \
        -x "*.pyc" \
        -x ".DS_Store" \
        -x "dist/*" \
        -x "*.log" \
        "$(basename "$SCRIPT_DIR")/"

    log_success "Created zip package: $package_path"
    echo "$package_path"
}

# Create installation script
create_installation_script() {
    local package_path="$1"
    local install_script="$SCRIPT_DIR/dist/install-from-package.sh"

    log_info "Creating installation script"

    cat > "$install_script" << 'EOF'
#!/bin/bash
# Universal Installation Script for MCP-CLI Toolkit
# Works on any device with bash

set -euo pipefail

PACKAGE_FILE="$1"
INSTALL_DIR="${2:-$HOME/.local/lib/mcp-cli-toolkit}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

main() {
    if [[ -z "$PACKAGE_FILE" ]]; then
        echo "Usage: $0 <package_file> [install_directory]"
        echo "Example: $0 mcp-cli-toolkit-1.0.0.tar.gz"
        exit 1
    fi

    if [[ ! -f "$PACKAGE_FILE" ]]; then
        log_error "Package file not found: $PACKAGE_FILE"
        exit 1
    fi

    log_info "Installing MCP-CLI Toolkit..."
    log_info "Package: $PACKAGE_FILE"
    log_info "Install directory: $INSTALL_DIR"

    # Create install directory
    mkdir -p "$INSTALL_DIR"

    # Extract package
    log_info "Extracting package..."
    if [[ "$PACKAGE_FILE" == *.tar.gz ]]; then
        tar -xzf "$PACKAGE_FILE" -C "$INSTALL_DIR"
    elif [[ "$PACKAGE_FILE" == *.zip ]]; then
        unzip -q "$PACKAGE_FILE" -d "$INSTALL_DIR"
    else
        log_error "Unsupported package format"
        exit 1
    fi

    # Make scripts executable
    find "$INSTALL_DIR" -name "*.sh" -exec chmod +x {} \;

    # Run installation
    log_info "Running installation..."
    cd "$INSTALL_DIR"
    ./install.sh

    log_success "Installation completed!"
    echo ""
    echo "Next steps:"
    echo "1. Restart your shell or run: source ~/.bashrc"
    echo "2. Activate toolkit: source $INSTALL_DIR/activate.sh"
    echo "3. Run: mcp-cli-toolkit status"
}

main "$@"
EOF

    chmod +x "$install_script"
    log_success "Created installation script: $install_script"
    echo "$install_script"
}

# Create portable activation script
create_portable_activation() {
    local activation_script="$SCRIPT_DIR/dist/activate-anywhere.sh"

    log_info "Creating portable activation script"

    cat > "$activation_script" << 'EOF'
#!/bin/bash
# Portable Activation Script for MCP-CLI Toolkit
# Can be run from any directory on any device

set -euo pipefail

# Find the toolkit root (directory containing this script)
TOOLKIT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_DIR="$HOME/.config/mcp-cli-toolkit"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

print_header() {
    echo -e "${BOLD}${BLUE}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "        $1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${NC}"
}

main() {
    print_header "🚀 MCP-CLI Toolkit Activation"

    log_info "Toolkit root: $TOOLKIT_ROOT"
    log_info "Configuration directory: $CONFIG_DIR"

    # Create config directory if it doesn't exist
    mkdir -p "$CONFIG_DIR/profiles"

    # Create environment configuration if it doesn't exist
    if [[ ! -f "$CONFIG_DIR/.env" ]]; then
        log_info "Creating environment configuration..."

        cat > "$CONFIG_DIR/.env" << EOF
# MCP-CLI Toolkit Environment Configuration
export MCP_CLI_TOOLKIT_HOME="$TOOLKIT_ROOT"
export MCP_CLI_VENV_PATH="$TOOLKIT_ROOT/lib/venv"
export MCP_CLI_CONFIG_PATH="$CONFIG_DIR"
export MCP_CLI_LOG_LEVEL="INFO"
export MCP_CLI_DEFAULT_PROFILE="default"

# Add to PATH if not already present
if [[ ":$PATH:" != *":$TOOLKIT_ROOT/bin:"* ]]; then
    export PATH="$TOOLKIT_ROOT/bin:$PATH"
fi
EOF
        log_success "Created environment configuration"
    fi

    # Create default profile if it doesn't exist
    if [[ ! -f "$CONFIG_DIR/profiles/default.json" ]]; then
        log_info "Creating default profile..."

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

    # Source environment configuration
    source "$CONFIG_DIR/.env"

    print_header "✅ Activation Completed!"

    echo ""
    echo "📋 Toolkit Information:"
    echo "   Version: 1.0.0"
    echo "   Root: $TOOLKIT_ROOT"
    echo "   Profile: $MCP_CLI_DEFAULT_PROFILE"
    echo "   Config: $CONFIG_DIR"
    echo ""
    echo "🛠️  Available Commands:"
    echo "   mcp-cli-toolkit status    - Show toolkit status"
    echo "   mcp-cli-toolkit list      - List available tools"
    echo "   mcp-cli-toolkit mcp       - Run MCP servers"
    echo "   mcp-cli-toolkit cli       - Run CLI scripts"
    echo "   mcp-config profile        - Manage profiles"
    echo "   mcp-config env           - Manage environment"
    echo ""
    echo "💡 Quick Start:"
    echo "   mcp-cli-toolkit list"
    echo "   mcp-cli-toolkit mcp hybrid_deployment"
    echo "   mcp-cli-toolkit cli deploy_best_practices.sh --help"
    echo ""
    echo "🔧 Need Help?"
    echo "   mcp-cli-toolkit help"
    echo "   cat $TOOLKIT_ROOT/README.md"
    echo ""

    log_success "MCP-CLI Toolkit is ready to use!"
}

# Run main function
main "$@"
EOF

    chmod +x "$activation_script"
    log_success "Created portable activation script: $activation_script"
    echo "$activation_script"
}

# Create distribution manifest
create_manifest() {
    local manifest_file="$SCRIPT_DIR/dist/manifest.json"

    log_info "Creating distribution manifest"

    cat > "$manifest_file" << EOF
{
  "name": "mcp-cli-toolkit",
  "version": "$VERSION",
  "description": "Universal toolkit for MCP servers and CLI tools",
  "created": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')",
  "package_type": "source",
  "supported_platforms": [
    "linux",
    "macos",
    "windows"
  ],
  "dependencies": {
    "python": "3.7+",
    "bash": "4.0+",
    "curl": "7.0+",
    "jq": "1.6+"
  },
  "installation": {
    "script": "install.sh",
    "activation": "activate-anywhere.sh",
    "configuration": "automatic"
  },
  "components": {
    "mcp_servers": [
      "hybrid_deployment_server.py",
      "hybrid_testing_server.py"
    ],
    "cli_scripts": [
      "deploy_best_practices.sh",
      "deploy.sh",
      "deployment_monitoring.py",
      "deployment_error_tracker.py",
      "deployment_error_prevention.py",
      "pre_deployment_validator.py"
    ],
    "binaries": [
      "mcp-cli-toolkit",
      "mcp-config"
    ]
  },
  "features": [
    "Multi-cloud deployment",
    "Comprehensive testing",
    "Security validation",
    "Performance monitoring",
    "Error prevention",
    "Agent agnostic",
    "Cross-platform"
  ]
}
EOF

    log_success "Created manifest: $manifest_file"
    echo "$manifest_file"
}

# Main packaging function
main() {
    print_header "📦 MCP-CLI Toolkit Packaging v$VERSION"

    log_info "Script directory: $SCRIPT_DIR"
    log_info "Toolkit version: $VERSION"

    # Check if required tools are available
    local package_format=$(detect_package_format)
    if [[ "$package_format" == "none" ]]; then
        log_error "No packaging tools found (tar or zip required)"
        exit 1
    fi

    log_info "Detected package format: $package_format"

    # Create packages
    local packages=()

    if [[ "$package_format" == "tar" ]]; then
        packages+=("$(create_tar_package)")
    fi

    if [[ "$package_format" == "zip" ]]; then
        packages+=("$(create_zip_package)")
    fi

    # Create installation script
    if [[ ${#packages[@]} -gt 0 ]]; then
        local install_script
        install_script="$(create_installation_script "${packages[0]}")"
        packages+=("$install_script")
    fi

    # Create portable activation
    local activation_script
    activation_script="$(create_portable_activation)"
    packages+=("$activation_script")

    # Create manifest
    local manifest
    manifest="$(create_manifest)"
    packages+=("$manifest")

    print_header "📦 Packaging Completed!"

    echo ""
    echo "📋 Created packages:"
    for package in "${packages[@]}"; do
        local size
        size=$(du -h "$package" | cut -f1)
        echo "   📦 $(basename "$package") ($size)"
    done
    echo ""
    echo "🚀 Distribution:"
    echo "   - Copy the 'dist' directory to any device"
    echo "   - Run: ./dist/activate-anywhere.sh"
    echo "   - Or: ./dist/install-from-package.sh <package_file>"
    echo ""
    echo "💡 Usage on any device:"
    echo "   # 1. Copy dist/ to device"
    echo "   # 2. cd dist/"
    echo "   # 3. ./activate-anywhere.sh"
    echo "   # 4. mcp-cli-toolkit status"
    echo ""

    log_success "Packaging completed successfully!"
}

# Run main function
main "$@"