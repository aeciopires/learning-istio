package releasenotes

import "strings"

// classifyByKeywords inspects a heading/summary string and returns the
// best-matching Category. fallback is used when nothing matches.
func classifyByKeywords(text string, fallback Category) Category {
	t := strings.ToLower(text)
	switch {
	case containsAny(t, "helm", "install", "chart", "cni ", "cni:", "cni-"):
		return CategoryInstallation
	case containsAny(t, "crd", "gateway api", " api ", "apiversion", "resource"):
		return CategoryConfig
	case containsAny(t, "metric", "telemetry", "tracing", "trace", "grafana", "prometheus", "logging", "access log"):
		return CategoryMetrics
	case containsAny(t, "deprecat", "removed", "remove "):
		return CategoryDeprecated
	case containsAny(t, "ztunnel", "waypoint", "ambient", "istiod", "ingress gateway", "egress gateway", "sidecar inject"):
		return CategoryComponent
	case containsAny(t, "security", "auth", "cve", "tls", "mtls", "rbac", "jwt"):
		return CategorySecurity
	default:
		return fallback
	}
}

func containsAny(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}
