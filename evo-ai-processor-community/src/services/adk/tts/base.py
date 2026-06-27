from abc import ABC, abstractmethod
import typing

class TTSProvider(ABC):
    """Base class for all Text-to-Speech providers."""
    
    @abstractmethod
    async def generate_speech(self, text: str, config: typing.Dict[str, typing.Any]) -> bytes:
        """
        Generate audio bytes from text.
        
        Args:
            text: The text to convert to speech.
            config: A dictionary containing provider-specific configuration.
            
        Returns:
            The raw audio bytes (ogg/opus format for WhatsApp compatibility).
        """
        pass
