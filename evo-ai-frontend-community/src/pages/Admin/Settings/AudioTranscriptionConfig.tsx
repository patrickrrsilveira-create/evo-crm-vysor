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
  Switch,
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

function createAudioTranscriptionSchema(_t: (key: string) => string) {
  return z.object({
    AUDIO_TRANSCRIPTION_ENABLED: z.union([z.boolean(), z.string()]).optional(),
    AUDIO_TRANSCRIPTION_API_URL: z.string().optional(),
    AUDIO_TRANSCRIPTION_API_SECRET: z.string().optional().nullable(),
    AUDIO_TRANSCRIPTION_MODEL: z.string().optional(),
  });
}

type AudioTranscriptionFormData = z.infer<ReturnType<typeof createAudioTranscriptionSchema>>;

const DEFAULTS: AudioTranscriptionFormData = {
  AUDIO_TRANSCRIPTION_ENABLED: false,
  AUDIO_TRANSCRIPTION_API_URL: '',
  AUDIO_TRANSCRIPTION_API_SECRET: null,
  AUDIO_TRANSCRIPTION_MODEL: '',
};

const SECRET_FIELDS = ['AUDIO_TRANSCRIPTION_API_SECRET'];

const AUDIO_PROVIDERS = [
  { id: 'openai', name: 'OpenAI', url: 'https://api.openai.com/v1', defaultModel: 'whisper-1' },
  { id: 'groq', name: 'Groq', url: 'https://api.groq.com/openai/v1', defaultModel: 'whisper-large-v3-turbo' },
  { id: 'openrouter', name: 'OpenRouter', url: 'https://openrouter.ai/api/v1', defaultModel: 'qwen/qwen3-asr-flash-2026-02-10' },
  { id: 'custom', name: 'Personalizado', url: '', defaultModel: '' },
];

function isSecretMasked(value: unknown): boolean {
  return typeof value === 'string' && value.includes('••••');
}

function toBool(value: unknown): boolean {
  if (typeof value === 'boolean') return value;
  if (typeof value === 'string') return value === 'true';
  return false;
}

function buildFormValues(data: Record<string, unknown>): AudioTranscriptionFormData {
  const formValues: Record<string, unknown> = { ...DEFAULTS };
  for (const [key, value] of Object.entries(data)) {
    if (SECRET_FIELDS.includes(key)) {
      formValues[key] = isSecretMasked(value) ? '' : (value ?? '');
    } else {
      formValues[key] = value ?? formValues[key] ?? '';
    }
  }
  return formValues as AudioTranscriptionFormData;
}

