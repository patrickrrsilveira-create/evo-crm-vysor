import { useState, useEffect, useRef } from 'react';
import { Image as ImageIcon, Plus, Loader2, Trash2, UploadCloud, Copy } from 'lucide-react';
import { Button } from '@evoapi/design-system/button';
import { toast } from 'sonner';
import { agentMediaService, AgentMedia } from '@/services/agentMedia';
import { formatFileSize } from '@/utils/fileUtils';

interface MediaSectionProps {
  agentId: string;
}

const MediaSection = ({ agentId }: MediaSectionProps) => {
  const [mediaList, setMediaList] = useState<AgentMedia[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isUploading, setIsUploading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<string | null>(null);
  
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const data = await agentMediaService.list(agentId);
      setMediaList(data);
    } catch (error) {
      console.error('Error fetching agent media:', error);
      toast.error('Erro ao carregar mídias do agente');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [agentId]);

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    try {
      await agentMediaService.upload(agentId, file);
      toast.success('Mídia adicionada com sucesso!');
      await fetchData();
    } catch (error) {
      console.error('Error uploading media:', error);
      toast.error('Erro ao fazer upload da mídia');
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleDelete = async (filename: string) => {
    setDeletingFile(filename);
    try {
      await agentMediaService.delete(agentId, filename);
      toast.success('Mídia removida com sucesso! Arquivo deletado do servidor.');
      setMediaList(prev => prev.filter(m => m.filename !== filename));
    } catch (error) {
      console.error('Error deleting media:', error);
      toast.error('Erro ao excluir mídia');
    } finally {
      setDeletingFile(null);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Nome do arquivo copiado!');
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between pb-4 border-b">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-primary/10">
            <ImageIcon className="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 className="text-lg font-semibold">
              Mídias (Vídeos e Imagens)
            </h3>
            <p className="text-sm text-muted-foreground">
              Faça upload de vídeos e imagens e peça para a IA enviar usando o nome do arquivo.
            </p>
          </div>
        </div>
        
        <div>
          <input 
            type="file" 
            className="hidden" 
            ref={fileInputRef}
            onChange={handleFileUpload}
            accept="image/*,video/*"
          />
          <Button 
            onClick={() => fileInputRef.current?.click()} 
            className="gap-2"
            disabled={isUploading}
          >
            {isUploading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <UploadCloud className="h-4 w-4" />
            )}
            {isUploading ? 'Enviando...' : 'Fazer Upload'}
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      ) : mediaList.length === 0 ? (
        <div className="flex flex-col items-center justify-center p-8 text-center border-2 border-dashed rounded-lg border-muted">
          <ImageIcon className="h-12 w-12 text-muted-foreground mb-4 opacity-50" />
          <h3 className="text-lg font-medium">Nenhuma Mídia Encontrada</h3>
          <p className="text-sm text-muted-foreground mt-1 mb-4">
            Faça upload de vídeos ou imagens. Depois, basta escrever no prompt da IA:<br/>
            <strong>"Envie o vídeo NOME_DO_ARQUIVO"</strong>
          </p>
          <Button variant="outline" className="gap-2" onClick={() => fileInputRef.current?.click()}>
            <UploadCloud className="h-4 w-4" />
            Fazer Primeiro Upload
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {mediaList.map((media) => {
            const isDeleting = deletingFile === media.filename;

            return (
              <div 
                key={media.filename} 
                className="flex flex-col p-4 border rounded-lg hover:border-primary/50 transition-colors bg-card"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1 min-w-0 mr-4">
                    <h4 className="font-medium truncate" title={media.filename}>
                      {media.filename}
                    </h4>
                    <p className="text-xs text-muted-foreground mt-1">
                      Tamanho: {formatBytes(media.size)}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
                    onClick={() => handleDelete(media.filename)}
                    disabled={isDeleting}
                  >
                    {isDeleting ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Trash2 className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                
                <div className="mt-4 pt-3 border-t flex items-center justify-between">
                  <p className="text-xs text-muted-foreground truncate max-w-[200px]">
                    {media.url}
                  </p>
                  <Button 
                    variant="outline" 
                    size="sm" 
                    className="h-7 text-xs gap-1"
                    onClick={() => copyToClipboard(media.filename)}
                  >
                    <Copy className="h-3 w-3" /> Copiar Nome
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default MediaSection;
