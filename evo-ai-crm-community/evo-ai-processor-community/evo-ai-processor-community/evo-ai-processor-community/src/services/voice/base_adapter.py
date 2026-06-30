from abc import ABC, abstractmethod
from typing import AsyncGenerator, Dict, Any, Optional

class BaseTTSAdapter(ABC):
    """
    Abstract base class for all Text-to-Speech (TTS) providers.
    Ensures a unified interface for the VoiceService.
    """

    @abstractmethod
    async def generate_audio_stream(
        self, 
        text: str, 
        voice_profile: Optional[Dict[str, Any]] = None
    ) -> AsyncGenerator[bytes, None]:
        """
        Generates audio chunks for the given text.

        Args:
            text (str): The text to synthesize (after humanizer processing).
            voice_profile (dict, optional): Provider-specific voice configuration 
                                          (e.g., voice_id, pitch, rate, language).
        
        Yields:
            bytes: Audio chunks as they are generated.
        """
        pass
        
        # We must have an explicit yield to make it an async generator typing-wise
        yield b""
