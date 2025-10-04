from sentence_transformers import SentenceTransformer
import numpy as np
import json
import time
from datetime import datetime
from typing import Dict, List, Any, Optional, Tuple
import asyncio
import os
import redis

class MemoAIMemory:
    """Lightweight vector memory system for AI agent experiences"""

    def __init__(self, redis_client=None, model_name: str = 'all-MiniLM-L6-v2'):
        self.redis_client = redis_client
        self.model_name = model_name

        # Initialize the embedding model
        try:
            self.encoder = SentenceTransformer(model_name)
        except Exception as e:
            print(f"Failed to load model {model_name}: {e}")
            # Fallback to a simpler model
            self.encoder = SentenceTransformer('paraphrase-MiniLM-L3-v2')

        # Initialize simple numpy-based storage
        self.embedding_dim = self.encoder.get_sentence_embedding_dimension()
        self.embeddings = []
        self.memory_store = {}
        self.memory_counter = 0

        # Load existing memories if available
        self._load_persistent_memories()

    def _load_persistent_memories(self):
        """Load memories from persistent storage"""
        try:
            if os.path.exists('factory_memories.json'):
                with open('factory_memories.json', 'r') as f:
                    data = json.load(f)

                self.memory_store = data.get('memories', {})
                self.memory_counter = data.get('counter', 0)

                # Rebuild embeddings list
                if self.memory_store:
                    for mem_id, memory in self.memory_store.items():
                        embedding = memory.get('embedding', [])
                        if embedding:
                            self.embeddings.append(np.array(embedding))

        except Exception as e:
            print(f"Failed to load persistent memories: {e}")

    def _save_persistent_memories(self):
        """Save memories to persistent storage"""
        try:
            data = {
                'memories': self.memory_store,
                'counter': self.memory_counter,
                'timestamp': datetime.now().isoformat()
            }

            with open('factory_memories.json', 'w') as f:
                json.dump(data, f, indent=2)

        except Exception as e:
            print(f"Failed to save persistent memories: {e}")

    async def remember(self, key: str, content: Dict, metadata: Dict = None) -> int:
        """Store experience in vector memory"""
        try:
            # Create text representation
            text_parts = [key]
            if isinstance(content, dict):
                text_parts.extend([f"{k}: {v}" for k, v in content.items()])
            else:
                text_parts.append(str(content))

            text = " | ".join(text_parts)

            # Create embedding
            embedding = self.encoder.encode([text])[0]

            # Create memory entry
            memory_id = self.memory_counter
            self.memory_counter += 1

            memory_entry = {
                "id": memory_id,
                "key": key,
                "content": content,
                "text": text,
                "embedding": embedding.tolist(),
                "timestamp": datetime.now().isoformat(),
                "metadata": metadata or {}
            }

            self.memory_store[memory_id] = memory_entry
            self.embeddings.append(embedding)

            # Store in Redis if available
            if self.redis_client:
                await self._store_in_redis(memory_id, memory_entry)

            # Save to persistent storage periodically
            if memory_id % 10 == 0:
                self._save_persistent_memories()

            return memory_id

        except Exception as e:
            print(f"Failed to store memory: {e}")
            return -1

    async def recall(self, query: str, k: int = 5, threshold: float = 0.1) -> List[Dict]:
        """Retrieve relevant memories"""
        try:
            # Create query embedding
            query_embedding = self.encoder.encode([query])[0]

            if not self.embeddings:
                return []

            # Calculate similarities with all stored embeddings
            similarities = []
            for i, embedding in enumerate(self.embeddings):
                # Calculate cosine similarity
                similarity = np.dot(query_embedding, embedding) / (
                    np.linalg.norm(query_embedding) * np.linalg.norm(embedding)
                )
                similarities.append((i, similarity))

            # Sort by similarity (highest first)
            similarities.sort(key=lambda x: x[1], reverse=True)

            memories = []
            for idx, similarity in similarities:
                if similarity < threshold:
                    continue

                memory_id = list(self.memory_store.keys())[idx]
                memory = self.memory_store[memory_id].copy()
                memory["similarity"] = float(similarity)
                memory["distance"] = 1.0 - similarity
                memories.append(memory)

                if len(memories) >= k:
                    break

            return memories

        except Exception as e:
            print(f"Failed to recall memories: {e}")
            return []

    async def search_by_key(self, key_pattern: str, limit: int = 10) -> List[Dict]:
        """Search memories by key pattern"""
        try:
            matching_memories = []

            for memory_id, memory in self.memory_store.items():
                if key_pattern.lower() in memory["key"].lower():
                    memory_copy = memory.copy()
                    memory_copy["similarity"] = 1.0  # Exact match
                    matching_memories.append(memory_copy)

            # Sort by timestamp (newest first)
            matching_memories.sort(key=lambda x: x["timestamp"], reverse=True)
            return matching_memories[:limit]

        except Exception as e:
            print(f"Failed to search by key: {e}")
            return []

    async def get_memory_stats(self) -> Dict:
        """Get memory system statistics"""
        try:
            return {
                "total_memories": len(self.memory_store),
                "embeddings_count": len(self.embeddings),
                "embedding_dimension": self.embedding_dim,
                "model_name": self.model_name,
                "last_updated": datetime.now().isoformat()
            }
        except Exception as e:
            print(f"Failed to get stats: {e}")
            return {}

    async def clear_memory(self, older_than_days: int = None) -> int:
        """Clear memories, optionally older than specified days"""
        try:
            if older_than_days is None:
                # Clear all memories
                count = len(self.memory_store)
                self.memory_store.clear()
                self.embeddings.clear()
                self.memory_counter = 0
            else:
                # Clear old memories
                cutoff_time = datetime.now().timestamp() - (older_than_days * 24 * 60 * 60)
                count = 0

                to_remove = []
                new_embeddings = []

                for memory_id, memory in self.memory_store.items():
                    memory_time = datetime.fromisoformat(memory["timestamp"]).timestamp()
                    if memory_time < cutoff_time:
                        to_remove.append(memory_id)
                        count += 1
                    else:
                        # Keep this memory and its embedding
                        embedding_index = list(self.memory_store.keys()).index(memory_id)
                        if embedding_index < len(self.embeddings):
                            new_embeddings.append(self.embeddings[embedding_index])

                # Remove old memories
                for memory_id in to_remove:
                    del self.memory_store[memory_id]

                # Update embeddings list
                self.embeddings = new_embeddings

            # Save changes
            self._save_persistent_memories()

            return count

        except Exception as e:
            print(f"Failed to clear memory: {e}")
            return 0

    async def _store_in_redis(self, memory_id: int, memory_entry: Dict):
        """Store memory in Redis"""
        try:
            if self.redis_client:
                await self.redis_client.hset(
                    f"memo:{memory_id}",
                    mapping={
                        "key": memory_entry["key"],
                        "content": json.dumps(memory_entry["content"]),
                        "timestamp": memory_entry["timestamp"],
                        "text": memory_entry["text"]
                    }
                )
        except Exception as e:
            print(f"Failed to store in Redis: {e}")

    async def get_similar_experiences(self, current_situation: str, limit: int = 3) -> List[Dict]:
        """Get experiences similar to current situation"""
        return await self.recall(current_situation, k=limit)

    async def store_success_pattern(self, task: str, solution: Dict, outcome: str):
        """Store a successful pattern for future reference"""
        content = {
            "task": task,
            "solution": solution,
            "outcome": outcome,
            "type": "success_pattern"
        }

        await self.remember(f"success_{task}", content, {
            "pattern_type": "success",
            "reusability": "high"
        })

    async def store_failure_pattern(self, task: str, error: str, solution_attempted: Dict):
        """Store a failure pattern to avoid in the future"""
        content = {
            "task": task,
            "error": error,
            "solution_attempted": solution_attempted,
            "type": "failure_pattern"
        }

        await self.remember(f"failure_{task}", content, {
            "pattern_type": "failure",
            "reusability": "avoid"
        })

# Example usage and testing
async def test_memory_system():
    """Test the memory system"""
    # Initialize memory system
    memory = MemoAIMemory()

    # Store some experiences
    await memory.store_success_pattern(
        "user_authentication",
        {"method": "jwt", "secure": True},
        "successful_implementation"
    )

    await memory.store_failure_pattern(
        "database_connection",
        "connection_timeout",
        {"host": "localhost", "timeout": 30}
    )

    # Recall similar experiences
    similar = await memory.recall("user login system", k=3)
    print(f"Found {len(similar)} similar experiences")

    for mem in similar:
        print(f"  - {mem['key']}: {mem['content']['type']}")

    # Get stats
    stats = await memory.get_memory_stats()
    print(f"Memory stats: {stats}")

if __name__ == "__main__":
    asyncio.run(test_memory_system())