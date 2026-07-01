import api from './core/api';

export interface CampaignPayload {
  account_id: number;
  name: string;
  inbox_ids: number[];
  audience: {
    mode: 'static' | 'dynamic';
    segment_filter: Record<string, any>;
    contacts: Array<{ contact_id: number; recipient: string; timezone: string }>;
  };
  variants: Array<{
    body: string;
    media_url?: string;
    media_type?: string;
    weight?: number;
  }>;
}

export interface Campaign {
  id: string;
  account_id: number;
  name: string;
  status: string;
  total_recipients: number;
  sent_count: number;
  created_at: string;
}

class CampaignService {
  async create(payload: CampaignPayload): Promise<Campaign> {
    const response = await api.post('/campaigns', payload);
    return response.data;
  }

  async list(accountId: number, page = 1, pageSize = 20): Promise<{ data: Campaign[]; total: number }> {
    const response = await api.get('/campaigns', {
      params: { account_id: accountId, page, page_size: pageSize },
    });
    return response.data;
  }

  async get(id: string): Promise<Campaign> {
    const response = await api.get(`/campaigns/${id}`);
    return response.data;
  }

  async start(id: string): Promise<{ status: string }> {
    const response = await api.post(`/campaigns/${id}/start`);
    return response.data;
  }

  async pause(id: string): Promise<{ status: string }> {
    const response = await api.post(`/campaigns/${id}/pause`);
    return response.data;
  }

  async cancel(id: string): Promise<{ status: string }> {
    const response = await api.post(`/campaigns/${id}/cancel`);
    return response.data;
  }

  async stats(id: string): Promise<Record<string, number>> {
    const response = await api.get(`/campaigns/${id}/stats`);
    return response.data.stats;
  }

  async delete(id: string): Promise<{ deleted: boolean }> {
    const response = await api.delete(`/campaigns/${id}`);
    return response.data;
  }
}

export const campaignsService = new CampaignService();
