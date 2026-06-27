from fastapi import APIRouter, UploadFile, File, HTTPException
from pathlib import Path
import shutil
from typing import List, Dict, Any
from src.config.settings import settings
import logging

logger = logging.getLogger(__name__)

router = APIRouter()

@router.post("/{agent_bot_id}/media", tags=["Agent Media"])
async def upload_agent_media(agent_bot_id: str, file: UploadFile = File(...)):
    try:
        static_folder = Path("static") / "agents" / agent_bot_id
        static_folder.mkdir(parents=True, exist_ok=True)
        
        file_path = static_folder / file.filename
        
        with open(file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)
            
        return {
            "status": "success", 
            "filename": file.filename,
            "url": f"{settings.APP_URL}/static/agents/{agent_bot_id}/{file.filename}"
        }
    except Exception as e:
        logger.error(f"Error uploading media for agent {agent_bot_id}: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/{agent_bot_id}/media", tags=["Agent Media"])
async def list_agent_media(agent_bot_id: str):
    static_folder = Path("static") / "agents" / agent_bot_id
    
    if not static_folder.exists():
        return {"media": []}
        
    files = []
    try:
        for f in static_folder.iterdir():
            if f.is_file():
                files.append({
                    "filename": f.name, 
                    "url": f"{settings.APP_URL}/static/agents/{agent_bot_id}/{f.name}", 
                    "size": f.stat().st_size
                })
        return {"media": files}
    except Exception as e:
        logger.error(f"Error listing media for agent {agent_bot_id}: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.delete("/{agent_bot_id}/media/{filename}", tags=["Agent Media"])
async def delete_agent_media(agent_bot_id: str, filename: str):
    file_path = Path("static") / "agents" / agent_bot_id / filename
    
    if file_path.exists():
        try:
            file_path.unlink()
            return {"status": "success", "message": f"File {filename} deleted permanently"}
        except Exception as e:
            logger.error(f"Error deleting media {filename} for agent {agent_bot_id}: {e}")
            raise HTTPException(status_code=500, detail=str(e))
            
    raise HTTPException(status_code=404, detail="File not found")
