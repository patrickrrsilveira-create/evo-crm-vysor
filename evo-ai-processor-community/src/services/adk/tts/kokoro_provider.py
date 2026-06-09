import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class KokoroProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice") or "af_heart"
        api_url = config.get("api_url", "https://api.kokoro.io/v1/audio/speech")

        if token and "|" in token:
            api_url, token = token.split("|", 1)

        if not token:
            raise ValueError("Kokoro provider requires an API token (apiKey).")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout) as client:
            json_payload = {
                "input": text,
                "voice": voice_id,
                "model": "kokoro-v1"
            }
                
            response = await client.post(
                api_url,
                headers={
                    "Authorization": f"Bearer {token}",
                    "Content-Type": "application/json"
                },
                json=json_payload,
            )

            if not response.is_success:
                error_detail = response.text
                raise Exception(f"Kokoro API error: {response.status_code} - {error_detail}")

            return response.content
