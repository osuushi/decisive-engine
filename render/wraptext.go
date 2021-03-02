package render

import "strings"

// Wrap a string to a given width. Returns a slice wher for each line in the
// original string, there is a slice of wrapped lines. If a word is longer than
// an entire line, it is split.
//
// Tabs are not handled.
func wrapText(s string, width int) [][]string {
	// Consider original lines to be paragraphs
	paragraphs := strings.Split(s, "\n")
	output := make([][]string, len(paragraphs))
	for i, p := range paragraphs {
		output[i] = wrapLine(p, width)
	}
	return output
}

func lastIndexOf(haystack []rune, needle rune) int {
	for i := len(haystack) - 1; i >= 0; i-- {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

// Helper for wrapText, which only handles single-line strings.
func wrapLine(s string, width int) []string {
	runes := []rune(s)
	output := make([]string, 0, 1)

	// Consume runes until the last line
	for len(runes) > width {
		// First get an entire line plus one character. Note that because of the
		// loop condition, there will always be one more character
		lineRunes := runes[:width+1]
		lastSpaceIndex := lastIndexOf(lineRunes, ' ')

		if lastSpaceIndex == -1 {
			// Special case: the word is too big and must be split across lines
			output = append(output, string(lineRunes[:width]))
			// Notice that we don't consume the following character, since it isn't a space.
			runes = runes[width:]
		} else {
			// Get everything up to the space
			output = append(output, string(lineRunes[:lastSpaceIndex]))
			// Advance, consuming the space.
			runes = runes[lastSpaceIndex+1:]
		}
	}

	// Last line case (if there is one)
	if len(runes) > 0 {
		output = append(output, string(runes))
	}
	return output
}
