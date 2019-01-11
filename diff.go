package main

import (
	"fmt"
	"strings"

	"gopkg.in/d4l3k/messagediff.v1"
)

func computeDiff(old, new map[string]interface{}) (string, bool) {
	var diffs strings.Builder
	equal := true

	diff, _ := messagediff.DeepDiff(old, new)

	for p, a := range diff.Added {
		equal = false
		s := fmt.Sprintf("\t+ %v = %v\n", p.String(), a)
		diffs.WriteString(formatter.Bold(formatter.Green(s)).String())
	}

	for p, m := range diff.Modified {
		equal = false
		oldVal := old
		s := fmt.Sprintf("\t~ %s = %v => %v\n", p.String(), oldVal[pathToKey(p)], m)
		diffs.WriteString(formatter.Bold(formatter.Red(s)).String())
	}

	return diffs.String(), equal
}

func pathToKey(p *messagediff.Path) string {
	return strings.Trim(p.String(), "[\"]")
}
