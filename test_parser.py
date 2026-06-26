import re
import logging
from typing import Optional, Dict

logger = logging.getLogger(__name__)

def extract_tool_call_from_text(text: str) -> tuple[str, Optional[Dict[str, str]]]:
    """
    Extracts LiteLlm tool call syntax from text.
    Format: <tool_call:transfer_conversation{reason:<|"|>...<|"|>,target_agent_id:<|"|>...<|"|>}>
    Returns: (cleaned_text, tool_call_dict)
    """
    pattern = r'<tool_call:transfer_conversation\{(.*?)\}>'
    match = re.search(pattern, text)
    
    if not match:
        return text, None
        
    cleaned_text = re.sub(pattern, '', text).strip()
    args_str = match.group(1)
    
    # Parse args_str: reason:<|"|>...<|"|>,target_agent_id:<|"|>...<|"|>
    args = {}
    
    # Split by , but we need to be careful about commas inside the values
    # Actually, let's just use regex to extract key:<|"|>value<|"|>
    arg_pattern = r'(\w+):<\|"\|>(.*?)<\|"\|>'
    for arg_match in re.finditer(arg_pattern, args_str):
        args[arg_match.group(1)] = arg_match.group(2)
        
    # Fallback if the delimiter is just normal quotes
    if not args:
        arg_pattern2 = r'(\w+):"(.*?)"'
        for arg_match in re.finditer(arg_pattern2, args_str):
            args[arg_match.group(1)] = arg_match.group(2)
            
    return cleaned_text, args

# Test
text = 'Vou encaminhar você agora para o Especialista Ganader. <tool_call:transfer_conversation{reason:<|"|>Cliente demonstrou interesse no Ganader.<|"|>,target_agent_id:<|"|>Especialista Ganader<|"|>}>'
cleaned, args = extract_tool_call_from_text(text)
print(f"Cleaned: {cleaned}")
print(f"Args: {args}")
