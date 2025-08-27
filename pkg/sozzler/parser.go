package sozzler

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

type RecipeParser struct{}

type ParseError struct{}

func (e *ParseError) Error() string {
	return "Parse Error"
}

func (rp *RecipeParser) ParseComponent(r io.Reader) (*Component, error) {
	var s scanner.Scanner

	var numerator, denominator int
	var words []string

	s.Init(r)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		text := s.TokenText()

		if text == "/" {
			continue
		}

		n, err := strconv.Atoi(text)
		if err != nil {
			words = append(words, text)
			continue
		}
		if numerator == 0 {
			numerator = n
			continue
		}
		denominator = n
		continue
	}

	q, err := ParseQuantity(fmt.Sprintf("%d/%d", numerator, denominator))
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

	return &c, nil
}

func (rp *RecipeParser) Parse(r io.Reader) (*Recipe, error) {
	return nil, &ParseError{}
}
