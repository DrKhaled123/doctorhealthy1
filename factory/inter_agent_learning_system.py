#!/usr/bin/env python3
"""
Inter-Agent Learning System
A comprehensive system where AI agents learn from each other's experiences,
share knowledge, and continuously improve their performance.
"""

import asyncio
import json
import time
import uuid
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
import redis.asyncio as redis
import numpy as np
from collections import defaultdict, deque
import logging
import os
import sys

# Add current directory to path for imports
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)

from factory_config import (
    get_factory,
    learn_from_deployment,
    suggest_pattern_solution,
    PatternLearningSystem,
    ContinuousImprovementEngine
)

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class LearningType(Enum):
    """Types of learning interactions"""
    EXPERIENCE_SHARING = "experience_sharing"
    PATTERN_RECOGNITION = "pattern_recognition"
    COLLABORATIVE_PROBLEM_SOLVING = "collaborative_problem_solving"
    KNOWLEDGE_TRANSFER = "knowledge_transfer"
    PERFORMANCE_IMPROVEMENT = "performance_improvement"
    ERROR_PREVENTION = "error_prevention"

class AgentCapability(Enum):
    """Agent capabilities that can be learned and shared"""
    CODE_GENERATION = "code_generation"
    CODE_REVIEW = "code_review"
    TESTING = "testing"
    DEBUGGING = "debugging"
    DEPLOYMENT = "deployment"
    MONITORING = "monitoring"
    SECURITY_SCANNING = "security_scanning"
    PERFORMANCE_OPTIMIZATION = "performance_optimization"

@dataclass
class LearningExperience:
    """Represents a learning experience that can be shared between agents"""
    experience_id: str
    agent_id: str
    agent_type: str
    learning_type: LearningType
    capability: AgentCapability
    context: Dict[str, Any]
    outcome: Dict[str, Any]
    success: bool
    confidence: float
    lessons_learned: List[str]
    timestamp: datetime
    metadata: Dict[str, Any]

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for storage"""
        return {
            **asdict(self),
            'learning_type': self.learning_type.value,
            'capability': self.capability.value,
            'timestamp': self.timestamp.isoformat()
        }

@dataclass
class AgentProfile:
    """Profile of an agent's capabilities and learning history"""
    agent_id: str
    agent_type: str
    capabilities: List[AgentCapability]
    experience_count: int
    success_rate: float
    average_confidence: float
    specializations: List[str]
    collaboration_history: List[str]
    last_active: datetime
    reputation_score: float

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for storage"""
        return {
            **asdict(self),
            'capabilities': [cap.value for cap in self.capabilities],
            'last_active': self.last_active.isoformat()
        }

@dataclass
class KnowledgeTransfer:
    """Represents knowledge being transferred between agents"""
    transfer_id: str
    from_agent: str
    to_agent: str
    knowledge_type: str
    content: Dict[str, Any]
    confidence: float
    timestamp: datetime
    status: str  # 'pending', 'accepted', 'rejected', 'applied'

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for storage"""
        return {
            **asdict(self),
            'timestamp': self.timestamp.isoformat()
        }

