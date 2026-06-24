"""
Transfer Conversation Tool (A2A Handoff)

This tool allows AI agents to transfer conversations to other AI agents
using the atomic handoff architecture.
"""

from typing import Optional, Dict, Any, List
from google.adk.tools import FunctionTool, ToolContext
from src.utils.logger import setup_logger
from src.services.handoff_service import transfer_conversation as atomic_transfer
from src.services.event_service import get_event_service
from src.services.adk.tools.evo_crm.transfer_to_human import _extract_conversation_id_from_metadata

logger = setup_logger(__name__)

def create_transfer_conversation_tool(
    transfer_rules: Optional[List[Dict[str, Any]]] = None
) -> FunctionTool:
    """Create the transfer_conversation tool for A2A handoffs.
    
    This tool safely and atomically transfers a conversation to a new AI agent,
    generating a session trace, context updates, and event hooks.
    """
    
    rules = transfer_rules or []
    agent_options = []
    for rule in rules:
        if rule.get("transferTo") == "agent" and rule.get("agentId"):
            agent_options.append(f"- Agent: '{rule.get('name', 'Unknown')}', UUID: '{rule.get('agentId')}'")
    
    agent_options_str = "\n".join(agent_options) if agent_options else "No specific agent routing rules configured."

    async def transfer_conversation(
        target_agent_id: str,
        reason: str,
        conversation_id: Optional[str] = None,
        tool_context: Optional[ToolContext] = None,
    ) -> Dict[str, Any]:
        f"""Transfer the current conversation to another AI Agent.
        
        Use this tool when the current conversation needs to be handled by 
        a different, specialized AI agent (e.g. Sales, Finance, Support).
        
        Available agents to transfer to based on your configuration:
{agent_options_str}

        Args:
            target_agent_id: The UUID of the destination agent to transfer to (MUST be a valid UUID from the available agents list).
            reason: Explanation of why the transfer is happening.
            conversation_id: Auto-extracted conversation ID.
            tool_context: The tool context containing session information (auto-provided).
            
        Returns:
            Dictionary with transfer status.
        """
        try:
            effective_conversation_id = conversation_id
            if not effective_conversation_id and tool_context:
                effective_conversation_id = _extract_conversation_id_from_metadata(tool_context)
                
            if not effective_conversation_id:
                return {
                    "status": "error",
                    "message": "conversation_id could not be extracted from context."
                }

            logger.info(f"Initiating A2A transfer for conversation {effective_conversation_id} to agent {target_agent_id}")

            # 1. Atomic Database Transfer
            transfer_result = await atomic_transfer(
                conversation_id=effective_conversation_id,
                to_agent_id=target_agent_id,
                reason=reason
            )
            
            # 2. Publish Handoff Event
            event_service = get_event_service()
            event_service.publish_handoff_event(
                conversation_id=effective_conversation_id,
                payload=transfer_result
            )
            
            # 3. Request Agent Interruption
            # Return a special flag that signals the engine to stop generating tokens
            return {
                "status": "success",
                "message": f"Conversation transferred successfully to agent {target_agent_id}. The current agent must stop responding.",
                "__system_instruction": "HALT_EXECUTION",
                "transfer_details": transfer_result
            }

        except Exception as e:
            logger.error(f"Error during A2A transfer: {str(e)}")
            return {
                "status": "error",
                "message": f"Failed to transfer conversation: {str(e)}"
            }
            
    transfer_conversation.__name__ = "transfer_conversation"
    
    return FunctionTool(func=transfer_conversation)
