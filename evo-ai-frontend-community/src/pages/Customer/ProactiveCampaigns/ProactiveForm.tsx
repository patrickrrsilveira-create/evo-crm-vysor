import { useState, useEffect } from 'react';
import { Button, Input, Label, Textarea, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@evoapi/design-system';
import { Save, ArrowLeft, UploadCloud, MessageSquare, Clock, Target, Bot } from 'lucide-react';
import { useNavigate, useParams } from 'react-router-dom';
import { proactiveService, ProactiveCampaign } from '@/services/proactive/proactiveService';
import { labelsService } from '@/services/contacts/labelsService';
import AgentsService from '@/services/channels/agentsService';
import { AgentChannel } from '@/types/channels/inbox';
import { Label as ContactLabel } from '@/types/settings';
export default function ProactiveForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditing = !!id;

  const [agents, setAgents] = useState<AgentChannel[]>([]);
  const [labels, setLabels] = useState<ContactLabel[]>([]);

  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<Partial<ProactiveCampaign>>({
    name: '',
    trigger_type: 'LABEL_ADDED',
    trigger_target: '',
    delay_hours: 0,
    agent_id: undefined,
    message_template: '',
    status: 'ACTIVE'
  });

  useEffect(() => {
    if (isEditing) {
      fetchCampaign();
    }
    fetchOptions();
  }, [id]);

  const fetchOptions = async () => {
    try {
      const [agentsData, labelsRes] = await Promise.all([
        AgentsService.getAll(),
        labelsService.getLabels({ per_page: 100 })
      ]);
      setAgents(agentsData);
      setLabels(labelsRes.data || []);
    } catch (e) {
      console.error('Failed to fetch options', e);
    }
  };

  const fetchCampaign = async () => {
    try {
      const data = await proactiveService.getCampaign(id!);
      setFormData(data);
    } catch (error) {
      console.error('Failed to load campaign', error);
    }
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      if (isEditing) {
        await proactiveService.updateCampaign(id!, formData);
      } else {
        await proactiveService.createCampaign(formData as any);
      }
      navigate('/marketing');
    } catch (error) {
      console.error('Failed to save', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      try {
        const url = await proactiveService.uploadAttachment(file);
        setFormData(prev => ({ ...prev, attachment_url: url }));
      } catch (error) {
        console.error('Failed to upload', error);
      }
    }
  };

  return (
    <div className="flex flex-col h-full bg-background p-6 max-w-4xl mx-auto w-full">
      <div className="flex items-center gap-4 mb-8 border-b pb-4">
        <Button variant="ghost" size="icon" onClick={() => navigate('/marketing')}>
          <ArrowLeft className="w-5 h-5 text-muted-foreground" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-foreground">
            {isEditing ? 'Editar Automação' : 'Nova Automação de Leads'}
          </h1>
          <p className="text-muted-foreground text-sm">Configure as regras de envio programado.</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Coluna Esquerda: Configurações do Gatilho */}
        <div className="space-y-6 bg-card p-6 rounded-lg border shadow-sm">
          <h2 className="text-lg font-semibold flex items-center gap-2 border-b pb-2">
            <Target className="w-5 h-5 text-primary" /> Regras de Disparo
          </h2>
          
          <div className="space-y-2">
            <Label>Nome da Campanha (Interno)</Label>
            <Input 
              placeholder="Ex: Recuperar Mornos - Produto X" 
              value={formData.name}
              onChange={e => setFormData({...formData, name: e.target.value})}
            />
          </div>

          <div className="space-y-2">
            <Label>Tipo de Gatilho</Label>
            <Select 
              value={formData.trigger_type} 
              onValueChange={(val: any) => setFormData({...formData, trigger_type: val})}
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione o gatilho" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="LABEL_ADDED">Quando receber Etiqueta</SelectItem>
                <SelectItem value="PIPELINE_STAGE_ENTERED">Quando entrar no Funil</SelectItem>
                <SelectItem value="SCHEDULED_DATE">Data Agendada (Transmissão)</SelectItem>
                <SelectItem value="CONTACT_CREATED">Novo Contato Criado</SelectItem>
                <SelectItem value="CONVERSATION_OPENED">Nova Conversa Aberta</SelectItem>
                <SelectItem value="CONVERSATION_RESOLVED">Conversa Resolvida</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Alvo (Configuração do Gatilho)</Label>
            {formData.trigger_type === 'SCHEDULED_DATE' ? (
              <Input 
                type="datetime-local"
                value={formData.trigger_target}
                onChange={e => setFormData({...formData, trigger_target: e.target.value})}
              />
            ) : formData.trigger_type === 'LABEL_ADDED' ? (
              <Select 
                value={formData.trigger_target} 
                onValueChange={(val: string) => setFormData({...formData, trigger_target: val})}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a etiqueta" />
                </SelectTrigger>
                <SelectContent>
                  {labels.map(label => (
                    <SelectItem key={label.id} value={label.title}>{label.title}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            ) : formData.trigger_type === 'CONTACT_CREATED' || formData.trigger_type === 'CONVERSATION_OPENED' || formData.trigger_type === 'CONVERSATION_RESOLVED' ? (
               <Input 
                placeholder="Não exige alvo (Opcional)" 
                disabled
                value=""
              />
            ) : (
              <Input 
                placeholder="Ex: ganader_morno" 
                value={formData.trigger_target}
                onChange={e => setFormData({...formData, trigger_target: e.target.value})}
              />
            )}
          </div>

          <div className="space-y-2">
            <Label className="flex items-center gap-2">
              <Clock className="w-4 h-4" /> Tempo de Espera (Horas)
            </Label>
            <Input 
              type="number" 
              min="0"
              placeholder="Quantas horas aguardar?" 
              value={formData.delay_hours}
              onChange={e => setFormData({...formData, delay_hours: parseInt(e.target.value) || 0})}
            />
            <p className="text-xs text-muted-foreground">Ex: 72 horas = 3 dias. Coloque 0 para envio imediato (se for agendado).</p>
          </div>
        </div>

        {/* Coluna Direita: Conteúdo da Mensagem */}
        <div className="space-y-6 bg-card p-6 rounded-lg border shadow-sm">
          <h2 className="text-lg font-semibold flex items-center gap-2 border-b pb-2">
            <MessageSquare className="w-5 h-5 text-primary" /> Conteúdo
          </h2>

          <div className="space-y-2">
            <Label className="flex items-center gap-2">
              <Bot className="w-4 h-4" /> Agente Remetente
            </Label>
            <Select 
              value={formData.agent_id ? formData.agent_id.toString() : ''} 
              onValueChange={(val: string) => setFormData({...formData, agent_id: parseInt(val)})}
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione o agente" />
              </SelectTrigger>
              <SelectContent>
                {agents.map(agent => (
                  <SelectItem key={agent.id} value={agent.id.toString()}>{agent.name || `Agente ${agent.id}`}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Mensagem (Texto Livre)</Label>
            <Textarea 
              placeholder="Oi, vi que você teve interesse no produto..." 
              className="h-32 resize-none"
              value={formData.message_template}
              onChange={e => setFormData({...formData, message_template: e.target.value})}
            />
          </div>

          <div className="space-y-2 pt-2">
            <Label>Anexo (Vídeo, Imagem ou Áudio)</Label>
            <div className="flex items-center gap-4">
              <Label htmlFor="file-upload" className="cursor-pointer">
                <div className="flex items-center gap-2 px-4 py-2 bg-secondary text-secondary-foreground hover:bg-secondary/80 rounded-md transition-colors text-sm font-medium">
                  <UploadCloud className="w-4 h-4" />
                  Subir Arquivo
                </div>
              </Label>
              <Input 
                id="file-upload" 
                type="file" 
                className="hidden" 
                onChange={handleFileUpload}
              />
              {formData.attachment_url && (
                <span className="text-sm text-green-600 truncate max-w-[200px]">
                  Arquivo anexado com sucesso!
                </span>
              )}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Este arquivo será enviado junto com a mensagem.</p>
          </div>
        </div>
      </div>

      <div className="flex justify-end gap-3 mt-8">
        <Button variant="outline" onClick={() => navigate('/marketing')}>Cancelar</Button>
        <Button onClick={handleSave} disabled={loading} className="gap-2">
          <Save className="w-4 h-4" />
          {loading ? 'Salvando...' : 'Salvar Automação'}
        </Button>
      </div>
    </div>
  );
}
