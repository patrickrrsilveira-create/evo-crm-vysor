from typing import AsyncGenerator, Dict, Any, Optional

from .base_adapter import BaseTTSAdapter
from .google_adapter import GoogleTTSAdapter
from .elevenlabs_adapter import ElevenLabsAdapter

class VoiceService:
    """
    The Maestro (Factory) for Text-To-Speech generation.
    Routes requests to the correct adapter based on the Voice Profile.
    """

    def __init__(self):
        # We instantiate adapters lazily or inject them.
        # For now, we will hold references to the classes.
        self.adapters = {
            "google": GoogleTTSAdapter,
            "elevenlabs": ElevenLabsAdapter,
        }
        
        self._instances = {}

    def _get_adapter(self, provider: str) -> BaseTTSAdapter:
        if provider not in self.adapters:
            raise ValueError(f"TTS Provider '{provider}' is not supported.")
            
        if provider not in self._instances:
            self._instances[provider] = self.adapters[provider]()
            
        return self._instances[provider]

    async def stream_audio(
        self, 
        text: str, 
        voice_profile: Optional[Dict[str, Any]] = None
    ) -> AsyncGenerator[bytes, None]:
        """
        Generates and streams audio for a given text using the configured provider.
        """
        profile = voice_profile or {}
        # Default to Google if not specified, but could be 'elevenlabs', etc.
        provider = profile.get("provider", "google").lower()
        
        adapter = self._get_adapter(provider)
        
        async for chunk in adapter.generate_audio_stream(text, profile):
            yield chunk
