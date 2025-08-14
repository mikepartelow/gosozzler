package sozzler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type FractionFloat64 float64

func (f FractionFloat64) String() string {
	if float64(f) == 0 {
		return ""
	}

	fractionMap := map[float64]string{
		0.5:  "1/2", // "½"
		0.25: "1/4", // "¼"
		0.75: "3/4", // "¾"
	}

	val := float64(f)

	intPart := int(val)
	fracPart := val - float64(intPart)

	if fancy, ok := fractionMap[fracPart]; ok {
		if intPart == 0 {
			return fancy
		}
		return fmt.Sprintf("%d %s", intPart, fancy)
	}

	return fmt.Sprint(val)
}

func (f *FractionFloat64) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("quantity: unsupported YAML type %q", value.Tag)
	}

	v, err := parseFraction(s)
	if err != nil {
		return fmt.Errorf("quantity: %w", err)
	}
	*f = FractionFloat64(v)
	return nil
}

func parseFraction(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fraction %q", s)
	}
	num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil || den == 0 {
		return 0, fmt.Errorf("invalid fraction %q", s)
	}
	return num / den, nil
}

type Component struct {
	Ingredient string          `yaml:"ingredient"`
	Quantity   FractionFloat64 `yaml:"quantity"`
	Units      string          `yaml:"unit"`
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
		if ci.Units == cj.Units {
			if ci.Quantity == cj.Quantity {
				return ci.Ingredient < cj.Ingredient
			}
			return ci.Quantity > cj.Quantity
		}
		if ci.Units == "" {
			return false
		}
		if cj.Units == "" {
			return true
		}
		// FIXME: have some conversion table that normalizes all quantities into grams
		return false
	})

	return components
}
