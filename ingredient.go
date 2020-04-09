package frosting

import (
	"context"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type RunFn func(ctx context.Context) error

type IngredientOption func(ing *Ingredient)

type IngredientFn func() *Ingredient

type Ingredient struct {
	name               string
	runFn              RunFn
	aliases            []string
	short              string
	long               string
	example            string
	dependencies       []string
	serialDependencies []string
	flagsFn            func(flagSet *flag.FlagSet)
	command            *cobra.Command
	ready              bool
	ran                bool
}

func (ing *Ingredient) Name() string {
	return ing.name
}

func (ing *Ingredient) Dependencies() []string {
	return ing.dependencies
}

func (ing *Ingredient) SerialDependencies() []string {
	return ing.serialDependencies
}

func WithDependencies(deps ...string) IngredientOption {
	return func(ing *Ingredient) {
		ing.dependencies = append(ing.dependencies, deps...)
	}
}

func WithSerialDependencies(deps ...string) IngredientOption {
	return func(ing *Ingredient) {
		ing.serialDependencies = append(ing.serialDependencies, deps...)
	}
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
