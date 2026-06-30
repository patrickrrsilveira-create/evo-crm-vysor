import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Input,
  Label,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Textarea,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@evoapi/design-system';
import { toast } from 'sonner';
import { Loader2, Lock, LockOpen, X } from 'lucide-react';
import { useLanguage } from '@/hooks/useLanguage';
import { adminConfigService } from '@/services/admin/adminConfigService';
import { extractError } from '@/utils/apiHelpers';
import type { AdminConfigData } from '@/types/admin/adminConfig';
import ModelSelector from '@/components/ai_agents/ModelSelector';

// --- Schema factory with i18n ---

function createOpenAISchema(_t: (key: string) => string) {
  return z.object({
    OPENAI_API_URL: z.string().optional(),
    OPENAI_API_SECRET: z.string().optional().nullable(),
    OPENAI_MODEL: z.string().optional(),
    OPENAI_PROMPT_REPLY: z.string().optional(),
    OPENAI_PROMPT_SUMMARY: z.string().optional(),
    OPENAI_PROMPT_REPHRASE: z.string().optional(),
    OPENAI_PROMPT_FIX_GRAMMAR: z.string().optional(),
    OPENAI_PROMPT_SHORTEN: z.string().optional(),
    OPENAI_PROMPT_EXPAND: z.string().optional(),
    OPENAI_PROMPT_FRIENDLY: z.string().optional(),
    OPENAI_PROMPT_FORMAL: z.string().optional(),
    OPENAI_PROMPT_SIMPLIFY: z.string().optional(),
  });
}

type OpenAIFormData = z.infer<ReturnType<typeof createOpenAISchema>>;

const DEFAULTS: OpenAIFormData = {
  OPENAI_API_URL: '',
  OPENAI_API_SECRET: null,
  OPENAI_MODEL: '',
  OPENAI_PROMPT_REPLY: '',
  OPENAI_PROMPT_SUMMARY: '',
  OPENAI_PROMPT_REPHRASE: '',
  OPENAI_PROMPT_FIX_GRAMMAR: '',
  OPENAI_PROMPT_SHORTEN: '',
  OPENAI_PROMPT_EXPAND: '',
  OPENAI_PROMPT_FRIENDLY: '',
  OPENAI_PROMPT_FORMAL: '',
  OPENAI_PROMPT_SIMPLIFY: '',
};

const SECRET_FIELDS = ['OPENAI_API_SECRET'];

const AI_PROVIDERS = [
  { id: 'openai', name: 'OpenAI', url: 'https://api.openai.com/v1', defaultModel: 'gpt-4o' },
  { id: 'openrouter', name: 'OpenRouter', url: 'https://openrouter.ai/api/v1', defaultModel: 'moonshotai/kimi-k2.6:free' },
  { id: 'groq', name: 'Groq', url: 'https://api.groq.com/openai/v1', defaultModel: 'llama3-70b-8192' },
  { id: 'together', name: 'Together AI', url: 'https://api.together.xyz/v1', defaultModel: 'meta-llama/Llama-3-70b-chat-hf' },
  { id: 'deepseek', name: 'DeepSeek', url: 'https://api.deepseek.com/v1', defaultModel: 'deepseek-chat' },
  { id: 'moonshot', name: 'Moonshot AI', url: 'https://api.moonshot.cn/v1', defaultModel: 'moonshot-v1-8k' },
  { id: 'custom', name: 'Personalizado', url: '', defaultModel: '' },
];

const PROMPT_FIELDS = [
  'OPENAI_PROMPT_REPLY',
  'OPENAI_PROMPT_SUMMARY',
  'OPENAI_PROMPT_REPHRASE',
  'OPENAI_PROMPT_FIX_GRAMMAR',
  'OPENAI_PROMPT_SHORTEN',
  'OPENAI_PROMPT_EXPAND',
  'OPENAI_PROMPT_FRIENDLY',
  'OPENAI_PROMPT_FORMAL',
  'OPENAI_PROMPT_SIMPLIFY',
] as const;

