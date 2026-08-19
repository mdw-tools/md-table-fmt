package main

import (
	"strings"
	"unicode/utf8"
)

type alignment int

const (
	alignDefault alignment = iota
	alignLeft
	alignRight
	alignCenter
)

// Format returns the markdown input with each table rewritten so that the
// columns are vertically aligned. All other content passes through unchanged.
func Format(input string) string {
	lines := strings.Split(input, "\n")
	var output []string
	var fence *fenceInfo
	for index := 0; index < len(lines); {
		line := lines[index]
		if fence != nil {
			output = append(output, line)
			if fence.closes(line) {
				fence = nil
			}
			index++
			continue
		}
		if opened := openFence(line); opened != nil {
			fence = opened
			output = append(output, line)
			index++
			continue
		}
		if end, ok := tableExtent(lines, index); ok {
			output = append(output, formatTable(lines[index:end])...)
			index = end
			continue
		}
		output = append(output, line)
		index++
	}
	return strings.Join(output, "\n")
}

type fenceInfo struct {
	marker rune
	length int
}

// openFence reports whether the line opens a fenced code block: up to three
// spaces of indentation followed by at least three backticks or tildes.
func openFence(line string) (result *fenceInfo) {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return nil
	}
	if trimmed == "" {
		return nil
	}
	marker := rune(trimmed[0])
	if marker != '`' && marker != '~' {
		return nil
	}
	length := 0
	for _, character := range trimmed {
		if character != marker {
			break
		}
		length++
	}
	if length < 3 {
		return nil
	}
	if marker == '`' && strings.ContainsRune(trimmed[length:], '`') {
		return nil
	}
	return &fenceInfo{marker: marker, length: length}
}

// closes reports whether the line closes this fence: the same marker,
// at least as many repetitions, and nothing else but spaces.
func (this *fenceInfo) closes(line string) bool {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return false
	}
	body := strings.TrimRight(trimmed, " ")
	if utf8.RuneCountInString(body) < this.length {
		return false
	}
	for _, character := range body {
		if character != this.marker {
			return false
		}
	}
	return true
}

// tableExtent reports whether a table starts at lines[start]. When it does,
// end is the index of the first line after the table.
func tableExtent(lines []string, start int) (end int, ok bool) {
	if start+1 >= len(lines) {
		return start, false
	}
	header := lines[start]
	if !containsUnescapedPipe(header) {
		return start, false
	}
	delimiter := lines[start+1]
	if !isDelimiterRow(delimiter) {
		return start, false
	}
	if len(splitCells(header)) != len(splitCells(delimiter)) {
		return start, false
	}
	end = start + 2
	for end < len(lines) {
		line := lines[end]
		if strings.TrimSpace(line) == "" {
			break
		}
		if !containsUnescapedPipe(line) {
			break
		}
		if openFence(line) != nil {
			break
		}
		end++
	}
	return end, true
}

// isDelimiterRow reports whether the line is a table delimiter row:
// it contains an unescaped pipe and every cell is dashes with optional
// leading and trailing colons.
func isDelimiterRow(line string) bool {
	if !containsUnescapedPipe(line) {
		return false
	}
	for _, cell := range splitCells(line) {
		if !isDelimiterCell(cell) {
			return false
		}
	}
	return true
}

func isDelimiterCell(cell string) bool {
	body := strings.TrimSuffix(strings.TrimPrefix(cell, ":"), ":")
	if body == "" {
		return false
	}
	return strings.Trim(body, "-") == ""
}

func delimiterAlignment(cell string) alignment {
	leading := strings.HasPrefix(cell, ":")
	trailing := strings.HasSuffix(cell, ":")
	switch {
	case leading && trailing:
		return alignCenter
	case trailing:
		return alignRight
	case leading:
		return alignLeft
	default:
		return alignDefault
	}
}

func containsUnescapedPipe(line string) bool {
	escaped := false
	for _, character := range line {
		if escaped {
			escaped = false
			continue
		}
		if character == '\\' {
			escaped = true
			continue
		}
		if character == '|' {
			return true
		}
	}
	return false
}

