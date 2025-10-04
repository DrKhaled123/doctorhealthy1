# MCP-CLI Toolkit

🚀 **Universal Toolkit for MCP Servers and CLI Tools**

A portable, agent-agnostic toolkit that can be used by any AI agent, on any device, at any time. Combines Model Context Protocol (MCP) servers with traditional CLI tools for comprehensive development and deployment workflows.

## 🌟 Features

- **🤖 MCP Servers**: Ready-to-use Model Context Protocol servers
- **🛠️ CLI Tools**: Comprehensive command-line utilities
- **📦 Portable**: Works on any device with zero configuration
- **👥 Agent Agnostic**: Compatible with any AI agent or system
- **⚙️ Configurable**: Flexible profiles and environment management
- **🔧 Easy Installation**: One-command setup on any platform

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Available Tools](#available-tools)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)

## 🚀 Quick Start

### 1. Install the Toolkit

```bash
# Download and install
curl -fsSL https://raw.githubusercontent.com/your-repo/mcp-cli-toolkit/main/install.sh | bash

# Or clone and install
git clone https://github.com/your-repo/mcp-cli-toolkit.git
cd mcp-cli-toolkit
./install.sh
```

### 2. Activate the Toolkit

```bash
# Activate in current session
source mcp-cli-toolkit/activate.sh

# Or use the main command
mcp-cli-toolkit status
```

### 3. Use the Tools

```bash
# List available tools
mcp-cli-toolkit list

# Run MCP server
mcp-cli-toolkit mcp hybrid_deployment

# Run CLI script
mcp-cli-toolkit cli deploy_best_practices.sh --help
```

## 📦 Installation

### System Requirements

- **Python 3.7+** (with pip)
- **curl** (for downloading)
- **jq** (for JSON processing)
- **git** (optional, for version control)

### Automated Installation

```bash
# Full installation with all dependencies
./install.sh
```

### Manual Installation

```bash
# 1. Create directories
mkdir -p ~/.local/bin ~/.local/lib/mcp-cli-toolkit ~/.config/mcp-cli-toolkit/profiles

# 2. Copy files
cp -r mcp-cli-toolkit/* ~/.local/lib/mcp-cli-toolkit/
cp mcp-cli-toolkit/bin/mcp-cli-toolkit ~/.local/bin/
cp mcp-cli-toolkit/bin/mcp-config ~/.local/bin/

# 3. Make executable
chmod +x ~/.local/bin/mcp-cli-toolkit
chmod +x ~/.local/bin/mcp-config

# 4. Configure environment
source ~/.local/lib/mcp-cli-toolkit/install.sh configure_environment
```

### Platform-Specific Notes

#### Linux (Ubuntu/Debian)
```bash
# Install system dependencies
sudo apt update
sudo apt install -y python3 python3-pip python3-venv curl jq git
```

#### macOS
```bash
# Install system dependencies
brew install python3 curl jq git
```

#### Windows (WSL/Git Bash)
```bash
# Install system dependencies
sudo apt update
sudo apt install -y python3 python3-pip python3-venv curl jq git
```

## 🎯 Usage

### Basic Commands

```bash
# Show toolkit status
mcp-cli-toolkit status

# List available tools
mcp-cli-toolkit list

# Show help
mcp-cli-toolkit help
```

### MCP Server Operations

```bash
# Start deployment server
mcp-cli-toolkit mcp hybrid_deployment

# Start testing server
mcp-cli-toolkit mcp hybrid_testing

# Pass arguments to MCP server
mcp-cli-toolkit mcp hybrid_deployment --port 8080 --debug
```

### CLI Script Operations

```bash
# Run deployment script
mcp-cli-toolkit cli deploy_best_practices.sh validate

# Run with full path
mcp-cli-toolkit cli deploy.sh --help

# Run monitoring script
mcp-cli-toolkit cli deployment_monitoring.py --status
```

### Profile Management

```bash
# List profiles
mcp-config profile list

# Create new profile
mcp-config profile create production --description "Production environment"

# Show profile details
mcp-config profile show default

# Delete profile
mcp-config profile delete old-profile
```

### Environment Configuration

```bash
# Show current environment
mcp-config env show

# Set environment variable
mcp-config env set MCP_CLI_LOG_LEVEL DEBUG

# Set custom config path
mcp-config env set MCP_CLI_CONFIG_PATH /custom/path
```

## 🛠️ Available Tools

### MCP Servers

#### 🚀 Hybrid Deployment Server (`hybrid_deployment_server.py`)
Advanced deployment automation with cloud platform support.

**Features:**
- Multi-cloud deployment (AWS, GCP, Azure, Coolify)
- Docker image building
- SSH tunnel management
- Health check automation
- API access for various platforms

**Usage:**
```bash
mcp-cli-toolkit mcp hybrid_deployment
```

**Available Tools:**
- `deploy_to_cloud` - Deploy to cloud platforms
- `access_coolify_api` - Access Coolify self-hosted platform
- `access_web_service` - Generic web service access
- `build_docker_image` - Build Docker images
- `manage_ssh_tunnel` - SSH tunnel management
- `health_check_deployment` - Comprehensive health checks

#### 🧪 Hybrid Testing Server (`hybrid_testing_server.py`)
Comprehensive testing framework for all development stages.

**Features:**
- Multi-language unit testing
- Integration testing
- Performance testing
- Security testing
- API endpoint validation

**Usage:**
```bash
mcp-cli-toolkit mcp hybrid_testing
```

**Available Tools:**
- `run_unit_tests` - Execute unit tests
- `run_integration_tests` - Execute integration tests
- `run_performance_tests` - Execute performance tests
- `run_security_tests` - Execute security tests
- `validate_api_endpoints` - Validate API endpoints

### CLI Scripts

#### 📋 Deployment Best Practices (`deploy_best_practices.sh`)
Comprehensive deployment validation and automation.

**Features:**
- Environment validation
- Security checks
- Performance optimization
- Nginx configuration
- Docker validation
- Error prevention
- Monitoring integration

**Usage:**
```bash
mcp-cli-toolkit cli deploy_best_practices.sh validate
mcp-cli-toolkit cli deploy_best_practices.sh deploy
mcp-cli-toolkit cli deploy_best_practices.sh status
```

#### 🚢 Coolify Deployment (`deploy.sh`)
Specialized deployment script for Coolify platform.

**Features:**
- SSH tunnel management
- API connectivity testing
- Deployment monitoring
- Health verification
- Status reporting

**Usage:**
```bash
mcp-cli-toolkit cli deploy.sh
```

#### 📊 Monitoring Scripts
- `deployment_monitoring.py` - Real-time deployment monitoring
- `deployment_error_tracker.py` - Error tracking and analysis
- `deployment_error_prevention.py` - Proactive error prevention
- `pre_deployment_validator.py` - Pre-deployment validation

## ⚙️ Configuration

### Profile Structure

```json
{
  "name": "production",
  "description": "Production environment profile",
  "version": "1.0.0",
  "created": "2024-01-01T00:00:00Z",
  "mcp_servers": {
    "hybrid_deployment": {
      "enabled": true,
      "path": "mcp-servers/hybrid_deployment_server.py",
      "description": "Production deployment server"
    }
  },
  "cli_scripts": {
    "deploy_best_practices": {
      "enabled": true,
      "path": "cli-scripts/deploy_best_practices.sh",
      "description": "Production deployment script"
    }
  },
  "settings": {
    "log_level": "INFO",
    "parallel_execution": true,
    "timeout": 300,
    "retry_attempts": 3
  }
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_CLI_TOOLKIT_HOME` | Toolkit root directory | Auto-detected |
| `MCP_CLI_VENV_PATH` | Python virtual environment | `$TOOLKIT_ROOT/lib/venv` |
| `MCP_CLI_CONFIG_PATH` | Configuration directory | `$HOME/.config/mcp-cli-toolkit` |
| `MCP_CLI_LOG_LEVEL` | Logging level | `INFO` |
| `MCP_CLI_DEFAULT_PROFILE` | Default profile | `default` |

### Custom Configuration

```bash
# Create custom profile
mcp-config profile create myprofile --description "My custom setup"

# Configure environment
mcp-config env set MCP_CLI_LOG_LEVEL DEBUG
mcp-config env set MCP_CLI_DEFAULT_PROFILE myprofile
```

## 💡 Examples

### Example 1: Development Workflow

```bash
# 1. Activate toolkit
source mcp-cli-toolkit/activate.sh

# 2. Check status
mcp-cli-toolkit status

# 3. Run tests
mcp-cli-toolkit mcp hybrid_testing

# 4. Deploy with validation
mcp-cli-toolkit cli deploy_best_practices.sh validate
mcp-cli-toolkit cli deploy_best_practices.sh deploy
```

### Example 2: Production Deployment

```bash
# 1. Load production profile
mcp-cli-toolkit profile production

# 2. Validate environment
mcp-cli-toolkit cli deploy_best_practices.sh validate

# 3. Deploy to Coolify
mcp-cli-toolkit cli deploy.sh

# 4. Monitor deployment
mcp-cli-toolkit cli deployment_monitoring.py --watch
```

### Example 3: Custom MCP Integration

```bash
# 1. Create custom profile
mcp-config profile create custom --description "Custom MCP setup"

# 2. Configure MCP server
mcp-config env set MCP_CLI_CUSTOM_SERVER "my_server.py"

# 3. Run custom server
mcp-cli-toolkit mcp my_server
```

## 🔧 Advanced Usage

### Creating Custom Profiles

```bash
# Create development profile
mcp-config profile create development \
  --description "Development environment with debug logging"

# Create production profile
mcp-config profile create production \
  --description "Production environment with monitoring"
```

### Environment Optimization

```bash
# Optimize for performance
mcp-config env set MCP_CLI_LOG_LEVEL WARN
mcp-config env set MCP_CLI_PARALLEL_EXECUTION true

# Optimize for debugging
mcp-config env set MCP_CLI_LOG_LEVEL DEBUG
mcp-config env set MCP_CLI_TIMEOUT 600
```

### Integration with AI Agents

```bash
# For any AI agent, simply:
# 1. Install the toolkit
# 2. Activate the environment
# 3. Use any of the available tools

# Example agent integration:
echo "Installing MCP-CLI Toolkit..."
./install.sh

echo "Activating toolkit..."
source activate.sh

echo "Available tools:"
mcp-cli-toolkit list
```

## 🐛 Troubleshooting

### Common Issues

**1. Command not found**
```bash
# Ensure toolkit is in PATH
export PATH="$HOME/.local/bin:$PATH"

# Or use full path
~/.local/bin/mcp-cli-toolkit status
```

**2. Python dependencies missing**
```bash
# Install Python dependencies
cd mcp-cli-toolkit
python3 -m venv lib/venv
source lib/venv/bin/activate
pip install -r requirements.txt
```

**3. Profile not found**
```bash
# List available profiles
mcp-config profile list

# Create new profile
mcp-config profile create default
```

**4. MCP server connection failed**
```bash
# Check if server is running
ps aux | grep python

# Check server logs
tail -f /tmp/mcp-server.log
```

### Getting Help

```bash
# General help
mcp-cli-toolkit help

# Configuration help
mcp-config --help

# Profile management help
mcp-config profile --help

# Environment help
mcp-config env --help
```

### Debug Mode

```bash
# Enable debug logging
mcp-config env set MCP_CLI_LOG_LEVEL DEBUG

# Run with verbose output
mcp-cli-toolkit --verbose status
```

## 🤝 Contributing

### Adding New Tools

1. **MCP Servers**: Add to `mcp-servers/` directory
2. **CLI Scripts**: Add to `cli-scripts/` directory
3. **Documentation**: Update this README
4. **Configuration**: Update default profile

### Development Setup

```bash
# Clone repository
git clone https://github.com/your-repo/mcp-cli-toolkit.git
cd mcp-cli-toolkit

# Install in development mode
./install.sh --dev

# Run tests
python3 -m pytest tests/

# Add new tool
cp template.py mcp-servers/my_new_server.py
```

### Testing

```bash
# Run all tests
make test

# Run specific test
python3 -m pytest tests/test_deployment.py

# Integration tests
python3 -m pytest tests/integration/
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Model Context Protocol** - For the MCP server framework
- **AI Agent Community** - For inspiration and feedback
- **Open Source Contributors** - For tools and utilities

## 📞 Support

- 📧 **Email**: support@mcp-cli-toolkit.com
- 💬 **Discord**: [Join our community](https://discord.gg/mcp-cli-toolkit)
- 🐛 **Issues**: [GitHub Issues](https://github.com/your-repo/mcp-cli-toolkit/issues)
- 📚 **Docs**: [Full Documentation](https://docs.mcp-cli-toolkit.com)

---

<div align="center">

**Made with ❤️ for the AI Agent community**

[⭐ Star this repo](https://github.com/your-repo/mcp-cli-toolkit) | [🐛 Report issues](https://github.com/your-repo/mcp-cli-toolkit/issues) | [📖 Documentation](https://docs.mcp-cli-toolkit.com)

</div>