package sozzler

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

type RecipeParser struct{}

var ParseError = errors.New("Parse Error")

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
				return nil, ParseError
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
			return nil, ParseError
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
		return nil, ParseError
	}

	return &c, nil
}

func (rp *RecipeParser) Parse(r io.Reader) (*Recipe, error) {
	return nil, ParseError
}
