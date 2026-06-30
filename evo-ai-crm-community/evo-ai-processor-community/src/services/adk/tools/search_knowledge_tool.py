import json
import logging
from google.adk.tools import FunctionTool

logger = logging.getLogger(__name__)

from src.services.database_service import DatabaseService

def create_search_knowledge_tool(agent_id: str, db: DatabaseService) -> FunctionTool:
    """Create a tool to search the agent's native knowledge bases."""
    
    async def search_knowledge(query: str, limit: int = 5) -> str:
        """
        Search the native CRM Knowledge Base (Base de Conhecimento) for manuals, pricing, and internal documents.
        Always use this tool when the user asks about pricing, technical manuals, or specific company knowledge.
        
        Args:
        query: The search query to look for in the knowledge base.
        limit: Number of text chunks to return. Default 5.
        
        Returns:
        The extracted knowledge text or a message indicating no information was found.
        """
        try:
            from src.services.knowledge_service import KnowledgeService
            knowledge_service = KnowledgeService(db)
            
            # This calls the service which already checks knowledge_base_agent_bots table
            result = await knowledge_service.search_agent_knowledge(
                agent_bot_id=agent_id,
                query=query,
                limit=limit
            )
            
            if not result:
                return json.dumps({"status": "no_results", "message": "Nenhuma informação relevante encontrada na Base de Conhecimento para esta busca."})
                
            return json.dumps({"status": "success", "content": result})
        except Exception as e:
            logger.error(f"Error searching knowledge base: {e}")
            return json.dumps({"status": "error", "message": f"An error occurred: {str(e)}"})
            
    search_knowledge.__name__ = "search_knowledge"
    return FunctionTool(func=search_knowledge)
