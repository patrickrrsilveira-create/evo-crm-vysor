import re

class SSMLHumanizer:
    """
    Parses "dirty" LLM text (with hesitations, ellipses, conversational markers)
    and converts it into SSML prosody and breaks for TTS engines that support it.
    """

    @staticmethod
    def humanize_text_for_google(text: str) -> str:
        """
        Converts text with conversational cues into Google TTS SSML.
        """
        # Replace ellipses with short pauses
        text = text.replace("...", '<break time="300ms"/>')
        
        # Replace multiple commas with medium pauses (if LLM is hesitant)
        text = re.sub(r',\\s*,', '<break time="400ms"/>', text)
        
        # Common hesitation words - add a tiny break after them to sound like thinking
        hesitations = ["hmmm", "hmmm,", "ahn", "ahn,", "ééé", "ééé,"]
        for h in hesitations:
            # Case insensitive replace but keep original casing is hard without re.sub, 
            # so we'll do a simple regex substitution.
            text = re.sub(f'(?i)\\b{h}\\b', f'{h} <break time="200ms"/>', text)
            
        # Clean any stray markdown formatting that the LLM might have ignored instructions on
        text = text.replace("*", "").replace("#", "")

        # Wrap in SSML tags with adjusted pitch and rate for natural speech
        # Note: pitch and rate are also handled at the synthesis config level, 
        # but we can do it here for specific phrases if needed.
        return f"<speak>{text}</speak>"

    @staticmethod
    def humanize_text_for_elevenlabs(text: str) -> str:
        """
        ElevenLabs relies on punctuation rather than SSML.
        """
        # Ensure ellipses have spaces around them to force a natural sigh/pause
        text = re.sub(r'\\.{3}', ' ... ', text)
        
        # Clean markdown
        text = text.replace("*", "").replace("#", "")
        return text
