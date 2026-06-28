import httpx
import asyncio

async def run():
    async with httpx.AsyncClient() as client:
        response = await client.post('https://openrouter.ai/api/v1/audio/speech', json={'model':'hexgrad/kokoro-82m','input':'test','voice':'af_heart'})
        print("Status:", response.status_code)
        print("Content-Type:", response.headers.get('content-type'))
        print("Body:", response.text[:200])

asyncio.run(run())
