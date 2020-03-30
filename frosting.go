package frosting

import (
	"errors"
	mapset "github.com/deckarep/golang-set"
	"github.com/golang-collections/go-datastructures/queue"
)

type ClientInfo struct {
	rootNamespace *NamespaceInfo
	ingredientFns []IngredientFn
}

type Client interface {
	Execute(args ...string)
	RootNamespace() *NamespaceInfo
}

func (c *ClientInfo) Execute(args ...string) {
	// ingredients := gatherIngredients(c.rootNamespace)

	//resolver := newIngredientDependencyResolver(ingredients)
}



func gatherIngredients(ns *NamespaceInfo) []*IngredientInfo {
	ingredients := []*IngredientInfo{}

	for _, ingFn := range ns.ingredientFns {
		ingredients = append(ingredients, ingFn())
	}

	for _, ns := range ns.Children() {
		ingredients = append(ingredients, gatherIngredients(ns)...)
	}

	return ingredients
}

func (c *ClientInfo) RootNamespace() *NamespaceInfo {
	return c.rootNamespace
}

var _ Client = &ClientInfo{}

func MustNew(namespaces ...*NamespaceInfo) *ClientInfo {
	rootNs := MustNewNamespace("root", []IngredientFn{}, namespaces...)

	c := &ClientInfo{
		rootNamespace: rootNs,
	}

	return c
}
