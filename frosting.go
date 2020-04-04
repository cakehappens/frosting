package frosting

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	//"reflect"
	//"runtime"
	//"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/oklog/run"
	//"gopkg.in/eapache/queue.v1"
	"k8s.io/kubectl/pkg/util/templates"
)

type ClientInfo struct {
	defaultIngredient *Ingredient
	ingredientGroups  []*IngredientGroup
}

func newRootCommand(binaryName string) *cobra.Command {
	cmd := &cobra.Command{
		Use: binaryName,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				runHelp(cmd, args)
				return nil
			}

			fmt.Printf("running: %s\n", args)
			return nil
		},
	}

	groups := templates.CommandGroups{}

	cmd.AddCommand()

	return cmd
}

func New(binaryName string) *ClientInfo {
	return &ClientInfo{
		defaultIngredient: &Ingredient{
			cobraCommand: newRootCommand(binaryName),
		},
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func (c *ClientInfo) MustAddIngredientGroups(ingGrps ...*IngredientGroup) {
	ingGrpNsMap := make(map[string]bool)
	ingGrpHeaderMap := make(map[string]bool)

	for _, ingGrp := range ingGrps {
		if ingGrp.Header == "" {
			panic(fmt.Errorf("ingredient groups must have a non-empty header. found empty header for: %+v", ingGrp))
		}

		if ingGrpNsMap[ingGrp.Namespace] {
			panic(fmt.Errorf("duplicate namespace: %s", ingGrp.Namespace))
		}

		if ingGrpHeaderMap[ingGrp.Header] {
			panic(fmt.Errorf("duplicate header: %s", ingGrp.Header))
		}

		ingGrpNsMap[ingGrp.Namespace] = true
		ingGrpHeaderMap[ingGrp.Header] = true
		c.ingredientGroups = append(c.ingredientGroups, ingGrp)
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
			rootC := c.defaultIngredient.cobraCommand
			rootC.SetArgs(os.Args[1:])
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
