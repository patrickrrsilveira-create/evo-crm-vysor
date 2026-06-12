import uuid
import subprocess
import tempfile
import os
from google.adk.tools import ToolContext
from google.genai import types
from google.adk.tools import FunctionTool
import logging
from typing import Dict, Any, Optional

from src.services.adk.tts.factory import get_tts_provider

logger = logging.getLogger(__name__)


def _convert_to_ogg_opus(audio_bytes: bytes) -> bytes:
    """Convert audio bytes (MP3/WAV/PCM) to OGG/Opus using ffmpeg.

    Returns the converted bytes, or the original bytes if ffmpeg is unavailable.
    """
    tmp_in_path = ""
    tmp_out_path = ""
    try:
        with tempfile.NamedTemporaryFile(suffix=".mp3", delete=False) as tmp_in:
            tmp_in.write(audio_bytes)
            tmp_in_path = tmp_in.name

        tmp_out_path = tmp_in_path.replace(".mp3", ".ogg")

        result = subprocess.run(
            [
                "ffmpeg", "-y",
                "-i", tmp_in_path,
                "-c:a", "libopus",
                "-b:a", "32k",
                "-ar", "48000",
                "-ac", "1",
                "-application", "voip",
                tmp_out_path,
            ],
            capture_output=True,
            timeout=30,
        )

        if result.returncode == 0 and os.path.exists(tmp_out_path):
            with open(tmp_out_path, "rb") as f:
                converted = f.read()
            logger.info(
                f"[TTS] Converted {len(audio_bytes)} bytes MP3 -> "
                f"{len(converted)} bytes OGG/Opus"
            )
            return converted
        else:
            stderr_snippet = result.stderr[:200] if result.stderr else b"(no stderr)"
            logger.warning(
                f"[TTS] ffmpeg conversion failed (rc={result.returncode}): "
                f"{stderr_snippet}"
            )
            return audio_bytes

    except FileNotFoundError:
        logger.warning("[TTS] ffmpeg not found, returning raw audio bytes (MP3)")
        return audio_bytes
    except Exception as e:
        logger.warning(f"[TTS] Conversion error: {e}, returning raw audio bytes")
        return audio_bytes
    finally:
        # Cleanup temp files
        for path in [tmp_in_path, tmp_out_path]:
            try:
                if path and os.path.exists(path):
                    os.unlink(path)
            except OSError:
                pass


def create_text_to_speech_tool(config: Dict[str, Any]) -> FunctionTool:
    """Create the text_to_speech tool for LoopAgent.

    Args:
        config: The configuration dict, must contain 'provider' (e.g. 'elevenlabs', 'fish')
                and provider-specific settings.
    """
    provider_name = config.get("provider", "elevenlabs")

    async def text_to_speech(
        text: str,
        tool_context: Optional[ToolContext] = None,
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

            # Convert to OGG/Opus for WhatsApp native voice message compatibility
            audio_bytes = _convert_to_ogg_opus(audio_bytes)

            # Generate unique filename
            filename = f"speech_{uuid.uuid4().hex[:8]}.ogg"

            # Create Part and save to artifacts
            audio_blob = types.Blob(mime_type="audio/ogg", data=audio_bytes)
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
    CRITICAL: You MUST use this tool to answer if the user requires audio. DO NOT answer with text alone.
    
    Args:
        text: The exact text of your response to generate speech from.
    """

    return FunctionTool(func=text_to_speech)
