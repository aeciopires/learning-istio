package releasenotes

import "strings"

// mdSection is one "## Heading" block of a Hugo/markdown release-notes page.
type mdSection struct {
	Title string
	Body  string
}

// splitMarkdownSections strips YAML frontmatter (if present) and splits
// the remaining markdown body on "## " (H2) headings. Content before the
// first H2 (typically a lede paragraph) is dropped, matching how Istio's
// announcement pages are structured (frontmatter, lede, then one H2 per
// notable change).
func splitMarkdownSections(raw string) []mdSection {
	body := stripFrontmatter(raw)
	lines := strings.Split(body, "\n")

	var sections []mdSection
	var cur *mdSection
	var buf []string

	flush := func() {
		if cur != nil {
			cur.Body = strings.TrimSpace(strings.Join(buf, "\n"))
			sections = append(sections, *cur)
		}
		buf = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			flush()
			title := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			cur = &mdSection{Title: title}
			continue
		}
		if strings.HasPrefix(line, "# ") {
			// Top-level H1: not a change section, ignore its body.
			flush()
			cur = nil
			continue
		}
		if cur != nil {
			buf = append(buf, line)
		}
	}
	flush()
	return sections
}

// bomAndWhitespace lists the leading runes stripped before checking for a
// "---" frontmatter fence: a UTF-8 byte-order-mark (some editors/CDNs
// prepend one; written as  to avoid embedding a literal BOM byte
// sequence in the source file, which Go's parser rejects) plus ordinary
// whitespace.
const bomAndWhitespace = "\uFEFF \n\r\t"

func stripFrontmatter(raw string) string {
	s := strings.TrimLeft(raw, bomAndWhitespace)
	if !strings.HasPrefix(s, "---") {
		return raw
	}
	rest := s[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return raw
	}
	after := rest[idx+len("\n---"):]
	if nl := strings.IndexByte(after, '\n'); nl != -1 {
		return after[nl+1:]
	}
	return ""
}
