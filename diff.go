package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"gopkg.in/d4l3k/messagediff.v1"
)

func computeDiff(old, new map[string]interface{}) (string, bool) {
	var diffs strings.Builder
	equal := true

	diff, _ := messagediff.DeepDiff(old, new)

	for p, a := range diff.Added {
		equal = false
		s := fmt.Sprintf("\t+ %s = %s\n", p.String(), a)
		diffs.WriteString(formatter.Bold(formatter.Green(s)).String())
	}

	for p, m := range diff.Modified {
		equal = false
		oldVal := old
		s := fmt.Sprintf("\t~ %s = %s => %s\n", p.String(), oldVal[pathToKey(p)], m)
		diffs.WriteString(formatter.Bold(formatter.Red(s)).String())
	}

	return diffs.String(), equal
}

func pathToKey(p *messagediff.Path) string {
	return strings.Trim(p.String(), "[\"]")
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func newFormatter(noColor bool) aurora.Aurora {
	var formatter aurora.Aurora
	if !isTerminal() {
		formatter = aurora.NewAurora(false)
	} else {
		formatter = aurora.NewAurora(true)
	}
	return formatter
}

func isTerminal() bool {
	fd := os.Stdout.Fd()
	switch {
	case isatty.IsTerminal(fd):
		return true
	case isatty.IsCygwinTerminal(fd):
		return true
	default:
		return false
	}
}
