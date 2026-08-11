package semverx

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in      string
		want    Version
		wantErr bool
	}{
		{"1.10", Version{Major: 1, Minor: 10, Raw: "1.10"}, false},
		{"1.10.0", Version{Major: 1, Minor: 10, Patch: 0, HasPatch: true, Raw: "1.10.0"}, false},
		{"v1.29.3", Version{Major: 1, Minor: 29, Patch: 3, HasPatch: true, Raw: "v1.29.3"}, false},
		{"release-1.30", Version{Major: 1, Minor: 30, Raw: "release-1.30"}, false},
		{"1", Version{}, true},
		{"1.2.3.4", Version{}, true},
		{"a.b", Version{}, true},
		{"", Version{}, true},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("Parse(%q): expected error, got %+v", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q): unexpected error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Parse(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestVersionString(t *testing.T) {
	v, _ := Parse("1.10")
	if got := v.String(); got != "1.10" {
		t.Errorf("String() = %q, want %q", got, "1.10")
	}
	v2, _ := Parse("1.10.2")
	if got := v2.String(); got != "1.10.2" {
		t.Errorf("String() = %q, want %q", got, "1.10.2")
	}
	if got := v2.MinorString(); got != "1.10" {
		t.Errorf("MinorString() = %q, want %q", got, "1.10")
	}
}

func TestCompare(t *testing.T) {
	a, _ := Parse("1.10.2")
	b, _ := Parse("1.10.3")
	c, _ := Parse("1.11")
	if Compare(a, b) >= 0 {
		t.Errorf("expected 1.10.2 < 1.10.3")
	}
	if CompareMinor(a, c) >= 0 {
		t.Errorf("expected 1.10 < 1.11")
	}
	if Compare(a, a) != 0 {
		t.Errorf("expected equal versions to compare 0")
	}
}

func TestExpandRange_MinorOnly(t *testing.T) {
	from, _ := Parse("1.10")
	to, _ := Parse("1.12")
	plans, err := ExpandRange(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plans) != 3 {
		t.Fatalf("expected 3 minor series, got %d: %+v", len(plans), plans)
	}
	wantMinors := []string{"1.10", "1.11", "1.12"}
	for i, p := range plans {
		if p.MinorString() != wantMinors[i] {
			t.Errorf("plans[%d] = %s, want %s", i, p.MinorString(), wantMinors[i])
		}
		if p.IncludePatchDetail {
			t.Errorf("plans[%d]: expected no patch detail for minor-only range, got %+v", i, p)
		}
	}
}

func TestExpandRange_SingleMinorPatchScoped(t *testing.T) {
	from, _ := Parse("1.29.0")
	to, _ := Parse("1.29.3")
	plans, err := ExpandRange(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("expected 1 minor series, got %d: %+v", len(plans), plans)
	}
	p := plans[0]
	if p.MinorString() != "1.29" {
		t.Fatalf("expected minor 1.29, got %s", p.MinorString())
	}
	if !p.IncludePatchDetail {
		t.Fatalf("expected patch detail to be included, got %+v", p)
	}
	if p.FromPatch != 0 || p.ToPatch != 3 {
		t.Errorf("expected FromPatch=0 ToPatch=3, got FromPatch=%d ToPatch=%d", p.FromPatch, p.ToPatch)
	}
}

func TestExpandRange_WideRangeWithPatchBoundaries(t *testing.T) {
	from, _ := Parse("1.10.5")
	to, _ := Parse("1.12.1")
	plans, err := ExpandRange(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plans) != 3 {
		t.Fatalf("expected 3 minor series, got %d", len(plans))
	}

	first := plans[0]
	if !first.IncludePatchDetail || first.FromPatch != 5 || first.ToPatch != -1 {
		t.Errorf("first minor: expected FromPatch=5 ToPatch=-1 IncludePatchDetail=true, got %+v", first)
	}

	mid := plans[1]
	if mid.IncludePatchDetail {
		t.Errorf("middle minor: expected no patch detail, got %+v", mid)
	}

	last := plans[2]
	if !last.IncludePatchDetail || last.FromPatch != 0 || last.ToPatch != 1 {
		t.Errorf("last minor: expected FromPatch=0 ToPatch=1 IncludePatchDetail=true, got %+v", last)
	}
}

func TestExpandRange_ReversedOrder(t *testing.T) {
	from, _ := Parse("1.30")
	to, _ := Parse("1.28")
	plans, err := ExpandRange(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"1.28", "1.29", "1.30"}
	for i, p := range plans {
		if p.MinorString() != want[i] {
			t.Errorf("plans[%d] = %s, want %s", i, p.MinorString(), want[i])
		}
	}
}

func TestExpandRange_SameVersion(t *testing.T) {
	v, _ := Parse("1.30.2")
	plans, err := ExpandRange(v, v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(plans))
	}
	if plans[0].FromPatch != 2 || plans[0].ToPatch != 2 {
		t.Errorf("expected FromPatch=ToPatch=2, got %+v", plans[0])
	}
}
