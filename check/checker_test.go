package check

import (
	"testing"
)

type test struct {
	word        string
	misspelled  bool
	suggested   []string
}

var (
	alphabet = *NewAlphabet("1234")
	words = []string{
		"1",
		"12",
		"123",
		"1234",
	}
)

func runTests(t *testing.T, checker Checker, initialWords []string, tests []test) {
	dict := NewTrie(alphabet)
	for _, word := range(initialWords) {
		dict.Add(word)
	}

	for _, test := range tests {
		res, sug := checker.IsMisspelled(test.word, dict)
		if res != test.misspelled {
			t.Errorf("expected misspelled=%v, received %v for word '%v'.", test.misspelled, res, test.word)
		}
		if len(sug) != len(test.suggested) {
			t.Errorf("expected suggested words %v, received %v.", test.suggested, sug)
		} else {
			for s := range(sug) {
				if sug[s] != test.suggested[s] {
					t.Errorf("expected suggested='%s', received '%s'.", test.suggested[s], sug[s])
				}
			}
		}
	}
}

func TestMinLengthChecker(t *testing.T) {
	var tests = []test{
		{"", false, nil},
		{"1", false, nil},
		{"11", false, nil},
		{"111", true, nil},
		{"1111", true, nil},
	}

	mlc := MinLengthChecker(3)
	runTests(t, mlc, words, tests)
}

func TestStrictChecker(t *testing.T) {
	var tests = []test{
		{"1", false, nil},
		{"12", false, nil},
		{"123", false, nil},
		{"122", true, nil},
		{"1234", false, nil},
		{"0", true, nil},
	}

	sc := StrictChecker{}
	runTests(t, sc, words, tests)
}

func TestDeltaCheckerInsertions(t *testing.T) {
	// single insertions
	var tests = []test{
		{"1", false, nil},
		{"3", false, nil},
		{"1233", false, nil},
		{"2", true, []string{"12"}},
		{"134", true, []string{"1234"}},
	}

	dc := DeltaChecker{AllowedIns:   1}

	runTests(t, dc, words, tests)

	// multiple insertions
	tests = []test{
		{"3", true, []string{"123"}},
		{"4", false, nil},
		{"13", true, []string{"123"}},
	}
	dc.AllowedIns = 2
	runTests(t, dc, words, tests)
}

func TestDeltaCheckerDeletions(t *testing.T) {
	// single deletions
	var tests = []test{
		{"1", false, nil},
		{"3", false, nil},
		{"1233", true, []string{"123"}},
		{"122", true, []string{"12"}},
		{"12344", true, []string{"1234"}},
		{"11234", true, []string{"1234"}},
	}

	dc := DeltaChecker{AllowedDel: 1}

	runTests(t, dc, words, tests)

	// multiple deletions
	tests = []test{
		{"13", true, []string{"1"}},
		{"133", true, []string{"1"}},
		{"1333", false, nil},
		{"11233", true, []string{"123"}},
		{"112333", false, nil},
	}
	dc.AllowedDel = 2
	runTests(t, dc, words, tests)
}

func TestDeltaCheckerSwaps(t *testing.T) {
	// single swaps
	var tests = []test{
		{"1", false, nil},
		{"2", false, nil},
		{"3", false, nil},
		{"1233", false, nil},
		{"21", true, []string{"12"}},
		{"321", false, nil},
		{"1243", true, []string{"1234"}},
	}

	dc := DeltaChecker{AllowedSwaps: 1}
	runTests(t, dc, words, tests)

	// multiple swaps
	tests = []test{
		{"21", true, []string{"12"}},
		{"321", false, nil},
		{"1423", true, []string{"1234"}},
		{"4123", false, nil},
	}
	dc.AllowedSwaps = 2
	runTests(t, dc, words, tests)
}

func TestDeltaCheckerMods(t *testing.T) {
	var tests = []test{
		{"12344", false, nil},
		{"333", false, nil},
		{"133", true, []string{"123"}},
	}

	dc := DeltaChecker{AllowedMods: 1}
	runTests(t, dc, words, tests)
}

func TestDeltaCheckerCombine(t *testing.T) {
	var tests = []test{
		{"1", false, nil},
		{"132", true, []string{"123"}},        // swapped chars           -> misspelled
		{"214", true, []string{"12"}},         // swapped + deletion      -> misspelled
		{"12334", true, []string{"1234"}},     // deletion                -> misspelled
		{"111", false, nil},                   // two deletions           -> not misspelled
		{"23", true, []string{"123"}},         // insertion               -> misspelled
		{"2234", true, []string{"1234"}},      // deletion + insertion    -> misspelled
		{"34", false, nil},                    // two insertions          -> not misspelled
	}

	dc := DeltaChecker{
		AllowedIns:   1,
		AllowedDel:   1,
		AllowedSwaps: 1}

	runTests(t, dc, words, tests)
}

func TestUnionChecker(t *testing.T) {
	var tests = []test{
		{"12", false, nil},
		{"132", true, []string{"123"}},
		{"134", true, []string{"1234"}},
	}

	dcOne := DeltaChecker{AllowedSwaps: 1}
	dcTwo := DeltaChecker{AllowedIns: 1}

	uc := unionChecker{[]Checker{dcOne, dcTwo}}

	runTests(t, uc, words, tests)
}

func TestIntersectChecker(t *testing.T) {
	var tests = []test{
		{"1224", true, []string{"1234"}},
		{"12344", false, nil},
		{"333", false, nil},
		{"133", false, nil},
	}

	lc := MinLengthChecker(4)
	mc := DeltaChecker{AllowedMods: 1}

	ic := intersectChecker{[]Checker{lc, mc}}

	runTests(t, ic, words, tests)
}
