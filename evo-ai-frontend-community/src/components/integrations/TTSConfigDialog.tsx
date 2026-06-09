import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Input,
  Label,
  Button,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Slider,
} from '@evoapi/design-system';
import { Loader2 } from 'lucide-react';
import { useLanguage } from '@/hooks/useLanguage';

interface TTSConfig {
  provider: 'elevenlabs' | 'fish' | 'cartesia' | 'kokoro' | 'voxtral';
  apiKey: string;
  respondInAudio: 'when_client_asks' | 'always' | 'never';
  voice: string;
  // ElevenLabs specific
  stability?: number;
  similarity?: number;
}

interface TTSConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSave: (config: TTSConfig) => void;
  onDeactivate?: () => void;
  initialConfig?: Partial<TTSConfig>;
}

interface Voice {
  voice_id: string;
  name: string;
}

const TTSConfigDialog = ({
  open,
  onOpenChange,
  onSave,
  onDeactivate,
  initialConfig,
}: TTSConfigDialogProps) => {
  const { t } = useLanguage('aiAgents');

  const [config, setConfig] = useState<TTSConfig>({
    provider: initialConfig?.provider || 'elevenlabs',
    apiKey: initialConfig?.apiKey || '',
    respondInAudio: initialConfig?.respondInAudio || 'when_client_asks',
    voice: initialConfig?.voice || '',
    stability: initialConfig?.stability ?? 80,
    similarity: initialConfig?.similarity ?? 50,
  });

  const [availableVoices, setAvailableVoices] = useState<Voice[]>([]);
  const [loadingVoices, setLoadingVoices] = useState(false);
  const [voicesError, setVoicesError] = useState(false);

  useEffect(() => {
    if (initialConfig) {
      setConfig({
        provider: initialConfig.provider || 'elevenlabs',
        apiKey: initialConfig.apiKey || '',
        respondInAudio: initialConfig.respondInAudio || 'when_client_asks',
        voice: initialConfig.voice || '',
        stability: initialConfig.stability ?? 80,
        similarity: initialConfig.similarity ?? 50,
      });
    }
  }, [initialConfig]);

  // Fetch voices based on provider
  useEffect(() => {
    const fetchVoices = async () => {
      if (!config.apiKey || config.apiKey.length < 10) {
        setAvailableVoices([]);
        return;
      }

      setLoadingVoices(true);
      setVoicesError(false);

      try {
        if (config.provider === 'elevenlabs') {
          const response = await fetch('https://api.elevenlabs.io/v1/voices', {
            method: 'GET',
            headers: { 'xi-api-key': config.apiKey },
          });

          if (!response.ok) throw new Error('Failed to fetch ElevenLabs voices');

          const data = await response.json();
          const voices = data.voices.map((voice: any) => ({
            voice_id: voice.voice_id,
            name: voice.name,
          }));

          setAvailableVoices(voices);

          if (!config.voice && voices.length > 0) {
            setConfig((prev) => ({ ...prev, voice: voices[0].voice_id }));
          }
        } else if (config.provider === 'fish') {
          // Fish Audio may not have a simple voices endpoint without user auth or we may just allow manual input
          // For now, if we can't fetch Fish Audio voices easily via API key, we allow manual entry.
          // Since it's typically custom references, let's leave it empty and let the user paste the ID.
          setAvailableVoices([]);
        }
      } catch (error) {
        console.error('Error fetching voices:', error);
        setVoicesError(true);
        setAvailableVoices([]);
      } finally {
        setLoadingVoices(false);
      }
    };

    fetchVoices();
  }, [config.apiKey, config.provider]);

  const handleSave = () => {
    onSave(config);
    onOpenChange(false);
  };

  const handleDeactivate = () => {
    if (onDeactivate) {
      onDeactivate();
    }
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {t('edit.integrations.tts.configTitle') || 'Configurar Motor TTS'}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {/* Provedor TTS */}
          <div className="space-y-2">
            <Label>Provedor TTS</Label>
            <Select
              value={config.provider}
              onValueChange={(value: 'elevenlabs' | 'fish' | 'cartesia' | 'kokoro' | 'voxtral') =>
                setConfig({ ...config, provider: value, voice: '' })
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="elevenlabs">ElevenLabs</SelectItem>
                <SelectItem value="fish">Fish Audio</SelectItem>
                <SelectItem value="cartesia">Cartesia</SelectItem>
                <SelectItem value="kokoro">Kokoro</SelectItem>
                <SelectItem value="voxtral">Voxtral</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* API Key */}
          <div className="space-y-2">
            <Label htmlFor="apiKey">API Key</Label>
            <Input
              id="apiKey"
              type="password"
              placeholder={`Insira sua API Key do provedor ${config.provider}`}
              value={config.apiKey}
              onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
            />
          </div>

          {/* Quando responder em áudio */}
          {config.apiKey && !voicesError && (
            <div className="space-y-3">
              <Label>Quando responder em áudio:</Label>
              <Select
                value={config.respondInAudio}
                onValueChange={(value: 'when_client_asks' | 'always' | 'never') =>
                  setConfig({ ...config, respondInAudio: value })
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="when_client_asks">Quando a pergunta do cliente for em áudio</SelectItem>
                  <SelectItem value="always">Responder sempre em áudio</SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          {/* Qual voz deseja usar */}
          {config.apiKey && !voicesError && config.provider === 'elevenlabs' && (
            <div className="space-y-3">
              <Label>Voz (ElevenLabs):</Label>
              <div className="flex gap-2">
                <Select
                  value={config.voice}
                  onValueChange={(value) => setConfig({ ...config, voice: value })}
                  disabled={loadingVoices || availableVoices.length === 0}
                >
                  <SelectTrigger className="flex-1">
                    {loadingVoices ? (
                      <span className="flex items-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin" /> Carregando vozes...
                      </span>
                    ) : (
                      <SelectValue placeholder="Selecione uma voz" />
                    )}
                  </SelectTrigger>
                  <SelectContent>
                    {availableVoices.map((voice) => (
                      <SelectItem key={voice.voice_id} value={voice.voice_id}>
                        {voice.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          )}

          {/* Manual Voice ID Entry (For Fish, Cartesia, Kokoro, Voxtral) */}
          {config.apiKey && config.provider !== 'elevenlabs' && (
            <div className="space-y-2">
              <Label htmlFor="voiceId">ID da Voz (Model Reference ID)</Label>
              <Input
                id="voiceId"
                type="text"
                placeholder={`Insira o ID do modelo da voz (${config.provider})`}
                value={config.voice}
                onChange={(e) => setConfig({ ...config, voice: e.target.value })}
              />
            </div>
          )}

          {config.apiKey && voicesError && config.provider === 'elevenlabs' && (
            <div className="p-4 border border-destructive/50 rounded-lg bg-destructive/10">
              <p className="text-sm text-destructive">
                Erro ao carregar vozes - Verifique se a API Key está correta.
              </p>
            </div>
          )}

          {/* Estabilidade (ElevenLabs) */}
          {config.apiKey && !voicesError && config.provider === 'elevenlabs' && (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Estabilidade:</Label>
                <span className="text-sm font-medium">{config.stability}%</span>
              </div>
              <Slider
                value={[config.stability || 80]}
                onValueChange={(value) => setConfig({ ...config, stability: value[0] })}
                min={0}
                max={100}
                step={1}
                className="w-full"
              />
            </div>
          )}

          {/* Similaridade (ElevenLabs) */}
          {config.apiKey && !voicesError && config.provider === 'elevenlabs' && (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Similaridade:</Label>
                <span className="text-sm font-medium">{config.similarity}%</span>
              </div>
              <Slider
                value={[config.similarity || 50]}
                onValueChange={(value) => setConfig({ ...config, similarity: value[0] })}
                min={0}
                max={100}
                step={1}
                className="w-full"
              />
            </div>
          )}

          {/* Botões de ação */}
          <div className="flex flex-col gap-3 pt-4">
            <Button
              onClick={handleSave}
              disabled={!config.apiKey || !config.voice}
              className="w-full"
            >
              APLICAR CONFIGURAÇÕES
            </Button>

            {onDeactivate && (
              <Button
                variant="ghost"
                onClick={handleDeactivate}
                className="w-full text-destructive hover:text-destructive/80"
              >
                Desativar integração
              </Button>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default TTSConfigDialog;
