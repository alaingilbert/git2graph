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
			t.Logf("Columns, ID: %s, Expected column: %d, Actual column: %d", nodeID, expectedColumn, actualColumn)
		}
	}
}

func pprintPoints(points []*PointTest) string {
	s := make([]string, 0)
	for _, p := range points {
		s = append(s, p.String())
	}
	return "[" + strings.Join(s, ",") + "]"
}

func pprintPoints1(points []IPoint) string {
	s := make([]string, 0)
	for _, p := range points {
		s = append(s, p.String())
	}
	return "[" + strings.Join(s, ",") + "]"
}

func validatePaths(t *testing.T, expectedPaths []map[string]PathTest, data []*Node) {
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
				t.Logf("ID: %s, Expected vs Actual:\n%v\n%v", nodeID, pprintPoints(expectedPath), pprintPoints1(parentPath.Points))
				return
			}
			for pathIdx, pathNode := range parentPath.Points {
				if !pathNode.Equal(expectedPath[pathIdx]) {
					t.Fail()
					t.Logf("ID: %s, Expected path: %s, Actual path: %s", nodeID, expectedPath[pathIdx].String(), pathNode.String())
					t.Logf("ID: %s, Expected vs Actual:\n%v\n%v", nodeID, pprintPoints(expectedPath), pprintPoints1(parentPath.Points))
				}
			}
		}
	}
}

