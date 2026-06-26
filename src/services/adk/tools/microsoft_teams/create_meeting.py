"""Microsoft Teams meeting creation tool."""

from typing import Optional, Dict, Any, List
from datetime import datetime
from google.adk.tools import FunctionTool, ToolContext
import traceback

from .base import MicrosoftTeamsClient
from src.utils.logger import setup_logger

logger = setup_logger(__name__)


def create_create_teams_meeting_tool(
    agent_id: Optional[str] = None,
    teams_config: Optional[Dict[str, Any]] = None,
    credentials_config: Optional[Dict[str, Any]] = None,
    db=None
) -> FunctionTool:
    """
    Create a tool for creating Microsoft Teams meetings.

    Args:
        agent_id: Optional default agent ID
        teams_config: Microsoft Teams configuration from agent.config.integrations
        credentials_config: Microsoft Teams credentials from agent.config.integrations
        db: Database session for direct database access

    Returns:
        FunctionTool for creating Microsoft Teams meetings
    """
    client = MicrosoftTeamsClient(db=db)

    async def create_teams_meeting(
        summary: str,
        start_date: str,
        end_date: str,
        description: str = "",
        attendees: Optional[List[str]] = None,
        tool_context: Optional[ToolContext] = None,
    ) -> Dict[str, Any]:
        """
        Create a new Microsoft Teams meeting.

        This tool creates a meeting in the configured agent's Outlook calendar
        and generates a Microsoft Teams meeting link. It respects the agent's business
        hours and scheduling constraints.

        Use this tool when:
        - A customer agrees to a specific time slot for a meeting
        - You need to schedule an appointment with a customer
        - The customer asks to book a meeting

        IMPORTANT: Always check availability first before creating a meeting to avoid conflicts.

        Args:
            summary: Short title for the meeting
            start_date: Start date/time in ISO format
            end_date: End date/time in ISO format
            description: Detailed description or agenda for the meeting
            attendees: List of email addresses of people to invite
            tool_context: Tool execution context

        Returns:
            Dictionary with creation result, including the Teams meeting link
        """
        try:
            logger.info(f"Creating Teams meeting: '{summary}' from {start_date} to {end_date}")

            effective_agent_id = agent_id

            if not effective_agent_id:
                return {
                    "status": "error",
                    "message": "Agent ID is required but was not provided"
                }

            if not teams_config:
                return {
                    "status": "error",
                    "message": "Microsoft Teams integration not configured for this agent"
                }

            if not credentials_config:
                # We don't strictly need credentials_config for Webhook
                pass

            try:
                start_dt = datetime.fromisoformat(start_date.replace('Z', '+00:00'))
                end_dt = datetime.fromisoformat(end_date.replace('Z', '+00:00'))
            except ValueError as e:
                return {
                    "status": "error",
                    "message": f"Invalid date format: {str(e)}. Use ISO format like '2024-01-15T09:00:00'"
                }

            if end_dt <= start_dt:
                return {
                    "status": "error",
                    "message": "End date must be after start date"
                }

            if "settings" in teams_config:
                config = teams_config["settings"]
            else:
                config = teams_config
                
            webhook_url = config.get("webhookUrl")
            if not webhook_url:
                return {
                    "status": "error",
                    "message": "URL do Webhook do n8n não configurada para este agente."
                }

            def get_config_value(key: str, default: Any) -> Any:
                value = config.get(key, default)
                if isinstance(value, dict) and "value" in value:
                    extracted_value = value["value"]
                    unit = value.get("unit")

                    if key in ["minAdvanceTime", "maxDistance"] and unit == "hours":
                        return extracted_value
                    elif key in ["minAdvanceTime", "maxDistance"] and unit == "weeks":
                        return extracted_value * 24 * 7
                    elif key == "maxDuration" and unit == "hours":
                        return extracted_value * 60
                    elif key == "maxDuration" and unit == "minutes":
                        return extracted_value

                    return extracted_value
                return value

            business_hours = get_config_value("businessHours", {})
            within_business_hours = client.is_within_business_hours(start_dt, business_hours)

            timezone = get_config_value("timezone", "America/Sao_Paulo")
            min_advance_time = get_config_value("minAdvanceTime", 0)
            advance_time_ok = client.validate_advance_time(start_dt, min_advance_time, timezone)

            duration_minutes = int((end_dt - start_dt).total_seconds() / 60)
            max_duration = get_config_value("maxDuration", 0)
            duration_ok = client.validate_max_duration(duration_minutes, max_duration)

            if not (within_business_hours and advance_time_ok and duration_ok):
                reasons = []
                if not within_business_hours:
                    reasons.append("time is outside business hours")
                if not advance_time_ok:
                    reasons.append(f"must be scheduled at least {min_advance_time} hours in advance")
                if not duration_ok:
                    reasons.append(f"duration exceeds maximum of {max_duration} minutes")

                return {
                    "status": "error",
                    "message": f"Cannot create meeting: {', '.join(reasons)}",
                    "checks": {
                        "within_business_hours": within_business_hours,
                        "sufficient_advance_time": advance_time_ok,
                        "duration_within_limit": duration_ok
                    }
                }

            result = await client.create_meeting(
                credentials_config,
                config,
                summary,
                start_dt,
                end_dt,
                description,
                attendees
            )

            if result["status"] == "success":
                logger.info(f"Successfully created Teams meeting '{summary}'")
                return {
                    "status": "success",
                    "message": "Meeting created successfully",
                    "meeting_link": result.get("meeting_link"),
                    "event_id": result.get("event_id"),
                    "start_time": result.get("start"),
                    "end_time": result.get("end")
                }
            else:
                logger.error(f"Failed to create Teams meeting: {result.get('message')}")
                return result

        except Exception as e:
            logger.error(f"Unexpected error in create_teams_meeting: {str(e)}")
            logger.error(traceback.format_exc())
            return {
                "status": "error",
                "message": f"Failed to create meeting: {str(e)}"
            }

    create_teams_meeting.__name__ = "create_teams_meeting"

    config = teams_config.get("settings", teams_config) if teams_config else {}

    def get_config_value(key: str, default: Any) -> Any:
        value = config.get(key, default)
        if isinstance(value, dict) and "value" in value:
            extracted_value = value["value"]
            unit = value.get("unit")
            if key in ["minAdvanceTime", "maxDistance"] and unit == "hours":
                return extracted_value
            elif key in ["minAdvanceTime", "maxDistance"] and unit == "weeks":
                return extracted_value * 24 * 7
            elif key == "maxDuration" and unit == "hours":
                return extracted_value * 60
            elif key == "maxDuration" and unit == "minutes":
                return extracted_value
            return extracted_value
        return value

    business_hours = get_config_value("businessHours", {})
    min_advance_time = get_config_value("minAdvanceTime", 0)
    max_duration = get_config_value("maxDuration", 0)
    timezone = get_config_value("timezone", "America/Sao_Paulo")

    bh_description = ""
    if business_hours and business_hours.get("enabled"):
        bh_description = "\n\nBUSINESS HOURS CONFIGURED:\n"
        day_names = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
        for day_name in day_names:
            day_config = business_hours.get(day_name, {})
            if day_config and day_config.get("enabled"):
                start = day_config.get("start", "09:00")
                end = day_config.get("end", "18:00")
                bh_description += f"- {day_name.capitalize()}: {start} to {end}\n"

    constraints = []
    if min_advance_time > 0:
        constraints.append(f"- Meetings must be scheduled at least {min_advance_time} hours in advance")
    if max_duration > 0:
        constraints.append(f"- Maximum meeting duration: {max_duration} minutes")

    constraints_description = ""
    if constraints:
        constraints_description = "\n\nSCHEDULING CONSTRAINTS:\n" + "\n".join(constraints)

    create_teams_meeting.__doc__ = f"""Create a new Microsoft Teams meeting.

This tool creates a meeting in the configured agent's Outlook calendar and generates a Teams meeting link.
{bh_description}{constraints_description}

IMPORTANT:
1. ALWAYS check availability first using check_teams_availability before creating a meeting.
2. ALWAYS ask the customer for their email address if you don't have it, to add them to the attendees list.

Args:
    summary (str): Short title for the meeting (e.g., 'Consultation with John Doe')
    start_date (str): Start date and time in ISO format (e.g., '2024-01-15T09:00:00') in timezone {timezone}
    end_date (str): End date and time in ISO format (e.g., '2024-01-15T10:00:00') in timezone {timezone}
    description (str, optional): Detailed description or agenda for the meeting
    attendees (list, optional): List of email addresses of people to invite

Examples:
- Schedule a 1-hour meeting: summary='Product Demo', start_date='2024-01-16T14:00:00', end_date='2024-01-16T15:00:00', attendees=['customer@example.com']
"""

    return create_teams_meeting
