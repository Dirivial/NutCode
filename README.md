# NutCode

A text editor written in Go :)

## Goal(s)

In short, the goal is to write a text editor which can read/write to text files.
I very much doubt that I will add syntax highlighting.

- [x] Rope data structure (might need to work on this a bit more)

- [ ] Basic Insert/Delete

  - [x] Insert characters
  - [ ] Delete characters

- [ ] Basic navigation

  - [x] Left/Right
  - [ ] Up
  - [x] Down

- [ ] Display some special characters

  - [x] Newline
  - [ ] Tab

- [ ] Reading text files
- [ ] Writing to text files

- [ ] Displaying line numbers

  - [ ] Highlight current line
  - [ ] Relative numbers

- [ ] Cursor should remember column

- [ ] Vim bindings (subset)

  - [ ] Normal mode
  - [ ] Insert mode
  - [ ] Visual mode
  - [ ] Command mode (maybe)

- [ ] Undo
- [ ] Redo

## Dependencies

I'm using [tcell (note: v2)](https://github.com/gdamore/tcell) to manage
writing to/from the terminal.
