---
name: md-table-fmt
description: Format the tables in markdown files with the md-table-fmt tool. Use when the user asks to format, align, or fix markdown tables, or when a markdown file has misaligned table columns.
---

# Format markdown tables

Use the `md-table-fmt` tool to align the columns of every table in a markdown
file. The tool changes only tables. All other content passes through
unchanged.

## Step 1: Locate the tool

Check for an installed binary:

```
command -v md-table-fmt
```

If the binary is not on the PATH, build it from this repository:

```
go build -o /tmp/md-table-fmt github.com/mdw-tools/md-table-fmt
```

Outside this repository, install it instead:

```
go install github.com/mdw-tools/md-table-fmt@latest
```

Use the path of the binary you found or built in the steps below.

## Step 2: Format each file

The tool is a Unix-style filter. It reads markdown from standard input and
writes formatted markdown to standard output. It takes no arguments and no
flags.

To format a file in place, write to a temporary file first, then replace the
original:

```
md-table-fmt < notes.md > notes.md.tmp && mv notes.md.tmp notes.md
```

**WARNING:** Never redirect a file into itself (`md-table-fmt < notes.md >
notes.md`). The shell truncates the file before the tool reads it, and the
content is lost.

To preview the result without a change to the file:

```
md-table-fmt < notes.md | diff notes.md -
```

An empty diff means the file is already formatted.

## Step 3: Verify

The output is idempotent. A second pass makes no further changes. Confirm a
formatted file with:

```
md-table-fmt < notes.md | diff notes.md - && echo "formatted"
```

Then read the changed regions of the file to confirm the tables look correct.

## What the tool does

- It pads each cell so the pipe characters (`|`) align vertically.
- It honors the alignment markers (`:--`, `:-:`, `--:`) in the delimiter row.
- An escaped pipe (`\|`) does not split a cell.
- Content inside fenced code blocks (``` or `~~~`) passes through unchanged,
  even when it looks like a table.
- A short row gains empty cells. The extra cells of a long row merge into the
  last column. No content is lost.

## Limits

- The tool does not repair broken markdown. A table must already parse as a
  table (a header row, then a delimiter row of dashes). Fix structural
  problems by hand first, then run the tool.
- The tool formats whole files only. To format one table, pipe just that
  table through the tool and paste the result back.
