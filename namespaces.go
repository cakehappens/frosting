package frosting

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/golang-collections/go-datastructures/queue"
)

type Namespace interface {
	Name() string
	ID() string
	Parent() *NamespaceInfo
	Children() []*NamespaceInfo
	MustAddIngredient(ing IngredientFn)
}

type NamespaceInfo struct {
	name string
	id string
	parent *NamespaceInfo
	children mapset.Set
	ingredientFns []IngredientFn
}

func (n *NamespaceInfo) Name() string {
	return n.name
}

func (n *NamespaceInfo) Parent() *NamespaceInfo {
	return n.parent
}

func (n *NamespaceInfo) Children() []*NamespaceInfo {
	children := []*NamespaceInfo{}
	for child := range n.children.Iter() {
		children = append(children, child.(*NamespaceInfo))
	}

	return children
}

func (n *NamespaceInfo) MustAddIngredient(ing IngredientFn) {
	if ing == nil {
		panic(errors.New("ingredient cannot be nil"))
	}

	n.ingredientFns = append(n.ingredientFns, ing)
}

func (n *NamespaceInfo) ID() string {
	return n.id
}

var _ Namespace = &NamespaceInfo{}

func MustNewNamespace(name string, ingredientFns []IngredientFn, children ...*NamespaceInfo) *NamespaceInfo {
	if name == "" {
		panic(errors.New("namespace name cannot be empty"))
	}

	n := &NamespaceInfo{
		name: name,
		id: newULID(),
		ingredientFns: ingredientFns,
	}

	for _, child := range children {
		n.children.Add(child)
	}

	return n
}

func resolveNamespaces(root *NamespaceInfo) ([]*NamespaceInfo, error) {
	namespaces := []*NamespaceInfo{}

	visited := make(map[string]bool)
	nsQueue := queue.New(1)

	visited[root.ID()] = true
	namespaces = append(namespaces, root)

	nsQueue.Put(root.Children())

	for nsQueue.Len() > 0 {
		dequeuedNamespaces, err := nsQueue.Get(1)
		if err != nil {
			return nil, fmt.Errorf("problem dequeuing namespace: %w", err)
		}

		ns := dequeuedNamespaces[0].(*NamespaceInfo)
		if _, ok := visited[ns.id]; ok {
			return nil, errors.New("found circular namespace reference")
		}

		namespaces = append(namespaces, ns)
		visited[ns.id] = true

		nsQueue.Put(ns.Children())
	}

	return namespaces, nil
}