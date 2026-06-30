package workflow_test

import (
	"context"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/workflow"
	"testing"
)

func TestWorkflowDAGExecution(t *testing.T) {
	wf := workflow.NewWorkflow("wf-test-1")

	// Estado compartilhado
	state := map[string]interface{}{
		"count": 0,
	}

	// Criar nós
	nodeA := &workflow.Node{
		ID:   "A",
		Type: workflow.NodeAgent,
		Action: func(ctx context.Context, s map[string]interface{}) error {
			s["count"] = s["count"].(int) + 1
			return nil
		},
	}

	nodeB := &workflow.Node{
		ID:   "B",
		Type: workflow.NodeCondition,
		Action: func(ctx context.Context, s map[string]interface{}) error {
			s["count"] = s["count"].(int) + 10
			return nil
		},
	}

	wf.AddNode(nodeA)
	wf.AddNode(nodeB)

	// Aresta A -> B
	err := wf.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Erro ao adicionar aresta: %v", err)
	}

	wf.StartID = "A"

	// Executar
	err = wf.Execute(context.Background(), state)
	if err != nil {
		t.Fatalf("Erro ao executar DAG: %v", err)
	}

	// count deveria ser 11 (A = +1, B = +10)
	if state["count"] != 11 {
		t.Errorf("Esperado count=11, recebido %d", state["count"])
	}
}
