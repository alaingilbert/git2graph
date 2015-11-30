package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var colors []string
var debugMode bool = false

var defaultColors []string = []string{"#5aa1be", "#c065b8", "#c0ab5f", "#59bc95", "#7a63be", "#c0615b", "#73bb5e", "#6ee585", "#7088e8", "#eb77a3"}

// Types
const (
	PIPE       = iota // 0: |
	MERGE_BACK = iota // 1: ┘
	FORK       = iota // 2: ┐
	MERGE_TO   = iota // 3: ┌
)

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
	Id                string                 `json:"id"`
	Parents           []string               `json:"parents"`
	Column            int                    `json:"column"`
	ParentsPaths      map[string]Path        `json:"-"`
	FinalParentsPaths []Path                 `json:"parents_paths"`
	Idx               int                    `json:"idx"`
	Children          []string               `json:"-"`
	Color             string                 `json:"color"`
	FirstInRow        bool                   `json:"-"`
	Debug             []string               `json:"debug,omitempty"`
	NbMoveDown        int                    `json:"-"`
	InitialNode       map[string]interface{} `json:"initial_node"`
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

func (node *OutputNode) ColumnDefined() bool {
	return node.Column != -1
}

func (node *OutputNode) HasBiggerParentDefined(index map[string]*OutputNode) bool {
	found := false
	for _, pId := range node.Parents {
		p := index[pId]
		if p.Column > node.Column {
			found = true
			break
		}
	}
	return found
}

func (node *OutputNode) SetPathColor(parentId, color string) {
	tmp := node.ParentsPaths[parentId]
	tmp.Color = color
	node.ParentsPaths[parentId] = tmp
}

func (node *OutputNode) GetPathPoint(parentId string, idx int) Point {
	if idx < 0 {
		return node.ParentsPaths[parentId].Path[len(node.ParentsPaths[parentId].Path)+idx]
	} else {
		return node.ParentsPaths[parentId].Path[idx]
	}
}

func (node *OutputNode) GetPathHeightAtIdx(parentId string, lookupIdx int) (height int) {
	height = -1
	firstPoint := node.GetPathPoint(parentId, 0)
	lastPoint := node.GetPathPoint(parentId, -1)
	if lookupIdx < firstPoint.Y || lookupIdx > lastPoint.Y {
		return
	}
	for _, point := range node.ParentsPaths[parentId].Path {
		if point.Y <= lookupIdx {
			height = point.X
		}
	}
	return
}

func (node *OutputNode) PathLength(parentId string) int {
	return len(node.ParentsPaths[parentId].Path)
}

func serializeOutput(out []map[string]interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}

func getInputNodesFromJson(inputJson string) (nodes []map[string]interface{}, err error) {
	dec := json.NewDecoder(strings.NewReader(inputJson))
	err = dec.Decode(&nodes)
	if err != nil {
		return
	}
	for _, node := range nodes {
		parents := make([]string, 0)
		for _, parent := range node["parents"].([]interface{}) {
			parents = append(parents, parent.(string))
		}
		node["parents"] = parents
	}
	return
}

