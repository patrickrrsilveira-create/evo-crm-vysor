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

  // ---- Legacy: knowledge page lists bots linked via agent_bots (CRM legacy bots) ----
  getAgentBots: async (id: number | string): Promise<any> => {
    const response = await api.get(`/knowledge_bases/${id}/agent_bots`);
    return response.data;
  },

  getAiAgents: async (id: number | string): Promise<any> => {
    const response = await api.get(`/knowledge_bases/${id}/ai_agents`);
    return response.data;
  },

  // ---- New: linking evo_core AI agents to knowledge bases ----

  // Get all knowledge bases linked to an evo_core agent
  getAgentLinkedBases: async (agentId: number | string): Promise<KnowledgeBase[]> => {
    try {
      const response = await api.get(`/ai_agents/${agentId}/knowledge_bases`);
      return (response.data as any)?.data || response.data || [];
    } catch (e) {
      console.error('Error fetching linked bases:', e);
      return [];
    }
  },

  // Link an evo_core agent to a knowledge base
  linkAgentBot: async (baseId: number | string, agentId: number | string): Promise<any> => {
    const response = await api.post(`/ai_agents/${agentId}/knowledge_bases`, { knowledge_base_id: baseId });
    return response.data;
  },

  // Unlink an evo_core agent from a knowledge base
  unlinkAgentBot: async (baseId: number | string, agentId: number | string): Promise<any> => {
    const response = await api.delete(`/ai_agents/${agentId}/knowledge_bases/${baseId}`);
    return response.data;
  },

  deleteDocument: async (baseId: number | string, documentId: number | string): Promise<any> => {
    const response = await api.delete(`/knowledge_bases/${baseId}/knowledge_documents/${documentId}`);
    return response.data;
  }
};
