import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class CartesiaProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice")

        if not token:
            raise ValueError("Cartesia provider requires an API token (apiKey).")
        if not voice_id:
            raise ValueError("Cartesia provider requires a voice_id.")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout) as client:
            json_payload = {
                "model_id": "sonic-english",
                "transcript": text,
                "voice": {
                    "mode": "id",
                    "id": voice_id
                },
                "output_format": {
                    "container": "raw",
                    "encoding": "pcm_s16le",
                    "sample_rate": 24000
                }
            }
                
            response = await client.post(
                "https://api.cartesia.ai/tts/bytes",
                headers={
                    "X-API-Key": token,
                    "Cartesia-Version": "2024-06-10",
                    "Content-Type": "application/json"
                },
                json=json_payload,
            )

            if not response.is_success:
                error_detail = response.text
                raise Exception(f"Cartesia API error: {response.status_code} - {error_detail}")

            return response.content
