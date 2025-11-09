package server

import "strings"

func RenderMarkdownToTerminal(markdown string) string {
	var outLines []string
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "# ") {
			outLines = append(outLines, AnsiFgGreen+AnsiUnderline+renderCode(line[2:])+AnsiReset)
			continue
		}
		if strings.HasPrefix(line, "## ") {
			outLines = append(outLines, AnsiFgCyan+AnsiUnderline+renderCode(line[3:])+AnsiReset)
			continue
		}
		outLines = append(outLines, renderCode(line))
	}
	return strings.Join(outLines, "\n")
}

func renderCode(line string) string {
	var outLine strings.Builder
	inCode := false
	for _, char := range line {
		if char != '`' {
			outLine.WriteRune(char)
			continue
		}
		if inCode {
			outLine.WriteString(AnsiReset)
			inCode = false
			continue
		}
		inCode = true
		outLine.WriteString(AnsiFgYellow)
	}
	return outLine.String()
}
