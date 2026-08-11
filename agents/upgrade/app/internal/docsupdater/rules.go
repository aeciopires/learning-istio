package docsupdater

import "regexp"

// change records one applied rewrite, for reporting/dry-run diffs.
type change struct {
	Description string
	Before      string
	After       string
}

// applyRule runs re against content, replacing every match with
// replacement (a regexp.ReplaceAllString template, e.g. "${1}1.31").
// It returns the new content and, if anything changed, a change record
// summarizing the first match (representative of all matches, since a
// given rule typically fires on one logical value repeated verbatim).
func applyRule(content, description string, re *regexp.Regexp, replacement string) (string, *change) {
	if !re.MatchString(content) {
		return content, nil
	}
	first := re.FindString(content)
	newContent := re.ReplaceAllString(content, replacement)
	if newContent == content {
		return content, nil
	}
	after := re.FindString(newContent)
	return newContent, &change{Description: description, Before: first, After: after}
}

// istioLinkRe matches istio.io documentation links pinned to a specific
// "vMAJOR.MINOR" release, e.g. "istio.io/v1.30/docs/...". Links using
// "istio.io/latest/..." are intentionally not matched.
var istioLinkRe = regexp.MustCompile(`istio\.io/v\d+\.\d+/`)

func rewriteIstioLinks(content, istioMinor string) (string, *change) {
	if istioMinor == "" {
		return content, nil
	}
	return applyRule(content, "istio.io docs links pinned to vMAJOR.MINOR", istioLinkRe, "istio.io/v"+istioMinor+"/")
}
