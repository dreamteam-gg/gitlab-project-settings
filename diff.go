package main

import (
	"fmt"
	"strings"

	"gopkg.in/d4l3k/messagediff.v1"
)

type changes struct {
	Added    []string
	Modified []string
	Removed  []string
}

// computeDiff returns updated and added diff, removed diff, removed list, bool for change
func computeDiff(old, new map[string]interface{}) (string, string, changes, bool) {
	var diffs strings.Builder
	var removedDiffs strings.Builder
	var changedItems changes

	equal := true

	diff, _ := messagediff.DeepDiff(old, new)

	for p, a := range diff.Added {
		equal = false
		s := fmt.Sprintf("\t+ %v = %v\n", p.String(), a)
		diffs.WriteString(formatter.Bold(formatter.Green(s)).String())
		changedItems.Added = append(changedItems.Added, pathToKey(p))
	}

	for p, m := range diff.Modified {
		oldVal := old
		// for comparing lists of maps which messagediff always sets as modified
		if fmt.Sprint(oldVal[pathToKey(p)]) == fmt.Sprint(m) {
			continue
		}
		equal = false
		s := fmt.Sprintf("\t~ %s = %v => %v\n", p.String(), oldVal[pathToKey(p)], m)
		diffs.WriteString(formatter.Bold(formatter.Brown(s)).String())
		changedItems.Modified = append(changedItems.Modified, pathToKey(p))
	}

	for p, r := range diff.Removed {
		s := fmt.Sprintf("\t- %v = %v\n", p.String(), r)
		removedDiffs.WriteString(formatter.Bold(formatter.Red(s)).String())
		changedItems.Removed = append(changedItems.Removed, pathToKey(p))
	}

	return diffs.String(), removedDiffs.String(), changedItems, equal
}

func pathToKey(p *messagediff.Path) string {
	return strings.Trim(p.String(), "[\"]")
}
