# minishell

```
$ ls
README.md
go.mod
main.go
```

```
$ ls | head -n 1
```

### Features

Redirection operators:

- `>` redirect stdout (overwrite). Same as `1>`
- `>>` redirect stdout (append). Same as `2>`
- `<` redirect stdin from a file

Pipe operators:

- `|` - redirect stdout to another command's stdin

## Questions

Q: Why does `ls` print in a single line by default, but when redirected with
`ls | head -n 1` it prints only a single "item"?

A: It uses `isatty(3)` to detect whether it's printing to a terminal or not.
