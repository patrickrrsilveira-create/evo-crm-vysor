import sys
sys.path.append('evo-ai-processor-community')
import asyncio
from src.config.database import async_session_maker
from src.models.agent import Agent
from sqlalchemy import select

async def main():
    async with async_session_maker() as s:
        res = await s.execute(select(Agent).where(Agent.name=='ganader'))
        agent = res.scalars().first()
        if agent:
            print('Enabled tools:', agent.config.get('enabled_tools', 'NOT_SET') if agent.config else 'NO_CONFIG')
        else:
            print('Agent not found')

asyncio.run(main())
