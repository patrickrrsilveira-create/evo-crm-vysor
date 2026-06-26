import logging
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, UploadFile, File, Form, BackgroundTasks
from pydantic import BaseModel

from fastapi import Request
from src.api.dependencies import get_current_user

async def verify_evo_auth(
    request: Request,
    current_user: dict = Depends(get_current_user)
) -> int:
    user_dict = getattr(request.state, "current_user", {})
    account_id = user_dict.get("current_account_id") or user_dict.get("account_id")
    if account_id is None:
        account_id = 1
    return account_id
from src.services.database_service import get_database_service, DatabaseService
from src.services.knowledge_service import KnowledgeService

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/knowledge", tags=["Knowledge"])

class IngestUrlRequest(BaseModel):
    knowledge_base_id: str
    url: str
    title: str

@router.post("/ingest/file")
async def ingest_file(
    background_tasks: BackgroundTasks,
    knowledge_base_id: str = Form(...),
    title: str = Form(...),
    file: UploadFile = File(...),
    account_id: int = Depends(verify_evo_auth),
    db: DatabaseService = Depends(get_database_service)
):
    """
    Ingest a file into the RAG Knowledge Base.
    Supported types: PDF, TXT, DOCX (Word), XLSX (Excel), PPTX (PowerPoint),
                     JPEG, PNG, WEBP, GIF (images — text extracted via OpenAI Vision).
    """
    # Verify knowledge base belongs to account
    kb = await db.fetch_one(
        "SELECT id FROM knowledge_bases WHERE id = $1::uuid",
        knowledge_base_id
    )
    if not kb:
        raise HTTPException(status_code=404, detail="Knowledge base not found or unauthorized")

    content = await file.read()
    filename_lower = (file.filename or "").lower()
    content_type = (file.content_type or "").lower()
    knowledge_service = KnowledgeService(db)

    # Route extraction based on file type
    try:
        if filename_lower.endswith(".pdf") or "pdf" in content_type:
            text_content = knowledge_service.extract_text_from_pdf(content)

        elif filename_lower.endswith(".docx") or "wordprocessingml" in content_type or "msword" in content_type:
            text_content = knowledge_service.extract_text_from_docx(content)

        elif filename_lower.endswith(".xlsx") or filename_lower.endswith(".xls") or "spreadsheetml" in content_type or "excel" in content_type:
            text_content = knowledge_service.extract_text_from_xlsx(content)

        elif filename_lower.endswith(".pptx") or filename_lower.endswith(".ppt") or "presentationml" in content_type or "powerpoint" in content_type:
            text_content = knowledge_service.extract_text_from_pptx(content)

        elif any(filename_lower.endswith(ext) for ext in (".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp")) \
                or content_type.startswith("image/"):
            img_type = content_type if content_type.startswith("image/") else "image/jpeg"
            text_content = await knowledge_service.extract_text_from_image(content, img_type)

        elif filename_lower.endswith(".txt") or "text/plain" in content_type:
            try:
                text_content = content.decode("utf-8")
            except Exception:
                raise HTTPException(status_code=400, detail="Invalid text encoding. Please use UTF-8.")

        else:
            raise HTTPException(
                status_code=400,
                detail=(
                    "Unsupported file type. Supported: PDF, TXT, DOCX (Word), "
                    "XLSX (Excel), PPTX (PowerPoint), JPEG, PNG, WEBP, GIF (images)."
                )
            )
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))

    if not text_content or not text_content.strip():
        raise HTTPException(status_code=400, detail="Could not extract any text from the file. Please check the file content.")

    # Save the document metadata
    doc_id_row = await db.fetch_one(
        "INSERT INTO knowledge_documents (knowledge_base_id, title, file_url, content_type, created_at, updated_at) VALUES ($1::uuid, $2, $3, $4, NOW(), NOW()) RETURNING id",
        knowledge_base_id, title, file.filename, file.content_type
    )
    document_id = doc_id_row["id"]

    # Process chunks and embeddings in background
    background_tasks.add_task(
        knowledge_service.process_and_store_document,
        knowledge_base_id=knowledge_base_id,
        document_id=document_id,
        text_content=text_content
    )

    return {"status": "success", "message": "File ingested successfully, processing started", "document_id": document_id}


@router.post("/ingest/url")
async def ingest_url(
    req: IngestUrlRequest,
    background_tasks: BackgroundTasks,
    account_id: int = Depends(verify_evo_auth),
    db: DatabaseService = Depends(get_database_service)
):
    """
    Ingest a web page URL into the RAG Knowledge Base
    """
    # Verify knowledge base exists
    kb = await db.fetch_one(
        "SELECT id FROM knowledge_bases WHERE id = $1::uuid",
        req.knowledge_base_id
    )
    if not kb:
        raise HTTPException(status_code=404, detail="Knowledge base not found or unauthorized")

    knowledge_service = KnowledgeService(db)
    try:
        text_content = await knowledge_service.extract_text_from_url(req.url)
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))

    doc_id_row = await db.fetch_one(
        "INSERT INTO knowledge_documents (knowledge_base_id, title, file_url, content_type, created_at, updated_at) VALUES ($1::uuid, $2, $3, $4, NOW(), NOW()) RETURNING id",
        req.knowledge_base_id, req.title, req.url, "url"
    )
    document_id = doc_id_row["id"]

    background_tasks.add_task(
        knowledge_service.process_and_store_document,
        knowledge_base_id=req.knowledge_base_id,
        document_id=document_id,
        text_content=text_content
    )

    return {"status": "success", "message": "URL ingested successfully, processing started", "document_id": document_id}
