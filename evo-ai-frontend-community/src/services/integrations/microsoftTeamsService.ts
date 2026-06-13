import api from '@/services/core/agentProcessorApi';

export interface MicrosoftTeamsConfig {
  provider: string;
  email: string;
  connected: boolean;
  settings: any;
}

const MicrosoftTeamsService = {
  /**
   * Save Microsoft Teams configuration for an agent
   */
  async saveConfiguration(
    agentId: string,
    config: Partial<MicrosoftTeamsConfig>
  ): Promise<{ success: boolean }> {
    try {
      const { data } = await api.put(
        `/agents/${agentId}/integrations/microsoft-teams`,
        config
      );
      return data;
    } catch (error) {
      console.error('MicrosoftTeamsService.saveConfiguration error:', error);
      throw error;
    }
  },

  /**
   * Disconnect Microsoft Teams integration from an agent
   */
  async disconnect(agentId: string): Promise<{ success: boolean }> {
    try {
      const { data } = await api.delete(
        `/agents/${agentId}/integrations/microsoft-teams`
      );
      return data;
    } catch (error) {
      console.error('MicrosoftTeamsService.disconnect error:', error);
      throw error;
    }
  },
};

export default MicrosoftTeamsService;
