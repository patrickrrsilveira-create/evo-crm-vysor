import { useEffect, useState } from 'react';
import { Button, Input, Table, TableBody, TableCell, TableHead, TableHeader, TableRow, Badge } from '@evoapi/design-system';
import { Plus, Edit2, Trash2, Copy, Play, Pause, Search } from 'lucide-react';
import { proactiveService, ProactiveCampaign } from '@/services/proactive/proactiveService';
import { useNavigate } from 'react-router-dom';

export default function ProactiveList() {
  const navigate = useNavigate();
  const [campaigns, setCampaigns] = useState<ProactiveCampaign[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  const fetchCampaigns = async () => {
    try {
      setLoading(true);
      // In a real scenario, handle pagination and actual API call
      const data = await proactiveService.getCampaigns();
      setCampaigns(data || []);
    } catch (error) {
      console.error('Failed to fetch proactive campaigns', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCampaigns();
  }, []);

  const handleToggleStatus = async (campaign: ProactiveCampaign) => {
    const newStatus = campaign.status === 'ACTIVE' ? 'PAUSED' : 'ACTIVE';
    try {
      await proactiveService.updateCampaign(campaign.id, { status: newStatus });
      fetchCampaigns();
    } catch (error) {
      console.error('Failed to toggle status', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Tem certeza que deseja excluir esta automação?')) {
      try {
        await proactiveService.deleteCampaign(id);
        fetchCampaigns();
      } catch (error) {
        console.error('Failed to delete', error);
      }
    }
  };

  const handleClone = async (id: string) => {
    try {
      await proactiveService.cloneCampaign(id);
      fetchCampaigns();
    } catch (error) {
      console.error('Failed to clone', error);
    }
  };

  const filteredCampaigns = campaigns.filter(c => 
    c.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getTriggerName = (type: string) => {
    const map: Record<string, string> = {
      'LABEL_ADDED': 'Etiqueta Adicionada',
      'PIPELINE_STAGE_ENTERED': 'Mudança de Funil',
      'SCHEDULED_DATE': 'Data Agendada'
    };
    return map[type] || type;
  };

  return (
    <div className="flex flex-col h-full bg-background p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Gestão de Leads & Marketing</h1>
          <p className="text-muted-foreground mt-1">
            Crie automações baseadas em tempo para recuperar leads e gerenciar o pós-venda.
          </p>
        </div>
        <Button onClick={() => navigate('/campaigns/new')} className="gap-2">
          <Plus className="w-4 h-4" />
          Nova Automação
        </Button>
      </div>

      <div className="flex items-center space-x-2 mb-4 bg-card p-4 rounded-lg border shadow-sm">
        <Search className="w-5 h-5 text-muted-foreground" />
        <Input
          placeholder="Buscar automação por nome..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="max-w-md border-0 bg-transparent focus-visible:ring-0 px-0"
        />
      </div>

      <div className="border rounded-lg bg-card overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Nome da Automação</TableHead>
              <TableHead>Gatilho</TableHead>
              <TableHead>Alvo</TableHead>
              <TableHead>Espera</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Ações</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center h-24 text-muted-foreground">
                  Carregando automações...
                </TableCell>
              </TableRow>
            ) : filteredCampaigns.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center h-24 text-muted-foreground">
                  Nenhuma automação encontrada.
                </TableCell>
              </TableRow>
            ) : (
              filteredCampaigns.map((campaign) => (
                <TableRow key={campaign.id}>
                  <TableCell className="font-medium">{campaign.name}</TableCell>
                  <TableCell>{getTriggerName(campaign.trigger_type)}</TableCell>
                  <TableCell><Badge variant="secondary">{campaign.trigger_target}</Badge></TableCell>
                  <TableCell>{campaign.delay_hours} horas</TableCell>
                  <TableCell>
                    <Badge variant={campaign.status === 'ACTIVE' ? 'default' : 'outline'} className={campaign.status === 'ACTIVE' ? 'bg-green-500 hover:bg-green-600' : ''}>
                      {campaign.status === 'ACTIVE' ? 'Ativo' : 'Pausado'}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="icon" onClick={() => handleToggleStatus(campaign)} title={campaign.status === 'ACTIVE' ? 'Pausar' : 'Ativar'}>
                        {campaign.status === 'ACTIVE' ? <Pause className="w-4 h-4 text-orange-500" /> : <Play className="w-4 h-4 text-green-500" />}
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => handleClone(campaign.id)} title="Clonar">
                        <Copy className="w-4 h-4 text-blue-500" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => navigate(`/campaigns/${campaign.id}/edit`)} title="Editar">
                        <Edit2 className="w-4 h-4 text-muted-foreground" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => handleDelete(campaign.id)} title="Excluir">
                        <Trash2 className="w-4 h-4 text-red-500" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
