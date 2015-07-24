package main

import (
	"testing"
)

func validateColumns(t *testing.T, expectedColumns []int, data []*OutputNode) {
	for idx, row := range data {
		if row.Column != expectedColumns[idx] {
			t.Fail()
			t.Logf("Id: %s, Expected column: %d, Actual column: %d", row.Id, expectedColumns[idx], row.Column)
		}
	}
}

func validatePaths(t *testing.T, expectedPaths []map[string]Path, data []*OutputNode) {
	for nodeIdx, node := range data {
		for _, parentId := range node.Parents {
			if len(expectedPaths) - 1 < nodeIdx {
				continue
			}
			if len(node.ParentsPaths[parentId].Path) != len(expectedPaths[nodeIdx][parentId].Path) {
				t.Fail()
				t.Logf("Id: %s, Expected nb paths: %d, Actual nb paths: %d", node.Id, len(expectedPaths[nodeIdx][parentId].Path), len(node.ParentsPaths[parentId].Path))
				t.Logf("Id: %s, Expected: %v, Actual: %v", node.Id, expectedPaths[nodeIdx][parentId], node.ParentsPaths[parentId])
				return
			}
			for pathIdx, pathNode := range node.ParentsPaths[parentId].Path {
				if pathNode != expectedPaths[nodeIdx][parentId].Path[pathIdx] {
					t.Fail()
					t.Logf("Id: %s, Expected path: %d, Actual path: %d", node.Id, expectedPaths[nodeIdx][parentId].Path[pathIdx], pathNode)
					t.Logf("Id: %s, Expected: %v, Actual: %v", node.Id, expectedPaths[nodeIdx][parentId].Path, node.ParentsPaths[parentId].Path)
				}
			}
		}
	}
}

func validateColors(t *testing.T, expectedPaths []map[string]Path, data []*OutputNode) {
	for nodeIdx, node := range data {
		for _, parentId := range node.Parents {
			if len(expectedPaths) - 1 < nodeIdx {
				continue
			}
			if expectedPaths[nodeIdx][parentId].Color != node.ParentsPaths[parentId].Color {
				t.Logf("Id: %s, Expected: %v, Actual: %v", node.Id, expectedPaths[nodeIdx][parentId].Color, node.ParentsPaths[parentId].Color)
				t.Fail()
			}
		}
	}
}

var customColors []string = []string{"color1", "color2", "color3", "color4", "color5", "color6"}

// 1
// |
// 2
// |
// 3
func Test1(t *testing.T) {
	// Initial input
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"3"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"3"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"3"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"5", "4"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"4"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"4"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"4"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"3", "4"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"4", "5"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"4"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"4", "5"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"4", "7"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"5", "6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"8"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{"8"}})
	inputNodes = append(inputNodes, InputNode{"7", []string{"8"}})
	inputNodes = append(inputNodes, InputNode{"8", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"4", "2"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"5", "3"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"8", "6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"7", "6"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{"7"}})
	inputNodes = append(inputNodes, InputNode{"7", []string{"8"}})
	inputNodes = append(inputNodes, InputNode{"8", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
	inputNodes := make([]InputNode, 0)
	inputNodes = append(inputNodes, InputNode{"1", []string{"3"}})
	inputNodes = append(inputNodes, InputNode{"2", []string{"3", "4"}})
	inputNodes = append(inputNodes, InputNode{"3", []string{"5"}})
	inputNodes = append(inputNodes, InputNode{"4", []string{"5", "6"}})
	inputNodes = append(inputNodes, InputNode{"5", []string{"6"}})
	inputNodes = append(inputNodes, InputNode{"6", []string{}})

	out, _ := buildTree(inputNodes, customColors)

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
