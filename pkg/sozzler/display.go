package sozzler

import (
	"fmt"
)

type Display interface {
	Error(string)
	List([]*Recipe)
	Show(*Recipe)
	String(string)
}

// StdoutDisplay
type StdoutDisplay struct {
	Plain bool
}

func (d *StdoutDisplay) Error(e string) {
	fmt.Println(e)
}

func (d *StdoutDisplay) List(recipes []*Recipe) {
	for _, r := range recipes {
		d.printName(r)
	}
}

func (d *StdoutDisplay) Show(recipe *Recipe) {
	d.printName(recipe)
	for _, c := range FancyOrder(recipe.Components) {
		if d.Plain {
			fmt.Println(c.Quantity, c.Units, c.Ingredient)
		} else {
			fmt.Println("  ", c.Quantity, c.Units, c.Ingredient)
		}
	}
	if !d.Plain {
		fmt.Println("---")
	}
	fmt.Println(recipe.Notes)
}

func (d *StdoutDisplay) printName(r *Recipe) {
	if d.Plain {
		fmt.Println(r.Name)
		rating := ""
		for i := 0; i < r.Rating; i++ {
			rating += "*"
		}
		fmt.Println(rating)
	} else {
		fmt.Println("🍸", r.Name, r.FancyRating())
	}
}

func (d *StdoutDisplay) String(s string) {
	fmt.Print(s)
}
