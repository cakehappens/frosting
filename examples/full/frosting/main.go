package main

import (
	"fmt"
	fr "github.com/cakehappens/frosting"
	"github.com/cakehappens/frosting/examples/full/frosting/greetings"
	"os"
)

func Build() *fr.IngredientInfo {
	return fr.MustNewIngredientInfo(
		"build",
		func(ing *fr.IngredientInfo) {
			ing.Fn = func() error {
				fmt.Println("Building...")
				return nil
			}
		},
	)
}

func main() {
	client := fr.MustNew(
		greetings.GreetingsNamespace,
	)
	client.RootNamespace().MustAddIngredients(Build)
	client.Execute(os.Args[1:]...)
}
