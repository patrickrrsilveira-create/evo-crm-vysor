"""
Service for handling conversation handoffs between agents using atomic transactions.
"""

import logging
import json
import uuid
import litellm
from pydantic import BaseModel
from typing import Dict, Any, Optional
from src.services.database_service import get_database_service
from src.services.global_config_service import GlobalConfigService

logger = logging.getLogger(__name__)

class HandoffError(Exception):
    pass

class OptimisticLockError(HandoffError):
    pass

class HandoffContext(BaseModel):
    summary: str
    entities: Dict[str, Any]

async def extract_conversation_context(conversation_id: str) -> Dict[str, Any]:
    """
    Extracts conversation context using the AI model configured in CRM's Global Settings.
    """
    try:
        db_service = get_database_service()
        pool = await db_service.get_pool()
        
        # Fetch latest session state
        query = """
            SELECT state 
            FROM evo_ai_agent_processor_sessions 
            WHERE user_id = $1 
            ORDER BY create_time DESC 
            LIMIT 1
        """
        async with pool.acquire() as connection:
            session = await connection.fetchrow(query, conversation_id)
            
        if not session or not session['state']:
            return {"summary": "Sem histórico disponível.", "entities": {}}
            
        state_json = session['state']
        if isinstance(state_json, str):
            state_json = json.loads(state_json)
            
        # Parse history
        history_text = ""
        for event in state_json:
            if isinstance(event, dict):
                actions = event.get('actions', {})
                if not isinstance(actions, dict):
                    # Sometimes actions is serialized as string in DB, or it's an object. 
                    # We'll just convert to string for safety.
                    actions_str = str(actions)
                    history_text += f"{actions_str}\n"
                    continue
                state_delta = actions.get('state_delta', {})
                if 'user_input' in state_delta:
                    history_text += f"User: {state_delta['user_input']}\n"
                if 'llm_response' in state_delta:
                    history_text += f"Agent: {state_delta['llm_response']}\n"
        
        if not history_text.strip():
            history_text = str(state_json)[:2000] # Fallback to raw state dump
            
        # Get AI config from GlobalConfigService
        config_service = GlobalConfigService()
        model = await config_service.get_config("OPENAI_MODEL") or "gpt-4o-mini"
        api_key = await config_service.get_config("OPENAI_API_SECRET")
        api_base = await config_service.get_config("OPENAI_API_URL")
        
        if not api_key:
            logger.warning("OPENAI_API_SECRET not found in Global Config. Context extraction may fail.")
            
        prompt = f"""You are transferring a conversation to a new agent.
Please read the following conversation history and extract:
1. A concise summary of the user's intent, current situation, and what needs to be done next.
2. Any important entities mentioned (e.g. name, email, product, issue).

History:
{history_text}
"""
        response = await litellm.acompletion(
            model=model,
            api_key=api_key,
            api_base=api_base,
            messages=[{"role": "user", "content": prompt}],
            response_format=HandoffContext,
            max_tokens=1024,
            temperature=0.3
        )
        
        # Parse Pydantic response
        content = response.choices[0].message.content
        if isinstance(content, str):
            parsed = json.loads(content)
            return parsed
        return {"summary": "Não foi possível extrair.", "entities": {}}
        
    except Exception as e:
        logger.error(f"Error extracting conversation context: {e}")
        return {"summary": "Erro na extração de contexto.", "entities": {}}

