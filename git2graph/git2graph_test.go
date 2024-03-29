package git2graph

import (
	"strings"
	"testing"
)

func validateColumns(t *testing.T, expectedColumns []int, data []*Node) {
	for idx, row := range data {
		expectedColumn := expectedColumns[idx]
		actualColumn := (*row)[gKey].([]any)[1]
		nodeID := (*row)[idKey]
		if actualColumn != expectedColumn {
			t.Fail()
			t.Logf("ID: %s, Expected column: %d, Actual column: %d", nodeID, expectedColumn, actualColumn)
		}
	}
}

func pprintPoints(points []*Point) string {
	s := make([]string, 0)
	for _, p := range points {
		s = append(s, p.String())
	}
	return "[" + strings.Join(s, ",") + "]"
}

func validatePaths(t *testing.T, expectedPaths []map[string]Path, data []*Node) {
	for nodeIdx, node := range data {
		for _, parentID := range (*node)[parentsKey].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			expectedPath := expectedPaths[nodeIdx][parentID].Points
			parentPath := (*node)[parentsPathsTestKey].(map[string]*Path)[parentID]
			nodeID := (*node)[idKey]
			if len(parentPath.Points) != len(expectedPath) {
				t.Fail()
				t.Logf("ID: %s, Expected nb paths: %d, Actual nb paths: %d", nodeID, len(expectedPath), len(parentPath.Points))
				t.Logf("ID: %s, Expected vs Actual:\n%v\n%v", nodeID, pprintPoints(expectedPath), pprintPoints(parentPath.Points))
				return
			}
			for pathIdx, pathNode := range parentPath.Points {
				if !pathNode.Equal(expectedPath[pathIdx]) {
					t.Fail()
					t.Logf("ID: %s, Expected path: %d, Actual path: %d", nodeID, expectedPath[pathIdx], pathNode)
					t.Logf("ID: %s, Expected vs Actual:\n%v\n%v", nodeID, pprintPoints(expectedPath), pprintPoints(parentPath.Points))
				}
			}
		}
	}
}

func validateColors(t *testing.T, expectedPaths []map[string]Path, data []*Node) {
	for nodeIdx, node := range data {
		nodeID := (*node)[idKey]
		for _, parentID := range (*node)[parentsKey].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			parentPath := (*node)[parentsPathsTestKey].(map[string]*Path)[parentID]
			expectedPath := expectedPaths[nodeIdx][parentID]
			if expectedPath.colorIdx != parentPath.colorIdx {
				t.Logf("ID: %s, Expected: %v, Actual: %v", nodeID, expectedPath.colorIdx, parentPath.colorIdx)
				t.Fail()
			}
		}
	}
}

var customColors = NewSimpleColorGen([]string{
	"color1",
	"color2",
	"color3",
	"color4",
	"color5",
	"color6",
	"color7",
	"color8",
	"color9",
	"color10",
})

func TestNotEnoughColors(t *testing.T) {
	var colors = NewSimpleColorGen([]string{
		"color1",
		"color2",
	})
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"3"}},
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{}},
	}
	out, _ := BuildTree(inputNodes, colors)
	if (*out[2])[gKey].([]any)[2] != "#000" {
		t.Fail()
	}
}

