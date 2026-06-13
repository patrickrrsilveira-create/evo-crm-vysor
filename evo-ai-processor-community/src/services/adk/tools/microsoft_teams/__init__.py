"""Microsoft Teams integration tools for AI agents."""

from .check_availability import create_check_teams_availability_tool
from .create_meeting import create_create_teams_meeting_tool

__all__ = [
    "create_check_teams_availability_tool",
    "create_create_teams_meeting_tool",
]
