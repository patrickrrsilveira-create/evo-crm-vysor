import api from '@/services/core/agentProcessorApi';
import type {
  GoogleCalendarConfig,
  GoogleCalendarItem,
  GoogleCalendarOAuthResponse,
  GoogleCalendarConnectionResponse
} from '@/types/integrations';
import { GOOGLE_OAUTH_GLOBAL_CONFIG } from '@/config/googleOAuth';

const GoogleCalendarService = {
  /**
   * Generate Google Calendar OAuth authorization URL
   */
  async generateAuthorization(agentId: string, email?: string): Promise<GoogleCalendarOAuthResponse> {
    try {
      const clientId = localStorage.getItem('GLOBAL_GOOGLE_CALENDAR_CLIENT_ID') || GOOGLE_OAUTH_GLOBAL_CONFIG.clientId || undefined;
      const clientSecret = localStorage.getItem('GLOBAL_GOOGLE_CALENDAR_CLIENT_SECRET') || GOOGLE_OAUTH_GLOBAL_CONFIG.clientSecret || undefined;

      const { data } = await api.post(
        `/agents/${agentId}/integrations/google-calendar/authorization`,
        { 
          email, 
          client_id: clientId,
          client_secret: clientSecret
        }
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.generateAuthorization error:', error);
      throw error;
    }
  },

  /**
   * Complete Google Calendar OAuth flow and get calendars
   */
  async completeAuthorization(
    agentId: string,
    code: string,
    state: string
  ): Promise<GoogleCalendarConnectionResponse> {
    try {
      const { data } = await api.post(
        `/agents/${agentId}/integrations/google-calendar/callback`,
        {
          code,
          state,
        }
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.completeAuthorization error:', error);
      throw error;
    }
  },

  /**
   * Get list of available calendars
   */
  async getCalendars(agentId: string): Promise<GoogleCalendarItem[]> {
    try {
      const { data } = await api.get(
        `/agents/${agentId}/integrations/google-calendar/calendars`
      );
      return data.calendars || [];
    } catch (error) {
      console.error('GoogleCalendarService.getCalendars error:', error);
      throw error;
    }
  },

  /**
   * Save Google Calendar configuration
   */
  async saveConfiguration(
    agentId: string,
    config: Partial<GoogleCalendarConfig>
  ): Promise<{ success: boolean }> {
    try {
      const { data } = await api.put(
        `/agents/${agentId}/integrations/google-calendar`,
        config
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.saveConfiguration error:', error);
      throw error;
    }
  },

  /**
   * Disconnect Google Calendar
   */
  async disconnect(agentId: string): Promise<{ success: boolean }> {
    try {
      const { data } = await api.delete(
        `/agents/${agentId}/integrations/google-calendar`
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.disconnect error:', error);
      throw error;
    }
  },

  /**
   * Check availability for a specific date/time
   */
  async checkAvailability(
    agentId: string,
    params: {
      calendarId: string;
      start: string; // ISO date
      end: string; // ISO date
    }
  ): Promise<{ available: boolean; slots?: Array<{ start: string; end: string }> }> {
    try {
      const { data } = await api.post(
        `/agents/${agentId}/integrations/google-calendar/availability`,
        params
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.checkAvailability error:', error);
      throw error;
    }
  },

  /**
   * Create a calendar event
   */
  async createEvent(
    agentId: string,
    event: {
      calendarId: string;
      summary: string;
      description?: string;
      start: string; // ISO date
      end: string; // ISO date
      attendees?: Array<{ email: string; name?: string }>;
      meetLink?: boolean;
    }
  ): Promise<{ success: boolean; eventId?: string; meetLink?: string }> {
    try {
      const { data } = await api.post(
        `/agents/${agentId}/integrations/google-calendar/events`,
        event
      );
      return data;
    } catch (error) {
      console.error('GoogleCalendarService.createEvent error:', error);
      throw error;
    }
  },

  /**
   * Get OAuth callback URL for the current domain
   */
  getOAuthCallbackUrl(): string {
    const baseUrl = window.location.origin;
    return `${baseUrl}/oauth/google-calendar/callback`;
  },
};

export default GoogleCalendarService;
