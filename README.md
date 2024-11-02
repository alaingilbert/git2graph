![Logo](img/logo.png)

Open source implementation of the [sourcetree](https://www.sourcetreeapp.com/) graph algorithm.

Given a list of commits and their parents, return a structure that tells you exactly how to draw the git graph.

It takes a json:

```json
[
  {"id": "1", "parents": ["3"], "non_related_attr": "non_related_value"},
  {"id": "2", "parents": ["3"]},
  {"id": "3", "parents": []}
]
```

and returns a structure that represent a git graph:

```json
[
  {"id": "1", "parents": ["3"], "non_related_attr": "non_related_value",
   "g": [0,0,"#5aa1be",[["#5aa1be",[[0,0,0],[0,2,0]]]]]},
  {"id": "2", "parents": ["3"],
   "g": [1,1,"#c065b8",[["#c065b8",[[1,1,0],[1,2,1],[0,2,0]]]]]},
  {"id": "3", "parents": [],
   "g": [2,0,"#5aa1be",[]]}
]
```

This structure can be directly rendered with D3.js, [you can try it out here.](http://alaingilbert.github.io/git2graph/)

![Logo](img/img1.png)

### Other examples

![Logo](img/img2.png)
![Logo](img/img5.png)
![Logo](img/img3.png)
![Logo](img/img4.png)

## How to use

### Inline

`git2graph -j '[{"id": 1, "parents": ["2"]}, ...]'`

### File

`git2graph -f path/to/file.json`

### Repository

`git2graph -r` (You must be in the repository directory)

### In code

```go
package main

import (
  "fmt"
  "git2graph"
)

func main() {
  in := []map[string]any{
    {"id": "1", "parents": []string{"3"}},
    {"id": "2", "parents": []string{"3"}},
    {"id": "3", "parents": []string{}},
  }

  // Use this to render a whole graph in a single SVG
  out, err := git2graph.Get(in)
  fmt.Println(out, err)

  // Use this if you want to render each rows individually (for html table)
  out, err = git2graph.GetRows(in)
  fmt.Println(out, err)
}
```

## See it in action

```
renderer/index.html
```

Use D3.js to render the graph represented by the output of Git2Graph.

## How to run

```
go run main.go -j '...'
```

Or

```
go install
git2graph -j '...'
```

## How to test
```
go test ./...
```

## How to contribute

- Fork the repo
- Create a new branch
- Make your changes
- Create new tests
- Append your name/email in main.go (contributors list)
- Make a pull request :)
