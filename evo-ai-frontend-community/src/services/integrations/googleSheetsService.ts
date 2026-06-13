import api from '@/services/core/agentProcessorApi';
import type {
  GoogleSheetsConfig,
  GoogleSheetsItem,
  GoogleSheetsOAuthResponse,
  GoogleSheetsConnectionResponse,
} from '@/types/integrations/googleSheets';
import { GOOGLE_OAUTH_GLOBAL_CONFIG } from '@/config/googleOAuth';

const GoogleSheetsService = {
  /**
   * Generate Google Sheets OAuth authorization URL
   */
  async generateAuthorization(agentId: string, email?: string): Promise<GoogleSheetsOAuthResponse> {
    try {
      const clientId = localStorage.getItem('GLOBAL_GOOGLE_SHEETS_CLIENT_ID') || GOOGLE_OAUTH_GLOBAL_CONFIG.clientId;

      if (!clientId) {
        throw new Error('Client ID não configurado nas Configurações Globais.');
      }

      const redirectUri = `${window.location.origin}/google-sheets/callback`;
      const scope = encodeURIComponent('https://www.googleapis.com/auth/spreadsheets');
      const state = encodeURIComponent(agentId);
      
      let url = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${encodeURIComponent(clientId)}&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=${scope}&access_type=offline&prompt=consent&state=${state}`;
      
      if (email) {
        url += `&login_hint=${encodeURIComponent(email)}`;
      }

      return { url };
    } catch (error) {
      console.error('GoogleSheetsService.generateAuthorization error:', error);
      throw error;
    }
  },

  /**
   * Complete Google Sheets OAuth flow and get spreadsheets
   */
  async completeAuthorization(
    agentId: string,
    code: string,
    state: string
  ): Promise<GoogleSheetsConnectionResponse> {
    try {
      const { data } = await api.post(
        `/agents/${agentId}/integrations/google-sheets/callback`,
        {
          code,
          state,
        }
      );
      return data;
    } catch (error) {
      console.error('GoogleSheetsService.completeAuthorization error:', error);
      throw error;
    }
  },

  /**
   * Get list of available spreadsheets
   */
  async getSpreadsheets(agentId: string): Promise<GoogleSheetsItem[]> {
    try {
      const { data } = await api.get(
        `/agents/${agentId}/integrations/google-sheets/spreadsheets`
      );
      return data.spreadsheets || [];
    } catch (error) {
      console.error('GoogleSheetsService.getSpreadsheets error:', error);
      throw error;
    }
  },

  /**
   * Save Google Sheets configuration
   */
  async saveConfiguration(
    agentId: string,
    config: Partial<GoogleSheetsConfig>
  ): Promise<{ success: boolean }> {
    try {
      const { data } = await api.put(
        `/agents/${agentId}/integrations/google-sheets`,
        config
      );
      return data;
    } catch (error) {
      console.error('GoogleSheetsService.saveConfiguration error:', error);
      throw error;
    }
  },

  /**
   * Disconnect Google Sheets
   */
  async disconnect(agentId: string): Promise<{ success: boolean }> {
    try {
      const { data } = await api.delete(
        `/agents/${agentId}/integrations/google-sheets`
      );
      return data;
    } catch (error) {
      console.error('GoogleSheetsService.disconnect error:', error);
      throw error;
    }
  },

  /**
   * Get OAuth callback URL for the current domain
   */
  getOAuthCallbackUrl(): string {
    const baseUrl = window.location.origin;
    return `${baseUrl}/oauth/google-sheets/callback`;
  },
};

export default GoogleSheetsService;
