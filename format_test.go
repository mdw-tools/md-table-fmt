package main

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, actual, expected string) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nexpected: %q\nactual:   %q", expected, actual)
	}
}

func TestFormatPassesPlainMarkdownThrough(t *testing.T) {
	input := "# Title\n\nSome *text* with | a stray pipe.\n\n- item 1\n- item 2\n"
	assertEqual(t, Format(input), input)
}

func TestFormatPreservesMissingTrailingNewline(t *testing.T) {
	input := "no newline at end"
	assertEqual(t, Format(input), input)
}

func TestFormatAlignsBasicTable(t *testing.T) {
	input := strings.Join([]string{
		"| Column A | Col B | Column C Is Very Verbose |",
		"|---|---|---|",
		"| `Hello` | 86-88 | Some very interesting text |",
		"| `Bye` | 146-147 | Lorem Ipsum |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| Column A | Col B   | Column C Is Very Verbose   |",
		"|----------|---------|----------------------------|",
		"| `Hello`  | 86-88   | Some very interesting text |",
		"| `Bye`    | 146-147 | Lorem Ipsum                |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatHonorsAlignmentMarkers(t *testing.T) {
	input := strings.Join([]string{
		"| Left | Center | Right |",
		"|:--|:-:|--:|",
		"| a | b | c |",
		"| longer | wider | rightmost |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| Left   | Center |     Right |",
		"|:-------|:------:|----------:|",
		"| a      |   b    |         c |",
		"| longer | wider  | rightmost |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatAccountsForEscapedPipes(t *testing.T) {
	input := strings.Join([]string{
		"| Expression | Meaning |",
		"|---|---|",
		"| `a \\| b` | a or b |",
		"| x | y |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| Expression | Meaning |",
		"|------------|---------|",
		"| `a \\| b`   | a or b  |",
		"| x          | y       |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatIgnoresTablesInBacktickFence(t *testing.T) {
	input := strings.Join([]string{
		"before",
		"```markdown",
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"```",
		"after",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatIgnoresTablesInTildeFence(t *testing.T) {
	input := strings.Join([]string{
		"~~~",
		"| a | b |",
		"|---|---|",
		"~~~",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatResumesAfterCodeFenceCloses(t *testing.T) {
	input := strings.Join([]string{
		"```",
		"| raw | table |",
		"|---|---|",
		"```",
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"```",
		"| raw | table |",
		"|---|---|",
		"```",
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatIgnoresShorterClosingFenceCandidate(t *testing.T) {
	input := strings.Join([]string{
		"````",
		"```",
		"| a | b |",
		"|---|---|",
		"````",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatPadsShortRowsWithEmptyCells(t *testing.T) {
	input := strings.Join([]string{
		"| a | b | c |",
		"|---|---|---|",
		"| 1 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b | c |",
		"|---|---|---|",
		"| 1 |   |   |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatMergesExtraCellsIntoLastColumn(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 2 | 3 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b     |",
		"|---|-------|",
		"| 1 | 2 | 3 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatRequiresMatchingDelimiterCellCount(t *testing.T) {
	input := strings.Join([]string{
		"| a | b | c |",
		"|---|---|",
		"| 1 | 2 | 3 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatStopsTableAtBlankLine(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"",
		"| not | a table row anymore",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
		"| not | a table row anymore",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatStopsTableAtLineWithoutPipe(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"plain paragraph text",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"plain paragraph text",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatPreservesTableIndentation(t *testing.T) {
	input := strings.Join([]string{
		"  | a | b |",
		"  |---|---|",
		"  | 1 | 22 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"  | a | b  |",
		"  |---|----|",
		"  | 1 | 22 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatHandlesTableWithoutOuterPipes(t *testing.T) {
	input := strings.Join([]string{
		"a | b",
		"--- | ---",
		"1 | 22",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatMeasuresUnicodeCellsByRuneCount(t *testing.T) {
	input := strings.Join([]string{
		"| Name | Range |",
		"|---|---|",
		"| héllo | 86–88 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| Name  | Range |",
		"|-------|-------|",
		"| héllo | 86–88 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatHandlesMultipleTables(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"",
		"text between",
		"",
		"| x | y |",
		"|---|---|",
		"| 333 | 4 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
		"text between",
		"",
		"| x   | y |",
		"|-----|---|",
		"| 333 | 4 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatDoesNotTreatThematicBreakAsDelimiterRow(t *testing.T) {
	input := strings.Join([]string{
		"| just | one line with pipes |",
		"---",
		"more text",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatHandlesTableAtEndOfFileWithoutNewline(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatHandlesEmptyCells(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"|  | 22 |",
		"| 1 |  |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"|   | 22 |",
		"| 1 |    |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatWidensNarrowColumnsToFitDelimiter(t *testing.T) {
	input := strings.Join([]string{
		"| a |",
		"|:-:|",
		"| b |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a |",
		"|:-:|",
		"| b |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatTreatsOverIndentedFenceMarkerAsText(t *testing.T) {
	input := strings.Join([]string{
		"    ```",
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"    ```",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"    ```",
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"    ```",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatTreatsShortBacktickRunAsText(t *testing.T) {
	input := strings.Join([]string{
		"`` not a fence ``",
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"`` not a fence ``",
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatRejectsBacktickFenceWithBacktickInInfoString(t *testing.T) {
	input := strings.Join([]string{
		"``` info `string`",
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"``` info `string`",
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatIgnoresOverIndentedClosingFenceCandidate(t *testing.T) {
	input := strings.Join([]string{
		"```",
		"    ```",
		"| a | b |",
		"|---|---|",
		"```",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatStopsTableAtFenceLineContainingPipe(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"``` info | pipe",
		"| raw | table |",
		"```",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"``` info | pipe",
		"| raw | table |",
		"```",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestFormatRejectsDelimiterRowWithEmptyCell(t *testing.T) {
	input := strings.Join([]string{
		"| a |",
		"|:|",
		"| b |",
		"",
	}, "\n")
	assertEqual(t, Format(input), input)
}

func TestFormatStopsTableAtLineWithOnlyEscapedPipes(t *testing.T) {
	input := strings.Join([]string{
		"| a | b |",
		"|---|---|",
		"| 1 | 22 |",
		"\\| escaped \\|",
		"",
	}, "\n")
	expected := strings.Join([]string{
		"| a | b  |",
		"|---|----|",
		"| 1 | 22 |",
		"\\| escaped \\|",
		"",
	}, "\n")
	assertEqual(t, Format(input), expected)
}

func TestSplitCells(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		expected []string
	}{
		{name: "outer pipes", line: "| a | b |", expected: []string{"a", "b"}},
		{name: "no outer pipes", line: "a | b", expected: []string{"a", "b"}},
		{name: "escaped pipe", line: "| a \\| b | c |", expected: []string{"a \\| b", "c"}},
		{name: "empty middle cell", line: "| a |  | c |", expected: []string{"a", "", "c"}},
		{name: "single cell", line: "| only |", expected: []string{"only"}},
		{name: "trailing escaped pipe", line: "| a \\| |", expected: []string{"a \\|"}},
		{name: "escaped backslash then pipe", line: "| a \\\\ | b |", expected: []string{"a \\\\", "b"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			actual := splitCells(test.line)
			if len(actual) != len(test.expected) {
				t.Fatalf("expected %d cells %q, got %d cells %q",
					len(test.expected), test.expected, len(actual), actual)
			}
			for c := range actual {
				if actual[c] != test.expected[c] {
					t.Errorf("cell %d: expected %q, got %q", c, test.expected[c], actual[c])
				}
			}
		})
	}
}
