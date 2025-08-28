package display

import (
	"fmt"
	"mp/sozzler/pkg/sozzler"
)

type StdoutDisplay struct {
	Plain bool
}

func (d *StdoutDisplay) Error(e string) {
	fmt.Println(e)
}

func (d *StdoutDisplay) List(recipes []*sozzler.Recipe) {
	for _, r := range recipes {
		d.printName(r)
	}
}

func (d *StdoutDisplay) Show(recipe *sozzler.Recipe) {
	d.printName(recipe)
	for _, c := range sozzler.FancyOrder(recipe.Components) {
		if d.Plain {
			fmt.Println(c.Quantity, c.Unit, c.Ingredient)
		} else {
			fmt.Println("  ", c.Quantity, c.Unit, c.Ingredient)
		}
	}
	if !d.Plain {
		fmt.Println("---")
	}
	fmt.Println(recipe.Notes)
}

func (d *StdoutDisplay) printName(r *sozzler.Recipe) {
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
