package git2graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var colors []Color

// DebugMode Debug mode
var DebugMode = false

// NoOutput No output
var NoOutput = false

var index = make(map[string]*OutputNode)

// Color color structure
type Color struct {
	ReleaseIdx int
	color      string
	InUse      bool
}

// DefaultColors Default colors
var DefaultColors = []Color{
	Color{-2, "#5aa1be", false},
	Color{-2, "#c065b8", false},
	Color{-2, "#c0ab5f", false},
	Color{-2, "#59bc95", false},
	Color{-2, "#7a63be", false},
	Color{-2, "#c0615b", false},
	Color{-2, "#73bb5e", false},
	Color{-2, "#6ee585", false},
	Color{-2, "#7088e8", false},
	Color{-2, "#eb77a3", false},
	Color{-2, "#c2e675", false},
	Color{-2, "#6fdfe9", false},
	Color{-2, "#d87de8", false},
	Color{-2, "#eab774", false},
}

func getColor(nodeIdx int) string {
	colorToTakeIdx := -1
	for idx, color := range colors {
		if nodeIdx >= color.ReleaseIdx+2 && !color.InUse {
			colorToTakeIdx = idx
			break
		}
	}
	if colorToTakeIdx == -1 {
		log.Error("Not enough colors")
		return "#000"
	}
	colors[colorToTakeIdx].InUse = true
	return colors[colorToTakeIdx].color
}

func releaseColor(color string, idx int) {
	for colorIdx, colorObj := range colors {
		if color == colorObj.color {
			colors[colorIdx].ReleaseIdx = idx
			colors[colorIdx].InUse = false
			break
		}
	}
}

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
	InitialNode       map[string]interface{} `json:"initial_node"`
	SubBranch         map[string]bool        `json:"-"`
}

func (node *OutputNode) append(parentID string, point Point) {
	tmp := node.ParentsPaths[parentID]
	tmp.Path = append(tmp.Path, point)
	node.ParentsPaths[parentID] = tmp
}

func (node *OutputNode) remove(parentID string, idx int) {
	tmp := node.ParentsPaths[parentID]
	tmp.Path = append(tmp.Path[:idx], tmp.Path[idx+1:]...)
	node.ParentsPaths[parentID] = tmp
}

func (node *OutputNode) insert(parentID string, idx int, point Point) {
	tmp := node.ParentsPaths[parentID]
	tmp.Path = append(tmp.Path, Point{})
	copy(tmp.Path[idx+1:], tmp.Path[idx:])
	tmp.Path[idx] = point
	node.ParentsPaths[parentID] = tmp
}

func (node *OutputNode) columnDefined() bool {
	return node.Column != -1
}

func (node *OutputNode) hasBiggerParentDefined() bool {
	found := false
	for _, parentNodeID := range node.Parents {
		parentNode := index[parentNodeID]
		if parentNode.Column > node.Column {
			found = true
			break
		}
	}
	return found
}

func (node *OutputNode) setPathColor(parentID, color string) {
	tmp := node.ParentsPaths[parentID]
	tmp.Color = color
	node.ParentsPaths[parentID] = tmp
}

func (node *OutputNode) getPathColor(parentID string) string {
	return node.ParentsPaths[parentID].Color
}

func (node *OutputNode) setPathSubBranch(parentID string) {
	node.SubBranch[parentID] = true
}

func (node *OutputNode) isPathSubBranch(parentID string) bool {
	return node.SubBranch[parentID]
}

func (node *OutputNode) getPathPoint(parentID string, idx int) Point {
	if idx < 0 {
		if len(node.ParentsPaths[parentID].Path)+idx < 0 {
			if index[parentID].Idx < index[node.Id].Idx {
				log.WithFields(log.Fields{
					"idx":       idx,
					"node id":   node.Id,
					"parent id": parentID,
				}).Error("Error in repo structure. parent idx < node idx")
				return Point{}
			}
			log.WithFields(log.Fields{
				"idx":       idx,
				"node id":   node.Id,
				"parent id": parentID,
			}).Error("1- Weird, need to investigate")
			return Point{}
		}
		return node.ParentsPaths[parentID].Path[len(node.ParentsPaths[parentID].Path)+idx]
	}
	return node.ParentsPaths[parentID].Path[idx]
}

