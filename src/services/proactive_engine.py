import logging
import asyncio
from datetime import datetime, timezone
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger

from src.services.global_config_service import GlobalConfigService
from src.config.database import SessionLocal
from sqlalchemy import text

logger = logging.getLogger(__name__)

class ProactiveEngineService:
    def __init__(self):
        self.scheduler = AsyncIOScheduler()
        self.global_config = GlobalConfigService()
        self.is_running = False

    def start(self):
        """Starts the proactive engine loop."""
        if self.is_running:
            return
        
        # Runs every hour
        self.scheduler.add_job(
            self.run_campaigns_loop,
            trigger=IntervalTrigger(minutes=60),
            id='proactive_campaign_loop',
            name='Scan and execute proactive campaigns',
            replace_existing=True
        )
        self.scheduler.start()
        self.is_running = True
        logger.info("Proactive Engine (Scheduler) started successfully.")

    async def run_campaigns_loop(self):
        """Main loop that scans the database and queues messages."""
        logger.info("Running proactive campaigns scan...")
        
        # We use a context manager for the DB session
        try:
            db = SessionLocal()
            
            # Fetch active campaigns
            query = text("SELECT id, account_id, trigger_target, delay_hours, message_template, agent_id, attachment_url FROM proactive_campaigns WHERE status = 'ACTIVE'")
            campaigns = db.execute(query).fetchall()
            
            for campaign in campaigns:
                logger.info(f"Processing Campaign ID: {campaign.id} for target: {campaign.trigger_target}")
                
                # Here we would query the Chatwoot API (or DB directly) 
                # to find conversations matching the label (trigger_target)
                # and with last_activity > delay_hours.
                # Push eligible conversations to a Delay Queue
                
                # Pseudocode for the actual delivery queue to prevent bans:
                # asyncio.create_task(self.queue_delivery(eligible_conversations, campaign))
                pass

        except Exception as e:
            logger.error(f"Error in proactive campaign loop: {e}")
        finally:
            db.close()

    async def queue_delivery(self, conversations, campaign):
        """Sends messages with a delay to avoid Meta/WhatsApp bans."""
        import httpx  # For sending async HTTP requests to N8N
        
        N8N_WEBHOOK_URL = "https://n8n.sua-empresa.com.br/webhook/enviar-video" # TODO: Substituir pela URL real do seu N8N
        
        for conv in conversations:
            try:
                # 1. Enviar mensagem de texto normal via CRM API
                # POST /api/v1/accounts/{account_id}/conversations/{conv_id}/messages
                # payload = {"content": campaign.message_template, "private": False}
                
                logger.info(f"Sent proactive message to conversation {conv['id']}")
                
                # 2. Se a campanha tiver uma mídia (attachment_url), dispara pro N8N!
                if hasattr(campaign, 'attachment_url') and campaign.attachment_url:
                    logger.info(f"Campaign {campaign.id} tem mídia associada. Disparando webhook do n8n para conv {conv['id']}...")
                    
                    # Pegamos o número do contato (telefone) do conversation. Exemplo: conv['meta']['sender']['phone_number']
                    # Neste mock, vamos supor que o numero venha no objeto
                    contact_number = conv.get('phone_number', '') 
                    
                    # Vamos tentar adivinhar o mediatype pela extensão ou enviar genérico
                    media_type = "video" if ".mp4" in campaign.attachment_url.lower() else "document"
                    
                    n8n_payload = {
                        "number": contact_number,
                        "mediatype": media_type,
                        "mimetype": "video/mp4" if media_type == "video" else "application/pdf",
                        "caption": "",
                        "media": campaign.attachment_url,
                        "fileName": "arquivo_campanha"
                    }
                    
                    async with httpx.AsyncClient() as client:
                        resp = await client.post(N8N_WEBHOOK_URL, json=n8n_payload)
                        logger.info(f"N8N disparado com sucesso: {resp.status_code}")
                
                # Apply idempotency label
                # POST /api/v1/accounts/{account_id}/conversations/{conv_id}/labels
                # payload = {"labels": [f"proact_done_{campaign.id}"]}
                
                # Anti-ban delay
                await asyncio.sleep(30)
            except Exception as e:
                logger.error(f"Delivery failed for conv {conv.get('id', 'unknown')}: {e}")

# Singleton instance
proactive_engine = ProactiveEngineService()
