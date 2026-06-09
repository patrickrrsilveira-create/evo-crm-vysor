import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@evoapi/design-system/card';
import { Button } from '@evoapi/design-system/button';

import { FileText, Link, Bot, Plus, BookOpen, Upload } from 'lucide-react';

import { KnowledgeBase, knowledgeBasesService } from '@/services/knowledgeBases';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { CreateKnowledgeBaseModal } from './CreateKnowledgeBaseModal';
import { UploadDocumentModal } from './UploadDocumentModal';
import { LinkAgentModal } from './LinkAgentModal';
const KnowledgeBasePage = () => {
  const [bases, setBases] = useState<KnowledgeBase[]>([]);
  const [selectedBase, setSelectedBase] = useState<KnowledgeBase | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [isLinkAgentModalOpen, setIsLinkAgentModalOpen] = useState(false);

  const fetchBases = async () => {
    setIsLoading(true);
    try {
      const response = await knowledgeBasesService.list();
      const data = (response as any).data || response;
      setBases(data);
      if (selectedBase) {
        const updatedSelected = data.find((b: KnowledgeBase) => b.id === selectedBase.id);
        if (updatedSelected) setSelectedBase(updatedSelected);
      }
    } catch (error) {
      toast.error('Erro ao carregar bases de conhecimento');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBases();
  }, []);

  return (
    <div className="flex flex-col gap-6 p-6 max-w-6xl mx-auto w-full h-full">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Base de Conhecimento (RAG)</h1>
          <p className="text-muted-foreground mt-1">
            Treine seus Agentes de IA com PDFs, documentos e sites.
          </p>
        </div>
        <Button className="gap-2" onClick={() => setIsCreateModalOpen(true)}>
          <Plus className="h-4 w-4" />
          Nova Base
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {!selectedBase ? (
          <>
            <Card className="col-span-2">
              <CardHeader>
                <CardTitle>Suas Bases de Dados</CardTitle>
                <CardDescription>Bases de conhecimento criadas para consulta das IAs.</CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <div className="flex justify-center p-8">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                  </div>
                ) : bases.length === 0 ? (
                  <div className="flex flex-col items-center justify-center p-8 text-center border-2 border-dashed rounded-lg border-muted">
                    <BookOpen className="h-12 w-12 text-muted-foreground mb-4 opacity-50" />
                    <h3 className="text-lg font-medium">Nenhuma Base Encontrada</h3>
                    <p className="text-sm text-muted-foreground mt-1 mb-4">
                      Crie uma base de conhecimento para fazer o upload de PDFs.
                    </p>
                    <Button variant="outline" className="gap-2" onClick={() => setIsCreateModalOpen(true)}>
                      <Plus className="h-4 w-4" />
                      Criar Primeira Base
                    </Button>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {bases.map((base) => (
                      <div key={base.id} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <div className="bg-primary/10 p-2 rounded-md text-primary">
                            <BookOpen className="h-5 w-5" />
                          </div>
                          <div>
                            <h4 className="font-medium">{base.name}</h4>
                            <p className="text-sm text-muted-foreground">
                              {base.documents_count} documentos
                            </p>
                          </div>
                        </div>
                        <Button variant="outline" size="sm" onClick={() => setSelectedBase(base)}>
                          Gerenciar
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Tipos de Documentos</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-center gap-3 text-sm">
                    <div className="bg-primary/10 p-2 rounded-md text-primary">
                      <FileText className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="font-medium">Arquivos PDF</p>
                      <p className="text-xs text-muted-foreground">Manuais, catálogos e ebooks.</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3 text-sm">
                    <div className="bg-primary/10 p-2 rounded-md text-primary">
                      <Link className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="font-medium">Links (URLs)</p>
                      <p className="text-xs text-muted-foreground">Páginas de sites e artigos.</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </>
        ) : (
          <Card className="col-span-3 flex-1 shadow-sm border border-border overflow-hidden">
            <CardHeader className="bg-muted/30 border-b pb-4 flex flex-row items-start justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Button variant="ghost" size="sm" className="-ml-3" onClick={() => setSelectedBase(null)}>
                    &larr; Voltar
                  </Button>
                  {selectedBase.name}
                </CardTitle>
                <CardDescription className="mt-1">
                  Gerencie os documentos e agentes vinculados a esta base.
                </CardDescription>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" className="gap-2" onClick={() => setIsLinkAgentModalOpen(true)}>
                  <Bot className="h-4 w-4" />
                  Vincular Agente
                </Button>
                <Button size="sm" className="gap-2" onClick={() => setIsUploadModalOpen(true)}>
                  <Upload className="h-4 w-4" />
                  Adicionar PDF / Link
                </Button>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <div className="flex h-[400px]">
                <div className="w-1/2 p-6 border-r overflow-y-auto">
                  <h3 className="font-medium flex items-center gap-2 mb-4">
                    <FileText className="h-4 w-4 text-primary" />
                    Documentos
                  </h3>
                  <div className="text-sm text-muted-foreground text-center py-8">
                    Esta seção exibirá os documentos ingeridos.<br />
                    Integração com PGVector concluída.
                  </div>
                </div>
                <div className="w-1/2 p-6 overflow-y-auto">
                  <h3 className="font-medium flex items-center gap-2 mb-4">
                    <Bot className="h-4 w-4 text-primary" />
                    Agentes Vinculados
                  </h3>
                  <div className="text-sm text-muted-foreground text-center py-8">
                    Esta seção exibirá os agentes vinculados.<br />
                    Pronto para consultar.
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
      
      <CreateKnowledgeBaseModal 
        isOpen={isCreateModalOpen} 
        onClose={() => setIsCreateModalOpen(false)} 
        onSuccess={fetchBases} 
      />

      {selectedBase && (
        <>
          <UploadDocumentModal
            isOpen={isUploadModalOpen}
            onClose={() => setIsUploadModalOpen(false)}
            onSuccess={fetchBases}
            knowledgeBaseId={selectedBase.id}
          />
          <LinkAgentModal
            isOpen={isLinkAgentModalOpen}
            onClose={() => setIsLinkAgentModalOpen(false)}
            onSuccess={fetchBases}
            knowledgeBaseId={selectedBase.id}
          />
        </>
      )}
    </div>
  );
};

export default KnowledgeBasePage;