func validateColors(t *testing.T, expectedPaths []map[string]PathTest, data []*Node) {
	for nodeIdx, node := range data {
		nodeID := (*node)[idKey]
		for _, parentID := range (*node)[parentsKey].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			parentPath := (*node)[parentsPathsTestKey].(map[string]*Path)[parentID]
			expectedPath := expectedPaths[nodeIdx][parentID]
			if expectedPath.colorIdx != parentPath.colorIdx {
				t.Logf("Colors, ID: %s -> %s, Expected: %v, Actual: %v", nodeID, parentID, expectedPath.colorIdx, parentPath.colorIdx)
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

var realColors = NewSimpleColorGen(DefaultColors)

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
	out, _ := buildTreeTest(inputNodes, colors, "", -1)
	if (*out[2])[gKey].([]any)[2] != "#000" {
		t.Fail()
	}
}

func TestGetInputNodesFromJson(t *testing.T) {
	json := `[{"id": "1", "parents": ["2"]}, {"id": "2", "parents": ["3"]}, {"id": "3", "parents": []}]`
	inputNodes, _ := GetInputNodesFromJSON([]byte(json))
	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{"2": PathTest{[]*PointTest{{0, 0, 0}, {0, 1, 0}}, 0}},
		{"3": PathTest{[]*PointTest{{0, 1, 0}, {0, 2, 0}}, 0}},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"2": {[]*PointTest{{0, 0, 0}, {0, 1, 0}}, 0},
		},
		{
			"3": {[]*PointTest{{0, 1, 0}, {0, 2, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{0, 2, 0}, {0, 4, 0}}, 0},
			"4": {[]*PointTest{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{2, 3, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{2, 2, 0}, {2, 3, 1}, {0, 3, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{0, 3, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 3},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, 0},
			"4": {[]*PointTest{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 3, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 3, 0}}, 1},
			"5": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, 1},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{0, 2, 0}, {0, 3, 0}}, 0},
			"5": {[]*PointTest{{0, 2, 0}, {2, 2, 2}, {2, 3, 1}, {1, 3, 0}, {1, 4, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{0, 3, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{0, 2, 0}, {0, 3, 0}}, 0},
			"7": {[]*PointTest{{0, 2, 0}, {2, 2, 2}, {2, 4, 1}, {1, 4, 0}, {1, 6, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{0, 3, 0}, {0, 4, 0}}, 0},
			"6": {[]*PointTest{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{0, 4, 0}, {0, 7, 0}}, 0},
		},
		{
			"8": {[]*PointTest{{2, 5, 0}, {2, 7, 1}, {0, 7, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
			"3": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{0, 3, 0}, {0, 7, 0}}, 0},
			"6": {[]*PointTest{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{1, 4, 0}, {2, 4, 2}, {2, 5, 0}}, 3},
			"7": {[]*PointTest{{1, 4, 0}, {1, 6, 0}}, 1},
		},
		{
			"7": {[]*PointTest{{2, 5, 0}, {2, 6, 1}, {1, 6, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, 0},
			"4": {[]*PointTest{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{0, 2, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 3, 0}, {0, 3, 3}, {0, 4, 0}}, 0},
			"6": {[]*PointTest{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, 1},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 1, 0, 1, 1, 0, 0, 2, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
			"2": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 1, 0, 2, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
			"4": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{0, 2, 0}, {0, 5, 0}}, 0},
			"4": {[]*PointTest{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{1, 4, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
			"7": {[]*PointTest{{1, 4, 0}, {2, 4, 2}, {2, 6, 0}}, 3},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 2, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
			"4": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{0, 2, 0}, {0, 4, 0}}, 0},
			"4": {[]*PointTest{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{2, 3, 0}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{0, 4, 0}, {0, 7, 0}}, 0},
			"7": {[]*PointTest{{0, 4, 0}, {3, 4, 2}, {3, 5, 1}, {2, 5, 0}, {2, 6, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{1, 5, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
			"7": {[]*PointTest{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, 3},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"7": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"7": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 6, 1}, {0, 6, 0}}, 3},
		},
		{
			"5": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{0, 5, 0}, {0, 6, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 4, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"9": {[]*PointTest{{2, 2, 0}, {2, 9, 0}}, 2},
		},
		{
			"7": {[]*PointTest{{3, 3, 0}, {3, 7, 0}}, 3},
		},
		{
			"11": {[]*PointTest{{1, 4, 0}, {1, 11, 0}}, 1},
			"6":  {[]*PointTest{{1, 4, 0}, {4, 4, 2}, {4, 6, 0}}, 4},
		},
		{
			"8": {[]*PointTest{{0, 5, 0}, {0, 8, 0}}, 0},
			"6": {[]*PointTest{{0, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"11": {[]*PointTest{{4, 6, 0}, {4, 8, 1}, {3, 8, 0}, {3, 10, 1}, {2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, 4},
		},
		{
			"8": {[]*PointTest{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, 3},
		},
		{
			"10": {[]*PointTest{{0, 8, 0}, {0, 10, 0}}, 0},
		},
		{
			"10": {[]*PointTest{{2, 9, 0}, {2, 10, 1}, {0, 10, 0}}, 2},
		},
		{
			"12": {[]*PointTest{{0, 10, 0}, {0, 12, 0}}, 0},
		},
		{
			"12": {[]*PointTest{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"4": {[]*PointTest{{3, 3, 0}, {3, 4, 1}, {0, 4, 0}}, 3},
		},
		{
			"5": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{1, 3, 0}, {1, 6, 0}}, 1},
			"5": {[]*PointTest{{1, 3, 0}, {2, 3, 2}, {2, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
		},
		{
			"7": {[]*PointTest{{0, 5, 0}, {0, 7, 0}}, 0},
		},
		{
			"7": {[]*PointTest{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"6": {[]*PointTest{{2, 2, 0}, {2, 6, 1}, {0, 6, 0}}, 2},
		},
		{
			"7": {[]*PointTest{{3, 3, 0}, {3, 6, 1}, {2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, 3},
		},
		{
			"7": {[]*PointTest{{1, 4, 0}, {1, 7, 0}}, 1},
			"6": {[]*PointTest{{1, 4, 0}, {2, 4, 2}, {2, 6, 1}, {0, 6, 0}}, 2},
		},
		{
			"6": {[]*PointTest{{0, 5, 0}, {0, 6, 0}}, 0},
		},
		{
			"8": {[]*PointTest{{0, 6, 0}, {0, 8, 0}}, 0},
		},
		{
			"8": {[]*PointTest{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, 2},
		},
		{
			"7": {[]*PointTest{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 7, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
			"5": {[]*PointTest{{0, 4, 0}, {2, 4, 2}, {2, 5, 0}}, 4},
		},
		{
			"6": {[]*PointTest{{2, 5, 0}, {2, 6, 1}, {0, 6, 0}}, 4},
		},
		{
			"8": {[]*PointTest{{0, 6, 0}, {0, 8, 0}}, 0},
		},
		{
			"8": {[]*PointTest{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 3},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 4, 1, 1, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"9": {[]*PointTest{{2, 2, 0}, {2, 9, 0}}, 2},
		},
		{
			"7": {[]*PointTest{{0, 3, 0}, {0, 7, 0}}, 0},
			"6": {[]*PointTest{{0, 3, 0}, {3, 3, 2}, {3, 6, 1}, {1, 6, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{4, 4, 0}, {4, 6, 1}, {1, 6, 0}}, 4},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 0}}, 1},
		},
		{
			"10": {[]*PointTest{{1, 6, 0}, {1, 10, 0}}, 1},
			"8":  {[]*PointTest{{1, 6, 0}, {3, 6, 2}, {3, 8, 0}}, 5},
		},
		{
			"11": {[]*PointTest{{0, 7, 0}, {0, 11, 0}}, 0},
			"8":  {[]*PointTest{{0, 7, 0}, {3, 7, 2}, {3, 8, 0}}, 5},
		},
		{
			"9": {[]*PointTest{{3, 8, 0}, {3, 9, 1}, {2, 9, 0}}, 5},
		},
		{
			"10": {[]*PointTest{{2, 9, 0}, {2, 10, 1}, {1, 10, 0}}, 2},
		},
		{
			"11": {[]*PointTest{{1, 10, 0}, {1, 11, 1}, {0, 11, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 2, 0, 2, 3, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"3": {[]*PointTest{{1, 1, 0}, {1, 3, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{2, 2, 0}, {2, 4, 0}}, 2},
		},
		{
			"9": {[]*PointTest{{1, 3, 0}, {1, 9, 0}}, 1},
			"7": {[]*PointTest{{1, 3, 0}, {3, 3, 2}, {3, 7, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{2, 4, 0}, {2, 6, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{0, 5, 0}, {0, 8, 0}}, 0},
			"7": {[]*PointTest{{0, 5, 0}, {3, 5, 2}, {3, 7, 0}}, 3},
		},
		{
			"9": {[]*PointTest{{2, 6, 0}, {2, 9, 1}, {1, 9, 0}}, 2},
			"7": {[]*PointTest{{2, 6, 0}, {3, 6, 2}, {3, 7, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, 3},
		},
		{
			"12": {[]*PointTest{{0, 8, 0}, {0, 12, 0}}, 0},
			"9":  {[]*PointTest{{0, 8, 0}, {1, 8, 2}, {1, 9, 0}}, 1},
		},
		{
			"11": {[]*PointTest{{1, 9, 0}, {1, 11, 0}}, 1},
			"10": {[]*PointTest{{1, 9, 0}, {2, 9, 2}, {2, 10, 0}}, 4},
		},
		{
			"11": {[]*PointTest{{2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, 4},
		},
		{
			"12": {[]*PointTest{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 1, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{0, 3, 0}, {0, 8, 0}}, 0},
			"5": {[]*PointTest{{0, 3, 0}, {2, 3, 2}, {2, 5, 1}, {1, 5, 0}}, 2},
		},
		{
			"5": {[]*PointTest{{1, 4, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*PointTest{{1, 5, 0}, {1, 7, 0}}, 1},
			"6": {[]*PointTest{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, 3},
		},
		{
			"7": {[]*PointTest{{2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 4, 2, 1, 1, 6, 3, 2, 5, 1, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*PointTest{{2, 2, 0}, {2, 7, 0}}, 2},
		},
		{
			"11": {[]*PointTest{{3, 3, 0}, {3, 11, 0}}, 3},
		},
		{
			"15": {[]*PointTest{{0, 4, 0}, {0, 15, 0}}, 0},
			"6":  {[]*PointTest{{0, 4, 0}, {4, 4, 2}, {4, 6, 0}}, 4},
		},
		{
			"8": {[]*PointTest{{1, 5, 0}, {1, 8, 0}}, 1},
			"6": {[]*PointTest{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"15": {[]*PointTest{{4, 6, 0}, {4, 14, 1}, {2, 14, 0}, {2, 15, 1}, {0, 15, 0}}, 4},
		},
		{
			"12": {[]*PointTest{{2, 7, 0}, {2, 12, 0}}, 2},
		},
		{
			"9":  {[]*PointTest{{1, 8, 0}, {1, 9, 0}}, 1},
			"13": {[]*PointTest{{1, 8, 0}, {5, 8, 2}, {5, 13, 0}}, 5},
		},
		{
			"14": {[]*PointTest{{1, 9, 0}, {1, 14, 0}}, 1},
			"10": {[]*PointTest{{1, 9, 0}, {6, 9, 2}, {6, 10, 0}}, 6},
		},
		{
			"14": {[]*PointTest{{6, 10, 0}, {6, 14, 1}, {1, 14, 0}}, 6},
		},
		{
			"14": {[]*PointTest{{3, 11, 0}, {3, 14, 1}, {1, 14, 0}}, 3},
		},
		{
			"14": {[]*PointTest{{2, 12, 0}, {2, 14, 1}, {1, 14, 0}}, 2},
		},
		{
			"14": {[]*PointTest{{5, 13, 0}, {5, 14, 1}, {1, 14, 0}}, 5},
		},
		{
			"16": {[]*PointTest{{1, 14, 0}, {1, 16, 0}}, 1},
		},
		{
			"17": {[]*PointTest{{0, 15, 0}, {0, 17, 0}}, 0},
		},
		{
			"17": {[]*PointTest{{1, 16, 0}, {1, 17, 1}, {0, 17, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"2": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
			"1": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"2": {[]*PointTest{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, 1},
		},
		{
			"3": {[]*PointTest{{0, 2, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{0, 3, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
			"5": {[]*PointTest{{0, 4, 0}, {1, 4, 2}, {1, 5, 0}}, 1},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 5, 6, 0, 3, 0, 7, 3, 5, 0, 4, 1, 0, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"7": {[]*PointTest{{0, 0, 0}, {0, 7, 0}}, 0},
		},
		{
			"15": {[]*PointTest{{1, 1, 0}, {1, 15, 0}}, 1},
		},
		{
			"17": {[]*PointTest{{2, 2, 0}, {2, 17, 0}}, 2},
		},
		{
			"8": {[]*PointTest{{3, 3, 0}, {3, 8, 0}}, 3},
		},
		{
			"18": {[]*PointTest{{4, 4, 0}, {4, 13, 1}, {3, 13, 0}, {3, 18, 0}}, 4},
		},
		{
			"12": {[]*PointTest{{5, 5, 0}, {5, 12, 0}}, 5},
		},
		{
			"20": {[]*PointTest{{6, 6, 0}, {6, 13, 1}, {5, 13, 0}, {5, 20, 0}}, 6},
		},
		{
			"9":  {[]*PointTest{{0, 7, 0}, {0, 9, 0}}, 0},
			"10": {[]*PointTest{{0, 7, 0}, {7, 7, 2}, {7, 10, 0}}, 7},
		},
		{
			"9":  {[]*PointTest{{3, 8, 0}, {0, 8, 3}, {0, 9, 0}}, 0},
			"11": {[]*PointTest{{3, 8, 0}, {3, 11, 0}}, 3},
		},
		{
			"13": {[]*PointTest{{0, 9, 0}, {0, 13, 0}}, 0},
			"14": {[]*PointTest{{0, 9, 0}, {8, 9, 2}, {8, 13, 1}, {7, 13, 0}, {7, 14, 1}, {4, 14, 0}}, 8},
		},
		{
			"21": {[]*PointTest{{7, 10, 0}, {7, 13, 1}, {6, 13, 0}, {6, 21, 0}}, 7},
		},
		{
			"13": {[]*PointTest{{3, 11, 0}, {3, 13, 1}, {0, 13, 0}}, 3},
		},
		{
			"14": {[]*PointTest{{5, 12, 0}, {5, 13, 1}, {4, 13, 0}, {4, 14, 0}}, 5},
		},
		{
			"16": {[]*PointTest{{0, 13, 0}, {0, 16, 0}}, 0},
			"15": {[]*PointTest{{0, 13, 0}, {1, 13, 2}, {1, 15, 0}}, 1},
		},
		{
			"19": {[]*PointTest{{4, 14, 0}, {4, 19, 0}}, 5},
		},
		{
			"26": {[]*PointTest{{1, 15, 0}, {1, 26, 0}}, 1},
		},
		{
			"27": {[]*PointTest{{0, 16, 0}, {0, 27, 0}}, 0},
		},
		{
			"25": {[]*PointTest{{2, 17, 0}, {2, 25, 0}}, 2},
		},
		{
			"24": {[]*PointTest{{3, 18, 0}, {3, 24, 0}}, 4},
		},
		{
			"23": {[]*PointTest{{4, 19, 0}, {4, 23, 0}}, 5},
		},
		{
			"22": {[]*PointTest{{5, 20, 0}, {5, 22, 0}}, 6},
		},
		{
			"22": {[]*PointTest{{6, 21, 0}, {6, 22, 1}, {5, 22, 0}}, 7},
		},
		{
			"23": {[]*PointTest{{5, 22, 0}, {5, 23, 1}, {4, 23, 0}}, 6},
		},
		{
			"24": {[]*PointTest{{4, 23, 0}, {4, 24, 1}, {3, 24, 0}}, 5},
		},
		{
			"25": {[]*PointTest{{3, 24, 0}, {3, 25, 1}, {2, 25, 0}}, 4},
		},
		{
			"26": {[]*PointTest{{2, 25, 0}, {2, 26, 1}, {1, 26, 0}}, 2},
		},
		{
			"27": {[]*PointTest{{1, 26, 0}, {1, 27, 1}, {0, 27, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0, 5, 6, 5, 4, 6, 6, 6, 2, 1, 0, 6, 0, 4, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
		},
		{
			"15": {[]*PointTest{{1, 1, 0}, {1, 15, 0}}, 1},
		},
		{
			"14": {[]*PointTest{{2, 2, 0}, {2, 14, 0}}, 2},
		},
		{
			"22": {[]*PointTest{{3, 3, 0}, {3, 18, 1}, {2, 18, 0}, {2, 22, 0}}, 3},
		},
		{
			"5":  {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0},
			"10": {[]*PointTest{{0, 4, 0}, {4, 4, 2}, {4, 10, 0}}, 4},
		},
		{
			"6": {[]*PointTest{{0, 5, 0}, {0, 6, 0}}, 0},
			"7": {[]*PointTest{{0, 5, 0}, {5, 5, 2}, {5, 7, 0}}, 5},
		},
		{
			"16": {[]*PointTest{{0, 6, 0}, {0, 16, 0}}, 0},
			"8":  {[]*PointTest{{0, 6, 0}, {6, 6, 2}, {6, 8, 0}}, 6},
		},
		{
			"9": {[]*PointTest{{5, 7, 0}, {5, 9, 0}}, 5},
		},
		{
			"11": {[]*PointTest{{6, 8, 0}, {6, 11, 0}}, 6},
		},
		{
			"16": {[]*PointTest{{5, 9, 0}, {5, 16, 1}, {0, 16, 0}}, 5},
		},
		{
			"18": {[]*PointTest{{4, 10, 0}, {4, 18, 1}, {0, 18, 0}}, 4},
		},
		{
			"12": {[]*PointTest{{6, 11, 0}, {6, 12, 0}}, 6},
		},
		{
			"13": {[]*PointTest{{6, 12, 0}, {6, 13, 0}}, 6},
			"16": {[]*PointTest{{6, 12, 0}, {0, 12, 3}, {0, 16, 0}}, 0},
		},
		{
			"21": {[]*PointTest{{6, 13, 0}, {6, 16, 1}, {5, 16, 0}, {5, 18, 1}, {3, 18, 0}, {3, 21, 0}}, 6},
		},
		{
			"18": {[]*PointTest{{2, 14, 0}, {2, 18, 1}, {0, 18, 0}}, 2},
		},
		{
			"23": {[]*PointTest{{1, 15, 0}, {1, 23, 0}}, 1},
		},
		{
			"18": {[]*PointTest{{0, 16, 0}, {0, 18, 0}}, 0},
			"17": {[]*PointTest{{0, 16, 0}, {6, 16, 2}, {6, 17, 0}}, 7},
		},
		{
			"18": {[]*PointTest{{6, 17, 0}, {6, 18, 1}, {0, 18, 0}}, 7},
		},
		{
			"20": {[]*PointTest{{0, 18, 0}, {0, 20, 0}}, 0},
			"19": {[]*PointTest{{0, 18, 0}, {4, 18, 2}, {4, 19, 0}}, 5},
		},
		{
			"20": {[]*PointTest{{4, 19, 0}, {4, 20, 1}, {0, 20, 0}}, 5},
		},
		{
			"24": {[]*PointTest{{0, 20, 0}, {0, 24, 0}}, 0},
		},
		{
			"22": {[]*PointTest{{3, 21, 0}, {3, 22, 1}, {2, 22, 0}}, 6},
		},
		{
			"23": {[]*PointTest{{2, 22, 0}, {2, 23, 1}, {1, 23, 0}}, 3},
		},
		{
			"24": {[]*PointTest{{1, 23, 0}, {1, 24, 1}, {0, 24, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 2},
			"4": {[]*PointTest{{2, 2, 0}, {1, 2, 3}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{0, 3, 0}, {0, 4, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"2": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
			"3": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {0, 3, 0}}, 2},
		},
		{
			"3": {[]*PointTest{{0, 2, 0}, {0, 3, 0}}, 0},
			"4": {[]*PointTest{{0, 2, 0}, {3, 2, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{0, 3, 0}, {0, 6, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 3, 1, 4, 2, 1, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"3": {[]*PointTest{{0, 0, 0}, {0, 3, 0}}, 0},
		},
		{
			"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1},
		},
		{
			"7": {[]*PointTest{{2, 2, 0}, {2, 7, 0}}, 2},
		},
		{
			"9": {[]*PointTest{{0, 3, 0}, {0, 9, 0}}, 0},
			"4": {[]*PointTest{{0, 3, 0}, {3, 3, 2}, {3, 4, 0}}, 3},
		},
		{
			"9": {[]*PointTest{{3, 4, 0}, {3, 9, 1}, {0, 9, 0}}, 3},
		},
		{
			"8": {[]*PointTest{{1, 5, 0}, {1, 8, 0}}, 1},
			"6": {[]*PointTest{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, 4},
		},
		{
			"10": {[]*PointTest{{4, 6, 0}, {4, 9, 1}, {3, 9, 0}, {3, 10, 1}, {1, 10, 0}}, 4},
		},
		{
			"10": {[]*PointTest{{2, 7, 0}, {2, 10, 1}, {1, 10, 0}}, 2},
		},
		{
			"10": {[]*PointTest{{1, 8, 0}, {1, 10, 0}}, 1},
		},
		{
			"13": {[]*PointTest{{0, 9, 0}, {0, 13, 0}}, 0},
		},
		{
			"12": {[]*PointTest{{1, 10, 0}, {1, 12, 0}}, 1},
			"11": {[]*PointTest{{1, 10, 0}, {2, 10, 2}, {2, 11, 0}}, 5},
		},
		{
			"12": {[]*PointTest{{2, 11, 0}, {2, 12, 1}, {1, 12, 0}}, 5},
		},
		{
			"13": {[]*PointTest{{1, 12, 0}, {1, 13, 1}, {0, 13, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 2, 1, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 0}}, 1},
			"2": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"3": {[]*PointTest{{2, 2, 0}, {2, 3, 0}}, 2},
		},
		{
			"4": {[]*PointTest{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, 2},
			"5": {[]*PointTest{{2, 3, 0}, {0, 3, 3}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 4, 0}, {1, 6, 0}}, 1},
		},
		{
			"7": {[]*PointTest{{0, 5, 0}, {0, 7, 0}}, 0},
		},
		{
			"7": {[]*PointTest{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 1},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 1, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
			"1": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"2": {[]*PointTest{{1, 1, 0}, {1, 2, 0}}, 1},
			"3": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {1, 3, 0}}, 2},
		},
		{},
		{
			"4": {[]*PointTest{{1, 3, 0}, {1, 4, 1}, {0, 4, 0}}, 2},
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

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 2, 0, 1, 0}

	expectedPaths := []map[string]PathTest{
		{
			"4": {[]*PointTest{{0, 0, 0}, {0, 4, 0}}, 0},
			"1": {[]*PointTest{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, 1},
		},
		{
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
			"2": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, 2},
		},
		{
			"3": {[]*PointTest{{2, 2, 0}, {2, 3, 0}}, 2},
			"5": {[]*PointTest{{2, 2, 0}, {3, 2, 2}, {3, 4, 1}, {1, 4, 0}, {1, 5, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{0, 4, 0}, {0, 6, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test37(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"6"}},
		{"id": "2", "parents": []string{"5"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"8", "7"}},
		{"id": "5", "parents": []string{"6"}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 3, 0, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{
			"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0},
		},
		{
			"6": {[]*PointTest{{1, 1, 0}, {1, 6, 1}, {0, 6, 0}}, 1},
		},
		{
			"5": {[]*PointTest{{2, 2, 0}, {2, 5, 1}, {0, 5, 0}}, 2},
		},
		{
			"4": {[]*PointTest{{3, 3, 0}, {3, 4, 0}}, 3},
		},
		{
			"7": {[]*PointTest{{3, 4, 0}, {4, 4, 2}, {4, 5, 1}, {3, 5, 0}, {3, 6, 1}, {2, 6, 0}, {2, 7, 1}, {0, 7, 0}}, 4},
			"8": {[]*PointTest{{3, 4, 0}, {3, 5, 1}, {2, 5, 0}, {2, 6, 1}, {1, 6, 0}, {1, 8, 1}, {0, 8, 0}}, 3},
		},
		{
			"6": {[]*PointTest{{0, 5, 0}, {0, 6, 0}}, 0},
		},
		{
			"7": {[]*PointTest{{0, 6, 0}, {0, 7, 0}}, 0},
		},
		{
			"8": {[]*PointTest{{0, 7, 0}, {0, 8, 0}}, 0},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test38(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"2"}},
		{"id": "2", "parents": []string{"6"}},
		{"id": "3", "parents": []string{"6"}},
		{"id": "4", "parents": []string{"7"}},
		{"id": "5", "parents": []string{}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{}},
	}

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 1, 2, 3, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{"5": {[]*PointTest{{0, 0, 0}, {0, 5, 0}}, 0}},
		{"2": {[]*PointTest{{1, 1, 0}, {1, 2, 0}}, 1}},
		{"6": {[]*PointTest{{1, 2, 0}, {1, 6, 1}, {0, 6, 0}}, 1}},
		{"6": {[]*PointTest{{2, 3, 0}, {2, 6, 1}, {0, 6, 0}}, 2}},
		{"7": {[]*PointTest{{3, 4, 0}, {3, 6, 1}, {1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 3}},
		{},
		{"7": {[]*PointTest{{0, 6, 0}, {0, 7, 0}}, 1}},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test39(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"6"}},
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"6"}},
		{"id": "3", "parents": []string{"7"}},
		{"id": "4", "parents": []string{"8"}},
		{"id": "5", "parents": []string{}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{}},
	}

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 1, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{"6": {[]*PointTest{{0, 0, 0}, {0, 6, 0}}, 0}},
		{"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1}},
		{"6": {[]*PointTest{{2, 2, 0}, {2, 6, 1}, {0, 6, 0}}, 2}},
		{"7": {[]*PointTest{{3, 3, 0}, {3, 6, 1}, {1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 3}},
		{"8": {[]*PointTest{{4, 4, 0}, {4, 6, 1}, {2, 6, 0}, {2, 7, 1}, {1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 4}},
		{},
		{"7": {[]*PointTest{{0, 6, 0}, {0, 7, 0}}, 0}},
		{"8": {[]*PointTest{{0, 7, 0}, {0, 8, 0}}, 0}},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test40(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"7"}},
		{"id": "1", "parents": []string{"5"}},
		{"id": "2", "parents": []string{"6"}},
		{"id": "3", "parents": []string{"8"}},
		{"id": "4", "parents": []string{"9"}},
		{"id": "5", "parents": []string{}},
		{"id": "6", "parents": []string{"7"}},
		{"id": "7", "parents": []string{"8"}},
		{"id": "8", "parents": []string{"9"}},
		{"id": "9", "parents": []string{}},
	}

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 1, 1, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{"7": {[]*PointTest{{0, 0, 0}, {0, 7, 0}}, 0}},
		{"5": {[]*PointTest{{1, 1, 0}, {1, 5, 0}}, 1}},
		{"6": {[]*PointTest{{2, 2, 0}, {2, 6, 1}, {1, 6, 0}}, 2}},
		{"8": {[]*PointTest{{3, 3, 0}, {3, 6, 1}, {2, 6, 0}, {2, 7, 1}, {1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, 3}},
		{"9": {[]*PointTest{{4, 4, 0}, {4, 6, 1}, {3, 6, 0}, {3, 7, 1}, {2, 7, 0}, {2, 8, 1}, {1, 8, 0}, {1, 9, 1}, {0, 9, 0}}, 4}},
		{},
		{"7": {[]*PointTest{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, 2}},
		{"8": {[]*PointTest{{0, 7, 0}, {0, 8, 0}}, 0}},
		{"9": {[]*PointTest{{0, 8, 0}, {0, 9, 0}}, 0}},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

// Test41 test the date-order bug where parent defined before node ends up with an infinite branch going down
func Test41(t *testing.T) {
	// Initial input
	inputNodes := []*Node{
		{"id": "0", "parents": []string{"2"}},
		{"id": "1", "parents": []string{"4", "3"}},
		{"id": "2", "parents": []string{"3", "1"}},
		{"id": "3", "parents": []string{"4"}},
		{"id": "4", "parents": []string{"5"}},
		{"id": "5", "parents": []string{}},
	}

	out, _ := buildTreeTest(inputNodes, customColors, "", -1)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 0}

	expectedPaths := []map[string]PathTest{
		{"2": {[]*PointTest{{0, 0, 0}, {0, 2, 0}}, 0}},
		{
			"3": {[]*PointTest{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {0, 3, 0}}, 2},
			"4": {[]*PointTest{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, 1},
		},
		{
			"3": {[]*PointTest{{0, 2, 0}, {0, 3, 0}}, 0},
			"1": {[]*PointTest{{0, 2, 0}, {3, 2, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {1, 5, 0}, {1, 6, 0}}, 3},
		},
		{"4": {[]*PointTest{{0, 3, 0}, {0, 4, 0}}, 0}},
		{"5": {[]*PointTest{{0, 4, 0}, {0, 5, 0}}, 0}},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func assertEq(t *testing.T, expected, actual any) {
	if actual != expected {
		t.Logf("Expected: %d, Actual: %d", expected, actual)
		t.Fail()
	}
}

func convertPoints[P IPoint](arr []P) (out []IPoint) {
	for _, v := range arr {
		out = append(out, v)
	}
	return
}

func TestPathHeight1(t *testing.T) {
	path := &Path{Points: convertPoints([]*PointTest{
		{0, 2, 0},
		{3, 2, 2},
		{3, 9, 1},
		{2, 9, 0},
		{2, 11, 1},
		{1, 11, 0},
	})}
	assertEq(t, -1, path.GetHeightAtIdx(1))
	assertEq(t, 3, path.GetHeightAtIdx(2))
	assertEq(t, 3, path.GetHeightAtIdx(3))
	assertEq(t, 2, path.GetHeightAtIdx(9))
	assertEq(t, 2, path.GetHeightAtIdx(10))
	assertEq(t, 1, path.GetHeightAtIdx(11))
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
		if _, err := buildTreeTest(inputNodes, customColors, "", -1); err != nil {
			b.Logf("Failed to build tree")
		}
	}
}

func TestCropPathAt(t *testing.T) {
	p := &Path{Points: convertPoints([]*PointTest{{0, 2, 0}, {3, 2, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {1, 5, 0}, {1, 6, 0}})}
	expected := convertPoints([]*PointTest{{2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {1, 5, 0}, {1, 6, 0}})
	np := cropPathAt(p, 4, 10)
	testPoints(t, expected, np.Points)

	p = &Path{Points: convertPoints([]*PointTest{{0, 0, 0}, {3, 0, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {1, 5, 0}, {1, 6, 0}})}
	expected = convertPoints([]*PointTest{{3, 1, 0}, {3, 3, 1}, {2, 3, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {1, 5, 0}, {1, 6, 0}})
	np = cropPathAt(p, 1, 10)
	testPoints(t, expected, np.Points)
}

func testPoints(t *testing.T, expected, points []IPoint) {
	for i, p := range points {
		if i >= len(expected) {
			t.Logf("fail1")
			t.Fail()
		} else if !p.Equal(expected[i]) {
			t.Logf("fail2 %v != %v", expected[i], p)
			t.Fail()
		}
	}
}

func TestExpandPath(t *testing.T) {
	path := &Path{Points: convertPoints([]*PointTest{{0, 0, 0}, {0, 3, 0}})}
	out := expandPath(path)
	expected := convertPoints([]*PointTest{{0, 0, 0}, {0, 1, 0}, {0, 2, 0}, {0, 3, 0}})
	testPoints(t, expected, out.Points)

	path = &Path{Points: convertPoints([]*PointTest{{0, 0, 0}, {1, 0, 1}, {1, 3, 0}})}
	out = expandPath(path)
	expected = convertPoints([]*PointTest{{0, 0, 0}, {1, 0, 1}, {1, 1, 0}, {1, 2, 0}, {1, 3, 0}})
	testPoints(t, expected, out.Points)
}
