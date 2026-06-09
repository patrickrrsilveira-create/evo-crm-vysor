import React, { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@evoapi/design-system/dialog';
import { Button } from '@evoapi/design-system/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@evoapi/design-system/select';
import { Label } from '@evoapi/design-system/label';
import AgentBotsService from '@/services/channels/agentBotsService';
import { toast } from 'sonner';
import api from "@/services/core/api";

interface LinkAgentModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  knowledgeBaseId: number;
}

export const LinkAgentModal: React.FC<LinkAgentModalProps> = ({
  isOpen,
  onClose,
  onSuccess,
  knowledgeBaseId
}) => {
  const [agents, setAgents] = useState<any[]>([]);
  const [selectedAgentId, setSelectedAgentId] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchAgents();
    }
  }, [isOpen]);

  const fetchAgents = async () => {
    setIsFetching(true);
    try {
      const response = await AgentBotsService.getAll();
      setAgents((response as any).data || response);
    } catch (error) {
      toast.error('Erro ao carregar agentes');
    } finally {
      setIsFetching(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedAgentId) return;

    setIsLoading(true);
    try {
      await api.post(`/api/v1/knowledge_bases/${knowledgeBaseId}/agent_bots`, {
        agent_bot_id: selectedAgentId
      });
      toast.success('Agente vinculado com sucesso!');
      onSuccess();
      onClose();
      setSelectedAgentId('');
    } catch (error) {
      toast.error('Erro ao vincular agente');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Vincular Agente à Base</DialogTitle>
          <DialogDescription>
            Selecione um agente para dar a ele acesso a esta base de conhecimento.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="agent">Agente (Bot)</Label>
            <Select value={selectedAgentId} onValueChange={setSelectedAgentId} disabled={isFetching || isLoading}>
              <SelectTrigger>
                <SelectValue placeholder={isFetching ? "Carregando..." : "Selecione um agente"} />
              </SelectTrigger>
              <SelectContent>
                {agents.map((agent) => (
                  <SelectItem key={agent.id} value={agent.id.toString()}>
                    {agent.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
              Cancelar
            </Button>
            <Button type="submit" disabled={!selectedAgentId || isLoading}>
              {isLoading ? 'Vinculando...' : 'Vincular'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};
