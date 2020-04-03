package frosting

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/eapache/queue.v1"
)

type RunFn func(ctx context.Context) error

type Ingredient struct {
	Name         string
	RunFn        RunFn
	Dependencies func() []*Ingredient

	cobraCommand *cobra.Command
}

func (ing *Ingredient) MustSetDependencies(deps ...*Ingredient) {
	for _, dep := range deps {
		if dep == ing {
			panic(errors.New("ingredient cannot be a dependency of itself"))
		}
	}

	ing.Dependencies = func() []*Ingredient {
		return deps
	}
}

// this probably only is needed from root ingredients
func (ing *Ingredient) resolveAllDeps() ([]*Ingredient, error) {
	var deps []*Ingredient
	visited := make(map[string]bool)

	ingQueue := queue.New()

	ingQueue.Add(ing)

	for ingQueue.Length() > 0 {
		dep := ingQueue.Remove().(*Ingredient)

		if visited[dep.Name] {
			return nil, fmt.Errorf("found circular reference via: %s", dep.Name)
		}

		deps = append(deps, dep)
		visited[dep.Name] = true

		for _, d := range dep.Dependencies() {
			ingQueue.Add(d)
		}
	}

	return deps, nil
}
