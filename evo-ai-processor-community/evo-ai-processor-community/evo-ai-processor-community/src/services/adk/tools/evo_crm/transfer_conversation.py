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
    
    # Build agent options text from transfer rules
    rules = transfer_rules or []
    agent_lines = []
    for rule in rules:
        if rule.get("transferTo") == "agent" and rule.get("agentId"):
            name = rule.get("name") or rule.get("agentName") or "Unknown"
            agent_id = rule.get("agentId")
            instructions = rule.get("instructions", "")
            line = f"- Agent '{name}' (UUID: {agent_id})"
            if instructions:
                line += f" — {instructions}"
            agent_lines.append(line)
    
    if agent_lines:
        agents_doc = "Available agents:\n" + "\n".join(agent_lines)
    else:
        agents_doc = "No specific agent routing rules configured."

    async def transfer_conversation(
        target_agent_id: str,
        reason: str,
        conversation_id: Optional[str] = None,
        tool_context: Optional[ToolContext] = None,
    ) -> Dict[str, Any]:
        """Transfer the current conversation to another AI Agent."""
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

    # Dynamically set the docstring (f-string docstrings don't work in Python)
    transfer_conversation.__doc__ = (
        "Transfer the current conversation to another AI Agent.\n\n"
        "Use this tool when the current conversation needs to be handled by "
        "a different, specialized AI agent (e.g. Sales, Finance, Support).\n\n"
        f"{agents_doc}\n\n"
        "If the agent you want to transfer to is NOT listed above, you can still transfer to them by passing their EXACT name as the target_agent_id.\n\n"
        "Args:\n"
        "    target_agent_id: The UUID or EXACT NAME of the agent to transfer to.\n"
        "    reason: A brief explanation of why the conversation is being transferred.\n"
        "    conversation_id: Auto-extracted conversation ID (do not provide)."
    )
    transfer_conversation.__name__ = "transfer_conversation"
    
    return FunctionTool(func=transfer_conversation)
