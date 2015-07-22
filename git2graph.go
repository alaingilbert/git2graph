package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

type InputNode struct {
	Id      string   `json:"id"`
	Parents []string `json:"parents"`
}

// Type:
// 0: |
// 1: ┘
// 2: ┐
// 3: ┌
type Point struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Type int `json:"type"`
}

type Path struct {
	Id   string  `json:"id"`
	Path []Point `json:"path"`
}

type OutputNode struct {
	Id                string             `json:"id"`
	Parents           []string           `json:"parents"`
	Column            int                `json:"column"`
	ParentsPaths      map[string][]Point `json:"-"`
	FinalParentsPaths []Path             `json:"parents_paths"`
	Idx               int                `json:"idx"`
	Children          []string           `json:"-"`
}

func serializeOutput(out []*OutputNode) ([]byte, error) {
	for _, node := range out {
		for parentId, path := range node.ParentsPaths {
			node.FinalParentsPaths = append(node.FinalParentsPaths, Path{parentId, path})
		}
	}
	treeBytes, err := json.Marshal(&out)
	return treeBytes, err
}

func getInputNodesFromJson(inputJson string) (nodes []InputNode, err error) {
	if err = json.Unmarshal([]byte(inputJson), &nodes); err != nil {
		return
	}
	return
}

func initNodes(inputNodes []InputNode) []*OutputNode {
	out := make([]*OutputNode, 0)
	for idx, node := range inputNodes {
		newNode := OutputNode{}
		newNode.Id = node.Id
		newNode.Parents = node.Parents
		newNode.Column = -1
		newNode.ParentsPaths = make(map[string][]Point)
		newNode.FinalParentsPaths = make([]Path, 0)
		newNode.Idx = idx
		newNode.Children = make([]string, 0)
		out = append(out, &newNode)
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

func initChildren(nodes []*OutputNode, index map[string]*OutputNode) {
	for _, node := range nodes {
		for _, parentId := range node.Parents {
			index[parentId].Children = append(index[parentId].Children, node.Id)
		}
	}
}

func setColumns(nodes []*OutputNode, index map[string]*OutputNode) {
	nextColumn := 0
	for _, node := range nodes {
		if node.Column == -1 {
			node.Column = nextColumn
			nextColumn++
		}

		for _, childId := range node.Children {
			child := index[childId]
			if child.Column > node.Column {
				nextColumn--
			}
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
				if parentIdx == 0 && node.Column < parent.Column {
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
			node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Point{node.Column, node.Idx, 0})
			if node.Column > parent.Column {
				node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Point{node.Column, parent.Idx, 1})
			} else if node.Column < parent.Column {
				node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Point{parent.Column, node.Idx, 2})
			}
			node.ParentsPaths[parent.Id] = append(node.ParentsPaths[parent.Id], Point{parent.Column, parent.Idx, 0})
		}
	}
}

func buildTree(inputNodes []InputNode) ([]*OutputNode, error) {
	var nodes []*OutputNode = initNodes(inputNodes)
	var index map[string]*OutputNode = initIndex(nodes)

	initChildren(nodes, index)
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

	treeBytes, err := serializeOutput(out)
	if err != nil {
		return
	}
	tree = string(treeBytes)
	return
}

func bootstrap(c *cli.Context) {
	var inputJson string
	jsonFlag := c.String("json")
	fileFlag := c.String("file")
	if jsonFlag != "" {
		inputJson = jsonFlag
	} else if fileFlag != "" {
		bytes, err := ioutil.ReadFile(fileFlag)
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f, file",
			Usage: "File",
		},
		cli.StringFlag{
			Name:  "j, json",
			Usage: "Json input",
		},
	}
	app.Action = bootstrap
	app.Run(os.Args)
}
