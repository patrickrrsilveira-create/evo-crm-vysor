"""
Service for publishing events to other components (e.g., CRM, Swarm Engine).
"""
import json
import logging
from typing import Dict, Any
import redis
from src.config.redis import get_redis_config

import os
import requests

logger = logging.getLogger(__name__)

class EventService:
    def __init__(self):
        redis_config = get_redis_config()
        self.redis = redis.Redis(
            host=redis_config["host"],
            port=redis_config["port"],
            db=redis_config["db"],
            password=redis_config.get("password"),
            ssl=redis_config.get("ssl", False)
        )
    
    def publish_handoff_event(self, conversation_id: str, payload: Dict[str, Any]):
        """Publish a handoff event to the Redis pubsub channel and trigger Rails Webhook."""
        try:
            event = {
                "event": "conversation.handoff",
                "conversation_id": conversation_id,
                "data": payload
            }
            # 1. Publish to Redis (for legacy or other microservices)
            channel = "evoai.events.conversations"
            self.redis.publish(channel, json.dumps(event))
            logger.info(f"Published handoff event for conversation {conversation_id} to Redis")
            
            # 2. Trigger Rails Webhook for UI/ActionCable Update
            webhook_url = f"{os.getenv('EVO_AI_CRM_URL', 'http://localhost:3000')}/webhooks/agent_processor/handoff"
            headers = {
                "Authorization": f"Bearer {os.getenv('EVOAI_CRM_API_TOKEN', '')}",
                "Content-Type": "application/json"
            }
            response = requests.post(
                webhook_url,
                json={"conversation_id": conversation_id, "payload": payload},
                headers=headers,
                timeout=5.0
            )
            
            if response.status_code in (200, 201):
                logger.info(f"Successfully triggered handoff webhook for conversation {conversation_id}")
            else:
                logger.warning(f"Handoff webhook failed with status {response.status_code}: {response.text}")
                
        except Exception as e:
            logger.error(f"Failed to publish handoff event: {str(e)}")

_event_service = None

def get_event_service() -> EventService:
    global _event_service
    if _event_service is None:
        _event_service = EventService()
    return _event_service
