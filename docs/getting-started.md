# Architecture

## Quick Start

`frosting` is made up of `Ingredient`s (pun intended). An ingredient is analogous to a `target` in a make file.

### Library Code Example

📌 _Let the caller define the name of the target_

🏷️ _Ingredients are ** **globally unique by name** **, to allow different ingredients to have the same dependencies, but ensuring we don't call the dependent ingredient more than once._

```go
func NewBuildIngredient(name string, options ...frosting.IngredientOption) *frosting.Ingredient {
	return frosting.MustNewIngredient(
		name,
		func(ctx context.Context, ing *frosting.Ingredient) error {
			fmt.Println("Building...")
			return nil
		},
		append([]frosting.IngredientOption{
			frosting.WithHelpDescriptions("buildShort", "buildLong"),
		}, options...)...,
	)
}
```

Next, mix ingredients into frosting by adding them to a group. This provides a nice interface for getting help about related ingredients.
Ingredients are still ** **globally unique by name** ** regardless of group.

```
$ frost help
Add a little frosting to your operational tasks..

Basic Commands (Beginner):
  build          Build Things
  test           Test Things
```

```go
func main() {
	f := frosting.New("frost")

	{
		build := NewBuildIngredient("build")

		// test depends on build
		test := NewTestIngredient(
			"test",
			frosting.WithDependencies(build),
		)

		f.Group(
			"Basic Commands (Beginner):",
			build,
			test,
		)
	}

	f.Execute(os.Args[1:]...)
}
```

👀 For more information, check out the [examples](../examples)!
