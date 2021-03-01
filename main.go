package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/osuushi/decisive-engine/render"
	"github.com/osuushi/decisive-engine/template"
)

func main() {
	input := os.Args[1]
	fmt.Println("Input:", input)
	fmt.Println()
	template, err := template.Parse(input)
	row := render.NewRow(template, 100)
	pretty.Println(row.Render(map[string]interface{}{
		"foo": "142",
		"bar": "world",
	}))
	pretty.Println(row.Render(map[string]interface{}{
		"foo": "long",
		"bar": "world",
	}))
	fmt.Println(err)
}
