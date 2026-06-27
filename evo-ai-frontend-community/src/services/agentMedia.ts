import api from './api';

export interface AgentMedia {
  filename: string;
  url: string;
  size: number;
}

export const agentMediaService = {
  list: async (agentId: string): Promise<AgentMedia[]> => {
    try {
      const response = await api.get(`/adk/agent_media/${agentId}/media`);
      return response.data.media || [];
    } catch (error) {
      console.error('Error fetching agent media:', error);
      throw error;
    }
  },

  upload: async (agentId: string, file: File) => {
    try {
      const formData = new FormData();
      formData.append('file', file);
      
      const response = await api.post(`/adk/agent_media/${agentId}/media`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data;
    } catch (error) {
      console.error('Error uploading agent media:', error);
      throw error;
    }
  },

  delete: async (agentId: string, filename: string) => {
    try {
      const response = await api.delete(`/adk/agent_media/${agentId}/media/${filename}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting agent media:', error);
      throw error;
    }
  }
};
