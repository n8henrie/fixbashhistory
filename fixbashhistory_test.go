package main

import (
	"strings"
	"testing"
)

func TestSortLines(t *testing.T) {
	type test struct {
		lines, sortedLines string
	}

	tests := []test{
		{
			`#22345
I come second
#32345
I come third
#12345
I come first`,
			`#12345
I come first
#22345
I come second
#32345
I come third`,
		}}

	for _, point := range tests {
		commands := makeCommandArr(strings.Split(point.lines, "\n"))
		sortCommands(&commands)
		n := commandsToString(commands)
		if n != point.sortedLines {
			t.Error("For", point.lines, "expected", point.sortedLines, "got", n)
		}
	}
}

func TestDedupLines(t *testing.T) {
	type test struct {
		lines, dedupLines string
	}

	tests := []test{
		{
			`#12345
I should be under timestamp 22345
#22345
I should be under timestamp 22345
#32345
I should be under timestamp 42345
#42345
I should be under timestamp 42345`,
			`#22345
I should be under timestamp 22345
#42345
I should be under timestamp 42345`,
		}}

	for _, point := range tests {
		commands := makeCommandArr(strings.Split(point.lines, "\n"))
		n := commandsToString(dedupCommands(commands))
		if n != point.dedupLines {
			t.Error("For", point.lines, "expected", point.dedupLines, "got", n)
		}
	}
}

func TestPackage(t *testing.T) {
	type test struct {
		lines, sortedDedupLines string
	}
	var tests = []test{
		{
			`#22345
this that
#12345
the other
#42345
this that
#52345
the other
#32345
foo bar
#32345
baz qux`,
			`#32345
foo bar
#32345
baz qux
#42345
this that
#52345
the other`,
		},
		{
			`#22345
$('other stuff') # and a comment
#12345
/bin/echo "${stuff}"`,
			`#12345
/bin/echo "${stuff}"
#22345
$('other stuff') # and a comment`,
		},
	}

	for _, point := range tests {
		commands := makeCommandArr(strings.Split(point.lines, "\n"))
		sortCommands(&commands)
		dedup := dedupCommands(commands)
		sortCommands(&dedup)
		n := commandsToString(dedup)
		if n != point.sortedDedupLines {
			t.Error("For", point.lines, "expected", point.sortedDedupLines, "got", n)
		}
	}
}

func BenchmarkPackage(b *testing.B) {
	// run the Foo function b.N times
	// NB: do not use n or b.N in the loop
	benchData := `#22345
this that
#12345
the other
#42345
this that
#52345
the other
#32345
foo bar
#32345
baz qux`

	for n := 0; n < b.N; n++ {
		commands := makeCommandArr(strings.Split(benchData, "\n"))
		sortCommands(&commands)
		dedup := dedupCommands(commands)
		sortCommands(&dedup)
		commandsToString(dedup)
	}
}
