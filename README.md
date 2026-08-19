# md-table-fmt

A Unix-style filter that formats markdown tables with vertically aligned
columns. All other content passes through unchanged.

## Install

```
go install github.com/mdw-tools/md-table-fmt@latest
```

## Usage

The tool reads from standard input and writes to standard output.

```
md-table-fmt < notes.md
cat notes.md | md-table-fmt
```

## Behavior

The tool turns this:

```
| Column A | Col B | Column C Is Very Verbose |
|---|---|---|
| `Hello` | 86-88 | Some very interesting text |
| `Bye` | 146-147 | Lorem Ipsum |
```

into this:

```
| Column A | Col B   | Column C Is Very Verbose   |
|----------|---------|----------------------------|
| `Hello`  | 86-88   | Some very interesting text |
| `Bye`    | 146-147 | Lorem Ipsum                |
```

Details:

- The tool honors the alignment markers (`:--`, `:-:`, `--:`) in the
  delimiter row. It pads each cell to match.
- An escaped pipe (`\|`) does not split a cell.
- Content inside fenced code blocks (``` or `~~~`) passes through unchanged,
  even when it looks like a table.
- A short row gains empty cells. The extra cells of a long row merge into
  the last column. No content is lost.
- The output is idempotent: a second pass makes no further changes.

## Development

```
make test
```
