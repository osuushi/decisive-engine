package render

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestWrapLine(t *testing.T) {
	check := func(input string, width int, expected []string) {
		actual := wrapLine(input, width)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf(
				"wrapLine(%#v, %#v)\nExpected:\n%s\nActual:\n%s",
				input,
				width,
				pretty.Sprint(expected),
				pretty.Sprint(actual),
			)
		}
	}

	check("this is a test", 9, []string{"this is a", "test"})
	check("this is a test", 8, []string{"this is", "a test"})
	check("this is a test", 7, []string{"this is", "a test"})
	check("this is a test", 6, []string{"this", "is a", "test"})
	check("this is a test", 4, []string{"this", "is a", "test"})
	check("thís is a test", 4, []string{"thís", "is a", "test"})
	check("excellent, a test, how entertaining", 6,
		[]string{"excell", "ent, a", "test,", "how", "entert", "aining"})
}
