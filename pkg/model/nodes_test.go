package model_test

import (
	"fmt"
	"testing"

	"github.com/rivo/tview"

	"github.com/simonswine/kommissar/pkg/model"
)

type testCase struct {
	_resourceVersion int
}

func (t *testCase) resourceVersion() string {
	t._resourceVersion++
	return fmt.Sprintf("%d", t._resourceVersion)
}

func (t *testCase) newPodEvent() *model.Event {
	e := &model.Event{}
	e.Object.Kind = "Pod"
	e.Object.APIVersion = "v1"
	e.Object.Metadata.Name = "my-pod-1"
	e.Object.Metadata.Namespace = "default"
	e.Object.Metadata.ResourceVersion = t.resourceVersion()
	return e
}

func (t *testCase) newNodeEvent() *model.Event {
	e := &model.Event{}
	e.Object.Kind = "Node"
	e.Object.APIVersion = "v1"
	e.Object.Metadata.Name = "node-a"
	e.Object.Metadata.ResourceVersion = t.resourceVersion()
	return e
}

func TestNodeIngest(t *testing.T) {
	tc := &testCase{}
	n := model.NewRootNode()

	for pos, f := range []func() error{
		func() error {
			return n.Ingest(tc.newNodeEvent())
		},
		func() error {
			node := tc.newNodeEvent()
			node.Object.Metadata.Name = "a-node"
			return n.Ingest(node)
		},
		func() error {
			return n.Ingest(tc.newPodEvent())
		},
		func() error {
			return n.Ingest(tc.newPodEvent())
		},
		func() error {
			return n.Ingest(tc.newPodEvent())
		},
		func() error {
			return n.Ingest(tc.newPodEvent())
		},
	} {
		if err := f(); err != nil {
			t.Fatalf("error during example %d: %s", pos, err)
		}
	}

	tree := tview.NewTreeView().
		SetRoot(n.TreeNode).
		SetCurrentNode(n.TreeNode)

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		t.Fatalf("error creating GUI: %s", err)
	}
}
