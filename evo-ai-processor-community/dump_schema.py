import sys
import json
sys.path.insert(0, '/app')
from src.services.adk.tools.send_agent_media import create_send_agent_media_tool
import litellm.utils

tool = create_send_agent_media_tool('123')
schema = litellm.utils.function_to_dict(tool.func)
print(json.dumps(schema, indent=2))
