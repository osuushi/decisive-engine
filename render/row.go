package render

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/osuushi/decisive-engine/template"
	// "github.com/spf13/cast"
)

type Row struct {
	Template   template.Template
	InnerWidth int
	// Widths of each node
	Widths []int
}

func NewRow(template template.Template, innerWidth int) *Row {
	self := &Row{template, innerWidth, nil}
	self.init()
	return self
}

func (self *Row) init() {
	autoWidthIndexes := make([]int, 0)
	self.Widths = make([]int, len(self.Template))

	totalFixedWidth := 0
	// Add fixed widths
	for i, node := range self.Template {
		if node.IsField {
			if node.Width != 0 {
				totalFixedWidth += node.Width
				self.Widths[i] = node.Width
			} else {
				autoWidthIndexes = append(autoWidthIndexes, i)
			}
		} else {
			width := utf8.RuneCountInString(node.Value)
			totalFixedWidth += width
			self.Widths[i] = width
		}
	}
	// TODO: Figure out what to do if the entire row width is too small to account
	// for fixed widths

	// Allocate free width
	for len(autoWidthIndexes) > 0 {
		remaining := len(autoWidthIndexes)
		var index int
		index, autoWidthIndexes = autoWidthIndexes[0], autoWidthIndexes[1:]

		// Whatever fixed width is left, divide it by the remaining buckets and
		// round to get the next allocation. In the trivial case where
		// freeWidth/buckets is an integer, this divides the space up evenly between
		// them. If not, the error is spread as uniformly as possible among the buckets.
		//
		// Note that while the error is spread "flat", it is not distributed evenly
		// in terms of array position, instead tending to clump at the ends of the slice
		width := int(math.Round(float64(totalFixedWidth) / float64(remaining)))
		self.Widths[index] = width

		// Subtract off what we allocated.
		totalFixedWidth -= width
	}
}

func (self *Row) render(data map[string]interface{}) string {
	// Builder
	parts := make([]string, len(self.Template))

	for i, node := range self.Template {
		if node.IsField {
			parts[i] = self.RenderFieldNodeAtIndex(i, data)
		} else {
			parts[i] = node.Value
		}
	}

	return strings.Join(parts, "")
}

func (self *Row) RenderFieldNodeAtIndex(index int, data map[string]interface{}) string {
	node := self.Template[index]
	key := node.Value
}