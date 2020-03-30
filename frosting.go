package frosting

import "fmt"

type ClientInfo struct {
	rootNamespace *NamespaceInfo
	ingredientFns []IngredientFn
}

type Client interface {
	Execute(args ...string)
	RootNamespace() *NamespaceInfo
}

func (c *ClientInfo) Execute(args ...string) {
	fmt.Printf("resolving namespaces...\n")
	namespaces, err := resolveNamespaces(c.rootNamespace)
	if err != nil {
		panic(fmt.Errorf("unable to resolve namespaces: %w", err))
	}
	fmt.Printf("%d namespaces found...\n", len(namespaces))

	fmt.Printf("gathering ingredients...\n")
	ingredients := gatherIngredients(namespaces)
	fmt.Printf("%d ingredients found...\n", len(ingredients))

	resolver, err := newIngredientDependencyResolver(ingredients)
	if err != nil {
		panic(fmt.Errorf("unable to create ingredient resolver: %w", err))
	}

	fmt.Printf("resolving ingredients...\n")
	ingredientsResolved, err := resolver.Resolve()
	if err != nil {
		panic(fmt.Errorf("unable to resolve ingredients list and dependencies: %w", err))
	}
	fmt.Printf("%d resolved ingredients...\n", len(ingredientsResolved))

	for _, ing := range ingredientsResolved {
		fmt.Printf("%s:%s\n", ing.namespace.name, ing.Name)
	}
}

func gatherIngredients(namespaces []*NamespaceInfo) []*IngredientInfo {
	ingredients := []*IngredientInfo{}

	for _, ns := range namespaces {
		for _, ingFn := range ns.ingredientFns {
			ingredient := ingFn()

			// set where we found the ingredient
			ingredient.namespace = ns
			ingredients = append(ingredients, ingredient)
		}
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
