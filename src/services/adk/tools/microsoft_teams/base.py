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
        Check user availability.
        Since we are using a Webhook to an external system (n8n/API), 
        we assume available and let the external system handle logic.
        """
        return {
            "status": "success",
            "available": True,
            "events": []
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
        Create an online Microsoft Teams meeting via Webhook (n8n/API).
        """
        try:
            webhook_url = config.get("webhookUrl")
            if not webhook_url:
                return {
                    "status": "error",
                    "message": "URL do Webhook não configurada para este agente."
                }
                
            # Create payload for the webhook
            start_dt = start_time.isoformat()
            end_dt = end_time.isoformat()
            
            payload = {
                "subject": summary,
                "description": description,
                "start": start_dt,
                "end": end_dt,
                "attendees": attendees or []
            }
            
            async with httpx.AsyncClient() as client:
                response = await client.post(webhook_url, json=payload, timeout=30.0)
                response.raise_for_status()
                
                response_data = response.json()
                
                # Try common link field names
                meeting_link = response_data.get("link") or response_data.get("joinWebUrl") or response_data.get("url")
                
                if not meeting_link:
                    return {
                        "status": "error",
                        "message": "Webhook did not return a valid meeting link (expected 'link', 'joinWebUrl' or 'url')"
                    }
                
                return {
                    "status": "success",
                    "event_id": "webhook-generated",
                    "meeting_link": meeting_link,
                    "summary": summary,
                    "start": start_dt,
                    "end": end_dt,
                    "raw_response": response_data
                }
                
        except Exception as e:
            return {
                "status": "error",
                "message": f"Erro ao disparar webhook: {str(e)}"
            }
