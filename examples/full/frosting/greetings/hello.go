package greetings

import (
	"fmt"
	fr "github.com/cakehappens/frosting"
)

func Hello() *fr.IngredientInfo {
	return fr.MustNewIngredientInfo(
		"hello",
		func(ing *fr.IngredientInfo) {
			ing.Fn = func() error {
				fmt.Println("Hello World...")
				return nil
			}
		},
	)
}

func Goodbye() *fr.IngredientInfo {
	return fr.MustNewIngredientInfo(
		"goodbye",
		func(ing *fr.IngredientInfo) {
			ing.AddDependencies(
				Hello,
			)

			ing.Fn = func() error {
				fmt.Println("Goodbye World...")
				return nil
			}
		},
	)
}
