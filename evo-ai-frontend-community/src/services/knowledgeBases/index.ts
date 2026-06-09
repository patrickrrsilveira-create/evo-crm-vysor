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
  }
};
