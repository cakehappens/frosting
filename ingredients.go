package frosting

import (
	"errors"
	"github.com/deckarep/golang-set"
)

type IngredientGraph []*IngredientInfo

type IngredientFn func() *IngredientInfo

type IngredientInfo struct {
	Name                string
	Fn                  func() error

	id                  string
	dependencyFns		[]IngredientFn
	// set of IngredientInfo's
	dependencies        mapset.Set
	currentDependencies mapset.Set
}

func (ing *IngredientInfo) AddDependencies(fns ...IngredientFn) {
	for _, fn := range fns {
		ing.dependencyFns = append(ing.dependencyFns, fn)
	}
}

func MustNewIngredientInfo(name string, options ...func(ing *IngredientInfo)) *IngredientInfo {
	if name == "" {
		panic(errors.New("name cannot be empty"))
	}

	ing := &IngredientInfo{
		Name: name,
		id:   newULID(),
		dependencies: mapset.NewSet(),
	}

	for _, fn := range options {
		fn(ing)
	}

	if ing.Fn == nil {
		ing.Fn = func() error {
			return nil
		}
	}

	ing.currentDependencies = ing.dependencies.Clone()

	return ing
}

type Ingredient interface {
	Satisfied() bool
}

func (ing *IngredientInfo) Satisfied() bool {
	if ing.currentDependencies.Cardinality() == 0 {
		return true
	}

	return false
}

var _ Ingredient = &IngredientInfo{}
