import { useState, useEffect } from 'react';
import { useLanguage } from '@/hooks/useLanguage';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  Button,
  Input,
  Card,
  Switch,
  Label
} from '@evoapi/design-system';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@evoapi/design-system/select';
import ModelSelector from '@/components/ai_agents/ModelSelector';
import { Brain, AlertCircle, ExternalLink } from 'lucide-react';
import { OpenAIHook, OpenAIFormData, IntegrationHook } from '@/types/integrations';

interface OpenAIModalProps {
  hook?: IntegrationHook;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: Record<string, unknown>) => Promise<void>;
  isNew?: boolean;
  loading?: boolean;
}

export default function OpenAIModal({
  hook,
  open,
  onOpenChange,
  onSubmit,

  loading: submitting = false
}: OpenAIModalProps) {
  const { t } = useLanguage('integrations');
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<OpenAIFormData>({
    api_key: '',
    enable_audio_transcription: false,
    provider: 'openai',
    model: '',
    base_url: ''
  });

  const providers = [
    { id: 'openai', name: 'OpenAI' },
    { id: 'gemini', name: 'Google Gemini' },
    { id: 'anthropic', name: 'Anthropic' },
    { id: 'openrouter', name: 'OpenRouter' },
    { id: 'deepseek', name: 'DeepSeek' },
    { id: 'together_ai', name: 'Together AI' },
    { id: 'fireworks_ai', name: 'Fireworks AI' },
    { id: 'perplexity', name: 'Perplexity' },
    { id: 'bedrock', name: 'AWS Bedrock' },
    { id: 'vertex_ai', name: 'Google Vertex AI' },
    { id: 'custom_openai_compatible', name: 'Custom (OpenAI Compatible)' }
  ];

  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    const openaiHook = hook as OpenAIHook | undefined;
    if (openaiHook?.settings) {
      setFormData({
        api_key: openaiHook.settings.api_key || '',
        enable_audio_transcription: openaiHook.settings.enable_audio_transcription || false,
        provider: openaiHook.settings.provider || 'openai',
        model: openaiHook.settings.model || '',
        base_url: openaiHook.settings.base_url || ''
      });
    } else {
      setFormData({
        api_key: '',
        enable_audio_transcription: false,
        provider: 'openai',
        model: '',
        base_url: ''
      });
    }
    setErrors({});
  }, [hook, open]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.api_key.trim()) {
      newErrors.api_key = t('openai.modal.fields.apiKey.required');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setLoading(true);
    try {
      await onSubmit(formData as unknown as Record<string, unknown>);
    } catch {
      // Error is handled by parent component
    } finally {
      setLoading(false);
    }
  };

  const openOpenAIDoc = () => {
    window.open('https://platform.openai.com/api-keys', '_blank');
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Brain className="w-5 h-5" />
            {hook ? t('openai.modal.updateTitle') : t('openai.modal.title')}
          </DialogTitle>
          <DialogDescription>{t('openai.modal.description')}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* API Key Configuration */}
          <Card className="p-4">
            <h4 className="font-semibold mb-4">{t('openai.modal.apiConfig')}</h4>

            <div className="space-y-4">
              <div>
                <label htmlFor="provider" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Provedor *
                </label>
                <Select
                  value={formData.provider || 'openai'}
                  onValueChange={(val) => setFormData(prev => ({ ...prev, provider: val, model: '' }))}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Selecione um provedor" />
                  </SelectTrigger>
                  <SelectContent>
                    {providers.map(p => (
                      <SelectItem key={p.id} value={p.id}>{p.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label htmlFor="base_url" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  URL Base (Opcional)
                </label>
                <Input
                  id="base_url"
                  placeholder="Ex: https://api.openai.com/v1"
                  value={formData.base_url || ''}
                  onChange={(e) => setFormData(prev => ({ ...prev, base_url: e.target.value }))}
                />
              </div>

              <div>
                <label htmlFor="api_key" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  {t('openai.modal.fields.apiKey.label')} *
                </label>
                <Input
                  id="api_key"
                  type="password"
                  placeholder={t('openai.modal.fields.apiKey.placeholder')}
                  value={formData.api_key}
                  onChange={(e) => setFormData(prev => ({ ...prev, api_key: e.target.value }))}
                  className={errors.api_key ? 'border-red-500' : ''}
                />
                {errors.api_key && (
                  <p className="text-sm text-red-600 mt-1">{errors.api_key}</p>
                )}
                <p className="text-xs text-slate-500 mt-1">
                  {t('openai.modal.fields.apiKey.hint')}
                </p>
              </div>

              <div>
                <ModelSelector
                  label="Modelo"
                  value={formData.model || ''}
                  onChange={(val) => setFormData(prev => ({ ...prev, model: val }))}
                  apiKeys={[{ id: 'temp-id', name: 'Temp', provider: formData.provider || 'openai' } as any]}
                  apiKeyId="temp-id"
                  className="w-full"
                />
              </div>

              {/* Audio Transcription Toggle */}
              <div className="flex items-center justify-between p-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50">
                <div className="flex-1">
                  <Label htmlFor="enable_audio_transcription" className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('openai.modal.fields.enableAudioTranscription.label')}
                  </Label>
                  <p className="text-xs text-slate-500 mt-1">
                    {t('openai.modal.fields.enableAudioTranscription.description')}
                  </p>
                </div>
                <Switch
                  id="enable_audio_transcription"
                  checked={formData.enable_audio_transcription || false}
                  onCheckedChange={(checked) => setFormData(prev => ({ ...prev, enable_audio_transcription: checked }))}
                />
              </div>
            </div>
          </Card>

          {/* Features Info */}
          <Card className="p-4 bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800">
            <h4 className="font-semibold mb-3 text-blue-800 dark:text-blue-200">
              {t('openai.modal.features.title')}
            </h4>
            <ul className="text-sm text-blue-700 dark:text-blue-300 space-y-2">
              <li>• {t('openai.modal.features.suggestions')}</li>
              <li>• {t('openai.modal.features.summaries')}</li>
              <li>• {t('openai.modal.features.improvement')}</li>
              <li>• {t('openai.modal.features.correction')}</li>
              <li>• {t('openai.modal.features.labels')}</li>
            </ul>
          </Card>

          {/* Security Warning */}
          <Card className="p-4 bg-amber-50 dark:bg-amber-900/20 border-amber-200 dark:border-amber-800">
            <div className="flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-amber-600 dark:text-amber-400 mt-0.5 flex-shrink-0" />
              <div className="text-sm text-amber-700 dark:text-amber-300">
                <strong>{t('openai.modal.security.title')}</strong> {t('openai.modal.security.description')}
              </div>
            </div>
          </Card>

          {/* Actions */}
          <div className="flex justify-between pt-4 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={openOpenAIDoc}
              className="flex items-center gap-2"
            >
              <ExternalLink className="w-4 h-4" />
              {t('openai.modal.actions.help')}
            </Button>

            <div className="flex gap-3">
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                {t('openai.modal.actions.cancel')}
              </Button>
              <Button type="submit" disabled={loading || submitting}>
                {loading || submitting ? t('openai.modal.actions.saving') : (hook ? t('openai.modal.actions.update') : t('openai.modal.actions.configure'))}
              </Button>
            </div>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
