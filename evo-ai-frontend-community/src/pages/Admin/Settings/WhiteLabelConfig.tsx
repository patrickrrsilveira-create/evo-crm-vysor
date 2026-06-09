import { useState, useEffect, useCallback, useMemo, useRef } from 'react';
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
  CardDescription,
} from '@evoapi/design-system';
import { toast } from 'sonner';
import { Loader2, UploadCloud, Trash2, Image as ImageIcon, Link as LinkIcon, Palette, Building2 } from 'lucide-react';
import { adminConfigService } from '@/services/admin/adminConfigService';
import { extractError } from '@/utils/apiHelpers';
import type { AdminConfigData } from '@/types/admin/adminConfig';

function createPremiumSchema() {
  return z.object({
    COMPANY_NAME: z.string().optional().nullable(),
    APP_LOGO_URL: z.string().optional().nullable(),
    APP_LOGO_WIDTH: z.string().optional().nullable(),
    APP_LOGO_HEIGHT: z.string().optional().nullable(),
    APP_LOGIN_LOGO_URL: z.string().optional().nullable(),
    APP_LOGIN_LOGO_WIDTH: z.string().optional().nullable(),
    APP_LOGIN_LOGO_HEIGHT: z.string().optional().nullable(),
    SIDEBAR_COPYRIGHT_TEXT: z.string().optional().nullable(),
    SUPPORT_LINK: z.string().optional().nullable(),
    DOCS_LINK: z.string().optional().nullable(),
  });
}

type PremiumFormData = z.infer<ReturnType<typeof createPremiumSchema>>;

const DEFAULTS: PremiumFormData = {
  COMPANY_NAME: '',
  APP_LOGO_URL: '',
  APP_LOGO_WIDTH: '',
  APP_LOGO_HEIGHT: '',
  APP_LOGIN_LOGO_URL: '',
  APP_LOGIN_LOGO_WIDTH: '',
  APP_LOGIN_LOGO_HEIGHT: '',
  SIDEBAR_COPYRIGHT_TEXT: '',
  SUPPORT_LINK: '',
  DOCS_LINK: '',
};

function buildFormValues(data: Record<string, unknown>): PremiumFormData {
  return {
    COMPANY_NAME: (data.COMPANY_NAME as string) ?? '',
    APP_LOGO_URL: (data.APP_LOGO_URL as string) ?? '',
    APP_LOGO_WIDTH: (data.APP_LOGO_WIDTH as string) ?? '',
    APP_LOGO_HEIGHT: (data.APP_LOGO_HEIGHT as string) ?? '',
    APP_LOGIN_LOGO_URL: (data.APP_LOGIN_LOGO_URL as string) ?? '',
    APP_LOGIN_LOGO_WIDTH: (data.APP_LOGIN_LOGO_WIDTH as string) ?? '',
    APP_LOGIN_LOGO_HEIGHT: (data.APP_LOGIN_LOGO_HEIGHT as string) ?? '',
    SIDEBAR_COPYRIGHT_TEXT: (data.SIDEBAR_COPYRIGHT_TEXT as string) ?? '',
    SUPPORT_LINK: (data.SUPPORT_LINK as string) ?? '',
    DOCS_LINK: (data.DOCS_LINK as string) ?? '',
  };
}

const MAX_IMAGE_WIDTH = 800;

