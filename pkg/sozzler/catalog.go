package sozzler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type RecipeCatalog struct {
	Recipes []*Recipe
}

func (rc *RecipeCatalog) Find(name string) (*Recipe, bool) {
	name = strings.ToLower(name)
	for _, r := range rc.Recipes {
		if strings.ToLower(r.Name) == name {
			return r, true
		}
	}
	return nil, false
}

func (rc *RecipeCatalog) Load(recipesDir string) error {
	entries, err := os.ReadDir(recipesDir)
	if err != nil {
		return fmt.Errorf("couldn't list recipes in directory %q: %w", recipesDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(recipesDir, entry.Name())
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("couldn't open recipe file %q: %w", filename, err)
		}

		var recipe Recipe
		err = yaml.NewDecoder(file).Decode(&recipe)
		_ = file.Close()

		if err != nil {
			return fmt.Errorf("error decoding recipe %q: %w", filename, err)
		}

		rc.Recipes = append(rc.Recipes, &recipe)
	}
	return nil
}

type MatchResult struct {
	Predicate Predicate
	Match     string
}

func (rc *RecipeCatalog) Search(predicates []Predicate) map[*Recipe][]MatchResult {
	results := make(map[*Recipe][]MatchResult)

	for _, r := range rc.Recipes {
		for _, p := range predicates {
			if match, ok := p.Match(r); ok {
				if _, ok := results[r]; !ok {
					results[r] = make([]MatchResult, 0)
				}
				results[r] = append(results[r], MatchResult{
					Predicate: p,
					Match:     match,
				})
			}
		}
	}

	return results
}
