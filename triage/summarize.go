package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"regexp"
	"strings"

	"k8s.io/test-infra/triage/editdistance"
)

var (
	flakeReasonDateRE = regexp.MustCompile(
		`[A-Z][a-z]{2}, \d+ \w+ 2\d{3} [\d.-: ]*([-+]\d+)?|` +
			`\w{3} \d{1,2} \d+:\d+:\d+(\.\d+)?|(\d{4}-\d\d-\d\d.|.\d{4} )\d\d:\d\d:\d\d(.\d+)?`)
	// Find random noisy strings that should be replaced with renumbered strings, for more similar messages.
	flakeReasonOrdinalRE = regexp.MustCompile(
		`0x[0-9a-fA-F]+` + // hex constants
			`|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}` + // IPs
			`|[0-9a-fA-F]{8}-\S{4}-\S{4}-\S{4}-\S{12}(-\d+)?` + // UUIDs + trailing digits
			`|[0-9a-f]{14,32}`) // hex garbage
	nameTagRE = regexp.MustCompile(`\[.*?\]|\{.*?\}`)

	ngramCountsCache = map[string][]int{}
)

// normalize reduces excess entropy in tracebacks to make clustering easier.
// This includes:
// - blanking dates and timestamps
// - renumbering unique information like
//     - pointer addresses
//     - UUIDs
//     - IP addresses
// - sorting randomly ordered map[] strings.
func normalize(s string) string {
	// blank out dates
	s = flakeReasonDateRE.ReplaceAllString(s, "")

	// do alpha conversion-- rename random garbage strings (hex pointer values, node names, etc)
	// into 'UNIQ1', 'UNIQ2', etc.
	matches := make(map[string]string)
	repl := func(m string) string {
		if matches[m] == "" {
			matches[m] = fmt.Sprintf("UNIQ%d", len(matches))
		}
		return matches[m]
	}

	s = flakeReasonOrdinalRE.ReplaceAllStringFunc(s, repl)

	if len(s) > 10000 { // for long strings, remove repeated lines!
		lines := strings.Split(s, "\n")
		newLines := make([]string, 0, len(lines))
		lastLine := ""
		for _, line := range lines {
			if line != lastLine {
				newLines = append(newLines, line)
			}
			lastLine = line
		}
		s = strings.Join(newLines, "\n")
	}

	if len(s) > 10000 { // ridiculously long test output
		s = s[:5000] + "\n...[truncated]...\n" + s[len(s)-5000:]
	}

	return s
}

// Given a test name, remove [...]/{...}.
// Matches code in testgrid and kubernetes/hack/update_owners.py.
func normalizeName(name string) string {
	name = nameTagRE.ReplaceAllString(name, "")
	return strings.TrimSpace(strings.Join(strings.Fields(name), " "))
}

// Convert a string into a histogram of frequencies for different byte combinations.
// This can be used as a heuristic to estimate edit distance between two strings in constant time.
// Instead of counting each ngram individually, they are hashed into buckets. This makes the output count size constant.
func makeNgramCounts(s string) []int {
	if b, ok := ngramCountsCache[s]; ok {
		return b
	}
	counts := [64]int{}
	for x := 0; x < len(s)-3; x++ {
		counts[crc32.ChecksumIEEE([]byte(s[x:x+4]))&63]++
	}
	ngramCountsCache[s] = counts[:]
	return counts[:]
}

/* Compute a heuristic lower-bound edit distance using ngram counts.

An insert/deletion/substitution can cause up to 4 ngrams to differ:

abcdefg => abcefg
(abcd, bcde, cdef, defg) => (abce, bcef, cefg)

This will underestimate the edit distance in many cases:
- ngrams hashing into the same bucket will get confused
- a large-scale transposition will barely disturb ngram frequencies,
but will have a very large effect on edit distance.

It is useful to avoid more expensive precise computations when they are
guaranteed to exceed some limit (being a lower bound), or as a proxy when
the exact edit distance computation is too expensive (for long inputs).
*/
func ngramEditDist(a, b string) int {
	countsA := makeNgramCounts(a)
	countsB := makeNgramCounts(b)
	sum := 0
	for i, ca := range countsA {
		cb := countsB[i]
		if ca < cb {
			sum += cb - ca
		} else {
			sum += ca - cb
		}
	}
	return sum
}

func makeNgramCountsDigest(s string) string {
	return "TODO"
}

//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterRunes = []rune("abc")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	n := 0
	for {
		aLen := rand.Intn(20) + 5
		bLen := rand.Intn(20) + 5
		a := RandStringRunes(aLen)
		b := RandStringRunes(bLen)
		brDist := editdistance.BerghelRoachDistance(a, b, 5)
		//lvDist := editdistance.LevenshteinDistance(a, b, 5)
		lvDist := brDist
		if brDist > 5 {
			if lvDist <= 5 {
				panic("FUCK")
			}
		} else if brDist != lvDist {
			panic("UG")
		}
		if brDist < 5 {
			n++
			if n&1023 == 0 {
				fmt.Println(a, b, brDist, lvDist)
			}
			if brDist != lvDist {
				panic("uguu")
			}
		}
	}
	//fmt.Println("WERUPU", editdistance.BerghelRoachDistance("foo", "football", 89))
}
