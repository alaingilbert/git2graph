package git2graph

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// DebugMode Debug mode
var DebugMode = false

// NoOutput No output
var NoOutput = false

// Color structure
type Color struct {
	ReleaseIdx int
	color      string
	InUse      bool
}

// DefaultColors Default colors
var DefaultColors = []string{
	"#5aa1be",
	"#c065b8",
	"#c0ab5f",
	"#59bc95",
	"#7a63be",
	"#c0615b",
	"#73bb5e",
	"#6ee585",
	"#7088e8",
	"#eb77a3",
	"#c2e675",
	"#6fdfe9",
	"#d87de8",
	"#eab774",
	"#be82fb",
	"#72d7fc",
	"#adfb82",
}

func getColor(colors []Color, nodeIdx int) string {
	colorToTakeIdx := -1
	for idx, color := range colors {
		if nodeIdx >= color.ReleaseIdx && !color.InUse {
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

func releaseColor(colors []Color, color string, idx int) {
	for colorIdx, colorObj := range colors {
		if color == colorObj.color {
			colors[colorIdx].ReleaseIdx = idx + 2
			colors[colorIdx].InUse = false
			break
		}
	}
}

// Types
const (
	PIPE       pointType = iota // 0: |
	MERGE_BACK                  // 1: ┘
	FORK                        // 2: ┐
	MERGE_TO                    // 3: ┌
)

type pointType uint8

func (p pointType) IsMergeTo() bool { return p == MERGE_TO }

// Point TODO
type Point struct {
	X    int       `json:"x"`
	Y    int       `json:"y"`
	Type pointType `json:"type"`
}

// Path TODO
type Path struct {
	ID    string  `json:"id"`
	Path  []Point `json:"path"`
	Color string  `json:"color"`
}

// OutputNode TODO
type OutputNode struct {
	ID                string         `json:"id"`
	Parents           []string       `json:"parents"`
	Column            int            `json:"column"`
	FinalParentsPaths []Path         `json:"parents_paths"`
	Idx               int            `json:"idx"`
	Color             string         `json:"color"`
	Debug             []string       `json:"debug,omitempty"`
	InitialNode       map[string]any `json:"initial_node"`
	parentsPaths      map[string]Path
	children          []string
	firstInRow        bool
	subBranch         map[string]bool
}

func (node *OutputNode) addDebug(msg string) {
	if DebugMode {
		node.Debug = append(node.Debug, msg)
	}
}

// append a point to a parent path if it is not a duplicate
func (node *OutputNode) noDupAppend(parentID string, point Point) {
	parentPath := node.parentsPaths[parentID]
	if len(parentPath.Path) > 0 && parentPath.Path[len(parentPath.Path)-1] == point {
		return
	}
	node.append(parentID, point)
}

// insert a point to a parent path if it is not a duplicate
func (node *OutputNode) noDupInsert(parentID string, idx int, point Point) {
	parentPath := node.parentsPaths[parentID]
	if parentPath.Path[idx-1] == point {
		return
	}
	node.insert(parentID, idx, point)
}

func (node *OutputNode) append(parentID string, point Point) {
	parentPath := node.parentsPaths[parentID]
	parentPath.Path = append(parentPath.Path, point)
	node.parentsPaths[parentID] = parentPath
}

func (node *OutputNode) remove(parentID string, idx int) {
	parentPath := node.parentsPaths[parentID]
	parentPath.Path = append(parentPath.Path[:idx], parentPath.Path[idx+1:]...)
	node.parentsPaths[parentID] = parentPath
}

func (node *OutputNode) insert(parentID string, idx int, point Point) {
	parentPath := node.parentsPaths[parentID]
	parentPath.Path = append(parentPath.Path, Point{})
	copy(parentPath.Path[idx+1:], parentPath.Path[idx:])
	parentPath.Path[idx] = point
	node.parentsPaths[parentID] = parentPath
}

func (node *OutputNode) columnDefined() bool {
	return node.Column != -1
}

func (node *OutputNode) hasBiggerParentDefined(index *nodesCache) bool {
	found := false
	for _, parentNodeID := range node.Parents {
		parentNode := index.Get(parentNodeID)
		if parentNode.Column > node.Column {
			found = true
			break
		}
	}
	return found
}

func (node *OutputNode) firstInBranch(index *nodesCache) bool {
	for _, parentNodeID := range node.Parents {
		parentNode := index.Get(parentNodeID)
		if !parentNode.columnDefined() || parentNode.Column == node.Column {
			return false
		}
	}
	return true
}

func (node *OutputNode) hasOlderParent(index *nodesCache, idx int) bool {
	found := false
	for _, parentNodeID := range node.Parents {
		parentNode := index.Get(parentNodeID)
		if parentNode.Idx > idx {
			found = true
			break
		}
	}
	return found
}

func (node *OutputNode) setPathColor(parentID, color string) {
	parentPath := node.parentsPaths[parentID]
	parentPath.Color = color
	node.parentsPaths[parentID] = parentPath
}

func (node *OutputNode) getPathColor(parentID string) string {
	return node.parentsPaths[parentID].Color
}

func (node *OutputNode) setPathSubBranch(parentID string) {
	node.subBranch[parentID] = true
}

func (node *OutputNode) isPathSubBranch(parentID string) bool {
	return node.subBranch[parentID]
}

func (node *OutputNode) getPathPoint(parentID string, idx int) (out Point) {
	path := node.parentsPaths[parentID].Path
	pathLen := len(path)
	if idx < 0 {
		rotatedIdx := pathLen + idx
		if rotatedIdx < 0 {
			fields := log.Fields{"idx": idx, "node id": node.ID, "parent id": parentID}
			log.WithFields(fields).Error("Weird, need to investigate")
			return
		}
		idx = rotatedIdx
	}
	return path[idx]
}

// GetPathHeightAtIdx Get the path X at Idx
func (node *OutputNode) GetPathHeightAtIdx(parentID string, lookupIdx int) (height int) {
	height = -1
	firstPoint := node.getPathPoint(parentID, 0)
	lastPoint := node.getPathPoint(parentID, -1)
	if lookupIdx < firstPoint.Y || lookupIdx > lastPoint.Y {
		return
	}
	for _, point := range node.parentsPaths[parentID].Path {
		if point.Y <= lookupIdx {
			height = point.X
		}
	}
	return
}

func (node *OutputNode) pathLength(parentID string) int {
	return len(node.parentsPaths[parentID].Path)
}

// SerializeOutput Json encode object
func SerializeOutput(out []map[string]any) {
	if !NoOutput {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(out); err != nil {
			log.Error("Could not encode json")
		}
	}
}

// GetInputNodesFromJSON Get nodes from json object
func GetInputNodesFromJSON(inputJSON []byte) (nodes []map[string]any, err error) {
	dec := json.NewDecoder(bytes.NewReader(inputJSON))
	err = dec.Decode(&nodes)
	if err != nil {
		return
	}
	for _, node := range nodes {
		parents := make([]string, 0)
		nodeParents, ok := node["parents"]
		if !ok {
			log.Fatal("malformed json input, node missing parents property")
		}
		for _, parent := range nodeParents.([]any) {
			parents = append(parents, parent.(string))
		}
		node["parents"] = parents
	}
	return
}

func initNodes(inputNodes []map[string]any) []*OutputNode {
	out := make([]*OutputNode, 0)
	for idx, node := range inputNodes {
		id, ok := node["id"].(string)
		if !ok {
			log.Fatal("id property must be a string")
		}
		parents, ok := node["parents"].([]string)
		if !ok {
			log.Fatal("parents property must be an array of string")
		}
		newNode := OutputNode{}
		newNode.InitialNode = node
		newNode.ID = id
		newNode.Parents = parents
		newNode.Column = -1
		newNode.parentsPaths = make(map[string]Path)
		newNode.FinalParentsPaths = make([]Path, 0)
		newNode.Idx = idx
		newNode.children = make([]string, 0)
		newNode.Debug = make([]string, 0)
		newNode.subBranch = make(map[string]bool)
		out = append(out, &newNode)
	}
	return out
}

func initIndex(nodes []*OutputNode) *nodesCache {
	index := newNodesCache()
	for _, node := range nodes {
		// Remove bad parents (parents that are before children)
		for idx := len(node.Parents) - 1; idx >= 0; idx-- {
			if index.Has(node.Parents[idx]) {
				node.Parents = append(node.Parents[:idx], node.Parents[idx+1:]...)
			}
		}
		index.Set(node.ID, node)
	}
	return index
}

func initChildren(index *nodesCache, nodes []*OutputNode) {
	for _, node := range nodes {
		for _, parentID := range node.Parents {
			n := index.Get(parentID)
			n.children = append(n.children, node.ID)
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

type nodesCache struct {
	m map[string]*OutputNode
}

func newNodesCache() *nodesCache {
	return &nodesCache{m: make(map[string]*OutputNode)}
}

func (n *nodesCache) Get(key string) *OutputNode {
	return n.m[key]
}

func (n *nodesCache) Set(key string, node *OutputNode) {
	n.m[key] = node
}

func (n *nodesCache) Has(key string) bool {
	_, ok := n.m[key]
	return ok
}

func setColumns(index *nodesCache, colors []Color, nodes []*OutputNode) {
	followingNodesWithChildrenBeforeIdx := newStringSet()
	nextColumn := -1
	incrCol := func() int {
		nextColumn++
		return nextColumn
	}
	for _, node := range nodes {
		// Set column if not defined
		if !node.columnDefined() {
			node.Column = incrCol()
			node.Color = getColor(colors, node.Idx)
		}

		// Cache the following node with child before the current node
		for _, parentID := range node.Parents {
			followingNodesWithChildrenBeforeIdx.Add(parentID)
		}
		followingNodesWithChildrenBeforeIdx.Remove(node.ID)

		// Each child that are merging
		processedNodes := make(map[string]map[string]bool)
		for _, childID := range node.children {
			child := index.Get(childID)
			secondToLastPoint := child.getPathPoint(node.ID, -2)
			if node.Column < secondToLastPoint.X {
				secondPoint := child.getPathPoint(node.ID, 1)
				childIsSubBranch := child.isPathSubBranch(node.ID)
				childHasOlderParent := child.hasOlderParent(index, node.Idx)
				if !childIsSubBranch &&
					!(childHasOlderParent && secondPoint.Type.IsMergeTo()) {
					nextColumn--
				}

				if !child.firstInRow && !childIsSubBranch && !childHasOlderParent {
					child.setPathColor(node.ID, child.Color)
				}
				releaseColor(colors, child.getPathColor(node.ID), node.Idx)

				// Insert before the last element
				pos := child.pathLength(node.ID) - 1
				child.noDupInsert(node.ID, pos, Point{secondToLastPoint.X, node.Idx, MERGE_BACK})

				// Nodes that are following the current node
				for followingNodeID := range followingNodesWithChildrenBeforeIdx.Items {
					followingNode := index.Get(followingNodeID)
					// Following nodes that have a child before the current node
					for _, followingNodeChildID := range followingNode.children {
						followingNodeChild := index.Get(followingNodeChildID)
						// Following node child has a path that is higher than the current path being merged
						if followingNodeChild.Idx < node.Idx &&
							followingNodeChild.GetPathHeightAtIdx(followingNode.ID, node.Idx) > secondToLastPoint.X {

							// Index to delete is the one before last
							idxRemove := followingNodeChild.pathLength(followingNode.ID) - 1
							if idxRemove < 0 || processedNodes[followingNode.ID][followingNodeChild.ID] {
								continue
							}
							// Remove second before last node has same Y, remove the before last node
							for followingNodeChild.getPathPoint(followingNode.ID, idxRemove).Y == followingNodeChild.getPathPoint(followingNode.ID, idxRemove-1).Y {
								followingNodeChild.remove(followingNode.ID, idxRemove-1)
								idxRemove--
							}

							// Calculate nb of merging nodes
							targetColumn := followingNodeChild.GetPathHeightAtIdx(followingNode.ID, node.Idx)
							nbNodesMergingBack := 0
							for _, childID := range node.children {
								child := index.Get(childID)
								childIsSubBranch := child.isPathSubBranch(node.ID)
								childHasOlderParent := child.hasOlderParent(index, node.Idx)
								secondToLastPoint := child.getPathPoint(node.ID, -2)
								secondPoint := child.getPathPoint(node.ID, 1)
								if node.Column < secondToLastPoint.X &&
									secondToLastPoint.X < targetColumn &&
									!childIsSubBranch &&
									!(childHasOlderParent && secondPoint.Type.IsMergeTo()) {
									nbNodesMergingBack++
								}
							}

							pathPointX := followingNodeChild.getPathPoint(followingNode.ID, idxRemove-1).X
							followingNodeChild.remove(followingNode.ID, idxRemove)
							followingNodeChild.noDupAppend(followingNode.ID, Point{pathPointX, node.Idx, MERGE_BACK})
							followingNodeChild.noDupAppend(followingNode.ID, Point{pathPointX - nbNodesMergingBack, node.Idx, PIPE})
							if followingNode.Column <= secondToLastPoint.X {
								followingNodeChild.noDupAppend(followingNode.ID, Point{pathPointX - nbNodesMergingBack, followingNode.Idx, MERGE_BACK})
							} else if processedNodes[followingNode.ID] == nil {
								followingNode.Column -= nbNodesMergingBack
							}
							followingNodeChild.noDupAppend(followingNode.ID, Point{followingNode.Column, followingNode.Idx, PIPE})
							if processedNodes[followingNode.ID] == nil {
								processedNodes[followingNode.ID] = make(map[string]bool)
							}
							processedNodes[followingNode.ID][followingNodeChild.ID] = true
						}
					}
				}
			}
		}

		for parentIdx, parentID := range node.Parents {
			parent := index.Get(parentID)

			node.noDupAppend(parent.ID, Point{node.Column, node.Idx, PIPE})

			if !parent.columnDefined() {
				firstParent := index.Get(node.Parents[0])
				column := node.Column
				color := node.Color
				if parentIdx > 0 && (parentIdx > 1 || firstParent.Column >= node.Column || firstParent.Idx > node.Idx+1) {
					column = incrCol()
					color = getColor(colors, node.Idx)
					node.noDupAppend(parent.ID, Point{column, node.Idx, FORK})
					node.firstInRow = true
				}
				parent.Column = column
				parent.Color = color
				node.setPathColor(parent.ID, parent.Color)
			} else if node.Column < parent.Column && parentIdx == 0 {
				for _, childID := range parent.children {
					child := index.Get(childID)
					if idxRemove := child.pathLength(parent.ID) - 1; idxRemove > 0 {
						child.remove(parent.ID, idxRemove)
						child.noDupAppend(parent.ID, Point{child.getPathPoint(parent.ID, idxRemove-1).X, parent.Idx, MERGE_BACK})
						child.noDupAppend(parent.ID, Point{node.Column, parent.Idx, PIPE})
					}
				}
				parent.Column = node.Column
				parent.Color = node.Color
				node.setPathColor(parent.ID, node.Color)
			} else if node.Column < parent.Column {
				node.setPathSubBranch(parent.ID)
				node.noDupAppend(parent.ID, Point{parent.Column, node.Idx, FORK})
				node.setPathColor(parent.ID, parent.Color)
			} else if node.Column > parent.Column {
				if node.hasBiggerParentDefined(index) || (parentIdx == 0 && (parent.Idx > node.Idx+1 || node.firstInBranch(index))) {
					node.noDupAppend(parent.ID, Point{node.Column, parent.Idx, MERGE_BACK})
					node.setPathColor(parent.ID, node.Color)
				} else {
					node.noDupAppend(parent.ID, Point{parent.Column, node.Idx, MERGE_TO})
					node.setPathColor(parent.ID, parent.Color)
				}
			}
			node.noDupAppend(parent.ID, Point{parent.Column, parent.Idx, PIPE})
		}
	}
}

// Get TODO
func Get(inputNodes []map[string]any) ([]map[string]any, error) {
	myColors := DefaultColors
	nodes, err := BuildTree(inputNodes, myColors)
	for _, node := range nodes {
		delete(node, "parentsPaths")
	}
	return nodes, err
}

// GetPaginated TODO
func GetPaginated(inputNodes []map[string]any, from, size int) ([]map[string]any, error) {
	myColors := DefaultColors
	nodes, err := BuildTree(inputNodes, myColors)
	for _, node := range nodes {
		delete(node, "parentsPaths")
	}
	return nodes[from : from+size], err
}

// BuildTree TODO
func BuildTree(inputNodes []map[string]any, myColors []string) ([]map[string]any, error) {
	colors := make([]Color, 0)
	for _, colorStr := range myColors {
		colors = append(colors, Color{color: colorStr})
	}

	nodes := initNodes(inputNodes)
	index := initIndex(nodes)

	initChildren(index, nodes)
	setColumns(index, colors, nodes)

	for _, node := range nodes {
		for parentID, path := range node.parentsPaths {
			node.FinalParentsPaths = append(node.FinalParentsPaths, Path{parentID, path.Path, path.Color})
		}
	}
	finalStruct := make([]map[string]any, 0)
	for _, node := range nodes {
		finalNode := map[string]any{}
		for key, value := range node.InitialNode {
			finalNode[key] = value
		}
		finalNode["parentsPaths"] = node.parentsPaths // Kept for tests
		finalNode["id"] = node.ID
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
func GetInputNodesFromFile(filePath string) (nodes []map[string]any, err error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	nodes, err = GetInputNodesFromJSON(fileBytes)
	if err != nil {
		return
	}
	return
}

func deleteEmpty(s []string) []string {
	r := make([]string, 0)
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// GetInputNodesFromRepo TODO
func GetInputNodesFromRepo(seqIds bool) (nodes []map[string]any, err error) {
	startOfCommit := "@@@@@@@@@@"
	outBytes, err := exec.Command("git", "log", "--pretty=tformat:"+startOfCommit+"%n%H%n%aN%n%aE%n%at%n%ai%n%P%n%T%n%s", "--date=local", "--branches", "--remotes").Output()
	if err != nil {
		return
	}
	outString := string(outBytes)
	lines := strings.Split(outString, "\n")
	ids := 0
	shaMap := make(map[string]string)
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
		node := map[string]any{}
		if seqIds {
			id := strconv.Itoa(ids)
			shaMap[sha] = id
			node["id"] = id
		} else {
			node["id"] = sha
		}
		node["parents"] = parents
		nodes = append(nodes, node)
		ids++
		if lines[i] != startOfCommit {
			break
		}
	}
	if seqIds {
		for _, node := range nodes {
			mappedParents := make([]string, 0)
			for _, parentSha := range node["parents"].([]string) {
				mappedParents = append(mappedParents, shaMap[parentSha])
			}
			node["parents"] = mappedParents
		}
	}
	return
}
