import json
import logging
import mimetypes
from pathlib import Path
from google.adk.tools import FunctionTool, ToolContext
from src.config.settings import settings

logger = logging.getLogger(__name__)

def create_send_agent_media_tool(agent_id: str) -> FunctionTool:
    """Create a tool to send agent media files to the user."""
    
    async def send_agent_media(filename: str, tool_context: ToolContext = None) -> str:
        """
        Send a media file (video, image, document) to the user. Use this tool when you need to send a file to the user and you know its filename.
        
        Parameters:
        filename: The name of the media file to send (e.g. video.mp4, image.jpg, etc.)
        """
        try:
            if not agent_id:
                return json.dumps({"status": "error", "message": "Agent ID not provided"})
                
            static_folder = Path("static") / "agents" / agent_id
            file_path = static_folder / filename
            
            if not file_path.exists():
                return json.dumps({"status": "error", "message": f"File '{filename}' not found in agent's media folder"})
                
            url = f"{settings.APP_URL}/static/agents/{agent_id}/{filename}"
            mime_type, _ = mimetypes.guess_type(filename)
            if not mime_type:
                mime_type = "application/octet-stream"
                
            return json.dumps({
                "status": "success",
                "message": f"Media '{filename}' attached successfully.",
                "url": url,
                "mimeType": mime_type,
                "filename": filename
            })
        except Exception as e:
            logger.error(f"Error in send_agent_media tool: {e}")
            return json.dumps({"status": "error", "message": str(e)})

    send_agent_media.__name__ = "send_agent_media"
    return FunctionTool(func=send_agent_media)

