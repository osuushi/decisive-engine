package parser

import (
	"fmt"
	"strconv"
	"unicode"
)

// Parse a row format.
// Example: "Title: @title{15} Description: @desc"
// @ prefixed implies an interpolated field.
// {} indicates padding width
// If no padding width is given, that field expands to fill space
//   if multiple fields have no padding value, they will divide up
//   the free space
// @@ is used to print a literal @ (cannot appear in field key)
// Two literals cannot be adjacent or have their own padding value
// All spaces are included in literals.
// Extra spaces between adjacent fields are ignored; use padding instead

type Node struct {
	// If false, this is a literal string
	IsField bool
	// For fields, the value is the key (no @ included)
	// For literals, it is the literal value (with escaped @s unescaped)
	Value string
	// The amount of padding (for fields only)
	// zero means fill
	Padding int
}

func Parse(s string) ([]Node, error) {
	if len(s) == 0 {
		return make([]Node, 0), nil
	}

	firstNode, remainder, err := parseNode(s)
	if err != nil {
		return nil, err
	}
	nodes := []Node{*firstNode}

	remainingNodes, err := Parse(remainder)
	if err != nil {
		return nil, err
	}

	return append(nodes, remainingNodes...), nil
}

func parseNode(s string) (*Node, string, error) {
	if s[0] == '@' && len(s) > 1 && s[1] != '@' {
		return parseFieldNode(s)
	} else {
		return parseLiteralNode(s)
	}
}

func parseFieldNode(s string) (*Node, string, error) {
	result := ""
	var i int
	node := Node{IsField: true}
	for i = 1; i < len(s); i++ { // skip leading @
		r := s[i]
		switch r {
		case '{':
			if i >= len(s)-2 {
				return nil, "", fmt.Errorf("Unexpected end of string; expected padding followed by }")
			}
			// Chop off curly and parse
			remainder := s[i+1:]
			padding, remainder, err := parseFieldPadding(remainder)
			if err != nil {
				return nil, "", err
			}
			node.Value = result
			node.Padding = padding
			return &node, remainder, nil
		case ' ':
			node.Value = result
			return &node, s[i+1:], nil
		case '@':
			node.Value = result
			return &node, s[i:], nil
		default:
			result += string(r)
		}
	}
	// This is the end of the input
	node.Value = result
	return &node, "", nil
}

func parseFieldPadding(s string) (int, string, error) {
	result := ""
	for i, r := range s {
		// Handle close brace case
		if r == '}' {
			if result == "" {
				return 0, "", fmt.Errorf("Unexpected empty padding. Expected digits.")
			}
			padding, _ := strconv.Atoi(result)
			return padding, s[i+1:], nil
		}

		if unicode.IsDigit(r) {
			result += string(r)
		} else {
			return 0, "", fmt.Errorf("Unexpected character '%s'; expected digit", string(r))
		}
	}
	return 0, "", fmt.Errorf("Unexpected end of input in padding string")
}

func parseLiteralNode(s string) (*Node, string, error) {
	result := ""
	var i int
	node := Node{IsField: false}
	for i = 0; i < len(s); i++ {
		r := s[i]
		switch r {
		case '@':
			// Look ahead to see if this is an escape
			if i == len(s)-1 {
				return nil, "", fmt.Errorf("Unexpected end of input")
			} else if s[i+1] == '@' {
				// Escaped atsign
				result += string(r)
				// Skip next character
				i++
			} else {
				// Not an escape; field started; stop
				node.Value = result
				return &node, s[i:], nil
			}
		default:
			result += string(r)
		}
	}
	// Making it here means we're at the end of the input
	node.Value = result
	return &node, "", nil
}
