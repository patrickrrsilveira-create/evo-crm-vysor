import React, { useState, useEffect } from 'react';
import { CampaignForm } from '../components/Campaigns/CampaignForm';
import { campaignsService, CampaignPayload, Campaign } from '../services/campaignsService';

export const CampaignsPage: React.FC = () => {
  const [showForm, setShowForm] = useState(false);
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [loading, setLoading] = useState(false);
  const accountId = 1; // Get from current user context

  useEffect(() => {
    loadCampaigns();
  }, []);

  const loadCampaigns = async () => {
    try {
      const { data } = await campaignsService.list(accountId);
      setCampaigns(data);
    } catch (error) {
      console.error('Failed to load campaigns', error);
    }
  };

  const handleSubmit = async (payload: CampaignPayload) => {
    setLoading(true);
    try {
      await campaignsService.create(payload);
      setShowForm(false);
      loadCampaigns();
    } catch (error) {
      console.error('Failed to create campaign', error);
    } finally {
      setLoading(false);
    }
  };

  const handleStart = async (id: string) => {
    try {
      await campaignsService.start(id);
      loadCampaigns();
    } catch (error) {
      console.error('Failed to start campaign', error);
    }
  };

  const handlePause = async (id: string) => {
    try {
      await campaignsService.pause(id);
      loadCampaigns();
    } catch (error) {
      console.error('Failed to pause campaign', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure?')) return;
    try {
      await campaignsService.delete(id);
      loadCampaigns();
    } catch (error) {
      console.error('Failed to delete campaign', error);
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Campanhas</h1>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          {showForm ? 'Cancelar' : '+ Nova Campanha'}
        </button>
      </div>

      {showForm && (
        <div className="mb-6">
          <CampaignForm onSubmit={handleSubmit} loading={loading} />
        </div>
      )}

      <div className="grid gap-4">
        {campaigns.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            Nenhuma campanha criada. Crie uma para começar!
          </div>
        ) : (
          campaigns.map((campaign) => (
            <div key={campaign.id} className="border rounded-lg p-4 hover:shadow-md">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <h3 className="font-semibold text-lg">{campaign.name}</h3>
                  <p className="text-sm text-gray-600">
                    Status: <span className="font-medium">{campaign.status}</span>
                  </p>
                  <p className="text-sm text-gray-600 mt-1">
                    Enviados: {campaign.sent_count} / {campaign.total_recipients}
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    {new Date(campaign.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="flex gap-2">
                  {campaign.status === 'draft' && (
                    <button
                      onClick={() => handleStart(campaign.id)}
                      className="px-3 py-1 bg-green-600 text-white rounded text-sm hover:bg-green-700"
                    >
                      Iniciar
                    </button>
                  )}
                  {campaign.status === 'running' && (
                    <button
                      onClick={() => handlePause(campaign.id)}
                      className="px-3 py-1 bg-yellow-600 text-white rounded text-sm hover:bg-yellow-700"
                    >
                      Pausar
                    </button>
                  )}
                  {['draft', 'paused'].includes(campaign.status) && (
                    <button
                      onClick={() => handleDelete(campaign.id)}
                      className="px-3 py-1 bg-red-600 text-white rounded text-sm hover:bg-red-700"
                    >
                      Deletar
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
