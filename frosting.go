package frosting

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
	"os"
	"os/signal"
	//"reflect"
	//"runtime"
	//"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/oklog/run"
)

type ClientInfo struct {
	rootCommand                *cobra.Command
	ingredientGroups           []*IngredientGroup
	ingredients                map[string]*Ingredient
	commands                   map[string]*cobra.Command
	commandsPerIngredientGroup map[*IngredientGroup][]*cobra.Command
}

func newRootCommand(binaryName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   binaryName,
		Short: "Frosting!",
		Long:  "CakeHappens/Frosting!",
		Run:   runHelp,
	}

	return cmd
}

func New(binaryName string) *ClientInfo {
	return &ClientInfo{
		rootCommand:                newRootCommand(binaryName),
		ingredients:                make(map[string]*Ingredient),
		commands:                   make(map[string]*cobra.Command),
		commandsPerIngredientGroup: make(map[*IngredientGroup][]*cobra.Command),
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func (c *ClientInfo) MustAddIngredientGroups(ingGrps ...*IngredientGroup) {
	for _, ingGrp := range ingGrps {
		if ingGrp == nil {
			panic(errors.New("cannot add nil ingredientGroup"))
		}

		c.ingredientGroups = append(c.ingredientGroups, ingGrps...)

		for _, ing := range ingGrp.ingredients {
			c.ingredients[ing.name] = ing
		}
	}
}

func (c *ClientInfo) Ingredients() map[string]*Ingredient {
	return c.ingredients
}

func (c *ClientInfo) Commands() map[string]*cobra.Command {
	if len(c.commands) == 0 {
		for _, ingGroup := range c.ingredientGroups {
			for _, ing := range ingGroup.ingredients {
				cmd := createCommandFromIngredient(ing)
				c.commandsPerIngredientGroup[ingGroup] = append(c.commandsPerIngredientGroup[ingGroup], cmd)
				c.commands[ing.name] = cmd
			}
		}
	}

	return c.commands
}

func (c *ClientInfo) CommandsPerIngredientGroup() map[*IngredientGroup][]*cobra.Command {
	return c.commandsPerIngredientGroup
}

func (c *ClientInfo) createCommandFromIngredient(ing *Ingredient) *cobra.Command {
	return &cobra.Command{
		Use:     ing.name,
		Aliases: ing.aliases,
		Short:   ing.short,
		Long:    ing.long,
		Example: ing.example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

func (c *ClientInfo) Execute(args ...string) {
	ctx, cancel := context.WithCancel(context.Background())

	runGroup := run.Group{}
	{
		cancelInterrupt := make(chan struct{})
		runGroup.Add(
			createSignalWatcher(ctx, cancelInterrupt, cancel),
			func(error) {
				close(cancelInterrupt)
			})
	}
	{
		runGroup.Add(func() error {
			rootC := c.rootCommand
			rootC.SetArgs(os.Args[1:])

			cmdGroups := templates.CommandGroups{}

			for ingGroup, cmds := range c.commandsPerIngredientGroup {
				rootC.AddCommand(cmds...)
				cmdGroups = append(cmdGroups, templates.CommandGroup{
					Message:  ingGroup.header,
					Commands: cmds,
				})
			}

			// add all commands from the cmdGroups as subcommands to the root
			cmdGroups.Add(rootC)
			templates.ActsAsRootCommand(rootC, nil, cmdGroups...)

			return rootC.ExecuteContext(ctx)
		}, func(error) {
			cancel()
		})
	}

	err := runGroup.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit reason: %s\n", err)
		os.Exit(1)
	}

	color.New(color.FgGreen).Fprintln(os.Stderr, "Done!")
}

// This function just sits and waits for ctrl-C
func createSignalWatcher(ctx context.Context, cancelInterruptChan <-chan struct{}, cancel context.CancelFunc) func() error {
	return func() error {
		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			err := errors.New(fmt.Sprintf("received signal %s", sig))
			fmt.Fprintf(os.Stderr, "%s\n", err)
			signal.Stop(c)
			cancel()
			return err
		case <-ctx.Done():
			return nil
		case <-cancelInterruptChan:
			return nil
		}
	}
}
