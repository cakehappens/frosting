package main

import (
	"context"
	"fmt"
	"github.com/cakehappens/frosting"
)

func NewBuildIngredient() *frosting.Ingredient {
	ing := &frosting.Ingredient{
		Name: "build",
		RunFn: func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
		},
	}

	ing.MustSetDependencies(
		NewTestIngredient(),
	)

	return ing
}

func NewTestIngredient() *frosting.Ingredient {
	return &frosting.Ingredient{
		Name: "test",
		RunFn: func(ctx context.Context) error {
			fmt.Println("Testing...")
			return nil
		},
	}
}

func main() {
	f := frosting.New("fr")
	f.MustAddIngredientGroups(
		&frosting.IngredientGroup{
			Header:    "stuff",
			Namespace: "",
			Ingredients: []*frosting.Ingredient{
				NewBuildIngredient(),
			},
		},
	)

	f.Execute("foo")
}
