"""Microsoft Teams availability checking tool."""

from typing import Optional, Dict, Any
from datetime import datetime, timedelta
from google.adk.tools import FunctionTool, ToolContext
import traceback

from .base import MicrosoftTeamsClient
from src.utils.logger import setup_logger

logger = setup_logger(__name__)


def create_check_teams_availability_tool(
    agent_id: Optional[str] = None,
    teams_config: Optional[Dict[str, Any]] = None,
    credentials_config: Optional[Dict[str, Any]] = None,
    db=None
) -> FunctionTool:
    """
    Create a tool for checking Microsoft Teams availability.

    Args:
        agent_id: Optional default agent ID
        teams_config: Microsoft Teams configuration from agent.config.integrations
        credentials_config: Microsoft Teams credentials from agent.config.integrations
        db: Database session for direct database access

    Returns:
        FunctionTool for checking calendar availability
    """
    client = MicrosoftTeamsClient(db=db)

    async def check_teams_availability(
        start_date: str,
        end_date: str,
        find_slots: bool = False,
        slot_duration: int = 60,
        tool_context: Optional[ToolContext] = None,
    ) -> Dict[str, Any]:
        """
        Check Microsoft Teams (Outlook Calendar) availability for a given time range.

        This tool checks if there are any events scheduled in the specified time range,
        respecting the agent's business hours configuration. It can also find available
        time slots within a date range.

        Use this tool when:
        - A customer asks about available times for a meeting
        - You need to check if a specific time slot is free
        - You want to suggest available time slots to a customer

        Args:
            start_date: Start date/time in ISO format
            end_date: End date/time in ISO format
            find_slots: Whether to return available time slots
            slot_duration: Duration of each slot in minutes
            tool_context: Tool execution context

        Returns:
            Dictionary with availability status and details
        """
        try:
            logger.info(f"Checking Teams availability from {start_date} to {end_date}, find_slots={find_slots}")

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
                return {
                    "status": "error",
                    "message": "Microsoft Teams credentials not configured for this agent"
                }

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

            user_principal_name = config.get("user_principal_name")
            if not user_principal_name:
                return {
                    "status": "error",
                    "message": "User Principal Name (UPN/Email) not configured for this agent."
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

            if find_slots:
                return await _find_available_slots(
                    client,
                    credentials_config,
                    start_dt,
                    end_dt,
                    slot_duration,
                    user_principal_name,
                    config
                )

            result = await client.check_availability(
                credentials_config,
                user_principal_name,
                start_dt,
                end_dt
            )

            if result["status"] == "error":
                return result

            business_hours = get_config_value("businessHours", {})
            within_business_hours = client.is_within_business_hours(start_dt, business_hours)

            timezone = get_config_value("timezone", "America/Sao_Paulo")
            min_advance_time = get_config_value("minAdvanceTime", 0)
            advance_time_ok = client.validate_advance_time(start_dt, min_advance_time, timezone)

            duration_minutes = int((end_dt - start_dt).total_seconds() / 60)
            max_duration = get_config_value("maxDuration", 0)
            duration_ok = client.validate_max_duration(duration_minutes, max_duration)

            is_available = (
                result["available"]
                and within_business_hours
                and advance_time_ok
                and duration_ok
            )

            response = {
                "status": "success",
                "available": is_available,
                "start_time": start_date,
                "end_time": end_date,
                "duration_minutes": duration_minutes,
                "checks": {
                    "no_conflicts": result["available"],
                    "within_business_hours": within_business_hours,
                    "sufficient_advance_time": advance_time_ok,
                    "duration_within_limit": duration_ok
                }
            }

            if is_available:
                response["message"] = f"Time slot is available from {start_date} to {end_date}"
            else:
                reasons = []
                if not result["available"]:
                    reasons.append("there are existing events in this time slot")
                if not within_business_hours:
                    reasons.append("time is outside business hours")
                if not advance_time_ok:
                    reasons.append(f"must be scheduled at least {min_advance_time} hours in advance")
                if not duration_ok:
                    reasons.append(f"duration exceeds maximum of {max_duration} minutes")

                response["message"] = f"Time slot is not available: {', '.join(reasons)}"

            if not result["available"]:
                response["conflicting_events"] = [
                    {
                        "summary": event.get("subject", "Untitled"),
                        "start": event.get("start", {}).get("dateTime"),
                        "end": event.get("end", {}).get("dateTime")
                    }
                    for event in result.get("events", [])
                ]

            return response

        except Exception as e:
            logger.error(f"Unexpected error in check_teams_availability: {str(e)}")
            logger.error(traceback.format_exc())
            return {
                "status": "error",
                "message": f"Failed to check Teams availability: {str(e)}"
            }

    async def _find_available_slots(
        client: MicrosoftTeamsClient,
        credentials_config: Dict[str, Any],
        start_dt: datetime,
        end_dt: datetime,
        slot_duration: int,
        user_principal_name: str,
        config: Dict[str, Any]
    ) -> Dict[str, Any]:
        
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

        available_slots = []
        business_hours = get_config_value("businessHours", {})
        timezone = get_config_value("timezone", "America/Sao_Paulo")
        min_advance_time = get_config_value("minAdvanceTime", 0)
        max_duration = get_config_value("maxDuration", 0)

        if max_duration > 0 and slot_duration > max_duration:
            return {
                "status": "error",
                "message": f"Slot duration ({slot_duration} min) exceeds maximum duration ({max_duration} min)"
            }

        current_day = start_dt.replace(hour=0, minute=0, second=0, microsecond=0)
        end_day = end_dt.replace(hour=23, minute=59, second=59, microsecond=999999)

        while current_day <= end_day:
            day_names = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
            day_name = day_names[current_day.weekday()]
            day_config = business_hours.get(day_name, {}) if business_hours.get("enabled") else {}

            if not day_config or not day_config.get("enabled"):
                current_day += timedelta(days=1)
                continue

            start_time_str = day_config.get("start", "09:00")
            end_time_str = day_config.get("end", "18:00")

            hour, minute = map(int, start_time_str.split(":"))
            day_start = current_day.replace(hour=hour, minute=minute)

            hour, minute = map(int, end_time_str.split(":"))
            day_end = current_day.replace(hour=hour, minute=minute)

            if day_start < start_dt:
                day_start = start_dt
            if day_end > end_dt:
                day_end = end_dt

            slot_start = day_start
            while slot_start + timedelta(minutes=slot_duration) <= day_end:
                slot_end = slot_start + timedelta(minutes=slot_duration)

                if not client.validate_advance_time(slot_start, min_advance_time, timezone):
                    slot_start += timedelta(minutes=30)
                    continue

                result = await client.check_availability(
                    credentials_config,
                    user_principal_name,
                    slot_start,
                    slot_end
                )

                if result["status"] == "success" and result["available"]:
                    available_slots.append({
                        "start": slot_start.isoformat(),
                        "end": slot_end.isoformat(),
                        "duration_minutes": slot_duration
                    })

                slot_start += timedelta(minutes=30)

            current_day += timedelta(days=1)

        return {
            "status": "success",
            "message": f"Found {len(available_slots)} available time slots",
            "available_slots": available_slots,
            "search_range": {
                "start": start_dt.isoformat(),
                "end": end_dt.isoformat()
            },
            "slot_duration_minutes": slot_duration
        }

    check_teams_availability.__name__ = "check_teams_availability"

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

    check_teams_availability.__doc__ = f"""Check Microsoft Teams availability for a time range or find available time slots.

This tool can:
1. Check if a specific time slot is available
2. Find available time slots within a date range
{bh_description}{constraints_description}

IMPORTANT: Always respect the business hours and scheduling constraints above when suggesting meeting times to customers.

Args:
    start_date (str): Start date and time in ISO format (e.g., '2024-01-15T09:00:00') in timezone {timezone}
    end_date (str): End date and time in ISO format (e.g., '2024-01-15T10:00:00') in timezone {timezone}
    find_slots (bool, optional): If True, find available time slots instead of just checking if range is free (default: False)
    slot_duration (int, optional): Duration of each time slot in minutes (used when find_slots=True, default: 60)

Examples:
- Check if 2PM-3PM tomorrow is free: start_date='2024-01-16T14:00:00', end_date='2024-01-16T15:00:00'
- Find available 1-hour slots this week: start_date='2024-01-15T00:00:00', end_date='2024-01-21T23:59:59', find_slots=True, slot_duration=60
"""

    return check_teams_availability
