package frosting

import (
	"context"

	flag "github.com/spf13/pflag"
)

type RunFn func(ctx context.Context, ing *Ingredient) error

type IngredientOption func(ing *Ingredient)

type IngredientFn func() *Ingredient

type Ingredient struct {
	name         string
	runFn        RunFn
	aliases      []string
	short        string
	long         string
	example      string
	dependencies []*Ingredient
	flagsFn      func(flagSet *flag.FlagSet)
	ran          bool
}

func (ing *Ingredient) Name() string {
	return ing.name
}

func WithAliases(aliases ...string) IngredientOption {
	return func(ing *Ingredient) {
		ing.aliases = append(ing.aliases, aliases...)
	}
}

func WithHelpDescriptions(short, long string) IngredientOption {
	return func(ing *Ingredient) {
		ing.short = short
		ing.long = long
	}
}

func WithExampleHelp(example string) IngredientOption {
	return func(ing *Ingredient) {
		ing.example = example
	}
}

func WithDependencies(deps ...*Ingredient) IngredientOption {
	return func(ing *Ingredient) {
		ing.dependencies = deps
	}
}

func WithFlags(fn func(set *flag.FlagSet)) IngredientOption {
	return func(ing *Ingredient) {
		ing.flagsFn = fn
	}
}

func MustNewIngredient(name string, runFn RunFn, opts ...IngredientOption) *Ingredient {
	ing := &Ingredient{
		name:  name,
		runFn: runFn,
	}

	for _, o := range opts {
		o(ing)
	}

	return ing
}

type ingredientGroup struct {
	header      string
	ingredients []*Ingredient
}
