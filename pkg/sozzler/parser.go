package sozzler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"text/scanner"
)

type RecipeParser struct{}

var ErrParseError = errors.New("Parse Error")

func (rp *RecipeParser) ParseComponent(r io.Reader) (*Component, error) {
	var s scanner.Scanner

	var numerator, denominator int
	var words []string
	var slashed bool

	s.Init(r)
	s.Mode &^= scanner.ScanFloats // otherwise, 1egg is "1e" "gg"
	s.Mode |= scanner.ScanInts    // re-enable Int scanning disabled line above
	for tok := s.Scan(); tok != scanner.EOF && s.ErrorCount == 0; tok = s.Scan() {
		text := s.TokenText()

		if text == "/" {
			if slashed {
				return nil, ErrParseError
			}
			slashed = true
			continue
		}

		n, err := strconv.Atoi(text)
		if err != nil {
			// couldn't parse a digit, so it's a word
			words = append(words, text)
			continue
		}
		if len(words) != 0 {
			// it's an int, but we've already added a word, so input is like "1foo2"
			return nil, ErrParseError
		}
		if numerator == 0 {
			// numerator hasn't been set yet (probably)
			numerator = n
			continue
		}
		denominator = n
		continue
	}

	var q *Quantity
	var err error
	if denominator == 0 {
		if numerator == 0 {
			// ""
			q = &Quantity{}
		} else {
			// e.g. "1 foo"
			q, err = ParseQuantity(fmt.Sprint(numerator))
		}
	} else {
		q, err = ParseQuantity(fmt.Sprintf("%d/%d", numerator, denominator))
	}
	if err != nil {
		return nil, err
	}

	c := Component{Quantity: *q}

	if len(words) > 0 {
		if _, ok := knownUnits[words[0]]; ok {
			c.Unit = words[0]
			words = words[1:]
		}
	}

	c.Ingredient = strings.Join(words, " ")

	invalidComponent := Component{}
	if c == invalidComponent || c.Ingredient == "" {
		return nil, ErrParseError
	}

	return &c, nil
}

func (rp *RecipeParser) Parse(r io.Reader) (*Recipe, error) {
	scanner := bufio.NewScanner(r)

	var recipe Recipe
	var componentsDone bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(recipe.Components) > 0 {
				componentsDone = true
			}
			continue
		}
		if recipe.Name != "" && !componentsDone {
			if c, err := rp.ParseComponent(strings.NewReader(line)); err == nil {
				recipe.Components = append(recipe.Components, *c)
				continue
			}
		}
		if recipe.Name == "" {
			if len(recipe.Components) != 0 || componentsDone {
				// name can't come after components
				return nil, ErrParseError
			}
			recipe.Name = line
			continue
		}
		if len(recipe.Notes) != 0 {
			recipe.Notes += "\n"
		}
		recipe.Notes += line
		continue
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error during scanning: %v", err)
	}

	if recipe.Name == "" || len(recipe.Components) == 0 {
		return nil, ErrParseError
	}

	return &recipe, nil
}
