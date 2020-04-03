package greetings

import (
	"context"
	"fmt"
	"github.com/cakehappens/frosting"
)

const (
	Hello   = "hello"
	GoodBye = "goodbye"
)

func NewHelloIngredient() *frosting.Ingredient {
	return &frosting.Ingredient{
		Name: Hello,
		RunFn: func(ctx context.Context) error {
			fmt.Println("Hello World...")
			return nil
		},
	}
}

func NewGoodByeIngredient() *frosting.Ingredient {
	return &frosting.Ingredient{
		Name: GoodBye,
		RunFn: func(ctx context.Context) error {
			fmt.Println("GoodBye World...")
			return nil
		},
	}
}
