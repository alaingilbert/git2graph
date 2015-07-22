package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"github.com/codegangsta/cli"
	"encoding/json"
)

type InputNode struct {
	Id      string   `json:"id"`
	Parents []string `json:"parents"`
}

// Type:
// 0: |
// 1: /
// 2: \
type Path struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Type int `json:"type"`
}

type OutputNode struct {
	Id           string            `json:"id"`
	Parents      []string          `json:"parents"`
	Column       int               `json:"column"`
	ParentsPaths map[string][]Path `json:"parents_paths"`
	Idx          int               `json:"idx"`
}

func getInputNodesFromJson(inputJson string) (nodes []InputNode, err error) {
	if err = json.Unmarshal([]byte(inputJson), &nodes); err != nil {
		return
	}
	return
}

func initNodes(inputNodes []InputNode) ([]*OutputNode) {
	out := make([]*OutputNode, 0)
	for idx, node := range inputNodes {
		out = append(out, &OutputNode{node.Id, node.Parents, -1, make(map[string][]Path), idx})
	}
	return out
}

func initIndex(nodes []*OutputNode) map[string]*OutputNode {
	index := make(map[string]*OutputNode)
	for _, node := range nodes {
		index[node.Id] = node
	}
	return index
}

func setColumns(nodes []*OutputNode, index map[string]*OutputNode) {
	nextColumn := 0
	for _, node := range nodes {
		if node.Column == -1 {
			node.Column = nextColumn
			nextColumn++
		}

		for parentIdx, parentId := range node.Parents {
			parent := index[parentId]
			if parent.Column == -1 {
				if parentIdx == 0 {
					parent.Column = node.Column
				} else {
					parent.Column = nextColumn
					nextColumn++
				}
			} else {
				if parentIdx == 0 && node.Column <  parent.Column {
					parent.Column = node.Column
				}
			}
		}
	}
}

func setPaths(nodes []*OutputNode, index map[string]*OutputNode) {
	for _, node := range nodes {
		for _, parentId := range node.Parents {
			parent := index[parentId]
			node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Path{node.Column, node.Idx, 0})
			if node.Column > parent.Column {
				node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Path{node.Column, parent.Idx, 1})
			}
			node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Path{parent.Column, parent.Idx, 0})
		}
	}
}

func buildTree(inputNodes []InputNode) ([]*OutputNode, error) {
	var nodes []*OutputNode = initNodes(inputNodes)
	var index map[string]*OutputNode = initIndex(nodes)

	setColumns(nodes, index)
	setPaths(nodes, index)

	return nodes, nil
}

func BuildTreeJson(inputJson string) (tree string, err error) {
	nodes, err := getInputNodesFromJson(inputJson)
	if err != nil {
		return
	}

	out, err := buildTree(nodes)
	if err != nil {
		return
	}

	treeBytes, err := json.Marshal(&out)
	if err != nil {
		return
	}
	tree = string(treeBytes)
	return
}

func bootstrap(c *cli.Context) {
	var inputJson string
	if c.String("json") != "" {
		inputJson = c.String("json")
	} else if c.String("input") != "" {
		bytes, err := ioutil.ReadFile(c.String("input"))
		if err != nil {
			fmt.Println(err)
			return
		}
		inputJson = string(bytes)
	} else {
		cli.ShowAppHelp(c)
		return
	}

	out, err := BuildTreeJson(inputJson)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

func main() {
	var authors []cli.Author
	// Collaborators, add your name here :)
	authors = append(authors, cli.Author{"Alain Gilbert", "alain.gilbert.15@gmail.com"})

	app := cli.NewApp()
	app.Authors = authors
	app.Version = "0.0.0"
	app.Name = "git2graph"
	app.Usage = "Take a git tree, make a graph structure"
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "f, file",
			Usage: "File",
		},
		cli.StringFlag {
			Name: "j, json",
			Usage: "Json input",
		},
	}
	app.Action = bootstrap
	app.Run(os.Args)
}