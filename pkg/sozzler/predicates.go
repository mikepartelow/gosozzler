package sozzler

import "strings"

type Predicate interface {
	Match(*Recipe) (string, bool)
	Name() string
}

// IngredientPredicate

type IngredientPredicate struct {
	ingredient string
}

func (ip *IngredientPredicate) Match(candidate *Recipe) (string, bool) {
	for _, c := range candidate.Components {
		if strings.Contains(strings.ToLower(c.Ingredient), ip.ingredient) {
			return c.Ingredient, true
		}
	}
	return "", false
}

func (ip *IngredientPredicate) Name() string {
	return "Ingredient"
}

func NewIngredientPredicate(ingredient string) Predicate {
	return &IngredientPredicate{
		ingredient: strings.TrimSpace(strings.ToLower(ingredient)),
	}
}

// NamePredicate

type NamePredicate struct {
	name string
}

func (ip *NamePredicate) Match(candidate *Recipe) (string, bool) {
	if strings.Contains(strings.ToLower(candidate.Name), ip.name) {
		return candidate.Name, true
	}
	return "", false
}

func (ip *NamePredicate) Name() string {
	return "Name"
}

func NewNamePredicate(name string) Predicate {
	return &NamePredicate{
		name: strings.TrimSpace(strings.ToLower(name)),
	}
}

// RatingPredicate

type RatingPredicate struct {
	rating int
}

func (rp *RatingPredicate) Match(candidate *Recipe) (string, bool) {
	if candidate.Rating >= rp.rating {
		return candidate.FancyRating(), true
	}
	return "", false
}

func (rp *RatingPredicate) Name() string {
	return "Rating"
}

func NewRatingPredicate(rating int) Predicate {
	return &RatingPredicate{
		rating: rating,
	}
}
