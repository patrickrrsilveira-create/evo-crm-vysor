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
} from '@evoapi/design-system';
import BrandIcon from '@/components/BrandIcon';
import {
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
  Clock,
  Settings,
  User,
  Building2,
  FileText,
  Mail,
  Loader2,
  CalendarCheck,
  CalendarClock,
  Repeat,
  Zap,
} from 'lucide-react';
import { useLanguage } from '@/hooks/useLanguage';
import { toast } from 'sonner';
import MicrosoftTeamsService from '@/services/integrations/microsoftTeamsService';

interface MicrosoftTeamsConfig {
  provider: string;
  email: string;
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
  const { t } = useLanguage('aiAgents');

  const [isConnecting, setIsConnecting] = useState(false);

  const [config, setConfig] = useState<MicrosoftTeamsConfig>({
    provider: 'microsoft_teams',
    email: initialConfig?.email || '',
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
      simultaneousBookings: {
        enabled: initialConfig?.settings?.simultaneousBookings?.enabled || false,
        limit: initialConfig?.settings?.simultaneousBookings?.limit || 1,
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
      bookingFields: initialConfig?.settings?.bookingFields || [
        { id: '1', name: 'name', label: 'Nome', enabled: false, required: false },
        { id: '2', name: 'company', label: 'Empresa', enabled: true, required: false },
        { id: '3', name: 'subject', label: 'Assunto', enabled: true, required: false },
        { id: '4', name: 'duration', label: 'Duração', enabled: false, required: false },
        { id: '5', name: 'email', label: 'E-mail', enabled: true, required: false },
        { id: '6', name: 'summary', label: 'Resumo', enabled: false, required: false },
      ],
    },
  });

  // Sync config when initialConfig changes
  useEffect(() => {
    if (initialConfig) {
      setConfig({
        provider: 'microsoft_teams',
        email: initialConfig?.email || '',
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
          simultaneousBookings: {
            enabled: initialConfig?.settings?.simultaneousBookings?.enabled || false,
            limit: initialConfig?.settings?.simultaneousBookings?.limit || 1,
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
          bookingFields: initialConfig?.settings?.bookingFields || [
            { id: '1', name: 'name', label: 'Nome', enabled: false, required: false },
            { id: '2', name: 'company', label: 'Empresa', enabled: true, required: false },
            { id: '3', name: 'subject', label: 'Assunto', enabled: true, required: false },
            { id: '4', name: 'duration', label: 'Duração', enabled: false, required: false },
            { id: '5', name: 'email', label: 'E-mail', enabled: true, required: false },
            { id: '6', name: 'summary', label: 'Resumo', enabled: false, required: false },
          ],
        },
      });
    }
  }, [initialConfig]);

