import { useState, useEffect, useCallback, useMemo } from 'react';
import { useForm } from 'react-hook-form';
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
} from '@evoapi/design-system';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';
import { adminConfigService } from '@/services/admin/adminConfigService';
import { extractError } from '@/utils/apiHelpers';
import type { AdminConfigData } from '@/types/admin/adminConfig';

function createWhiteLabelSchema() {
  return z.object({
    COMPANY_NAME: z.string().optional().nullable(),
    APP_LOGO_URL: z.string().url('Deve ser uma URL válida').optional().nullable().or(z.literal('')),
    APP_PRIMARY_COLOR: z.string().regex(/^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/, 'Deve ser um código HEX (ex: #0044FF)').optional().nullable().or(z.literal('')),
  });
}

type WhiteLabelFormData = z.infer<ReturnType<typeof createWhiteLabelSchema>>;

const DEFAULTS: WhiteLabelFormData = {
  COMPANY_NAME: '',
  APP_LOGO_URL: '',
  APP_PRIMARY_COLOR: '',
};

function buildFormValues(data: Record<string, unknown>): WhiteLabelFormData {
  return {
    COMPANY_NAME: (data.COMPANY_NAME as string) ?? '',
    APP_LOGO_URL: (data.APP_LOGO_URL as string) ?? '',
    APP_PRIMARY_COLOR: (data.APP_PRIMARY_COLOR as string) ?? '',
  };
}

export default function WhiteLabelConfig() {
  // const { t } = useLanguage('adminSettings');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const schema = useMemo(() => createWhiteLabelSchema(), []);

  const {
    register,
    handleSubmit,
    reset,
  } = useForm<WhiteLabelFormData>({
    resolver: zodResolver(schema),
    defaultValues: DEFAULTS,
  });

  const loadConfig = useCallback(async () => {
    setLoading(true);
    try {
      const data = await adminConfigService.getConfig('white_label');
      reset(buildFormValues(data));
    } catch {
      toast.error('Erro ao carregar configurações de White-Label');
    } finally {
      setLoading(false);
    }
  }, [reset]);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const onSubmit = async (formData: WhiteLabelFormData) => {
    setSaving(true);
    try {
      const payload: Record<string, unknown> = {
        COMPANY_NAME: formData.COMPANY_NAME || '',
        APP_LOGO_URL: formData.APP_LOGO_URL || '',
        APP_PRIMARY_COLOR: formData.APP_PRIMARY_COLOR || '',
      };

      const data = await adminConfigService.saveConfig('white_label', payload as AdminConfigData);
      reset(buildFormValues(data));
      toast.success('Configurações salvas! Recarregue a página para aplicar a nova marca.');
      
      // Auto-reload to apply global CSS injected variables
      setTimeout(() => window.location.reload(), 1500);
    } catch (error) {
      const errorInfo = extractError(error);
      toast.error('Erro ao salvar as configurações', {
        description: errorInfo.message,
      });
    } finally {
      setSaving(false);
    }
  };

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
        <h2 className="text-xl font-semibold text-sidebar-foreground">Personalização da Marca</h2>
        <p className="text-sm text-sidebar-foreground/70 mt-1">Configure o White-Label do CRM. Altere a logo, cores e nome do sistema.</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Aparência Global</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            <div className="space-y-2">
              <Label htmlFor="COMPANY_NAME">Nome da Empresa</Label>
              <Input
                id="COMPANY_NAME"
                placeholder="Ex: Vysor Tech"
                {...register('COMPANY_NAME')}
              />
              <p className="text-xs text-muted-foreground">O nome que aparecerá no cabeçalho das abas e no lugar de Evo CRM.</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="APP_LOGO_URL">URL da Logo (PNG, SVG, JPG)</Label>
              <Input
                id="APP_LOGO_URL"
                placeholder="https://exemplo.com/sua-logo.png"
                {...register('APP_LOGO_URL')}
              />
              <p className="text-xs text-muted-foreground">URL pública da sua marca. Ficará no cabeçalho e na tela de Login.</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="APP_PRIMARY_COLOR">Cor Primária (HEX)</Label>
              <Input
                id="APP_PRIMARY_COLOR"
                placeholder="#0044FF"
                {...register('APP_PRIMARY_COLOR')}
              />
              <p className="text-xs text-muted-foreground">A cor dos botões e destaques principais (em formato Hexadecimal).</p>
            </div>

            <div className="pt-2">
              <Button type="submit" disabled={saving}>
                {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {saving ? 'Salvando...' : 'Salvar Alterações'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
