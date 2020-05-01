package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cakehappens/frosting"
)

func NewBuildIngredient(name string, options ...frosting.IngredientOption) *frosting.Ingredient {
	return frosting.MustNewIngredient(
		name,
		func(ctx context.Context, ing *frosting.Ingredient) error {
			fmt.Println("Building...")
			return nil
		},
		append([]frosting.IngredientOption{
			frosting.WithHelpDescriptions("buildShort", "buildLong"),
		}, options...)...,
	)
}

func NewTestIngredient(name string, options ...frosting.IngredientOption) *frosting.Ingredient {
	return frosting.MustNewIngredient(
		name,
		func(ctx context.Context, ing *frosting.Ingredient) error {
			fmt.Println("Testing...")
			return nil
		},
		options...,
	)
}

func main() {
	f := frosting.New("frost")

	{
		build := NewBuildIngredient("build")

		// test depends on build
		test := NewTestIngredient(
			"test",
			frosting.WithDependencies(build),
		)

		f.Group(
			"Basic Commands (Beginner):",
			build,
			test,
		)
	}

	f.Execute(os.Args[1:]...)
}