  const handleConnectTeams = async () => {
    if (!config.email) {
      toast.error('Por favor, insira o e-mail (UPN) do atendente');
      return;
    }

    setIsConnecting(true);
    try {
      // Because we use Server-to-Server flow globally, the user just specifies their email.
      // We ping the backend to verify the user exists or simply just save it.
      await MicrosoftTeamsService.saveConfiguration(agentId, { ...config, connected: true });
      setConfig({ ...config, connected: true });
      onSave({ ...config, connected: true });
      toast.success('Conectado ao Microsoft Teams com sucesso!');
    } catch (error: any) {
      console.error('Error connecting to MS Teams:', error);
      toast.error('Erro ao salvar as configurações. Verifique as credenciais globais.');
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

  const updateBusinessHours = (day: string, field: 'enabled' | 'start' | 'end', value: boolean | string) => {
    setConfig((prev) => {
      const currentDayData = prev.settings?.businessHours?.[day] || { enabled: false, start: '08:00', end: '18:00' };
      return {
        ...prev,
        settings: {
          ...prev.settings,
          businessHours: {
            ...prev.settings?.businessHours,
            [day]: {
              enabled: field === 'enabled' ? (value as boolean) : currentDayData.enabled,
              start: field === 'start' ? (value as string) : currentDayData.start,
              end: field === 'end' ? (value as string) : currentDayData.end,
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
                  O sistema usa credenciais de aplicativo globais. Apenas informe o e-mail do usuário que receberá as reuniões e compromissos.
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="email">
                  E-mail do Atendente (User Principal Name)
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="usuario@suaempresa.onmicrosoft.com"
                  value={config.email}
                  onChange={(e) => setConfig({ ...config, email: e.target.value })}
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

            <TabsContent value="general" className="space-y-6">
              {/* Email Display */}
              <div className="space-y-3">
                <Label>E-mail do Atendente</Label>
                <Input value={config.email} disabled />
                <p className="text-xs text-muted-foreground">A conta que hospeda as reuniões do Teams.</p>
              </div>

              {/* Min Advance Time */}
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <Label>Tempo mínimo de antecedência</Label>
                  </div>
                  <Switch
                    checked={config.settings?.minAdvanceTime?.enabled}
                    onCheckedChange={(checked) =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          minAdvanceTime: { ...config.settings?.minAdvanceTime, enabled: checked },
                        },
                      })
                    }
                  />
                </div>
                {config.settings?.minAdvanceTime?.enabled && (
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      min="1"
                      value={config.settings?.minAdvanceTime?.value || 1}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          settings: {
                            ...config.settings,
                            minAdvanceTime: {
                              enabled: config.settings?.minAdvanceTime?.enabled || false,
                              value: parseInt(e.target.value) || 1,
                              unit: config.settings?.minAdvanceTime?.unit || 'hours',
                            },
                          },
                        })
                      }
                      className="w-24"
                    />
                    <Select
                      value={config.settings?.minAdvanceTime?.unit || 'hours'}
                      onValueChange={(value: 'hours' | 'days' | 'weeks') =>
                        setConfig({
                          ...config,
                          settings: {
                            ...config.settings,
                            minAdvanceTime: {
                              enabled: config.settings?.minAdvanceTime?.enabled || false,
                              value: config.settings?.minAdvanceTime?.value || 1,
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
                        <SelectItem value="hours">Horas</SelectItem>
                        <SelectItem value="days">Dias</SelectItem>
                        <SelectItem value="weeks">Semanas</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                )}
              </div>

              {/* Max Distance */}
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <CalendarCheck className="h-4 w-4 text-muted-foreground" />
                    <Label>Distância máxima permitida</Label>
                  </div>
                  <Switch
                    checked={config.settings?.maxDistance?.enabled}
                    onCheckedChange={(checked) =>
                      setConfig({
                        ...config,
                        settings: {
                          ...config.settings,
                          maxDistance: { ...config.settings?.maxDistance, enabled: checked },
                        },
                      })
                    }
                  />
                </div>
                {config.settings?.maxDistance?.enabled && (
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      min="1"
                      value={config.settings?.maxDistance?.value || 1}
                      onChange={(e) =>
                        setConfig({
                          ...config,
                          settings: {
                            ...config.settings,
                            maxDistance: {
                              enabled: config.settings?.maxDistance?.enabled || false,
                              value: parseInt(e.target.value) || 1,
                              unit: config.settings?.maxDistance?.unit || 'weeks',
                            },
                          },
                        })
                      }
                      className="w-24"
                    />
                    <Select
                      value={config.settings?.maxDistance?.unit || 'weeks'}
                      onValueChange={(value: 'days' | 'weeks' | 'months') =>
                        setConfig({
                          ...config,
                          settings: {
                            ...config.settings,
                            maxDistance: {
                              enabled: config.settings?.maxDistance?.enabled || false,
                              value: config.settings?.maxDistance?.value || 1,
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
                        <SelectItem value="days">Dias</SelectItem>
                        <SelectItem value="weeks">Semanas</SelectItem>
                        <SelectItem value="months">Meses</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                )}
              </div>
            </TabsContent>

            <TabsContent value="schedule" className="space-y-6">
              {/* Always Open */}
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Repeat className="h-4 w-4 text-muted-foreground" />
                    <Label>Sempre aberto</Label>
                  </div>
                </div>
                <Switch
                  checked={config.settings?.alwaysOpen}
                  onCheckedChange={(checked) =>
                    setConfig({
                      ...config,
                      settings: { ...config.settings, alwaysOpen: checked },
                    })
                  }
                />
              </div>

              {/* Business Hours */}
              {!config.settings?.alwaysOpen && (
                <div className="space-y-4">
                  <Label>Horários de atendimento</Label>
                  <div className="space-y-3">
                    {DAYS_OF_WEEK.map((day) => {
                      const dayData = config.settings?.businessHours?.[day];
                      return (
                        <div key={day} className="flex items-center gap-3">
                          <div className="flex items-center gap-2 w-32">
                            <input
                              type="checkbox"
                              checked={dayData?.enabled}
                              onChange={(e) => updateBusinessHours(day, 'enabled', e.target.checked)}
                              className="rounded"
                            />
                            <Label className="text-sm capitalize">{day}</Label>
                          </div>
                          {dayData?.enabled && (
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

            <TabsContent value="settings" className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <CalendarCheck className="h-4 w-4 text-muted-foreground" />
                    <Label>Consulta de horários livre</Label>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    O agente pode consultar a agenda para encontrar horários livres no Microsoft Teams
                  </p>
                </div>
                <Switch
                  checked={config.settings?.allowAvailabilityCheck}
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