export default function AudioTranscriptionConfig() {
  const { t } = useLanguage('adminSettings');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [secretModified, setSecretModified] = useState<Record<string, boolean>>({});
  const [secretConfigured, setSecretConfigured] = useState<Record<string, boolean>>({});

  const schema = useMemo(() => createAudioTranscriptionSchema(t), [t]);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    control,
    watch,
    formState: { errors },
  } = useForm<AudioTranscriptionFormData>({
    resolver: zodResolver(schema),
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
      const data = await adminConfigService.getConfig('audio_transcription');
      updateSecretStatus(data);
      reset(buildFormValues(data));
    } catch (error) {
      toast.error(t('audioTranscription.messages.loadError'));
    } finally {
      setLoading(false);
    }
  }, [reset, t]);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const onSubmit = async (formData: AudioTranscriptionFormData) => {
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

      const data = await adminConfigService.saveConfig('audio_transcription', payload as AdminConfigData);
      updateSecretStatus(data);
      reset(buildFormValues(data));

      toast.success(t('audioTranscription.messages.saveSuccess'));
    } catch (error) {
      const errorInfo = extractError(error);
      toast.error(t('audioTranscription.messages.saveError'), {
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
    setValue(fieldName as keyof AudioTranscriptionFormData, '');
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
          {...register(fieldName as keyof AudioTranscriptionFormData, {
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
        <h2 className="text-xl font-semibold text-sidebar-foreground">{t('audioTranscription.title')}</h2>
        <p className="text-sm text-sidebar-foreground/70 mt-1">{t('audioTranscription.description')}</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t('audioTranscription.connection.cardTitle')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-5">
            <Controller
              name="AUDIO_TRANSCRIPTION_ENABLED"
              control={control}
              render={({ field }) => (
                <div className="flex items-center justify-between p-3 rounded-lg border border-sidebar-border bg-sidebar-accent/30">
                  <div>
                    <Label htmlFor="AUDIO_TRANSCRIPTION_ENABLED" className="text-sm font-medium">
                      {t('audioTranscription.connection.fields.enabled')}
                    </Label>
                    <p className="text-xs text-muted-foreground mt-1">
                      Habilita a transcrição automática de áudios recebidos usando o provedor configurado.
                    </p>
                  </div>
                  <Switch
                    id="AUDIO_TRANSCRIPTION_ENABLED"
                    checked={toBool(field.value)}
                    onCheckedChange={field.onChange}
                  />
                </div>
              )}
            />

            <div className="space-y-2">
              <Label>Provedor de IA (Atalho)</Label>
              <Select 
                value={
                  AUDIO_PROVIDERS.find(p => p.url !== '' && watch('AUDIO_TRANSCRIPTION_API_URL')?.includes(p.url))?.id || 'custom'
                } 
                onValueChange={(val) => {
                  if (val === 'custom') {
                    setValue('AUDIO_TRANSCRIPTION_API_URL', '');
                    setValue('AUDIO_TRANSCRIPTION_MODEL', '');
                  } else {
                    const provider = AUDIO_PROVIDERS.find(p => p.id === val);
                    if (provider) {
                      setValue('AUDIO_TRANSCRIPTION_API_URL', provider.url, { shouldValidate: true, shouldDirty: true });
                      if (provider.defaultModel) {
                        setValue('AUDIO_TRANSCRIPTION_MODEL', provider.defaultModel, { shouldValidate: true, shouldDirty: true });
                      }
                    }
                  }
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione um provedor..." />
                </SelectTrigger>
                <SelectContent>
                  {AUDIO_PROVIDERS.map((provider) => (
                    <SelectItem key={provider.id} value={provider.id}>
                      {provider.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground">
                Selecione o provedor para preencher a URL base e sugerir um modelo de transcrição automaticamente.
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="AUDIO_TRANSCRIPTION_API_URL">{t('audioTranscription.connection.fields.apiUrl')}</Label>
              <Input
                id="AUDIO_TRANSCRIPTION_API_URL"
                placeholder={t('audioTranscription.connection.placeholders.apiUrl')}
                {...register('AUDIO_TRANSCRIPTION_API_URL')}
              />
              {errors.AUDIO_TRANSCRIPTION_API_URL && (
                <p className="text-xs text-destructive">{errors.AUDIO_TRANSCRIPTION_API_URL.message}</p>
              )}
            </div>

            {renderSecretField('AUDIO_TRANSCRIPTION_API_SECRET', t('audioTranscription.connection.fields.apiSecret'), t('audioTranscription.connection.placeholders.apiSecret'))}

            <div className="space-y-2">
              <Controller
                name="AUDIO_TRANSCRIPTION_MODEL"
                control={control}
                render={({ field }) => (
                  <ModelSelector
                    id="AUDIO_TRANSCRIPTION_MODEL"
                    value={field.value || ''}
                    onChange={field.onChange}
                    apiKeys={[]}
                    label={t('audioTranscription.connection.fields.model')}
                    showLabel={true}
                    className="w-full"
                    error={errors.AUDIO_TRANSCRIPTION_MODEL?.message}
                  />
                )}
              />
            </div>
          </CardContent>
        </Card>

        <div className="pt-2">
          <Button type="submit" disabled={saving} aria-label={t('audioTranscription.save')}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {saving ? t('audioTranscription.saving') : t('audioTranscription.save')}
          </Button>
        </div>
      </form>
    </div>
  );
}
