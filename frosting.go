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
	rootCommand      *cobra.Command
	ingredientGroups []*ingredientGroup
	ingredients      map[string]*Ingredient
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
		rootCommand: newRootCommand(binaryName),
		ingredients: make(map[string]*Ingredient),
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func (c *ClientInfo) Group(header string, ingredients ...*Ingredient) {
	grp := &ingredientGroup{
		header:      header,
		ingredients: ingredients,
	}

	c.ingredientGroups = append(c.ingredientGroups, grp)

	for _, ing := range ingredients {
		if _, ok := c.ingredients[ing.name]; ok {
			panic(fmt.Errorf("ingredient with name already added: %s", ing.name))
		}

		c.ingredients[ing.name] = ing
	}
}

//func mergeIngredientMaps(ingredientMaps ...map[string]*Ingredient) (map[string]*Ingredient, error) {
//	newMap := make(map[string]*Ingredient)
//
//	for _, ingMap := range ingredientMaps {
//		for name, ing := range ingMap {
//			if _, ok := newMap[name]; ok {
//				return nil, fmt.Errorf("ingredient with name already added: %s", name)
//			}
//
//			newMap[name] = ing
//		}
//	}
//
//	return newMap, nil
//}

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

			//for ingGroup, cmds := range c.commandsPerIngredientGroup {
			//	rootC.AddCommand(cmds...)
			//	cmdGroups = append(cmdGroups, templates.CommandGroup{
			//		Message:  ingGroup.header,
			//		Commands: cmds,
			//	})
			//}

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
