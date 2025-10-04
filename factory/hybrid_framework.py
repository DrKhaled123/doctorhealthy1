#!/usr/bin/env python3
"""
Hybrid Framework for Factory Orchestrator
Combines GitHub, SWE-Agent, and ChatGDB Light integrations
"""

import os
import json
import time
import asyncio
from typing import Dict, List, Any, Optional
from datetime import datetime

class HybridFramework:
    """Hybrid framework combining all advanced integrations"""
    
    def __init__(self):
        self.error_count = 0
        self.components = {
            'github': False,
            'swe_agent': False,
            'debugger': True,
            'memory': True,
            'orchestrator': True
        }
    
    async def run_hybrid_analysis(self, specification: str) -> Dict:
        """Run hybrid analysis with all components"""
        print(f"🔄 Running Hybrid Analysis: {specification}")
        
        results = {
            'specification': specification,
            'timestamp': datetime.now().isoformat(),
            'components_used': [],
            'success': True,
            'error_count': 0
        }
        
        # Memory analysis
        try:
            from memo_ai_memory import MemoAIMemory
            memory = MemoAIMemory()
            similar = await memory.recall(specification, k=3)
            results['memory_analysis'] = {
                'similar_experiences': len(similar),
                'success': True
            }
            results['components_used'].append('memory')
        except Exception as e:
            results['memory_analysis'] = {'error': str(e)}
            results['error_count'] += 1
        
        # GitHub integration
        try:
            from github_integration import GitHubIntegration
            github = GitHubIntegration()
            if github.token and github.repository:
                results['github_ready'] = True
                results['components_used'].append('github')
            else:
                results['github_ready'] = False
        except Exception as e:
            results['github_analysis'] = {'error': str(e)}
            results['error_count'] += 1
        
        # SWE-Agent integration
        try:
            from swe_agent_integration import SWEAgentIntegration
            swe_agent = SWEAgentIntegration()
            if os.path.exists(swe_agent.swe_agent_path):
                results['swe_agent_ready'] = True
                results['components_used'].append('swe_agent')
            else:
                results['swe_agent_ready'] = False
        except Exception as e:
            results['swe_agent_analysis'] = {'error': str(e)}
            results['error_count'] += 1
        
        # Debugging integration
        try:
            from chatgdb_light_integration import ChatGDBLight
            debugger = ChatGDBLight()
            results['debugger_ready'] = True
            results['components_used'].append('debugger')
        except Exception as e:
            results['debugger_analysis'] = {'error': str(e)}
            results['error_count'] += 1
        
        results['total_components'] = len(results['components_used'])
        results['success_rate'] = (results['total_components'] / 5) * 100
        
        print(f"✅ Hybrid analysis completed: {results['total_components']}/5 components ready")
        return results

# Example usage
async def demo_hybrid_framework():
    """Demo the hybrid framework"""
    framework = HybridFramework()
    results = await framework.run_hybrid_analysis("Implement user authentication")
    
    print("📊 Hybrid Framework Results:")
    print(f"  Components ready: {results['total_components']}/5")
    print(f"  Success rate: {results['success_rate']:.1f}%")
    print(f"  Components used: {', '.join(results['components_used'])}")
    
    return results

if __name__ == "__main__":
    results = asyncio.run(demo_hybrid_framework())
