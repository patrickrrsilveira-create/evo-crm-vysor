import asyncio
from src.services.database_service import DatabaseService
from src.config.settings import settings
import json

async def main():
    db = DatabaseService(settings.DATABASE_URL)
    query = """
        SELECT name, serialized_value FROM installation_configs 
        WHERE name IN ('OPENAI_API_SECRET', 'OPENAI_API_KEY', 'OPENAI_API_URL')
    """
    results = await db.fetch_all(query)
    for row in results:
        val = row['serialized_value']
        print(f"{row['name']} = {repr(val)} (type: {type(val)})")

if __name__ == '__main__':
    asyncio.run(main())
