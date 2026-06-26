import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class VoxtralProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice") or "default"
        api_url = config.get("api_url", "https://api.voxtral.ai/v1/tts")

        if token and "|" in token:
            api_url, token = token.split("|", 1)

        if not token:
            raise ValueError("Voxtral provider requires an API token (apiKey).")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout, verify=False) as client:
            json_payload = {
                "text": text,
                "voice_id": voice_id
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
                raise Exception(f"Voxtral API error: {response.status_code} - {error_detail}")

            return response.content
