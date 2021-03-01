package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/osuushi/decisive-engine/parser"
)

func main() {
	input := os.Args[1]
	fmt.Println("Input:", input)
	fmt.Println()
	result, err := parser.Parse(input)
	pretty.Println(result)
	fmt.Println(err)
}
