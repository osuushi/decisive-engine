package main

import (
	"fmt"

	"github.com/osuushi/decisive-engine/parser"
)

func main() {
	result, err := parser.Parse("@foo{3} bar: @baz{5} @atsign@")
	fmt.Printf("%#v\n", result)
	fmt.Println(err)
}
