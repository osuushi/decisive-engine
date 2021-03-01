package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/osuushi/decisive-engine/template"
)

func main() {
	input := os.Args[1]
	fmt.Println("Input:", input)
	fmt.Println()
	result, err := template.Parse(input)
	pretty.Println(result)
	fmt.Println(err)
}
