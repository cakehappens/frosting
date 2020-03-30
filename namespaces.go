package frosting

import (
	"errors"
	"gopkg.in/eapache/queue.v1"

	mapset "github.com/deckarep/golang-set"
)

type Namespace interface {
	Name() string
	ID() string
	Parent() *NamespaceInfo
	Children() []*NamespaceInfo
	MustAddIngredients(ingFns ...IngredientFn)
}

type NamespaceInfo struct {
	name          string
	id            string
	parent        *NamespaceInfo
	children      mapset.Set
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

func (n *NamespaceInfo) MustAddIngredients(ing ...IngredientFn) {
	for _, ingFn := range ing {
		if ingFn == nil {
			panic(errors.New("ingredient cannot be nil"))
		}

		n.ingredientFns = append(n.ingredientFns, ingFn)
	}
}

func (n *NamespaceInfo) ID() string {
	return n.id
}

var _ Namespace = &NamespaceInfo{}

func MustNewNamespace(name string, ingredientFns []IngredientFn, children ...*NamespaceInfo) *NamespaceInfo {
	if name == "" {
		panic(errors.New("namespace name cannot be empty"))
	}

	if ingredientFns == nil {
		ingredientFns = []IngredientFn{}
	}

	n := &NamespaceInfo{
		name:          name,
		id:            newULID(),
		ingredientFns: ingredientFns,
		children:      mapset.NewSet(),
	}

	for _, child := range children {
		n.children.Add(child)
	}

	return n
}

func resolveNamespaces(root *NamespaceInfo) ([]*NamespaceInfo, error) {
	namespaces := []*NamespaceInfo{}

	visited := make(map[string]bool)
	nsQueue := queue.New()

	visited[root.ID()] = true
	namespaces = append(namespaces, root)

	for _, ns := range root.Children() {
		nsQueue.Add(ns)
	}

	for nsQueue.Length() > 0 {
		ns := nsQueue.Remove().(*NamespaceInfo)

		if _, ok := visited[ns.id]; ok {
			return nil, errors.New("found circular namespace reference")
		}

		namespaces = append(namespaces, ns)
		visited[ns.id] = true

		for _, ns := range ns.Children() {
			nsQueue.Add(ns)
		}
	}

	return namespaces, nil
}
