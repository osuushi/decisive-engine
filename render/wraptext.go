package render

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/osuushi/decisive-engine/template"
)

var wordSplitPattern *regexp.Regexp

func init() {
	var err error
	wordSplitPattern, err = regexp.Compile("\\s+")
	if err != nil {
		panic(err)
	}
}

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

// Take the output of wrapText and apply an alignment to it. All output strings
// inside the structure will be exactly `width` when returned
func alignText(paragraphs [][]string, width int, alignment template.Alignment) [][]string {
	output := make([][]string, len(paragraphs))
	for i, p := range paragraphs {
		output[i] = alignParagraph(p, width, alignment)
	}
	return output
}

func alignParagraph(lines []string, width int, alignment template.Alignment) []string {
	output := make([]string, len(lines))
	i := 0

	// Last line gets special treatment
	for ; i < len(lines)-1; i++ {
		output[i] = alignLine(lines[i], width, alignment)
	}

	// If justified, last line is left aligned
	if alignment == template.AlignmentJustify {
		alignment = template.AlignmentLeft
	}

	output[i] = alignLine(lines[i], width, alignment)
	return output
}

func alignLine(s string, width int, alignment template.Alignment) string {
	switch alignment {
	case template.AlignmentDefault, template.AlignmentLeft:
		return alignLineLeft(s, width)
	case template.AlignmentRight:
		return alignLineRight(s, width)
	case template.AlignmentCenter:
		return alignLineCenter(s, width)
	case template.AlignmentJustify:
		return alignLineJustify(s, width)
	default:
		panic(fmt.Sprintf("Unknown alignment value %v", alignment))
	}
}

func alignLineLeft(s string, width int) string {
	return s + strings.Repeat(" ", spacesNeeded(s, width))
}

func alignLineRight(s string, width int) string {
	return strings.Repeat(" ", spacesNeeded(s, width)) + s
}

func alignLineCenter(s string, width int) string {
	needed := spacesNeeded(s, width)
	// Note that in the case of odd spacing, we have to put the remainder
	// somewhere
	left := needed / 2
	right := needed - left
	return strings.Join([]string{
		strings.Repeat(" ", left),
		s,
		strings.Repeat(" ", right),
	}, "")
}

func alignLineJustify(s string, width int) string {
	// Split into words
	words := wordSplitPattern.Split(s, -1)

	gapCount := len(words) - 1

	parts := make([]string, len(words)+gapCount)

	totalLength := 0
	for _, word := range words {
		totalLength += utf8.RuneCountInString(word)
	}

	// Fractional width we'd like a gap to be; errors will diffuse into this value
	idealGapWidth := float64(width-totalLength) / float64(gapCount)
	lastError := 0.0
	// Divide amount of free width between gaps between words
	var i int
	for i = 0; i < len(words)-1; i++ {
		parts[2*i] = words[i]
		// Diffuse error
		correctedWidth := idealGapWidth - lastError
		actualGapWidth := int(math.Round(correctedWidth))
		lastError = float64(actualGapWidth) - correctedWidth
		parts[2*i+1] = strings.Repeat(" ", actualGapWidth)
	}
	parts[2*i] = words[i]

	return strings.Join(parts, "")
}

func spacesNeeded(s string, width int) int {
	return width - utf8.RuneCountInString(s)
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
