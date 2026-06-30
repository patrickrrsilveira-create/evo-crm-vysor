package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

// Config representa a estrutura JSON da coluna config do banco
type AgentConfig struct {
	SubAgents []SubAgent `json:"sub_agents"`
}

type SubAgent struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Role      string `json:"role"`
	Goal      string `json:"goal"`
	Backstory string `json:"backstory"`
}

// ParseAgentToDAG lê a configuração JSON de um agente e constrói o Grafo DAG.
func ParseAgentToDAG(agent *models.Agent) (*Workflow, error) {
	if agent.Type == "llm" || agent.Type == "task" {
		// Agente simples de 1 nó
		wf := NewWorkflow(agent.ID.String())
		node := &Node{
			ID:   agent.ID.String(),
			Type: NodeAgent,
		}
		wf.AddNode(node)
		wf.StartID = node.ID
		return wf, nil
	}

	// Faz o Parse JSON da configuração para extrair sub_agents
	var config AgentConfig
	rawJSON, err := agent.Config.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("falha ao ler config json: %v", err)
	}

	if err := json.Unmarshal(rawJSON, &config); err != nil {
		return nil, fmt.Errorf("falha ao extrair sub_agents do JSON: %v", err)
	}

	wf := NewWorkflow(agent.ID.String())

	// Adiciona todos os nós (SubAgentes) na DAG
	for _, sub := range config.SubAgents {
		node := &Node{
			ID:   sub.ID,
			Type: NodeAgent,
		}
		wf.AddNode(node)
	}

	// Constrói as arestas (Edges) baseado no tipo de roteamento
	switch agent.Type {
	case "sequential":
		// Conecta A -> B -> C -> D
		for i := 0; i < len(config.SubAgents)-1; i++ {
			wf.AddEdge(config.SubAgents[i].ID, config.SubAgents[i+1].ID)
		}
		if len(config.SubAgents) > 0 {
			wf.StartID = config.SubAgents[0].ID
		}

	case "parallel":
		// Adiciona um JoinNode no final
		joinNodeID := "JOIN-" + agent.ID.String()
		wf.AddNode(&Node{ID: joinNodeID, Type: NodeJoin})

		// Começamos por um Dummy Start para ramificar (ou assumimos o Worker processa todos)
		// Para simplificar, o primeiro SubAgent inicia a todos em paralelo
		// Mas em DAG puro, teríamos um nó raiz que aponta para todos.
		startNodeID := "START-" + agent.ID.String()
		wf.AddNode(&Node{ID: startNodeID, Type: NodeCondition})
		wf.StartID = startNodeID

		for _, sub := range config.SubAgents {
			wf.AddEdge(startNodeID, sub.ID)
			wf.AddEdge(sub.ID, joinNodeID)
		}

	case "a2a", "loop":
		// TODO: Lógica complexa com Conditionals dinâmicos (Agents decidem quem chamar)
		return nil, fmt.Errorf("roteamento %s requer motor dinâmico e LLM Routing (ainda não implementado)", agent.Type)
	}

	return wf, nil
}
