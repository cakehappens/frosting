package frosting

import (
	gofakeit "github.com/brianvoe/gofakeit/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	gofakeit.Seed(0)
}

func Test_defaultDependencyResolver_Dequeue(t *testing.T) {
	type testData struct {
		name      string
		arrange   func(r *defaultDependencyResolver)
		preassert func(t *testing.T, r *defaultDependencyResolver)
		assert    func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool)
	}

	tests := []testData{
		{
			name:      "when ready is empty, return nil, false",
			arrange:   func(r *defaultDependencyResolver) {},
			preassert: func(t *testing.T, r *defaultDependencyResolver) {},
			assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
				assert.Nil(t, ing, "there should be no ready ingredients")
				assert.False(t, ok, "no ingredient was returned")
			},
		},
		func() testData {
			const name = "when ready has 1 item, return item, true"

			expectedIngredient := &Ingredient{name: gofakeit.Word()}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[expectedIngredient.name] = expectedIngredient
					r.ready.Enqueue(expectedIngredient)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					if assert.NotNil(t, ing, "because there an ingredient ready") {
						assert.Equal(t, expectedIngredient, ing)
					}

					assert.True(t, ok, "because an ingredient was returned")
				},
			}
		}(),
		func() testData {
			const name = "when ready has multiple items, return item, true"

			expectedIngredient := &Ingredient{name: "a"}
			ing2 := &Ingredient{name: "b"}
			ing3 := &Ingredient{name: "c"}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[expectedIngredient.name] = expectedIngredient
					r.ingredients[ing2.name] = ing2
					r.ingredients[ing3.name] = ing3

					r.ready.Enqueue(expectedIngredient)
					r.ready.Enqueue(ing2)
					r.ready.Enqueue(ing3)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					if assert.NotNil(t, ing, "because there an ingredient ready") {
						assert.Equal(t, expectedIngredient, ing)
					}

					assert.True(t, ok, "because an ingredient was returned")
				},
			}
		}(),
		func() testData {
			const name = "when there are ingredients but none of them are ready, returns nil, false"

			ing1 := &Ingredient{name: "a"}
			ing2 := &Ingredient{name: "b"}
			ing3 := &Ingredient{name: "c"}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[ing1.name] = ing1
					r.ingredients[ing2.name] = ing2
					r.ingredients[ing3.name] = ing3
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					assert.Nil(t, ing, "because there are no ingredients ready")
					assert.False(t, ok, "because no ingredient was returned")
				},
			}
		}(),
		func() testData {
			const name = "dequeued item is removed from ready"

			ing1 := &Ingredient{name: "a"}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[ing1.name] = ing1
					r.ready.Enqueue(ing1)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {
					assert.Equal(t, 1, r.ready.Length(), "because item was enqueued")
				},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					assert.True(t, r.ready.IsEmpty(), "because item was dequeued")
				},
			}
		}(),
		func() testData {
			const name = "dequeued item is added to inflight"

			ing1 := &Ingredient{name: "a"}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[ing1.name] = ing1
					r.ready.Enqueue(ing1)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {
					assert.Equal(t, 0, r.inflight.Cardinality(), "because nothing has been dequeued")
				},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					if assert.Equal(t, 1, r.inflight.Cardinality(), "because item was dequeued") {
						assert.Contains(t, r.inflight.Ingredients(), ing1)
					}
				},
			}
		}(),
		func() testData {
			const name = "dequeued item is removed from unready"

			ing1 := &Ingredient{name: "a"}

			return testData{
				name: name,
				arrange: func(r *defaultDependencyResolver) {
					r.ingredients[ing1.name] = ing1
					r.unready.Add(ing1)
					r.ready.Enqueue(ing1)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {
					assert.Equal(t, 1, r.unready.Cardinality(), "because ingredient started as unready")
				},
				assert: func(t *testing.T, r *defaultDependencyResolver, ing *Ingredient, ok bool) {
					assert.Equal(t, 0, r.unready.Cardinality(), "because item transitioned to inflight")
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newDefaultDependencyResolver()
			tt.arrange(r)
			tt.preassert(t, r)
			got1, got2 := r.Dequeue()
			tt.assert(t, r, got1, got2)
		})
	}
}

