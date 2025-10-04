# Factory Orchestrator System

A comprehensive automated development pipeline system with multiple specialized agents for coding, testing, fixing, and deployment.

## 🚀 Features

- **Multi-Agent Architecture**: Specialized agents for different development stages
- **Memory System**: Vector-based memory for learning from experiences
- **Browser Testing**: Automated web application testing with Playwright and Selenium
- **Code Autofix**: Automatic code fixing and formatting
- **Redis Integration**: High-performance caching and state management
- **Pipeline Orchestration**: Complete development workflow automation

## 📁 Project Structure

```
factory/
├── factory_orchestrator.py     # Main pipeline orchestrator
├── browser_testing_agent.py    # Web application testing
├── autofix_agent.py           # Code fixing and formatting
├── memo_ai_memory.py          # Vector memory system
├── orchestrator.py            # LightAgent integration
├── demo_complete_system_fixed.py # System demo
├── start_factory.sh           # Startup script
├── requirements.txt           # Python dependencies
└── factory_config.json        # Configuration
```

## 🛠️ Installation & Setup

### Prerequisites
- Python 3.13+
- Redis Server
- Homebrew (macOS)

### Quick Start

1. **Install Dependencies**
   ```bash
   cd factory
   pip install -r requirements.txt
   ```

2. **Start Redis**
   ```bash
   brew services start redis
   ```

3. **Run the Factory**
   ```bash
   # Option 1: Use the startup script
   ./start_factory.sh

   # Option 2: Run directly
   python3 factory_orchestrator.py "Implement user login system"
   ```

## 🎯 Components

### 1. Factory Orchestrator (`factory_orchestrator.py`)
Main pipeline orchestrator that manages the complete development workflow:

- **Stages**: Coding → Review → Autofix → Testing → Error Solving → Validation → Deployment
- **Redis Integration**: State management and caching
- **Memory Integration**: Learning from past experiences
- **Error Handling**: Robust error recovery and retry mechanisms

### 2. Browser Testing Agent (`browser_testing_agent.py`)
Comprehensive web application testing:

- **Playwright**: Modern async testing (primary)
- **Selenium**: Legacy browser support (fallback)
- **Test Generation**: Automatic test case generation
- **Screenshot Capture**: Visual regression testing
- **Multi-Browser**: Cross-browser compatibility

### 3. Autofix Agent (`autofix_agent.py`)
Automatic code fixing and formatting:

- **Syntax Error Fixes**: Common syntax issues
- **Import Error Resolution**: Missing dependencies
- **Code Formatting**: PEP8 and Black formatting
- **Security Fixes**: Common security vulnerabilities
- **Type Hints**: Automatic type annotation

### 4. Memory System (`memo_ai_memory.py`)
Vector-based memory for AI learning:

- **Sentence Transformers**: Advanced text embeddings
- **Similarity Search**: Find relevant past experiences
- **Persistent Storage**: JSON-based memory storage
- **Redis Integration**: Fast memory retrieval
- **Pattern Learning**: Success/failure pattern recognition

### 5. Light Agent Orchestrator (`orchestrator.py`)
Lightweight multi-agent coordination:

- **Agent Management**: Register and manage multiple agents
- **Task Distribution**: Intelligent task routing
- **Backend Integration**: Redis-based state sharing
- **Workflow Orchestration**: Complex workflow management

## 📖 Usage Examples

### Basic Factory Run
```python
from factory_orchestrator import FactoryOrchestrator
import asyncio

async def main():
    orchestrator = FactoryOrchestrator()
    await orchestrator.initialize()
    result = await orchestrator.run_pipeline("Create a REST API")
    print(f"Result: {result.success}")

asyncio.run(main())
```

### Memory System Usage
```python
from memo_ai_memory import MemoAIMemory

memory = MemoAIMemory()
await memory.store_success_pattern("api_design", {"method": "REST"}, "success")
similar = await memory.recall("web service", k=3)
```

### Browser Testing
```python
from browser_testing_agent import BrowserTestingAgent

agent = BrowserTestingAgent()
test_cases = await agent.generate_test_cases("web application")
results = await agent.test_web_application("http://localhost:3000", test_cases)
```

### Code Autofix
```python
from autofix_agent import AutofixAgent

agent = AutofixAgent()
result = await agent.fix_code(problematic_code)
```

## 🔧 Configuration

### Factory Configuration (`factory_config.json`)
```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "db": 0
  },
  "memory": {
    "max_memories": 1000,
    "similarity_threshold": 0.8
  },
  "debug": false
}
```

### Environment Variables
- `REDIS_HOST`: Redis server hostname
- `REDIS_PORT`: Redis server port
- `FACTORY_ENV`: Environment (development/production)

## 🧪 Testing

Run the complete system demo:
```bash
python3 demo_complete_system_fixed.py
```

Expected output:
```
🚀 Factory Orchestrator Complete System Demo
=============================================
✅ Memory System: PASSED
✅ Autofix Agent: PASSED
✅ Browser Agent: PASSED
✅ Factory Orchestrator: PASSED

🎉 All tests passed! Factory system is ready.
```

## 🔍 Monitoring

### Redis Monitoring
```bash
redis-cli monitor
```

### System Status
The startup script provides real-time system status every 30 seconds.

### Logs
- Application logs: `logs/`
- Screenshots: `screenshots/`
- Reports: `reports/`

## 🚀 Deployment

### Production Deployment
1. Configure production Redis
2. Set environment variables
3. Use production startup script
4. Enable monitoring and alerting

### Docker Deployment
```dockerfile
FROM python:3.13-slim
COPY factory/ /app/
WORKDIR /app
RUN pip install -r requirements.txt
CMD ["python3", "factory_orchestrator.py"]
```

## 📊 Performance

- **Memory Efficient**: Optimized memory usage with Redis caching
- **Fast Testing**: Parallel test execution with Playwright
- **Scalable**: Multi-agent architecture supports horizontal scaling
- **Reliable**: Comprehensive error handling and recovery

## 🔒 Security

- **Input Validation**: All inputs are validated and sanitized
- **Dependency Security**: Regular security updates
- **Access Control**: Configurable access permissions
- **Audit Logging**: Complete audit trail

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 🆘 Support

For issues and questions:
1. Check the demo output for common issues
2. Review the configuration files
3. Check Redis connectivity
4. Verify Python dependencies

---

**Built with ❤️ for autonomous software development**
