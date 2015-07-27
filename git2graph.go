package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

var colors []string

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
	Id    string  `json:"id"`
	Path  []Point `json:"path"`
	Color string  `json:"color"`
}

type OutputNode struct {
	Id                string          `json:"id"`
	Parents           []string        `json:"parents"`
	Column            int             `json:"column"`
	ParentsPaths      map[string]Path `json:"-"`
	FinalParentsPaths []Path          `json:"parents_paths"`
	Idx               int             `json:"idx"`
	Children          []string        `json:"-"`
	Color             string          `json:"color"`
	FirstInRow        bool            `json:"-"`
}

func (node *OutputNode) Append(parentId string, point Point) {
	tmp := node.ParentsPaths[parentId]
	tmp.Path = append(tmp.Path, point)
	node.ParentsPaths[parentId] = tmp
}

func (node *OutputNode) Remove(parentId string, idx int) {
	tmp := node.ParentsPaths[parentId]
	tmp.Path = append(tmp.Path[:idx], tmp.Path[idx+1:]...)
	node.ParentsPaths[parentId] = tmp
}

func (node *OutputNode) Insert(parentId string, idx int, point Point) {
	tmp := node.ParentsPaths[parentId]
	tmp.Path = append(tmp.Path, Point{})
	copy(tmp.Path[idx+1:], tmp.Path[idx:])
	tmp.Path[idx] = point
	node.ParentsPaths[parentId] = tmp
}

func (node *OutputNode) SetPathColor(parentId, color string) {
	tmp := node.ParentsPaths[parentId]
	tmp.Color = color
	node.ParentsPaths[parentId] = tmp
}

func serializeOutput(out []*OutputNode) ([]byte, error) {
	for _, node := range out {
		for parentId, path := range node.ParentsPaths {
			node.FinalParentsPaths = append(node.FinalParentsPaths, Path{parentId, path.Path, path.Color})
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
		newNode.ParentsPaths = make(map[string]Path)
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
			node.Color, colors = colors[0], colors[1:]
			nextColumn++
		}

		for _, childId := range node.Children {
			child := index[childId]
			isType3 := child.ParentsPaths[node.Id].Path[len(child.ParentsPaths[node.Id].Path)-2].Type == 3
			if node.Column < child.Column && !isType3 {
				nextColumn--

				if child.Parents[0] != node.Id || len(child.Parents) <= 1 {
					if !child.FirstInRow {
						child.SetPathColor(node.Id, child.Color)
					}
					colors = append(colors[:1], append([]string{child.Color}, colors[1:]...)...)

					// Insert before the last element
					pos := len(child.ParentsPaths[node.Id].Path) - 1
					point := Point{child.ParentsPaths[node.Id].Path[pos-1].X, node.Idx, 1}
					child.Insert(node.Id, pos, point)

					for followingNodeIdx, followingNode := range nodes {
						if followingNodeIdx > node.Idx {
							if followingNode.Column > child.Column {
								if followingNode.Column > child.ParentsPaths[node.Id].Path[len(child.ParentsPaths[node.Id].Path)-2].X {

								for _, followingNodeChildId := range followingNode.Children {
									followingNodeChild := index[followingNodeChildId]

									idxRemove := len(followingNodeChild.ParentsPaths[followingNode.Id].Path) - 1
									followingNodeChild.Remove(followingNode.Id, idxRemove)

									pos := len(followingNodeChild.ParentsPaths[followingNode.Id].Path) - 1
									followingNodeChild.Append(followingNode.Id, Point{followingNodeChild.ParentsPaths[followingNode.Id].Path[pos].X, node.Idx, 1})
									followingNodeChild.Append(followingNode.Id, Point{followingNode.Column - 1, node.Idx, 0})
									followingNodeChild.Append(followingNode.Id, Point{followingNode.Column - 1, followingNode.Idx, 0})
								}

								followingNode.Column--
								}
							}
						}
					}
				}
			}
		}

		for parentIdx, parentId := range node.Parents {
			parent := index[parentId]

			node.Append(parent.Id, Point{node.Column, node.Idx, 0})

			if parent.Column == -1 {
				if parentIdx == 0 || (parentIdx == 1 && index[node.Parents[0]].Column < node.Column) {
					parent.Column = node.Column
					parent.Color = node.Color
					node.SetPathColor(parent.Id, parent.Color)
				} else {
					parent.Column = nextColumn
					parent.Color, colors = colors[0], colors[1:]
					node.Append(parent.Id, Point{parent.Column, node.Idx, 2})
					node.SetPathColor(parent.Id, parent.Color)
					node.FirstInRow = true
					nextColumn++
				}
			} else {
				if node.Column < parent.Column && parentIdx == 0 {
					for _, childId := range parent.Children {
						child := index[childId]
						idxRemove := len(child.ParentsPaths[parent.Id].Path) - 1
						if idxRemove > 0 {
							if child.ParentsPaths[parent.Id].Path[idxRemove].Type != 2 {
								child.Remove(parent.Id, idxRemove)
							}
							child.Append(parent.Id, Point{node.Column, parent.Idx, 0})
						}
					}
					parent.Column = node.Column
					parent.Color = node.Color
					node.SetPathColor(parent.Id, node.Color)
				} else if node.Column < parent.Column && parentIdx > 0 {
					node.Append(parent.Id, Point{parent.Column, node.Idx, 2})
					node.SetPathColor(parent.Id, parent.Color)
				} else if node.Column > parent.Column {
					if len(node.Parents) > 1 {
						node.Append(parent.Id, Point{parent.Column, node.Idx, 3})
						node.SetPathColor(parent.Id, parent.Color)
					}
				}
			}

			node.Append(parent.Id, Point{parent.Column, parent.Idx, 0})

		}
	}
}

func buildTree(inputNodes []InputNode, myColors []string) ([]*OutputNode, error) {
	colors = myColors
	var nodes []*OutputNode = initNodes(inputNodes)
	var index map[string]*OutputNode = initIndex(nodes)

	initChildren(nodes, index)
	setColumns(nodes, index)

	return nodes, nil
}

func BuildTreeJson(inputJson string, myColors []string) (tree string, err error) {
	nodes, err := getInputNodesFromJson(inputJson)
	if err != nil {
		return
	}

	out, err := buildTree(nodes, myColors)
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

	myColors := []string{"#5aa1be", "#c065b8", "#c0ab5f", "#59bc95", "#c0615b", "#7a63be", "#73bb5e", "#6ee585", "#7088e8", "#eb77a3"}

	out, err := BuildTreeJson(inputJson, myColors)
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
