// Package semverx implements a small, dependency-free semantic-version
// parser and range-expander tailored to how Istio and Envoy publish
// release notes: per MINOR series (e.g. "upgrade notes for 1.30") and per
// PATCH/hotfix release (e.g. "bug fixes in 1.30.2").
//
// It intentionally supports both "1.10" (minor-only) and "1.29.3"
// (minor+patch) inputs, because upgrade planning needs to work for both
// a wide minor-to-minor range (1.10 -> 1.30) and a narrow patch/hotfix
// range (1.29.0 -> 1.29.3).
package semverx

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a parsed Major.Minor[.Patch] version. HasPatch tells apart
// "1.10" (minor only, patch unspecified) from "1.10.0" (patch explicitly 0).
type Version struct {
	Major    int
	Minor    int
	Patch    int
	HasPatch bool
	Raw      string
}

// Parse accepts "1.10", "1.10.2", with an optional leading "v"/"V" and
// optional leading "release-" prefix (as used in istio branch names).
func Parse(s string) (Version, error) {
	raw := s
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "release-")
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimPrefix(s, "V")

	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return Version{}, fmt.Errorf("semverx: invalid version %q: expected MAJOR.MINOR or MAJOR.MINOR.PATCH", raw)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return Version{}, fmt.Errorf("semverx: invalid major component in %q", raw)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return Version{}, fmt.Errorf("semverx: invalid minor component in %q", raw)
	}

	v := Version{Major: major, Minor: minor, Raw: raw}
	if len(parts) == 3 {
		patch, err := strconv.Atoi(parts[2])
		if err != nil || patch < 0 {
			return Version{}, fmt.Errorf("semverx: invalid patch component in %q", raw)
		}
		v.Patch = patch
		v.HasPatch = true
	}
	return v, nil
}

// MinorString returns "MAJOR.MINOR", e.g. "1.30".
func (v Version) MinorString() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// String returns "MAJOR.MINOR.PATCH" if a patch was specified, else "MAJOR.MINOR".
func (v Version) String() string {
	if v.HasPatch {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return v.MinorString()
}

// minorKey returns a value that is monotonically increasing across
// Major.Minor combinations, suitable for ordering/iteration.
func (v Version) minorKey() int64 {
	return int64(v.Major)*100000 + int64(v.Minor)
}

// CompareMinor compares only the Major.Minor components: -1, 0, 1.
func CompareMinor(a, b Version) int {
	ak, bk := a.minorKey(), b.minorKey()
	switch {
	case ak < bk:
		return -1
	case ak > bk:
		return 1
	default:
		return 0
	}
}

// Compare compares Major.Minor.Patch fully. A version without an explicit
// patch is treated as patch 0 for ordering purposes.
func Compare(a, b Version) int {
	if c := CompareMinor(a, b); c != 0 {
		return c
	}
	ap, bp := a.Patch, b.Patch
	switch {
	case ap < bp:
		return -1
	case ap > bp:
		return 1
	default:
		return 0
	}
}

// MinorPlan describes what to fetch/report for one Major.Minor series
// within an expanded range.
type MinorPlan struct {
	Major int
	Minor int

	// IncludePatchDetail is true when the caller explicitly bounded this
	// minor series by a patch number on at least one side of the range
	// (e.g. "1.29.0"->"1.29.3", or the first/last minor of a wider range
	// when --from/--to included a patch). When false, only minor-level
	// upgrade/breaking-change notes are relevant for this series.
	IncludePatchDetail bool

	// FromPatch is the first patch (inclusive) to consider for this
	// minor series. Always 0 unless the range's lower bound landed on
	// this exact minor with an explicit patch.
	FromPatch int

	// ToPatch is the last patch (inclusive) to consider for this minor
	// series, or -1 meaning "unbounded / up to the latest known patch".
	ToPatch int
}

// MinorString returns "MAJOR.MINOR" for this plan entry.
func (p MinorPlan) MinorString() string {
	return fmt.Sprintf("%d.%d", p.Major, p.Minor)
}

// ExpandRange walks every Major.Minor series from "from" to "to"
// (inclusive) and returns one MinorPlan per series, correctly scoping
// patch-level detail to the boundary series when the caller supplied an
// explicit patch number. from/to may be given in either order; the
// smaller one is always treated as the lower bound.
func ExpandRange(from, to Version) ([]MinorPlan, error) {
	if CompareMinor(from, to) > 0 {
		from, to = to, from
	}

	var plans []MinorPlan
	major, minor := from.Major, from.Minor
	for {
		cur := Version{Major: major, Minor: minor}
		isFirst := cur.Major == from.Major && cur.Minor == from.Minor
		isLast := cur.Major == to.Major && cur.Minor == to.Minor

		plan := MinorPlan{Major: major, Minor: minor, ToPatch: -1}
		if isFirst && from.HasPatch {
			plan.FromPatch = from.Patch
			plan.IncludePatchDetail = true
		}
		if isLast && to.HasPatch {
			plan.ToPatch = to.Patch
			plan.IncludePatchDetail = true
		}
		plans = append(plans, plan)

		if isLast {
			break
		}
		minor++
		// Safety valve: never loop forever on malformed input.
		if len(plans) > 10000 {
			return nil, fmt.Errorf("semverx: refusing to expand a range spanning more than 10000 minor series (from=%s to=%s)", from, to)
		}
	}
	return plans, nil
}
