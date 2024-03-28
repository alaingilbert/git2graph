package git2graph

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// CycleColorGen ...
type CycleColorGen struct {
	colors []string
}

// NewCycleColorGen creates a new CycleColorGen
func NewCycleColorGen(colors []string) *CycleColorGen {
	return &CycleColorGen{colors: colors}
}

func (c *CycleColorGen) GetColor(idx int) string {
	return c.colors[idx%len(c.colors)]
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
func (p pointType) IsFork() bool    { return p == Fork }

// Node is the raw information for a commit
type Node map[string]any

// Path defines how to draw a line in between a parent and child nodes
type Path struct {
	Points   []*Point
	colorIdx int
}

// Return the current length of the path
func (p *Path) len() int {
	return len(p.Points)
}

// Return either or not a path is valid (has at least 2 points)
func (p *Path) isValid() bool {
	return p.len() >= 2
}

// Return either or not the path is of type "Fork"
func (p *Path) isFork() bool {
	return p.isValid() && p.second().Type.IsFork()
}

// Return either or not the path is of type "MergeTo"
func (p *Path) isMergeTo() bool {
	return p.isValid() && p.second().Type.IsMergeTo()
}

func (p *Path) get(idx int) (out *Point) {
	if idx < 0 {
		idx = p.len() + idx
		if idx < 0 {
			log.Fatal("Weird, need to investigate")
		}
	}
	return p.Points[idx]
}
func (p *Path) first() *Point        { return p.get(0) }
func (p *Path) second() *Point       { return p.get(1) }
func (p *Path) last() *Point         { return p.get(-1) }
func (p *Path) secondToLast() *Point { return p.get(-2) }
func (p *Path) removeLast()          { p.remove(p.len() - 1) }
func (p *Path) removeSecondToLast()  { p.remove(p.len() - 2) }
func (p *Path) remove(idx int) {
	p.Points = append(p.Points[:idx], p.Points[idx+1:]...)
}

// Point is one part of a path
type Point struct {
	X    int
	Y    int
	Type pointType
}

func (p *Point) String() string {
	return fmt.Sprintf("{%d,%d,%d}", p.X, p.Y, p.Type)
}

func (p *Point) Equal(other *Point) bool {
	return p.X == other.X && p.Y == other.Y && p.Type == other.Type
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
func (n *internalNode) noDupAppend(parentID string, point *Point) {
	parentPath := n.pathTo(parentID)
	parentPathLen := parentPath.len()
	if parentPathLen > 0 && parentPath.Points[parentPathLen-1].Equal(point) {
		return
	}
	n.append(parentID, point)
}

// insert a point to a parent path if it is not a duplicate
func (n *internalNode) noDupInsert(parentID string, idx int, point *Point) {
	parentPath := n.pathTo(parentID)
	if idx < 0 {
		idx = parentPath.len() + idx
	}
	if parentPath.Points[idx-1].Equal(point) {
		return
	}
	n.insert(parentID, idx, point)
}

func (n *internalNode) append(parentID string, point *Point) {
	parentPath := n.pathTo(parentID)
	parentPath.Points = append(parentPath.Points, point)
}

func (n *internalNode) insert(parentID string, idx int, point *Point) {
	parentPath := n.pathTo(parentID)
	parentPath.Points = append(parentPath.Points, &Point{})
	copy(parentPath.Points[idx+1:], parentPath.Points[idx:])
	parentPath.Points[idx] = point
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
	n.pathTo(parentID).colorIdx = color
}

func (n *internalNode) getPathColor(parentID string) int {
	return n.pathTo(parentID).colorIdx
}

// A subbranch, is when the child node is in the middle of another branch
// See test_022.png node #4 (zero-indexed)
func (n *internalNode) isPathSubBranch(parentID string) bool {
	return n.pathTo(parentID).isFork() && !n.isFirstOfBranch()
}

const (
	idKey               = "id"
	parentsKey          = "parents"
	gKey                = "g"
	parentsPathsTestKey = "parentsPaths"
)

// GetPathHeightAtIdx Get the path X at Idx
func (n *internalNode) GetPathHeightAtIdx(parentID string, lookupIdx int) (height int) {
	height = -1
	parentPath := n.pathTo(parentID)
	firstPoint := parentPath.first()
	lastPoint := parentPath.last()
	if lookupIdx < firstPoint.Y || lookupIdx > lastPoint.Y {
		return
	}
	for _, point := range parentPath.Points {
		if point.Y <= lookupIdx {
			height = point.X
		}
	}
	return
}

// A merging node is one that come from a higher column, but is not a sub-branch and is not a MergeTo
func (n *internalNode) nbNodesMergingBack(index *nodesCache, maxX int) (nbNodesMergingBack int) {
	nodeID := n.ID
	for _, childID := range n.children {
		child := index.Get(childID)
		childIsSubBranch := child.isPathSubBranch(nodeID)
		secondToLastPoint := child.pathTo(nodeID).secondToLast()
		if n.Column < secondToLastPoint.X && secondToLastPoint.X < maxX &&
			!childIsSubBranch &&
			!child.pathTo(nodeID).isMergeTo() {
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
	Items map[string]struct{}
}

func newStringSet() stringSet {
	s := stringSet{}
	s.Items = make(map[string]struct{})
	return s
}

func (s *stringSet) Add(ins []string) {
	for _, in := range ins {
		s.Items[in] = struct{}{}
	}
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
		followingNodesWithChildrenBeforeIdx.Add(node.Parents)
		followingNodesWithChildrenBeforeIdx.Remove(node.ID)

		// Each child that are merging
		// For each node, we need to check each child.
		// For each child that is merging back, we need to alter paths that are passing over
		// and decrement their column.
		processedNodesInst := newProcessedNodes()
		processedNodesInst1 := newProcessedNodes()
		for _, childID := range node.children {
			child := index.Get(childID)
			pathToNode := child.pathTo(node.ID)
			secondToLastPoint := pathToNode.secondToLast()
			if node.Column < secondToLastPoint.X || len(node.Parents) == 0 {
				childIsSubBranch := child.isPathSubBranch(node.ID)
				if !childIsSubBranch && !pathToNode.isMergeTo() {
					nextColumn--
				}

				childHasOlderParent := child.hasOlderParent(index, node.Idx)
				if !child.isFirstOfBranch() && !childIsSubBranch && !childHasOlderParent {
					child.setPathColor(node.ID, child.ColorIdx)
				}
				releaseColor(colorsMan, child.getPathColor(node.ID), node.Idx)

				// Insert before the last element
				if node.Column != child.Column {
					child.noDupInsert(node.ID, -1, &Point{secondToLastPoint.X, node.Idx, MergeBack})
				}

				// Nodes that are following the current node
				for followingNodeID := range followingNodesWithChildrenBeforeIdx.Items {
					followingNode := index.Get(followingNodeID)
					// Following nodes that have a child before the current node
					for _, followingNodeChildID := range followingNode.children {
						followingNodeChild := index.Get(followingNodeChildID)
						pathToFollowingNode := followingNodeChild.pathTo(followingNode.ID)
						if followingNodeChild.Idx < node.Idx &&
							pathToFollowingNode.len() > 0 && !processedNodesInst.HasChild(followingNode.ID, followingNodeChild.ID) {
							// Following node child has a path that is higher than the current path being merged
							targetColumn := followingNodeChild.GetPathHeightAtIdx(followingNode.ID, node.Idx)
							if targetColumn > secondToLastPoint.X {
								// Remove second before last node has same Y, remove the before last node
								for pathToFollowingNode.last().Y == pathToFollowingNode.secondToLast().Y {
									pathToFollowingNode.removeSecondToLast()
								}
								pathToFollowingNode.removeLast()

								// Calculate nb of merging nodes
								nbNodesMergingBack := 0
								y := node.Idx
								if len(node.Parents) == 0 {
									y++
									nbNodesMergingBack++
								}
								nbNodesMergingBack += nodes[y].nbNodesMergingBack(index, targetColumn)
								shouldMoveNode := followingNode.Column > secondToLastPoint.X && !processedNodesInst1.HasNode(followingNode.ID)
								if shouldMoveNode {
									followingNode.Column -= nbNodesMergingBack
								}
								pathPointX := pathToFollowingNode.last().X
								followingNodeChild.noDupAppend(followingNode.ID, &Point{pathPointX, y, MergeBack})
								followingNodeChild.noDupAppend(followingNode.ID, &Point{pathPointX - nbNodesMergingBack, y, Pipe})
								followingNodeChild.noDupAppend(followingNode.ID, &Point{followingNode.Column, followingNode.Idx, Pipe})
								if shouldMoveNode {
									// If we move the node, we need to ensure that all paths going to that node now goes to the new column
									for _, c := range followingNode.children {
										path := index.Get(c).pathTo(followingNode.ID)
										if path.len() > 0 {
											path.last().X = followingNode.Column
										}
									}
									processedNodesInst1.Set(followingNode.ID, "")
								}
								processedNodesInst.Set(followingNode.ID, followingNodeChild.ID)
							}
						}
					}
				}
			}
		}

		for parentIdx, parentID := range node.Parents {
			parent := index.Get(parentID)
			node.noDupAppend(parent.ID, &Point{node.Column, node.Idx, Pipe})
			if !parent.columnDefined() {
				if parentIdx > 0 && !node.pathTo(node.Parents[0]).isMergeTo() {
					parent.Column = incrCol()
					parent.ColorIdx = getColor(colorsMan, node.Idx)
					node.noDupAppend(parent.ID, &Point{parent.Column, node.Idx, Fork})
					node.setFirstOfBranch()
				} else {
					parent.Column = node.Column
					parent.ColorIdx = node.ColorIdx
				}
				node.setPathColor(parent.ID, parent.ColorIdx)
			} else if node.Column < parent.Column {
				if parentIdx == 0 {
					for _, childID := range parent.children {
						child := index.Get(childID)
						pathToParent := child.pathTo(parent.ID)
						if idxRemove := pathToParent.len() - 1; idxRemove > 0 {
							pathToParent.remove(idxRemove)
							child.noDupAppend(parent.ID, &Point{pathToParent.get(idxRemove - 1).X, parent.Idx, MergeBack})
							child.noDupAppend(parent.ID, &Point{node.Column, parent.Idx, Pipe})
						}
					}
					parent.Column = node.Column
					parent.ColorIdx = node.ColorIdx
					node.setPathColor(parent.ID, node.ColorIdx)
				} else {
					node.noDupAppend(parent.ID, &Point{parent.Column, node.Idx, Fork})
					node.setPathColor(parent.ID, parent.ColorIdx)
				}
			} else if node.Column > parent.Column {
				if node.hasBiggerParentDefined(index) || (parentIdx == 0 && (parent.Idx > node.Idx+1 || node.firstInBranch(index))) {
					node.noDupAppend(parent.ID, &Point{node.Column, parent.Idx, MergeBack})
					node.setPathColor(parent.ID, node.ColorIdx)
				} else {
					node.noDupAppend(parent.ID, &Point{parent.Column, node.Idx, MergeTo})
					node.setPathColor(parent.ID, parent.ColorIdx)
				}
			}
			node.noDupAppend(parent.ID, &Point{parent.Column, parent.Idx, Pipe})
		}
	}
}

// Get generates the props to turn the input into a graph drawable
func Get(inputNodes []Node) ([]Node, error) {
	nodes, err := BuildTree(inputNodes, NewCycleColorGen(DefaultColors))
	for _, node := range nodes {
		delete(node, parentsPathsTestKey)
	}
	return nodes, err
}

// GetPaginated same as Get but only return the nodes for the asked page
func GetPaginated(inputNodes []Node, from, size int) ([]Node, error) {
	nodes, err := BuildTree(inputNodes, NewCycleColorGen(DefaultColors))
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
			path := make([][]any, len(n.Points))
			for pointIdx, point := range n.Points {
				path[pointIdx] = []any{point.X, point.Y, point.Type}
			}
			finalParentsPaths[i] = []any{colorGen.GetColor(n.colorIdx), path}
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
func GetInputNodesFromRepo(seqIds bool, parentsOf string) (nodes []Node, err error) {
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
		if parentsOf != "" && strings.HasPrefix(sha, parentsOf) {
			log.Errorf("%v: %v", sha, parents)
			os.Exit(0)
		}
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