func Test_defaultDependencyResolver_Length(t *testing.T) {
	type testData struct {
		name      string
		arrange   func(r *defaultDependencyResolver)
		preassert func(t *testing.T, r *defaultDependencyResolver)
		assert    func(t *testing.T, r *defaultDependencyResolver, got int)
	}

	tests := []testData{
		func() testData {
			return testData{
				name: "when everything is empty, returns 0",
				arrange: func(r *defaultDependencyResolver) {
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {
				},
				assert: func(t *testing.T, r *defaultDependencyResolver, got int) {
					assert.Equal(t, 0, got, "because everything is empty")
				},
			}
		}(),
		func() testData {
			return testData{
				name: "sums ready, unready, and inflight",
				arrange: func(r *defaultDependencyResolver) {
					a := &Ingredient{name: "a"}
					b := &Ingredient{name: "b"}
					c := &Ingredient{name: "c"}
					d := &Ingredient{name: "d"}
					e := &Ingredient{name: "e"}
					f := &Ingredient{name: "f"}
					g := &Ingredient{name: "g"}
					h := &Ingredient{name: "h"}
					i := &Ingredient{name: "i"}

					r.unready.Add(a)
					r.inflight.Add(b)
					r.inflight.Add(c)
					r.inflight.Add(d)
					r.ready.Enqueue(e)
					r.ready.Enqueue(f)
					r.ready.Enqueue(g)
					r.ready.Enqueue(h)
					r.ready.Enqueue(i)
				},
				preassert: func(t *testing.T, r *defaultDependencyResolver) {

				},
				assert: func(t *testing.T, r *defaultDependencyResolver, got int) {
					assert.Equal(t, 9, got, "because there is 1 unready + 3 inflight + 5 ready")
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newDefaultDependencyResolver()
			tt.arrange(r)
			tt.preassert(t, r)
			got := r.Length()
			tt.assert(t, r, got)
		})
	}
}

func Test_defaultDependencyResolver_NotifyComplete(t *testing.T) {
	type testData struct {
		name      string
		args      *Ingredient
		arrange   func(r *defaultDependencyResolver)
		preassert func(t *testing.T, r *defaultDependencyResolver)
		assert    func(t *testing.T, r *defaultDependencyResolver)
	}

	//assert := assert.New(t)

	tests := []testData{
		//func() testData {
		//	return testData{
		//		name: "",
		//	}
		//}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newDefaultDependencyResolver()
			tt.arrange(r)
			tt.preassert(t, r)
			r.NotifyComplete(tt.args)
			tt.assert(t, r)
		})
	}
}

func Test_defaultDependencyResolver_Load(t *testing.T) {
	type testData struct {
		name      string
		args      map[string]*Ingredient
		arrange   func(r *defaultDependencyResolver)
		preassert func(t *testing.T, r *defaultDependencyResolver)
		assert    func(t *testing.T, r *defaultDependencyResolver, got error)
	}

	tests := []testData{
		//func() testData {
		//	return testData{}
		//}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newDefaultDependencyResolver()
			tt.arrange(r)
			tt.preassert(t, r)
			err := r.Load(tt.args)
			tt.assert(t, r, err)
		})
	}
}

