import React, { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@evoapi/design-system/dialog';
import { Button } from '@evoapi/design-system/button';
import { Input } from '@evoapi/design-system/input';
import { Label } from '@evoapi/design-system/label';
import { Textarea } from '@evoapi/design-system/textarea';
import { knowledgeBasesService } from '@/services/knowledgeBases';
import { toast } from 'sonner';

interface CreateKnowledgeBaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: (baseId?: number) => void;
}

export const CreateKnowledgeBaseModal: React.FC<CreateKnowledgeBaseModalProps> = ({
  isOpen,
  onClose,
  onSuccess
}) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsLoading(true);
    try {
      const result = await knowledgeBasesService.create({ name, description });
      const createdBase = (result as any).data || result;
      toast.success('Base de conhecimento criada com sucesso!');
      onSuccess(createdBase?.id);
      onClose();
      setName('');
      setDescription('');
    } catch (error) {
      toast.error('Erro ao criar base de conhecimento');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Nova Base de Conhecimento</DialogTitle>
          <DialogDescription>
            Crie uma base de conhecimento para treinar seus agentes de IA.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="name">Nome da Base</Label>
            <Input
              id="name"
              placeholder="Ex: Manuais Técnicos de Produtos"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={isLoading}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Descrição (Opcional)</Label>
            <Textarea
              id="description"
              placeholder="Descreva o propósito desta base..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isLoading}
            />
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
              Cancelar
            </Button>
            <Button type="submit" disabled={!name.trim() || isLoading}>
              {isLoading ? 'Criando...' : 'Criar Base'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};
