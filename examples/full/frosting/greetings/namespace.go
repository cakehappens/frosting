package greetings

import (
	fr "github.com/cakehappens/frosting"
)

var GreetingsNamespace = fr.MustNewNamespace(
	"greetings",
	[]fr.IngredientFn{
		Hello,
		Goodbye,
	},
)
