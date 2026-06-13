import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { Card, Button, Input, Label } from '@evoapi/design-system';
import { adminConfigService } from '@/services/admin/adminConfigService';
import type { AdminConfigData } from '@/types/admin/adminConfig';
import BaseHeader from '@/components/base/BaseHeader';
import BrandIcon from '@/components/BrandIcon';

export default function MicrosoftTeamsGlobalPage() {
  const navigate = useNavigate();
  const [clientId, setClientId] = useState('');
  const [clientSecret, setClientSecret] = useState('');
  const [tenantId, setTenantId] = useState('');
  const [ownerId, setOwnerId] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    const loadConfig = async () => {
      try {
        const data = await adminConfigService.getConfig('microsoft_teams');
        setClientId((data.MICROSOFT_TEAMS_CLIENT_ID as string) || '');
        setClientSecret((data.MICROSOFT_TEAMS_CLIENT_SECRET as string) || '');
        setTenantId((data.MICROSOFT_TEAMS_TENANT_ID as string) || '');
        setOwnerId((data.MICROSOFT_TEAMS_MEETING_OWNER_ID as string) || '');
      } catch (error) {
        console.error('Error loading global microsoft_teams config:', error);
      } finally {
        setIsLoading(false);
      }
    };
    loadConfig();
  }, []);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const payload = {
        MICROSOFT_TEAMS_CLIENT_ID: clientId,
        MICROSOFT_TEAMS_CLIENT_SECRET: clientSecret,
        MICROSOFT_TEAMS_TENANT_ID: tenantId,
        MICROSOFT_TEAMS_MEETING_OWNER_ID: ownerId,
      } as AdminConfigData;
      
      await adminConfigService.saveConfig('microsoft_teams', payload);
      toast.success('Configurações globais do Teams salvas com sucesso!');
      navigate('/settings/integrations');
    } catch (error) {
      toast.error('Erro ao salvar as configurações. Verifique se você é um Super Admin.');
      console.error(error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDisconnect = async () => {
    setIsSaving(true);
    try {
      const payload = {
        MICROSOFT_TEAMS_CLIENT_ID: null,
        MICROSOFT_TEAMS_CLIENT_SECRET: null,
        MICROSOFT_TEAMS_TENANT_ID: null,
        MICROSOFT_TEAMS_MEETING_OWNER_ID: null,
      } as AdminConfigData;
      
      await adminConfigService.saveConfig('microsoft_teams', payload);
      setClientId('');
      setClientSecret('');
      setTenantId('');
      setOwnerId('');
      toast.success('Credenciais globais removidas com sucesso!');
    } catch (error) {
      toast.error('Erro ao remover as configurações.');
      console.error(error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="h-full flex flex-col p-4">
      <BaseHeader
        title="Microsoft Teams (Global)"
        subtitle="Configuração global Server-to-Server (App Permissions) do Microsoft Teams para todos os agentes."
        secondaryActions={[
          {
            label: 'Voltar',
            onClick: () => navigate('/settings/integrations'),
            variant: 'outline'
          }
        ]}
      />

      <div className="mt-6 flex-1 overflow-auto">
        <Card className="max-w-2xl p-6">
          <div className="flex items-center gap-4 mb-6">
            <div className="p-3 bg-primary/10 rounded-full flex items-center justify-center">
              <span className="text-2xl font-bold text-primary">T</span>
            </div>
            <div>
              <h2 className="text-xl font-semibold">Credenciais Microsoft Graph (Client Credentials)</h2>
              <p className="text-sm text-muted-foreground">
                Configurações para permissão global (Application permissions) no Azure AD para leitura/escrita de calendários e reuniões.
              </p>
            </div>
          </div>

          <div className="space-y-6">
            <div className="bg-slate-100 dark:bg-slate-800 p-3 rounded-md text-sm border border-slate-200 dark:border-slate-700">
              <p className="font-medium mb-1">Permissões Necessárias no Azure AD:</p>
              <ul className="text-xs text-muted-foreground list-disc pl-4 space-y-1">
                <li>Calendars.ReadWrite</li>
                <li>OnlineMeetings.ReadWrite.All</li>
              </ul>
              <p className="text-xs mt-2 text-muted-foreground">
                Certifique-se de conceder o <b>Admin Consent</b> para estas permissões no portal do Azure.
              </p>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="tenant_id">Tenant ID (Diretório)</Label>
                <Input
                  id="tenant_id"
                  value={tenantId}
                  onChange={(e) => setTenantId(e.target.value)}
                  placeholder="Ex: 8a7f9c..."
                  className="mt-1"
                />
              </div>

              <div>
                <Label htmlFor="owner_id">Organizador Padrão (User Principal Name ou Object ID)</Label>
                <Input
                  id="owner_id"
                  value={ownerId}
                  onChange={(e) => setOwnerId(e.target.value)}
                  placeholder="Ex: atendimento@suaempresa.com.br"
                  className="mt-1"
                />
                <p className="text-xs text-muted-foreground mt-1">
                  E-mail corporativo ou ID do usuário da Microsoft que será o "Dono" das reuniões geradas.
                </p>
              </div>

              <div>
                <Label htmlFor="client_id">Client ID (ID do Aplicativo)</Label>
                <Input
                  id="client_id"
                  value={clientId}
                  onChange={(e) => setClientId(e.target.value)}
                  placeholder="Ex: 5b3..."
                  className="mt-1"
                />
              </div>

              <div>
                <Label htmlFor="client_secret">Client Secret (Valor do Segredo do Cliente)</Label>
                <Input
                  id="client_secret"
                  type="password"
                  value={clientSecret}
                  onChange={(e) => setClientSecret(e.target.value)}
                  placeholder="Seu client secret do Azure AD"
                  className="mt-1"
                />
              </div>
            </div>

            <div className="flex justify-between items-center pt-4 border-t">
              <Button variant="destructive" onClick={handleDisconnect} disabled={!clientId && !clientSecret && !tenantId}>
                Limpar Credenciais
              </Button>
              <div className="flex gap-3">
                <Button variant="outline" onClick={() => navigate('/settings/integrations')}>
                  Cancelar
                </Button>
                <Button onClick={handleSave} disabled={isSaving || isLoading}>
                  {isSaving ? 'Salvando...' : 'Salvar Configurações'}
                </Button>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
