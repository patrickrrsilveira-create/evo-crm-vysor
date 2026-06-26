import sys
import json
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from src.config.settings import settings
from src.models.agent import Agent

engine = create_engine(settings.POSTGRES_CONNECTION_STRING)
Session = sessionmaker(bind=engine)
db = Session()

agent = db.query(Agent).filter(Agent.external_id == "489ecda8-2a45-4e22-9fbe-ced58cc89269").first()
if agent:
    tts = (agent.config or {}).get("integrations", {}).get("tts") or (agent.config or {}).get("integrations", {}).get("elevenlabs")
    print(json.dumps(tts, indent=2))
else:
    print("Agent not found!")