function isSecretMasked(value: unknown): boolean {
  return typeof value === 'string' && value.includes('••••');
}

function buildFormValues(data: Record<string, unknown>): OpenAIFormData {
  const formValues: Record<string, unknown> = { ...DEFAULTS };
  for (const [key, value] of Object.entries(data)) {
    if (SECRET_FIELDS.includes(key)) {
      formValues[key] = isSecretMasked(value) ? '' : (value ?? '');
    } else {
      formValues[key] = value ?? formValues[key] ?? '';
    }
  }
  return formValues as OpenAIFormData;
}

export default function OpenAIConfig() {
  const { t } = useLanguage('adminSettings');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [secretModified, setSecretModified] = useState<Record<string, boolean>>({});
  const [secretConfigured, setSecretConfigured] = useState<Record<string, boolean>>({});

  const openaiSchema = useMemo(() => createOpenAISchema(t), [t]);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    control,
    watch,
    formState: { errors },
  } = useForm<OpenAIFormData>({
    resolver: zodResolver(openaiSchema),
    defaultValues: DEFAULTS,
  });

  const updateSecretStatus = (data: Record<string, unknown>) => {
    const configured: Record<string, boolean> = {};
    for (const key of SECRET_FIELDS) {
      configured[key] = isSecretMasked(data[key]);
    }
    setSecretConfigured(configured);
    setSecretModified({});
  };

  const loadConfig = useCallback(async () => {
    setLoading(true);
    try {
      const data = await adminConfigService.getConfig('openai');
      updateSecretStatus(data);
      reset(buildFormValues(data));
    } catch (error) {
      toast.error(t('openai.messages.loadError'));
    } finally {
      setLoading(false);
    }
  }, [reset, t]);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const onSubmit = async (formData: OpenAIFormData) => {
    setSaving(true);
    try {
      const payload: Record<string, unknown> = {};
      for (const [key, value] of Object.entries(formData)) {
        if (SECRET_FIELDS.includes(key)) {
          if (!secretModified[key] || value === '') {
            payload[key] = null;
          } else {
            payload[key] = value;
          }
        } else {
          payload[key] = value;
        }
      }

      const data = await adminConfigService.saveConfig('openai', payload as AdminConfigData);
      updateSecretStatus(data);
      reset(buildFormValues(data));

      toast.success(t('openai.messages.saveSuccess'));
    } catch (error) {
      const errorInfo = extractError(error);
      toast.error(t('openai.messages.saveError'), {
        description: errorInfo.message,
      });
    } finally {
      setSaving(false);
    }
  };

  const handleSecretChange = (fieldName: string, value: string) => {
    setSecretModified((prev) => ({ ...prev, [fieldName]: value.length > 0 }));
  };

  const handleClearSecret = (fieldName: string) => {
    setValue(fieldName as keyof OpenAIFormData, '');
    setSecretModified((prev) => ({ ...prev, [fieldName]: true }));
  };

  const renderSecretField = (fieldName: string, label: string, placeholder: string) => (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <Label htmlFor={fieldName}>{label}</Label>
        {!secretModified[fieldName] && (
          secretConfigured[fieldName] ? (
            <span className="inline-flex items-center gap-1 text-xs text-green-600">
              <Lock className="h-3 w-3" />
              {t('openai.secretConfigured')}
            </span>
          ) : (
            <span className="inline-flex items-center gap-1 text-xs text-sidebar-foreground/50">
              <LockOpen className="h-3 w-3" />
              {t('openai.secretNotConfigured')}
            </span>
          )
        )}
      </div>
      <div className="relative">
        <Input
          id={fieldName}
          type="password"
          autoComplete="off"
          placeholder={placeholder}
          {...register(fieldName as keyof OpenAIFormData, {
            onChange: (e: React.ChangeEvent<HTMLInputElement>) => handleSecretChange(fieldName, e.target.value),
          })}
        />
        {secretConfigured[fieldName] && !secretModified[fieldName] && (
          <button
            type="button"
            onClick={() => handleClearSecret(fieldName)}
            className="absolute right-2 top-1/2 -translate-y-1/2 p-1 rounded hover:bg-sidebar-accent text-sidebar-foreground/50 hover:text-sidebar-foreground"
            title={t('openai.clearSecret')}
            aria-label={t('openai.clearSecret')}
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
    </div>
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <div className="mb-6">
        <h2 className="text-xl font-semibold text-sidebar-foreground">{t('openai.title')}</h2>
        <p className="text-sm text-sidebar-foreground/70 mt-1">{t('openai.description')}</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Connection Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t('openai.connection.cardTitle')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-5">
            <div className="space-y-2">
              <Label>Provedor de IA (Atalho)</Label>
              <Select 
                value={
                  AI_PROVIDERS.find(p => p.url !== '' && watch('OPENAI_API_URL')?.includes(p.url))?.id || 'custom'
                } 
                onValueChange={(val) => {
                  if (val === 'custom') {
                    setValue('OPENAI_API_URL', '');
                    setValue('OPENAI_MODEL', '');
                  } else {
                    const provider = AI_PROVIDERS.find(p => p.id === val);
                    if (provider) {
                      setValue('OPENAI_API_URL', provider.url, { shouldValidate: true, shouldDirty: true });
                      if (provider.defaultModel) {
                        setValue('OPENAI_MODEL', provider.defaultModel, { shouldValidate: true, shouldDirty: true });
                      }
                    }
                  }
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione um provedor..." />
                </SelectTrigger>
                <SelectContent>
                  {AI_PROVIDERS.map((provider) => (
                    <SelectItem key={provider.id} value={provider.id}>
                      {provider.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground">
                Selecione o provedor para preencher a URL base e sugerir um modelo automaticamente.
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="OPENAI_API_URL">{t('openai.connection.fields.apiUrl')}</Label>
              <Input
                id="OPENAI_API_URL"
                placeholder={t('openai.connection.placeholders.apiUrl')}
                {...register('OPENAI_API_URL')}
              />
              {errors.OPENAI_API_URL && (
                <p className="text-xs text-destructive">{errors.OPENAI_API_URL.message}</p>
              )}
            </div>

            {renderSecretField('OPENAI_API_SECRET', t('openai.connection.fields.apiSecret'), t('openai.connection.placeholders.apiSecret'))}

            <div className="space-y-2">
              <Controller
                name="OPENAI_MODEL"
                control={control}
                render={({ field }) => (
                  <ModelSelector
                    id="OPENAI_MODEL"
                    value={field.value || ''}
                    onChange={field.onChange}
                    apiKeys={[]} // Pass empty array so it shows availableModels
                    label={t('openai.connection.fields.model')}
                    showLabel={true}
                    className="w-full"
                    error={errors.OPENAI_MODEL?.message}
                  />
                )}
              />
            </div>

            <div className="p-4 bg-muted/50 rounded-lg border text-sm text-muted-foreground mt-4">
              <strong>Nota:</strong> A configuração de Transcrição de Áudio foi movida para um menu dedicado. Acesse a seção <a href="/settings/admin/audio-transcription" className="text-primary hover:underline font-medium">Transcrição de Áudio</a> no menu lateral para configurar seu provedor de Speech-to-Text independentemente.
            </div>
          </CardContent>
        </Card>

        {/* AI Prompts */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t('openai.prompts.cardTitle')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-5">
            {PROMPT_FIELDS.map((fieldName) => (
              <div key={fieldName} className="space-y-2">
                <Label htmlFor={fieldName}>{t(`openai.prompts.fields.${fieldName}`)}</Label>
                <Textarea
                  id={fieldName}
                  rows={4}
                  placeholder={t(`openai.prompts.placeholders.${fieldName}`)}
                  {...register(fieldName)}
                />
              </div>
            ))}
          </CardContent>
        </Card>

        <div className="pt-2">
          <Button type="submit" disabled={saving} aria-label={t('openai.save')}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {saving ? t('openai.saving') : t('openai.save')}
          </Button>
        </div>
      </form>
    </div>
  );
}
