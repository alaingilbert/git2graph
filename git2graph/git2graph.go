package git2graph

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NoOutput No output
var NoOutput = false

// Color structure
type Color struct {
	releaseIdx int
	color      string
	inUse      bool
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

type IColorGenerator interface {
	GetColor(idx int) string
}

// SimpleColorGen is a color generator that take a static colors array, or return black when running out of colors
type SimpleColorGen struct {
	colors []string
}

// NewSimpleColorGen creates a new SimpleColorGen
func NewSimpleColorGen(colors []string) *SimpleColorGen {
	return &SimpleColorGen{colors: colors}
}

func (c *SimpleColorGen) GetColor(idx int) string {
	if idx >= len(c.colors) {
		log.Error("Not enough colors")
		return "#000"
	}
	return c.colors[idx]
}

type colorsManager struct {
	g IColorGenerator
	m map[int]*Color
}

func newColorsManager(colorGen IColorGenerator) *colorsManager {
	return &colorsManager{g: colorGen, m: make(map[int]*Color)}
}

func getColor(colorsManager *colorsManager, nodeIdx int) int {
	var color *Color
	i := 0
	for {
		var ok bool
		color, ok = colorsManager.m[i]
		if !ok {
			color = &Color{color: colorsManager.g.GetColor(i), releaseIdx: -1}
			colorsManager.m[i] = color
			break
		}
		if nodeIdx >= color.releaseIdx && !color.inUse {
			break
		}
		i++
	}
	color.inUse = true
	return i
}

func releaseColor(colorsMan *colorsManager, colorIdx int, idx int) {
	for i := range colorsMan.m {
		if i == colorIdx {
			color := colorsMan.m[i]
			color.releaseIdx = idx + 2
			color.inUse = false
			break
		}
	}
}

// Types; to understand these constants, you need to read the graph from top to bottom.
// Fork is when a node fork into two paths (top -> bottom)
// MergeBack is when a branch merge back into a branch on its right
// MergeTo is when a branch on the right merge into a branch on its left
const (
	Pipe      pointType = iota // 0: |
	MergeBack                  // 1: ┘
	Fork                       // 2: ┐
	MergeTo                    // 3: ┌
)

type pointType uint8

func (p pointType) IsMergeTo() bool { return p == MergeTo }

// Node is the raw information for a commit
type Node map[string]any

// Path defines how to draw a line in between a parent and child nodes
type Path struct {
	ID       string
	Path     []Point
	ColorIdx int
}

// Return either or not a path is valid (has at least 2 points)
func (p Path) isValid() bool {
	return len(p.Path) >= 2
}

// Return either or not the path is of type "Fork"
func (p Path) isFork() bool {
	return p.isValid() && p.Path[SecondPt].Type == Fork
}

// Return either or not the path is of type "MergeTo"
func (p Path) isMergeTo() bool {
	return p.isValid() && p.Path[SecondPt].Type == MergeTo
}

// Point is one part of a path
type Point struct {
	X    int       `json:"x"`
	Y    int       `json:"y"`
	Type pointType `json:"type"`
}

// parents are the node below the current node
// children are the nodes above the current node
// A node only ever have at most 2 parents.
type internalNode struct {
	InitialNode   Node
	ID            string
	Idx           int
	Column        int
	ColorIdx      int
	firstOfBranch bool
	Parents       []string
	children      []string
	parentsPaths  map[string]*Path
}

// A node is a "firstOfBranch" if there is a path to a parent that needs a new color,
// and the commit is the first commit in that new branch.
func (n *internalNode) isFirstOfBranch() bool {
	return n.firstOfBranch
}

func (n *internalNode) setFirstOfBranch() {
	n.firstOfBranch = true
}

func (n *internalNode) pathTo(parentID string) *Path {
	parentPath, ok := n.parentsPaths[parentID]
	if !ok {
		parentPath = &Path{}
		n.parentsPaths[parentID] = parentPath
	}
	return parentPath
}

// append a point to a parent path if it is not a duplicate
func (n *internalNode) noDupAppend(parentID string, point Point) {
	parentPath := n.pathTo(parentID)
	if len(parentPath.Path) > 0 && parentPath.Path[len(parentPath.Path)-1] == point {
		return
	}
	n.append(parentID, point)
}

// insert a point to a parent path if it is not a duplicate
func (n *internalNode) noDupInsert(parentID string, idx int, point Point) {
	parentPath := n.parentsPaths[parentID]
	if idx < 0 {
		idx = len(parentPath.Path) + idx
	}
	if parentPath.Path[idx-1] == point {
		return
	}
	n.insert(parentID, idx, point)
}

func (n *internalNode) append(parentID string, point Point) {
	parentPath := n.pathTo(parentID)
	parentPath.Path = append(parentPath.Path, point)
}

func (n *internalNode) remove(parentID string, idx int) {
	parentPath := n.parentsPaths[parentID]
	parentPath.Path = append(parentPath.Path[:idx], parentPath.Path[idx+1:]...)
	n.parentsPaths[parentID] = parentPath
}

func (n *internalNode) insert(parentID string, idx int, point Point) {
	parentPath := n.parentsPaths[parentID]
	parentPath.Path = append(parentPath.Path, Point{})
	copy(parentPath.Path[idx+1:], parentPath.Path[idx:])
	parentPath.Path[idx] = point
	n.parentsPaths[parentID] = parentPath
}

func (n *internalNode) columnDefined() bool {
	return n.Column != -1
}

func (n *internalNode) firstInBranch(index *nodesCache) bool {
	for _, parentNodeID := range n.Parents {
		parentNode := index.Get(parentNodeID)
		if !parentNode.columnDefined() || parentNode.Column == n.Column {
			return false
		}
	}
	return true
}

func (n *internalNode) hasBiggerParentDefined(index *nodesCache) bool {
	for _, parentNodeID := range n.Parents {
		parentNode := index.Get(parentNodeID)
		if parentNode.Column > n.Column {
			return true
		}
	}
	return false
}

// Return either or not the node has a parent that has higher "Idx" than the one in parameter
func (n *internalNode) hasOlderParent(index *nodesCache, idx int) bool {
	for _, parentNodeID := range n.Parents {
		parentNode := index.Get(parentNodeID)
		if parentNode.Idx > idx {
			return true
		}
	}
	return false
}

func (n *internalNode) setPathColor(parentID string, color int) {
	parentPath := n.parentsPaths[parentID]
	parentPath.ColorIdx = color
	n.parentsPaths[parentID] = parentPath
}

func (n *internalNode) getPathColor(parentID string) int {
	return n.parentsPaths[parentID].ColorIdx
}

// Return either or not the path to a parent is of type "MergeTo"
func (n *internalNode) isMergeTo(parentID string) bool {
	return n.parentsPaths[parentID].isMergeTo()
}

// A subbranch, is when the child node is in the middle of another branch
// See test_022.png node #4 (zero-indexed)
func (n *internalNode) isPathSubBranch(parentID string) bool {
	return n.parentsPaths[parentID].isFork() && !n.isFirstOfBranch()
}

const (
	FirstPt        = 0
	SecondPt       = 1
	LastPt         = -1
	SecondToLastPt = -2
)

const (
	idKey               = "id"
	parentsKey          = "parents"
	gKey                = "g"
	parentsPathsTestKey = "parentsPaths"
)

// Return the point at idx in the path to parentID
// 0 return the first point
// 1 return the second point
// -1 return the last point
// -2 return the second to last point
func (n *internalNode) getPathPoint(parentID string, idx int) (out Point) {
	path := n.parentsPaths[parentID].Path
	pathLen := len(path)
	if idx < 0 {
		rotatedIdx := pathLen + idx
		if rotatedIdx < 0 {
			fields := log.Fields{"idx": idx, "n id": n.ID, "parent id": parentID}
			log.WithFields(fields).Error("Weird, need to investigate")
			return
		}
		idx = rotatedIdx
	}
	return path[idx]
}

// Return either or not the path to the parent is a MergeTo
func (n *internalNode) pathIsMergeTo(parentID string) bool {
	return n.getPathPoint(parentID, SecondPt).Type.IsMergeTo()
}

// GetPathHeightAtIdx Get the path X at Idx
func (n *internalNode) GetPathHeightAtIdx(parentID string, lookupIdx int) (height int) {
	height = -1
	firstPoint := n.getPathPoint(parentID, FirstPt)
	lastPoint := n.getPathPoint(parentID, LastPt)
	if lookupIdx < firstPoint.Y || lookupIdx > lastPoint.Y {
		return
	}
	for _, point := range n.parentsPaths[parentID].Path {
		if point.Y <= lookupIdx {
			height = point.X
		}
	}
	return
}

func (n *internalNode) pathLength(parentID string) int {
	parentPath := n.pathTo(parentID)
	return len(parentPath.Path)
}

// A merging node is one that come from a higher column, but is not a sub-branch and is not a MergeTo
func (n *internalNode) nbNodesMergingBack(index *nodesCache, maxX int) (nbNodesMergingBack int) {
	nodeID := n.ID
	for _, childID := range n.children {
		child := index.Get(childID)
		childIsSubBranch := child.isPathSubBranch(nodeID)
		secondToLastPoint := child.getPathPoint(nodeID, SecondToLastPt)
		if n.Column < secondToLastPoint.X && secondToLastPoint.X < maxX &&
			!childIsSubBranch &&
			!child.pathIsMergeTo(nodeID) {
			nbNodesMergingBack++
		}
	}
	return
}

// SerializeOutput Json encode object
func SerializeOutput(out []Node) {
	if !NoOutput {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(out); err != nil {
			log.Error("Could not encode json")
		}
	}
}

// GetInputNodesFromJSON Get nodes from json object
func GetInputNodesFromJSON(inputJSON []byte) (nodes []Node, err error) {
	dec := json.NewDecoder(bytes.NewReader(inputJSON))
	err = dec.Decode(&nodes)
	if err != nil {
		return
	}
	for _, node := range nodes {
		parents := make([]string, 0)
		nodeParents, ok := node[parentsKey]
		if !ok {
			log.Fatal("malformed json input, node missing parents property")
		}
		for _, parent := range nodeParents.([]any) {
			parents = append(parents, parent.(string))
		}
		node[parentsKey] = parents
	}
	return
}

func initNodes(inputNodes []Node) []*internalNode {
	out := make([]*internalNode, 0)
	for idx, node := range inputNodes {
		id, ok := node[idKey].(string)
		if !ok {
			log.Fatal("id property must be a string")
		}
		parents, ok := node[parentsKey].([]string)
		if !ok {
			log.Fatal("parents property must be an array of string")
		}
		newNode := internalNode{}
		newNode.InitialNode = node
		newNode.ID = id
		newNode.Idx = idx
		newNode.Column = -1
		newNode.Parents = parents
		newNode.parentsPaths = make(map[string]*Path)
		newNode.children = make([]string, 0)
		out = append(out, &newNode)
	}
	return out
}

func initIndex(nodes []*internalNode) *nodesCache {
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

func initChildren(index *nodesCache, nodes []*internalNode) {
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
	m map[string]*internalNode
}

func newNodesCache() *nodesCache {
	return &nodesCache{m: make(map[string]*internalNode)}
}

func (n *nodesCache) Get(key string) *internalNode {
	return n.m[key]
}

func (n *nodesCache) Set(key string, node *internalNode) {
	n.m[key] = node
}

func (n *nodesCache) Has(key string) bool {
	_, ok := n.m[key]
	return ok
}

type processedNodes struct {
	m map[string]map[string]bool
}

func newProcessedNodes() *processedNodes {
	return &processedNodes{m: make(map[string]map[string]bool)}
}

func (p *processedNodes) HasNode(nodeID string) bool {
	return p.m[nodeID] != nil
}

func (p *processedNodes) HasChild(nodeID, childID string) bool {
	return p.m[nodeID][childID]
}

func (p *processedNodes) Set(nodeID, childID string) {
	if p.m[nodeID] == nil {
		p.m[nodeID] = make(map[string]bool)
	}
	p.m[nodeID][childID] = true
}

func setColumns(index *nodesCache, colorsMan *colorsManager, nodes []*internalNode) {
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
			node.ColorIdx = getColor(colorsMan, node.Idx)
		}

		// Cache the following node with child before the current node
		for _, parentID := range node.Parents {
			followingNodesWithChildrenBeforeIdx.Add(parentID)
		}
		followingNodesWithChildrenBeforeIdx.Remove(node.ID)

		// Each child that are merging
		processedNodesInst := newProcessedNodes()
		for _, childID := range node.children {
			child := index.Get(childID)
			secondToLastPoint := child.getPathPoint(node.ID, SecondToLastPt)
			if node.Column < secondToLastPoint.X {
				secondPoint := child.getPathPoint(node.ID, SecondPt)
				childIsSubBranch := child.isPathSubBranch(node.ID)
				if !childIsSubBranch && !secondPoint.Type.IsMergeTo() {
					nextColumn--
				}

				childHasOlderParent := child.hasOlderParent(index, node.Idx)
				if !child.isFirstOfBranch() && !childIsSubBranch && !childHasOlderParent {
					child.setPathColor(node.ID, child.ColorIdx)
				}
				releaseColor(colorsMan, child.getPathColor(node.ID), node.Idx)

				// Insert before the last element
				child.noDupInsert(node.ID, -1, Point{secondToLastPoint.X, node.Idx, MergeBack})

				// Nodes that are following the current node
				for followingNodeID := range followingNodesWithChildrenBeforeIdx.Items {
					followingNode := index.Get(followingNodeID)
					// Following nodes that have a child before the current node
					for _, followingNodeChildID := range followingNode.children {
						followingNodeChild := index.Get(followingNodeChildID)
						// Index to delete is the one before last
						idxRemove := followingNodeChild.pathLength(followingNode.ID) - 1
						if followingNodeChild.Idx < node.Idx &&
							idxRemove >= 0 && !processedNodesInst.HasChild(followingNode.ID, followingNodeChild.ID) {
							// Following node child has a path that is higher than the current path being merged
							targetColumn := followingNodeChild.GetPathHeightAtIdx(followingNode.ID, node.Idx)
							if targetColumn > secondToLastPoint.X {
								// Remove second before last node has same Y, remove the before last node
								for followingNodeChild.getPathPoint(followingNode.ID, idxRemove).Y == followingNodeChild.getPathPoint(followingNode.ID, idxRemove-1).Y {
									followingNodeChild.remove(followingNode.ID, idxRemove-1)
									idxRemove--
								}
								followingNodeChild.remove(followingNode.ID, idxRemove)
								idxRemove--

								// Calculate nb of merging nodes
								nbNodesMergingBack := node.nbNodesMergingBack(index, targetColumn)
								if followingNode.Column > secondToLastPoint.X && !processedNodesInst.HasNode(followingNode.ID) {
									followingNode.Column -= nbNodesMergingBack
								}
								pathPointX := followingNodeChild.getPathPoint(followingNode.ID, idxRemove).X
								followingNodeChild.noDupAppend(followingNode.ID, Point{pathPointX, node.Idx, MergeBack})
								followingNodeChild.noDupAppend(followingNode.ID, Point{pathPointX - nbNodesMergingBack, node.Idx, Pipe})
								followingNodeChild.noDupAppend(followingNode.ID, Point{followingNode.Column, followingNode.Idx, Pipe})
								processedNodesInst.Set(followingNode.ID, followingNodeChild.ID)
							}
						}
					}
				}
			}
		}

		for parentIdx, parentID := range node.Parents {
			parent := index.Get(parentID)

			node.noDupAppend(parent.ID, Point{node.Column, node.Idx, Pipe})

			if !parent.columnDefined() {
				firstParent := index.Get(node.Parents[0])
				column := node.Column
				color := node.ColorIdx
				if parentIdx > 0 && !node.isMergeTo(firstParent.ID) {
					column = incrCol()
					color = getColor(colorsMan, node.Idx)
					node.noDupAppend(parent.ID, Point{column, node.Idx, Fork})
					node.setFirstOfBranch()
				}
				parent.Column = column
				parent.ColorIdx = color
				node.setPathColor(parent.ID, parent.ColorIdx)
			} else if parentIdx == 0 && node.Column < parent.Column {
				for _, childID := range parent.children {
					child := index.Get(childID)
					if idxRemove := child.pathLength(parent.ID) - 1; idxRemove > 0 {
						child.remove(parent.ID, idxRemove)
						child.noDupAppend(parent.ID, Point{child.getPathPoint(parent.ID, idxRemove-1).X, parent.Idx, MergeBack})
						child.noDupAppend(parent.ID, Point{node.Column, parent.Idx, Pipe})
					}
				}
				parent.Column = node.Column
				parent.ColorIdx = node.ColorIdx
				node.setPathColor(parent.ID, node.ColorIdx)
			} else if node.Column < parent.Column {
				node.noDupAppend(parent.ID, Point{parent.Column, node.Idx, Fork})
				node.setPathColor(parent.ID, parent.ColorIdx)
			} else if node.Column > parent.Column {
				if node.hasBiggerParentDefined(index) || (parentIdx == 0 && (parent.Idx > node.Idx+1 || node.firstInBranch(index))) {
					node.noDupAppend(parent.ID, Point{node.Column, parent.Idx, MergeBack})
					node.setPathColor(parent.ID, node.ColorIdx)
				} else {
					node.noDupAppend(parent.ID, Point{parent.Column, node.Idx, MergeTo})
					node.setPathColor(parent.ID, parent.ColorIdx)
				}
			}
			node.noDupAppend(parent.ID, Point{parent.Column, parent.Idx, Pipe})
		}
	}
}

