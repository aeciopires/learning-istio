package releasenotes

import "strings"

// rstSection is one underlined heading of an Envoy version_history RST
// file, with its top-level bullet list flattened to plain strings.
type rstSection struct {
	Title   string
	Bullets []string
}

// parseRSTSections parses the small, consistent subset of RST used by
// Envoy's docs/root/version_history/*.rst files: a heading line
// immediately followed by a line of repeated underline characters
// (=, -, ~, ^, ", ', `), with "* " top-level bullets underneath
// (optionally wrapped onto indented continuation lines).
func parseRSTSections(raw string) []rstSection {
	lines := strings.Split(raw, "\n")

	headingAt := make(map[int]bool)
	skipLine := make(map[int]bool)
	for i := 1; i < len(lines); i++ {
		if !isUnderline(lines[i]) {
			continue
		}
		prev := lines[i-1]
		prevTrimmed := strings.TrimSpace(prev)
		if prevTrimmed == "" || strings.HasPrefix(strings.TrimLeft(prev, " "), "*") {
			continue
		}
		headingAt[i-1] = true
		skipLine[i] = true
	}

	var sections []rstSection
	var cur *rstSection
	var bulletBuf strings.Builder
	hasBullet := false

	flushBullet := func() {
		if cur != nil && hasBullet {
			text := strings.Join(strings.Fields(bulletBuf.String()), " ")
			if text != "" {
				cur.Bullets = append(cur.Bullets, text)
			}
		}
		bulletBuf.Reset()
		hasBullet = false
	}
	flushSection := func() {
		flushBullet()
		if cur != nil {
			sections = append(sections, *cur)
		}
		cur = nil
	}

	for i, line := range lines {
		if skipLine[i] {
			continue
		}
		if headingAt[i] {
			flushSection()
			cur = &rstSection{Title: strings.TrimSpace(line)}
			continue
		}
		if cur == nil {
			continue
		}

		trimmedLeft := strings.TrimLeft(line, " ")
		leadingSpaces := len(line) - len(trimmedLeft)

		switch {
		case strings.HasPrefix(trimmedLeft, "* ") && leadingSpaces <= 1:
			flushBullet()
			bulletBuf.WriteString(strings.TrimPrefix(trimmedLeft, "* "))
			hasBullet = true
		case strings.TrimSpace(line) == "":
			// Blank line: a following "* " will flush; otherwise ignore.
		default:
			if hasBullet {
				bulletBuf.WriteString(" ")
				bulletBuf.WriteString(trimmedLeft)
			}
			// Narrative text outside any bullet (section intro) is dropped.
		}
	}
	flushSection()
	return sections
}

func isUnderline(s string) bool {
	t := strings.TrimSpace(s)
	if len(t) < 3 {
		return false
	}
	first := rune(t[0])
	if !strings.ContainsRune("=-~^\"'`", first) {
		return false
	}
	for _, r := range t {
		if r != first {
			return false
		}
	}
	return true
}

func normalizeHeading(h string) string {
	return strings.ToLower(strings.TrimSpace(h))
}
