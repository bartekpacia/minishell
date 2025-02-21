# minishell

```console
$ go build -o msh .
```

Try these examples, they should work :)

```
$ ls
```

```
$ ls | head -n 1
```

```
$ cat | rev
```

### Features

Redirection operators:

- `>` redirect stdout (overwrite). Same as `1>`
- `>>` redirect stdout (append). Same as `2>`
- `<` redirect stdin from a file

Pipe operators:

- `|` - redirect stdout to another command's stdin
