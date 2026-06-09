import httpx
import logging
from typing import Dict, Any
from .base import TTSProvider

logger = logging.getLogger(__name__)

class KokoroProvider(TTSProvider):
    async def generate_speech(self, text: str, config: Dict[str, Any]) -> bytes:
        token = config.get("apiKey")
        voice_id = config.get("voice_id") or config.get("voice") or "af_heart"

        # Priority 1: separate api_url field sent directly by the frontend
        api_url = config.get("api_url") or ""

        # Priority 2: pipe-encoded format inside apiKey field ("api_url|api_key")
        if token and "|" in token:
            api_url, token = token.split("|", 1)

        # Priority 3: fall back to default Kokoro endpoint
        if not api_url:
            api_url = "https://api.kokoro.io/v1/audio/speech"

        if not token:
            raise ValueError("Kokoro provider requires an API token (apiKey).")

        # Auto-select model based on endpoint
        model_name = "kokoro-v1"
        if "openrouter" in api_url.lower():
            model_name = "hexgrad/kokoro-82m"

        logger.info(f"[Kokoro TTS] url={api_url} | model={model_name} | voice={voice_id}")

        timeout = httpx.Timeout(60.0)
        async with httpx.AsyncClient(timeout=timeout, verify=False) as client:
            json_payload = {
                "input": text,
                "voice": voice_id,
                "model": model_name,
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
                raise Exception(f"Kokoro API error: {response.status_code} - {error_detail}")

            return response.content
