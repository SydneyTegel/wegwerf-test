// Package main ist der Einstiegspunkt von wegwerf-test.
package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var whitespaceRe = regexp.MustCompile(`([dhms])\s+(\d)`)

var durationRe = regexp.MustCompile(`(\d+)([dhms])`)

// ParseDuration parses a duration string like "1h30m" or "2d" and returns the
// total number of seconds. Supported units: d (86400), h (3600), m (60), s (1).
func ParseDuration(s string) (int64, error) {
	s = strings.TrimSpace(s)
	// Entferne nur Whitespace nach einer Einheit und vor der nächsten Zahl.
	// Andere Whitespaces (z.B. "1 h") bleiben erhalten und werden von der Gap-Prüfung erkannt.
	s = whitespaceRe.ReplaceAllString(s, "$1$2")

	if s == "" {
		return 0, errors.New("leere Eingabe")
	}

	positions := durationRe.FindAllStringSubmatchIndex(s, -1)
	if len(positions) == 0 {
		return 0, fmt.Errorf("ungültige Einheit in %q", s)
	}

	// Check for gaps (unknown units like negative signs or unrecognized chars).
	for i := range positions {
		start := positions[i][0]
		end := positions[i][1]
		if i == 0 && start > 0 {
			return 0, fmt.Errorf("ungültige Zeichen am Anfang von %q", s)
		}
		if i > 0 && positions[i-1][1] != start {
			return 0, fmt.Errorf("ungültige Zeichen in %q", s)
		}
		if end != len(s) && i == len(positions)-1 {
			return 0, fmt.Errorf("ungültige Zeichen am Ende von %q", s)
		}
	}

	seen := make(map[byte]bool)
	var total int64

	for _, m := range positions {
		unit := s[m[4]]
		if seen[unit] {
			return 0, fmt.Errorf("doppelte Einheit %c", unit)
		}
		seen[unit] = true

		num, err := strconv.ParseInt(s[m[2]:m[3]], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("falsche Zahl in %q", s)
		}

		var mult int64
		switch unit {
		case 'd':
			mult = 86400
		case 'h':
			mult = 3600
		case 'm':
			mult = 60
		case 's':
			mult = 1
		}

		// Überlauf-Prüfung für Multiplikation
		if mult != 0 && num > math.MaxInt64/mult {
			return 0, fmt.Errorf("Überlauf: %d%c zu groß", num, unit)
		}
		product := num * mult

		// Überlauf-Prüfung für Addition
		if math.MaxInt64-total < product {
			return 0, fmt.Errorf("Überlauf: Summe überschreitet MaxInt64")
		}
		total += product
	}

	// Unterlauf defensiv prüfen
	if total < 0 {
		return 0, fmt.Errorf("Unterlauf: Ergebnis zu klein")
	}

	return total, nil
}

// Hello returns a greeting message.
func Hello() string {
	return "hello"
}

func main() {
	fmt.Println(Hello())
}
