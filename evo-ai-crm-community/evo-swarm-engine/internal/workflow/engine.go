package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// NodeType define o tipo de nó no DAG
type NodeType string

const (
	NodeAgent     NodeType = "AGENT"
	NodeCondition NodeType = "CONDITION"
	NodeJoin      NodeType = "JOIN"
)

// Node representa um passo no Workflow
type Node struct {
	ID        string
	Type      NodeType
	Action    func(ctx context.Context, state map[string]interface{}) error
	NextNodes []string // IDs of nodes to execute after this one
}

// Workflow (DAG Engine) representa o grafo de execução
type Workflow struct {
	ID      string
	Nodes   map[string]*Node
	StartID string
	mu      sync.RWMutex
}

// NewWorkflow cria uma nova engine de workflow baseada em DAG
func NewWorkflow(id string) *Workflow {
	return &Workflow{
		ID:    id,
		Nodes: make(map[string]*Node),
	}
}

// AddNode adiciona um nó ao Workflow
func (w *Workflow) AddNode(node *Node) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Nodes[node.ID] = node
}

// AddEdge cria uma aresta (direcionamento) entre dois nós
func (w *Workflow) AddEdge(fromID, toID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	fromNode, exists := w.Nodes[fromID]
	if !exists {
		return fmt.Errorf("node %s not found", fromID)
	}
	if _, exists := w.Nodes[toID]; !exists {
		return fmt.Errorf("node %s not found", toID)
	}

	fromNode.NextNodes = append(fromNode.NextNodes, toID)
	return nil
}

// Execute roda o Workflow (DAG) a partir de um nó inicial, propagando o estado
func (w *Workflow) Execute(ctx context.Context, state map[string]interface{}) error {
	w.mu.RLock()
	startNode, exists := w.Nodes[w.StartID]
	w.mu.RUnlock()

	if !exists {
		return errors.New("start node not found")
	}

	return w.executeNode(ctx, startNode, state)
}

func (w *Workflow) executeNode(ctx context.Context, node *Node, state map[string]interface{}) error {
	// Executa a lógica do nó
	if node.Action != nil {
		if err := node.Action(ctx, state); err != nil {
			return fmt.Errorf("node %s failed: %w", node.ID, err)
		}
	}

	// Se for folha, acabou
	if len(node.NextNodes) == 0 {
		return nil
	}

	// TODO: Suportar execução paralela real (WaitGroups) para N nós seguintes e JoinNodes
	// Por agora, implementação simples sequencial (Depth-First)
	for _, nextID := range node.NextNodes {
		w.mu.RLock()
		nextNode := w.Nodes[nextID]
		w.mu.RUnlock()

		if err := w.executeNode(ctx, nextNode, state); err != nil {
			return err
		}
	}

	return nil
}
