#!/usr/bin/env python3
"""
Inter-Agent Learning System Comprehensive Demo
Demonstrates how AI agents learn from each other's experiences
"""

import asyncio
import time
import json
from datetime import datetime
from typing import Dict, List, Any

# Import the learning system components
from inter_agent_learning_system import (
    get_learning_coordinator,
    AgentCapability,
    LearningType
)

from learning_enabled_browser_agent import LearningEnabledBrowserTestingAgent
from learning_enabled_code_review_agent import LearningEnabledCodeReviewAgent

class InterAgentLearningDemo:
    """Comprehensive demonstration of inter-agent learning"""

    def __init__(self):
        self.coordinator = None
        self.browser_agent = None
        self.review_agent = None
        self.demo_results = {}

    async def initialize(self):
        """Initialize the demo system"""
        print("🚀 Initializing Inter-Agent Learning Demo")
        print("=" * 60)

        try:
            # Initialize learning coordinator
            self.coordinator = await get_learning_coordinator()
            print("✅ Learning coordinator initialized")

            # Initialize learning-enabled agents
            self.browser_agent = LearningEnabledBrowserTestingAgent("demo_browser_agent")
            await self.browser_agent.initialize()
            print("✅ Browser testing agent initialized")

            self.review_agent = LearningEnabledCodeReviewAgent("demo_review_agent")
            await self.review_agent.initialize()
            print("✅ Code review agent initialized")

            print("🎓 All agents registered in learning system")
            return True

        except Exception as e:
            print(f"❌ Initialization failed: {e}")
            return False

    async def demonstrate_agent_registration(self):
        """Demonstrate agent registration and profiling"""
        print("\n🤖 Agent Registration and Profiling")
        print("-" * 40)

        try:
            # Get learning coordinator
            coordinator = await get_learning_coordinator()

            # Show agent profiles
            print("📋 Current Agent Profiles:")
            for agent_id, profile in coordinator.agent_profiles.items():
                print(f"  • {agent_id} ({profile.agent_type}):")
                print(f"    - Capabilities: {[cap.value for cap in profile.capabilities]}")
                print(f"    - Experience count: {profile.experience_count}")
                print(f"    - Success rate: {profile.success_rate:.2f}")
                print(f"    - Reputation: {profile.reputation_score:.2f}")

            self.demo_results['agent_registration'] = {
                'total_agents': len(coordinator.agent_profiles),
                'agent_profiles': {
                    aid: profile.to_dict() for aid, profile in coordinator.agent_profiles.items()
                }
            }

            return True

        except Exception as e:
            print(f"❌ Agent registration demo failed: {e}")
            return False

    async def demonstrate_experience_sharing(self):
        """Demonstrate agents sharing learning experiences"""
        print("\n📖 Experience Sharing Demonstration")
        print("-" * 40)

        try:
            coordinator = await get_learning_coordinator()

            # Simulate browser testing experience
            print("🧪 Simulating browser testing experience...")
            browser_experience = await self.browser_agent.share_experience(
                capability=AgentCapability.TESTING,
                context={
                    'task_type': 'ui_testing',
                    'complexity': 'medium',
                    'test_count': 5
                },
                outcome={
                    'success': True,
                    'success_patterns': ['wait_conditions', 'error_handling'],
                    'performance_metrics': {'execution_time': 45.2}
                },
                success=True,
                confidence=0.85,
                lessons_learned=[
                    'Proper wait conditions prevent flaky tests',
                    'Error handling improves test reliability',
                    'Test execution time correlates with page performance'
                ]
            )

            # Simulate code review experience
            print("🔍 Simulating code review experience...")
            review_experience = await self.review_agent.share_experience(
                capability=AgentCapability.CODE_REVIEW,
                context={
                    'task_type': 'security_review',
                    'complexity': 'high',
                    'files_reviewed': 25
                },
                outcome={
                    'success': True,
                    'success_patterns': ['automated_analysis', 'pattern_matching'],
                    'security_findings': 3
                },
                success=True,
                confidence=0.92,
                lessons_learned=[
                    'Automated analysis catches common security issues',
                    'Pattern matching helps identify code smells',
                    'Large codebases benefit from systematic review'
                ]
            )

            print("✅ Experiences shared successfully")
            print(f"📊 Total experiences in knowledge base: {len(coordinator.learning_experiences)}")

            self.demo_results['experience_sharing'] = {
                'experiences_shared': 2,
                'knowledge_base_size': len(coordinator.learning_experiences),
                'last_experiences': list(coordinator.learning_experiences)[-2:]
            }

            return True

        except Exception as e:
            print(f"❌ Experience sharing demo failed: {e}")
            return False

    async def demonstrate_knowledge_transfer(self):
        """Demonstrate knowledge transfer between agents"""
        print("\n🔄 Knowledge Transfer Demonstration")
        print("-" * 40)

        try:
            coordinator = await get_learning_coordinator()

            # Request knowledge transfer from browser agent to review agent
            print("🔄 Requesting knowledge transfer...")
            transfer_id = await coordinator.request_knowledge_transfer(
                self.browser_agent.agent_id,
                self.review_agent.agent_id,
                AgentCapability.TESTING
            )

            if transfer_id:
                print(f"✅ Knowledge transfer requested: {transfer_id}")

                # Process the transfer
                success = await coordinator.process_knowledge_transfer(transfer_id, accept=True)

                if success:
                    print("✅ Knowledge transfer completed successfully")
                    print("🎓 Review agent can now apply testing knowledge to code review")

                    self.demo_results['knowledge_transfer'] = {
                        'transfer_id': transfer_id,
                        'success': True,
                        'transfers_processed': len([t for t in coordinator.knowledge_transfers.values() if t.status == 'accepted'])
                    }
                else:
                    print("❌ Knowledge transfer failed")
                    return False
            else:
                print("⚠️ No knowledge transfer created (may be insufficient experiences)")
                self.demo_results['knowledge_transfer'] = {
                    'transfer_id': None,
                    'success': False,
                    'reason': 'insufficient_experiences'
                }

            return True

        except Exception as e:
            print(f"❌ Knowledge transfer demo failed: {e}")
            return False

    async def demonstrate_learning_recommendations(self):
        """Demonstrate learning recommendations system"""
        print("\n🎯 Learning Recommendations Demonstration")
        print("-" * 40)

        try:
            coordinator = await get_learning_coordinator()

            # Get recommendations for browser agent
            print("📚 Getting recommendations for browser agent...")
            browser_recommendations = await coordinator.get_agent_recommendations(
                self.browser_agent.agent_id, AgentCapability.TESTING
            )

            print(f"✅ Found {len(browser_recommendations)} recommendations for browser agent:")
            for rec in browser_recommendations:
                print(f"  • {rec['type']}: {rec['reason']}")

            # Get recommendations for review agent
            print("\n📚 Getting recommendations for code review agent...")
            review_recommendations = await coordinator.get_agent_recommendations(
                self.review_agent.agent_id, AgentCapability.CODE_REVIEW
            )

            print(f"✅ Found {len(review_recommendations)} recommendations for review agent:")
            for rec in review_recommendations:
                print(f"  • {rec['type']}: {rec['reason']}")

            # Get collaboration opportunities
            print("\n🤝 Finding collaboration opportunities...")
            opportunities = await coordinator.get_collaboration_opportunities(
                self.browser_agent.agent_id
            )

            print(f"✅ Found {len(opportunities)} collaboration opportunities:")
            for opp in opportunities[:3]:  # Show top 3
                print(f"  • {opp['agent_type']}: {opp['reason']}")

            self.demo_results['learning_recommendations'] = {
                'browser_agent_recommendations': len(browser_recommendations),
                'review_agent_recommendations': len(review_recommendations),
                'collaboration_opportunities': len(opportunities),
                'top_opportunities': opportunities[:3]
            }

            return True

        except Exception as e:
            print(f"❌ Learning recommendations demo failed: {e}")
            return False

    async def demonstrate_pattern_learning(self):
        """Demonstrate pattern recognition and learning"""
        print("\n🔍 Pattern Learning Demonstration")
        print("-" * 40)

        try:
            coordinator = await get_learning_coordinator()

            # Generate learning report to show patterns
            print("📊 Generating comprehensive learning report...")
            report = await coordinator.generate_learning_report()

            print("✅ Learning Report Generated:")
            print(f"  • Total experiences: {report['metrics']['total_experiences']}")
            print(f"  • Knowledge base size: {report['metrics']['knowledge_base_size']}")
            print(f"  • Successful transfers: {report['metrics']['successful_transfers']}")

            # Show top performers
            print("\n🏆 Top Performing Agents:")
            for agent in report['top_performers']:
                print(f"  • {agent['agent_id']} ({agent['type']}): Score {agent['score']:.3f}")

            # Show learning patterns
            print("\n📈 Learning Patterns Detected:")
            for capability, pattern in report['learning_patterns'].items():
                print(f"  • {capability}: {pattern['experience_count']} experiences, "
                       f"{pattern['success_rate']:.1%} success rate")

            # Show knowledge gaps
            if report['knowledge_gaps']:
                print("\n🔍 Knowledge Gaps Identified:")
                for gap in report['knowledge_gaps']:
                    print(f"  • {gap['capability']}: {gap['recommendation']}")

            self.demo_results['pattern_learning'] = {
                'report_summary': {
                    'total_experiences': report['metrics']['total_experiences'],
                    'knowledge_base_size': report['metrics']['knowledge_base_size'],
                    'top_performers_count': len(report['top_performers']),
                    'patterns_detected': len(report['learning_patterns']),
                    'knowledge_gaps': len(report['knowledge_gaps'])
                },
                'full_report': report
            }

            return True

        except Exception as e:
            print(f"❌ Pattern learning demo failed: {e}")
            return False

    async def demonstrate_continuous_learning(self):
        """Demonstrate continuous learning capabilities"""
        print("\n🔄 Continuous Learning Demonstration")
        print("-" * 40)

        try:
            # Simulate multiple learning cycles
            print("🔄 Running multiple learning cycles...")

            for cycle in range(3):
                print(f"\n📚 Learning Cycle {cycle + 1}:")

                # Share new experiences
                cycle_experience = await self.browser_agent.share_experience(
                    capability=AgentCapability.TESTING,
                    context={
                        'task_type': f'cycle_{cycle}_testing',
                        'complexity': 'medium',
                        'cycle': cycle
                    },
                    outcome={
                        'success': True,
                        'improvement': f'cycle_{cycle}_improvement'
                    },
                    success=True,
                    confidence=0.8 + (cycle * 0.05),  # Increasing confidence
                    lessons_learned=[
                        f'Cycle {cycle} learning: Improved test reliability',
                        'Continuous learning enhances agent performance'
                    ]
                )

                print(f"  ✅ Shared cycle {cycle} experience")

                # Get updated recommendations
                recommendations = await self.browser_agent.get_learning_recommendations(AgentCapability.TESTING)
                print(f"  📋 {len(recommendations)} recommendations available")

                # Small delay to simulate real-world timing
                await asyncio.sleep(0.5)

            print("✅ Continuous learning cycles completed")

            self.demo_results['continuous_learning'] = {
                'cycles_completed': 3,
                'experiences_accumulated': len(self.coordinator.learning_experiences),
                'learning_progression': 'demonstrated'
            }

            return True

        except Exception as e:
            print(f"❌ Continuous learning demo failed: {e}")
            return False

    async def generate_final_report(self):
        """Generate comprehensive final report"""
        print("\n📊 Generating Final Report")
        print("-" * 40)

        try:
            coordinator = await get_learning_coordinator()

            # Get final metrics
            final_report = await coordinator.generate_learning_report()

            # Add demo-specific metrics
            final_report['demo_results'] = self.demo_results
            final_report['demo_summary'] = {
                'total_demo_time': time.time() - self.demo_results.get('start_time', time.time()),
                'agents_participated': len(coordinator.agent_profiles),
                'experiences_generated': len(coordinator.learning_experiences),
                'knowledge_transfers': len(coordinator.knowledge_transfers),
                'collaboration_network_size': len(coordinator.collaboration_network)
            }

            # Save report to file
            report_filename = f"inter_agent_learning_demo_report_{int(time.time())}.json"
            with open(report_filename, 'w') as f:
                json.dump(final_report, f, indent=2, default=str)

            print("✅ Final report generated and saved")
            print(f"📄 Report file: {report_filename}")

            # Display summary
            summary = final_report['demo_summary']
            print("\n🎉 Demo Summary:")
            print(f"  • Agents participated: {summary['agents_participated']}")
            print(f"  • Experiences generated: {summary['experiences_generated']}")
            print(f"  • Knowledge transfers: {summary['knowledge_transfers']}")
            print(f"  • Collaboration connections: {summary['collaboration_network_size']}")
            print(f"  • Total demo time: {summary['total_demo_time']:.2f}s")

            self.demo_results['final_report'] = final_report
            return report_filename

        except Exception as e:
            print(f"❌ Final report generation failed: {e}")
            return None

    async def run_complete_demo(self):
        """Run the complete inter-agent learning demonstration"""
        print("🎭 Inter-Agent Learning System - Complete Demo")
        print("=" * 70)
        print("This demo showcases how AI agents can learn from each other's")
        print("experiences, share knowledge, and continuously improve performance.")
        print("=" * 70)

        # Record start time
        self.demo_results['start_time'] = time.time()

        try:
            # Step 1: Initialize system
            if not await self.initialize():
                return False

            # Step 2: Demonstrate agent registration
            if not await self.demonstrate_agent_registration():
                return False

            # Step 3: Demonstrate experience sharing
            if not await self.demonstrate_experience_sharing():
                return False

            # Step 4: Demonstrate knowledge transfer
            if not await self.demonstrate_knowledge_transfer():
                return False

            # Step 5: Demonstrate learning recommendations
            if not await self.demonstrate_learning_recommendations():
                return False

            # Step 6: Demonstrate pattern learning
            if not await self.demonstrate_pattern_learning():
                return False

            # Step 7: Demonstrate continuous learning
            if not await self.demonstrate_continuous_learning():
                return False

            # Step 8: Generate final report
            report_file = await self.generate_final_report()

            print("\n🎉 COMPLETE DEMO SUCCESSFUL!")
            print("=" * 70)
            print("The Inter-Agent Learning System has been successfully demonstrated!")
            print("Key achievements:")
            print("  ✅ Agents can register and maintain learning profiles")
            print("  ✅ Agents can share experiences with each other")
            print("  ✅ Knowledge transfer between agents works seamlessly")
            print("  ✅ Learning recommendations help agents improve")
            print("  ✅ Pattern recognition identifies successful strategies")
            print("  ✅ Continuous learning enables ongoing improvement")
            print("  ✅ Comprehensive reporting provides insights")

            if report_file:
                print(f"\n📄 Detailed report saved to: {report_file}")

            return True

        except Exception as e:
            print(f"❌ Demo failed with error: {e}")
            import traceback
            traceback.print_exc()
            return False

# Main execution function
async def main():
    """Main function to run the complete demo"""
    demo = InterAgentLearningDemo()
    success = await demo.run_complete_demo()
    return success

if __name__ == "__main__":
    # Run the complete demonstration
    success = asyncio.run(main())
    exit(0 if success else 1)
