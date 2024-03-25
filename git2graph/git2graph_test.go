package git2graph

import (
	"testing"
)

func validateColumns(t *testing.T, expectedColumns []int, data []map[string]any) {
	for idx, row := range data {
		if row["column"] != expectedColumns[idx] {
			t.Fail()
			t.Logf("ID: %s, Expected column: %d, Actual column: %d", row["id"], expectedColumns[idx], row["column"])
		}
	}
}

func validatePaths(t *testing.T, expectedPaths []map[string]Path, data []map[string]any) {
	for nodeIdx, node := range data {
		for _, parentID := range node["parents"].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			if len(node["parentsPaths"].(map[string]Path)[parentID].Path) != len(expectedPaths[nodeIdx][parentID].Path) {
				t.Fail()
				t.Logf("ID: %s, Expected nb paths: %d, Actual nb paths: %d", node["id"], len(expectedPaths[nodeIdx][parentID].Path), len(node["parentsPaths"].(map[string]Path)[parentID].Path))
				t.Logf("ID: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentID], node["parentsPaths"].(map[string]Path)[parentID])
				return
			}
			for pathIdx, pathNode := range node["parentsPaths"].(map[string]Path)[parentID].Path {
				if pathNode != expectedPaths[nodeIdx][parentID].Path[pathIdx] {
					t.Fail()
					t.Logf("ID: %s, Expected path: %d, Actual path: %d", node["id"], expectedPaths[nodeIdx][parentID].Path[pathIdx], pathNode)
					t.Logf("ID: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentID].Path, node["parentsPaths"].(map[string]Path)[parentID].Path)
				}
			}
		}
	}
}

func validateColors(t *testing.T, expectedPaths []map[string]Path, data []map[string]any) {
	for nodeIdx, node := range data {
		for _, parentID := range node["parents"].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			if expectedPaths[nodeIdx][parentID].Color != node["parentsPaths"].(map[string]Path)[parentID].Color {
				t.Logf("ID: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentID].Color, node["parentsPaths"].(map[string]Path)[parentID].Color)
				t.Fail()
			}
		}
	}
}

var customColors = []string{
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
}

func TestDebug(t *testing.T) {
	DebugMode = true
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{}})
	out, _ := BuildTree(inputNodes, customColors)
	// Ensure nodes have debug property
	if len(out[0]["debug"].([]string)) <= 0 {
		t.Fail()
	}
}

