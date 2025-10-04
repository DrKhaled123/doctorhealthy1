#!/bin/bash

# Factory Orchestrator Startup Script
# This script starts the complete factory system with all components

echo "🚀 Starting Factory Orchestrator System..."
echo "============================================="

# Check if Redis is running
echo "📊 Checking Redis status..."
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis is not running. Starting Redis..."
    brew services start redis
    sleep 2
fi

# Verify Redis is accessible
if redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is running and accessible"
else
    echo "❌ Failed to start Redis. Please start it manually: brew services start redis"
    exit 1
fi

# Set up Python environment
echo "🐍 Setting up Python environment..."
cd "$(dirname "$0")"

# Check if virtual environment exists
if [ ! -d "venv" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source venv/bin/activate

# Install/update dependencies
echo "📚 Installing Python dependencies..."
pip install -r requirements.txt

# Create necessary directories
echo "📁 Creating output directories..."
mkdir -p logs
mkdir -p reports
mkdir -p screenshots

# Set environment variables
export REDIS_HOST=localhost
export REDIS_PORT=6379
export FACTORY_ENV=production

echo ""
echo "============================================="
echo "🎯 Factory Components Ready!"
echo "============================================="
echo "Components available:"
echo "  • Factory Orchestrator (factory_orchestrator.py)"
echo "  • Browser Testing Agent (browser_testing_agent.py)"
echo "  • Autofix Agent (autofix_agent.py)"
echo "  • Memory System (memo_ai_memory.py)"
echo "  • Light Agent Orchestrator (orchestrator.py)"
echo ""
echo "📖 Usage Examples:"
echo "  python3 factory_orchestrator.py 'Implement user login system'"
echo "  python3 orchestrator.py"
echo "  python3 browser_testing_agent.py"
echo ""
echo "🌐 Redis is running on localhost:6379"
echo "📊 Monitor with: redis-cli monitor"
echo "============================================="

# Keep the script running to show status
echo "Press Ctrl+C to stop all services"
echo ""

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down Factory system..."
    echo "✅ Thank you for using Factory Orchestrator!"
    exit 0
}

# Set trap for cleanup
trap cleanup SIGINT SIGTERM

# Show system status every 30 seconds
while true; do
    echo "📊 System Status (Redis: $(redis-cli ping 2>/dev/null || echo 'disconnected'))"
    sleep 30
done
