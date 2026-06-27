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
            # For simplicity in this blueprint, we use raw SQL or SQLAlchemy to fetch campaigns
            query = text("SELECT id, account_id, trigger_target, delay_hours, message_template, agent_id FROM proactive_campaigns WHERE status = 'ACTIVE'")
            campaigns = db.execute(query).fetchall()
            
            for campaign in campaigns:
                logger.info(f"Processing Campaign ID: {campaign.id} for target: {campaign.trigger_target}")
                
                # Here we would query the Chatwoot API (or DB directly) 
                # to find conversations matching the label (trigger_target)
                # and with last_activity > delay_hours.
                # Example:
                # 1. Hit GET /api/v1/accounts/{account_id}/conversations?labels={trigger_target}
                # 2. Filter by updated_at
                # 3. Exclude conversations with label "proact_done_{campaign.id}"
                # 4. Push eligible conversations to a Delay Queue (Redis or asyncio.sleep queue)
                
                # Pseudocode for the actual delivery queue to prevent bans:
                # asyncio.create_task(self.queue_delivery(eligible_conversations, campaign))
                pass

        except Exception as e:
            logger.error(f"Error in proactive campaign loop: {e}")
        finally:
            db.close()

    async def queue_delivery(self, conversations, campaign):
        """Sends messages with a delay to avoid Meta/WhatsApp bans."""
        for conv in conversations:
            try:
                # Send Message via CRM API
                # POST /api/v1/accounts/{account_id}/conversations/{conv_id}/messages
                # payload = {"content": campaign.message_template, "private": False}
                
                # Apply idempotency label
                # POST /api/v1/accounts/{account_id}/conversations/{conv_id}/labels
                # payload = {"labels": [f"proact_done_{campaign.id}"]}
                
                logger.info(f"Sent proactive message to conversation {conv['id']}")
                
                # Anti-ban delay
                await asyncio.sleep(30)
            except Exception as e:
                logger.error(f"Delivery failed for conv {conv['id']}: {e}")

# Singleton instance
proactive_engine = ProactiveEngineService()
