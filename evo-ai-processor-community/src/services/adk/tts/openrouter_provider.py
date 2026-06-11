import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class OpenRouterProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice") or "af_heart"
        api_url = config.get("api_url") or "https://openrouter.ai/api/v1/audio/speech"
        model_name = config.get("model") or "hexgrad/kokoro-82m"

        api_url = api_url.strip() if api_url else ""
        token = token.strip() if token else ""
        model_name = model_name.strip() if model_name else ""

        if not token:
            raise ValueError("OpenRouter provider requires an API token (apiKey).")

        logger.info(f"[OpenRouter TTS] url={api_url} | model={model_name} | voice={voice_id}")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout, verify=False) as client:
            json_payload = {
                "input": text,
                "voice": voice_id,
                "model": model_name,
                "response_format": "mp3"
            }

            response = await client.post(
                api_url,
                headers={
                    "Authorization": f"Bearer {token}",
                    "Content-Type": "application/json",
                },
                json=json_payload,
            )

            if not response.is_success:
                error_detail = response.text
                raise Exception(f"OpenRouter API error: {response.status_code} - {error_detail}")

            return response.content
