import { useState, useEffect } from 'react';
import { BookOpen, Plus, Loader2, Link as LinkIcon, Unlink } from 'lucide-react';
import { useLanguage } from '@/hooks/useLanguage';
import { Button } from '@evoapi/design-system/button';
import { toast } from 'sonner';
import { KnowledgeBase, knowledgeBasesService } from '@/services/knowledgeBases';
import { CreateKnowledgeBaseModal } from '../../../Settings/KnowledgeBase/CreateKnowledgeBaseModal';

interface KnowledgeBaseSectionProps {
  agentId: string;
}

const KnowledgeBaseSection = ({ agentId }: KnowledgeBaseSectionProps) => {
  const { t } = useLanguage('aiAgents');
  const [bases, setBases] = useState<KnowledgeBase[]>([]);
  const [linkedBases, setLinkedBases] = useState<number[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLinking, setIsLinking] = useState<number | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      // Carregar todas as bases
      const basesRes = await knowledgeBasesService.list();
      const allBases = (basesRes as any).data || basesRes;
      setBases(allBases);

      // Carregar bases vinculadas a este agente
      const linked = await knowledgeBasesService.getAgentLinkedBases(agentId);
      setLinkedBases(linked.map(b => b.id));
    } catch (error) {
      console.error('Error fetching knowledge bases:', error);
      toast.error('Erro ao carregar bases de conhecimento');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [agentId]);

  const handleToggleLink = async (baseId: number, isLinked: boolean) => {
    setIsLinking(baseId);
    try {
      if (isLinked) {
        await knowledgeBasesService.unlinkAgentBot(baseId, agentId);
        setLinkedBases(prev => prev.filter(id => id !== baseId));
        toast.success('Base desvinculada com sucesso!');
      } else {
        await knowledgeBasesService.linkAgentBot(baseId, agentId);
        setLinkedBases(prev => [...prev, baseId]);
        toast.success('Base vinculada com sucesso!');
      }
    } catch (error) {
      console.error('Error toggling link:', error);
      toast.error(isLinked ? 'Erro ao desvincular base' : 'Erro ao vincular base');
    } finally {
      setIsLinking(null);
    }
  };

  const handleCreateSuccess = async (newBaseId?: number) => {
    setIsLoading(true);
    try {
      const basesRes = await knowledgeBasesService.list();
      const allBases = (basesRes as any).data || basesRes;
      setBases(allBases);

      if (newBaseId) {
        await knowledgeBasesService.linkAgentBot(newBaseId, agentId);
        setLinkedBases(prev => [...prev, newBaseId]);
        toast.success('Base criada e vinculada com sucesso!');
      } else {
        await fetchData();
      }
    } catch (error) {
      console.error('Error after creation:', error);
      toast.error('Base criada, mas houve erro ao vincular automaticamente.');
      await fetchData();
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between pb-4 border-b">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-primary/10">
            <BookOpen className="h-5 w-5 text-primary" />
          </div>
          <div>
            <h3 className="text-lg font-semibold">
              {t('knowledge.title') || 'Base de Conhecimento'}
            </h3>
            <p className="text-sm text-muted-foreground">
              {t('knowledge.subtitle') || 'Conecte documentos e informações para o agente consultar (RAG)'}
            </p>
          </div>
        </div>
        <Button onClick={() => setIsCreateModalOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          {t('knowledge.createNew') || 'Nova Base'}
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      ) : bases.length === 0 ? (
        <div className="flex flex-col items-center justify-center p-8 text-center border-2 border-dashed rounded-lg border-muted">
          <BookOpen className="h-12 w-12 text-muted-foreground mb-4 opacity-50" />
          <h3 className="text-lg font-medium">Nenhuma Base Encontrada</h3>
          <p className="text-sm text-muted-foreground mt-1 mb-4">
            Você ainda não criou nenhuma base de conhecimento.
          </p>
          <Button variant="outline" className="gap-2" onClick={() => setIsCreateModalOpen(true)}>
            <Plus className="h-4 w-4" />
            Criar Primeira Base
          </Button>
        </div>
      ) : (
        <div className="space-y-4">
          <p className="text-sm font-medium mb-2">Selecione as bases que este agente terá acesso:</p>
          {bases.map((base) => {
            const isLinked = linkedBases.includes(base.id);
            const isProcessing = isLinking === base.id;

            return (
              <div 
                key={base.id} 
                className={`flex items-center justify-between p-4 border rounded-lg transition-colors ${
                  isLinked ? 'border-primary bg-primary/5' : 'hover:border-primary/50'
                }`}
              >
                <div className="flex items-center gap-3">
                  <div className={`p-2 rounded-md ${isLinked ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'}`}>
                    <BookOpen className="h-5 w-5" />
                  </div>
                  <div>
                    <h4 className="font-medium">{base.name}</h4>
                    <p className="text-sm text-muted-foreground">
                      {base.documents_count} documentos
                    </p>
                  </div>
                </div>
                <div>
                  <Button
                    variant={isLinked ? 'destructive' : 'default'}
                    size="sm"
                    className="gap-2 w-[140px]"
                    onClick={() => handleToggleLink(base.id, isLinked)}
                    disabled={isProcessing}
                  >
                    {isProcessing ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : isLinked ? (
                      <>
                        <Unlink className="h-4 w-4" /> Desvincular
                      </>
                    ) : (
                      <>
                        <LinkIcon className="h-4 w-4" /> Vincular
                      </>
                    )}
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <CreateKnowledgeBaseModal 
        isOpen={isCreateModalOpen} 
        onClose={() => setIsCreateModalOpen(false)} 
        onSuccess={handleCreateSuccess} 
      />
    </div>
  );
};

export default KnowledgeBaseSection;
