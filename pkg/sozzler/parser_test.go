package sozzler_test

import (
	"fmt"
	"mp/sozzler/pkg/sozzler"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestComponents(t *testing.T) {
	testCases := []struct {
		given         string
		wantComponent *sozzler.Component
		wantErr       error
	}{
		// happy paths
		{
			given:         "1/8 egg",
			wantComponent: component("egg", "1/8", ""),
		},
		{
			given:         "1/8 hard boiled egg",
			wantComponent: component("hard boiled egg", "1/8", ""),
		},
		{
			given:         "1/8 ounce hard boiled egg",
			wantComponent: component("hard boiled egg", "1/8", "ounce"),
		},
		{
			given:         "1/8 oz hard boiled egg",
			wantComponent: component("hard boiled egg", "1/8", "oz"),
		},
		{
			given:         "1 oz hard boiled egg",
			wantComponent: component("hard boiled egg", "1", "oz"),
		},
		{
			given:         "oz hard boiled egg",
			wantComponent: component("hard boiled egg", "", "oz"),
		},
		{
			given:         "oz egg",
			wantComponent: component("egg", "", "oz"),
		},
		{
			given:         "1oz egg",
			wantComponent: component("egg", "1", "oz"),
		},
		{
			// weird, but correct
			given:         "1/2x egg",
			wantComponent: component("x egg", "1/2", ""),
		},
		{
			// weird, but correct
			given:         "2x egg",
			wantComponent: component("x egg", "2", ""),
		},
		{
			// weird, but correct
			given:         "1x",
			wantComponent: component("x", "1", ""),
		},
		{
			// weird, but correct
			given:         "1egg",
			wantComponent: component("egg", "1", ""),
		},
		{
			// weird, but correct
			given:         "1boiled egg",
			wantComponent: component("boiled egg", "1", ""),
		},

		{
			given:         "1 egg",
			wantComponent: component("egg", "1", ""),
		},
		{
			given:         "1/3 egg",
			wantComponent: component("egg", "1/3", ""),
		},
		{
			given:         "1 poached egg",
			wantComponent: component("poached egg", "1", ""),
		},
		{
			given:         "1/3 poached egg",
			wantComponent: component("poached egg", "1/3", ""),
		},
		{
			given:         "101 poached eggs",
			wantComponent: component("poached eggs", "101", ""),
		},
		{
			given:         "2/1 poached eggs",
			wantComponent: component("poached eggs", "2/1", ""),
		},
		{
			given:         "poached egg",
			wantComponent: component("poached egg", "", ""),
		},
		{
			given:         "egg",
			wantComponent: component("egg", "", ""),
		},
		{
			given:         "egg",
			wantComponent: component("egg", "", ""),
		},
		// sad paths
		{
			given:   "",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "/",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1/2",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1//2 egg",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1/2/ egg",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1/2/3 egg",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "ounce",
			wantErr: sozzler.ErrParseError,
		},
		{
			given:   "1.5 oz egg",
			wantErr: sozzler.ErrParseError,
		},
	}
	for i, tC := range testCases {
		t.Run(fmt.Sprintf("%d: %s", i, tC.given), func(t *testing.T) {
			parser := sozzler.RecipeParser{}

			c, err := parser.ParseComponent(strings.NewReader(tC.given))
			assert.ErrorIs(t, err, tC.wantErr)
			assert.Equal(t, tC.wantComponent, c)
		})
	}
}

func TestInvalidRecipe(t *testing.T) {
	for _, suffix := range []string{"0", "1", "0"} {
		file, err := os.Open(fmt.Sprintf("testdata/invalid%s.txt", suffix))
		require.NoError(t, err)
		defer func() { _ = file.Close() }()

		parser := sozzler.RecipeParser{}

		_, err = parser.Parse(file)
		assert.ErrorIs(t, err, sozzler.ErrParseError)
	}
}

func TestValidRecipe(t *testing.T) {
	for _, suffix := range []string{"0", "1", "2", "3", "4"} {
		var want sozzler.Recipe

		fileW, err := os.Open(fmt.Sprintf("testdata/recipe%s.yaml", suffix))
		require.NoError(t, err)
		defer func() { _ = fileW.Close() }()

		err = yaml.NewDecoder(fileW).Decode(&want)
		require.NoError(t, err)

		fileG, err := os.Open(fmt.Sprintf("testdata/recipe%s.txt", suffix))
		require.NoError(t, err)
		defer func() { _ = fileG.Close() }()

		p := sozzler.RecipeParser{}
		got, err := p.Parse(fileG)
		assert.NoError(t, err)

		assertSame(t, &want, got)
	}
}

func assertSame(t *testing.T, want *sozzler.Recipe, got *sozzler.Recipe) {
	t.Helper()
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Notes, got.Notes)
	assert.Equal(t, len(want.Components), len(got.Components))

	for _, wantC := range want.Components {
		foundIt := false
		for _, gotC := range got.Components {
			if wantC == gotC {
				foundIt = true
				break
			}
		}
		assert.True(t, foundIt, fmt.Sprintf("wanted %v in %v", wantC, got.Components))
	}
}

func must[T any](thing T, err error) T {
	if err != nil {
		panic(err)
	}
	return thing
}

func component(ingredient, quantity, unit string) *sozzler.Component {
	return &sozzler.Component{
		Ingredient: ingredient,
		Quantity:   *must(sozzler.ParseQuantity(quantity)),
		Unit:       unit,
	}
}
