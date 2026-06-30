import typing
import logging
import base64
import asyncio
import re
import time

import aiohttp

from .base import TTSProvider
from src.utils.ssml_humanizer import humanize_text_to_ssml

logger = logging.getLogger(__name__)

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------
_MAX_CONCURRENT_CHUNKS = 4          # Semaphore cap to avoid Google rate-limits
_MIN_CHUNK_LENGTH = 50              # Merge chunks shorter than this
_MAX_CHUNK_LENGTH = 300             # Subdivide chunks longer than this (by comma)
_SENTENCE_BOUNDARY = re.compile(
    r'(?<=[.!?…])\s+|(?<=\n)\s*',   # Split on ". " / "! " / "? " / "…" / newlines
    flags=re.UNICODE,
)
_COMMA_BOUNDARY = re.compile(
    r'(?<=[,;:])\s+',               # Fallback split for oversized chunks
    flags=re.UNICODE,
)
_API_TIMEOUT = aiohttp.ClientTimeout(total=25)
_GOOGLE_TTS_URL = "https://texttospeech.googleapis.com/v1/text:synthesize"


class GoogleProvider(TTSProvider):
    """High-performance Google Cloud TTS provider.

    Uses **Parallel Sentence Chunking** to drastically reduce perceived latency:
    1. Splits the input text into sentence-sized chunks.
    2. Fires all chunks concurrently (capped by a semaphore).
    3. Reassembles the resulting OGG/Opus byte streams in order.

    For short inputs (single sentence), falls back to a single request with
    zero overhead — identical to the previous implementation.
    """

    # ------------------------------------------------------------------
    # Public interface (TTSProvider contract)
    # ------------------------------------------------------------------
    async def generate_speech(
        self,
        text: str,
        config: typing.Dict[str, typing.Any],
    ) -> bytes:
        api_key = config.get("apiKey")
        if not api_key:
            raise ValueError("Google Cloud TTS API key is not configured.")

        voice_id = config.get("voice_id") or config.get("voice") or "pt-BR-Neural2-A"
        parts = voice_id.split("-")
        language_code = f"{parts[0]}-{parts[1]}" if len(parts) >= 2 else "pt-BR"

        url = f"{_GOOGLE_TTS_URL}?key={api_key}"

        chunks = self._split_into_chunks(text)
        t0 = time.monotonic()

        if len(chunks) == 1:
            # ---- Fast path: single chunk, zero overhead ----
            audio = await self._synthesize_chunk(
                url=url,
                text=chunks[0],
                voice_id=voice_id,
                language_code=language_code,
            )
            elapsed = time.monotonic() - t0
            logger.info(
                f"[GoogleTTS] Single-chunk synthesis completed in {elapsed:.2f}s "
                f"({len(chunks[0])} chars → {len(audio)} bytes)"
            )
            return audio

        # ---- Parallel path: fire all chunks concurrently ----
        semaphore = asyncio.Semaphore(_MAX_CONCURRENT_CHUNKS)

        async def _guarded_synthesize(idx: int, chunk: str) -> typing.Tuple[int, bytes]:
            async with semaphore:
                t_chunk = time.monotonic()
                audio = await self._synthesize_chunk(
                    url=url,
                    text=chunk,
                    voice_id=voice_id,
                    language_code=language_code,
                )
                logger.info(
                    f"[GoogleTTS] Chunk {idx+1}/{len(chunks)} done in "
                    f"{time.monotonic() - t_chunk:.2f}s "
                    f"({len(chunk)} chars → {len(audio)} bytes)"
                )
                return idx, audio

        results = await asyncio.gather(
            *[_guarded_synthesize(i, c) for i, c in enumerate(chunks)],
            return_exceptions=True,
        )

        # Check for failures — if ANY chunk fails, raise immediately
        for r in results:
            if isinstance(r, BaseException):
                raise r

        # Sort by original index and concatenate OGG byte streams
        results.sort(key=lambda x: x[0])
        combined = b"".join(audio for _, audio in results)

        elapsed = time.monotonic() - t0
        logger.info(
            f"[GoogleTTS] Parallel synthesis completed in {elapsed:.2f}s "
            f"({len(chunks)} chunks, {len(text)} chars → {len(combined)} bytes)"
        )
        return combined

    # ------------------------------------------------------------------
    # Single-chunk REST call (reusable for both paths)
    # ------------------------------------------------------------------
    async def _synthesize_chunk(
        self,
        url: str,
        text: str,
        voice_id: str,
        language_code: str,
    ) -> bytes:
        """Send a single text segment to Google Cloud TTS and return raw audio bytes."""
        formatted = humanize_text_to_ssml(text)
        input_payload = (
            {"ssml": formatted}
            if formatted.startswith("<speak>")
            else {"text": formatted}
        )

        payload = {
            "input": input_payload,
            "voice": {
                "languageCode": language_code,
                "name": voice_id,
            },
            "audioConfig": {
                "audioEncoding": "OGG_OPUS",
            },
        }

        async with aiohttp.ClientSession() as session:
            async with session.post(
                url, json=payload, timeout=_API_TIMEOUT
            ) as response:
                if response.status != 200:
                    error_text = await response.text()
                    logger.error(
                        f"[GoogleTTS] API error {response.status}: {error_text[:500]}"
                    )
                    raise Exception(
                        f"Google TTS API error: {response.status}"
                    )

                data = await response.json()
                audio_b64 = data.get("audioContent")

                if not audio_b64:
                    logger.error(
                        f"[GoogleTTS] Response missing audioContent: "
                        f"{str(data)[:300]}"
                    )
                    raise Exception(
                        "Missing audioContent in Google TTS response"
                    )

                return base64.b64decode(audio_b64)

    # ------------------------------------------------------------------
    # Text chunking engine
    # ------------------------------------------------------------------
    @staticmethod
    def _split_into_chunks(text: str) -> typing.List[str]:
        """Split text into sentence-sized chunks optimized for parallel TTS.

        Rules:
        1. Split on sentence boundaries (. ! ? … and newlines).
        2. Merge fragments shorter than _MIN_CHUNK_LENGTH with the previous chunk
           to prevent robotic-sounding micro-pauses.
        3. Subdivide chunks longer than _MAX_CHUNK_LENGTH at commas/semicolons.
        4. If a single sentence exceeds _MAX_CHUNK_LENGTH even after comma-split,
           keep it as-is (Google handles up to 5000 chars).
        """
        if not text or not text.strip():
            return [text or ""]

        # Step 1: Split on sentence boundaries
        raw_parts = _SENTENCE_BOUNDARY.split(text.strip())
        raw_parts = [p.strip() for p in raw_parts if p and p.strip()]

        if not raw_parts:
            return [text.strip()]

        # Step 2: Merge short fragments with previous chunk
        merged: typing.List[str] = []
        for part in raw_parts:
            if merged and len(part) < _MIN_CHUNK_LENGTH:
                merged[-1] = f"{merged[-1]} {part}"
            else:
                merged.append(part)

        # Step 3: Subdivide oversized chunks at commas
        final: typing.List[str] = []
        for chunk in merged:
            if len(chunk) <= _MAX_CHUNK_LENGTH:
                final.append(chunk)
                continue

            # Try splitting at commas/semicolons
            sub_parts = _COMMA_BOUNDARY.split(chunk)
            sub_parts = [s.strip() for s in sub_parts if s and s.strip()]

            if len(sub_parts) <= 1:
                # No commas found — keep the long chunk as-is
                final.append(chunk)
                continue

            # Re-merge comma-split parts that are too small
            buffer = ""
            for sp in sub_parts:
                candidate = f"{buffer} {sp}".strip() if buffer else sp
                if len(candidate) > _MAX_CHUNK_LENGTH and buffer:
                    final.append(buffer)
                    buffer = sp
                else:
                    buffer = candidate
            if buffer:
                final.append(buffer)

        return final if final else [text.strip()]
