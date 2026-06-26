import os
from typing import AsyncGenerator, Dict, Any, Optional
import base64

from .base_adapter import BaseTTSAdapter
from .humanizer import SSMLHumanizer

try:
    from google.cloud import texttospeech
    import google.auth
    HAS_GOOGLE_TTS = True
except ImportError:
    HAS_GOOGLE_TTS = False

class GoogleTTSAdapter(BaseTTSAdapter):
    def __init__(self):
        if not HAS_GOOGLE_TTS:
            raise RuntimeError("google-cloud-texttospeech is not installed.")
        
        # This will use the default credentials in the environment (GOOGLE_APPLICATION_CREDENTIALS)
        self.client = texttospeech.TextToSpeechAsyncClient()

    async def generate_audio_stream(
        self, 
        text: str, 
        voice_profile: Optional[Dict[str, Any]] = None
    ) -> AsyncGenerator[bytes, None]:
        
        # Prepare the SSML formatted text
        ssml_text = SSMLHumanizer.humanize_text_for_google(text)
        
        synthesis_input = texttospeech.SynthesisInput(ssml=ssml_text)

        # Use profile or default to Brazilian Portuguese Neural2-B (Male) or Neural2-A (Female)
        profile = voice_profile or {}
        language_code = profile.get("language_code", "pt-BR")
        voice_name = profile.get("voice_id", "pt-BR-Neural2-A")  # High quality female by default
        
        # Configure Voice
        voice = texttospeech.VoiceSelectionParams(
            language_code=language_code,
            name=voice_name
        )

        # Configure Audio (We want OGG_OPUS to send directly to WhatsApp / Frontend without re-encoding)
        # Apply standard Pitch and Speaking Rate for realism
        pitch = profile.get("pitch", -1.5)
        speaking_rate = profile.get("speaking_rate", 1.10)
        
        audio_config = texttospeech.AudioConfig(
            audio_encoding=texttospeech.AudioEncoding.OGG_OPUS,
            pitch=pitch,
            speaking_rate=speaking_rate
        )

        # Generate Audio
        response = await self.client.synthesize_speech(
            input=synthesis_input, voice=voice, audio_config=audio_config
        )

        # For Google TTS standard API (not streaming API), it returns the full audio content
        # We yield it as a single chunk. In the future, we can upgrade to streaming synthesis.
        yield response.audio_content