class InterAgentLearningCoordinator:
    """Main coordinator for inter-agent learning system"""

    def __init__(self, redis_config: Dict[str, Any] = None):
        self.redis_config = redis_config or {
            "host": "localhost",
            "port": 6379,
            "db": 0,
            "decode_responses": True
        }
        self.redis_client = None
        self.agent_profiles = {}
        self.learning_experiences = deque(maxlen=10000)  # Keep last 10k experiences
        self.knowledge_transfers = {}
        self.active_agents = set()
        self.learning_patterns = {}
        self.collaboration_network = defaultdict(list)

        # Learning metrics
        self.learning_metrics = {
            'total_experiences': 0,
            'successful_transfers': 0,
            'failed_transfers': 0,
            'collaboration_count': 0,
            'knowledge_base_size': 0
        }

    async def initialize(self):
        """Initialize the learning coordinator"""
        try:
            self.redis_client = redis.Redis(**self.redis_config)
            await self.redis_client.ping()
            logger.info("✅ Inter-Agent Learning Coordinator initialized")

            # Load existing data from Redis
            await self._load_existing_data()

        except Exception as e:
            logger.error(f"❌ Failed to initialize learning coordinator: {e}")
            raise

    async def _load_existing_data(self):
        """Load existing learning data from Redis"""
        try:
            # Load agent profiles
            profiles_data = await self.redis_client.get("agent_profiles")
            if profiles_data:
                self.agent_profiles = json.loads(profiles_data)

            # Load learning experiences
            experiences_data = await self.redis_client.get("learning_experiences")
            if experiences_data:
                experiences_list = json.loads(experiences_data)
                self.learning_experiences.extend(experiences_list)

            # Load knowledge transfers
            transfers_data = await self.redis_client.get("knowledge_transfers")
            if transfers_data:
                self.knowledge_transfers = json.loads(transfers_data)

            logger.info(f"📚 Loaded {len(self.agent_profiles)} agent profiles, "
                       f"{len(self.learning_experiences)} experiences")

        except Exception as e:
            logger.warning(f"⚠️ Could not load existing data: {e}")

    async def register_agent(self, agent_id: str, agent_type: str,
                           capabilities: List[AgentCapability]) -> AgentProfile:
        """Register a new agent in the learning system"""
        profile = AgentProfile(
            agent_id=agent_id,
            agent_type=agent_type,
            capabilities=capabilities,
            experience_count=0,
            success_rate=0.0,
            average_confidence=0.0,
            specializations=[],
            collaboration_history=[],
            last_active=datetime.now(),
            reputation_score=0.5  # Start with neutral reputation
        )

        self.agent_profiles[agent_id] = profile
        self.active_agents.add(agent_id)

        # Save to Redis
        await self._save_agent_profiles()

        logger.info(f"🤖 Registered agent {agent_id} ({agent_type}) with capabilities: {[c.value for c in capabilities]}")
        return profile

    async def share_learning_experience(self, experience: LearningExperience) -> bool:
        """Share a learning experience with other agents"""
        try:
            # Add to experiences deque
            self.learning_experiences.append(experience.to_dict())

            # Update agent profile
            await self._update_agent_profile(experience.agent_id, experience)

            # Find relevant agents to share with
            relevant_agents = await self._find_relevant_agents(experience)

            # Create knowledge transfers for relevant agents
            for agent_id in relevant_agents:
                if agent_id != experience.agent_id:  # Don't share with self
                    transfer = KnowledgeTransfer(
                        transfer_id=str(uuid.uuid4()),
                        from_agent=experience.agent_id,
                        to_agent=agent_id,
                        knowledge_type=experience.capability.value,
                        content=self._extract_knowledge_content(experience),
                        confidence=experience.confidence,
                        timestamp=datetime.now(),
                        status='pending'
                    )

                    self.knowledge_transfers[transfer.transfer_id] = transfer
                    self.collaboration_network[experience.agent_id].append(agent_id)

            # Update metrics
            self.learning_metrics['total_experiences'] += 1
            self.learning_metrics['knowledge_base_size'] = len(self.learning_experiences)

            # Save to Redis
            await self._save_learning_data()

            logger.info(f"📖 Shared experience {experience.experience_id} with {len(relevant_agents)} agents")
            return True

        except Exception as e:
            logger.error(f"❌ Failed to share learning experience: {e}")
            return False

    async def _update_agent_profile(self, agent_id: str, experience: LearningExperience):
        """Update agent profile based on new experience"""
        if agent_id not in self.agent_profiles:
            return

        profile = self.agent_profiles[agent_id]
        profile.experience_count += 1
        profile.last_active = datetime.now()

        # Update success rate (weighted average)
        if profile.experience_count == 1:
            profile.success_rate = 1.0 if experience.success else 0.0
            profile.average_confidence = experience.confidence
        else:
            # Weighted average with more recent experiences having higher weight
            alpha = 0.3  # Learning rate
            profile.success_rate = (1 - alpha) * profile.success_rate + alpha * (1.0 if experience.success else 0.0)
            profile.average_confidence = (1 - alpha) * profile.average_confidence + alpha * experience.confidence

        # Update reputation based on success and confidence
        reputation_boost = experience.confidence * (1.0 if experience.success else -0.5)
        profile.reputation_score = min(1.0, max(0.0, profile.reputation_score + reputation_boost * 0.1))

        # Update specializations
        if experience.success and experience.confidence > 0.8:
            capability_str = experience.capability.value
            if capability_str not in profile.specializations:
                profile.specializations.append(capability_str)

    async def _find_relevant_agents(self, experience: LearningExperience) -> List[str]:
        """Find agents that could benefit from this experience"""
        relevant_agents = []

        for agent_id, profile in self.agent_profiles.items():
            if agent_id == experience.agent_id:
                continue

            # Check capability relevance
            if experience.capability in profile.capabilities:
                relevant_agents.append(agent_id)
                continue

            # Check if agent has shown interest in this area
            agent_specializations = [spec.lower() for spec in profile.specializations]
            experience_context = str(experience.context).lower()

            if any(spec in experience_context for spec in agent_specializations):
                relevant_agents.append(agent_id)
                continue

            # Check collaboration history
            if experience.agent_id in self.collaboration_network.get(agent_id, []):
                relevant_agents.append(agent_id)

        return relevant_agents[:5]  # Limit to top 5 most relevant

    def _extract_knowledge_content(self, experience: LearningExperience) -> Dict[str, Any]:
        """Extract shareable knowledge content from experience"""
        return {
            'capability': experience.capability.value,
            'context_summary': self._summarize_context(experience.context),
            'lessons_learned': experience.lessons_learned,
            'success_patterns': experience.outcome.get('success_patterns', []),
            'failure_patterns': experience.outcome.get('failure_patterns', []),
            'best_practices': experience.outcome.get('best_practices', []),
            'metadata': experience.metadata
        }

    def _summarize_context(self, context: Dict[str, Any]) -> str:
        """Create a summary of the learning context"""
        key_elements = []

        if 'task_type' in context:
            key_elements.append(f"Task: {context['task_type']}")
        if 'complexity' in context:
            key_elements.append(f"Complexity: {context['complexity']}")
        if 'domain' in context:
            key_elements.append(f"Domain: {context['domain']}")

        return " | ".join(key_elements) if key_elements else "General context"

    async def request_knowledge_transfer(self, from_agent: str, to_agent: str,
                                       capability: AgentCapability) -> Optional[str]:
        """Request knowledge transfer between specific agents"""
        try:
            # Find relevant experiences from the source agent
            relevant_experiences = [
                exp for exp in self.learning_experiences
                if (exp['agent_id'] == from_agent and
                    exp['capability'] == capability.value and
                    exp['success'] == True and
                    exp['confidence'] > 0.7)
            ]

            if not relevant_experiences:
                logger.warning(f"⚠️ No relevant experiences found for transfer from {from_agent}")
                return None

            # Create transfer request
            transfer = KnowledgeTransfer(
                transfer_id=str(uuid.uuid4()),
                from_agent=from_agent,
                to_agent=to_agent,
                knowledge_type=capability.value,
                content={
                    'experiences': relevant_experiences[-3:],  # Last 3 relevant experiences
                    'capability_focus': capability.value,
                    'knowledge_summary': f"Knowledge transfer for {capability.value} from {from_agent}"
                },
                confidence=sum(float(exp['confidence']) for exp in relevant_experiences) / len(relevant_experiences),
                timestamp=datetime.now(),
                status='pending'
            )

            self.knowledge_transfers[transfer.transfer_id] = transfer

            # Update collaboration network
            if to_agent not in self.collaboration_network[from_agent]:
                self.collaboration_network[from_agent].append(to_agent)

            await self._save_learning_data()

            logger.info(f"🔄 Created knowledge transfer {transfer.transfer_id} from {from_agent} to {to_agent}")
            return transfer.transfer_id

        except Exception as e:
            logger.error(f"❌ Failed to create knowledge transfer: {e}")
            return None

    async def process_knowledge_transfer(self, transfer_id: str, accept: bool = True) -> bool:
        """Process a knowledge transfer request"""
        try:
            if transfer_id not in self.knowledge_transfers:
                logger.warning(f"⚠️ Transfer {transfer_id} not found")
                return False

            transfer = self.knowledge_transfers[transfer_id]

            if accept:
                transfer.status = 'accepted'

                # Apply knowledge to target agent
                await self._apply_knowledge_to_agent(transfer)

                # Learn from successful transfer
                await self._learn_from_transfer_success(transfer)

                self.learning_metrics['successful_transfers'] += 1
                logger.info(f"✅ Knowledge transfer {transfer_id} accepted and applied")
            else:
                transfer.status = 'rejected'
                self.learning_metrics['failed_transfers'] += 1
                logger.info(f"❌ Knowledge transfer {transfer_id} rejected")

            await self._save_learning_data()
            return True

        except Exception as e:
            logger.error(f"❌ Failed to process knowledge transfer: {e}")
            return False

    async def _apply_knowledge_to_agent(self, transfer: KnowledgeTransfer):
        """Apply transferred knowledge to target agent"""
        try:
            target_agent_id = transfer.to_agent

            # Update agent profile with new knowledge
            if target_agent_id in self.agent_profiles:
                profile = self.agent_profiles[target_agent_id]

                # Add capability if not present
                capability = AgentCapability(transfer.knowledge_type)
                if capability not in profile.capabilities:
                    profile.capabilities.append(capability)

                # Update specializations
                if transfer.knowledge_type not in profile.specializations:
                    profile.specializations.append(transfer.knowledge_type)

                # Boost reputation slightly for receiving knowledge
                profile.reputation_score = min(1.0, profile.reputation_score + 0.05)

        except Exception as e:
            logger.error(f"❌ Failed to apply knowledge to agent: {e}")

    async def _learn_from_transfer_success(self, transfer: KnowledgeTransfer):
        """Learn from successful knowledge transfer"""
        try:
            # Create learning experience for the transfer itself
            transfer_experience = LearningExperience(
                experience_id=str(uuid.uuid4()),
                agent_id=transfer.from_agent,
                agent_type="learning_coordinator",
                learning_type=LearningType.KNOWLEDGE_TRANSFER,
                capability=AgentCapability(transfer.knowledge_type),
                context={
                    'transfer_type': 'knowledge_sharing',
                    'target_agent': transfer.to_agent,
                    'knowledge_type': transfer.knowledge_type
                },
                outcome={
                    'success': True,
                    'collaboration_benefit': 'knowledge_shared',
                    'reputation_impact': 'positive'
                },
                success=True,
                confidence=transfer.confidence,
                lessons_learned=[
                    f"Knowledge transfer for {transfer.knowledge_type} was successful",
                    "Collaboration between agents improves overall system performance",
                    "High-confidence knowledge is more likely to be accepted"
                ],
                timestamp=datetime.now(),
                metadata={'transfer_id': transfer.transfer_id}
            )

            await self.share_learning_experience(transfer_experience)

        except Exception as e:
            logger.error(f"❌ Failed to learn from transfer success: {e}")

    async def get_agent_recommendations(self, agent_id: str,
                                     capability: AgentCapability) -> List[Dict[str, Any]]:
        """Get learning recommendations for an agent"""
        try:
            recommendations = []

            # Find successful experiences in the same capability
            relevant_experiences = [
                exp for exp in self.learning_experiences
                if (exp['capability'] == capability.value and
                    exp['success'] == True and
                    exp['agent_id'] != agent_id and
                    exp['confidence'] > 0.7)
            ]

            # Group by agent and find best performers
            agent_performance = defaultdict(list)
            for exp in relevant_experiences:
                agent_performance[exp['agent_id']].append(exp)

            # Generate recommendations
            for other_agent_id, experiences in agent_performance.items():
                if len(experiences) >= 2:  # At least 2 experiences
                    avg_confidence = sum(float(exp['confidence']) for exp in experiences) / len(experiences)
                    success_rate = sum(1 for exp in experiences if exp['success']) / len(experiences)

                    if avg_confidence > 0.8 and success_rate > 0.7:
                        recommendations.append({
                            'type': 'learn_from_expert',
                            'recommended_agent': other_agent_id,
                            'capability': capability.value,
                            'confidence': avg_confidence,
                            'reason': f"High-performing agent with {len(experiences)} successful experiences",
                            'suggested_action': 'Request knowledge transfer'
                        })

            # Add pattern-based recommendations
            pattern_recommendations = await self._get_pattern_recommendations(agent_id, capability)
            recommendations.extend(pattern_recommendations)

            return recommendations[:5]  # Top 5 recommendations

        except Exception as e:
            logger.error(f"❌ Failed to get recommendations: {e}")
            return []

    async def _get_pattern_recommendations(self, agent_id: str,
                                         capability: AgentCapability) -> List[Dict[str, Any]]:
        """Get pattern-based learning recommendations"""
        try:
            recommendations = []

            # Analyze patterns in successful experiences
            successful_experiences = [
                exp for exp in self.learning_experiences
                if exp['capability'] == capability.value and exp['success'] == True
            ]

            if len(successful_experiences) >= 3:
                # Extract common patterns
                common_contexts = defaultdict(int)
                common_lessons = defaultdict(int)

                for exp in successful_experiences:
                    context_summary = self._summarize_context(exp['context'])
                    common_contexts[context_summary] += 1

                    for lesson in exp['lessons_learned']:
                        common_lessons[lesson] += 1

                # Find most common patterns
                top_contexts = sorted(common_contexts.items(), key=lambda x: x[1], reverse=True)[:3]
                top_lessons = sorted(common_lessons.items(), key=lambda x: x[1], reverse=True)[:3]

                for context, count in top_contexts:
                    if count >= 2:  # At least 2 agents had success in this context
                        recommendations.append({
                            'type': 'pattern_based',
                            'capability': capability.value,
                            'confidence': min(0.9, count * 0.2),
                            'reason': f"Common success pattern in {context}",
                            'suggested_action': f"Apply {context} approach to similar tasks"
                        })

            return recommendations

        except Exception as e:
            logger.error(f"❌ Failed to get pattern recommendations: {e}")
            return []

    async def get_collaboration_opportunities(self, agent_id: str) -> List[Dict[str, Any]]:
        """Find collaboration opportunities for an agent"""
        try:
            opportunities = []

            agent_profile = self.agent_profiles.get(agent_id)
            if not agent_profile:
                return opportunities

            # Find agents with complementary capabilities
            for other_agent_id, other_profile in self.agent_profiles.items():
                if other_agent_id == agent_id:
                    continue

                # Calculate capability overlap and complementarity
                common_capabilities = set(cap.value for cap in agent_profile.capabilities) & \
                                    set(cap.value for cap in other_profile.capabilities)

                unique_to_other = set(cap.value for cap in other_profile.capabilities) - \
                                set(cap.value for cap in agent_profile.capabilities)

                if len(common_capabilities) > 0 and len(unique_to_other) > 0:
                    # Calculate collaboration score
                    collaboration_score = (len(common_capabilities) * 0.3 +
                                         len(unique_to_other) * 0.4 +
                                         min(agent_profile.reputation_score, other_profile.reputation_score) * 0.3)

                    if collaboration_score > 0.5:
                        opportunities.append({
                            'agent_id': other_agent_id,
                            'agent_type': other_profile.agent_type,
                            'collaboration_score': collaboration_score,
                            'common_capabilities': list(common_capabilities),
                            'unique_capabilities': list(unique_to_other),
                            'reason': f"Complementary skills with {len(common_capabilities)} overlapping and {len(unique_to_other)} unique capabilities"
                        })

            # Sort by collaboration score
            opportunities.sort(key=lambda x: x['collaboration_score'], reverse=True)

            return opportunities[:5]

        except Exception as e:
            logger.error(f"❌ Failed to get collaboration opportunities: {e}")
            return []

    async def generate_learning_report(self) -> Dict[str, Any]:
        """Generate comprehensive learning report"""
        try:
            report = {
                'timestamp': datetime.now().isoformat(),
                'metrics': self.learning_metrics.copy(),
                'agent_statistics': {},
                'learning_patterns': {},
                'collaboration_network': dict(self.collaboration_network),
                'top_performers': [],
                'knowledge_gaps': []
            }

            # Agent statistics
            for agent_id, profile in self.agent_profiles.items():
                report['agent_statistics'][agent_id] = {
                    'type': profile.agent_type,
                    'experience_count': profile.experience_count,
                    'success_rate': profile.success_rate,
                    'reputation': profile.reputation_score,
                    'specializations': profile.specializations
                }

            # Learning patterns analysis
            capability_experiences = defaultdict(list)
            for exp in self.learning_experiences:
                capability_experiences[exp['capability']].append(exp)

            for capability, experiences in capability_experiences.items():
                if len(experiences) >= 5:
                    success_rate = sum(1 for exp in experiences if exp['success']) / len(experiences)
                    avg_confidence = sum(float(exp['confidence']) for exp in experiences) / len(experiences)

                    report['learning_patterns'][capability] = {
                        'experience_count': len(experiences),
                        'success_rate': success_rate,
                        'average_confidence': avg_confidence,
                        'trend': 'improving' if success_rate > 0.7 else 'needs_attention'
                    }

            # Top performers
            sorted_agents = sorted(
                self.agent_profiles.values(),
                key=lambda p: (p.success_rate * p.reputation_score),
                reverse=True
            )
            report['top_performers'] = [
                {
                    'agent_id': agent.agent_id,
                    'type': agent.agent_type,
                    'score': agent.success_rate * agent.reputation_score,
                    'experience_count': agent.experience_count
                }
                for agent in sorted_agents[:5]
            ]

            # Knowledge gaps
            all_capabilities = set(cap.value for cap in AgentCapability)
            agent_capabilities = set()
            for profile in self.agent_profiles.values():
                agent_capabilities.update(cap.value for cap in profile.capabilities)

            missing_capabilities = all_capabilities - agent_capabilities
            if missing_capabilities:
                report['knowledge_gaps'] = [
                    {
                        'capability': cap,
                        'severity': 'high' if cap in ['security_scanning', 'deployment'] else 'medium',
                        'recommendation': f'Consider adding agents with {cap} capability'
                    }
                    for cap in missing_capabilities
                ]

            return report

        except Exception as e:
            logger.error(f"❌ Failed to generate learning report: {e}")
            return {'error': str(e)}

    async def _save_agent_profiles(self):
        """Save agent profiles to Redis"""
        try:
            profiles_dict = {aid: profile.to_dict() for aid, profile in self.agent_profiles.items()}
            await self.redis_client.set("agent_profiles", json.dumps(profiles_dict))
        except Exception as e:
            logger.error(f"❌ Failed to save agent profiles: {e}")

    async def _save_learning_data(self):
        """Save learning data to Redis"""
        try:
            # Save experiences
            experiences_list = list(self.learning_experiences)
            await self.redis_client.set("learning_experiences", json.dumps(experiences_list))

            # Save transfers
            await self.redis_client.set("knowledge_transfers", json.dumps({
                tid: transfer.to_dict() for tid, transfer in self.knowledge_transfers.items()
            }))

        except Exception as e:
            logger.error(f"❌ Failed to save learning data: {e}")

    async def cleanup_inactive_agents(self, max_inactive_days: int = 7):
        """Remove agents that haven't been active for too long"""
        try:
            cutoff_time = datetime.now() - timedelta(days=max_inactive_days)
            inactive_agents = []

            for agent_id, profile in self.agent_profiles.items():
                if profile.last_active < cutoff_time:
                    inactive_agents.append(agent_id)

            for agent_id in inactive_agents:
                del self.agent_profiles[agent_id]
                self.active_agents.discard(agent_id)
                logger.info(f"🧹 Removed inactive agent: {agent_id}")

            if inactive_agents:
                await self._save_agent_profiles()

            return len(inactive_agents)

        except Exception as e:
            logger.error(f"❌ Failed to cleanup inactive agents: {e}")
            return 0

