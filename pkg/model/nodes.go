package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/rivo/tview"
)

var (
	// ErrReferenceNil is returned when a Node contains a nil reference
	ErrReferenceNil = errors.New("reference is nil")

	// ErrReferenceNotNodeTypes is returned when an interface is not
	// implementing the NodeTypes interface
	ErrReferenceNotNodeTypes = errors.New(
		"reference is not implementing NodeTypes",
	)
)

func resolveNodeType(intf interface{}) (NodeTypes, error) {
	if intf == nil {
		return nil, ErrReferenceNil
	}
	if nt := intf.(NodeTypes); nt != nil {
		return nt, nil
	}

	return nil, ErrReferenceNotNodeTypes
}

func ingestFor(e *Event, tn *tview.TreeNode) error {
	nt, err := resolveNodeType(tn.GetReference())
	if err != nil {
		nt = &NodeRoot{}
		tn.SetReference(nt)
	}
	children := tn.GetChildren()
	for pos := range children {
		c := children[pos].GetReference().(NodeTypes)
		if c == nil {
			return fmt.Errorf("error: unexpected reference")
		}
		if c.Match(e) {
			return ingestFor(e, children[pos])
		}
	}

	newTreeNode := tview.NewTreeNode("")
	newNodeType := nt.Next(e)
	if newNodeType == nil {
		log.Printf("")
		return nil
	}
	newTreeNode.SetReference(newNodeType)
	newNodeType.Label(newTreeNode)
	tn.AddChild(newTreeNode)
	return ingestFor(e, newTreeNode)
}

type NodeTypes interface {
	Next(e *Event) NodeTypes
	Match(e *Event) bool
	Label(*tview.TreeNode)
}

type NodeRoot struct {
	*tview.TreeNode
}

func (n *NodeRoot) Ingest(e *Event) error {
	return ingestFor(e, n.TreeNode)
}

func (*NodeRoot) Match(e *Event) bool {
	return true
}

func (n *NodeRoot) Next(e *Event) NodeTypes {
	return &NodeAPIVersion{
		APIVersion: e.Object.APIVersion,
		parent:     n,
	}
}

func (n *NodeRoot) Label(t *tview.TreeNode) {
	t.SetText("*")
}

type NodeAPIVersion struct {
	APIVersion string
	parent     NodeTypes
}

func (n *NodeAPIVersion) Next(e *Event) NodeTypes {
	return &NodeKind{
		Kind:   e.Object.Kind,
		parent: n,
	}
}

func (n *NodeAPIVersion) Match(e *Event) bool {
	return n.APIVersion == e.Object.APIVersion
}

func (n *NodeAPIVersion) Label(t *tview.TreeNode) {
	t.SetText(n.APIVersion)
}

type NodeKind struct {
	Kind   string
	parent NodeTypes
}

func (n *NodeKind) Next(e *Event) NodeTypes {
	if e.Object.Metadata.Namespace != "" {
		return &NodeNamespace{
			Namespace: e.Object.Metadata.Namespace,
			parent:    n,
		}
	}
	return &NodeName{
		Name:   e.Object.Metadata.Name,
		parent: n,
	}
}

func (n *NodeKind) Match(e *Event) bool {
	return n.Kind == e.Object.Kind
}

func (n *NodeKind) Label(t *tview.TreeNode) {
	t.SetText(n.Kind)
}

type NodeNamespace struct {
	Namespace string
	parent    NodeTypes
}

func (n *NodeNamespace) Next(e *Event) NodeTypes {
	return &NodeName{
		Name:   e.Object.Metadata.Name,
		parent: n,
	}
}
func (n *NodeNamespace) Match(e *Event) bool {
	return n.Namespace == e.Object.Metadata.Namespace
}

func (n *NodeNamespace) Label(t *tview.TreeNode) {
	t.SetText(n.Namespace)
}

type NodeName struct {
	Name     string
	parent   NodeTypes
	Versions []*Event
}

func (n *NodeName) Next(e *Event) NodeTypes {
	return &NodeResourceVersion{
		ResourceVersion: e.Object.Metadata.ResourceVersion,
		parent:          n,
	}
}

func (n *NodeName) Match(e *Event) bool {
	return n.Name == e.Object.Metadata.Name
}

func (n *NodeName) Label(t *tview.TreeNode) {
	t.SetText(n.Name)
}

type NodeResourceVersion struct {
	ResourceVersion string
	parent          NodeTypes
	event           *Event
}

func (n *NodeResourceVersion) Next(e *Event) NodeTypes {
	n.event = e
	return nil
}

func (n *NodeResourceVersion) Match(e *Event) bool {
	return n.ResourceVersion == e.Object.Metadata.ResourceVersion
}

func (n *NodeResourceVersion) Label(t *tview.TreeNode) {
	t.SetText(fmt.Sprintf("#%s", n.ResourceVersion))
}

func NewRootNode() *NodeRoot {
	return &NodeRoot{
		tview.NewTreeNode("").SetIndent(0),
	}
}
