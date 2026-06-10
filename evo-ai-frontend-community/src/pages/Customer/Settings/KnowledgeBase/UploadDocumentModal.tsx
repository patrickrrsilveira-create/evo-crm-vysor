import React, { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@evoapi/design-system/dialog';
import { Button } from '@evoapi/design-system/button';
import { Input } from '@evoapi/design-system/input';
import { Label } from '@evoapi/design-system/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@evoapi/design-system/tabs';
import { toast } from 'sonner';
import { Upload, Link as LinkIcon } from 'lucide-react';
import api from '@/services/core/api';

interface UploadDocumentModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  knowledgeBaseId: number;
}

export const UploadDocumentModal: React.FC<UploadDocumentModalProps> = ({
  isOpen,
  onClose,
  onSuccess,
  knowledgeBaseId
}) => {
  const [activeTab, setActiveTab] = useState<'file' | 'url'>('file');
  const [title, setTitle] = useState('');
  const [url, setUrl] = useState('');
  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error('O título é obrigatório');
      return;
    }

    setIsLoading(true);
    try {
      if (activeTab === 'file') {
        if (!file) {
          toast.error('Selecione um arquivo PDF ou texto');
          setIsLoading(false);
          return;
        }
        
        const formData = new FormData();
        formData.append('knowledge_base_id', knowledgeBaseId.toString());
        formData.append('title', title);
        formData.append('file', file);

        await api.post('/knowledge/ingest/file', formData);
      } else {
        if (!url.trim()) {
          toast.error('A URL é obrigatória');
          setIsLoading(false);
          return;
        }

        await api.post('/knowledge/ingest/url', {
          knowledge_base_id: knowledgeBaseId,
          title,
          url
        });
      }

      toast.success('Documento enviado! O processamento foi iniciado em segundo plano.');
      onSuccess();
      onClose();
      resetForm();
    } catch (error) {
      toast.error('Erro ao enviar documento. Verifique o formato ou tente novamente.');
    } finally {
      setIsLoading(false);
    }
  };

  const resetForm = () => {
    setTitle('');
    setUrl('');
    setFile(null);
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => {
      if (!open) resetForm();
      onClose();
    }}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Adicionar Conhecimento</DialogTitle>
          <DialogDescription>
            Faça upload de um arquivo ou forneça um link para treinar a IA.
          </DialogDescription>
        </DialogHeader>
        
        <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'file' | 'url')} className="mt-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="file" className="gap-2">
              <Upload className="h-4 w-4" />
              Arquivo (PDF/TXT)
            </TabsTrigger>
            <TabsTrigger value="url" className="gap-2">
              <LinkIcon className="h-4 w-4" />
              Link (Web)
            </TabsTrigger>
          </TabsList>
          
          <form onSubmit={handleSubmit} className="mt-4 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">Título do Documento</Label>
              <Input
                id="title"
                placeholder="Ex: Manual do Produto 2024"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={isLoading}
                required
              />
            </div>
            
            <TabsContent value="file" className="space-y-4 mt-0">
              <div className="space-y-2">
                <Label htmlFor="file">Arquivo</Label>
                <Input
                  id="file"
                  type="file"
                  accept=".pdf,.txt"
                  onChange={(e) => setFile(e.target.files?.[0] || null)}
                  disabled={isLoading}
                  className="cursor-pointer"
                />
                <p className="text-xs text-muted-foreground">Tamanho máximo: 10MB. Formatos: PDF, TXT.</p>
              </div>
            </TabsContent>
            
            <TabsContent value="url" className="space-y-4 mt-0">
              <div className="space-y-2">
                <Label htmlFor="url">URL da Página Web</Label>
                <Input
                  id="url"
                  type="url"
                  placeholder="https://exemplo.com/faq"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  disabled={isLoading}
                />
                <p className="text-xs text-muted-foreground">A página deve ser acessível publicamente.</p>
              </div>
            </TabsContent>

            <DialogFooter className="pt-4">
              <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading || (!file && activeTab === 'file') || (!url && activeTab === 'url')}>
                {isLoading ? 'Enviando...' : 'Adicionar'}
              </Button>
            </DialogFooter>
          </form>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
};