// splitCells splits a table line into trimmed cell contents. It ignores
// escaped pipes and it drops the empty segments that the outer pipes create.
func splitCells(line string) (results []string) {
	trimmed := strings.TrimSpace(line)
	var segments []string
	var current strings.Builder
	escaped := false
	lastWasSplit := false
	for _, character := range trimmed {
		if escaped {
			current.WriteRune(character)
			escaped = false
			lastWasSplit = false
			continue
		}
		if character == '\\' {
			current.WriteRune(character)
			escaped = true
			lastWasSplit = false
			continue
		}
		if character == '|' {
			segments = append(segments, current.String())
			current.Reset()
			lastWasSplit = true
			continue
		}
		current.WriteRune(character)
		lastWasSplit = false
	}
	segments = append(segments, current.String())
	if len(segments) > 1 && segments[len(segments)-1] == "" && lastWasSplit {
		segments = segments[:len(segments)-1]
	}
	if len(segments) > 1 && segments[0] == "" && strings.HasPrefix(trimmed, "|") {
		segments = segments[1:]
	}
	for _, segment := range segments {
		results = append(results, strings.TrimSpace(segment))
	}
	return results
}

// formatTable rewrites the lines of one table with vertically aligned columns.
// The first line is the header, the second is the delimiter row, and the rest
// are body rows. The indentation of the header applies to every line.
func formatTable(lines []string) (results []string) {
	indent := lines[0][:len(lines[0])-len(strings.TrimLeft(lines[0], " "))]
	columns := len(splitCells(lines[1]))
	alignments := make([]alignment, columns)
	for index, cell := range splitCells(lines[1]) {
		alignments[index] = delimiterAlignment(cell)
	}
	var rows [][]string
	for index, line := range lines {
		if index == 1 {
			continue
		}
		rows = append(rows, fitToColumns(splitCells(line), columns))
	}
	widths := make([]int, columns)
	for index := range widths {
		widths[index] = 1
	}
	for _, row := range rows {
		for index, cell := range row {
			width := utf8.RuneCountInString(cell)
			if width > widths[index] {
				widths[index] = width
			}
		}
	}
	results = append(results, indent+renderRow(rows[0], widths, alignments))
	results = append(results, indent+renderDelimiter(widths, alignments))
	for _, row := range rows[1:] {
		results = append(results, indent+renderRow(row, widths, alignments))
	}
	return results
}

// fitToColumns pads a short row with empty cells and merges the extra cells
// of a long row into the last column, so no content is lost.
func fitToColumns(cells []string, columns int) (results []string) {
	for len(cells) < columns {
		cells = append(cells, "")
	}
	if len(cells) > columns {
		merged := strings.Join(cells[columns-1:], " | ")
		cells = append(cells[:columns-1], merged)
	}
	return cells
}

func renderRow(cells []string, widths []int, alignments []alignment) string {
	var builder strings.Builder
	builder.WriteString("|")
	for index, cell := range cells {
		builder.WriteString(" ")
		builder.WriteString(pad(cell, widths[index], alignments[index]))
		builder.WriteString(" |")
	}
	return builder.String()
}

func renderDelimiter(widths []int, alignments []alignment) string {
	var builder strings.Builder
	builder.WriteString("|")
	for index, width := range widths {
		switch alignments[index] {
		case alignLeft:
			builder.WriteString(":" + strings.Repeat("-", width+1))
		case alignRight:
			builder.WriteString(strings.Repeat("-", width+1) + ":")
		case alignCenter:
			builder.WriteString(":" + strings.Repeat("-", width) + ":")
		default:
			builder.WriteString(strings.Repeat("-", width+2))
		}
		builder.WriteString("|")
	}
	return builder.String()
}

func pad(cell string, width int, align alignment) string {
	gap := width - utf8.RuneCountInString(cell)
	if gap <= 0 {
		return cell
	}
	switch align {
	case alignRight:
		return strings.Repeat(" ", gap) + cell
	case alignCenter:
		left := gap / 2
		return strings.Repeat(" ", left) + cell + strings.Repeat(" ", gap-left)
	default:
		return cell + strings.Repeat(" ", gap)
	}
}
