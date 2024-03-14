# SmolPad

A text editor written in Go :)

## Goal(s)

In short, the goal is to write a text editor which can read/write to text files.
I very much doubt that I will add syntax highlighting.

- [x] Rope data structure (might need to work on this a bit more)

- [x] Basic Insert/Delete

  - [x] Insert characters
  - [x] Delete characters

- [x] Basic navigation

  - [x] Left/Right
  - [x] Up
  - [x] Down

- [x] Handle special characters

  - [x] Newline
  - [x] Tab

- [x] Reading text files
- [x] Writing to text files
- [x] Scrolling on large text content

  - [x] Y-axis
  - [x] X-axis

- [x] Displaying line numbers

  - [x] Highlight current line
  - [x] Relative numbers

## Dependencies

I'm using [tcell (note: v2)](https://github.com/gdamore/tcell) to manage
writing to/from the terminal.
