"""Google Calendar availability checking tool."""

from typing import Optional, Dict, Any
from datetime import datetime, timedelta
from google.adk.tools import FunctionTool, ToolContext
import traceback

from .base import GoogleCalendarClient
from src.utils.logger import setup_logger

logger = setup_logger(__name__)


def create_check_availability_tool(
    agent_id: Optional[str] = None,
    calendar_config: Optional[Dict[str, Any]] = None,
    credentials_config: Optional[Dict[str, Any]] = None,
    db=None
) -> FunctionTool:
    """
    Create a tool for checking Google Calendar availability.

    Args:
        agent_id: Optional default agent ID
        calendar_config: Google Calendar configuration from agent.config.integrations
        credentials_config: Google Calendar credentials from agent.config.integrations
        db: Database session for direct database access

    Returns:
        FunctionTool for checking calendar availability
    """
    client = GoogleCalendarClient(db=db)

    async def check_calendar_availability(
        start_date: str,
        end_date: str,
        calendar_id: str = "primary",
        find_slots: bool = False,
        slot_duration: int = 60,
        tool_context: Optional[ToolContext] = None,
    ) -> Dict[str, Any]:
        """
        Check Google Calendar availability for a given time range.

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
            calendar_id: Which calendar to check
            find_slots: Whether to return available time slots
            slot_duration: Duration of each slot in minutes
            tool_context: Tool execution context

        Returns:
            Dictionary with availability status and details
        """
        try:
            logger.info(f"Checking calendar availability from {start_date} to {end_date}, find_slots={find_slots}")
            logger.info(f"[CalendarConfig] Full config keys: {list(calendar_config.keys()) if calendar_config else 'None'}")
            logger.info(f"[CalendarConfig] businessHours raw: {(calendar_config.get('settings') or calendar_config).get('businessHours', 'NOT FOUND')}")
            logger.debug(f"Calendar config received: {calendar_config}")
            logger.debug(f"Credentials available: {bool(credentials_config)}")

            # Use agent_id from closure (passed to create_check_availability_tool)
            effective_agent_id = agent_id

            # Validate required parameters
            if not effective_agent_id:
                return {
                    "status": "error",
                    "message": "Agent ID is required but was not provided"
                }

            # Validate configs provided
            if not calendar_config:
                return {
                    "status": "error",
                    "message": "Google Calendar integration not configured for this agent"
                }

            if not credentials_config:
                return {
                    "status": "error",
                    "message": "Google Calendar credentials not configured for this agent"
                }

            # Parse dates
            try:
                start_dt = datetime.fromisoformat(start_date.replace('Z', '+00:00'))
                end_dt = datetime.fromisoformat(end_date.replace('Z', '+00:00'))
            except ValueError as e:
                return {
                    "status": "error",
                    "message": f"Invalid date format: {str(e)}. Use ISO format like '2024-01-15T09:00:00'"
                }

            # Validate date range
            if end_dt <= start_dt:
                return {
                    "status": "error",
                    "message": "End date must be after start date"
                }

            # Use configs from closure (passed from agent.config.integrations)
            # Support both flat and nested structures
            if "settings" in calendar_config:
                config = calendar_config["settings"]
            else:
                # Config values are directly in calendar_config
                config = calendar_config

            # Helper to extract value from config (handles both dict and direct values)
            def get_config_value(key: str, default: Any) -> Any:
                value = config.get(key, default)
                # If value is a dict with 'value' key, extract it and convert units
                if isinstance(value, dict) and "value" in value:
                    if "enabled" in value and not value["enabled"]:
                        return default

                    extracted_value = value["value"]
                    unit = value.get("unit")

                    # Convert time units to appropriate format
                    if key in ["minAdvanceTime", "maxDistance"] and unit == "hours":
                        return extracted_value  # Already in hours
                    elif key in ["minAdvanceTime", "maxDistance"] and unit == "weeks":
                        return extracted_value * 24 * 7  # Convert weeks to hours
                    elif key == "maxDuration" and unit == "hours":
                        return extracted_value * 60  # Convert hours to minutes
                    elif key == "maxDuration" and unit == "minutes":
                        return extracted_value  # Already in minutes

                    return extracted_value
                return value

            # If find_slots is True, find available slots
            if find_slots:
                return await _find_available_slots(
                    client,
                    credentials_config,
                    start_dt,
                    end_dt,
                    slot_duration,
                    calendar_id,
                    config
                )

            # Otherwise, just check if the specific range is available
            result = await client.check_availability(
                credentials_config,
                start_dt,
                end_dt,
                calendar_id
            )

            if result["status"] == "error":
                return result

            # Check business hours
            business_hours = get_config_value("businessHours", {})
            within_business_hours = client.is_within_business_hours(start_dt, business_hours)

            # Check minimum advance time
            timezone = get_config_value("timezone", "America/Sao_Paulo")
            min_advance_time = get_config_value("minAdvanceTime", 0)
            advance_time_ok = client.validate_advance_time(start_dt, min_advance_time, timezone)

            # Check maximum duration
            duration_minutes = int((end_dt - start_dt).total_seconds() / 60)
            max_duration = get_config_value("maxDuration", 0)
            duration_ok = client.validate_max_duration(duration_minutes, max_duration)

            # Build response
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

            # Add human-readable message
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

            # Add conflicting events if any
            if not result["available"]:
                response["conflicting_events"] = [
                    {
                        "summary": event.get("summary", "Untitled"),
                        "start": event.get("start", {}).get("dateTime"),
                        "end": event.get("end", {}).get("dateTime")
                    }
                    for event in result.get("events", [])
                ]

            logger.info(f"Availability check completed: available={is_available}")
            return response

        except Exception as e:
            logger.error(f"Unexpected error in check_calendar_availability: {str(e)}")
            logger.error(traceback.format_exc())
            return {
                "status": "error",
                "message": f"Failed to check calendar availability: {str(e)}"
            }

    async def _find_available_slots(
        client: GoogleCalendarClient,
        credentials_config: Dict[str, Any],
        start_dt: datetime,
        end_dt: datetime,
        slot_duration: int,
        calendar_id: str,
        config: Dict[str, Any]
    ) -> Dict[str, Any]:
        """
        Find available time slots within a date range.

        Uses a single freebusy query to fetch all busy intervals, then computes
        available slots locally — much faster than one API call per slot.
        """
        # Helper to extract value from config (handles both dict and direct values)
        def get_config_value(key: str, default: Any) -> Any:
            value = config.get(key, default)
            if isinstance(value, dict) and "value" in value:
                if "enabled" in value and not value["enabled"]:
                    return default

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
        timezone = get_config_value("timezone", "America/Sao_Paulo")
        min_advance_time = get_config_value("minAdvanceTime", 0)
        max_duration = get_config_value("maxDuration", 0)

        # business_hours_enabled: True if explicitly enabled OR if there are day configs present
        # (handles cases where 'enabled' key is missing but days are configured)
        if isinstance(business_hours, dict):
            if "enabled" in business_hours:
                business_hours_enabled = bool(business_hours["enabled"])
            else:
                # Infer enabled from presence of day configs
                day_keys = {"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
                business_hours_enabled = any(
                    k in business_hours and business_hours[k]
                    for k in day_keys
                )
        else:
            business_hours_enabled = False

        logger.info(
            f"[FindSlots] business_hours raw={business_hours}, enabled={business_hours_enabled}, "
            f"timezone={timezone}, min_advance={min_advance_time}h, max_duration={max_duration}min"
        )

        if max_duration > 0 and slot_duration > max_duration:
            return {
                "status": "error",
                "message": f"Slot duration ({slot_duration} min) exceeds maximum duration ({max_duration} min)"
            }

        # --- Step 1: Fetch ALL busy intervals with a single freebusy query ---
        try:
            service = client.get_calendar_service(credentials_config)
            freebusy_body = {
                "timeMin": start_dt.isoformat() + "Z",
                "timeMax": end_dt.isoformat() + "Z",
                "items": [{"id": calendar_id}],
            }
            fb_result = service.freebusy().query(body=freebusy_body).execute()
            busy_raw = fb_result.get("calendars", {}).get(calendar_id, {}).get("busy", [])
        except Exception as e:
            logger.error(f"freebusy query failed: {e}")
            return {"status": "error", "message": f"Google Calendar API error: {str(e)}"}

        # Convert busy intervals to datetime pairs (timezone-naive for comparison)
        busy_intervals = []
        for b in busy_raw:
            b_start = datetime.fromisoformat(b["start"].replace("Z", "+00:00")).replace(tzinfo=None)
            b_end = datetime.fromisoformat(b["end"].replace("Z", "+00:00")).replace(tzinfo=None)
            busy_intervals.append((b_start, b_end))
        busy_intervals.sort(key=lambda x: x[0])

        logger.info(f"Fetched {len(busy_intervals)} busy intervals from freebusy query")

        # --- Step 2: Build business-hours windows per day ---
        DEFAULT_BUSINESS_START = "08:00"
        DEFAULT_BUSINESS_END = "18:00"
        DEFAULT_BUSINESS_DAYS = {"monday", "tuesday", "wednesday", "thursday", "friday"}
        day_names = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]

        current_day = start_dt.replace(hour=0, minute=0, second=0, microsecond=0)
        end_day = end_dt.replace(hour=23, minute=59, second=59, microsecond=999999)
        logger.info(f"Searching slots from {current_day} to {end_day}")

        # Minimum allowed time (for advance-time validation)
        from zoneinfo import ZoneInfo
        tz = ZoneInfo(timezone)
        now_naive = datetime.now(tz).replace(tzinfo=None)
        min_start = now_naive + timedelta(hours=min_advance_time) if min_advance_time > 0 else None

        available_slots = []

        while current_day <= end_day:
            day_name = day_names[current_day.weekday()]

            if business_hours_enabled:
                day_config = business_hours.get(day_name, {})
                if not day_config or not day_config.get("enabled"):
                    logger.debug(f"Skipping {day_name} - not a configured business day")
                    current_day += timedelta(days=1)
                    continue
                start_time_str = day_config.get("start", DEFAULT_BUSINESS_START)
                end_time_str = day_config.get("end", DEFAULT_BUSINESS_END)
            else:
                if day_name not in DEFAULT_BUSINESS_DAYS:
                    logger.debug(f"Skipping {day_name} - not a default business day (Mon-Fri)")
                    current_day += timedelta(days=1)
                    continue
                start_time_str = DEFAULT_BUSINESS_START
                end_time_str = DEFAULT_BUSINESS_END

            sh, sm = map(int, start_time_str.split(":"))
            eh, em = map(int, end_time_str.split(":"))
            day_start = current_day.replace(hour=sh, minute=sm, second=0, microsecond=0)
            day_end = current_day.replace(hour=eh, minute=em, second=0, microsecond=0)

            # Clamp to requested range
            if day_start < start_dt:
                day_start = start_dt.replace(tzinfo=None) if start_dt.tzinfo else start_dt
            if day_end > end_dt:
                day_end = end_dt.replace(tzinfo=None) if end_dt.tzinfo else end_dt

            logger.debug(f"Checking {day_name} ({current_day.date()}): {start_time_str} - {end_time_str}")

            # --- Step 3: Walk slots, skipping busy intervals without extra API calls ---
            slot_start = day_start
            busy_idx = 0  # index into sorted busy_intervals (minor optimisation)

            while slot_start + timedelta(minutes=slot_duration) <= day_end:
                slot_end = slot_start + timedelta(minutes=slot_duration)

                # Advance time check
                if min_start and slot_start < min_start:
                    slot_start += timedelta(minutes=30)
                    continue

                # Check if slot overlaps any busy interval
                is_free = True
                for b_start, b_end in busy_intervals:
                    if b_start >= slot_end:
                        break  # sorted list — no more overlaps possible
                    if b_end > slot_start:
                        is_free = False
                        # Jump past this busy block (round up to next 30-min boundary)
                        mins_past = int((b_end - slot_start).total_seconds() / 60)
                        jump = ((mins_past + 29) // 30) * 30  # ceil to 30-min grid
                        slot_start += timedelta(minutes=max(jump, 30))
                        break

                if is_free:
                    available_slots.append({
                        "start": slot_start.isoformat(),
                        "end": slot_end.isoformat(),
                        "duration_minutes": slot_duration
                    })
                    slot_start += timedelta(minutes=30)

            current_day += timedelta(days=1)

        logger.info(f"Found {len(available_slots)} available time slots in range {start_dt} to {end_dt}")

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

    # Set function metadata

    check_calendar_availability.__name__ = "check_calendar_availability"

    # Build dynamic docstring with config constraints
    config = calendar_config.get("settings", calendar_config) if calendar_config else {}

    def get_config_value(key: str, default: Any) -> Any:
        value = config.get(key, default)
        if isinstance(value, dict) and "value" in value:
            if "enabled" in value and not value["enabled"]:
                return default

            extracted_value = value["value"]
            unit = value.get("unit")

            # Convert time units to appropriate format
            if key in ["minAdvanceTime", "maxDistance"] and unit == "hours":
                return extracted_value  # Already in hours
            elif key in ["minAdvanceTime", "maxDistance"] and unit == "weeks":
                return extracted_value * 24 * 7  # Convert weeks to hours
            elif key == "maxDuration" and unit == "hours":
                return extracted_value * 60  # Convert hours to minutes
            elif key == "maxDuration" and unit == "minutes":
                return extracted_value  # Already in minutes

            return extracted_value
        return value

    # Extract configuration
    business_hours = get_config_value("businessHours", {})
    min_advance_time = get_config_value("minAdvanceTime", 0)
    max_duration = get_config_value("maxDuration", 0)
    timezone = get_config_value("timezone", "America/Sao_Paulo")

    # Build business hours description
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

    # Build constraints description
    constraints = []
    if min_advance_time > 0:
        constraints.append(f"- Meetings must be scheduled at least {min_advance_time} hours in advance")
    if max_duration > 0:
        constraints.append(f"- Maximum meeting duration: {max_duration} minutes")

    constraints_description = ""
    if constraints:
        constraints_description = "\n\nSCHEDULING CONSTRAINTS:\n" + "\n".join(constraints)

    check_calendar_availability.__doc__ = f"""Check Google Calendar availability for a time range or find available time slots.

This tool has TWO modes:

MODE 1 — FIND AVAILABLE SLOTS (use this when customer asks "what times are available?", "quando você tem horários livres?", "quais horários posso marcar?", "me mostra os horários disponíveis"):
  → Call with find_slots=True and a date range (e.g., the next 7 days)
  → The tool returns a list of available_slots the customer can choose from
  → ALWAYS use find_slots=True when you need to present options to the customer
  → Example: start_date='2024-01-15T00:00:00', end_date='2024-01-22T23:59:59', find_slots=True, slot_duration=60

MODE 2 — CHECK SPECIFIC SLOT (use this when customer already gave you a specific time):
  → Call with find_slots=False (default) and the exact start/end times
  → The tool returns whether that specific slot is available
  → Example: start_date='2024-01-16T14:00:00', end_date='2024-01-16T15:00:00'
{bh_description}{constraints_description}

CRITICAL RULES:
- When customer asks WHAT TIMES are available → ALWAYS use find_slots=True
- When customer GIVES a specific time → use find_slots=False to confirm
- Always respect business hours and constraints above
- Present slots to customer in a friendly, readable format (e.g., "Segunda-feira 18/01 às 09:00", "Terça-feira 19/01 às 14:00")
- Show at most 5-8 slots to avoid overwhelming the customer

Args:
    start_date (str): Start date/time ISO format in timezone {timezone} (e.g., '2024-01-15T09:00:00')
    end_date (str): End date/time ISO format in timezone {timezone} (e.g., '2024-01-15T10:00:00')
    calendar_id (str, optional): Calendar ID to check (default: 'primary')
    find_slots (bool): True = list available slots for customer to choose; False = check one specific time (default: False)
    slot_duration (int, optional): Duration of each slot in minutes when find_slots=True (default: 60)
"""

    return check_calendar_availability
