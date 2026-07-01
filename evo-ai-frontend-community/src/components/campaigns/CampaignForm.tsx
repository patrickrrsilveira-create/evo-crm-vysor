import React, { useState, useEffect } from 'react';
import { useAppDataStore } from '../../store/appDataStore';

interface CampaignFormProps {
  onSubmit?: (data: any) => Promise<void>;
  loading?: boolean;
}

interface SelectedInboxes {
  [key: string]: boolean;
}

export const CampaignForm: React.FC<CampaignFormProps> = ({ onSubmit, loading = false }) => {
  const { inboxes, fetchInboxes } = useAppDataStore();
  const [formData, setFormData] = useState({
    name: '',
    selectedInboxes: {} as SelectedInboxes,
    variants: [{ body: '', mediaUrl: '', mediaType: '', weight: 1 }],
    audienceMode: 'static',
    segmentFilter: {},
  });

  useEffect(() => {
    fetchInboxes();
  }, []);

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({ ...prev, name: e.target.value }));
  };

  const handleInboxToggle = (inboxId: string) => {
    setFormData(prev => ({
      ...prev,
      selectedInboxes: {
        ...prev.selectedInboxes,
        [inboxId]: !prev.selectedInboxes[inboxId],
      },
    }));
  };

  const handleVariantChange = (index: number, field: string, value: string | number) => {
    const newVariants = [...formData.variants];
    newVariants[index] = { ...newVariants[index], [field]: value };
    setFormData(prev => ({ ...prev, variants: newVariants }));
  };

  const addVariant = () => {
    setFormData(prev => ({
      ...prev,
      variants: [...prev.variants, { body: '', mediaUrl: '', mediaType: '', weight: 1 }],
    }));
  };

  const removeVariant = (index: number) => {
    setFormData(prev => ({
      ...prev,
      variants: prev.variants.filter((_, i) => i !== index),
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const selectedInboxIds = Object.entries(formData.selectedInboxes)
      .filter(([, selected]) => selected)
      .map(([id]) => parseInt(id));

    const payload = {
      account_id: 1, // Will get from current user context
      name: formData.name,
      inbox_ids: selectedInboxIds,
      audience: {
        mode: formData.audienceMode,
        segment_filter: formData.segmentFilter,
        contacts: [], // Will be filled by segment resolver
      },
      variants: formData.variants.filter(v => v.body.trim()),
    };

    if (onSubmit) {
      await onSubmit(payload);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 p-6 bg-white rounded-lg shadow">
      <div>
        <label className="block text-sm font-medium mb-2">Campaign Name</label>
        <input
          type="text"
          value={formData.name}
          onChange={handleNameChange}
          placeholder="New Campaign"
          className="w-full px-3 py-2 border rounded-md"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">Select Channels (Multi-check)</label>
        <div className="space-y-2 border rounded-md p-3">
          {inboxes?.map(inbox => (
            <label key={inbox.id} className="flex items-center">
              <input
                type="checkbox"
                checked={formData.selectedInboxes[inbox.id] || false}
                onChange={() => handleInboxToggle(inbox.id)}
                className="rounded"
              />
              <span className="ml-2 flex-1">
                {inbox.name} ({inbox.channel_type}) {inbox.phone_number && ` - ${inbox.phone_number}`}
              </span>
            </label>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">Content Variants</label>
        <div className="space-y-3">
          {formData.variants.map((variant, index) => (
            <div key={index} className="border rounded-md p-3">
              <textarea
                value={variant.body}
                onChange={e => handleVariantChange(index, 'body', e.target.value)}
                placeholder="Message body (supports {Olá|Hi|E aí})"
                className="w-full px-3 py-2 border rounded-md mb-2"
                rows={3}
              />
              <div className="grid grid-cols-2 gap-2">
                <input
                  type="text"
                  value={variant.mediaUrl}
                  onChange={e => handleVariantChange(index, 'mediaUrl', e.target.value)}
                  placeholder="Media URL (optional)"
                  className="px-3 py-2 border rounded-md"
                />
                <select
                  value={variant.mediaType}
                  onChange={e => handleVariantChange(index, 'mediaType', e.target.value)}
                  className="px-3 py-2 border rounded-md"
                >
                  <option value="">No media</option>
                  <option value="video/mp4">Video (MP4)</option>
                  <option value="image/jpeg">Image (JPEG)</option>
                  <option value="image/png">Image (PNG)</option>
                </select>
              </div>
              <div className="flex justify-between items-center mt-2">
                <label className="text-sm">Weight (rotation):</label>
                <input
                  type="number"
                  min="1"
                  value={variant.weight}
                  onChange={e => handleVariantChange(index, 'weight', parseInt(e.target.value))}
                  className="w-16 px-2 py-1 border rounded-md"
                />
              </div>
              {formData.variants.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeVariant(index)}
                  className="mt-2 text-sm text-red-600 hover:text-red-800"
                >
                  Remove variant
                </button>
              )}
            </div>
          ))}
        </div>
        <button
          type="button"
          onClick={addVariant}
          className="mt-2 px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded-md text-sm"
        >
          + Add Variant
        </button>
      </div>

      <div className="flex gap-3 justify-end">
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
        >
          {loading ? 'Creating...' : 'Create Campaign'}
        </button>
      </div>
    </form>
  );
};
