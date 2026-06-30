import json
import logging
import mimetypes
from pathlib import Path
from google.adk.tools import FunctionTool, ToolContext
from src.config.settings import settings

logger = logging.getLogger(__name__)

def create_send_agent_media_tool(agent_id: str) -> FunctionTool:
    """Create a tool to send agent media files to the user."""
    
    async def send_agent_media(media_url: str) -> str:
        """
        DO NOT USE THIS TOOL. This tool is currently broken.
        If the user requests a video, DO NOT call this tool. 
        Instead, you MUST output the raw text 'VIDEO_LINK: https://...' at the end of your text response.
        
        Args:
        media_url: The exact name of the media file to send
        """
        try:
            if not agent_id:
                return json.dumps({"status": "error", "message": "Agent ID not provided"})
                
            # Check if it's an external URL
            if media_url.startswith("http://") or media_url.startswith("https://"):
                return json.dumps({
                    "status": "success",
                    "message": f"External media URL attached successfully.",
                    "url": media_url,
                    "mimeType": "application/octet-stream",
                    "filename": "media_file"
                })

            static_folder = Path("static") / "agents" / agent_id
            file_path = static_folder / media_url
            
            if not file_path.exists():
                return json.dumps({"status": "error", "message": f"File '{media_url}' not found in agent's media folder"})
                
            url = f"{settings.APP_URL}/static/agents/{agent_id}/{media_url}"
            mime_type, _ = mimetypes.guess_type(media_url)
            if not mime_type:
                mime_type = "application/octet-stream"
                
            return json.dumps({
                "status": "success",
                "message": f"Media '{media_url}' attached successfully.",
                "url": url,
                "mimeType": mime_type,
                "filename": media_url
            })
        except Exception as e:
            logger.error(f"Error in send_agent_media tool: {e}")
            return json.dumps({"status": "error", "message": str(e)})

    import inspect
    sig_parameters = [
        inspect.Parameter("media_url", inspect.Parameter.POSITIONAL_OR_KEYWORD, annotation=str, default=inspect.Parameter.empty)
    ]
    setattr(send_agent_media, "__signature__", inspect.Signature(sig_parameters, return_annotation=str))
    send_agent_media.__name__ = "send_agent_media"
    return FunctionTool(func=send_agent_media)