func TestGetInputNodesFromJson(t *testing.T) {
	json := `[{"id": "1", "parents": ["2"]}, {"id": "2", "parents": ["3"]}, {"id": "3", "parents": []}]`
	inputNodes, _ := GetInputNodesFromJSON([]byte(json))
	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]Path{
		{"2": Path{[]*Point{{0, 0, 0}, {0, 1, 0}}, 0}},
		{"3": Path{[]*Point{{0, 1, 0}, {0, 2, 0}}, 0}},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func TestGetInputNodesFromJsonWithBadJson(t *testing.T) {
	json := `[{"id": "1", "parents": ["2"]}, {"id": "2", "parents": ["3"]}, {"id": "3", "parents": []}`
	_, err := GetInputNodesFromJSON([]byte(json))
	if err == nil {
		t.Fail()
	}
}

// 1
// |
// 2
// |
// 3
func Test1(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"2": {[]*Point{{0, 0, 0}, {0, 1, 0}}, 0},
		},
		{
			"3": {[]*Point{{0, 1, 0}, {0, 2, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// | 2
// |/
// 3
func Test2(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// |/
// 3
func Test3(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// 3 |
// |\|
// | |\
// | | 4
// | |/
// |/
// 5
func Test4(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"5", "4"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*Point{{0, 2, 0}, {0, 4, 0}}, 0},
			"4": {[]*Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*Point{{2, 3, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// | 2
// | | 3
// | |/
// |/
// 4
// | 5
// |/
// 6
func Test5(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"4"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, 1},
		},
		{
			"4": {[]*Point{{2, 2, 0}, {2, 3, 1}, {0, 3, 0}}, 2},
		},
		{
			"6": {[]*Point{{0, 3, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 3},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// |/|
// 3 |
// | 4
// |/
// 5
func Test6(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"3", "4"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, 0},
			"4": {[]*Point{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*Point{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 3, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// 3 |\
// | 4 |
// | |/
// |/|
// 5 |
// |/
// 6
func Test7(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"4", "5"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 3, 0}}, 1},
			"5": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"5": {[]*Point{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, 1},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// 3 |
// |\|
// | |\
// |/ /
// 4 |
// | 5
// |/
// 6
func Test8(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"4"}},
		{"id": "3", "parents": []string{"4", "5"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, 1},
		},
		{
			"4": {[]*Point{{0, 2, 0}, {0, 3, 0}}, 0},
			"5": {[]*Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 1}, {1, 3, 0}, {1, 4, 0}}, 2},
		},
		{
			"6": {[]*Point{{0, 3, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// 3 |
// |\|
// | |\
// 4 | |
// |\| |
// | |\|
// | | |\
// |/ / /
// 5 | |
// | | 6
// | 7 |
// | |/
// |/
// 8
func Test9(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3", "2"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"4", "7"}},
		{"id": "4", "parents": []string{"5", "6"}},
		{"id": "5", "parents": []string{"8"}},
		{"id": "6", "parents": []string{"8"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*Point{{0, 2, 0}, {0, 3, 0}}, 0},
			"7": {[]*Point{{0, 2, 0}, {2, 2, 2}, {2, 4, 1}, {1, 4, 0}, {1, 6, 0}}, 2},
		},
		{
			"5": {[]*Point{{0, 3, 0}, {0, 4, 0}}, 0},
			"6": {[]*Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, 3},
		},
		{
			"8": {[]*Point{{0, 4, 0}, {0, 7, 0}}, 0},
		},
		{
			"8": {[]*Point{{2, 5, 0}, {2, 7, 1}, {0, 7, 0}}, 3},
		},
		{
			"8": {[]*Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// | |\
// | | 3
// 4 | |
// |\| |
// | |\|
// | | |\
// | |/ /
// | 5 |
// | |\|
// | | 6
// | |/
// | 7
// |/
// 8
func Test10(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"4", "2"}},
		{"id": "2", "parents": []string{"5", "3"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"8", "6"}},
		{"id": "5", "parents": []string{"7", "6"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
			"3": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
		},
		{
			"8": {[]*Point{{0, 3, 0}, {0, 7, 0}}, 0},
			"6": {[]*Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, 3},
		},
		{
			"6": {[]*Point{{1, 4, 0}, {2, 4, 2}, {2, 5, 0}}, 3},
			"7": {[]*Point{{1, 4, 0}, {1, 6, 0}}, 1},
		},
		{
			"7": {[]*Point{{2, 5, 0}, {2, 6, 1}, {1, 6, 0}}, 3},
		},
		{
			"8": {[]*Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |-2
// 3 |
// |-4
// 5 |
// |/
// 6
func Test11(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"3", "4"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"5", "6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, 0},
			"4": {[]*Point{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*Point{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 3, 0}, {0, 3, 3}, {0, 4, 0}}, 0},
			"6": {[]*Point{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, 1},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |-2
// 3 |
// |\|
// | |\
// | | 4
// | |/
// |/|
// 5 |
// | 6
// |/
// 7
func Test12(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"3", "6"}},
		{"id": "3", "parents": []string{"5", "4"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"7"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// |\
// | 2
// | |
// | 3
// 4 |\
// |\| |
// | 5 |
// | |\|
// | 6 |\
// 7 | | |
// | | |/
// | |/|
// |/| |
// 8 | |
// | | 9
// | |/
// |/
// 10
func Test13(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"4", "2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{"5", "9"}},
		{"id": "4", "parents": []string{"7", "5"}},
		{"id": "5", "parents": []string{"6", "8"}},
		{"id": "6", "parents": []string{"10"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{"10"}},
		{"id": "9", "parents": []string{"10"}},
		{"id": "10", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 1, 0, 1, 1, 0, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
			"2": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// | 2
// 3 |\
// |\| |
// | |\|
// | | 4
// | |/
// | 5
// 6 |\
// |\| |
// | |\|
// | | 7
// | |/
// |/
// 8
func Test14(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"5", "4"}},
		{"id": "3", "parents": []string{"6", "4"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"8", "7"}},
		{"id": "6", "parents": []string{"8", "7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
			"4": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, 2},
		},
		{
			"6": {[]*Point{{0, 2, 0}, {0, 5, 0}}, 0},
			"4": {[]*Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*Point{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
		},
		{
			"8": {[]*Point{{1, 4, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
			"7": {[]*Point{{1, 4, 0}, {2, 4, 2}, {2, 6, 0}}, 3},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// | 2
// 3 |\
// |\| |
// | |\|
// | | 4
// 5 | |
// |\| |
// | |\|
// | | |\
// | |/ /
// | 6 |
// | |\|
// | | 7
// | |/
// |/
// 8
func Test15(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"6", "4"}},
		{"id": "3", "parents": []string{"5", "4"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"8", "7"}},
		{"id": "6", "parents": []string{"8", "7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
			"4": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*Point{{0, 2, 0}, {0, 4, 0}}, 0},
			"4": {[]*Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"6": {[]*Point{{2, 3, 0}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"8": {[]*Point{{0, 4, 0}, {0, 7, 0}}, 0},
			"7": {[]*Point{{0, 4, 0}, {3, 4, 2}, {3, 5, 1}, {2, 5, 0}, {2, 6, 0}}, 3},
		},
		{
			"8": {[]*Point{{1, 5, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
			"7": {[]*Point{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, 3},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// 1
// | 2
// | | 3
// | |/ 4
// |/| /
// 5 |/
// | 6
// |/
// 7
func Test16(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"6"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"7"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"6": {[]*Point{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"7": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"7": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test17(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"4"}},
		{"id": "3", "parents": []string{"6"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"6": {[]*Point{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 6, 1}, {0, 6, 0}}, 3},
		},
		{
			"5": {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*Point{{0, 5, 0}, {0, 6, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test18(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"5"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 2},
		},
		{
			"5": {[]*Point{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test19(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"9"}},
		{"id": "3", "parents": []string{"7"}},
		{"id": "4", "parents": []string{"11", "6"}},
		{"id": "5", "parents": []string{"8", "6"}},
		{"id": "6", "parents": []string{"11"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{"10"}},
		{"id": "9", "parents": []string{"10"}},
		{"id": "10", "parents": []string{"12"}},
		{"id": "11", "parents": []string{"12"}},
		{"id": "12", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 4, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": {[]*Point{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"9": {[]*Point{{2, 2, 0}, {2, 9, 0}}, 2},
		},
		{
			"7": {[]*Point{{3, 3, 0}, {3, 7, 0}}, 3},
		},
		{
			"11": {[]*Point{{1, 4, 0}, {1, 11, 0}}, 1},
			"6":  {[]*Point{{1, 4, 0}, {4, 4, 2}, {4, 6, 0}}, 4},
		},
		{
			"8": {[]*Point{{0, 5, 0}, {0, 8, 0}}, 0},
			"6": {[]*Point{{0, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"11": {[]*Point{{4, 6, 0}, {4, 8, 1}, {3, 8, 0}, {3, 10, 1}, {2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, 4},
		},
		{
			"8": {[]*Point{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, 3},
		},
		{
			"10": {[]*Point{{0, 8, 0}, {0, 10, 0}}, 0},
		},
		{
			"10": {[]*Point{{2, 9, 0}, {2, 10, 1}, {0, 10, 0}}, 2},
		},
		{
			"12": {[]*Point{{0, 10, 0}, {0, 12, 0}}, 0},
		},
		{
			"12": {[]*Point{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test20(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"4": {[]*Point{{3, 3, 0}, {3, 4, 1}, {0, 4, 0}}, 3},
		},
		{
			"5": {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test21(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"6", "5"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"7"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"6": {[]*Point{{1, 3, 0}, {1, 6, 0}}, 1},
			"5": {[]*Point{{1, 3, 0}, {2, 3, 2}, {2, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"5": {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
		},
		{
			"7": {[]*Point{{0, 5, 0}, {0, 7, 0}}, 0},
		},
		{
			"7": {[]*Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test22(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"6"}},
		{"id": "3", "parents": []string{"7"}},
		{"id": "4", "parents": []string{"7", "6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{"8"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": {[]*Point{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"6": {[]*Point{{2, 2, 0}, {2, 6, 1}, {0, 6, 0}}, 2},
		},
		{
			"7": {[]*Point{{3, 3, 0}, {3, 6, 1}, {2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, 3},
		},
		{
			"7": {[]*Point{{1, 4, 0}, {1, 7, 0}}, 1},
			"6": {[]*Point{{1, 4, 0}, {2, 4, 2}, {2, 6, 1}, {0, 6, 0}}, 2},
		},
		{
			"6": {[]*Point{{0, 5, 0}, {0, 6, 0}}, 0},
		},
		{
			"8": {[]*Point{{0, 6, 0}, {0, 8, 0}}, 0},
		},
		{
			"8": {[]*Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test23(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"4"}},
		{"id": "3", "parents": []string{"7"}},
		{"id": "4", "parents": []string{"6", "5"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{"8"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"7": {[]*Point{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 7, 0}}, 3},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
			"5": {[]*Point{{0, 4, 0}, {2, 4, 2}, {2, 5, 0}}, 4},
		},
		{
			"6": {[]*Point{{2, 5, 0}, {2, 6, 1}, {0, 6, 0}}, 4},
		},
		{
			"8": {[]*Point{{0, 6, 0}, {0, 8, 0}}, 0},
		},
		{
			"8": {[]*Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 3},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test24(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"3"}},
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"9"}},
		{"id": "3", "parents": []string{"7", "6"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{"10", "8"}},
		{"id": "7", "parents": []string{"11", "8"}},
		{"id": "8", "parents": []string{"9"}},
		{"id": "9", "parents": []string{"10"}},
		{"id": "10", "parents": []string{"11"}},
		{"id": "11", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 4, 1, 1, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"9": {[]*Point{{2, 2, 0}, {2, 9, 0}}, 2},
		},
		{
			"7": {[]*Point{{0, 3, 0}, {0, 7, 0}}, 0},
			"6": {[]*Point{{0, 3, 0}, {3, 3, 2}, {3, 6, 1}, {1, 6, 0}}, 3},
		},
		{
			"6": {[]*Point{{4, 4, 0}, {4, 6, 1}, {1, 6, 0}}, 4},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 0}}, 1},
		},
		{
			"10": {[]*Point{{1, 6, 0}, {1, 10, 0}}, 1},
			"8":  {[]*Point{{1, 6, 0}, {3, 6, 2}, {3, 8, 0}}, 5},
		},
		{
			"11": {[]*Point{{0, 7, 0}, {0, 11, 0}}, 0},
			"8":  {[]*Point{{0, 7, 0}, {3, 7, 2}, {3, 8, 0}}, 5},
		},
		{
			"9": {[]*Point{{3, 8, 0}, {3, 9, 1}, {2, 9, 0}}, 5},
		},
		{
			"10": {[]*Point{{2, 9, 0}, {2, 10, 1}, {1, 10, 0}}, 2},
		},
		{
			"11": {[]*Point{{1, 10, 0}, {1, 11, 1}, {0, 11, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test25(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"3"}},
		{"id": "2", "parents": []string{"4"}},
		{"id": "3", "parents": []string{"9", "7"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"8", "7"}},
		{"id": "6", "parents": []string{"9", "7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{"12", "9"}},
		{"id": "9", "parents": []string{"11", "10"}},
		{"id": "10", "parents": []string{"11"}},
		{"id": "11", "parents": []string{"12"}},
		{"id": "12", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 2, 0, 2, 3, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": {[]*Point{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"3": {[]*Point{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"4": {[]*Point{{2, 2, 0}, {2, 4, 0}}, 2},
		},
		{
			"9": {[]*Point{{1, 3, 0}, {1, 9, 0}}, 1},
			"7": {[]*Point{{1, 3, 0}, {3, 3, 2}, {3, 7, 0}}, 3},
		},
		{
			"6": {[]*Point{{2, 4, 0}, {2, 6, 0}}, 2},
		},
		{
			"8": {[]*Point{{0, 5, 0}, {0, 8, 0}}, 0},
			"7": {[]*Point{{0, 5, 0}, {3, 5, 2}, {3, 7, 0}}, 3},
		},
		{
			"9": {[]*Point{{2, 6, 0}, {2, 9, 1}, {1, 9, 0}}, 2},
			"7": {[]*Point{{2, 6, 0}, {3, 6, 2}, {3, 7, 0}}, 3},
		},
		{
			"8": {[]*Point{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, 3},
		},
		{
			"12": {[]*Point{{0, 8, 0}, {0, 12, 0}}, 0},
			"9":  {[]*Point{{0, 8, 0}, {1, 8, 2}, {1, 9, 0}}, 1},
		},
		{
			"11": {[]*Point{{1, 9, 0}, {1, 11, 0}}, 1},
			"10": {[]*Point{{1, 9, 0}, {2, 9, 2}, {2, 10, 0}}, 4},
		},
		{
			"11": {[]*Point{{2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, 4},
		},
		{
			"12": {[]*Point{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test26(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"3"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"8", "5"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"7", "6"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"8": {[]*Point{{0, 3, 0}, {0, 8, 0}}, 0},
			"5": {[]*Point{{0, 3, 0}, {2, 3, 2}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"5": {[]*Point{{1, 4, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*Point{{1, 5, 0}, {1, 7, 0}}, 1},
			"6": {[]*Point{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, 3},
		},
		{
			"7": {[]*Point{{2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, 3},
		},
		{
			"8": {[]*Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test27(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"7"}},
		{"id": "3", "parents": []string{"11"}},
		{"id": "4", "parents": []string{"15", "6"}},
		{"id": "5", "parents": []string{"8", "6"}},
		{"id": "6", "parents": []string{"15"}},
		{"id": "7", "parents": []string{"12"}},
		{"id": "8", "parents": []string{"9", "13"}},
		{"id": "9", "parents": []string{"14", "10"}},
		{"id": "10", "parents": []string{"14"}},
		{"id": "11", "parents": []string{"14"}},
		{"id": "12", "parents": []string{"14"}},
		{"id": "13", "parents": []string{"14"}},
		{"id": "14", "parents": []string{"16"}},
		{"id": "15", "parents": []string{"17"}},
		{"id": "16", "parents": []string{"17"}},
		{"id": "17", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 4, 2, 1, 1, 6, 3, 2, 5, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*Point{{2, 2, 0}, {2, 7, 0}}, 2},
		},
		{
			"11": {[]*Point{{3, 3, 0}, {3, 11, 0}}, 3},
		},
		{
			"15": {[]*Point{{0, 4, 0}, {0, 15, 0}}, 0},
			"6":  {[]*Point{{0, 4, 0}, {4, 4, 2}, {4, 6, 0}}, 4},
		},
		{
			"8": {[]*Point{{1, 5, 0}, {1, 8, 0}}, 1},
			"6": {[]*Point{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"15": {[]*Point{{4, 6, 0}, {4, 14, 1}, {2, 14, 0}, {2, 15, 1}, {0, 15, 0}}, 4},
		},
		{
			"12": {[]*Point{{2, 7, 0}, {2, 12, 0}}, 2},
		},
		{
			"9":  {[]*Point{{1, 8, 0}, {1, 9, 0}}, 1},
			"13": {[]*Point{{1, 8, 0}, {5, 8, 2}, {5, 13, 0}}, 5},
		},
		{
			"14": {[]*Point{{1, 9, 0}, {1, 14, 0}}, 1},
			"10": {[]*Point{{1, 9, 0}, {6, 9, 2}, {6, 10, 0}}, 6},
		},
		{
			"14": {[]*Point{{6, 10, 0}, {6, 14, 1}, {1, 14, 0}}, 6},
		},
		{
			"14": {[]*Point{{3, 11, 0}, {3, 14, 1}, {1, 14, 0}}, 3},
		},
		{
			"14": {[]*Point{{2, 12, 0}, {2, 14, 1}, {1, 14, 0}}, 2},
		},
		{
			"14": {[]*Point{{5, 13, 0}, {5, 14, 1}, {1, 14, 0}}, 5},
		},
		{
			"16": {[]*Point{{1, 14, 0}, {1, 16, 0}}, 1},
		},
		{
			"17": {[]*Point{{0, 15, 0}, {0, 17, 0}}, 0},
		},
		{
			"17": {[]*Point{{1, 16, 0}, {1, 17, 1}, {0, 17, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test28(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"2", "1"}},
		{"id": "1", "parents": []string{"2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"6", "5"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"2": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
			"1": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"2": {[]*Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
		},
		{
			"3": {[]*Point{{0, 2, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*Point{{0, 3, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
			"5": {[]*Point{{0, 4, 0}, {1, 4, 2}, {1, 5, 0}}, 1},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test29(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"7"}},
		{"id": "1", "parents": []string{"15"}},
		{"id": "2", "parents": []string{"17"}},
		{"id": "3", "parents": []string{"8"}},
		{"id": "4", "parents": []string{"18"}},
		{"id": "5", "parents": []string{"12"}},
		{"id": "6", "parents": []string{"20"}},
		{"id": "7", "parents": []string{"9", "10"}},
		{"id": "8", "parents": []string{"9", "11"}},
		{"id": "9", "parents": []string{"13", "14"}},
		{"id": "10", "parents": []string{"21"}},
		{"id": "11", "parents": []string{"13"}},
		{"id": "12", "parents": []string{"14"}},
		{"id": "13", "parents": []string{"16", "15"}},
		{"id": "14", "parents": []string{"19"}},
		{"id": "15", "parents": []string{"26"}},
		{"id": "16", "parents": []string{"27"}},
		{"id": "17", "parents": []string{"25"}},
		{"id": "18", "parents": []string{"24"}},
		{"id": "19", "parents": []string{"23"}},
		{"id": "20", "parents": []string{"22"}},
		{"id": "21", "parents": []string{"22"}},
		{"id": "22", "parents": []string{"23"}},
		{"id": "23", "parents": []string{"24"}},
		{"id": "24", "parents": []string{"25"}},
		{"id": "25", "parents": []string{"26"}},
		{"id": "26", "parents": []string{"27"}},
		{"id": "27", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 5, 6, 0, 3, 0, 7, 3, 5, 0, 4, 1, 0, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"7": {[]*Point{{0, 0, 0}, {0, 7, 0}}, 0},
		},
		{
			"15": {[]*Point{{1, 1, 0}, {1, 15, 0}}, 1},
		},
		{
			"17": {[]*Point{{2, 2, 0}, {2, 17, 0}}, 2},
		},
		{
			"8": {[]*Point{{3, 3, 0}, {3, 8, 0}}, 3},
		},
		{
			"18": {[]*Point{{4, 4, 0}, {4, 13, 1}, {3, 13, 0}, {3, 18, 0}}, 4},
		},
		{
			"12": {[]*Point{{5, 5, 0}, {5, 12, 0}}, 5},
		},
		{
			"20": {[]*Point{{6, 6, 0}, {6, 13, 1}, {5, 13, 0}, {5, 20, 0}}, 6},
		},
		{
			"9":  {[]*Point{{0, 7, 0}, {0, 9, 0}}, 0},
			"10": {[]*Point{{0, 7, 0}, {7, 7, 2}, {7, 10, 0}}, 7},
		},
		{
			"9":  {[]*Point{{3, 8, 0}, {0, 8, 3}, {0, 9, 0}}, 0},
			"11": {[]*Point{{3, 8, 0}, {3, 11, 0}}, 3},
		},
		{
			"13": {[]*Point{{0, 9, 0}, {0, 13, 0}}, 0},
			"14": {[]*Point{{0, 9, 0}, {8, 9, 2}, {8, 13, 1}, {7, 13, 0}, {7, 14, 1}, {4, 14, 0}}, 8},
		},
		{
			"21": {[]*Point{{7, 10, 0}, {7, 13, 1}, {6, 13, 0}, {6, 21, 0}}, 7},
		},
		{
			"13": {[]*Point{{3, 11, 0}, {3, 13, 1}, {0, 13, 0}}, 3},
		},
		{
			"14": {[]*Point{{5, 12, 0}, {5, 13, 1}, {4, 13, 0}, {4, 14, 0}}, 5},
		},
		{
			"16": {[]*Point{{0, 13, 0}, {0, 16, 0}}, 0},
			"15": {[]*Point{{0, 13, 0}, {1, 13, 2}, {1, 15, 0}}, 1},
		},
		{
			"19": {[]*Point{{4, 14, 0}, {4, 19, 0}}, 5},
		},
		{
			"26": {[]*Point{{1, 15, 0}, {1, 26, 0}}, 1},
		},
		{
			"27": {[]*Point{{0, 16, 0}, {0, 27, 0}}, 0},
		},
		{
			"25": {[]*Point{{2, 17, 0}, {2, 25, 0}}, 2},
		},
		{
			"24": {[]*Point{{3, 18, 0}, {3, 24, 0}}, 4},
		},
		{
			"23": {[]*Point{{4, 19, 0}, {4, 23, 0}}, 5},
		},
		{
			"22": {[]*Point{{5, 20, 0}, {5, 22, 0}}, 6},
		},
		{
			"22": {[]*Point{{6, 21, 0}, {6, 22, 1}, {5, 22, 0}}, 7},
		},
		{
			"23": {[]*Point{{5, 22, 0}, {5, 23, 1}, {4, 23, 0}}, 6},
		},
		{
			"24": {[]*Point{{4, 23, 0}, {4, 24, 1}, {3, 24, 0}}, 5},
		},
		{
			"25": {[]*Point{{3, 24, 0}, {3, 25, 1}, {2, 25, 0}}, 4},
		},
		{
			"26": {[]*Point{{2, 25, 0}, {2, 26, 1}, {1, 26, 0}}, 2},
		},
		{
			"27": {[]*Point{{1, 26, 0}, {1, 27, 1}, {0, 27, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test30(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4"}},
		{"id": "1", "parents": []string{"15"}},
		{"id": "2", "parents": []string{"14"}},
		{"id": "3", "parents": []string{"22"}},
		{"id": "4", "parents": []string{"5", "10"}},
		{"id": "5", "parents": []string{"6", "7"}},
		{"id": "6", "parents": []string{"16", "8"}},
		{"id": "7", "parents": []string{"9"}},
		{"id": "8", "parents": []string{"11"}},
		{"id": "9", "parents": []string{"16"}},
		{"id": "10", "parents": []string{"18"}},
		{"id": "11", "parents": []string{"12"}},
		{"id": "12", "parents": []string{"13", "16"}},
		{"id": "13", "parents": []string{"21"}},
		{"id": "14", "parents": []string{"18"}},
		{"id": "15", "parents": []string{"23"}},
		{"id": "16", "parents": []string{"18", "17"}},
		{"id": "17", "parents": []string{"18"}},
		{"id": "18", "parents": []string{"20", "19"}},
		{"id": "19", "parents": []string{"20"}},
		{"id": "20", "parents": []string{"24"}},
		{"id": "21", "parents": []string{"22"}},
		{"id": "22", "parents": []string{"23"}},
		{"id": "23", "parents": []string{"24"}},
		{"id": "24", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0, 5, 6, 5, 4, 6, 6, 6, 2, 1, 0, 6, 0, 4, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"15": {[]*Point{{1, 1, 0}, {1, 15, 0}}, 1},
		},
		{
			"14": {[]*Point{{2, 2, 0}, {2, 14, 0}}, 2},
		},
		{
			"22": {[]*Point{{3, 3, 0}, {3, 18, 1}, {2, 18, 0}, {2, 22, 0}}, 3},
		},
		{
			"5":  {[]*Point{{0, 4, 0}, {0, 5, 0}}, 0},
			"10": {[]*Point{{0, 4, 0}, {4, 4, 2}, {4, 10, 0}}, 4},
		},
		{
			"6": {[]*Point{{0, 5, 0}, {0, 6, 0}}, 0},
			"7": {[]*Point{{0, 5, 0}, {5, 5, 2}, {5, 7, 0}}, 5},
		},
		{
			"16": {[]*Point{{0, 6, 0}, {0, 16, 0}}, 0},
			"8":  {[]*Point{{0, 6, 0}, {6, 6, 2}, {6, 8, 0}}, 6},
		},
		{
			"9": {[]*Point{{5, 7, 0}, {5, 9, 0}}, 5},
		},
		{
			"11": {[]*Point{{6, 8, 0}, {6, 11, 0}}, 6},
		},
		{
			"16": {[]*Point{{5, 9, 0}, {5, 16, 1}, {0, 16, 0}}, 5},
		},
		{
			"18": {[]*Point{{4, 10, 0}, {4, 18, 1}, {0, 18, 0}}, 4},
		},
		{
			"12": {[]*Point{{6, 11, 0}, {6, 12, 0}}, 6},
		},
		{
			"13": {[]*Point{{6, 12, 0}, {6, 13, 0}}, 6},
			"16": {[]*Point{{6, 12, 0}, {0, 12, 3}, {0, 16, 0}}, 0},
		},
		{
			"21": {[]*Point{{6, 13, 0}, {6, 16, 1}, {5, 16, 0}, {5, 18, 1}, {3, 18, 0}, {3, 21, 0}}, 6},
		},
		{
			"18": {[]*Point{{2, 14, 0}, {2, 18, 1}, {0, 18, 0}}, 2},
		},
		{
			"23": {[]*Point{{1, 15, 0}, {1, 23, 0}}, 1},
		},
		{
			"18": {[]*Point{{0, 16, 0}, {0, 18, 0}}, 0},
			"17": {[]*Point{{0, 16, 0}, {6, 16, 2}, {6, 17, 0}}, 7},
		},
		{
			"18": {[]*Point{{6, 17, 0}, {6, 18, 1}, {0, 18, 0}}, 7},
		},
		{
			"20": {[]*Point{{0, 18, 0}, {0, 20, 0}}, 0},
			"19": {[]*Point{{0, 18, 0}, {4, 18, 2}, {4, 19, 0}}, 5},
		},
		{
			"20": {[]*Point{{4, 19, 0}, {4, 20, 1}, {0, 20, 0}}, 5},
		},
		{
			"24": {[]*Point{{0, 20, 0}, {0, 24, 0}}, 0},
		},
		{
			"22": {[]*Point{{3, 21, 0}, {3, 22, 1}, {2, 22, 0}}, 6},
		},
		{
			"23": {[]*Point{{2, 22, 0}, {2, 23, 1}, {1, 23, 0}}, 3},
		},
		{
			"24": {[]*Point{{1, 23, 0}, {1, 24, 1}, {0, 24, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test31(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"3"}},
		{"id": "1", "parents": []string{"4"}},
		{"id": "2", "parents": []string{"5", "4"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 2},
			"4": {[]*Point{{2, 2, 0}, {1, 2, 3}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*Point{{0, 3, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test32(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"2"}},
		{"id": "1", "parents": []string{"5", "3"}},
		{"id": "2", "parents": []string{"3", "4"}},
		{"id": "3", "parents": []string{"6"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"2": {[]*Point{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
			"3": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {0, 3, 0}}, 2},
		},
		{
			"3": {[]*Point{{0, 2, 0}, {0, 3, 0}}, 0},
			"4": {[]*Point{{0, 2, 0}, {3, 2, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 0}}, 3},
		},
		{
			"6": {[]*Point{{0, 3, 0}, {0, 6, 0}}, 0},
		},
		{
			"5": {[]*Point{{2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test33(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"3"}},
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"7"}},
		{"id": "3", "parents": []string{"9", "4"}},
		{"id": "4", "parents": []string{"9"}},
		{"id": "5", "parents": []string{"8", "6"}},
		{"id": "6", "parents": []string{"10"}},
		{"id": "7", "parents": []string{"10"}},
		{"id": "8", "parents": []string{"10"}},
		{"id": "9", "parents": []string{"13"}},
		{"id": "10", "parents": []string{"12", "11"}},
		{"id": "11", "parents": []string{"12"}},
		{"id": "12", "parents": []string{"13"}},
		{"id": "13", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 3, 1, 4, 2, 1, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": {[]*Point{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"5": {[]*Point{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*Point{{2, 2, 0}, {2, 7, 0}}, 2},
		},
		{
			"9": {[]*Point{{0, 3, 0}, {0, 9, 0}}, 0},
			"4": {[]*Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 0}}, 3},
		},
		{
			"9": {[]*Point{{3, 4, 0}, {3, 9, 1}, {0, 9, 0}}, 3},
		},
		{
			"8": {[]*Point{{1, 5, 0}, {1, 8, 0}}, 1},
			"6": {[]*Point{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"10": {[]*Point{{4, 6, 0}, {4, 9, 1}, {3, 9, 0}, {3, 10, 1}, {1, 10, 0}}, 4},
		},
		{
			"10": {[]*Point{{2, 7, 0}, {2, 10, 1}, {1, 10, 0}}, 2},
		},
		{
			"10": {[]*Point{{1, 8, 0}, {1, 10, 0}}, 1},
		},
		{
			"13": {[]*Point{{0, 9, 0}, {0, 13, 0}}, 0},
		},
		{
			"12": {[]*Point{{1, 10, 0}, {1, 12, 0}}, 1},
			"11": {[]*Point{{1, 10, 0}, {2, 10, 2}, {2, 11, 0}}, 5},
		},
		{
			"12": {[]*Point{{2, 11, 0}, {2, 12, 1}, {1, 12, 0}}, 5},
		},
		{
			"13": {[]*Point{{1, 12, 0}, {1, 13, 1}, {0, 13, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test34(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"4", "2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{"4", "5"}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"7"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 2, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": {[]*Point{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 0}}, 1},
			"2": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"3": {[]*Point{{2, 2, 0}, {2, 3, 0}}, 2},
		},
		{
			"4": {[]*Point{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
			"5": {[]*Point{{2, 3, 0}, {0, 3, 3}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*Point{{1, 4, 0}, {1, 6, 0}}, 1},
		},
		{
			"7": {[]*Point{{0, 5, 0}, {0, 7, 0}}, 0},
		},
		{
			"7": {[]*Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test35(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4", "1"}},
		{"id": "1", "parents": []string{"2", "3"}},
		{"id": "2", "parents": []string{}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 1, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
			"1": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"2": {[]*Point{{1, 1, 0}, {1, 2, 0}}, 1},
			"3": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {1, 3, 0}}, 2},
		},
		{
			"4": {[]*Point{{1, 3, 0}, {1, 4, 1}, {0, 4, 0}}, 2},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test36(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"4", "1"}},
		{"id": "1", "parents": []string{"4", "2"}},
		{"id": "2", "parents": []string{"3", "5"}},
		{"id": "3", "parents": []string{}},
		{"id": "4", "parents": []string{"6"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{}},
	}

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {[]*Point{{0, 0, 0}, {0, 4, 0}}, 0},
			"1": {[]*Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
			"2": {[]*Point{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"3": {[]*Point{{2, 2, 0}, {2, 3, 0}}, 2},
			"5": {[]*Point{{2, 2, 0}, {3, 2, 2}, {3, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 3},
		},
		{
			"6": {[]*Point{{0, 4, 0}, {0, 6, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func assertEq(t *testing.T, expected, actual any) {
	if actual != expected {
		t.Fail()
	}
}

func TestPathHeight1(t *testing.T) {
	path := &Path{Points: []*Point{
		{x: 0, y: 2, typ: 0},
		{x: 3, y: 2, typ: 2},
		{x: 3, y: 9, typ: 1},
		{x: 2, y: 9, typ: 0},
		{x: 2, y: 11, typ: 1},
		{x: 1, y: 11, typ: 0},
	}}
	assertEq(t, -1, path.GetHeightAtIdx(1))
	assertEq(t, 3, path.GetHeightAtIdx(2))
	assertEq(t, 3, path.GetHeightAtIdx(3))
	assertEq(t, 2, path.GetHeightAtIdx(9))
	assertEq(t, 2, path.GetHeightAtIdx(10))
	assertEq(t, 1, path.GetHeightAtIdx(11))
	assertEq(t, -1, path.GetHeightAtIdx(1000))
}

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inputNodes := []*Node{
			{"id": "0", "parents": []string{"4"}},
			{"id": "1", "parents": []string{"15"}},
			{"id": "2", "parents": []string{"14"}},
			{"id": "3", "parents": []string{"22"}},
			{"id": "4", "parents": []string{"5", "10"}},
			{"id": "5", "parents": []string{"6", "7"}},
			{"id": "6", "parents": []string{"16", "8"}},
			{"id": "7", "parents": []string{"9"}},
			{"id": "8", "parents": []string{"11"}},
			{"id": "9", "parents": []string{"16"}},
			{"id": "10", "parents": []string{"18"}},
			{"id": "11", "parents": []string{"12"}},
			{"id": "12", "parents": []string{"13", "16"}},
			{"id": "13", "parents": []string{"21"}},
			{"id": "14", "parents": []string{"18"}},
			{"id": "15", "parents": []string{"23"}},
			{"id": "16", "parents": []string{"18", "17"}},
			{"id": "17", "parents": []string{"18"}},
			{"id": "18", "parents": []string{"20", "19"}},
			{"id": "19", "parents": []string{"20"}},
			{"id": "20", "parents": []string{"24"}},
			{"id": "21", "parents": []string{"22"}},
			{"id": "22", "parents": []string{"23"}},
			{"id": "23", "parents": []string{"24"}},
			{"id": "24", "parents": []string{}},
		}
		if _, err := BuildTree(inputNodes, customColors); err != nil {
			b.Logf("Failed to build tree")
		}
	}
}
