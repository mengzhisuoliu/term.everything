package main

import (
	"fmt"
	"strings"
	"unicode"
)

func genEnums(intf Interface) string {
	if len(intf.Enums) == 0 {
		return ""
	}

	var b strings.Builder
	for _, en := range intf.Enums {
		enumType := enumName(intf.Name, en.Name)

		fmt.Fprintf(&b, "type %s uint32\n\n", enumType)

		if len(en.Entries) > 0 {
			b.WriteString("const (\n")
			for _, e := range en.Entries {
				constName := sanitizeConstName(e.Name)
				fullName := enumType + "_" + constName
				fmt.Fprintf(&b, "    %s %s = %s\n", fullName, enumType, e.Value)
			}
			b.WriteString(")\n\n")
		}
	}

	return b.String()
}
func sanitizeConstName(s string) string {
	if s == "" {
		return "_"
	}
	var out strings.Builder
	for i, r := range s {
		switch {
		case i == 0 && unicode.IsDigit(r):
			out.WriteRune('_')
			out.WriteRune(r)
		case (i == 0 && (unicode.IsLetter(r) || r == '_')) ||
			(i > 0 && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_')):
			out.WriteRune(r)
		default:
			out.WriteRune('_')
		}
	}
	id := out.String()
	switch id {
	case "break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var":
		return id + "_"
	}
	return id
}