func TestNotEnoughColors(t *testing.T) {
	var colors = []string{
		"color1",
		"color2",
	}
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{}})
	out, _ := BuildTree(inputNodes, colors)
	if out[2]["color"] != "#000" {
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
		{"2": Path{"2", []Point{{0, 0, 0}, {0, 1, 0}}, "color1"}},
		{"3": Path{"3", []Point{{0, 1, 0}, {0, 2, 0}}, "color1"}},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"2": Path{"2", []Point{{0, 0, 0}, {0, 1, 0}}, "color1"},
		},
		{
			"3": Path{"3", []Point{{0, 1, 0}, {0, 2, 0}}, "color1"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, "color2"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, "color2"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{0, 2, 0}, {0, 4, 0}}, "color1"},
			"4": Path{"4", []Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{2, 3, 0}, {2, 4, 1}, {0, 4, 0}}, "color3"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": {"4", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"4": {"4", []Point{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{2, 2, 0}, {2, 3, 1}, {0, 3, 0}}, "color3"},
		},
		{
			"6": {"6", []Point{{0, 3, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, "color4"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, "color1"},
			"4": Path{"4", []Point{{1, 1, 0}, {1, 3, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{0, 2, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 3, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 3, 0}}, "color2"},
			"5": Path{"5", []Point{{1, 1, 0}, {2, 1, 2}, {2, 4, 1}, {0, 4, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{0, 2, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, "color2"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 3, 1}, {0, 3, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{0, 2, 0}, {0, 3, 0}}, "color1"},
			"5": Path{"5", []Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 1}, {1, 3, 0}, {1, 4, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{0, 3, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, "color3"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{0, 2, 0}, {0, 3, 0}}, "color1"},
			"7": Path{"7", []Point{{0, 2, 0}, {2, 2, 2}, {2, 4, 1}, {1, 4, 0}, {1, 6, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{0, 3, 0}, {0, 4, 0}}, "color1"},
			"6": Path{"6", []Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{0, 4, 0}, {0, 7, 0}}, "color1"},
		},
		{
			"8": Path{"8", []Point{{2, 5, 0}, {2, 7, 1}, {0, 7, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, "color3"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5", "3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
			"3": Path{"3", []Point{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{0, 3, 0}, {0, 7, 0}}, "color1"},
			"6": Path{"6", []Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 1}, {2, 4, 0}, {2, 5, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{1, 4, 0}, {2, 4, 2}, {2, 5, 0}}, "color4"},
			"7": Path{"7", []Point{{1, 4, 0}, {1, 6, 0}}, "color2"},
		},
		{
			"7": Path{"7", []Point{{2, 5, 0}, {2, 6, 1}, {1, 6, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, "color2"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {0, 1, 3}, {0, 2, 0}}, "color1"},
			"4": Path{"4", []Point{{1, 1, 0}, {1, 3, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{0, 2, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 3, 0}, {0, 3, 3}, {0, 4, 0}}, "color1"},
			"6": Path{"6", []Point{{1, 3, 0}, {1, 5, 1}, {0, 5, 0}}, "color2"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4", "2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5", "9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"7", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6", "8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 1, 0, 1, 1, 0, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
			"2": Path{"2", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"6", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
			"4": Path{"4", []Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{0, 2, 0}, {0, 5, 0}}, "color1"},
			"4": Path{"4", []Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{1, 4, 0}, {1, 7, 1}, {0, 7, 0}}, "color2"},
			"7": Path{"7", []Point{{1, 4, 0}, {2, 4, 2}, {2, 6, 0}}, "color4"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"6", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 2, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
			"4": Path{"4", []Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{0, 2, 0}, {0, 4, 0}}, "color1"},
			"4": Path{"4", []Point{{0, 2, 0}, {2, 2, 2}, {2, 3, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{2, 3, 0}, {2, 5, 1}, {1, 5, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{0, 4, 0}, {0, 7, 0}}, "color1"},
			"7": Path{"7", []Point{{0, 4, 0}, {3, 4, 2}, {3, 5, 1}, {2, 5, 0}, {2, 6, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{1, 5, 0}, {1, 7, 1}, {0, 7, 0}}, "color2"},
			"7": Path{"7", []Point{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, "color4"},
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
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": Path{"5", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, "color4"},
		},
		{
			"7": Path{"7", []Point{{0, 4, 0}, {0, 6, 0}}, "color1"},
		},
		{
			"7": Path{"7", []Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test17(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 6, 1}, {0, 6, 0}}, "color4"},
		},
		{
			"5": Path{"5", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{0, 5, 0}, {0, 6, 0}}, "color1"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test18(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{3, 3, 0}, {3, 4, 1}, {2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 6, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, "color3"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test19(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"11", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 4, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": Path{"5", []Point{{0, 0, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
		},
		{
			"9": Path{"9", []Point{{2, 2, 0}, {2, 9, 0}}, "color3"},
		},
		{
			"7": Path{"7", []Point{{3, 3, 0}, {3, 7, 0}}, "color4"},
		},
		{
			"11": Path{"11", []Point{{1, 4, 0}, {1, 11, 0}}, "color2"},
			"6":  Path{"6", []Point{{1, 4, 0}, {4, 4, 2}, {4, 6, 0}}, "color5"},
		},
		{
			"8": Path{"8", []Point{{0, 5, 0}, {0, 8, 0}}, "color1"},
			"6": Path{"6", []Point{{0, 5, 0}, {4, 5, 2}, {4, 6, 0}}, "color5"},
		},
		{
			"11": Path{"11", []Point{{4, 6, 0}, {4, 8, 1}, {3, 8, 0}, {3, 10, 1}, {2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, "color5"},
		},
		{
			"8": Path{"8", []Point{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, "color4"},
		},
		{
			"10": Path{"10", []Point{{0, 8, 0}, {0, 10, 0}}, "color1"},
		},
		{
			"10": Path{"10", []Point{{2, 9, 0}, {2, 10, 1}, {0, 10, 0}}, "color3"},
		},
		{
			"12": Path{"12", []Point{{0, 10, 0}, {0, 12, 0}}, "color1"},
		},
		{
			"12": Path{"12", []Point{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test20(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 1}, {0, 5, 0}}, "color3"},
		},
		{
			"4": Path{"4", []Point{{3, 3, 0}, {3, 4, 1}, {0, 4, 0}}, "color4"},
		},
		{
			"5": Path{"5", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test21(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {1, 3, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 5, 1}, {0, 5, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{1, 3, 0}, {1, 6, 0}}, "color2"},
			"5": Path{"5", []Point{{1, 3, 0}, {2, 3, 2}, {2, 5, 1}, {0, 5, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"7": Path{"7", []Point{{0, 5, 0}, {0, 7, 0}}, "color1"},
		},
		{
			"7": Path{"7", []Point{{1, 6, 0}, {1, 7, 1}, {0, 7, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test22(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": Path{"5", []Point{{0, 0, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
		},
		{
			"6": Path{"6", []Point{{2, 2, 0}, {2, 6, 1}, {0, 6, 0}}, "color3"},
		},
		{
			"7": Path{"7", []Point{{3, 3, 0}, {3, 6, 1}, {2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, "color4"},
		},
		{
			"7": Path{"7", []Point{{1, 4, 0}, {1, 7, 0}}, "color2"},
			"6": Path{"6", []Point{{1, 4, 0}, {2, 4, 2}, {2, 6, 1}, {0, 6, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{0, 5, 0}, {0, 6, 0}}, "color1"},
		},
		{
			"8": Path{"8", []Point{{0, 6, 0}, {0, 8, 0}}, "color1"},
		},
		{
			"8": Path{"8", []Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test23(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{2, 2, 0}, {2, 4, 1}, {0, 4, 0}}, "color3"},
		},
		{
			"7": Path{"7", []Point{{3, 3, 0}, {3, 4, 1}, {1, 4, 0}, {1, 7, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 6, 0}}, "color1"},
			"5": Path{"5", []Point{{0, 4, 0}, {2, 4, 2}, {2, 5, 0}}, "color5"},
		},
		{
			"6": Path{"6", []Point{{2, 5, 0}, {2, 6, 1}, {0, 6, 0}}, "color5"},
		},
		{
			"8": Path{"8", []Point{{0, 6, 0}, {0, 8, 0}}, "color1"},
		},
		{
			"8": Path{"8", []Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, "color4"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test24(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"10", "8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"11", "8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 4, 1, 1, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
		},
		{
			"9": Path{"9", []Point{{2, 2, 0}, {2, 9, 0}}, "color3"},
		},
		{
			"7": Path{"7", []Point{{0, 3, 0}, {0, 7, 0}}, "color1"},
			"6": Path{"6", []Point{{0, 3, 0}, {3, 3, 2}, {3, 6, 1}, {1, 6, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{4, 4, 0}, {4, 6, 1}, {1, 6, 0}}, "color5"},
		},
		{
			"6": Path{"6", []Point{{1, 5, 0}, {1, 6, 0}}, "color2"},
		},
		{
			"10": Path{"10", []Point{{1, 6, 0}, {1, 10, 0}}, "color2"},
			"8":  Path{"8", []Point{{1, 6, 0}, {3, 6, 2}, {3, 8, 0}}, "color6"},
		},
		{
			"11": Path{"11", []Point{{0, 7, 0}, {0, 11, 0}}, "color1"},
			"8":  Path{"8", []Point{{0, 7, 0}, {3, 7, 2}, {3, 8, 0}}, "color6"},
		},
		{
			"9": Path{"9", []Point{{3, 8, 0}, {3, 9, 1}, {2, 9, 0}}, "color6"},
		},
		{
			"10": Path{"10", []Point{{2, 9, 0}, {2, 10, 1}, {1, 10, 0}}, "color3"},
		},
		{
			"11": Path{"11", []Point{{1, 10, 0}, {1, 11, 1}, {0, 11, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test25(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"9", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"9", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"12", "9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"11", "10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 2, 0, 2, 3, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"5": Path{"5", []Point{{0, 0, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"3": Path{"3", []Point{{1, 1, 0}, {1, 3, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{2, 2, 0}, {2, 4, 0}}, "color3"},
		},
		{
			"9": Path{"9", []Point{{1, 3, 0}, {1, 9, 0}}, "color2"},
			"7": Path{"7", []Point{{1, 3, 0}, {3, 3, 2}, {3, 7, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{2, 4, 0}, {2, 6, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{0, 5, 0}, {0, 8, 0}}, "color1"},
			"7": Path{"7", []Point{{0, 5, 0}, {3, 5, 2}, {3, 7, 0}}, "color4"},
		},
		{
			"9": Path{"9", []Point{{2, 6, 0}, {2, 9, 1}, {1, 9, 0}}, "color3"},
			"7": Path{"7", []Point{{2, 6, 0}, {3, 6, 2}, {3, 7, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{3, 7, 0}, {3, 8, 1}, {0, 8, 0}}, "color4"},
		},
		{
			"12": Path{"12", []Point{{0, 8, 0}, {0, 12, 0}}, "color1"},
			"9":  Path{"9", []Point{{0, 8, 0}, {1, 8, 2}, {1, 9, 0}}, "color2"},
		},
		{
			"11": Path{"11", []Point{{1, 9, 0}, {1, 11, 0}}, "color2"},
			"10": Path{"10", []Point{{1, 9, 0}, {2, 9, 2}, {2, 10, 0}}, "color5"},
		},
		{
			"11": Path{"11", []Point{{2, 10, 0}, {2, 11, 1}, {1, 11, 0}}, "color5"},
		},
		{
			"12": Path{"12", []Point{{1, 11, 0}, {1, 12, 1}, {0, 12, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test26(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"8", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 5, 1}, {1, 5, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{0, 3, 0}, {0, 8, 0}}, "color1"},
			"5": Path{"5", []Point{{0, 3, 0}, {2, 3, 2}, {2, 5, 1}, {1, 5, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{1, 4, 0}, {1, 5, 0}}, "color2"},
		},
		{
			"7": Path{"7", []Point{{1, 5, 0}, {1, 7, 0}}, "color2"},
			"6": Path{"6", []Point{{1, 5, 0}, {2, 5, 2}, {2, 6, 0}}, "color4"},
		},
		{
			"7": Path{"7", []Point{{2, 6, 0}, {2, 7, 1}, {1, 7, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{1, 7, 0}, {1, 8, 1}, {0, 8, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test27(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"15", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"15"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"9", "13"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"14", "10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "13", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "14", "parents": []string{"16"}})
	inputNodes = append(inputNodes, map[string]any{"id": "15", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]any{"id": "16", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]any{"id": "17", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 4, 2, 1, 1, 6, 3, 2, 5, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
		},
		{
			"7": Path{"7", []Point{{2, 2, 0}, {2, 7, 0}}, "color3"},
		},
		{
			"11": Path{"11", []Point{{3, 3, 0}, {3, 11, 0}}, "color4"},
		},
		{
			"15": Path{"15", []Point{{0, 4, 0}, {0, 15, 0}}, "color1"},
			"6":  Path{"6", []Point{{0, 4, 0}, {4, 4, 2}, {4, 6, 0}}, "color5"},
		},
		{
			"8": Path{"8", []Point{{1, 5, 0}, {1, 8, 0}}, "color2"},
			"6": Path{"6", []Point{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, "color5"},
		},
		{
			"15": Path{"15", []Point{{4, 6, 0}, {4, 14, 1}, {2, 14, 0}, {2, 15, 1}, {0, 15, 0}}, "color5"},
		},
		{
			"12": Path{"12", []Point{{2, 7, 0}, {2, 12, 0}}, "color3"},
		},
		{
			"9":  Path{"9", []Point{{1, 8, 0}, {1, 9, 0}}, "color2"},
			"13": Path{"13", []Point{{1, 8, 0}, {5, 8, 2}, {5, 13, 0}}, "color6"},
		},
		{
			"14": Path{"14", []Point{{1, 9, 0}, {1, 14, 0}}, "color2"},
			"10": Path{"10", []Point{{1, 9, 0}, {6, 9, 2}, {6, 10, 0}}, "color7"},
		},
		{
			"14": Path{"14", []Point{{6, 10, 0}, {6, 14, 1}, {1, 14, 0}}, "color7"},
		},
		{
			"14": Path{"14", []Point{{3, 11, 0}, {3, 14, 1}, {1, 14, 0}}, "color4"},
		},
		{
			"14": Path{"14", []Point{{2, 12, 0}, {2, 14, 1}, {1, 14, 0}}, "color3"},
		},
		{
			"14": Path{"14", []Point{{5, 13, 0}, {5, 14, 1}, {1, 14, 0}}, "color6"},
		},
		{
			"16": Path{"16", []Point{{1, 14, 0}, {1, 16, 0}}, "color2"},
		},
		{
			"17": Path{"17", []Point{{0, 15, 0}, {0, 17, 0}}, "color1"},
		},
		{
			"17": Path{"17", []Point{{1, 16, 0}, {1, 17, 1}, {0, 17, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test28(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"2", "1"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"2": Path{"2", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
			"1": Path{"1", []Point{{0, 0, 0}, {1, 0, 2}, {1, 1, 0}}, "color2"},
		},
		{
			"2": Path{"2", []Point{{1, 1, 0}, {1, 2, 1}, {0, 2, 0}}, "color2"},
		},
		{
			"3": Path{"3", []Point{{0, 2, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{0, 3, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 6, 0}}, "color1"},
			"5": Path{"5", []Point{{0, 4, 0}, {1, 4, 2}, {1, 5, 0}}, "color2"},
		},
		{
			"6": Path{"6", []Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test29(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"15"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"18"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"20"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"9", "10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"9", "11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"13", "14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"21"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"13"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "13", "parents": []string{"16", "15"}})
	inputNodes = append(inputNodes, map[string]any{"id": "14", "parents": []string{"19"}})
	inputNodes = append(inputNodes, map[string]any{"id": "15", "parents": []string{"26"}})
	inputNodes = append(inputNodes, map[string]any{"id": "16", "parents": []string{"27"}})
	inputNodes = append(inputNodes, map[string]any{"id": "17", "parents": []string{"25"}})
	inputNodes = append(inputNodes, map[string]any{"id": "18", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]any{"id": "19", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]any{"id": "20", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]any{"id": "21", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]any{"id": "22", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]any{"id": "23", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]any{"id": "24", "parents": []string{"25"}})
	inputNodes = append(inputNodes, map[string]any{"id": "25", "parents": []string{"26"}})
	inputNodes = append(inputNodes, map[string]any{"id": "26", "parents": []string{"27"}})
	inputNodes = append(inputNodes, map[string]any{"id": "27", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 5, 6, 0, 3, 0, 7, 3, 5, 0, 4, 1, 0, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"7": Path{"7", []Point{{0, 0, 0}, {0, 7, 0}}, "color1"},
		},
		{
			"15": Path{"15", []Point{{1, 1, 0}, {1, 15, 0}}, "color2"},
		},
		{
			"17": Path{"17", []Point{{2, 2, 0}, {2, 17, 0}}, "color3"},
		},
		{
			"8": Path{"8", []Point{{3, 3, 0}, {3, 8, 0}}, "color4"},
		},
		{
			"18": Path{"18", []Point{{4, 4, 0}, {4, 13, 1}, {3, 13, 0}, {3, 18, 0}}, "color5"},
		},
		{
			"12": Path{"12", []Point{{5, 5, 0}, {5, 12, 0}}, "color6"},
		},
		{
			"20": Path{"20", []Point{{6, 6, 0}, {6, 13, 1}, {5, 13, 0}, {5, 20, 0}}, "color7"},
		},
		{
			"9":  Path{"9", []Point{{0, 7, 0}, {0, 9, 0}}, "color1"},
			"10": Path{"10", []Point{{0, 7, 0}, {7, 7, 2}, {7, 10, 0}}, "color8"},
		},
		{
			"9":  Path{"9", []Point{{3, 8, 0}, {0, 8, 3}, {0, 9, 0}}, "color1"},
			"11": Path{"11", []Point{{3, 8, 0}, {3, 11, 0}}, "color4"},
		},
		{
			"13": Path{"13", []Point{{0, 9, 0}, {0, 13, 0}}, "color1"},
			"14": Path{"14", []Point{{0, 9, 0}, {8, 9, 2}, {8, 13, 1}, {7, 13, 0}, {7, 14, 1}, {4, 14, 0}}, "color9"},
		},
		{
			"21": Path{"21", []Point{{7, 10, 0}, {7, 13, 1}, {6, 13, 0}, {6, 21, 0}}, "color8"},
		},
		{
			"13": Path{"13", []Point{{3, 11, 0}, {3, 13, 1}, {0, 13, 0}}, "color4"},
		},
		{
			"14": Path{"14", []Point{{5, 12, 0}, {5, 13, 1}, {4, 13, 0}, {4, 14, 0}}, "color6"},
		},
		{
			"16": Path{"16", []Point{{0, 13, 0}, {0, 16, 0}}, "color1"},
			"15": Path{"15", []Point{{0, 13, 0}, {1, 13, 2}, {1, 15, 0}}, "color2"},
		},
		{
			"19": Path{"19", []Point{{4, 14, 0}, {4, 19, 0}}, "color6"},
		},
		{
			"26": Path{"26", []Point{{1, 15, 0}, {1, 26, 0}}, "color2"},
		},
		{
			"27": Path{"27", []Point{{0, 16, 0}, {0, 27, 0}}, "color1"},
		},
		{
			"25": Path{"25", []Point{{2, 17, 0}, {2, 25, 0}}, "color3"},
		},
		{
			"24": Path{"24", []Point{{3, 18, 0}, {3, 24, 0}}, "color5"},
		},
		{
			"23": Path{"23", []Point{{4, 19, 0}, {4, 23, 0}}, "color6"},
		},
		{
			"22": Path{"22", []Point{{5, 20, 0}, {5, 22, 0}}, "color7"},
		},
		{
			"22": Path{"22", []Point{{6, 21, 0}, {6, 22, 1}, {5, 22, 0}}, "color8"},
		},
		{
			"23": Path{"23", []Point{{5, 22, 0}, {5, 23, 1}, {4, 23, 0}}, "color7"},
		},
		{
			"24": Path{"24", []Point{{4, 23, 0}, {4, 24, 1}, {3, 24, 0}}, "color6"},
		},
		{
			"25": Path{"25", []Point{{3, 24, 0}, {3, 25, 1}, {2, 25, 0}}, "color5"},
		},
		{
			"26": Path{"26", []Point{{2, 25, 0}, {2, 26, 1}, {1, 26, 0}}, "color3"},
		},
		{
			"27": Path{"27", []Point{{1, 26, 0}, {1, 27, 1}, {0, 27, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test30(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"15"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5", "10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6", "7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"16", "8"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"16"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"18"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{"13", "16"}})
	inputNodes = append(inputNodes, map[string]any{"id": "13", "parents": []string{"21"}})
	inputNodes = append(inputNodes, map[string]any{"id": "14", "parents": []string{"18"}})
	inputNodes = append(inputNodes, map[string]any{"id": "15", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]any{"id": "16", "parents": []string{"18", "17"}})
	inputNodes = append(inputNodes, map[string]any{"id": "17", "parents": []string{"18"}})
	inputNodes = append(inputNodes, map[string]any{"id": "18", "parents": []string{"20", "19"}})
	inputNodes = append(inputNodes, map[string]any{"id": "19", "parents": []string{"20"}})
	inputNodes = append(inputNodes, map[string]any{"id": "20", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]any{"id": "21", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]any{"id": "22", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]any{"id": "23", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]any{"id": "24", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0, 5, 6, 5, 4, 6, 6, 6, 2, 1, 0, 6, 0, 4, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"4": Path{"4", []Point{{0, 0, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"15": Path{"15", []Point{{1, 1, 0}, {1, 15, 0}}, "color2"},
		},
		{
			"14": Path{"14", []Point{{2, 2, 0}, {2, 14, 0}}, "color3"},
		},
		{
			"22": Path{"22", []Point{{3, 3, 0}, {3, 18, 1}, {2, 18, 0}, {2, 22, 0}}, "color4"},
		},
		{
			"5":  Path{"5", []Point{{0, 4, 0}, {0, 5, 0}}, "color1"},
			"10": Path{"10", []Point{{0, 4, 0}, {4, 4, 2}, {4, 10, 0}}, "color5"},
		},
		{
			"6": Path{"6", []Point{{0, 5, 0}, {0, 6, 0}}, "color1"},
			"7": Path{"7", []Point{{0, 5, 0}, {5, 5, 2}, {5, 7, 0}}, "color6"},
		},
		{
			"16": Path{"16", []Point{{0, 6, 0}, {0, 16, 0}}, "color1"},
			"8":  Path{"8", []Point{{0, 6, 0}, {6, 6, 2}, {6, 8, 0}}, "color7"},
		},
		{
			"9": Path{"9", []Point{{5, 7, 0}, {5, 9, 0}}, "color6"},
		},
		{
			"11": Path{"11", []Point{{6, 8, 0}, {6, 11, 0}}, "color7"},
		},
		{
			"16": Path{"16", []Point{{5, 9, 0}, {5, 16, 1}, {0, 16, 0}}, "color6"},
		},
		{
			"18": Path{"18", []Point{{4, 10, 0}, {4, 18, 1}, {0, 18, 0}}, "color5"},
		},
		{
			"12": Path{"12", []Point{{6, 11, 0}, {6, 12, 0}}, "color7"},
		},
		{
			"13": Path{"13", []Point{{6, 12, 0}, {6, 13, 0}}, "color7"},
			"16": Path{"16", []Point{{6, 12, 0}, {0, 12, 3}, {0, 16, 0}}, "color1"},
		},
		{
			"21": Path{"21", []Point{{6, 13, 0}, {6, 16, 1}, {5, 16, 0}, {5, 18, 1}, {3, 18, 0}, {3, 21, 0}}, "color7"},
		},
		{
			"18": Path{"18", []Point{{2, 14, 0}, {2, 18, 1}, {0, 18, 0}}, "color3"},
		},
		{
			"23": Path{"23", []Point{{1, 15, 0}, {1, 23, 0}}, "color2"},
		},
		{
			"18": Path{"18", []Point{{0, 16, 0}, {0, 18, 0}}, "color1"},
			"17": Path{"17", []Point{{0, 16, 0}, {6, 16, 2}, {6, 17, 0}}, "color8"},
		},
		{
			"18": Path{"18", []Point{{6, 17, 0}, {6, 18, 1}, {0, 18, 0}}, "color8"},
		},
		{
			"20": Path{"20", []Point{{0, 18, 0}, {0, 20, 0}}, "color1"},
			"19": Path{"19", []Point{{0, 18, 0}, {4, 18, 2}, {4, 19, 0}}, "color6"},
		},
		{
			"20": Path{"20", []Point{{4, 19, 0}, {4, 20, 1}, {0, 20, 0}}, "color6"},
		},
		{
			"24": Path{"24", []Point{{0, 20, 0}, {0, 24, 0}}, "color1"},
		},
		{
			"22": Path{"22", []Point{{3, 21, 0}, {3, 22, 1}, {2, 22, 0}}, "color7"},
		},
		{
			"23": Path{"23", []Point{{2, 22, 0}, {2, 23, 1}, {1, 23, 0}}, "color4"},
		},
		{
			"24": Path{"24", []Point{{1, 23, 0}, {1, 24, 1}, {0, 24, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test31(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"5": Path{"5", []Point{{2, 2, 0}, {2, 4, 1}, {1, 4, 0}, {1, 5, 0}}, "color3"},
			"4": Path{"4", []Point{{2, 2, 0}, {1, 2, 3}, {1, 4, 1}, {0, 4, 0}}, "color2"},
		},
		{
			"4": Path{"4", []Point{{0, 3, 0}, {0, 4, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{0, 4, 0}, {0, 6, 0}}, "color1"},
		},
		{
			"6": Path{"6", []Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, "color3"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test32(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"5", "3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"3", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"2": Path{"2", []Point{{0, 0, 0}, {0, 2, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
			"3": Path{"3", []Point{{1, 1, 0}, {2, 1, 2}, {2, 3, 1}, {0, 3, 0}}, "color3"},
		},
		{
			"3": Path{"3", []Point{{0, 2, 0}, {0, 3, 0}}, "color1"},
			"4": Path{"4", []Point{{0, 2, 0}, {3, 2, 2}, {3, 3, 1}, {2, 3, 0}, {2, 4, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{0, 3, 0}, {0, 6, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{2, 4, 0}, {2, 5, 1}, {1, 5, 0}}, "color4"},
		},
		{
			"6": Path{"6", []Point{{1, 5, 0}, {1, 6, 1}, {0, 6, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test33(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]any, 0)
	inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"9", "4"}})
	inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"13"}})
	inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"12", "11"}})
	inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{"13"}})
	inputNodes = append(inputNodes, map[string]any{"id": "13", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 3, 1, 4, 2, 1, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		{
			"3": Path{"3", []Point{{0, 0, 0}, {0, 3, 0}}, "color1"},
		},
		{
			"5": Path{"5", []Point{{1, 1, 0}, {1, 5, 0}}, "color2"},
		},
		{
			"7": Path{"7", []Point{{2, 2, 0}, {2, 7, 0}}, "color3"},
		},
		{
			"9": Path{"9", []Point{{0, 3, 0}, {0, 9, 0}}, "color1"},
			"4": Path{"4", []Point{{0, 3, 0}, {3, 3, 2}, {3, 4, 0}}, "color4"},
		},
		{
			"9": Path{"9", []Point{{3, 4, 0}, {3, 9, 1}, {0, 9, 0}}, "color4"},
		},
		{
			"8": Path{"8", []Point{{1, 5, 0}, {1, 8, 0}}, "color2"},
			"6": Path{"6", []Point{{1, 5, 0}, {4, 5, 2}, {4, 6, 0}}, "color5"},
		},
		{
			"10": Path{"10", []Point{{4, 6, 0}, {4, 9, 1}, {3, 9, 0}, {3, 10, 1}, {1, 10, 0}}, "color5"},
		},
		{
			"10": Path{"10", []Point{{2, 7, 0}, {2, 10, 1}, {1, 10, 0}}, "color3"},
		},
		{
			"10": Path{"10", []Point{{1, 8, 0}, {1, 10, 0}}, "color2"},
		},
		{
			"13": Path{"13", []Point{{0, 9, 0}, {0, 13, 0}}, "color1"},
		},
		{
			"12": Path{"12", []Point{{1, 10, 0}, {1, 12, 0}}, "color2"},
			"11": Path{"11", []Point{{1, 10, 0}, {2, 10, 2}, {2, 11, 0}}, "color6"},
		},
		{
			"12": Path{"12", []Point{{2, 11, 0}, {2, 12, 1}, {1, 12, 0}}, "color6"},
		},
		{
			"13": Path{"13", []Point{{1, 12, 0}, {1, 13, 1}, {0, 13, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test34(t *testing.T) {
	// Initial input
	inputNodes := []map[string]any{
		{"id": "0", "parents": []string{"5"}},
		{"id": "1", "parents": []string{"4", "2"}},
		{"id": "2", "parents": []string{"3"}},
		{"id": "3", "parents": []string{"5", "4"}},
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
			"5": Path{"5", []Point{{0, 0, 0}, {0, 5, 0}}, "color1"},
		},
		{
			"4": Path{"4", []Point{{1, 1, 0}, {1, 4, 0}}, "color2"},
			"2": Path{"2", []Point{{1, 1, 0}, {2, 1, 2}, {2, 2, 0}}, "color3"},
		},
		{
			"3": Path{"3", []Point{{2, 2, 0}, {2, 3, 0}}, "color3"},
		},
		{
			"5": Path{"5", []Point{{2, 3, 0}, {0, 3, 3}, {0, 5, 0}}, "color1"},
			"4": Path{"4", []Point{{2, 3, 0}, {2, 4, 1}, {1, 4, 0}}, "color3"},
		},
		{
			"6": Path{"6", []Point{{1, 4, 0}, {1, 6, 0}}, "color2"},
		},
		{
			"7": Path{"7", []Point{{0, 5, 0}, {0, 7, 0}}, "color1"},
		},
		{
			"7": Path{"7", []Point{{1, 6, 0}, {0, 7, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func TestPathHeight1(t *testing.T) {
	index := make(map[string]*OutputNode)
	out := OutputNode{parentsPaths: map[string]Path{"1": {Path: []Point{
		{X: 0, Y: 2, Type: 0},
		{X: 3, Y: 2, Type: 2},
		{X: 3, Y: 9, Type: 1},
		{X: 2, Y: 9, Type: 0},
		{X: 2, Y: 11, Type: 1},
		{X: 1, Y: 11, Type: 0},
	}}}}
	if out.GetPathHeightAtIdx(index, "1", 1) != -1 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 2) != 3 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 3) != 3 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 9) != 2 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 10) != 2 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 11) != 1 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx(index, "1", 1000) != -1 {
		t.Fail()
	}
}

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inputNodes := make([]map[string]any, 0)
		inputNodes = append(inputNodes, map[string]any{"id": "0", "parents": []string{"4"}})
		inputNodes = append(inputNodes, map[string]any{"id": "1", "parents": []string{"15"}})
		inputNodes = append(inputNodes, map[string]any{"id": "2", "parents": []string{"14"}})
		inputNodes = append(inputNodes, map[string]any{"id": "3", "parents": []string{"22"}})
		inputNodes = append(inputNodes, map[string]any{"id": "4", "parents": []string{"5", "10"}})
		inputNodes = append(inputNodes, map[string]any{"id": "5", "parents": []string{"6", "7"}})
		inputNodes = append(inputNodes, map[string]any{"id": "6", "parents": []string{"16", "8"}})
		inputNodes = append(inputNodes, map[string]any{"id": "7", "parents": []string{"9"}})
		inputNodes = append(inputNodes, map[string]any{"id": "8", "parents": []string{"11"}})
		inputNodes = append(inputNodes, map[string]any{"id": "9", "parents": []string{"16"}})
		inputNodes = append(inputNodes, map[string]any{"id": "10", "parents": []string{"18"}})
		inputNodes = append(inputNodes, map[string]any{"id": "11", "parents": []string{"12"}})
		inputNodes = append(inputNodes, map[string]any{"id": "12", "parents": []string{"13", "16"}})
		inputNodes = append(inputNodes, map[string]any{"id": "13", "parents": []string{"21"}})
		inputNodes = append(inputNodes, map[string]any{"id": "14", "parents": []string{"18"}})
		inputNodes = append(inputNodes, map[string]any{"id": "15", "parents": []string{"23"}})
		inputNodes = append(inputNodes, map[string]any{"id": "16", "parents": []string{"18", "17"}})
		inputNodes = append(inputNodes, map[string]any{"id": "17", "parents": []string{"18"}})
		inputNodes = append(inputNodes, map[string]any{"id": "18", "parents": []string{"20", "19"}})
		inputNodes = append(inputNodes, map[string]any{"id": "19", "parents": []string{"20"}})
		inputNodes = append(inputNodes, map[string]any{"id": "20", "parents": []string{"24"}})
		inputNodes = append(inputNodes, map[string]any{"id": "21", "parents": []string{"22"}})
		inputNodes = append(inputNodes, map[string]any{"id": "22", "parents": []string{"23"}})
		inputNodes = append(inputNodes, map[string]any{"id": "23", "parents": []string{"24"}})
		inputNodes = append(inputNodes, map[string]any{"id": "24", "parents": []string{}})
		if _, err := BuildTree(inputNodes, customColors); err != nil {
			b.Logf("Failed to build tree")
		}
	}
}
