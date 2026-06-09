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
    return api.get<KnowledgeBase[]>('/api/v1/knowledge_bases');
  },
  
  get: async (id: number | string): Promise<ApiResponse<KnowledgeBase>> => {
    return api.get<KnowledgeBase>(`/api/v1/knowledge_bases/${id}`);
  },

  create: async (data: { name: string; description?: string }): Promise<ApiResponse<KnowledgeBase>> => {
    return api.post<KnowledgeBase>('/api/v1/knowledge_bases', { knowledge_base: data });
  },

  update: async (id: number | string, data: { name?: string; description?: string }): Promise<ApiResponse<KnowledgeBase>> => {
    return api.put<KnowledgeBase>(`/api/v1/knowledge_bases/${id}`, { knowledge_base: data });
  },

  delete: async (id: number | string): Promise<ApiResponse<void>> => {
    return api.delete<void>(`/api/v1/knowledge_bases/${id}`);
  }
};
