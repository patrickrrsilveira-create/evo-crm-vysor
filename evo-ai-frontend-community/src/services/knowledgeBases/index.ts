import api from '../core/api';
import { ApiResponse } from '@/types/core/api';

export interface KnowledgeBase {
  id: number;
  name: string;
  description: string;
  documents_count: number;
  created_at: string;
  updated_at: string;
}

export const knowledgeBasesService = {
  list: async (): Promise<ApiResponse<KnowledgeBase[]>> => {
    const response = await api.get<ApiResponse<KnowledgeBase[]>>('/knowledge_bases');
    return response.data;
  },
  
  get: async (id: number | string): Promise<ApiResponse<KnowledgeBase>> => {
    const response = await api.get<ApiResponse<KnowledgeBase>>(`/knowledge_bases/${id}`);
    return response.data;
  },

  create: async (data: { name: string; description?: string }): Promise<ApiResponse<KnowledgeBase>> => {
    const response = await api.post<ApiResponse<KnowledgeBase>>('/knowledge_bases', { knowledge_base: data });
    return response.data;
  },

  update: async (id: number | string, data: { name?: string; description?: string }): Promise<ApiResponse<KnowledgeBase>> => {
    const response = await api.put<ApiResponse<KnowledgeBase>>(`/knowledge_bases/${id}`, { knowledge_base: data });
    return response.data;
  },

  delete: async (id: number | string): Promise<ApiResponse<void>> => {
    const response = await api.delete<ApiResponse<void>>(`/knowledge_bases/${id}`);
    return response.data;
  },

  getDocuments: async (id: number | string): Promise<any> => {
    const response = await api.get(`/knowledge_bases/${id}/knowledge_documents`);
    return response.data;
  },

  getAgentBots: async (id: number | string): Promise<any> => {
    const response = await api.get(`/knowledge_bases/${id}/agent_bots`);
    return response.data;
  },

  getAgentLinkedBases: async (agentId: number | string): Promise<KnowledgeBase[]> => {
    try {
      const response = await api.get(`/agent_bots/${agentId}`);
      return response.data?.payload?.knowledge_bases || response.data?.knowledge_bases || [];
    } catch (e) {
      console.error('Error fetching linked bases:', e);
      return [];
    }
  },

  linkAgentBot: async (baseId: number | string, agentId: number | string): Promise<any> => {
    const response = await api.post(`/knowledge_bases/${baseId}/agent_bots`, { agent_bot_id: agentId });
    return response.data;
  },

  unlinkAgentBot: async (baseId: number | string, agentId: number | string): Promise<any> => {
    const response = await api.delete(`/knowledge_bases/${baseId}/agent_bots/${agentId}`);
    return response.data;
  },

  deleteDocument: async (baseId: number | string, documentId: number | string): Promise<any> => {
    const response = await api.delete(`/knowledge_bases/${baseId}/knowledge_documents/${documentId}`);
    return response.data;
  }
};
