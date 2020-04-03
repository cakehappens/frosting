package frosting

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMustNewIngredientInfo(t *testing.T) {
	type args struct {
		name    string
		options []func(ing *IngredientInfo)
	}
	doesNotPanicTests := []struct {
		name string
		args args
		want *IngredientInfo
	}{
		{
			name: "name provided no options",
			args: args{
				name: "test",
			},
		},
		{
			name: "nil options",
			args: args{
				name:    "test",
				options: nil,
			},
		},
		{
			name: "empty options",
			args: args{
				name:    "test",
				options: []func(ing *IngredientInfo){},
			},
		},
		{
			name: "options do nothing",
			args: args{
				name: "test",
				options: []func(ing *IngredientInfo){
					func(ing *IngredientInfo) {

					},
				},
			},
		},
	}
	for _, tt := range doesNotPanicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(
				t,
				func() {
					if got := MustNewIngredientInfo(tt.args.name, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("MustNewIngredientInfo() = %v, want %v", got, tt.want)
					}
				},
			)
		})
	}

	panicTests := []struct {
		name string
		args args
		want *IngredientInfo
	}{
		{
			name: "name is empty",
			args: args{
				name: "",
			},
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(
				t,
				func() {
					MustNewIngredientInfo(tt.args.name, tt.args.options...)
				},
			)
		})
	}
}