//GetPathHeightAtIdx Get the path X at Idx
func (node *OutputNode) GetPathHeightAtIdx(parentID string, lookupIdx int) (height int) {
	height = -1
	firstPoint := node.getPathPoint(parentID, 0)
	lastPoint := node.getPathPoint(parentID, -1)
	if lookupIdx < firstPoint.Y || lookupIdx > lastPoint.Y {
		return
	}
	for _, point := range node.ParentsPaths[parentID].Path {
		if point.Y <= lookupIdx {
			height = point.X
		}
	}
	return
}

func (node *OutputNode) pathLength(parentID string) int {
	return len(node.ParentsPaths[parentID].Path)
}

// SerializeOutput Json encode object
func SerializeOutput(out []map[string]interface{}) {
	if !NoOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(out)
	}
}

// GetInputNodesFromJSON Get nodes from json object
func GetInputNodesFromJSON(inputJSON string) (nodes []map[string]interface{}, err error) {
	dec := json.NewDecoder(strings.NewReader(inputJSON))
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
		newNode.SubBranch = make(map[string]bool)
		out = append(out, &newNode)
	}
	return out
}

func initIndex(nodes []*OutputNode) map[string]*OutputNode {
	index := make(map[string]*OutputNode)
	for _, node := range nodes {
		// Remove bad parents (parents that are before children)
		for idx := len(node.Parents) - 1; idx >= 0; idx-- {
			if index[node.Parents[idx]] != nil {
				node.Parents = append(node.Parents[:idx], node.Parents[idx+1:]...)
			}
		}
		index[node.Id] = node
	}
	return index
}

func initChildren(nodes []*OutputNode) {
	for _, node := range nodes {
		for _, parentID := range node.Parents {
			index[parentID].Children = append(index[parentID].Children, node.Id)
		}
	}
}

type stringSet struct {
	Items map[string]bool
}

func newStringSet() stringSet {
	s := stringSet{}
	s.Items = make(map[string]bool)
	return s
}

func (s *stringSet) Add(in string) {
	s.Items[in] = true
}

func (s *stringSet) Remove(in string) {
	delete(s.Items, in)
}

