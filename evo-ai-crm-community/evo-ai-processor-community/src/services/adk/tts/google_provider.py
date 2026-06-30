import typing
import logging
import base64
import aiohttp

from .base import TTSProvider
from src.utils.ssml_humanizer import humanize_text_to_ssml

logger = logging.getLogger(__name__)

class GoogleProvider(TTSProvider):
    """Google Cloud TTS implementation using the REST API."""
    
    async def generate_speech(self, text: str, config: typing.Dict[str, typing.Any]) -> bytes:
        api_key = config.get("apiKey")
        if not api_key:
            raise ValueError("Google Cloud TTS API key is not configured.")
            
        voice_id = config.get("voice_id") or config.get("voice") or "pt-BR-Neural2-A"
        
        # Extract languageCode from voice_id if possible (e.g. pt-BR from pt-BR-Neural2-A)
        parts = voice_id.split("-")
        language_code = f"{parts[0]}-{parts[1]}" if len(parts) >= 2 else "pt-BR"
            
        url = f"https://texttospeech.googleapis.com/v1/text:synthesize?key={api_key}"
        
        formatted_text = humanize_text_to_ssml(text)
        
        # Determine if we are sending plain text or SSML
        input_payload = {"ssml": formatted_text} if formatted_text.startswith("<speak>") else {"text": formatted_text}
        
        payload = {
            "input": input_payload,
            "voice": {
                "languageCode": language_code,
                "name": voice_id
            },
            "audioConfig": {
                "audioEncoding": "OGG_OPUS"
            }
        }
        
        async with aiohttp.ClientSession() as session:
            try:
                async with session.post(url, json=payload, timeout=aiohttp.ClientTimeout(total=30)) as response:
                    if response.status != 200:
                        error_text = await response.text()
                        logger.error(f"Google TTS API error: {response.status} - {error_text}")
                        raise Exception(f"Google TTS API error: {response.status}")
                        
                    data = await response.json()
                    audio_base64 = data.get("audioContent")
                    
                    if not audio_base64:
                        logger.error(f"Google TTS API response missing audioContent: {data}")
                        raise Exception("Missing audioContent in Google TTS response")
                        
                    return base64.b64decode(audio_base64)
            except Exception as e:
                logger.error(f"Error generating Google TTS speech: {str(e)}")
                raise