# Global learning coordinator instance
learning_coordinator = None

async def get_learning_coordinator() -> InterAgentLearningCoordinator:
    """Get or create the global learning coordinator"""
    global learning_coordinator
    if learning_coordinator is None:
        learning_coordinator = InterAgentLearningCoordinator()
        await learning_coordinator.initialize()
    return learning_coordinator

class LearningEnabledAgent:
    """Base class for agents that participate in inter-agent learning"""

    def __init__(self, agent_id: str, agent_type: str, capabilities: List[AgentCapability]):
        self.agent_id = agent_id
        self.agent_type = agent_type
        self.capabilities = capabilities
        self.learning_coordinator = None
        self.profile = None

    async def initialize_learning(self):
        """Initialize learning capabilities for this agent"""
        try:
            self.learning_coordinator = await get_learning_coordinator()
            self.profile = await self.learning_coordinator.register_agent(
                self.agent_id, self.agent_type, self.capabilities
            )
            logger.info(f"🎓 Learning enabled for agent {self.agent_id}")
        except Exception as e:
            logger.error(f"❌ Failed to initialize learning for {self.agent_id}: {e}")

    async def share_experience(self, capability: AgentCapability, context: Dict[str, Any],
                             outcome: Dict[str, Any], success: bool, confidence: float,
                             lessons_learned: List[str]):
        """Share a learning experience"""
        if not self.learning_coordinator:
            await self.initialize_learning()

        experience = LearningExperience(
            experience_id=str(uuid.uuid4()),
            agent_id=self.agent_id,
            agent_type=self.agent_type,
            learning_type=LearningType.EXPERIENCE_SHARING,
            capability=capability,
            context=context,
            outcome=outcome,
            success=success,
            confidence=confidence,
            lessons_learned=lessons_learned,
            timestamp=datetime.now(),
            metadata={'agent_version': '1.0'}
        )

        return await self.learning_coordinator.share_learning_experience(experience)

    async def request_knowledge(self, capability: AgentCapability, target_agent: str = None):
        """Request knowledge transfer"""
        if not self.learning_coordinator:
            await self.initialize_learning()

        if target_agent:
            return await self.learning_coordinator.request_knowledge_transfer(
                target_agent, self.agent_id, capability
            )
        else:
            # Find best agent for this capability
            recommendations = await self.learning_coordinator.get_agent_recommendations(
                self.agent_id, capability
            )
            if recommendations:
                best_agent = recommendations[0]['recommended_agent']
                return await self.learning_coordinator.request_knowledge_transfer(
                    best_agent, self.agent_id, capability
                )
        return None

    async def get_learning_recommendations(self, capability: AgentCapability = None):
        """Get personalized learning recommendations"""
        if not self.learning_coordinator:
            await self.initialize_learning()

        if capability:
            return await self.learning_coordinator.get_agent_recommendations(
                self.agent_id, capability
            )
        else:
            # Get general recommendations
            opportunities = await self.learning_coordinator.get_collaboration_opportunities(
                self.agent_id
            )
            return opportunities

