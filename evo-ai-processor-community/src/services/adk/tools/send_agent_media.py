from typing import Dict, Any, Type, Optional
from pydantic import BaseModel, Field
import os
from pathlib import Path
import mimetypes

from src.services.adk.tools.base import BaseTool
from src.config.settings import settings

class SendAgentMediaInput(BaseModel):
    filename: str = Field(..., description="The name of the media file to send (e.g. video.mp4, image.jpg, etc.)")

class SendAgentMediaTool(BaseTool):
    name = "send_agent_media"
    description = "Send a media file (video, image, document) to the user. Use this tool when you need to send a file to the user and you know its filename."
    args_schema: Type[BaseModel] = SendAgentMediaInput
    agent_id: Optional[str] = None
    
    def __init__(self, agent_id: Optional[str] = None):
        self.agent_id = agent_id
        
    async def _run(self, filename: str) -> Dict[str, Any]:
        if not self.agent_id:
            return {"status": "error", "message": "Agent ID not provided"}
            
        static_folder = Path("static") / "agents" / self.agent_id
        file_path = static_folder / filename
        
        if not file_path.exists():
            return {"status": "error", "message": f"File '{filename}' not found in agent's media folder"}
            
        url = f"{settings.APP_URL}/static/agents/{self.agent_id}/{filename}"
        mime_type, _ = mimetypes.guess_type(filename)
        if not mime_type:
            mime_type = "application/octet-stream"
            
        return {
            "status": "success",
            "message": f"Media '{filename}' attached successfully.",
            "url": url,
            "mimeType": mime_type,
            "filename": filename
        }
