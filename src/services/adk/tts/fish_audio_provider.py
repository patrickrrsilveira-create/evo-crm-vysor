import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class FishAudioProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice")

        if not token:
            raise ValueError("FishAudio provider requires an API token (apiKey).")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout) as client:
            json_payload = {
                "text": text,
                "format": "mp3"
            }
            if voice_id:
                json_payload["reference_id"] = voice_id
                
            response = await client.post(
                "https://api.fish.audio/v1/tts",
                headers={
                    "Authorization": f"Bearer {token}",
                    "Content-Type": "application/json",
                    "model": "s2-pro"
                },
                json=json_payload,
            )

            if not response.is_success:
                error_detail = response.text
                raise Exception(f"FishAudio API error: {response.status_code} - {error_detail}")

            return response.content
