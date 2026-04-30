package rapid

import (
	"testing"
	"unicode"
)

func TestRestrictedRuneGen_CharClass(t *testing.T) {
	printableASCII := &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0020, Hi: 0x007E, Stride: 1},
		},
	}

	tests := []struct {
		name    string
		pattern string
		allowed []*unicode.RangeTable
		wantMin rune
		wantMax rune
	}{
		{"digits only", "[0-9]", []*unicode.RangeTable{printableASCII}, '0', '9'},
		{"lowercase only", "[a-z]", []*unicode.RangeTable{printableASCII}, 'a', 'z'},
		{"uppercase only", "[A-Z]", []*unicode.RangeTable{printableASCII}, 'A', 'Z'},
		{"alphanumeric", "[a-zA-Z0-9]", []*unicode.RangeTable{printableASCII}, '0', 'z'},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen := StringMatchingWithRunes(tc.pattern, tc.allowed...)

			for seed := uint64(0); seed < 200; seed++ {
				st := NewSeededT(seed)
				v := gen.Draw(st, "rune")

				if len(v) != 1 {
					t.Fatalf("seed=%d: expected 1 char, got %d: %q", seed, len(v), v)
				}

				r := rune(v[0])
				if r < tc.wantMin || r > tc.wantMax {
					t.Errorf("seed=%d: got %q (0x%04X), want [%c-%c]",
						seed, r, r, tc.wantMin, tc.wantMax)
				}
			}
		})
	}
}

func TestStringMatchingWithRunes_Simple(t *testing.T) {
	pattern := `[0-9]+-[a-z]+-[0-9]+`

	printableASCII := &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0020, Hi: 0x007E, Stride: 1},
		},
	}

	gen := StringMatchingWithRunes(pattern, printableASCII)

	for seed := uint64(0); seed < 20; seed++ {
		st := NewSeededT(seed)
		v := gen.Draw(st, "simple")
		t.Logf("seed=%d value=%q", seed, v)
	}
}

func TestStringMatchingWithRunes_ARN(t *testing.T) {
	pattern := `arn:[a-z\d-]+:kinesisvideo:[a-z0-9-]+:[0-9]+:[a-z]+/[a-zA-Z0-9_.-]+/[0-9]+`

	printableASCII := &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0020, Hi: 0x007E, Stride: 1},
		},
	}

	gen := StringMatchingWithRunes(pattern, printableASCII)

	for seed := uint64(0); seed < 20; seed++ {
		st := NewSeededT(seed)
		v := gen.Draw(st, "arn")
		t.Logf("seed=%d value=%q", seed, v)
	}
}
