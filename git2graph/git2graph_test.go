package git2graph

import (
	"testing"
)

func validateColumns(t *testing.T, expectedColumns []int, data []map[string]interface{}) {
	for idx, row := range data {
		if row["column"] != expectedColumns[idx] {
			t.Fail()
			t.Logf("Id: %s, Expected column: %d, Actual column: %d", row["id"], expectedColumns[idx], row["column"])
		}
	}
}

func validatePaths(t *testing.T, expectedPaths []map[string]Path, data []map[string]interface{}) {
	for nodeIdx, node := range data {
		for _, parentId := range node["parents"].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			if len(node["parentsPaths"].(map[string]Path)[parentId].Path) != len(expectedPaths[nodeIdx][parentId].Path) {
				t.Fail()
				t.Logf("Id: %s, Expected nb paths: %d, Actual nb paths: %d", node["id"], len(expectedPaths[nodeIdx][parentId].Path), len(node["parentsPaths"].(map[string]Path)[parentId].Path))
				t.Logf("Id: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentId], node["parentsPaths"].(map[string]Path)[parentId])
				return
			}
			for pathIdx, pathNode := range node["parentsPaths"].(map[string]Path)[parentId].Path {
				if pathNode != expectedPaths[nodeIdx][parentId].Path[pathIdx] {
					t.Fail()
					t.Logf("Id: %s, Expected path: %d, Actual path: %d", node["id"], expectedPaths[nodeIdx][parentId].Path[pathIdx], pathNode)
					t.Logf("Id: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentId].Path, node["parentsPaths"].(map[string]Path)[parentId].Path)
				}
			}
		}
	}
}

func validateColors(t *testing.T, expectedPaths []map[string]Path, data []map[string]interface{}) {
	for nodeIdx, node := range data {
		for _, parentId := range node["parents"].([]string) {
			if len(expectedPaths)-1 < nodeIdx {
				continue
			}
			if expectedPaths[nodeIdx][parentId].Color != node["parentsPaths"].(map[string]Path)[parentId].Color {
				t.Logf("Id: %s, Expected: %v, Actual: %v", node["id"], expectedPaths[nodeIdx][parentId].Color, node["parentsPaths"].(map[string]Path)[parentId].Color)
				t.Fail()
			}
		}
	}
}

var customColors []Color = []Color{
	Color{-2, "color1", false},
	Color{-2, "color2", false},
	Color{-2, "color3", false},
	Color{-2, "color4", false},
	Color{-2, "color5", false},
	Color{-2, "color6", false},
	Color{-2, "color7", false},
	Color{-2, "color8", false},
}

// 1
// |
// 2
// |
// 3
func Test1(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{0, 1, 0}}, "color1"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 1, 0}, Point{0, 2, 0}}, "color1"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{1, 2, 1}, Point{0, 2, 0}}, "color2"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{1, 2, 1}, Point{0, 2, 0}}, "color2"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{0, 4, 0}}, "color1"},
			"4": Path{"4", []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 3, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 3, 0}, Point{2, 4, 1}, Point{0, 4, 0}}, "color3"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 3, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 3, 1}, Point{0, 3, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{2, 2, 0}, Point{2, 3, 1}, Point{0, 3, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 3, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 4, 0}, Point{1, 5, 1}, Point{0, 5, 0}}, "color4"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{0, 1, 3}, Point{0, 2, 0}}, "color1"},
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 3, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 3, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 3, 0}}, "color2"},
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{2, 1, 2}, Point{2, 4, 1}, Point{0, 4, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 3, 0}, Point{1, 5, 1}, Point{0, 5, 0}}, "color2"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 4, 0}, Point{0, 5, 0}}, "color1"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"4", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 3, 1}, Point{0, 3, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 2, 0}, Point{0, 3, 0}}, "color1"},
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 3, 1}, Point{1, 3, 0}, Point{1, 4, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 3, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 4, 0}, Point{1, 5, 1}, Point{0, 5, 0}}, "color3"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"4", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 2, 0}, Point{0, 3, 0}}, "color1"},
			"7": Path{"7", []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 4, 1}, Point{1, 4, 0}, Point{1, 6, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 3, 0}, Point{0, 4, 0}}, "color1"},
			"6": Path{"6", []Point{Point{0, 3, 0}, Point{3, 3, 2}, Point{3, 4, 1}, Point{2, 4, 0}, Point{2, 5, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 4, 0}, Point{0, 7, 0}}, "color1"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{2, 5, 0}, Point{2, 7, 1}, Point{0, 7, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 6, 0}, Point{1, 7, 1}, Point{0, 7, 0}}, "color3"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5", "3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 3, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 4, 0}}, "color2"},
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{2, 1, 2}, Point{2, 2, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{1, 4, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 3, 0}, Point{0, 7, 0}}, "color1"},
			"6": Path{"6", []Point{Point{0, 3, 0}, Point{3, 3, 2}, Point{3, 4, 1}, Point{2, 4, 0}, Point{2, 5, 0}}, "color4"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 4, 0}, Point{2, 4, 2}, Point{2, 5, 0}}, "color4"},
			"7": Path{"7", []Point{Point{1, 4, 0}, Point{1, 6, 0}}, "color2"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{2, 5, 0}, Point{2, 6, 1}, Point{1, 6, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 6, 0}, Point{1, 7, 1}, Point{0, 7, 0}}, "color2"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 1, 0, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{0, 1, 3}, Point{0, 2, 0}}, "color1"},
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 3, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 3, 0}, Point{0, 3, 3}, Point{0, 4, 0}}, "color1"},
			"6": Path{"6", []Point{Point{1, 3, 0}, Point{1, 5, 1}, Point{0, 5, 0}}, "color2"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 4, 0}, Point{0, 5, 0}}, "color1"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4", "2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5", "9"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"7", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6", "8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 1, 0, 1, 1, 0, 0, 2, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 3, 0}}, "color1"},
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"6", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 1, 0, 2, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 4, 0}}, "color2"},
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{2, 1, 2}, Point{2, 3, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 2, 0}, Point{0, 5, 0}}, "color1"},
			"4": Path{"4", []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 3, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 3, 0}, Point{2, 4, 1}, Point{1, 4, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 4, 0}, Point{1, 7, 1}, Point{0, 7, 0}}, "color2"},
			"7": Path{"7", []Point{Point{1, 4, 0}, Point{2, 4, 2}, Point{2, 6, 0}}, "color4"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"6", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5", "4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0, 1, 2, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 1, 0}, Point{1, 5, 0}}, "color2"},
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{2, 1, 2}, Point{2, 3, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 2, 0}, Point{0, 4, 0}}, "color1"},
			"4": Path{"4", []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 3, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{2, 3, 0}, Point{2, 5, 1}, Point{1, 5, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 4, 0}, Point{0, 7, 0}}, "color1"},
			"7": Path{"7", []Point{Point{0, 4, 0}, Point{3, 4, 2}, Point{3, 5, 1}, Point{2, 5, 0}, Point{2, 6, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 5, 0}, Point{1, 7, 1}, Point{0, 7, 0}}, "color2"},
			"7": Path{"7", []Point{Point{1, 5, 0}, Point{2, 5, 2}, Point{2, 6, 0}}, "color4"},
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
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 1, 0}, Point{1, 5, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{0, 4, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{3, 3, 0}, Point{3, 4, 1}, Point{2, 4, 0}, Point{2, 5, 1}, Point{1, 5, 0}}, "color4"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{0, 4, 0}, Point{0, 6, 0}}, "color1"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{1, 5, 0}, Point{1, 6, 1}, Point{0, 6, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test17(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{0, 4, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{3, 3, 0}, Point{3, 4, 1}, Point{1, 4, 0}, Point{1, 6, 1}, Point{0, 6, 0}}, "color4"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 4, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 5, 0}, Point{0, 6, 0}}, "color1"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test18(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{1, 4, 0}, Point{1, 5, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{3, 3, 0}, Point{3, 4, 1}, Point{2, 4, 0}, Point{2, 5, 1}, Point{1, 5, 0}}, "color4"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 4, 0}, Point{0, 6, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 5, 0}, Point{1, 6, 1}, Point{0, 6, 0}}, "color3"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test19(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"11", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "12", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 4, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 0, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 0}}, "color2"},
		},
		map[string]Path{
			"9": Path{"9", []Point{Point{2, 2, 0}, Point{2, 9, 0}}, "color3"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{3, 3, 0}, Point{3, 7, 0}}, "color4"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{1, 4, 0}, Point{1, 11, 0}}, "color2"},
			"6":  Path{"6", []Point{Point{1, 4, 0}, Point{4, 4, 2}, Point{4, 6, 0}}, "color5"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 5, 0}, Point{0, 8, 0}}, "color1"},
			"6": Path{"6", []Point{Point{0, 5, 0}, Point{4, 5, 2}, Point{4, 6, 0}}, "color5"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{4, 6, 0}, Point{4, 8, 1}, Point{3, 8, 0}, Point{3, 10, 1}, Point{2, 10, 0}, Point{2, 11, 1}, Point{1, 11, 0}}, "color5"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{3, 7, 0}, Point{3, 8, 1}, Point{0, 8, 0}}, "color4"},
		},
		map[string]Path{
			"10": Path{"10", []Point{Point{0, 8, 0}, Point{0, 10, 0}}, "color1"},
		},
		map[string]Path{
			"10": Path{"10", []Point{Point{2, 9, 0}, Point{2, 10, 1}, Point{0, 10, 0}}, "color3"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{0, 10, 0}, Point{0, 12, 0}}, "color1"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{1, 11, 0}, Point{1, 12, 1}, Point{0, 12, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test20(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{1, 4, 0}, Point{1, 5, 1}, Point{0, 5, 0}}, "color3"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{3, 3, 0}, Point{3, 4, 1}, Point{0, 4, 0}}, "color4"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 4, 0}, Point{0, 5, 0}}, "color1"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test21(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{1, 3, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 5, 1}, Point{0, 5, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 3, 0}, Point{1, 6, 0}}, "color2"},
			"5": Path{"5", []Point{Point{1, 3, 0}, Point{2, 3, 2}, Point{2, 5, 1}, Point{0, 5, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 4, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{0, 5, 0}, Point{0, 7, 0}}, "color1"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{1, 6, 0}, Point{1, 7, 1}, Point{0, 7, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test22(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 1, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 0, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 0}}, "color2"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{2, 2, 0}, Point{2, 6, 1}, Point{0, 6, 0}}, "color3"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{3, 3, 0}, Point{3, 6, 1}, Point{2, 6, 0}, Point{2, 7, 1}, Point{1, 7, 0}}, "color4"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{1, 4, 0}, Point{1, 7, 0}}, "color2"},
			"6": Path{"6", []Point{Point{1, 4, 0}, Point{2, 4, 2}, Point{2, 6, 1}, Point{0, 6, 0}}, "color3"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 5, 0}, Point{0, 6, 0}}, "color1"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 6, 0}, Point{0, 8, 0}}, "color1"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 7, 0}, Point{1, 8, 1}, Point{0, 8, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test23(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 2, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{2, 2, 0}, Point{2, 4, 1}, Point{0, 4, 0}}, "color3"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{3, 3, 0}, Point{3, 4, 1}, Point{1, 4, 0}, Point{1, 7, 0}}, "color4"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 4, 0}, Point{0, 6, 0}}, "color1"},
			"5": Path{"5", []Point{Point{0, 4, 0}, Point{2, 4, 2}, Point{2, 5, 0}}, "color5"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{2, 5, 0}, Point{2, 6, 1}, Point{0, 6, 0}}, "color5"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 6, 0}, Point{0, 8, 0}}, "color1"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 7, 0}, Point{1, 8, 1}, Point{0, 8, 0}}, "color4"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test24(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"10", "8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"11", "8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"9"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "11", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 4, 1, 1, 0, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 3, 0}}, "color1"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 5, 0}}, "color2"},
		},
		map[string]Path{
			"9": Path{"9", []Point{Point{2, 2, 0}, Point{2, 9, 0}}, "color3"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{0, 3, 0}, Point{0, 7, 0}}, "color1"},
			"6": Path{"6", []Point{Point{0, 3, 0}, Point{3, 3, 2}, Point{3, 6, 1}, Point{1, 6, 0}}, "color4"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{4, 4, 0}, Point{4, 6, 1}, Point{1, 6, 0}}, "color5"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 5, 0}, Point{1, 6, 0}}, "color2"},
		},
		map[string]Path{
			"10": Path{"10", []Point{Point{1, 6, 0}, Point{1, 10, 0}}, "color2"},
			"8":  Path{"8", []Point{Point{1, 6, 0}, Point{3, 6, 2}, Point{3, 8, 0}}, "color6"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{0, 7, 0}, Point{0, 11, 0}}, "color1"},
			"8":  Path{"8", []Point{Point{0, 7, 0}, Point{3, 7, 2}, Point{3, 8, 0}}, "color6"},
		},
		map[string]Path{
			"9": Path{"9", []Point{Point{3, 8, 0}, Point{3, 9, 1}, Point{2, 9, 0}}, "color6"},
		},
		map[string]Path{
			"10": Path{"10", []Point{Point{2, 9, 0}, Point{2, 10, 1}, Point{1, 10, 0}}, "color3"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{1, 10, 0}, Point{1, 11, 1}, Point{0, 11, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test25(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"9", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"9", "7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"12", "9"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"11", "10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "11", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "12", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 1, 2, 0, 2, 3, 0, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"5": Path{"5", []Point{Point{0, 0, 0}, Point{0, 5, 0}}, "color1"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{1, 1, 0}, Point{1, 3, 0}}, "color2"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{2, 2, 0}, Point{2, 4, 0}}, "color3"},
		},
		map[string]Path{
			"9": Path{"9", []Point{Point{1, 3, 0}, Point{1, 9, 0}}, "color2"},
			"7": Path{"7", []Point{Point{1, 3, 0}, Point{3, 3, 2}, Point{3, 7, 0}}, "color4"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{2, 4, 0}, Point{2, 6, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 5, 0}, Point{0, 8, 0}}, "color1"},
			"7": Path{"7", []Point{Point{0, 5, 0}, Point{3, 5, 2}, Point{3, 7, 0}}, "color4"},
		},
		map[string]Path{
			"9": Path{"9", []Point{Point{2, 6, 0}, Point{2, 9, 1}, Point{1, 9, 0}}, "color3"},
			"7": Path{"7", []Point{Point{2, 6, 0}, Point{3, 6, 2}, Point{3, 7, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{3, 7, 0}, Point{3, 8, 1}, Point{0, 8, 0}}, "color4"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{0, 8, 0}, Point{0, 12, 0}}, "color1"},
			"9":  Path{"9", []Point{Point{0, 8, 0}, Point{1, 8, 2}, Point{1, 9, 0}}, "color2"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{1, 9, 0}, Point{1, 11, 0}}, "color2"},
			"10": Path{"10", []Point{Point{1, 9, 0}, Point{2, 9, 2}, Point{2, 10, 0}}, "color5"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{2, 10, 0}, Point{2, 11, 1}, Point{1, 11, 0}}, "color5"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{1, 11, 0}, Point{1, 12, 1}, Point{0, 12, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test26(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"8", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"7", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 1, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 0, 0}, Point{0, 3, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{1, 1, 0}, Point{1, 4, 0}}, "color2"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{2, 2, 0}, Point{2, 5, 1}, Point{1, 5, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{0, 3, 0}, Point{0, 8, 0}}, "color1"},
			"5": Path{"5", []Point{Point{0, 3, 0}, Point{2, 3, 2}, Point{2, 5, 1}, Point{1, 5, 0}}, "color3"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 4, 0}, Point{1, 5, 0}}, "color2"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{1, 5, 0}, Point{1, 7, 0}}, "color2"},
			"6": Path{"6", []Point{Point{1, 5, 0}, Point{2, 5, 2}, Point{2, 6, 0}}, "color4"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{2, 6, 0}, Point{2, 7, 1}, Point{1, 7, 0}}, "color4"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 7, 0}, Point{1, 8, 1}, Point{0, 8, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test27(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"11"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"15", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"8", "6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"15"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"9", "13"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"14", "10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "11", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "12", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "13", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "14", "parents": []string{"16"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "15", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "16", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "17", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 0, 1, 4, 2, 1, 1, 6, 3, 2, 5, 1, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 0, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"5": Path{"5", []Point{Point{1, 1, 0}, Point{1, 5, 0}}, "color2"},
		},
		map[string]Path{
			"7": Path{"7", []Point{Point{2, 2, 0}, Point{2, 7, 0}}, "color3"},
		},
		map[string]Path{
			"11": Path{"11", []Point{Point{3, 3, 0}, Point{3, 11, 0}}, "color4"},
		},
		map[string]Path{
			"15": Path{"15", []Point{Point{0, 4, 0}, Point{0, 15, 0}}, "color1"},
			"6":  Path{"6", []Point{Point{0, 4, 0}, Point{4, 4, 2}, Point{4, 6, 0}}, "color5"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{1, 5, 0}, Point{1, 8, 0}}, "color2"},
			"6": Path{"6", []Point{Point{1, 5, 0}, Point{4, 5, 2}, Point{4, 6, 0}}, "color5"},
		},
		map[string]Path{
			"15": Path{"15", []Point{Point{4, 6, 0}, Point{4, 14, 1}, Point{2, 14, 0}, Point{2, 15, 1}, Point{0, 15, 0}}, "color5"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{2, 7, 0}, Point{2, 12, 0}}, "color3"},
		},
		map[string]Path{
			"9":  Path{"9", []Point{Point{1, 8, 0}, Point{1, 9, 0}}, "color2"},
			"13": Path{"13", []Point{Point{1, 8, 0}, Point{5, 8, 2}, Point{5, 13, 0}}, "color6"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{1, 9, 0}, Point{1, 14, 0}}, "color2"},
			"10": Path{"10", []Point{Point{1, 9, 0}, Point{6, 9, 2}, Point{6, 10, 0}}, "color7"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{6, 10, 0}, Point{6, 14, 1}, Point{1, 14, 0}}, "color7"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{3, 11, 0}, Point{3, 14, 1}, Point{1, 14, 0}}, "color4"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{2, 12, 0}, Point{2, 14, 1}, Point{1, 14, 0}}, "color3"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{5, 13, 0}, Point{5, 14, 1}, Point{1, 14, 0}}, "color6"},
		},
		map[string]Path{
			"16": Path{"16", []Point{Point{1, 14, 0}, Point{1, 16, 0}}, "color2"},
		},
		map[string]Path{
			"17": Path{"17", []Point{Point{0, 15, 0}, Point{0, 17, 0}}, "color1"},
		},
		map[string]Path{
			"17": Path{"17", []Point{Point{1, 16, 0}, Point{1, 17, 1}, Point{0, 17, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test28(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"2", "1"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"2"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"3"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"4"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"6", "5"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"6"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 0, 0, 0, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"2": Path{"2", []Point{Point{0, 0, 0}, Point{0, 2, 0}}, "color1"},
			"1": Path{"1", []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}}, "color2"},
		},
		map[string]Path{
			"2": Path{"2", []Point{Point{1, 1, 0}, Point{1, 2, 1}, Point{0, 2, 0}}, "color2"},
		},
		map[string]Path{
			"3": Path{"3", []Point{Point{0, 2, 0}, Point{0, 3, 0}}, "color1"},
		},
		map[string]Path{
			"4": Path{"4", []Point{Point{0, 3, 0}, Point{0, 4, 0}}, "color1"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{0, 4, 0}, Point{0, 6, 0}}, "color1"},
			"5": Path{"5", []Point{Point{0, 4, 0}, Point{1, 4, 2}, Point{1, 5, 0}}, "color2"},
		},
		map[string]Path{
			"6": Path{"6", []Point{Point{1, 5, 0}, Point{1, 6, 1}, Point{0, 6, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func Test29(t *testing.T) {
	// Initial input
	inputNodes := make([]map[string]interface{}, 0)
	inputNodes = append(inputNodes, map[string]interface{}{"id": "0", "parents": []string{"7"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "1", "parents": []string{"15"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "2", "parents": []string{"17"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "3", "parents": []string{"8"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "4", "parents": []string{"18"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "5", "parents": []string{"12"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "6", "parents": []string{"20"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "7", "parents": []string{"9", "10"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "8", "parents": []string{"9", "11"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "9", "parents": []string{"13", "14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "10", "parents": []string{"21"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "11", "parents": []string{"13"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "12", "parents": []string{"14"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "13", "parents": []string{"16", "15"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "14", "parents": []string{"19"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "15", "parents": []string{"26"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "16", "parents": []string{"27"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "17", "parents": []string{"25"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "18", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "19", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "20", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "21", "parents": []string{"22"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "22", "parents": []string{"23"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "23", "parents": []string{"24"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "24", "parents": []string{"25"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "25", "parents": []string{"26"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "26", "parents": []string{"27"}})
	inputNodes = append(inputNodes, map[string]interface{}{"id": "27", "parents": []string{}})

	out, _ := BuildTree(inputNodes, customColors)

	// Expected output
	expectedColumns := []int{0, 1, 2, 3, 4, 5, 6, 0, 3, 0, 7, 3, 5, 0, 4, 1, 0, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1, 0}

	expectedPaths := []map[string]Path{
		map[string]Path{
			"7": Path{"7", []Point{Point{0, 0, 0}, Point{0, 7, 0}}, "color1"},
		},
		map[string]Path{
			"15": Path{"15", []Point{Point{1, 1, 0}, Point{1, 15, 0}}, "color2"},
		},
		map[string]Path{
			"17": Path{"17", []Point{Point{2, 2, 0}, Point{2, 17, 0}}, "color3"},
		},
		map[string]Path{
			"8": Path{"8", []Point{Point{3, 3, 0}, Point{3, 8, 0}}, "color4"},
		},
		map[string]Path{
			"18": Path{"18", []Point{Point{4, 4, 0}, Point{4, 13, 1}, Point{3, 13, 0}, Point{3, 18, 0}}, "color5"},
		},
		map[string]Path{
			"12": Path{"12", []Point{Point{5, 5, 0}, Point{5, 12, 0}}, "color6"},
		},
		map[string]Path{
			"20": Path{"20", []Point{Point{6, 6, 0}, Point{6, 13, 1}, Point{5, 13, 0}, Point{5, 20, 0}}, "color7"},
		},
		map[string]Path{
			"9":  Path{"9", []Point{Point{0, 7, 0}, Point{0, 9, 0}}, "color1"},
			"10": Path{"10", []Point{Point{0, 7, 0}, Point{7, 7, 2}, Point{7, 10, 0}}, "color8"},
		},
		map[string]Path{
			"9":  Path{"9", []Point{Point{3, 8, 0}, Point{0, 8, 3}, Point{0, 9, 0}}, "color1"},
			"11": Path{"11", []Point{Point{3, 8, 0}, Point{3, 11, 0}}, "color4"},
		},
		map[string]Path{
			"13": Path{"13", []Point{Point{0, 9, 0}, Point{0, 13, 0}}, "color1"},
			"14": Path{"14", []Point{Point{0, 9, 0}, Point{8, 9, 2}, Point{8, 13, 1}, Point{7, 13, 0}, Point{7, 14, 1}, Point{4, 14, 0}}, "color9"},
		},
		map[string]Path{
			"21": Path{"21", []Point{Point{7, 10, 0}, Point{7, 13, 1}, Point{6, 13, 0}, Point{6, 21, 0}}, "color8"},
		},
		map[string]Path{
			"13": Path{"13", []Point{Point{3, 11, 0}, Point{3, 13, 1}, Point{0, 13, 0}}, "color4"},
		},
		map[string]Path{
			"14": Path{"14", []Point{Point{5, 12, 0}, Point{5, 13, 1}, Point{4, 13, 0}, Point{4, 14, 0}}, "color6"},
		},
		map[string]Path{
			"16": Path{"16", []Point{Point{0, 13, 0}, Point{0, 16, 0}}, "color1"},
			"15": Path{"15", []Point{Point{0, 13, 0}, Point{1, 13, 2}, Point{1, 15, 0}}, "color2"},
		},
		map[string]Path{
			"19": Path{"19", []Point{Point{4, 14, 0}, Point{4, 19, 0}}, "color6"},
		},
		map[string]Path{
			"26": Path{"26", []Point{Point{1, 15, 0}, Point{1, 26, 0}}, "color2"},
		},
		map[string]Path{
			"27": Path{"27", []Point{Point{0, 16, 0}, Point{0, 27, 0}}, "color1"},
		},
		map[string]Path{
			"25": Path{"25", []Point{Point{2, 17, 0}, Point{2, 25, 0}}, "color3"},
		},
		map[string]Path{
			"24": Path{"24", []Point{Point{3, 18, 0}, Point{3, 24, 0}}, "color5"},
		},
		map[string]Path{
			"23": Path{"23", []Point{Point{4, 19, 0}, Point{4, 23, 0}}, "color6"},
		},
		map[string]Path{
			"22": Path{"22", []Point{Point{5, 20, 0}, Point{5, 22, 0}}, "color7"},
		},
		map[string]Path{
			"22": Path{"22", []Point{Point{6, 21, 0}, Point{6, 22, 1}, Point{5, 22, 0}}, "color8"},
		},
		map[string]Path{
			"23": Path{"23", []Point{Point{5, 22, 0}, Point{5, 23, 1}, Point{4, 23, 0}}, "color7"},
		},
		map[string]Path{
			"24": Path{"24", []Point{Point{4, 23, 0}, Point{4, 24, 1}, Point{3, 24, 0}}, "color6"},
		},
		map[string]Path{
			"25": Path{"25", []Point{Point{3, 24, 0}, Point{3, 25, 1}, Point{2, 25, 0}}, "color5"},
		},
		map[string]Path{
			"26": Path{"26", []Point{Point{2, 25, 0}, Point{2, 26, 1}, Point{1, 26, 0}}, "color3"},
		},
		map[string]Path{
			"27": Path{"27", []Point{Point{1, 26, 0}, Point{1, 27, 1}, Point{0, 27, 0}}, "color2"},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
	validateColors(t, expectedPaths, out)
}

func TestPathHeight1(t *testing.T) {
	out := OutputNode{ParentsPaths: map[string]Path{"1": Path{Path: []Point{
		Point{X: 0, Y: 2, Type: 0},
		Point{X: 3, Y: 2, Type: 2},
		Point{X: 3, Y: 9, Type: 1},
		Point{X: 2, Y: 9, Type: 0},
		Point{X: 2, Y: 11, Type: 1},
		Point{X: 1, Y: 11, Type: 0},
	}}}}
	if out.GetPathHeightAtIdx("1", 1) != -1 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 2) != 3 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 3) != 3 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 9) != 2 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 10) != 2 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 11) != 1 {
		t.Fail()
	}
	if out.GetPathHeightAtIdx("1", 1000) != -1 {
		t.Fail()
	}
}