func Test_defaultDependencyResolver_AcceptanceTests(t *testing.T) {
	t.Run("single ingredient scenario", func(t *testing.T) {
		require := require.New(t)

		input := make(map[string]*Ingredient)
		ing := &Ingredient{name: "a"}
		input[ing.name] = ing

		r := newDefaultDependencyResolver()

		err := r.Load(input)
		require.NoError(err, "because there should be no circular dependency")
		require.Equal(1, r.Length(), "because one ingredient has been loaded")

		val, ok := r.Dequeue()

		require.Equal(ing, val, "because single ingredient dequeued")
		require.True(ok, "because ingredient was returned")

		require.Equal(1, r.Length(), "because one ingredient is in flight")

		r.NotifyComplete(val)

		require.Equal(0, r.Length(), "because one ingredient has completed")
	})

	t.Run("single ingredient depends on self", func(t *testing.T) {
		input := make(map[string]*Ingredient)
		ing := &Ingredient{name: "a"}
		ing.dependencies = append(ing.dependencies, ing)
		input[ing.name] = ing

		r := newDefaultDependencyResolver()

		err := r.Load(input)
		require.Error(t, err, "because there is a circular dependency")
	})

	t.Run("ingredient circle (circular dependency)", func(t *testing.T) {
		input := make(map[string]*Ingredient)

		ingA := &Ingredient{name: "a"}
		input[ingA.name] = ingA

		ingB := &Ingredient{name: "b"}
		input[ingB.name] = ingB

		ingC := &Ingredient{name: "c"}
		input[ingC.name] = ingC

		// A -> B
		ingA.dependencies = append(ingA.dependencies, ingB)

		// B -> C
		ingB.dependencies = append(ingB.dependencies, ingC)

		// C -> A
		ingC.dependencies = append(ingC.dependencies, ingA)

		r := newDefaultDependencyResolver()

		err := r.Load(input)
		require.Error(t, err, "because there is a circular dependency")
	})

	t.Run("ingredients tree with circular dependency", func(t *testing.T) {
		input := make(map[string]*Ingredient)

		ingA := &Ingredient{name: "a"}
		input[ingA.name] = ingA

		ingB := &Ingredient{name: "b"}
		input[ingB.name] = ingB

		ingC := &Ingredient{name: "c"}
		input[ingC.name] = ingC

		// A -> B
		ingA.dependencies = append(ingA.dependencies, ingB)

		// A -> C
		ingA.dependencies = append(ingA.dependencies, ingC)

		// C -> A
		ingC.dependencies = append(ingC.dependencies, ingA)

		r := newDefaultDependencyResolver()

		err := r.Load(input)
		require.Error(t, err, "because there is a circular dependency")
	})

	t.Run("ingredients tree, can dequeue multiple ready", func(t *testing.T) {
		require := require.New(t)

		input := make(map[string]*Ingredient)

		ingA := &Ingredient{name: "a"}
		input[ingA.name] = ingA

		ingB := &Ingredient{name: "b"}
		input[ingB.name] = ingB

		ingC := &Ingredient{name: "c"}
		input[ingC.name] = ingC

		// A -> B
		// A -> C
		ingA.dependencies = append(ingA.dependencies, []*Ingredient{ingB, ingC}...)

		r := newDefaultDependencyResolver()

		err := r.Load(input)
		total := len(input)
		require.Equal(3, total, "because we have 3 ingredients")
		require.Equal(total, r.Length(), "because we have 3 ingredients")
		numCompleted := 0

		require.NoError(err, "because there should be no circular dependency")
		require.Equal(total, r.Length(), "because %d ingredients have been loaded", total)

		require.Equal(2, r.Ready(), "because b & c are without dependencies")

		val1, ok := r.Dequeue()
		require.True(ok, "because the dequeue should be successful")
		require.NotEqual(ingA, val1, "because the first dequeue should be either b or c")

		val2, ok := r.Dequeue()
		require.True(ok, "because the dequeue should be successful")
		require.NotEqual(ingA, val1, "because the second dequeue should be either b or c")

		require.NotEqual(val1, val2, "because only b and c should have been dequeued")

		require.Equal(total, r.Length(), "because nothing has been completed", total)

		err = r.NotifyComplete(val1)
		require.NoError(err, "because the ingredient should be in-flight")
		numCompleted += 1

		require.Equal(total-numCompleted, r.Length(), "because %d ingredients have been completed", numCompleted)

		_, ok = r.Dequeue()
		require.False(ok, "because 'a' still has an incomplete dependency")

		err = r.NotifyComplete(val2)
		require.NoError(err, "because the ingredient should be in-flight")
		numCompleted += 1

		require.Equal(total-numCompleted, r.Length(), "because %d ingredients have been completed", numCompleted)

		val3, ok := r.Dequeue()
		require.True(ok, "because the dequeue should be successful")
		require.Equal(ingA, val3, "because 'a' should have been returned")

		err = r.NotifyComplete(val3)
		require.NoError(err, "because the ingredient should be in-flight")
		numCompleted += 1

		require.Equal(total-numCompleted, r.Length(), "because %d ingredients have been completed", numCompleted)

		require.Equal(0, r.Length(), "because all ingredients have completed")
	})
}
