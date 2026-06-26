import re
import html
import logging
from typing import Optional

logger = logging.getLogger(__name__)

def humanize_text_to_ssml(text: str) -> str:
    """
    Converts plain text into a highly humanized SSML (Speech Synthesis Markup Language)
    format optimized for Brazilian Portuguese (PT-BR) voices, particularly
    the Google Cloud TTS Chirp3 and Neural2 families.

    Args:
        text (str): The plain text to be converted to SSML.

    Returns:
        str: The formatted SSML string wrapped in <speak> and <prosody> tags.
             Returns the original text if it's already an SSML string or on error.
    """
    if not text:
        return ""

    try:
        # 1. Bypass if already SSML
        clean_text = text.strip()
        if clean_text.startswith("<speak>") or clean_text.startswith("<?xml"):
            return clean_text

        # 2. Escape HTML characters to avoid breaking the XML/SSML structure
        # (e.g., & -> &amp;, < -> &lt;, > -> &gt;)
        processed = html.escape(clean_text)

        # 3. Clean up Markdown artifacts (bold, italic, headers)
        processed = re.sub(r'(\*\*|\*|__|_|#+)', '', processed)

        # 4. Insert natural pauses for punctuation
        # Pause after sentence endings (. ! ?)
        processed = re.sub(r'([.!?]+)(?=\s|$)', r'\1 <break time="350ms"/>', processed)
        
        # Pause after commas and semicolons (, ;)
        processed = re.sub(r'([,;]+)(?=\s|$)', r'\1 <break time="200ms"/>', processed)

        # 5. Handle "Perfeito", "Claro", "Certo", "Entendi" at the start of sentences
        # giving them a slightly longer pause for a natural transition
        expressions = r'\b(Perfeito|Claro|Certo|Entendi|Sem problema|Entendido)\b\s*[.,!?;]'
        processed = re.sub(
            expressions,
            lambda m: f'{m.group(0)} <break time="300ms"/>',
            processed,
            flags=re.IGNORECASE
        )

        # 6. Optional: Pause before typical question markers (Qual, O que, Como, etc) 
        # at the end of paragraphs to simulate a human thinking before asking.
        # This is a bit aggressive for general TTS, so we stick to the punctuation breaks which already cover this.

        # 7. Wrap in the universal PT-BR humanizing prosody tags
        # rate="93%": Slows down the slightly fast/rushed native speed of PT-BR models.
        # pitch="-2%": Lowers the pitch slightly to remove the robotic/synthetic high frequencies.
        # volume="+0dB": Baseline volume.
        ssml = (
            '<speak>\n'
            '  <prosody rate="130%" pitch="-2%">\n'
            f'    {processed}\n'
            '  </prosody>\n'
            '</speak>'
        )

        return ssml

    except Exception as e:
        logger.error(f"Error humanizing text to SSML: {str(e)}")
        # In case of any regex or formatting error, fail gracefully returning the raw text
        return text
