package frosting

//import (
//	"errors"
//	"fmt"
//
//	mapset "github.com/deckarep/golang-set"
//)
//
//type IngredientGraph []Ingredient
//
//type IngredientDependencyResolver interface {
//	Resolve() (IngredientGraph, error)
//	CurrentReady() IngredientGraph
//}
//
//type ingredientDependencyResolver struct {
//	ingredients map[string]*IngredientInfo
//	resolved    IngredientGraph
//}
//
//func newIngredientDependencyResolver(ingredients IngredientGraph) (*ingredientDependencyResolver, error) {
//	r := &ingredientDependencyResolver{}
//
//	r.ingredients = make(map[string]*IngredientInfo)
//	r.resolved = IngredientGraph{}
//
//	for _, ing := range ingredients {
//		if ing.id == "" {
//			return nil, fmt.Errorf("ingredientInfo does not have id. Create only using NewIngredientInfo. name: %s", ing.name)
//		}
//		r.ingredients[ing.id] = ing
//	}
//
//	return r, nil
//}
//
//func (r *ingredientDependencyResolver) Resolve() (IngredientGraph, error) {
//	// if we've already resolved, don't re-resolve
//	if len(r.resolved) > 0 {
//		return r.resolved, nil
//	}
//
//	ingredientDependencies := make(map[string]mapset.Set)
//
//	// Populate the map
//	for _, ing := range r.ingredients {
//		ingredientDependencies[ing.id] = ing.dependencies.Clone()
//	}
//
//	// Iteratively find and remove nodes from the graph which have no dependencies.
//	// If at some point there are still nodes in the graph and we cannot find
//	// nodes without dependencies, that means we have a circular dependency
//	for len(ingredientDependencies) != 0 {
//		// Get all nodes from the graph which have no dependencies
//		readySet := mapset.NewSet()
//		for id, deps := range ingredientDependencies {
//			if deps.Cardinality() == 0 {
//				readySet.Add(id)
//			}
//		}
//
//		// If there aren't any ready nodes, then we have a circular dependency
//		if readySet.Cardinality() == 0 {
//			var g IngredientGraph
//			for id := range ingredientDependencies {
//				g = append(g, r.ingredients[id])
//			}
//
//			return g, errors.New("circular dependency found")
//		}
//
//		// Remove the ready nodes and add them to the resolved graph
//		for id := range readySet.Iter() {
//			delete(ingredientDependencies, id.(string))
//			r.resolved = append(r.resolved, r.ingredients[id.(string)])
//		}
//
//		// Also make sure to remove the ready nodes from the
//		// remaining node dependencies as well
//		for id, deps := range ingredientDependencies {
//			diff := deps.Difference(readySet)
//			ingredientDependencies[id] = diff
//		}
//	}
//
//	return r.resolved, nil
//}
//
//func (r *ingredientDependencyResolver) CurrentReady() IngredientGraph {
//	currentReady := IngredientGraph{}
//	for _, ing := range r.resolved {
//		if ing.currentDependencies.Cardinality() == 0 {
//			currentReady = append(currentReady, ing)
//		}
//	}
//
//	return currentReady
//}
//
//var _ IngredientDependencyResolver = &ingredientDependencyResolver{}
