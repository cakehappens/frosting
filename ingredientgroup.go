package frosting

type IngredientGroup struct {
	header      string
	ingredients []*Ingredient
}

func (grp *IngredientGroup) Header() string {
	return grp.header
}

func (grp *IngredientGroup) Ingredients() []*Ingredient {
	return grp.ingredients
}

type IngredientGroupOption func(grp *IngredientGroup)

type IngredientGroupFn func() *IngredientGroup

func Includes(ingredients ...*Ingredient) IngredientGroupOption {
	return func(grp *IngredientGroup) {
		grp.ingredients = append(grp.ingredients, ingredients...)
	}
}

func MustNewGroup(header string, opts ...IngredientGroupOption) *IngredientGroup {
	grp := &IngredientGroup{
		header: header,
	}

	for _, o := range opts {
		o(grp)
	}

	return grp
}