func setColumns(nodes []*OutputNode) {
	followingNodesWithChildrenBeforeIdx := newStringSet()
	nextColumn := 0
	for _, node := range nodes {
		// Set column if not defined
		if !node.columnDefined() {
			node.Column = nextColumn
			if DebugMode {
				node.Debug = append(node.Debug, fmt.Sprintf("Column set to %d", nextColumn))
			}
			node.Color = getColor(node.Idx)
			nextColumn++
			log.WithFields(log.Fields{
				"nextColumn": nextColumn,
				"operator":   "++",
				"created":    node.Id,
			}).Debug("new node ++")
		}

		// Cache the following node with child before the current node
		for _, parentID := range node.Parents {
			followingNodesWithChildrenBeforeIdx.Add(parentID)
		}
		followingNodesWithChildrenBeforeIdx.Remove(node.Id)

		// Each children that are merging
		processedNodes := make(map[string]map[string]bool)
		for _, childID := range node.Children {
			child := index[childID]
			if node.Column < child.getPathPoint(node.Id, -2).X {
				if !child.isPathSubBranch(node.Id) {
					nextColumn--
					log.WithFields(log.Fields{
						"nextColumn": nextColumn,
						"operator":   "--",
						"merging":    child.Id,
						"into":       node.Id,
						"sub":        child.isPathSubBranch(node.Id),
					}).Debug("node merging --")
					releaseColor(child.getPathColor(node.Id), node.Idx)
				}

				if !child.FirstInRow && !child.isPathSubBranch(node.Id) {
					child.setPathColor(node.Id, child.Color)
				}
				releaseColor(child.getPathColor(node.Id), node.Idx)

				// Insert before the last element
				pos := child.pathLength(node.Id) - 1
				point := Point{child.getPathPoint(node.Id, -2).X, node.Idx, MERGE_BACK}
				child.insert(node.Id, pos, point)

				// Nodes that are following the current node
				for followingNodeID := range followingNodesWithChildrenBeforeIdx.Items {
					followingNode := index[followingNodeID]
					if followingNode.Idx > node.Idx {
						// Following nodes that have a child before the current node
						for _, followingNodeChildID := range followingNode.Children {
							followingNodeChild := index[followingNodeChildID]
							if followingNodeChild.Idx < node.Idx {
								// Following node child has a path that is higher than the current path being merged
								if followingNodeChild.GetPathHeightAtIdx(followingNode.Id, node.Idx) > child.getPathPoint(node.Id, -2).X {

									// Index to delete is the one before last
									idxRemove := followingNodeChild.pathLength(followingNode.Id) - 1
									if idxRemove < 0 {
										continue
									}
									// Remove second before last node has same Y, remove the before last node
									if followingNodeChild.pathLength(followingNode.Id) > idxRemove &&
										followingNodeChild.getPathPoint(followingNode.Id, idxRemove).Y == followingNodeChild.getPathPoint(followingNode.Id, idxRemove-1).Y {
										followingNodeChild.remove(followingNode.Id, idxRemove-1)
										idxRemove--
									}

									// Calculate nb of merging nodes
									nbNodesMergingBack := 0
									for _, childID := range node.Children {
										child := index[childID]
										if node.Column < child.getPathPoint(node.Id, -2).X &&
											child.getPathPoint(node.Id, -2).X < followingNodeChild.GetPathHeightAtIdx(followingNode.Id, node.Idx) &&
											!child.isPathSubBranch(node.Id) {
											nbNodesMergingBack++
										}
									}

									if processedNodes[followingNode.Id] != nil && processedNodes[followingNode.Id][followingNodeChild.Id] {
										continue
									}
									tmp := followingNodeChild.getPathPoint(followingNode.Id, idxRemove-1).X
									followingNodeChild.remove(followingNode.Id, idxRemove)
									followingNodeChild.append(followingNode.Id, Point{tmp, node.Idx, MERGE_BACK})
									followingNodeChild.append(followingNode.Id, Point{tmp - 1 - (nbNodesMergingBack - 1), node.Idx, PIPE})
									if followingNode.Column > child.getPathPoint(node.Id, -2).X {
										if processedNodes[followingNode.Id] == nil {
											followingNodeChild.append(followingNode.Id, Point{followingNode.Column - (nbNodesMergingBack - 1) - 1, followingNode.Idx, PIPE})
											followingNode.Column -= nbNodesMergingBack
										} else {
											followingNodeChild.append(followingNode.Id, Point{followingNode.Column, followingNode.Idx, PIPE})
										}
										if DebugMode {
											followingNode.Debug = append(followingNode.Debug, fmt.Sprintf("Column minus %s, %s, %d, %d", followingNode.Id, child.Id, followingNode.Column, nbNodesMergingBack))
										}
									} else {
										followingNodeChild.append(followingNode.Id, Point{tmp - 1 - (nbNodesMergingBack - 1), followingNode.Idx, MERGE_BACK})
										followingNodeChild.append(followingNode.Id, Point{followingNode.Column, followingNode.Idx, PIPE})
									}
									if processedNodes[followingNode.Id] == nil {
										processedNodes[followingNode.Id] = make(map[string]bool)
									}
									processedNodes[followingNode.Id][followingNodeChild.Id] = true
								}
							}
						}
					}
				}
			}
		}

		for parentIdx, parentID := range node.Parents {
			parent := index[parentID]

			node.append(parent.Id, Point{node.Column, node.Idx, PIPE})

			if !parent.columnDefined() {
				if parentIdx == 0 || (parentIdx == 1 && index[node.Parents[0]].Column < node.Column && index[node.Parents[0]].Idx == node.Idx+1) {
					parent.Column = node.Column
					if DebugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("1- Column set to %d", node.Column))
					}
					parent.Color = node.Color
					node.setPathColor(parent.Id, parent.Color)
				} else {
					parent.Column = nextColumn
					if DebugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("2- Column set to %d", nextColumn))
					}
					parent.Color = getColor(node.Idx)
					node.append(parent.Id, Point{parent.Column, node.Idx, FORK})
					node.setPathColor(parent.Id, parent.Color)
					node.FirstInRow = true
					nextColumn++
					log.WithFields(log.Fields{
						"nextColumn": nextColumn,
						"operator":   "++",
						"node":       node.Id,
						"parent":     parent.Id,
					}).Debug("new parent undefined column++")

				}
			} else if parent.columnDefined() {
				if node.Column < parent.Column && parentIdx == 0 {
					for _, childID := range parent.Children {
						child := index[childID]
						idxRemove := child.pathLength(parent.Id) - 1
						if idxRemove > 0 {
							if child.getPathPoint(parent.Id, idxRemove).Type != FORK {
								child.remove(parent.Id, idxRemove)
							}
							pos := child.pathLength(parent.Id) - 1
							child.append(parent.Id, Point{child.getPathPoint(parent.Id, pos).X, parent.Idx, MERGE_BACK})
							child.append(parent.Id, Point{node.Column, parent.Idx, PIPE})
						}
					}
					parent.Column = node.Column
					if DebugMode {
						parent.Debug = append(parent.Debug, fmt.Sprintf("Column reset to %d", node.Column))
					}
					parent.Color = node.Color
					node.setPathColor(parent.Id, node.Color)
				} else if node.Column < parent.Column && parentIdx > 0 {
					node.setPathSubBranch(parent.Id)
					node.append(parent.Id, Point{parent.Column, node.Idx, FORK})
					node.setPathColor(parent.Id, parent.Color)
				} else if node.Column > parent.Column {
					if len(node.Parents) > 1 {
						if node.hasBiggerParentDefined() || (parentIdx == 0 && parent.Idx > node.Idx+1) {
							node.append(parent.Id, Point{node.Column, parent.Idx, MERGE_BACK})
							node.setPathColor(parent.Id, node.Color)
						} else {
							node.append(parent.Id, Point{parent.Column, node.Idx, MERGE_TO})
							node.setPathColor(parent.Id, parent.Color)
						}
					}
				}
			}

			node.append(parent.Id, Point{parent.Column, parent.Idx, PIPE})

		}
	}

	// Deduplicate path nodes
	for _, node := range nodes {
		for parentID, path := range node.ParentsPaths {
			previousPoint := Point{-1, -1, -1}
			for pointIdx := len(path.Path) - 1; pointIdx >= 0; pointIdx-- {
				point := path.Path[pointIdx]
				if point.X == previousPoint.X && point.Y == previousPoint.Y && point.Type == previousPoint.Type {
					parentPath := node.ParentsPaths[parentID]
					parentPath.Path = append(parentPath.Path[:pointIdx], parentPath.Path[pointIdx+1:]...)
					node.ParentsPaths[parentID] = parentPath
				}
				previousPoint = point
			}
		}
	}
}

