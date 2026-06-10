import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@evoapi/design-system/card';
import { Button } from '@evoapi/design-system/button';

import { FileText, Link, Bot, Plus, BookOpen, Upload, Trash2 } from 'lucide-react';

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
  const [documents, setDocuments] = useState<any[]>([]);
  const [agentBots, setAgentBots] = useState<any[]>([]);
  const [isFetchingDetails, setIsFetchingDetails] = useState(false);

  const fetchBases = async (newBaseId?: number) => {
    setIsLoading(true);
    try {
      const response = await knowledgeBasesService.list();
      const data = (response as any).data || response;
      setBases(data);
      if (newBaseId) {
        const newlyCreated = data.find((b: KnowledgeBase) => b.id === newBaseId);
        if (newlyCreated) setSelectedBase(newlyCreated);
      } else if (selectedBase) {
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

  const fetchBaseDetails = async (baseId: number) => {
    setIsFetchingDetails(true);
    try {
      const [docsResponse, agentsResponse] = await Promise.all([
        knowledgeBasesService.getDocuments(baseId),
        knowledgeBasesService.getAgentBots(baseId)
      ]);
      setDocuments((docsResponse as any).data || docsResponse);
      setAgentBots((agentsResponse as any).data || agentsResponse);
    } catch (error) {
      toast.error('Erro ao carregar detalhes da base');
    } finally {
      setIsFetchingDetails(false);
    }
  };

  useEffect(() => {
    if (selectedBase) {
      fetchBaseDetails(selectedBase.id);
    } else {
      setDocuments([]);
      setAgentBots([]);
    }
  }, [selectedBase]);

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
                        <div className="flex gap-2">
                          <Button variant="outline" size="sm" onClick={() => setSelectedBase(base)}>
                            Gerenciar
                          </Button>
                          <Button variant="ghost" size="sm" className="text-red-500 hover:text-red-700 hover:bg-red-50" onClick={async () => {
                            if (window.confirm('Tem certeza que deseja excluir esta base e todos os seus documentos?')) {
                              try {
                                await knowledgeBasesService.delete(base.id);
                                toast.success('Base excluída com sucesso');
                                fetchBases();
                              } catch (e) {
                                toast.error('Erro ao excluir base');
                              }
                            }
                          }}>
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
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
              {isFetchingDetails ? (
                <div className="flex h-[400px] items-center justify-center">
                  <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
              ) : (
                <div className="flex h-[400px]">
                  <div className="w-1/2 p-6 border-r overflow-y-auto">
                    <h3 className="font-medium flex items-center gap-2 mb-4">
                      <FileText className="h-4 w-4 text-primary" />
                      Documentos ({documents.length})
                    </h3>
                    {documents.length === 0 ? (
                      <div className="text-sm text-muted-foreground text-center py-8">
                        Nenhum documento adicionado ainda.<br />
                        Clique em "Adicionar PDF / Link".
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {documents.map((doc) => (
                          <div key={doc.id} className="flex flex-col gap-1 p-3 border rounded-md">
                            <div className="flex items-center gap-2">
                              {doc.content_type === 'url' ? <Link className="h-4 w-4 text-primary" /> : <FileText className="h-4 w-4 text-primary" />}
                              <span className="font-medium text-sm truncate" title={doc.title}>{doc.title}</span>
                            </div>
                            <span className="text-xs text-muted-foreground ml-6 truncate" title={doc.file_url}>{doc.file_url}</span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                  <div className="w-1/2 p-6 overflow-y-auto">
                    <h3 className="font-medium flex items-center gap-2 mb-4">
                      <Bot className="h-4 w-4 text-primary" />
                      Agentes Vinculados ({agentBots.length})
                    </h3>
                    {agentBots.length === 0 ? (
                      <div className="text-sm text-muted-foreground text-center py-8">
                        Nenhum agente vinculado.<br />
                        Clique em "Vincular Agente".
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {agentBots.map((agent) => (
                          <div key={agent.id} className="flex items-center gap-3 p-3 border rounded-md">
                            <div className="bg-primary/10 p-2 rounded-md text-primary">
                              <Bot className="h-4 w-4" />
                            </div>
                            <div>
                              <p className="font-medium text-sm">{agent.name}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              )}
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
