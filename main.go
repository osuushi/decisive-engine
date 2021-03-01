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

	pretty.Println(row.Render(nil))
	fmt.Println(err)
}
