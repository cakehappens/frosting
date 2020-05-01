package frosting

import (
	"errors"
	"fmt"
)

type DependencyResolver interface {
	Ready() int
	Dequeue() (*Ingredient, bool)
	Length() int
	NotifyComplete(ingredient *Ingredient)
}

type defaultDependencyResolver struct {
	ingredients map[string]*Ingredient
	// values are those dependents of the key
	dependentsGraph map[*Ingredient]*IngredientSet
	// values are those parents of the key
	parentsGraph map[*Ingredient]*IngredientSet
	unready      *IngredientSet
	inflight     *IngredientSet
	ready        *IngredientQueue
}

func newDefaultDependencyResolver() *defaultDependencyResolver {
	return &defaultDependencyResolver{
		ingredients:     make(map[string]*Ingredient),
		dependentsGraph: make(map[*Ingredient]*IngredientSet),
		parentsGraph:    make(map[*Ingredient]*IngredientSet),
		ready:           NewIngredientQueue(),
		unready:         NewIngredientSet(),
		inflight:        NewIngredientSet(),
	}
}

func (r *defaultDependencyResolver) Load(ingredients map[string]*Ingredient) error {
	// we'll use a temporary map to check for circular dependencies
	tmpIngredientDependencies := make(map[*Ingredient]*IngredientSet)
	for _, node := range ingredients {
		r.ingredients[node.name] = node

		dependencySet := NewIngredientSet()
		for _, dep := range node.dependencies {
			dependencySet.Add(dep)
			if parentSet, ok := r.parentsGraph[dep]; ok {
				parentSet.Add(node)
			} else {
				r.parentsGraph[dep] = NewIngredientSet(node)
			}
		}
		r.dependentsGraph[node] = dependencySet
		tmpIngredientDependencies[node] = dependencySet.Clone()

		if len(node.dependencies) == 0 {
			r.ready.Enqueue(node)
		} else {
			r.unready.Add(node)
		}
	}

	// check for circular dependency
	for len(tmpIngredientDependencies) != 0 {
		// readyset type is *Ingredient
		readySet := NewIngredientSet()
		for ing, deps := range tmpIngredientDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(ingredients[ing.name])
			}
		}

		if readySet.Cardinality() == 0 {
			return errors.New("circular dependency found in ingredients")
		}

		for _, ing := range readySet.Ingredients() {
			delete(tmpIngredientDependencies, ing)
		}

		for parent, deps := range tmpIngredientDependencies {
			diff := deps.Difference(readySet)
			tmpIngredientDependencies[parent] = diff
		}
	}

	return nil
}

func (r *defaultDependencyResolver) Ready() int {
	return r.ready.Length()
}

func (r *defaultDependencyResolver) Dequeue() (*Ingredient, bool) {
	// remove item from ready queue
	val, ok := r.ready.Dequeue()
	if !ok {
		return nil, false
	}
	// remove item from unready
	r.unready.Remove(val)

	// add item to inflight
	r.inflight.Add(val)

	return val, true
}

// Length returns number of items that are left to be completed
// Items can be in a state of ready, unready or inflight
// Once marked as complete via NotifyComplete
// they are no longer tracked
func (r *defaultDependencyResolver) Length() int {
	return r.ready.Length() + r.unready.Cardinality() + r.inflight.Cardinality()
}

// whenever a node is marked complete, we can look at that nodes that depend on it (parent)
// and check if there are any other dependencies for not yet satisfied
// if there are no dependencies for that parent, add that node to the ready set
func (r *defaultDependencyResolver) NotifyComplete(ingredient *Ingredient) error {
	// remove item from inflight
	if !r.inflight.Contains(ingredient) {
		return fmt.Errorf("tried to notify completion of an ingredient not in-flight: %s", ingredient.name)
	}
	r.inflight.Remove(ingredient)

	if parents, ok := r.parentsGraph[ingredient]; ok {
		for _, p := range parents.Ingredients() {
			if !r.dependentsGraph[p].Contains(ingredient) {
				return fmt.Errorf("something went horribly wrong. Dependent or Parent Graph is malformed. Tried to remove ingredient from parent dependencies")
			}
			r.dependentsGraph[p].Remove(ingredient)

			if r.dependentsGraph[p].Cardinality() == 0 {
				r.ready.Enqueue(p)
				r.unready.Remove(p)
			}
		}
	}

	return nil
}
