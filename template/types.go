package template

import "image/color"

type Node struct {
	// If false, this is a literal string
	IsField bool
	// For fields, the value is the key (no @ included)
	// For literals, it is the literal value (with escaped @s unescaped)
	Value string
	// The amount of width to pad to (for fields only)
	// zero means fill
	Width int
	// Formatting; nil for literals
	FieldFormatting *FieldFormatting
}

type Template []Node

type FieldFormatting struct {
	// Color for field node; nil means default
	Color *color.RGBA
	// Styling
	Bold, Italic, Underline bool
}
