from typing import Dict, Any, Type
from .base import TTSProvider
from .elevenlabs_provider import ElevenLabsProvider
from .fish_audio_provider import FishAudioProvider
from .cartesia_provider import CartesiaProvider
from .kokoro_provider import KokoroProvider
from .voxtral_provider import VoxtralProvider
from .openrouter_provider import OpenRouterProvider

from .google_provider import GoogleProvider

_PROVIDERS: Dict[str, Type[TTSProvider]] = {
    "elevenlabs": ElevenLabsProvider,
    "fishaudio": FishAudioProvider,
    "fish": FishAudioProvider,
    "cartesia": CartesiaProvider,
    "kokoro": KokoroProvider,
    "voxtral": VoxtralProvider,
    "openrouter": OpenRouterProvider,
    "google": GoogleProvider,
}

def get_tts_provider(provider_name: str) -> TTSProvider:
    """
    Factory function to get a TTS provider instance by name.
    
    Args:
        provider_name: The name of the TTS provider (e.g. 'elevenlabs', 'fishaudio')
        
    Returns:
        An instance of TTSProvider
        
    Raises:
        ValueError: If the provider is not supported.
    """
    provider_class = _PROVIDERS.get(provider_name.lower())
    if not provider_class:
        raise ValueError(f"Unsupported TTS provider: {provider_name}")
    return provider_class()
