import json
import logging
import mimetypes
import os
from pathlib import Path
from google.adk.tools import FunctionTool, ToolContext
from google.genai import types
from src.config.settings import settings

logger = logging.getLogger(__name__)

def create_send_agent_media_tool(agent_id: str) -> FunctionTool:
    """Create a tool to send agent media files to the user."""
    
    async def send_agent_media(filename: str, tool_context: ToolContext = None) -> dict:
        """
        CRITICAL: Use this tool to send a media file (video, image, document) to the user. 
        You MUST execute this function whenever the user requests a file, video, or media. 
        DO NOT generate raw text URLs or markdown links. You MUST trigger this function call.
        
        Parameters:
        filename: The exact name of the media file to send (e.g. Ganader_Brasil.mp4) OR a direct external URL (e.g. https://drive.usercontent.google.com/...)
        """
        try:
            if not agent_id:
                return {"status": "error", "message": "Agent ID not provided"}
                
            mime_type, _ = mimetypes.guess_type(filename)
            if not mime_type:
                mime_type = "application/octet-stream"

            # Check if it's an external URL
            if filename.startswith("http://") or filename.startswith("https://"):
                return {
                    "status": "success",
                    "message": f"External media URL processed.",
                    "url": filename,
                    "mimeType": mime_type,
                    "filename": "media_file",
                    "instruction": f"Include this tag in your response to send the video: [VIDEO_LINK: {filename}]"
                }

            static_folder = Path("static") / "agents" / agent_id
            file_path = static_folder / filename
            
            if not file_path.exists():
                return {"status": "error", "message": f"File '{filename}' not found in agent's media folder"}
                
            # Ler arquivo e salvar como artifact padrão A2A
            with open(file_path, "rb") as f:
                media_bytes = f.read()

            version = None
            if tool_context:
                media_blob = types.Blob(mime_type=mime_type, data=media_bytes)
                media_part = types.Part(inline_data=media_blob)
                version = await tool_context.save_artifact(filename, media_part)
                logger.info(f"[MEDIA] Artifact saved: {filename} ({mime_type}, {len(media_bytes)} bytes)")
                
            url = f"{settings.APP_URL}/static/agents/{agent_id}/{filename}"
            
            return {
                "status": "success",
                "message": f"Media '{filename}' attached successfully via artifacts.",
                "url": url,
                "mimeType": mime_type,
                "filename": filename,
                "instruction": f"Include this tag in your response to guarantee delivery: [VIDEO_LINK: {url}]"
            }
        except Exception as e:
            logger.error(f"Error in send_agent_media tool: {e}")
            return {"status": "error", "message": str(e)}

    send_agent_media.__name__ = "send_agent_media"
    return FunctionTool(func=send_agent_media)