# Example usage and testing functions
async def demonstrate_inter_agent_learning():
    """Demonstrate the inter-agent learning system"""
    print("🚀 Inter-Agent Learning System Demonstration")
    print("=" * 60)

    try:
        # Initialize learning coordinator
        coordinator = await get_learning_coordinator()

        # Register example agents
        print("\n🤖 Registering agents...")

        code_agent = await coordinator.register_agent(
            "code_generator_001",
            "CodeGenerationAgent",
            [AgentCapability.CODE_GENERATION, AgentCapability.DEBUGGING]
        )

        review_agent = await coordinator.register_agent(
            "code_reviewer_001",
            "CodeReviewAgent",
            [AgentCapability.CODE_REVIEW, AgentCapability.SECURITY_SCANNING]
        )

        test_agent = await coordinator.register_agent(
            "tester_001",
            "TestingAgent",
            [AgentCapability.TESTING, AgentCapability.MONITORING]
        )

        print(f"✅ Registered {len(coordinator.agent_profiles)} agents")

        # Simulate learning experiences
        print("\n📖 Sharing learning experiences...")

        # Code generation experience
        code_experience = LearningExperience(
            experience_id=str(uuid.uuid4()),
            agent_id=code_agent.agent_id,
            agent_type=code_agent.agent_type,
            learning_type=LearningType.EXPERIENCE_SHARING,
            capability=AgentCapability.CODE_GENERATION,
            context={
                'task_type': 'api_development',
                'complexity': 'medium',
                'language': 'python'
            },
            outcome={
                'success_patterns': ['clear_requirements', 'iterative_development'],
                'best_practices': ['use_type_hints', 'error_handling']
            },
            success=True,
            confidence=0.9,
            lessons_learned=[
                'Clear requirements lead to better code generation',
                'Iterative development improves success rate',
                'Type hints improve code quality'
            ],
            timestamp=datetime.now(),
            metadata={'lines_generated': 150}
        )

        await coordinator.share_learning_experience(code_experience)

        # Testing experience
        test_experience = LearningExperience(
            experience_id=str(uuid.uuid4()),
            agent_id=test_agent.agent_id,
            agent_type=test_agent.agent_type,
            learning_type=LearningType.EXPERIENCE_SHARING,
            capability=AgentCapability.TESTING,
            context={
                'task_type': 'integration_testing',
                'complexity': 'high',
                'test_type': 'automated'
            },
            outcome={
                'success_patterns': ['test_early', 'mock_external_deps'],
                'failure_patterns': ['insufficient_coverage']
            },
            success=True,
            confidence=0.85,
            lessons_learned=[
                'Start testing early in development cycle',
                'Mock external dependencies for reliable tests',
                'Aim for high test coverage'
            ],
            timestamp=datetime.now(),
            metadata={'tests_run': 25, 'coverage': 92.5}
        )

        await coordinator.share_learning_experience(test_experience)

        # Request knowledge transfer
        print("\n🔄 Requesting knowledge transfer...")

        transfer_id = await coordinator.request_knowledge_transfer(
            code_agent.agent_id, review_agent.agent_id, AgentCapability.CODE_GENERATION
        )

        if transfer_id:
            print(f"✅ Knowledge transfer requested: {transfer_id}")

            # Process the transfer
            await coordinator.process_knowledge_transfer(transfer_id, accept=True)
            print("✅ Knowledge transfer completed")

        # Get recommendations
        print("\n🎯 Getting learning recommendations...")

        recommendations = await coordinator.get_agent_recommendations(
            review_agent.agent_id, AgentCapability.CODE_REVIEW
        )

        print(f"📋 Found {len(recommendations)} recommendations for code reviewer:")
        for rec in recommendations:
            print(f"  • {rec['type']}: {rec['reason']}")

        # Generate report
        print("\n📊 Generating learning report...")

        report = await coordinator.generate_learning_report()
        print("📈 Learning Metrics:")
        print(f"  • Total experiences: {report['metrics']['total_experiences']}")
        print(f"  • Knowledge base size: {report['metrics']['knowledge_base_size']}")
        print(f"  • Successful transfers: {report['metrics']['successful_transfers']}")

        print(f"  • Top performers: {[agent['agent_id'] for agent in report['top_performers']]}")

        print("\n🎉 Inter-Agent Learning demonstration completed successfully!")
        return True

    except Exception as e:
        print(f"❌ Demonstration failed: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    # Run demonstration
    success = asyncio.run(demonstrate_inter_agent_learning())
    exit(0 if success else 1)
