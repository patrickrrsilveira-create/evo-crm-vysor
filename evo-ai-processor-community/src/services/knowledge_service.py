import logging
import io
import tiktoken
from typing import List, Dict, Any, Optional
from pypdf import PdfReader
from openai import AsyncOpenAI
import traceback
import httpx
from bs4 import BeautifulSoup

from src.services.database_service import DatabaseService
from src.services.global_config_service import GlobalConfigService

logger = logging.getLogger(__name__)

class KnowledgeService:
    def __init__(self, db: DatabaseService):
        self.db = db
        self.config_service = GlobalConfigService()
        
        # Tokenizer for splitting
        try:
            self.tokenizer = tiktoken.get_encoding("cl100k_base")
        except Exception as e:
            logger.warning(f"Could not load cl100k_base tokenizer: {e}. Falling back to r50k_base")
            self.tokenizer = tiktoken.get_encoding("r50k_base")

    async def get_openai_client(self) -> AsyncOpenAI:
        """Initialize the OpenAI async client using the global configuration API key."""
        api_key = await self.config_service.get_config("OPENAI_API_SECRET")
        if not api_key:
            # Fallback to OPENAI_API_KEY
            api_key = await self.config_service.get_config("OPENAI_API_KEY")
            
        if not api_key:
            raise ValueError("OpenAI API key is not configured in Global Configs.")
            
        return AsyncOpenAI(api_key=api_key)

    def extract_text_from_pdf(self, file_content: bytes) -> str:
        """Extract text from a PDF file."""
        try:
            pdf_file = io.BytesIO(file_content)
            reader = PdfReader(pdf_file)
            text = ""
            for i, page in enumerate(reader.pages):
                page_text = page.extract_text()
                if page_text:
                    text += f"\\n--- Page {i + 1} ---\\n" + page_text
            return text
        except Exception as e:
            logger.error(f"Error extracting PDF: {e}\\n{traceback.format_exc()}")
            raise ValueError(f"Failed to parse PDF file: {str(e)}")

    async def extract_text_from_url(self, url: str) -> str:
        """Fetch and extract visible text from a URL."""
        try:
            async with httpx.AsyncClient(timeout=10.0) as client:
                response = await client.get(url, follow_redirects=True)
                response.raise_for_status()
                
            soup = BeautifulSoup(response.text, "html.parser")
            
            # Remove scripts, styles, head, title, meta, etc.
            for element in soup(["script", "style", "head", "title", "meta", "[document]"]):
                element.extract()
                
            text = soup.get_text(separator=" ", strip=True)
            return text
        except Exception as e:
            logger.error(f"Error extracting URL {url}: {e}\\n{traceback.format_exc()}")
            raise ValueError(f"Failed to extract text from URL: {str(e)}")

    def chunk_text(self, text: str, max_tokens: int = 500, overlap_tokens: int = 50) -> List[str]:
        """Split text into chunks of maximum `max_tokens` with an overlap."""
        tokens = self.tokenizer.encode(text)
        chunks = []
        
        if not tokens:
            return chunks

        i = 0
        while i < len(tokens):
            chunk_tokens = tokens[i:i + max_tokens]
            chunk_text = self.tokenizer.decode(chunk_tokens)
            chunks.append(chunk_text)
            
            i += (max_tokens - overlap_tokens)
            
        return chunks

    async def generate_embeddings(self, chunks: List[str]) -> List[List[float]]:
        """Generate vector embeddings for a list of text chunks using OpenAI."""
        if not chunks:
            return []
            
        client = await self.get_openai_client()
        
        try:
            # Generate embeddings in batches if necessary, but OpenAI supports multiple inputs up to a certain limit.
            # We'll assume the chunks list isn't astronomically large. A robust implementation would batch this.
            response = await client.embeddings.create(
                model="text-embedding-ada-002",
                input=chunks
            )
            
            # Sort the data by index to ensure order is preserved
            sorted_data = sorted(response.data, key=lambda x: x.index)
            return [data.embedding for data in sorted_data]
            
        except Exception as e:
            logger.error(f"Error generating embeddings: {e}\\n{traceback.format_exc()}")
            raise ValueError(f"Failed to generate embeddings: {str(e)}")

    async def process_and_store_document(
        self, 
        knowledge_base_id: str, 
        document_id: str, 
        text_content: str
    ):
        """Process text, generate chunks, generate embeddings, and store in DB."""
        logger.info(f"Processing document {document_id} for KB {knowledge_base_id}")
        
        # 1. Chunk Text
        chunks = self.chunk_text(text_content, max_tokens=500, overlap_tokens=50)
        logger.info(f"Generated {len(chunks)} chunks for document {document_id}")
        
        if not chunks:
            logger.warning(f"No text extracted for document {document_id}")
            return
            
        # 2. Generate Embeddings
        embeddings = await self.generate_embeddings(chunks)
        
        # 3. Store chunks and embeddings in database
        pool = await self.db.get_pool()
        
        async with pool.acquire() as conn:
            async with conn.transaction():
                for chunk, embedding in zip(chunks, embeddings):
                    # In asyncpg, inserting pgvector is natively supported if pgvector extension is created.
                    # A standard list works. We cast the string formatted list to vector.
                    embedding_str = f"[{','.join(map(str, embedding))}]"
                    
                    await conn.execute(
                        "INSERT INTO knowledge_document_chunks (knowledge_document_id, content, embedding, created_at, updated_at) "
                        "VALUES ($1, $2, $3::vector, NOW(), NOW())",
                        document_id, chunk, embedding_str
                    )
                    
        logger.info(f"Successfully stored {len(chunks)} embedded chunks for document {document_id}")

    async def retrieve_knowledge(self, query: str, knowledge_base_ids: List[str], limit: int = 5) -> str:
        """Search for relevant chunks across given knowledge bases."""
        if not knowledge_base_ids:
            return ""

        logger.info(f"Retrieving knowledge for query: '{query}' in KBs: {knowledge_base_ids}")
        
        # Embed the query
        query_embeddings = await self.generate_embeddings([query])
        if not query_embeddings:
            return ""
            
        query_embedding = query_embeddings[0]
        embedding_str = f"[{','.join(map(str, query_embedding))}]"
        
        pool = await self.db.get_pool()
        
        # The vector operator <-> is L2 distance, <=> is cosine distance, <#> is inner product.
        # Since text-embedding-ada-002 is normalized, cosine distance <=> or inner product <#> work best.
        # We use <=> for cosine distance.
        query_sql = f"""
            SELECT c.content, c.embedding <=> $1::vector AS distance
            FROM knowledge_document_chunks c
            JOIN knowledge_documents d ON c.knowledge_document_id = d.id
            WHERE d.knowledge_base_id = ANY($2::uuid[])
            ORDER BY c.embedding <=> $1::vector
            LIMIT $3
        """
        
        results = await pool.fetch(query_sql, embedding_str, knowledge_base_ids, limit)
        
        if not results:
            return ""
            
        # Combine the content into a single string
        combined_text = "\n\n---\n\n".join([row["content"] for row in results])
        return combined_text

    async def search_agent_knowledge(self, agent_bot_id: str, query: str, limit: int = 5) -> str:
        """Search for relevant chunks across all knowledge bases attached to an agent bot."""
        pool = await self.db.get_pool()
        
        # Get all knowledge_base_ids for this agent_bot_id
        query_sql = "SELECT knowledge_base_id FROM knowledge_base_agent_bots WHERE agent_bot_id = $1"
        records = await pool.fetch(query_sql, agent_bot_id)
        
        knowledge_base_ids = [record['knowledge_base_id'] for record in records]
        
        if not knowledge_base_ids:
            return ""
            
        return await self.retrieve_knowledge(query, knowledge_base_ids, limit)

