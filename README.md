# NutCode

A text editor written in Go :)

## Goal(s)

In short, the goal is to write a text editor which can read/write to text files.
I very much doubt that I will add syntax highlighting.

- [ ] Rope data structure
- [ ] Reading text files
- [ ] Writing to text files
- [ ] Displaying line numbers
  - [ ] Highlight current line
  - [ ] Relative numbers
- [ ] Vim bindings (subset)
  - [ ] Normal mode
  - [ ] Insert mode
  - [ ] Visual mode
  - [ ] Command mode (maybe)
- [ ] Cursor should remember column

## Dependencies

I'm using [tcell (note: v2)](https://github.com/gdamore/tcell) to manage writing to/from the terminal.
