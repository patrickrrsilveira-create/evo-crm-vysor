import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  Input,
  Label,
  Button,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Switch,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@evoapi/design-system';

import {
  Calendar,
  Settings,
  Loader2,
  CalendarCheck,
  CalendarClock,
  Repeat,
  Zap,
} from 'lucide-react';

import { toast } from 'sonner';
import MicrosoftTeamsService from '@/services/integrations/microsoftTeamsService';

interface MicrosoftTeamsConfig {
  provider: string;
  webhookUrl: string;
  connected: boolean;
  settings: any;
}

interface MicrosoftTeamsConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSave: (config: MicrosoftTeamsConfig) => void;
  onDisconnect?: () => void;
  initialConfig?: Partial<MicrosoftTeamsConfig>;
  agentId: string;
}

const DAYS_OF_WEEK = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];

const MicrosoftTeamsConfigDialog = ({
  open,
  onOpenChange,
  onSave,
  onDisconnect,
  initialConfig,
  agentId,
}: MicrosoftTeamsConfigDialogProps) => {

  const [isConnecting, setIsConnecting] = useState(false);

  const [config, setConfig] = useState<MicrosoftTeamsConfig>({
    provider: 'microsoft_teams',
    webhookUrl: initialConfig?.webhookUrl || '',
    connected: initialConfig?.connected || false,
    settings: {
      minAdvanceTime: {
        enabled: initialConfig?.settings?.minAdvanceTime?.enabled || false,
        value: initialConfig?.settings?.minAdvanceTime?.value || 1,
        unit: initialConfig?.settings?.minAdvanceTime?.unit || 'hours',
      },
      maxDistance: {
        enabled: initialConfig?.settings?.maxDistance?.enabled || false,
        value: initialConfig?.settings?.maxDistance?.value || 1,
        unit: initialConfig?.settings?.maxDistance?.unit || 'weeks',
      },
      maxDuration: {
        value: initialConfig?.settings?.maxDuration?.value || 1,
        unit: initialConfig?.settings?.maxDuration?.unit || 'hours',
      },
      alwaysOpen: initialConfig?.settings?.alwaysOpen || false,
      businessHours: initialConfig?.settings?.businessHours || {
        monday: { enabled: true, start: '08:00', end: '18:00' },
        tuesday: { enabled: true, start: '08:00', end: '18:00' },
        wednesday: { enabled: true, start: '08:00', end: '18:00' },
        thursday: { enabled: true, start: '08:00', end: '18:00' },
        friday: { enabled: true, start: '08:00', end: '18:00' },
        saturday: { enabled: false, start: '08:00', end: '18:00' },
        sunday: { enabled: false, start: '08:00', end: '18:00' },
      },
      allowAvailabilityCheck: initialConfig?.settings?.allowAvailabilityCheck ?? true,
    },
  });

  useEffect(() => {
    if (initialConfig && open) {
      setConfig((prev) => ({
        ...prev,
        webhookUrl: initialConfig.webhookUrl || '',
        connected: initialConfig.connected || false,
        settings: {
          ...prev.settings,
          ...initialConfig.settings,
        }
      }));
    }
  }, [initialConfig, open]);

  const handleConnectTeams = async () => {
    if (!config.webhookUrl) {
      toast.error('Por favor, insira a URL do Webhook');
      return;
    }

    setIsConnecting(true);
    try {
      await MicrosoftTeamsService.saveConfiguration(agentId, { ...config, connected: true });
      setConfig({ ...config, connected: true });
      onSave({ ...config, connected: true });
      toast.success('Conectado ao Microsoft Teams com sucesso!');
    } catch (error: any) {
      console.error('Error connecting to MS Teams:', error);
      toast.error('Erro ao salvar as configurações.');
    } finally {
      setIsConnecting(false);
    }
  };

  const handleSave = async () => {
    try {
      await MicrosoftTeamsService.saveConfiguration(agentId, config);
      onSave(config);
      toast.success('Configurações salvas com sucesso!');
      onOpenChange(false);
    } catch (error) {
      console.error('Error saving Microsoft Teams configuration:', error);
      toast.error('Erro ao salvar configurações');
    }
  };

  const handleDisconnect = async () => {
    try {
      await MicrosoftTeamsService.disconnect(agentId);
      if (onDisconnect) {
        onDisconnect();
      }
      toast.success('Microsoft Teams desconectado com sucesso!');
      onOpenChange(false);
    } catch (error) {
      console.error('Error disconnecting Microsoft Teams:', error);
      toast.error('Erro ao desconectar Microsoft Teams');
    }
  };

  const updateBusinessHours = (day: string, field: 'enabled' | 'start' | 'end', value: any) => {
    setConfig((prev) => {
      const currentDayData = prev.settings?.businessHours?.[day] || { enabled: false, start: '08:00', end: '18:00' };
      return {
        ...prev,
        settings: {
          ...prev.settings,
          businessHours: {
            ...prev.settings?.businessHours,
            [day]: {
              ...currentDayData,
              [field]: value
            },
          },
        },
      };
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span className="flex items-center justify-center w-5 h-5 font-bold text-blue-600 bg-blue-100 rounded">T</span>
            Configurar Microsoft Teams
          </DialogTitle>
          <DialogDescription>
            Configure as opções de agendamento do agente via Microsoft Teams.
          </DialogDescription>
        </DialogHeader>

        {!config.connected ? (
          <div className="space-y-6 py-4">
            <div className="text-center space-y-4">
              <div className="flex justify-center">
                <div className="p-4 bg-primary/10 rounded-full">
                  <span className="flex items-center justify-center w-12 h-12 text-2xl font-bold text-blue-600 bg-blue-100 rounded">T</span>
                </div>
              </div>
              <div>
                <h3 className="text-lg font-semibold">
                  Conectar com Microsoft Teams
                </h3>
                <p className="text-sm text-muted-foreground">
                  Insira o webhook que o agente usará para confirmar ou notificar os agendamentos.
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="webhookUrl">
                  URL do Webhook (n8n ou API)
                </Label>
                <Input
                  id="webhookUrl"
                  type="url"
                  placeholder="https://seu-n8n.com/webhook/..."
                  value={config.webhookUrl}
                  onChange={(e) => setConfig({ ...config, webhookUrl: e.target.value })}
                />
              </div>

              <Button onClick={handleConnectTeams} disabled={isConnecting} className="w-full" size="lg">
                {isConnecting ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Conectando...
                  </>
                ) : (
                  <>
                    <Calendar className="mr-2 h-4 w-4" />
                    Conectar com MS Teams
                  </>
                )}
              </Button>
            </div>
          </div>
        ) : (
          <Tabs defaultValue="general" className="w-full">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="general">
                <Settings className="h-4 w-4 mr-2" />
                Geral
              </TabsTrigger>
              <TabsTrigger value="schedule">
                <CalendarClock className="h-4 w-4 mr-2" />
                Horários
              </TabsTrigger>
              <TabsTrigger value="settings">
                <Zap className="h-4 w-4 mr-2" />
                Configurações
              </TabsTrigger>
            </TabsList>

            <TabsContent value="general" className="space-y-6 pt-4">
              <div className="space-y-4">
                <Label>Duração Padrão da Reunião</Label>
                <div className="flex gap-2">
                  <Input
                    type="number"
                    min="1"
                    value={config.settings?.maxDuration?.value || 1}
                    onChange={(e) =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          maxDuration: {
                            value: parseInt(e.target.value) || 1,
                            unit: config.settings?.maxDuration?.unit || 'hours',
                          },
                        },
                      })
                    }
                    className="w-24"
                  />
                  <Select
                    value={config.settings?.maxDuration?.unit || 'hours'}
                    onValueChange={(value: 'minutes' | 'hours') =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          maxDuration: {
                            value: config.settings?.maxDuration?.value || 1,
                            unit: value,
                          },
                        },
                      })
                    }
                  >
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="minutes">Minutos</SelectItem>
                      <SelectItem value="hours">Horas</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Tempo mínimo de antecedência</Label>
                    <p className="text-sm text-muted-foreground">
                      Tempo mínimo antes de um evento poder ser agendado
                    </p>
                  </div>
                  <Switch
                    checked={config.settings?.minAdvanceTime?.enabled || false}
                    onCheckedChange={(checked) =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          minAdvanceTime: {
                            ...config.settings?.minAdvanceTime,
                            enabled: checked,
                          },
                        },
                      })
                    }
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Distância máxima no futuro</Label>
                    <p className="text-sm text-muted-foreground">
                      Até quanto tempo no futuro os clientes podem agendar
                    </p>
                  </div>
                  <Switch
                    checked={config.settings?.maxDistance?.enabled || false}
                    onCheckedChange={(checked) =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          maxDistance: {
                            ...config.settings?.maxDistance,
                            enabled: checked,
                          },
                        },
                      })
                    }
                  />
                </div>
              </div>
            </TabsContent>

            <TabsContent value="schedule" className="space-y-6 pt-4">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Repeat className="h-4 w-4 text-muted-foreground" />
                    <Label>Sempre aberto</Label>
                  </div>
                </div>
                <Switch
                  checked={config.settings?.alwaysOpen || false}
                  onCheckedChange={(checked) =>
                    setConfig({
                      ...config,
                      settings: { ...config.settings, alwaysOpen: checked },
                    })
                  }
                />
              </div>

              {!config.settings?.alwaysOpen && (
                <div className="space-y-4">
                  <Label>Horários de atendimento</Label>
                  <div className="space-y-3">
                    {DAYS_OF_WEEK.map((day) => {
                      const dayData = config.settings?.businessHours?.[day] || { enabled: false, start: '08:00', end: '18:00' };
                      return (
                        <div key={day} className="flex items-center gap-3">
                          <div className="flex items-center gap-2 w-32">
                            <input
                              type="checkbox"
                              checked={dayData.enabled}
                              onChange={(e) => updateBusinessHours(day, 'enabled', e.target.checked)}
                              className="rounded"
                            />
                            <Label className="text-sm capitalize">{day}</Label>
                          </div>
                          {dayData.enabled && (
                            <>
                              <Input
                                type="time"
                                value={dayData.start}
                                onChange={(e) => updateBusinessHours(day, 'start', e.target.value)}
                                className="w-32"
                              />
                              <span>-</span>
                              <Input
                                type="time"
                                value={dayData.end}
                                onChange={(e) => updateBusinessHours(day, 'end', e.target.value)}
                                className="w-32"
                              />
                            </>
                          )}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}
            </TabsContent>

            <TabsContent value="settings" className="space-y-6 pt-4">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <CalendarCheck className="h-4 w-4 text-muted-foreground" />
                    <Label>Consulta de horários livre</Label>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    O agente pode consultar a agenda para encontrar horários livres
                  </p>
                </div>
                <Switch
                  checked={config.settings?.allowAvailabilityCheck || false}
                  onCheckedChange={(checked) =>
                    setConfig({
                      ...config,
                      settings: { ...config.settings, allowAvailabilityCheck: checked },
                    })
                  }
                />
              </div>

              <div className="pt-6 border-t mt-6 flex justify-between items-center">
                <Button variant="destructive" onClick={handleDisconnect}>
                  Desconectar Integração
                </Button>
                <div className="space-x-3">
                  <Button variant="outline" onClick={() => onOpenChange(false)}>
                    Cancelar
                  </Button>
                  <Button onClick={handleSave}>Salvar Configurações</Button>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        )}
      </DialogContent>
    </Dialog>
  );
};

export default MicrosoftTeamsConfigDialog;
