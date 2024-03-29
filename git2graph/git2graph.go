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
type color struct {
	releaseIdx int
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
	m map[int]*color
}

func newColorsManager(colorGen IColorGenerator) *colorsManager {
	return &colorsManager{g: colorGen, m: make(map[int]*color)}
}

func (m *colorsManager) getColor(nodeIdx int) int {
	var clr *color
	i := 0
	for {
		var ok bool
		clr, ok = m.m[i]
		if !ok {
			clr = &color{releaseIdx: -1}
			m.m[i] = clr
			break
		}
		if nodeIdx >= clr.releaseIdx && !clr.inUse {
			break
		}
		i++
	}
	clr.inUse = true
	return i
}

func (m *colorsManager) releaseColor(colorIdx int, idx int) {
	for i := range m.m {
		if i == colorIdx {
			color := m.m[i]
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

func (p *Path) isEmpty() bool {
	return p.len() == 0
}

// Return either or not a path is valid (has at least 2 points)
func (p *Path) isValid() bool {
	return p.len() >= 2
}

func (p *Path) setColor(color int) {
	p.colorIdx = color
}

// Return either or not the path is of type "Fork"
func (p *Path) isFork() bool {
	return p.isValid() && p.second().Type.IsFork()
}

// Return either or not the path is of type "MergeTo"
func (p *Path) isMergeTo() bool {
	return p.isValid() && p.second().Type.IsMergeTo()
}

func rotateIdx(idx, length int) int {
	if idx < 0 {
		idx = length + idx
		if idx < 0 {
			log.Fatal("Weird, need to investigate")
		}
	}
	return idx
}

func (p *Path) get(idx int) (out *Point) {
	return p.Points[rotateIdx(idx, p.len())]
}
func (p *Path) first() *Point        { return p.get(0) }
func (p *Path) second() *Point       { return p.get(1) }
func (p *Path) last() *Point         { return p.get(-1) }
func (p *Path) secondToLast() *Point { return p.get(-2) }
func (p *Path) removeLast()          { p.remove(-1) }
func (p *Path) removeSecondToLast()  { p.remove(-2) }
func (p *Path) remove(idx int) {
	idx = rotateIdx(idx, p.len())
	p.Points = append(p.Points[:idx], p.Points[idx+1:]...)
}

// append a point to a path if it is not a duplicate
func (p *Path) noDupAppend(point *Point) {
	if p.isEmpty() || !p.last().Equal(point) {
		p.append(point)
	}
}

// insert a point to a path if it is not a duplicate
func (p *Path) noDupInsert(idx int, point *Point) {
	idx = rotateIdx(idx, p.len())
	if !p.Points[idx-1].Equal(point) {
		p.insert(idx, point)
	}
}

func (p *Path) append(point *Point) {
	p.Points = append(p.Points, point)
}
func (p *Path) insert(idx int, point *Point) {
	p.Points = append(p.Points, &Point{})
	copy(p.Points[idx+1:], p.Points[idx:])
	p.Points[idx] = point
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
	InitialNode   *Node
	ID            string
	Idx           int
	Column        int
	ColorIdx      int
	firstOfBranch bool
	Parents       []*internalNode
	children      []*internalNode
	parentsPaths  map[string]*Path
}

func (n *internalNode) isOrphan() bool {
	return len(n.Parents) == 0
}

// A node is a "firstOfBranch" if there is a path to a parent that needs a new color,
// and the commit is the first commit in that new branch.
func (n *internalNode) isFirstOfBranch() bool {
	return n.firstOfBranch
}

func (n *internalNode) setFirstOfBranch() {
	n.firstOfBranch = true
}

func (n *internalNode) pathTo(parent *internalNode) *Path {
	parentPath, ok := n.parentsPaths[parent.ID]
	if !ok {
		parentPath = &Path{}
		n.parentsPaths[parent.ID] = parentPath
	}
	return parentPath
}

func (n *internalNode) columnDefined() bool {
	return n.Column != -1
}

func (n *internalNode) firstInBranch() bool {
	for _, parentNode := range n.Parents {
		if !parentNode.columnDefined() || parentNode.Column == n.Column {
			return false
		}
	}
	return true
}

func (n *internalNode) hasBiggerParentDefined() bool {
	for _, parentNode := range n.Parents {
		if parentNode.Column > n.Column {
			return true
		}
	}
	return false
}

// Return either or not the node has a parent that has higher "Idx" than the one in parameter
func (n *internalNode) hasOlderParent(idx int) bool {
	for _, parentNode := range n.Parents {
		if parentNode.Idx > idx {
			return true
		}
	}
	return false
}

// A subbranch, is when the child node is in the middle of another branch
// See test_022.png node #4 (zero-indexed)
func (n *internalNode) isPathSubBranch(parent *internalNode) bool {
	return n.pathTo(parent).isFork() && !n.isFirstOfBranch()
}

const (
	idKey               = "id"
	parentsKey          = "parents"
	gKey                = "g"
	parentsPathsTestKey = "parentsPaths"
)

// GetPathHeightAtIdx Get the path X at Idx
func (n *internalNode) GetPathHeightAtIdx(parent *internalNode, lookupIdx int) (height int) {
	height = -1
	parentPath := n.pathTo(parent)
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
func (n *internalNode) nbNodesMergingBack(maxX int) (nbNodesMergingBack int) {
	for _, child := range n.children {
		path := child.pathTo(n)
		childIsSubBranch := child.isPathSubBranch(n)
		secondToLastPoint := path.secondToLast()
		if n.Column < secondToLastPoint.X && secondToLastPoint.X < maxX &&
			!childIsSubBranch &&
			!path.isMergeTo() {
			nbNodesMergingBack++
		}
	}
	return
}

// SerializeOutput Json encode object
func SerializeOutput(out []*Node) {
	if !NoOutput {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(out); err != nil {
			log.Error("Could not encode json")
		}
	}
}

// GetInputNodesFromJSON Get nodes from json object
func GetInputNodesFromJSON(inputJSON []byte) (nodes []*Node, err error) {
	dec := json.NewDecoder(bytes.NewReader(inputJSON))
	err = dec.Decode(&nodes)
	if err != nil {
		return
	}
	for _, node := range nodes {
		parents := make([]string, 0)
		nodeParents, ok := (*node)[parentsKey]
		if !ok {
			log.Fatal("malformed json input, node missing parents property")
		}
		for _, parent := range nodeParents.([]any) {
			parents = append(parents, parent.(string))
		}
		(*node)[parentsKey] = parents
	}
	return
}

func initNodes(inputNodes []*Node) []*internalNode {
	out := make([]*internalNode, len(inputNodes))
	index := make(map[string]*internalNode)
	for idx, node := range inputNodes {
		id, ok := (*node)[idKey].(string)
		if !ok {
			log.Fatal("id property must be a string")
		}
		newNode := &internalNode{}
		newNode.InitialNode = node
		newNode.ID = id
		newNode.Idx = idx
		newNode.Column = -1
		newNode.Parents = make([]*internalNode, 0)
		newNode.parentsPaths = make(map[string]*Path)
		newNode.children = make([]*internalNode, 0)
		out[idx] = newNode
		index[newNode.ID] = newNode
	}
	for _, node := range out {
		parents, ok := (*node.InitialNode)[parentsKey].([]string)
		if !ok {
			log.Fatal("parents property must be an array of string")
		}
		for _, parent := range parents {
			node.Parents = append(node.Parents, index[parent])
		}
	}
	return out
}

type stringSet map[*internalNode]struct{}

func newStringSet() stringSet {
	return make(map[*internalNode]struct{})
}

func (s *stringSet) Add(ins []*internalNode) {
	for _, in := range ins {
		(*s)[in] = struct{}{}
	}
}

func (s *stringSet) Remove(in *internalNode) {
	delete(*s, in)
}

func setColumns(colorsMan *colorsManager, nodes []*internalNode) {
	followingNodesWithChildrenBeforeIdx := newStringSet()
	nextColumn := -1
	incrCol := func() int {
		nextColumn++
		return nextColumn
	}
	for _, node := range nodes {
		// Add node as a child into parents
		for _, parent := range node.Parents {
			parent.children = append(parent.children, node)
		}

		// Set column if not defined
		if !node.columnDefined() {
			node.Column = incrCol()
			node.ColorIdx = colorsMan.getColor(node.Idx)
		}

		// Cache the following node with child before the current node
		followingNodesWithChildrenBeforeIdx.Add(node.Parents)
		followingNodesWithChildrenBeforeIdx.Remove(node)

		// Each child that are merging
		// For each node, we need to check each child.
		// For each child that is merging back, we need to alter paths that are passing over
		// and decrement their column.
		processedNodes := make(map[*internalNode]bool)
		for _, child := range node.children {
			pathToNode := child.pathTo(node)
			secondToLastPoint := pathToNode.secondToLast()
			if node.Column < secondToLastPoint.X || node.isOrphan() {
				childIsSubBranch := child.isPathSubBranch(node)
				if !childIsSubBranch && !pathToNode.isMergeTo() {
					nextColumn--
				}

				childHasOlderParent := child.hasOlderParent(node.Idx)
				if !child.isFirstOfBranch() && !childIsSubBranch && !childHasOlderParent {
					pathToNode.setColor(child.ColorIdx)
				}
				colorsMan.releaseColor(pathToNode.colorIdx, node.Idx)

				// Insert before the last element
				if node.Column != child.Column {
					pathToNode.noDupInsert(-1, &Point{secondToLastPoint.X, node.Idx, MergeBack})
				}

				// Nodes that are following the current node
				for followingNode := range followingNodesWithChildrenBeforeIdx {
					// Following nodes that have a child before the current node
					for _, followingNodeChild := range followingNode.children {
						pathToFollowingNode := followingNodeChild.pathTo(followingNode)
						if followingNodeChild.Idx < node.Idx &&
							!pathToFollowingNode.isEmpty() && !processedNodes[followingNodeChild] {
							// Following node child has a path that is higher than the current path being merged
							targetColumn := followingNodeChild.GetPathHeightAtIdx(followingNode, node.Idx)
							if targetColumn > secondToLastPoint.X {
								// Remove all nodes, that are next to the last node, that have the same Y as the last node
								for pathToFollowingNode.last().Y == pathToFollowingNode.secondToLast().Y {
									pathToFollowingNode.removeSecondToLast()
								}
								pathToFollowingNode.removeLast()

								// Calculate nb of merging nodes
								nbNodesMergingBack := 0
								nodeForMerge := node
								if node.isOrphan() {
									nodeForMerge = nodes[node.Idx+1]
									nbNodesMergingBack++
								}
								nbNodesMergingBack += nodeForMerge.nbNodesMergingBack(targetColumn)
								shouldMoveNode := followingNode.Column > secondToLastPoint.X && !processedNodes[followingNode]
								if shouldMoveNode {
									followingNode.Column -= nbNodesMergingBack
								}
								pathPointX := pathToFollowingNode.last().X
								pathToFollowingNode.noDupAppend(&Point{pathPointX, nodeForMerge.Idx, MergeBack})
								pathToFollowingNode.noDupAppend(&Point{pathPointX - nbNodesMergingBack, nodeForMerge.Idx, Pipe})
								pathToFollowingNode.noDupAppend(&Point{followingNode.Column, followingNode.Idx, Pipe})
								if shouldMoveNode {
									// If we move the node, we need to ensure that all paths going to that node now goes to the new column
									fixPathsToNode(followingNode)
									processedNodes[followingNode] = true
								}
								processedNodes[followingNodeChild] = true
							}
						}
					}
				}
			}
		}

		for parentIdx, parent := range node.Parents {
			nodePathToParent := node.pathTo(parent)
			nodePathToParent.noDupAppend(&Point{node.Column, node.Idx, Pipe})
			if !parent.columnDefined() {
				if parentIdx > 0 && !node.pathTo(node.Parents[0]).isMergeTo() {
					parent.Column = incrCol()
					parent.ColorIdx = colorsMan.getColor(node.Idx)
					nodePathToParent.noDupAppend(&Point{parent.Column, node.Idx, Fork})
					node.setFirstOfBranch()
				} else {
					parent.Column = node.Column
					parent.ColorIdx = node.ColorIdx
				}
				nodePathToParent.setColor(parent.ColorIdx)
			} else if node.Column < parent.Column {
				if parentIdx == 0 {
					for _, child := range parent.children {
						pathToParent := child.pathTo(parent)
						if pathToParent.isValid() {
							pathToParent.removeLast()
							pathToParent.noDupAppend(&Point{pathToParent.last().X, parent.Idx, MergeBack})
							pathToParent.noDupAppend(&Point{node.Column, parent.Idx, Pipe})
						}
					}
					parent.Column = node.Column
					parent.ColorIdx = node.ColorIdx
					nodePathToParent.setColor(node.ColorIdx)
				} else {
					nodePathToParent.noDupAppend(&Point{parent.Column, node.Idx, Fork})
					nodePathToParent.setColor(parent.ColorIdx)
				}
			} else if node.Column > parent.Column {
				if node.hasBiggerParentDefined() || (parentIdx == 0 && (parent.Idx > node.Idx+1 || node.firstInBranch())) {
					nodePathToParent.noDupAppend(&Point{node.Column, parent.Idx, MergeBack})
					nodePathToParent.setColor(node.ColorIdx)
				} else {
					nodePathToParent.noDupAppend(&Point{parent.Column, node.Idx, MergeTo})
					nodePathToParent.setColor(parent.ColorIdx)
				}
			}
			nodePathToParent.noDupAppend(&Point{parent.Column, parent.Idx, Pipe})
		}
	}
}

func fixPathsToNode(node *internalNode) {
	for _, child := range node.children {
		path := child.pathTo(node)
		if !path.isEmpty() {
			path.last().X = node.Column
		}
	}
}

// Get generates the props to turn the input into a graph drawable
func Get(inputNodes []*Node) ([]*Node, error) {
	nodes, err := BuildTree(inputNodes, NewCycleColorGen(DefaultColors))
	for _, node := range nodes {
		delete(*node, parentsPathsTestKey)
	}
	return nodes, err
}

// GetPaginated same as Get but only return the nodes for the asked page
func GetPaginated(inputNodes []*Node, from, size int) ([]*Node, error) {
	nodes, err := BuildTree(inputNodes, NewCycleColorGen(DefaultColors))
	for _, node := range nodes {
		delete(*node, parentsPathsTestKey)
	}
	return nodes[from : from+size], err
}

// BuildTree given an array of Node, execute the algorithm on it to generate the necessary properties
// to make it drawable as a graph.
func BuildTree(inputNodes []*Node, colorGen IColorGenerator) ([]*Node, error) {
	colorsMan := newColorsManager(colorGen)

	nodes := initNodes(inputNodes)

	setColumns(colorsMan, nodes)

	finalStruct := make([]*Node, len(nodes))
	for nodeIdx, node := range nodes {
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
		finalNode := node.InitialNode
		(*finalNode)[parentsPathsTestKey] = node.parentsPaths // Kept for tests
		(*finalNode)[gKey] = []any{node.Idx, node.Column, colorGen.GetColor(node.ColorIdx), finalParentsPaths}
		finalStruct[nodeIdx] = finalNode
	}

	return finalStruct, nil
}

// GetInputNodesFromFile creates an array of Node from json contained in a file
func GetInputNodesFromFile(filePath string) (nodes []*Node, err error) {
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
func GetInputNodesFromRepo(seqIds bool, parentsOf string) (nodes []*Node, err error) {
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
		node := &Node{}
		(*node)[idKey] = id
		(*node)[parentsKey] = parents
		nodes = append(nodes, node)
		ids++
		if lines[i] != startOfCommit {
			break
		}
	}
	if seqIds {
		for _, node := range nodes {
			mappedParents := make([]string, 0)
			for _, parentSha := range (*node)[parentsKey].([]string) {
				mappedParents = append(mappedParents, shaMap[parentSha])
			}
			(*node)[parentsKey] = mappedParents
		}
	}
	return
}
