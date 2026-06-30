import api from '../core/api';

export interface ProactiveCampaign {
  id: string;
  name: string;
  trigger_type: 'LABEL_ADDED' | 'PIPELINE_STAGE_ENTERED' | 'SCHEDULED_DATE' | 'CONTACT_CREATED' | 'CONVERSATION_OPENED' | 'CONVERSATION_RESOLVED';
  trigger_target: string;
  delay_hours: number;
  agent_id?: string | number;
  message_template: string;
  attachment_url?: string;
  status: 'DRAFT' | 'ACTIVE' | 'PAUSED';
  created_at?: string;
  last_run_at?: string;
}

export const proactiveService = {
  getCampaigns: async () => {
    const response = await api.get('/proactive-campaigns');
    return response.data;
  },

  getCampaign: async (id: string) => {
    const response = await api.get(`/proactive-campaigns/${id}`);
    return response.data;
  },

  createCampaign: async (data: Omit<ProactiveCampaign, 'id' | 'created_at' | 'last_run_at'>) => {
    const response = await api.post('/proactive-campaigns', data);
    return response.data;
  },

  updateCampaign: async (id: string, data: Partial<ProactiveCampaign>) => {
    const response = await api.put(`/proactive-campaigns/${id}`, data);
    return response.data;
  },

  deleteCampaign: async (id: string) => {
    const response = await api.delete(`/proactive-campaigns/${id}`);
    return response.data;
  },

  cloneCampaign: async (id: string) => {
    const response = await api.post(`/proactive-campaigns/${id}/clone`);
    return response.data;
  },
  
  uploadAttachment: async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    // Assuming a generic upload endpoint exists or creating one for this purpose
    const response = await api.post('/storage/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data.url;
  }
};