async def transfer_conversation(
    conversation_id: str,
    to_agent_id: str,
    reason: str = "Handoff requested",
    summary: str = "",
    entities: Optional[Dict[str, Any]] = None
) -> Dict[str, Any]:
    """
    Atomically transfers a conversation from the current agent to a new agent.
    Implements optimistic locking to prevent race conditions.
    """
    if not summary:
        # Extrai o contexto sincronamente (bloqueia o fluxo até concluir) 
        # antes de entrar na transação atômica do banco, para não segurar locks no DB.
        context = await extract_conversation_context(conversation_id)
        summary = context.get('summary', '')
        entities = context.get('entities', {})
    elif entities is None:
        entities = {}

    db_service = get_database_service()
    pool = await db_service.get_pool()
    
    async with pool.acquire() as connection:
        async with connection.transaction():
            # Resolve to_agent_id if it's not a valid UUID
            try:
                import uuid
                uuid.UUID(to_agent_id)
            except ValueError:
                logger.info(f"to_agent_id '{to_agent_id}' is not a UUID, attempting name resolution.")
                query_agents = "SELECT id, name FROM evo_core_agents"
                agents = await connection.fetch(query_agents)
                
                import re
                matched_id = None
                target_clean = re.sub(r'[^a-z0-9]', '', to_agent_id.lower())
                for agent in agents:
                    db_name_clean = re.sub(r'[^a-z0-9]', '', agent['name'].lower())
                    if db_name_clean in target_clean or target_clean in db_name_clean:
                        matched_id = str(agent['id'])
                        logger.info(f"Fuzzy matched agent name '{agent['name']}' to UUID: {matched_id}")
                        break
                        
                if matched_id:
                    to_agent_id = matched_id
                else:
                    available = [a['name'] for a in agents]
                    raise HandoffError(f"Agent with name '{to_agent_id}' not found in database. Available agents: {available}")
            
            # Resolve evo_core_agents ID to agent_bots ID
            query_bot = "SELECT id FROM agent_bots WHERE outgoing_url LIKE '%' || $1"
            bot_row = await connection.fetchrow(query_bot, to_agent_id)
            if not bot_row:
                raise HandoffError(f"Agent {to_agent_id} is not linked to any agent_bot.")
            bot_id = str(bot_row['id'])

            # 1. Fetch current conversation state
            query_conv = """
                SELECT id, active_agent_id, state_version 
                FROM conversations 
                WHERE id = $1
            """
            # UUID in asyncpg needs to be handled, we can cast string to uuid in SQL
            conv = await connection.fetchrow(query_conv, uuid.UUID(conversation_id))
            
            if not conv:
                raise HandoffError(f"Conversation {conversation_id} not found")
                
            current_agent_id = conv['active_agent_id']
            current_version = conv['state_version']
            
            # 2. Update conversation with Optimistic Locking
            update_conv = """
                UPDATE conversations 
                SET active_agent_id = $1,
                    state_version = state_version + 1,
                    state = 'ACTIVE',
                    transfer_lock = false,
                    updated_at = NOW()
                WHERE id = $2 AND state_version = $3
                RETURNING id
            """
            updated = await connection.fetchrow(
                update_conv, 
                uuid.UUID(bot_id), 
                uuid.UUID(conversation_id), 
                current_version
            )
            
            if not updated:
                raise OptimisticLockError("Conversation state changed during transfer. Please retry.")

            # 3. Create Transfer Record
            insert_transfer = """
                INSERT INTO conversation_transfers (id, conversation_id, from_agent_id, to_agent_id, reason, status, created_at, completed_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, 'completed', NOW(), NOW(), NOW())
            """
            transfer_id = uuid.uuid4()
            await connection.execute(
                insert_transfer,
                transfer_id,
                uuid.UUID(conversation_id),
                current_agent_id,
                uuid.UUID(bot_id),
                reason
            )

            # 4. Create New Agent Session
            insert_session = """
                INSERT INTO agent_sessions (id, conversation_id, agent_bot_id, summary, entities, state, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, 'ACTIVE', NOW(), NOW())
            """
            session_id = uuid.uuid4()
            import json
            await connection.execute(
                insert_session,
                session_id,
                uuid.UUID(conversation_id),
                uuid.UUID(bot_id),
                summary,
                json.dumps(entities)
            )

            # 5. Upsert Conversation Context
            # Using INSERT ... ON CONFLICT DO UPDATE
            upsert_context = """
                INSERT INTO conversation_contexts (id, conversation_id, summary, entities, created_at, updated_at)
                VALUES ($1, $2, $3, $4, NOW(), NOW())
                ON CONFLICT (conversation_id) DO UPDATE 
                SET summary = EXCLUDED.summary,
                    entities = EXCLUDED.entities,
                    updated_at = NOW()
            """
            # We need to make sure conversation_id is unique in contexts, or we just insert if no unique constraint.
            # Assuming conversation_id is unique per context. Wait, we didn't add unique constraint in migration!
            # If not unique constraint, we can delete existing and insert new, or just insert new and use latest.
            # Let's delete existing first to be safe, then insert.
            
            await connection.execute("DELETE FROM conversation_contexts WHERE conversation_id = $1", uuid.UUID(conversation_id))
            
            context_id = uuid.uuid4()
            await connection.execute(
                """
                INSERT INTO conversation_contexts (id, conversation_id, summary, entities, created_at, updated_at)
                VALUES ($1, $2, $3, $4, NOW(), NOW())
                """,
                context_id,
                uuid.UUID(conversation_id),
                summary,
                json.dumps(entities)
            )
            
            logger.info(f"Successfully transferred conversation {conversation_id} to agent {to_agent_id}")
            
            return {
                "success": True,
                "conversation_id": conversation_id,
                "transfer_id": str(transfer_id),
                "session_id": str(session_id),
                "to_agent_id": to_agent_id
            }
