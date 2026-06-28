import asyncio
import os
from google.adk.models.lite_llm import LiteLlm
from google.adk.tools import FunctionTool
from dotenv import load_dotenv

load_dotenv(".env")

def test_tool(arg1: str) -> str:
    """Test tool that does nothing."""
    return f"Called with {arg1}"

async def main():
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        print("No OpenAI API key")
        return
        
    model = LiteLlm(model="gpt-4o", api_key=api_key)
    tool = FunctionTool(func=test_tool)
    
    # We need to simulate the ADK runner or just call the model directly
    # Wait, how to call model directly?
    print(dir(model))
    
    # Let's just create an agent and run it
    from google.adk.agents import LlmAgent
    agent = LlmAgent(model=model, tools=[tool], instruction="Use the test_tool with arg1='hello'")
    
    response = await agent.run_async("Do it")
    print(response)

asyncio.run(main())
