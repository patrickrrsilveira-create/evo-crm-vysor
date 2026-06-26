import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { Card, Button, Input, Label } from '@evoapi/design-system';
import { adminConfigService } from '@/services/admin/adminConfigService';
import type { AdminConfigData } from '@/types/admin/adminConfig';
import BaseHeader from '@/components/base/BaseHeader';
import BrandIcon from '@/components/BrandIcon';

export default function GoogleSheetsGlobalPage() {
  const navigate = useNavigate();
  const [clientId, setClientId] = useState('');
  const [clientSecret, setClientSecret] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    const loadConfig = async () => {
      try {
        const data = await adminConfigService.getConfig('google_oauth');
        setClientId((data.GOOGLE_OAUTH_CLIENT_ID as string) || '');
        setClientSecret((data.GOOGLE_OAUTH_CLIENT_SECRET as string) || '');
      } catch (error) {
        console.error('Error loading global google_oauth config:', error);
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
        GOOGLE_OAUTH_CLIENT_ID: clientId,
        GOOGLE_OAUTH_CLIENT_SECRET: clientSecret,
      } as AdminConfigData;
      
      await adminConfigService.saveConfig('google_oauth', payload);
      toast.success('Configurações globais salvas com sucesso no banco de dados!');
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
        GOOGLE_OAUTH_CLIENT_ID: null,
        GOOGLE_OAUTH_CLIENT_SECRET: null,
      } as AdminConfigData;
      
      await adminConfigService.saveConfig('google_oauth', payload);
      setClientId('');
      setClientSecret('');
      toast.success('Credenciais globais removidas com sucesso do banco de dados!');
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
        title="Google Sheets (Global)"
        subtitle="Configuração global de credenciais OAuth do Google Sheets para todos os agentes."
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
            <div className="p-3 bg-primary/10 rounded-full">
              <BrandIcon id="google-sheets" size={32} className="h-8 w-8" />
            </div>
            <div>
              <h2 className="text-xl font-semibold">Credenciais OAuth Google</h2>
              <p className="text-sm text-muted-foreground">
                Configurações usadas quando os agentes não possuem credenciais próprias.
              </p>
            </div>
          </div>

          <div className="space-y-6">
            <div className="bg-slate-100 dark:bg-slate-800 p-3 rounded-md text-sm border border-slate-200 dark:border-slate-700">
              <p className="font-medium mb-1">URL de Redirecionamento Autorizado:</p>
              <code className="text-xs break-all block text-muted-foreground select-all">
                {window.location.origin}/google-sheets/callback
              </code>
              <p className="text-xs mt-1 text-muted-foreground">Copie esta URL e cole no painel de credenciais do Google Cloud Console.</p>
            </div>

            <div className="space-y-4">
              <div>
                <Label htmlFor="client_id">Client ID (Chave Usuário)</Label>
                <Input
                  id="client_id"
                  value={clientId}
                  onChange={(e) => setClientId(e.target.value)}
                  placeholder="Ex: 792603653049-5ba8ss...apps.googleusercontent.com"
                  className="mt-1"
                />
              </div>

              <div>
                <Label htmlFor="client_secret">Client Secret (Chave Secreta)</Label>
                <Input
                  id="client_secret"
                  type="password"
                  value={clientSecret}
                  onChange={(e) => setClientSecret(e.target.value)}
                  placeholder="Sua chave secreta do Google"
                  className="mt-1"
                />
              </div>
            </div>

            <div className="flex justify-between items-center pt-4 border-t">
              <Button variant="destructive" onClick={handleDisconnect} disabled={!clientId && !clientSecret}>
                Desconectar / Limpar
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
