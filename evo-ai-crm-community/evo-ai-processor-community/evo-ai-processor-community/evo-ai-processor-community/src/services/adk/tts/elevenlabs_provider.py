import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class ElevenLabsProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("eleven_labs_token") or config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice")
        model_id = config.get("model_id", "eleven_multilingual_v2")
        stability = config.get("stability", 80) / 100.0
        similarity_boost = config.get("similarity_boost", 50) / 100.0
        style = config.get("style", 0.0)
        use_speaker_boost = config.get("use_speaker_boost", True)

        if not token or not voice_id:
            raise ValueError("ElevenLabs provider requires an API token and a voice_id.")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout) as client:
            response = await client.post(
                f"https://api.elevenlabs.io/v1/text-to-speech/{voice_id}",
                headers={"xi-api-key": token},
                json={
                    "text": text,
                    "model_id": model_id,
                    "voice_settings": {
                        "stability": stability,
                        "similarity_boost": similarity_boost,
                        "style": style,
                        "use_speaker_boost": use_speaker_boost,
                    },
                },
            )

            if not response.is_success:
                error_detail = response.text
                raise Exception(f"ElevenLabs API error: {response.status_code} - {error_detail}")

            return response.content
