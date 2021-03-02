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
	pretty.Println(template)
	fmt.Println()
	row := render.NewRow(template, 150)
	pretty.Println(row.Render(map[string]interface{}{
		"foo": "This is some very long text that doesn't even fit on a single line",
		"bar": "world",
	}))
	pretty.Println(row.Render(map[string]interface{}{
		"foo": "this is some long text",
		"bar": "earth",
	}))
	fmt.Println(err)
}
