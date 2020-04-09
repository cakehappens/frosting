package main

import (
	"context"
	"fmt"
	"github.com/cakehappens/frosting"
)

const Build = "build"

func NewBuildIngredient() *frosting.Ingredient {
	return frosting.MustNew(
		Build,
		func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
		},
		//ingredient.WithDependencies(Test),
		frosting.WithHelpDescriptions("buildShort", "buildLong"),
	)
}

const Test = "test"

func NewTestIngredient() *frosting.Ingredient {
	return frosting.MustNew(
		Test,
		func(ctx context.Context) error {
			fmt.Println("Testing...")
			return nil
		},
	)
}

func main() {
	f := frosting.New("frost")
	f.MustAddIngredientGroups(
		frosting.MustNewGroup(
			"Main Stuff:",
			frosting.Includes(
				NewBuildIngredient(),
				NewTestIngredient(),
			),
		),
	)

	f.Execute("foo")
}