func initNodes(inputNodes []map[string]interface{}) []*OutputNode {
	out := make([]*OutputNode, 0)
	for idx, node := range inputNodes {
		newNode := OutputNode{}
		newNode.InitialNode = node
		newNode.Id = node["id"].(string)
		newNode.Parents = node["parents"].([]string)
		newNode.Column = -1
		newNode.ParentsPaths = make(map[string]Path)
		newNode.FinalParentsPaths = make([]Path, 0)
		newNode.Idx = idx
		newNode.Children = make([]string, 0)
		newNode.Debug = make([]string, 0)
		newNode.NbMoveDown = 0
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
		// Set column if not defined
		if !node.ColumnDefined() {
			node.Column = nextColumn
			if debugMode {
				node.Debug = append(node.Debug, fmt.Sprintf("Column set to %d", nextColumn))
			}
			node.Color, colors = colors[0], colors[1:]
			nextColumn++
		}

		// May not be useful
		nbNodesMergingBack := 0
		for _, childId := range node.Children {
			child := index[childId]
			isNodeMerging := child.GetPathPoint(node.Id, -2).Type == MERGE_TO
			if (node.Column+node.NbMoveDown) < child.Column && !isNodeMerging {
				nbNodesMergingBack++
			}
		}

		// Each children
		for _, childId := range node.Children {
			child := index[childId]
			isNodeMerging := child.GetPathPoint(node.Id, -2).Type == MERGE_TO
			if (node.Column+node.NbMoveDown) < child.Column && !isNodeMerging {
				nextColumn--

				if child.Parents[0] != node.Id || len(child.Parents) <= 1 {
					if !child.FirstInRow {
						child.SetPathColor(node.Id, child.Color)
					}
					colors = append(colors[:1], append([]string{child.Color}, colors[1:]...)...)

					// Insert before the last element
					pos := child.PathLength(node.Id) - 1
					point := Point{child.GetPathPoint(node.Id, -2).X, node.Idx, MERGE_BACK}
					child.Insert(node.Id, pos, point)

					// Nodes that are following the current node
					for followingNodeIdx, followingNode := range nodes {
						if followingNodeIdx > node.Idx {
							// Following nodes that have a child before the current node
							for _, followingNodeChildId := range followingNode.Children {
								followingNodeChild := index[followingNodeChildId]
								if followingNodeChild.Idx < node.Idx {
									// Following node child has a path that is higher than the current path being merged
									if followingNodeChild.GetPathPoint(followingNode.Id, -2).X > child.GetPathPoint(node.Id, -2).X {
										idxRemove := followingNodeChild.PathLength(followingNode.Id) - 1
										if idxRemove < 0 {
											continue
										}
										if followingNodeChild.PathLength(followingNode.Id) > idxRemove &&
											followingNodeChild.GetPathPoint(followingNode.Id, idxRemove).Y == followingNodeChild.GetPathPoint(followingNode.Id, idxRemove-1).Y {
											followingNodeChild.Remove(followingNode.Id, idxRemove-1)
											idxRemove -= 1
										}
										tmp := followingNodeChild.GetPathPoint(followingNode.Id, idxRemove-1).X
										followingNodeChild.Remove(followingNode.Id, idxRemove)
										followingNodeChild.Append(followingNode.Id, Point{tmp, node.Idx, MERGE_BACK})
										followingNodeChild.Append(followingNode.Id, Point{tmp - 1 - (nbNodesMergingBack - 1), node.Idx, PIPE})
										if followingNode.Column > child.GetPathPoint(node.Id, -2).X {
											followingNodeChild.Append(followingNode.Id, Point{followingNode.Column - 1, followingNode.Idx, PIPE})
										} else {
											followingNodeChild.Append(followingNode.Id, Point{tmp - 1 - (nbNodesMergingBack - 1), followingNode.Idx, MERGE_BACK})
											followingNodeChild.Append(followingNode.Id, Point{followingNode.Column, followingNode.Idx, PIPE})
										}
									}
								}
							}
							if followingNode.Column > child.GetPathPoint(node.Id, -2).X {
								followingNode.Column--
								followingNode.NbMoveDown++
								if debugMode {
									followingNode.Debug = append(followingNode.Debug, fmt.Sprintf("Node moved down, %s -> %s", child.Id, node.Id))
								}
							}
						}
					}
				}
			}
		}

		for parentIdx, parentId := range node.Parents {
			parent := index[parentId]

			node.Append(parent.Id, Point{node.Column, node.Idx, PIPE})

			if !parent.ColumnDefined() {
				if parentIdx == 0 || (parentIdx == 1 && index[node.Parents[0]].Column < node.Column) {
					parent.Column = node.Column
					if debugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("Column set to %d", node.Column))
					}
					parent.Color = node.Color
					node.SetPathColor(parent.Id, parent.Color)
				} else {
					parent.Column = nextColumn
					if debugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("Column set to %d", nextColumn))
					}
					parent.Color, colors = colors[0], colors[1:]
					node.Append(parent.Id, Point{parent.Column, node.Idx, FORK})
					node.SetPathColor(parent.Id, parent.Color)
					node.FirstInRow = true
					nextColumn++
				}
			} else if parent.ColumnDefined() {
				if node.Column < parent.Column && parentIdx == 0 {
					for _, childId := range parent.Children {
						child := index[childId]
						idxRemove := child.PathLength(parent.Id) - 1
						if idxRemove > 0 {
							if child.GetPathPoint(parent.Id, idxRemove).Type != FORK {
								child.Remove(parent.Id, idxRemove)
							}
							pos := child.PathLength(parent.Id) - 1
							child.Append(parent.Id, Point{child.GetPathPoint(parent.Id, pos).X, parent.Idx, MERGE_BACK})
							child.Append(parent.Id, Point{node.Column, parent.Idx, PIPE})
						}
					}
					parent.Column = node.Column
					if debugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("Column reset to %d", node.Column))
					}
					parent.Color = node.Color
					node.SetPathColor(parent.Id, node.Color)
				} else if node.Column < parent.Column && parentIdx > 0 {
					node.Append(parent.Id, Point{parent.Column, node.Idx, FORK})
					node.SetPathColor(parent.Id, parent.Color)
				} else if node.Column > parent.Column {
					if len(node.Parents) > 1 {
						if node.HasBiggerParentDefined(index) {
							node.Append(parent.Id, Point{node.Column, parent.Idx, MERGE_BACK})
							node.SetPathColor(parent.Id, node.Color)
						} else {
							node.Append(parent.Id, Point{parent.Column, node.Idx, MERGE_TO})
							node.SetPathColor(parent.Id, parent.Color)
						}
					}
				}
			}

			node.Append(parent.Id, Point{parent.Column, parent.Idx, PIPE})

		}
	}

	// Deduplicate path nodes
	for _, node := range nodes {
		for pathIdx, path := range node.ParentsPaths {
			previousPoint := Point{-1, -1, -1}
			for pointIdx, point := range path.Path {
				if point.X == previousPoint.X && point.Y == previousPoint.Y && point.Type == previousPoint.Type {
					tmp := node.ParentsPaths[pathIdx]
					tmp.Path = append(tmp.Path[:pointIdx], tmp.Path[pointIdx+1:]...)
					node.ParentsPaths[pathIdx] = tmp
				}
				previousPoint = point
			}
		}
	}
}

