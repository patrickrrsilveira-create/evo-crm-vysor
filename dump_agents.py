import asyncio
import os
import sys

# Add the project root to sys.path so we can import src
sys.path.insert(0, '/app')

from src.config.database import SessionLocal
from src.models.models import Agent

def dump_agents():
    db = SessionLocal()
    try:
        agents = db.query(Agent).all()
        print(f"Found {len(agents)} agents in evo_core_agents:")
        for agent in agents:
            print(f"- ID: {agent.id}, Name: '{agent.name}'")
            
        # Let's also check raw SQL for agent_bots if possible
        from sqlalchemy import text
        bots = db.execute(text("SELECT id, name FROM agent_bots")).fetchall()
        print(f"\nFound {len(bots)} agents in agent_bots:")
        for bot in bots:
            print(f"- ID: {bot.id}, Name: '{bot.name}'")
            
    finally:
        db.close()

if __name__ == "__main__":
    dump_agents()
