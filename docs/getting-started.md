# Architecture

## Quick Start

`frosting` is made up of `Ingredient`s (pun intended). An ingredient is analogous to a `target` in a make file.

### Library Code Example

📌 _Let the consumer define the name of the target_

🏷️ _Ingredients are ** **globally unique by name** **, to allow different ingredients to have the same dependencies, but ensuring we don't call the dependent ingredient more than once._

🛑 _Don't perform any initialization logic within the function that returns the ingredient. This function may be called multiple times, and the resulting ingredient struct may be tossed out based on the globally unique rule above_

```go
func NewBuildIngredient(name) *ingredient.Ingredient {
	return ingredient.MustNew(
		name,
		func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
		},
	)
}
```

Ingredients may be grouped as well, but regardless of group, they are still ** **globally unique by name** **. A group is only to assist when getting help.

```
$ frost help
Add a little frosting to your operational tasks..

Basic Commands (Beginner):
  build          Build Things
  test           Test Things
```

```go
func NewFooIngredientGroup() *frosting.IngredientGroup {
    return ingredient.MustNewGroup(
        "Basic Commands (Beginner):",
        ingredient.Includes(
            NewBuildIngredient("build"),
            NewTestIngredient("test"),
        ),
    ),
}
```

### Declaring Dependencies

The recommended way to define dependencies is to first declare the name of ingredient as a constant. Then use that constant when setting dependencies.

```go

const Build = "build"

func NewBuildIngredient() *ingredient.Ingredient {
	return ingredient.MustNew(
		Build,
		func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
        },
        ingredient.WithDependency(
            
        )
	)
}

const Test = "test"

func NewBuildIngredient(name) *ingredient.Ingredient {
	return ingredient.MustNew(
		name,
		func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
        },
        ingredient.WithDependencies(
            Build
        )
	)
}
```

You may also declare dependencies that should be run serially. These dependencies will implicitly be made dependencies of each other. The following will result in ingredients running in this order:

```
A -> B -> C -> This
```

```go
ingredient.WithSerialDependencies(
    A,
    B,
    C
)
```
