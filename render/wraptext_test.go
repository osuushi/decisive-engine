package render

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/osuushi/decisive-engine/template"
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

func TestAlignLineJustify(t *testing.T) {
	check := func(input string, width int, expected string) {
		actual := alignLineJustify(input, width)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf(
				"alignTextJustify(%#v, %#v)\nExpected:\n%#v\nActual:\n%#v",
				input,
				width,
				expected,
				actual,
			)
		}
	}
	check("hello world", 13, "hello   world")
	check(
		"Is this the real life? Is this just fantasy?", 60,
		"Is   this   the   real   life?   Is   this   just   fantasy?",
	)
	check(
		"Is this the real life? Is this just fantasy?", 58,
		"Is   this   the  real   life?   Is   this  just   fantasy?",
	)
	check(
		"Is this the real life? Is this just fantasy?", 55,
		"Is  this   the  real   life?  Is  this   just  fantasy?",
	)

	check(
		"Because I'm easy come easy go", 53,
		"Because      I'm      easy     come      easy      go",
	)

	check(
		"Because I'm easy come eásy go", 53,
		"Because      I'm      easy     come      eásy      go",
	)

	check(
		"This \t sentence  has    weird  spaces", 38,
		"This   sentence   has   weird   spaces",
	)
}

func TestAlignParagraph(t *testing.T) {
	check := func(
		input []string,
		width int,
		alignment template.Alignment,
		expected []string,
	) {
		actual := alignParagraph(input, width, alignment)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf(
				"alignParagraph(%#v, %#v)\nExpected:\n%s\nActual:\n%s",
				input,
				width,
				"|"+strings.Join(expected, "|\n|")+"|",
				"|"+strings.Join(actual, "|\n|")+"|",
			)
		}
	}
	input := []string{
		"Whan that Aprille with his shoures soote,",
		"The droghte of March hath perced to the roote,",
		"And bathed every veyne in swich licóur",
		"Of which vertú engendred is the flour;",
	}

	check(input, 46, template.AlignmentDefault, []string{
		"Whan that Aprille with his shoures soote,     ",
		"The droghte of March hath perced to the roote,",
		"And bathed every veyne in swich licóur        ",
		"Of which vertú engendred is the flour;        ",
	})

	check(input, 46, template.AlignmentLeft, []string{
		"Whan that Aprille with his shoures soote,     ",
		"The droghte of March hath perced to the roote,",
		"And bathed every veyne in swich licóur        ",
		"Of which vertú engendred is the flour;        ",
	})

	check(input, 46, template.AlignmentRight, []string{
		"     Whan that Aprille with his shoures soote,",
		"The droghte of March hath perced to the roote,",
		"        And bathed every veyne in swich licóur",
		"        Of which vertú engendred is the flour;",
	})

	check(input, 46, template.AlignmentJustify, []string{
		"Whan  that  Aprille with  his  shoures  soote,",
		"The droghte of March hath perced to the roote,",
		"And  bathed   every  veyne  in   swich  licóur",
		"Of which vertú engendred is the flour;        ",
	})
}
