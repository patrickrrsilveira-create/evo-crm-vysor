import uuid
from google.adk.tools import ToolContext
from google.genai import types
from google.adk.tools import FunctionTool
import logging
from typing import Dict, Any

from src.services.adk.tts.factory import get_tts_provider

logger = logging.getLogger(__name__)

def create_text_to_speech_tool(config: Dict[str, Any]) -> FunctionTool:
    """Create the text_to_speech tool for LoopAgent.

    Args:
        config: The configuration dict, must contain 'provider' (e.g. 'elevenlabs', 'fish')
                and provider-specific settings.
    """
    provider_name = config.get("provider", "elevenlabs")

    async def text_to_speech(
        text: str,
        tool_context: "ToolContext",
    ):
        """Generates speech from text using the configured TTS provider and stores it in artifacts."""
        try:
            if not text or not text.strip():
                return {"status": "failed", "error": "Text cannot be empty"}

            if len(text) > 5000:
                return {
                    "status": "failed",
                    "error": f"Text too long ({len(text)} characters). Maximum allowed is 5000 characters.",
                }

            logger.info(f"Generating speech with {provider_name} for text: {len(text)} chars")

            try:
                provider = get_tts_provider(provider_name)
            except ValueError as e:
                return {"status": "failed", "error": str(e)}

            audio_bytes = await provider.generate_speech(text, config)
            logger.info(f"Received audio response: {len(audio_bytes)} bytes")

            # Generate unique filename
            filename = f"speech_{uuid.uuid4().hex[:8]}.mp3"

            # Create Part and save to artifacts
            audio_blob = types.Blob(mime_type="audio/mpeg", data=audio_bytes)
            audio_part = types.Part(inline_data=audio_blob)
            version = await tool_context.save_artifact(filename, audio_part)

            result = {
                "status": "success",
                "detail": "Speech generated successfully and stored in artifacts.",
                "filename": filename,
                "version": str(version),
            }

            try:
                loaded_artifact = await tool_context.load_artifact(filename, version)
                if loaded_artifact and loaded_artifact.text and loaded_artifact.text.startswith("Artifact URL: "):
                    result["audio_url"] = loaded_artifact.text.replace("Artifact URL: ", "")
                else:
                    result["url_hint"] = loaded_artifact.text if loaded_artifact and loaded_artifact.text else f"Audio file saved as {filename} version {version}."
            except Exception as e:
                result["url_error"] = str(e)

            return result

        except Exception as e:
            logger.error(f"Error in text_to_speech tool ({provider_name}): {str(e)}")
            return {
                "status": "failed",
                "error": f"Text-to-speech error ({provider_name}): {str(e)}",
            }

    text_to_speech.__name__ = "text_to_speech"
    text_to_speech.__doc__ = f"""Generate speech from text using {provider_name}.
    
    Args:
        text: The text to generate speech from
        tool_context: The tool context containing session information
    """

    return text_to_speech