// Get TODO
func Get(inputNodes []map[string]interface{}) ([]map[string]interface{}, error) {
	myColors := DefaultColors
	nodes, err := BuildTree(inputNodes, myColors)
	for _, node := range nodes {
		delete(node, "parentsPaths")
	}
	return nodes, err
}

// GetPaginated TODO
func GetPaginated(inputNodes []map[string]interface{}, from, size int) ([]map[string]interface{}, error) {
	myColors := DefaultColors
	nodes, err := BuildTree(inputNodes, myColors)
	for _, node := range nodes {
		delete(node, "parentsPaths")
	}
	return nodes[from : from+size], err
}

// BuildTree TODO
func BuildTree(inputNodes []map[string]interface{}, myColors []Color) ([]map[string]interface{}, error) {
	colors = make([]Color, 0)
	for _, color := range myColors {
		colors = append(colors, color)
	}

	nodes := initNodes(inputNodes)
	index = initIndex(nodes)

	initChildren(nodes)
	setColumns(nodes)

	for _, node := range nodes {
		for parentID, path := range node.ParentsPaths {
			node.FinalParentsPaths = append(node.FinalParentsPaths, Path{parentID, path.Path, path.Color})
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
		if DebugMode {
			finalNode["debug"] = node.Debug
		}
		finalStruct = append(finalStruct, finalNode)
	}

	return finalStruct, nil
}

// GetInputNodesFromFile TODO
func GetInputNodesFromFile(filePath string) (nodes []map[string]interface{}, err error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	inputJSON := string(bytes)
	nodes, err = GetInputNodesFromJSON(inputJSON)
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

// GetInputNodesFromRepo TODO
func GetInputNodesFromRepo() (nodes []map[string]interface{}, err error) {
	startOfCommit := "@@@@@@@@@@"
	outBytes, err := exec.Command("git", "log", "--pretty=tformat:"+startOfCommit+"%n%H%n%aN%n%aE%n%at%n%ai%n%P%n%T%n%s", "--date=local", "--branches", "--remotes").Output()
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
		if lines[i] != startOfCommit {
			break
		}
	}
	return
}
