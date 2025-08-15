package sozzler_test

import (
	"fmt"
	"mp/sozzler/pkg/sozzler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredientPredicate(t *testing.T) {
	testCases := []struct {
		recipe sozzler.Recipe
		term   string
		wantI  string
		wantO  bool
	}{
		{
			recipe: sozzler.Recipe{},
			term:   "foo",
			wantI:  "",
			wantO:  false,
		},
		{
			recipe: sozzler.Recipe{
				Components: []sozzler.Component{
					{Ingredient: "zippy foo bar"},
				},
			},
			term:  "foo",
			wantI: "zippy foo bar",
			wantO: true,
		},
		{
			recipe: sozzler.Recipe{
				Components: []sozzler.Component{
					{Ingredient: "zippy foo bar"},
				},
			},
			term:  "banana",
			wantI: "",
			wantO: false,
		},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			gotI, gotO := sozzler.NewIngredientPredicate(tC.term).Match(&tC.recipe)
			assert.Equal(t, tC.wantI, gotI)
			assert.Equal(t, tC.wantO, gotO)
		})
	}
}

func TestNamePredicate(t *testing.T) {
	testCases := []struct {
		recipe sozzler.Recipe
		term   string
		wantI  string
		wantO  bool
	}{
		{
			recipe: sozzler.Recipe{},
			term:   "foo",
			wantI:  "",
			wantO:  false,
		},
		{
			recipe: sozzler.Recipe{Name: "ZippY FoO BaR"},
			term:   "foo",
			wantI:  "ZippY FoO BaR",
			wantO:  true,
		},
		{
			recipe: sozzler.Recipe{Name: "pineapple ice cream"},
			term:   "banana",
			wantI:  "",
			wantO:  false,
		},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			gotI, gotO := sozzler.NewNamePredicate(tC.term).Match(&tC.recipe)
			assert.Equal(t, tC.wantI, gotI)
			assert.Equal(t, tC.wantO, gotO)
		})
	}
}

func TestRatingPredicate(t *testing.T) {
	testCases := []struct {
		recipe sozzler.Recipe
		term   int
		wantI  string
		wantO  bool
	}{
		{
			recipe: sozzler.Recipe{},
			term:   2,
			wantI:  "",
			wantO:  false,
		},
		{
			recipe: sozzler.Recipe{Name: "foo", Rating: 3},
			term:   2,
			wantI:  "🫒🫒🫒",
			wantO:  true,
		},
		{
			recipe: sozzler.Recipe{Name: "bar", Rating: 2},
			term:   2,
			wantI:  "🫒🫒",
			wantO:  true,
		},

		{
			recipe: sozzler.Recipe{Rating: 1},
			term:   2,
			wantI:  "",
			wantO:  false,
		},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			gotI, gotO := sozzler.NewRatingPredicate(tC.term).Match(&tC.recipe)
			assert.Equal(t, tC.wantI, gotI)
			assert.Equal(t, tC.wantO, gotO)
		})
	}
}