func Get(inputNodes []map[string]interface{}) ([]map[string]interface{}, error) {
	myColors := defaultColors
	nodes, err := buildTree(inputNodes, myColors)
	for _, node := range nodes {
		delete(node, "parentsPaths")
	}
	return nodes, err
}

func buildTree(inputNodes []map[string]interface{}, myColors []string) ([]map[string]interface{}, error) {
	colors = myColors
	var nodes []*OutputNode = initNodes(inputNodes)
	var index map[string]*OutputNode = initIndex(nodes)

	initChildren(nodes, index)
	setColumns(nodes, index)

	for _, node := range nodes {
		for parentId, path := range node.ParentsPaths {
			node.FinalParentsPaths = append(node.FinalParentsPaths, Path{parentId, path.Path, path.Color})
		}
	}
	finalStruct := []map[string]interface{}{}
	for _, node := range nodes {
		finalNode := map[string]interface{}{}
		for key, value := range node.InitialNode {
			finalNode[key] = value
		}
		finalNode["parentsPaths"] = node.ParentsPaths // Kept for tests
		finalNode["id"] = node.Id
		finalNode["parents"] = node.Parents
		finalNode["column"] = node.Column
		finalNode["parents_paths"] = node.FinalParentsPaths
		finalNode["idx"] = node.Idx
		finalNode["color"] = node.Color
		if debugMode {
			finalNode["debug"] = node.Debug
		}
		finalStruct = append(finalStruct, finalNode)
	}

	return finalStruct, nil
}

func getInputNodesFromFile(filePath string) (nodes []map[string]interface{}, err error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	inputJson := string(bytes)
	nodes, err = getInputNodesFromJson(inputJson)
	if err != nil {
		return
	}
	return
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func getInputNodesFromRepo() (nodes []map[string]interface{}, err error) {
	START_OF_COMMIT := "@@@@@@@@@@"
	outBytes, err := exec.Command("git", "log", "--pretty=tformat:"+START_OF_COMMIT+"%n%H%n%aN%n%aE%n%at%n%ai%n%P%n%T%n%s", "--date=local", "--branches", "--remotes").Output()
	if err != nil {
		return
	}
	outString := string(outBytes)
	lines := strings.Split(outString, "\n")
	i := 0
	for i < len(lines) {
		if i >= len(lines) {
			break
		}
		i++
		sha := lines[i]
		//name := lines[i+1]
		//email := lines[i+2]
		//date := lines[i+3]
		//dateIso := lines[i+4]
		parents := strings.Split(lines[i+5], " ")
		parents = deleteEmpty(parents)
		//tree := lines[i+6]
		//subject := lines[i+7]
		i += 8
		node := map[string]interface{}{}
		node["id"] = sha
		node["parents"] = parents
		nodes = append(nodes, node)
		if lines[i] != START_OF_COMMIT {
			break
		}
	}
	return
}

func bootstrap(c *cli.Context) {
	var nodes []map[string]interface{}
	var err error
	jsonFlag := c.String("json")
	fileFlag := c.String("file")
	debugMode = c.Bool("debug")
	repoFlag := c.Bool("repo")
	repoLinearFlag := c.Bool("repo-linear")

	if repoFlag {
		nodes, err = getInputNodesFromRepo()
	} else if repoLinearFlag {
		nodes, err = getInputNodesFromRepo()
		serializeOutput(nodes)
		return
	} else if jsonFlag != "" {
		nodes, err = getInputNodesFromJson(jsonFlag)
	} else if fileFlag != "" {
		nodes, err = getInputNodesFromFile(fileFlag)
	} else {
		cli.ShowAppHelp(c)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	myColors := defaultColors

	out, err := buildTree(nodes, myColors)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, node := range out {
		delete(node, "parentsPaths")
	}

	serializeOutput(out)
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
		cli.BoolFlag{
			Name:  "d, debug",
			Usage: "Debug mode",
		},
		cli.BoolFlag{
			Name:  "r, repo",
			Usage: "Repository",
		},
		cli.BoolFlag{
			Name:  "l, repo-linear",
			Usage: "Repository linear history",
		},
	}
	app.Action = bootstrap
	app.Run(os.Args)
}
