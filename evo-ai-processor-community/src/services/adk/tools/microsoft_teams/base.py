"""Base client for Microsoft Teams integration tools."""

import os
from typing import Any, Dict, Optional, List
from datetime import datetime, timedelta
from zoneinfo import ZoneInfo
import httpx


class MicrosoftTeamsClient:
    """Client for Microsoft Teams / Graph API with integration configuration."""

    def __init__(self, db=None):
        """
        Initialize Microsoft Teams client.

        Args:
            db: Database session for direct database access (bypasses API sanitization)
        """
        self.db = db
        self._integration_cache: Dict[str, Dict[str, Any]] = {}
        self._token_cache: Dict[str, str] = {}

    async def get_integration(
        self,
        agent_id: str
    ) -> Optional[Dict[str, Any]]:
        """
        Fetch Microsoft Teams integration configuration for an agent directly from database.

        Args:
            agent_id: The agent ID

        Returns:
            Integration configuration with credentials and settings, or None if not found
        """
        cache_key = agent_id
        if cache_key in self._integration_cache:
            return self._integration_cache[cache_key]

        # Load directly from database (no sanitization)
        from src.services.agent_service import get_agent_integration_by_provider

        integration_config = await get_agent_integration_by_provider(
            self.db, agent_id, "microsoft_teams"
        )

        if integration_config:
            integration = {
                "provider": "microsoft_teams",
                "config": integration_config
            }
            self._integration_cache[cache_key] = integration
            return integration

        return None

    async def _get_access_token(self, credentials_dict: Dict[str, Any]) -> str:
        """
        Get Microsoft Graph API access token using client credentials flow.
        
        Args:
            credentials_dict: Dictionary with tenant_id, client_id, and client_secret
            
        Returns:
            Access token string
        """
        tenant_id = credentials_dict.get("tenant_id")
        client_id = credentials_dict.get("client_id")
        client_secret = credentials_dict.get("client_secret")
        
        cache_key = f"{tenant_id}_{client_id}"
        if cache_key in self._token_cache:
            # Note: In a production app we should check token expiration
            return self._token_cache[cache_key]
            
        url = f"https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/token"
        
        data = {
            "client_id": client_id,
            "scope": "https://graph.microsoft.com/.default",
            "client_secret": client_secret,
            "grant_type": "client_credentials"
        }
        
        async with httpx.AsyncClient() as client:
            response = await client.post(url, data=data)
            response.raise_for_status()
            
            token_data = response.json()
            access_token = token_data.get("access_token")
            self._token_cache[cache_key] = access_token
            return access_token

    def is_within_business_hours(
        self,
        dt: datetime,
        business_hours: Dict[str, Any]
    ) -> bool:
        """
        Check if a datetime falls within configured business hours.
        """
        if not business_hours or not business_hours.get("enabled"):
            return True

        day_names = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
        day_name = day_names[dt.weekday()]

        day_config = business_hours.get(day_name, {})
        if not day_config.get("enabled"):
            return False

        start_time = day_config.get("start", "09:00")
        end_time = day_config.get("end", "18:00")

        current_time = dt.strftime("%H:%M")

        return start_time <= current_time <= end_time

    def validate_advance_time(
        self,
        event_time: datetime,
        min_advance_time: int,
        timezone_str: str = "America/Sao_Paulo"
    ) -> bool:
        """
        Check if event respects minimum advance time.
        """
        if min_advance_time <= 0:
            return True

        tz = ZoneInfo(timezone_str)
        now_in_tz = datetime.now(tz)

        if event_time.tzinfo is None:
            event_time_aware = event_time.replace(tzinfo=tz)
        else:
            event_time_aware = event_time

        min_allowed_time = now_in_tz + timedelta(hours=min_advance_time)
        return event_time_aware >= min_allowed_time

    def validate_max_duration(
        self,
        duration_minutes: int,
        max_duration: int
    ) -> bool:
        """
        Check if event duration is within maximum allowed.
        """
        if max_duration <= 0:
            return True

        return duration_minutes <= max_duration

    async def check_availability(
        self,
        credentials_config: Dict[str, Any],
        user_principal_name: str,
        start_time: datetime,
        end_time: datetime
    ) -> Dict[str, Any]:
        """
        Check user availability for a time range using Microsoft Graph API.
        
        Args:
            credentials_config: Microsoft Teams credentials configuration
            user_principal_name: The UPN or email of the user to check
            start_time: Start of time range
            end_time: End of time range
            
        Returns:
            Dictionary with availability information
        """
        try:
            token = await self._get_access_token(credentials_config)
            
            # Use the CalendarView API to get events in the specified range
            url = f"https://graph.microsoft.com/v1.0/users/{user_principal_name}/calendarView"
            params = {
                "startDateTime": start_time.isoformat() + "Z",
                "endDateTime": end_time.isoformat() + "Z",
                "$select": "subject,start,end,showAs"
            }
            
            headers = {
                "Authorization": f"Bearer {token}",
                "Content-Type": "application/json"
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.get(url, params=params, headers=headers)
                response.raise_for_status()
                
                events_data = response.json()
                events = events_data.get("value", [])
                
                # Filter events that actually block time (Busy, Oof)
                blocking_events = [e for e in events if e.get("showAs") in ["busy", "oof"]]
                
                return {
                    "status": "success",
                    "available": len(blocking_events) == 0,
                    "events": blocking_events
                }
                
        except Exception as e:
            return {
                "status": "error",
                "message": f"Microsoft Graph API error: {str(e)}"
            }

    async def create_meeting(
        self,
        credentials_config: Dict[str, Any],
        config: Dict[str, Any],
        summary: str,
        start_time: datetime,
        end_time: datetime,
        description: str = "",
        attendees: Optional[List[str]] = None
    ) -> Dict[str, Any]:
        """
        Create an online Microsoft Teams meeting.
        
        Args:
            credentials_config: Microsoft Teams credentials configuration
            config: Microsoft Teams settings configuration (contains user_principal_name)
            summary: Event title
            start_time: Event start time
            end_time: Event end time
            description: Event description
            attendees: List of attendee email addresses
            
        Returns:
            Dictionary with creation result
        """
        try:
            user_principal_name = config.get("user_principal_name")
            if not user_principal_name:
                return {
                    "status": "error",
                    "message": "User Principal Name (UPN/Email) not configured for this agent."
                }
                
            token = await self._get_access_token(credentials_config)
            
            # Format attendees
            attendee_list = []
            if attendees:
                for email in attendees:
                    attendee_list.append({
                        "emailAddress": {
                            "address": email,
                            "name": email
                        },
                        "type": "required"
                    })
            
            # Create event with onlineMeeting
            url = f"https://graph.microsoft.com/v1.0/users/{user_principal_name}/events"
            
            # Ensure naive datetimes are converted to aware datetimes properly
            start_dt = start_time.strftime("%Y-%m-%dT%H:%M:%S")
            end_dt = end_time.strftime("%Y-%m-%dT%H:%M:%S")
            
            event_data = {
                "subject": summary,
                "body": {
                    "contentType": "HTML",
                    "content": description
                },
                "start": {
                    "dateTime": start_dt,
                    "timeZone": config.get("timezone", "America/Sao_Paulo")
                },
                "end": {
                    "dateTime": end_dt,
                    "timeZone": config.get("timezone", "America/Sao_Paulo")
                },
                "attendees": attendee_list,
                "isOnlineMeeting": True,
                "onlineMeetingProvider": "teamsForBusiness"
            }
            
            headers = {
                "Authorization": f"Bearer {token}",
                "Content-Type": "application/json"
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.post(url, json=event_data, headers=headers)
                response.raise_for_status()
                
                created_event = response.json()
                
                return {
                    "status": "success",
                    "event_id": created_event.get("id"),
                    "meeting_link": created_event.get("onlineMeeting", {}).get("joinUrl"),
                    "summary": created_event.get("subject"),
                    "start": created_event.get("start", {}).get("dateTime"),
                    "end": created_event.get("end", {}).get("dateTime"),
                    "raw_response": created_event
                }
                
        except Exception as e:
            return {
                "status": "error",
                "message": f"Microsoft Graph API error: {str(e)}"
            }
