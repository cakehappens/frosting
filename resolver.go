package frosting

import (
	"errors"
	mapset "github.com/deckarep/golang-set"
)

type ingredientDependencyResolver struct {
	ingredients map[string]*Ingredient
	resolved    []*Ingredient
}

func newIngredientDependencyResolver(ingredientList []*Ingredient) (*ingredientDependencyResolver, error) {
	r := &ingredientDependencyResolver{}

	r.ingredients = make(map[string]*Ingredient)
	r.resolved = []*Ingredient{}

	for _, ingredient := range ingredientList {
		r.ingredients[ingredient.name] = ingredient
	}

	return r, nil
}

func (r *ingredientDependencyResolver) Resolve() ([]*Ingredient, error) {
	// if we've already resolved, don't re-resolve
	if len(r.resolved) > 0 {
		return r.resolved, nil
	}

	ingredientDependencies := make(map[string]mapset.Set)

	// Populate the map
	for _, ingredient := range r.ingredients {
		ingredientDependencies[ingredient.name] = ingredient.dependencies
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	for len(ingredientDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for id, deps := range ingredientDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(id)
			}
		}

		// If there aren't any ready nodes, then we have a circular dependency
		if readySet.Cardinality() == 0 {
			var g []*Ingredient
			for id := range ingredientDependencies {
				g = append(g, r.ingredients[id])
			}

			return g, errors.New("circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		for id := range readySet.Iter() {
			delete(ingredientDependencies, id.(string))
			r.resolved = append(r.resolved, r.ingredients[id.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for id, deps := range ingredientDependencies {
			diff := deps.Difference(readySet)
			ingredientDependencies[id] = diff
		}
	}

	return r.resolved, nil
}
