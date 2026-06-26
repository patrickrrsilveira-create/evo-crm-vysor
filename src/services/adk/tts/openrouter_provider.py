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
            }
            
            # Allow passing language parameter for models like Kokoro
            language = config.get("language") or config.get("lang")
            
            # Default to Portuguese since the assistant mostly speaks Portuguese
            if not language:
                language = "pt"
                
            if language:
                json_payload["language"] = language
            # OpenRouter kokoro sometimes uses 'lang' inside 'extra_body' or directly
            # We'll pass it both ways just in case
            if language:
                json_payload["lang"] = language

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

            content_type = response.headers.get("content-type", "unknown")
            data = response.content
            
            # Use print to bypass logger filtering
            print(f"[OpenRouter TTS] content-type={content_type}, size={len(data)} bytes")
            
            if "application/json" in content_type.lower():
                import json
                import base64
                try:
                    json_data = response.json()
                    print(f"[OpenRouter TTS] JSON Payload: {str(json_data)[:200]}")
                    
                    if "error" in json_data:
                        raise Exception(f"OpenRouter returned error: {json_data['error']}")
                    
                    # Sometimes OpenRouter returns base64 inside "audio" or "data"
                    if "audio" in json_data:
                        data = base64.b64decode(json_data["audio"])
                    elif "data" in json_data and isinstance(json_data["data"], str):
                        data = base64.b64decode(json_data["data"])
                    elif "data" in json_data and isinstance(json_data["data"], list) and len(json_data["data"]) > 0:
                        # OpenAI style JSON wrapper (rare for speech, but possible)
                        first_item = json_data["data"][0]
                        if isinstance(first_item, dict) and "b64_json" in first_item:
                            data = base64.b64decode(first_item["b64_json"])
                    else:
                        raise Exception(f"OpenRouter returned JSON but no recognized audio payload. Keys: {list(json_data.keys())}")
                except Exception as e:
                    print(f"[OpenRouter TTS] JSON Parse Error: {e}")
                    raise Exception(f"Failed to parse OpenRouter JSON: {e}")
            
            # If the response is raw PCM, wrap it in a WAV header
            # OpenRouter Kokoro returns: audio/pcm;rate=24000;channels=1
            if "audio/pcm" in content_type.lower():
                import wave
                import io
                import re
                
                # Extract rate and channels from content_type, or use defaults
                rate = 24000
                channels = 1
                
                rate_match = re.search(r"rate=(\d+)", content_type.lower())
                if rate_match:
                    rate = int(rate_match.group(1))
                    
                ch_match = re.search(r"channels=(\d+)", content_type.lower())
                if ch_match:
                    channels = int(ch_match.group(1))
                
                print(f"[OpenRouter TTS] Wrapping raw PCM into WAV (rate={rate}, channels={channels})")
                wav_io = io.BytesIO()
                with wave.open(wav_io, "wb") as wav_file:
                    wav_file.setnchannels(channels)
                    wav_file.setsampwidth(2) # Assume 16-bit PCM
                    wav_file.setframerate(rate)
                    wav_file.writeframes(data)
                
                data = wav_io.getvalue()

            return data