// Get generates the props to turn the input into a graph drawable
func Get(inputNodes []Node) ([]Node, error) {
	nodes, err := BuildTree(inputNodes, NewSimpleColorGen(DefaultColors))
	for _, node := range nodes {
		delete(node, parentsPathsTestKey)
	}
	return nodes, err
}

// GetPaginated same as Get but only return the nodes for the asked page
func GetPaginated(inputNodes []Node, from, size int) ([]Node, error) {
	nodes, err := BuildTree(inputNodes, NewSimpleColorGen(DefaultColors))
	for _, node := range nodes {
		delete(node, parentsPathsTestKey)
	}
	return nodes[from : from+size], err
}

// BuildTree given an array of Node, execute the algorithm on it to generate the necessary properties
// to make it drawable as a graph.
func BuildTree(inputNodes []Node, colorGen IColorGenerator) ([]Node, error) {
	colorsMan := newColorsManager(colorGen)

	nodes := initNodes(inputNodes)
	index := initIndex(nodes)

	initChildren(index, nodes)
	setColumns(index, colorsMan, nodes)

	finalStruct := make([]Node, len(nodes))
	for nodeIdx, node := range nodes {
		finalNode := map[string]any{}
		for key, value := range node.InitialNode {
			finalNode[key] = value
		}
		finalParentsPaths := make([]any, len(node.parentsPaths))
		i := 0
		for _, n := range node.parentsPaths {
			path := make([][]any, len(n.Path))
			for pointIdx, point := range n.Path {
				path[pointIdx] = []any{point.X, point.Y, point.Type}
			}
			finalParentsPaths[i] = []any{colorGen.GetColor(n.ColorIdx), path}
			i++
		}
		finalNode[parentsPathsTestKey] = node.parentsPaths // Kept for tests
		finalNode[gKey] = []any{node.Idx, node.Column, colorGen.GetColor(node.ColorIdx), finalParentsPaths}
		finalStruct[nodeIdx] = finalNode
	}

	return finalStruct, nil
}

// GetInputNodesFromFile creates an array of Node from json contained in a file
func GetInputNodesFromFile(filePath string) (nodes []Node, err error) {
	fileBytes, err := os.ReadFile(filePath)
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

// GetInputNodesFromRepo creates an array of Node from a repository
func GetInputNodesFromRepo(seqIds bool) (nodes []Node, err error) {
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
		var id string
		if seqIds {
			id = strconv.Itoa(ids)
			shaMap[sha] = id
		} else {
			id = sha
		}
		node := Node{}
		node[idKey] = id
		node[parentsKey] = parents
		nodes = append(nodes, node)
		ids++
		if lines[i] != startOfCommit {
			break
		}
	}
	if seqIds {
		for _, node := range nodes {
			mappedParents := make([]string, 0)
			for _, parentSha := range node[parentsKey].([]string) {
				mappedParents = append(mappedParents, shaMap[parentSha])
			}
			node[parentsKey] = mappedParents
		}
	}
	return
}
