package sozzler

import (
	"sort"
)

type Component struct {
	Ingredient string   `yaml:"ingredient"`
	Quantity   Quantity `yaml:"quantity"`
	Unit       string   `yaml:"unit"`
}

type Recipe struct {
	Name       string      `yaml:"name"`
	Notes      string      `yaml:"text"`
	Components []Component `yaml:"components"`
	Rating     int         `yaml:"rating"`
}

func (r *Recipe) FancyRating() string {
	rating := ""
	for i := 0; i < r.Rating; i++ {
		rating += "🫒"
	}
	return rating
}

func FancyOrder(components []Component) []Component {
	sort.Slice(components, func(i, j int) bool {
		ci, cj := components[i], components[j]
		if ci.Unit == cj.Unit {
			if ci.Quantity == cj.Quantity {
				return ci.Ingredient < cj.Ingredient
			}
			return ci.Quantity.Float() > cj.Quantity.Float()
		}
		if ci.Unit == "" {
			return false
		}
		if cj.Unit == "" {
			return true
		}
		// FIXME: have some conversion table that normalizes all quantities into grams
		return false
	})

	return components
}
