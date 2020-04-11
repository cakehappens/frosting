package frosting

import "github.com/spf13/cobra"

type defaultCommandBuilder struct{}

func newSimpleCommandBuilder() *defaultCommandBuilder {
	return &defaultCommandBuilder{}
}

type CommandBuilder interface {
	Build(ingredient *Ingredient) (*cobra.Command, error)
}

func (c *defaultCommandBuilder) Build(ingredient *Ingredient) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     ingredient.name,
		Short:   ingredient.short,
		Long:    ingredient.long,
		Example: ingredient.example,
		Aliases: ingredient.aliases,
	}

	if ingredient.flagsFn != nil {
		ingredient.flagsFn(cmd.Flags())
	}

	return cmd, nil
}