export default function PremiumConfig() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const schema = useMemo(() => createPremiumSchema(), []);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
  } = useForm<PremiumFormData>({
    resolver: zodResolver(schema),
    defaultValues: DEFAULTS,
  });

  const currentLogoUrl = watch('APP_LOGO_URL');
  const currentLoginLogoUrl = watch('APP_LOGIN_LOGO_URL');

  const loadConfig = useCallback(async () => {
    setLoading(true);
    try {
      const data = await adminConfigService.getConfig('white_label');
      reset(buildFormValues(data));
    } catch {
      toast.error('Erro ao carregar configurações Premium');
    } finally {
      setLoading(false);
    }
  }, [reset]);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      toast.error('Por favor, selecione um arquivo de imagem válido.');
      return;
    }

    const reader = new FileReader();
    reader.onload = (event) => {
      const img = new Image();
      img.onload = () => {
        let width = img.width;
        let height = img.height;

        if (width > MAX_IMAGE_WIDTH) {
          height = Math.round((height * MAX_IMAGE_WIDTH) / width);
          width = MAX_IMAGE_WIDTH;
        }

        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext('2d');
        if (ctx) {
          ctx.drawImage(img, 0, 0, width, height);
          const compressedBase64 = canvas.toDataURL(file.type, 0.9);
          setValue('APP_LOGO_URL', compressedBase64, { shouldDirty: true });
        }
      };
      img.src = event.target?.result as string;
    };
    reader.readAsDataURL(file);
  };

  const handleLoginImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      toast.error('Por favor, selecione um arquivo de imagem válido.');
      return;
    }

    const reader = new FileReader();
    reader.onload = (event) => {
      const img = new Image();
      img.onload = () => {
        let width = img.width;
        let height = img.height;

        if (width > MAX_IMAGE_WIDTH) {
          height = Math.round((height * MAX_IMAGE_WIDTH) / width);
          width = MAX_IMAGE_WIDTH;
        }

        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext('2d');
        if (ctx) {
          ctx.drawImage(img, 0, 0, width, height);
          const compressedBase64 = canvas.toDataURL(file.type, 0.9);
          setValue('APP_LOGIN_LOGO_URL', compressedBase64, { shouldDirty: true });
        }
      };
      img.src = event.target?.result as string;
    };
    reader.readAsDataURL(file);
  };

  const removeImage = () => {
    setValue('APP_LOGO_URL', '', { shouldDirty: true });
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const removeLoginImage = () => {
    setValue('APP_LOGIN_LOGO_URL', '', { shouldDirty: true });
    // Note: We'll need a separate ref if we wanted to clear the input, 
    // but the value is managed via react-hook-form anyway.
  };

  const onSubmit = async (formData: PremiumFormData) => {
    setSaving(true);
    try {
      const payload: Record<string, unknown> = {
        COMPANY_NAME: formData.COMPANY_NAME || '',
        APP_LOGO_URL: formData.APP_LOGO_URL || '',
        APP_LOGO_WIDTH: formData.APP_LOGO_WIDTH || '',
        APP_LOGO_HEIGHT: formData.APP_LOGO_HEIGHT || '',
        APP_LOGIN_LOGO_URL: formData.APP_LOGIN_LOGO_URL || '',
        APP_LOGIN_LOGO_WIDTH: formData.APP_LOGIN_LOGO_WIDTH || '',
        APP_LOGIN_LOGO_HEIGHT: formData.APP_LOGIN_LOGO_HEIGHT || '',
        SIDEBAR_COPYRIGHT_TEXT: formData.SIDEBAR_COPYRIGHT_TEXT || '',
        SUPPORT_LINK: formData.SUPPORT_LINK || '',
        DOCS_LINK: formData.DOCS_LINK || '',
      };

      const data = await adminConfigService.saveConfig('white_label', payload as AdminConfigData);
      reset(buildFormValues(data));
      toast.success('Configurações Premium ativadas! Recarregando a plataforma...');
      
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
    <div className="max-w-4xl mx-auto space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="mb-8">
        <h2 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60 tracking-tight">
          Configurações Premium
        </h2>
        <p className="text-muted-foreground mt-2 text-lg">
          Personalize a identidade visual e os links da sua plataforma CRM para uma experiência única.
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          
          {/* IDENTIDADE VISUAL */}
          <Card className="border-sidebar-border/50 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
            <div className="h-1 w-full bg-gradient-to-r from-blue-500 to-indigo-500" />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-xl">
                <Palette className="h-5 w-5 text-indigo-500" /> Identidade Visual
              </CardTitle>
              <CardDescription>Defina o nome da sua marca</CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="COMPANY_NAME" className="flex items-center gap-2">
                  <Building2 className="h-4 w-4 text-muted-foreground" /> Nome da Empresa
                </Label>
                <Input
                  id="COMPANY_NAME"
                  placeholder="Ex: Minha Empresa"
                  className="bg-background/50 focus:bg-background transition-colors"
                  {...register('COMPANY_NAME')}
                />
              </div>
            </CardContent>
          </Card>

          {/* LOGOMARCA */}
          <Card className="border-sidebar-border/50 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
            <div className="h-1 w-full bg-gradient-to-r from-emerald-400 to-teal-500" />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-xl">
                <ImageIcon className="h-5 w-5 text-teal-500" /> Logo do Dashboard
              </CardTitle>
              <CardDescription>Logo exibida na barra lateral do painel</CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="space-y-3">
                <Label>Upload da Imagem</Label>
                {currentLogoUrl ? (
                  <div className="relative group rounded-lg border-2 border-dashed border-border bg-background/50 p-4 flex flex-col items-center justify-center gap-4 transition-all hover:bg-accent/50">
                    <img src={currentLogoUrl} alt="Logo Preview" className="max-h-24 max-w-full object-contain" />
                    <Button 
                      type="button" 
                      variant="destructive" 
                      size="sm" 
                      onClick={removeImage}
                      className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ) : (
                  <label className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-border rounded-lg cursor-pointer bg-background/50 hover:bg-accent/50 transition-colors group">
                    <div className="flex flex-col items-center justify-center pt-5 pb-6">
                      <UploadCloud className="w-8 h-8 mb-3 text-muted-foreground group-hover:text-primary transition-colors" />
                      <p className="mb-1 text-sm text-muted-foreground font-medium"><span className="text-primary font-semibold">Clique para enviar</span> ou arraste</p>
                      <p className="text-xs text-muted-foreground/70">PNG, JPG ou SVG (Máx. 800px)</p>
                    </div>
                    <input 
                      type="file" 
                      className="hidden" 
                      accept="image/*" 
                      onChange={handleImageUpload} 
                      ref={fileInputRef}
                    />
                  </label>
                )}
              </div>

              <div className="grid grid-cols-2 gap-4 pt-2">
                <div className="space-y-2">
                  <Label htmlFor="APP_LOGO_WIDTH" className="text-xs text-muted-foreground">Largura (px)</Label>
                  <Input
                    id="APP_LOGO_WIDTH"
                    type="number"
                    placeholder="Ex: 150"
                    className="bg-background/50 focus:bg-background"
                    {...register('APP_LOGO_WIDTH')}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="APP_LOGO_HEIGHT" className="text-xs text-muted-foreground">Altura (px)</Label>
                  <Input
                    id="APP_LOGO_HEIGHT"
                    type="number"
                    placeholder="Ex: 40"
                    className="bg-background/50 focus:bg-background"
                    {...register('APP_LOGO_HEIGHT')}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* LOGO DE LOGIN */}
          <Card className="border-sidebar-border/50 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
            <div className="h-1 w-full bg-gradient-to-r from-blue-400 to-cyan-500" />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-xl">
                <ImageIcon className="h-5 w-5 text-cyan-500" /> Logo de Login
              </CardTitle>
              <CardDescription>Logo exclusiva para a tela de autenticação</CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="space-y-3">
                <Label>Upload da Imagem</Label>
                {currentLoginLogoUrl ? (
                  <div className="relative group rounded-lg border-2 border-dashed border-border bg-background/50 p-4 flex flex-col items-center justify-center gap-4 transition-all hover:bg-accent/50">
                    <img src={currentLoginLogoUrl} alt="Login Logo Preview" className="max-h-24 max-w-full object-contain" />
                    <Button 
                      type="button" 
                      variant="destructive" 
                      size="sm" 
                      onClick={removeLoginImage}
                      className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ) : (
                  <label className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-border rounded-lg cursor-pointer bg-background/50 hover:bg-accent/50 transition-colors group">
                    <div className="flex flex-col items-center justify-center pt-5 pb-6">
                      <UploadCloud className="w-8 h-8 mb-3 text-muted-foreground group-hover:text-primary transition-colors" />
                      <p className="mb-1 text-sm text-muted-foreground font-medium"><span className="text-primary font-semibold">Clique para enviar</span> ou arraste</p>
                      <p className="text-xs text-muted-foreground/70">PNG, JPG ou SVG (Máx. 800px)</p>
                    </div>
                    <input 
                      type="file" 
                      className="hidden" 
                      accept="image/*" 
                      onChange={handleLoginImageUpload} 
                    />
                  </label>
                )}
              </div>

              <div className="grid grid-cols-2 gap-4 pt-2">
                <div className="space-y-2">
                  <Label htmlFor="APP_LOGIN_LOGO_WIDTH" className="text-xs text-muted-foreground">Largura (px)</Label>
                  <Input
                    id="APP_LOGIN_LOGO_WIDTH"
                    type="number"
                    placeholder="Ex: 250"
                    className="bg-background/50 focus:bg-background"
                    {...register('APP_LOGIN_LOGO_WIDTH')}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="APP_LOGIN_LOGO_HEIGHT" className="text-xs text-muted-foreground">Altura (px)</Label>
                  <Input
                    id="APP_LOGIN_LOGO_HEIGHT"
                    type="number"
                    placeholder="Ex: 80"
                    className="bg-background/50 focus:bg-background"
                    {...register('APP_LOGIN_LOGO_HEIGHT')}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* RODAPÉ E LINKS */}
          <Card className="border-sidebar-border/50 shadow-sm hover:shadow-md transition-shadow overflow-hidden md:col-span-2">
            <div className="h-1 w-full bg-gradient-to-r from-violet-500 to-fuchsia-500" />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-xl">
                <LinkIcon className="h-5 w-5 text-fuchsia-500" /> Rodapé & Links
              </CardTitle>
              <CardDescription>Personalize os textos e URLs exibidos no rodapé do sistema</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="SIDEBAR_COPYRIGHT_TEXT">Texto de Copyright</Label>
                  <Input
                    id="SIDEBAR_COPYRIGHT_TEXT"
                    placeholder="© 2026 Minha Empresa. Todos os direitos reservados."
                    className="bg-background/50 focus:bg-background transition-colors"
                    {...register('SIDEBAR_COPYRIGHT_TEXT')}
                  />
                  <p className="text-xs text-muted-foreground">O texto legal que aparece no final da barra lateral.</p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="SUPPORT_LINK">Link de Suporte</Label>
                  <Input
                    id="SUPPORT_LINK"
                    placeholder="https://api.whatsapp.com/send?phone=..."
                    className="bg-background/50 focus:bg-background transition-colors"
                    {...register('SUPPORT_LINK')}
                  />
                  <p className="text-xs text-muted-foreground">URL para onde o botão "Preciso de suporte" deve apontar.</p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="DOCS_LINK">Link da Documentação</Label>
                  <Input
                    id="DOCS_LINK"
                    placeholder="https://docs.minhaempresa.com"
                    className="bg-background/50 focus:bg-background transition-colors"
                    {...register('DOCS_LINK')}
                  />
                  <p className="text-xs text-muted-foreground">URL para onde o botão "Documentação" deve apontar.</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="flex justify-end gap-4 pt-4 pb-12">
          <Button 
            type="button" 
            variant="outline"
            disabled={saving} 
            size="lg"
            onClick={() => {
              reset(DEFAULTS);
              toast.info('Configurações redefinidas para o padrão. Não esqueça de Salvar.');
            }}
          >
            Restaurar Padrões
          </Button>

          <Button 
            type="submit" 
            disabled={saving} 
            size="lg"
            className="w-full md:w-auto shadow-lg shadow-primary/20 hover:shadow-xl hover:shadow-primary/30 transition-all font-medium px-8"
          >
            {saving ? (
              <>
                <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                Aplicando Magia Premium...
              </>
            ) : (
              'Salvar Configurações Premium'
            )}
          </Button>
        </div>
      </form>
    </div>
  );
}
