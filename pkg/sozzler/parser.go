package sozzler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type MarkdownParser struct {
}

func (p *MarkdownParser) Parse(r io.Reader) (*Recipe, error) {
	scanner := bufio.NewScanner(r)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	var name string
	var components []Component
	var notesLines []string

	var consumed int

	// 1) Name: first H1 "# "
	for _, ln := range lines {
		consumed += 1
		ln = cleanLine(ln)
		if strings.HasPrefix(ln, "# ") {
			name = strings.TrimSpace(strings.TrimPrefix(ln, "# "))
			break
		}
	}

	// 2) Components: bullet lines starting with "- "
	inComponents := false
	for _, raw := range lines[consumed:len(lines)] {
		ln := cleanLine(raw)
		if ln == "" {
			// empty lines don't toggle state, but once we've seen a component list,
			// an empty line may signal the end of a contiguous list.
			if inComponents {
				// keep going; we only stop when we hit a non-bullet, non-empty after components started.
			}
			continue
		}

		if strings.HasPrefix(ln, "- ") {
			inComponents = true
			item := strings.TrimSpace(strings.TrimPrefix(ln, "- "))

			// tokenization
			toks := strings.Fields(item)
			if len(toks) == 0 {
				continue
			}

			// quantity
			qty, err := ParseQuantity(toks[0])
			if err != nil {
				// not a measure line; treat entire bullet as note
				notesLines = append(notesLines, strings.TrimPrefix(ln, "- "))
				continue
			}

			// Look for "of" to split units vs ingredient
			ofIdx := -1
			for i := 1; i < len(toks); i++ {
				if strings.ToLower(toks[i]) == "of" {
					ofIdx = i
					break
				}
			}

			var units string
			var ingredient string

			switch {
			case ofIdx >= 0:
				// units are toks[1:ofIdx], ingredient toks[ofIdx+1:]
				if ofIdx > 1 {
					units = strings.Join(toks[1:ofIdx], " ")
				}
				if ofIdx+1 < len(toks) {
					ingredient = strings.Join(toks[ofIdx+1:], " ")
				}
			default:
				// No "of": if toks[1] is a known unit, treat it as unit; else, no unit.
				if len(toks) >= 2 && isUnitWord(toks[1]) {
					units = toks[1]
					if len(toks) >= 3 {
						ingredient = strings.Join(toks[2:], " ")
					}
				} else {
					// No unit; ingredient is rest
					if len(toks) >= 2 {
						ingredient = strings.Join(toks[1:], " ")
					}
				}
			}

			components = append(components, Component{
				Ingredient: ingredient,
				Quantity:   *qty,
				Units:      units,
			})
			continue
		}

		// Non-bullet
		if inComponents {
			// We've exited the contiguous bullet list; remaining lines are notes unless a separator '---'
			// fallthrough to notes handling
		}
		if strings.HasPrefix(ln, "---") {
			// treat as explicit end separator; skip the separator line itself
			continue
		}
		notesLines = append(notesLines, ln)
	}

	// 3) Notes: join remaining text blocks (after components), trim extra blank lines
	notes := strings.Join(notesLines, "\n")
	notes = strings.TrimSpace(notes)

	return &Recipe{
		Name:       name,
		Notes:      notes,
		Components: components,
		Rating:     0,
	}, nil
}

func isUnitWord(w string) bool {
	w = strings.ToLower(strings.TrimSpace(w))
	unitSet := map[string]struct{}{
		"oz": {}, "ounce": {}, "ounces": {},
		"ml": {}, "cl": {}, "l": {},
		"tsp": {}, "teaspoon": {}, "teaspoons": {},
		"tbsp": {}, "tablespoon": {}, "tablespoons": {},
		"cup": {}, "cups": {},
		"dash": {}, "dashes": {},
		"pinch": {}, "pinches": {},
		"g": {}, "gram": {}, "grams": {},
		"kg": {}, "lb": {}, "pound": {}, "pounds": {},
		"slice": {}, "slices": {},
		"wheel": {}, "wheels": {},
	}
	_, ok := unitSet[w]
	return ok
}

func cleanLine(s string) string {
	// Normalize whitespace and remove trailing punctuation that isn't helpful
	s = strings.TrimSpace(s)
	// Replace fancy dashes with hyphen and normalize spaces
	s = strings.Map(func(r rune) rune {
		if r == '–' || r == '—' {
			return '-'
		}
		// keep printable runes
		if r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(s)
}
