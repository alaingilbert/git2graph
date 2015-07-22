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

func validatePaths(t *testing.T, expectedPaths []map[string][]Point, data []*OutputNode) {
	for nodeIdx, node := range data {
		for _, parentId := range node.Parents {
			for pathIdx, pathNode := range node.ParentsPaths[parentId] {
				if pathNode != expectedPaths[nodeIdx][parentId][pathIdx] {
					t.Fail()
					t.Logf("Id: %s, Expected path: %d, Actual path: %d", node.Id, expectedPaths[nodeIdx][parentId][pathIdx], pathNode)
				}
			}
		}
	}
}

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

	out, _ := buildTree(inputNodes)

	// Expected output
	expectedColumns := []int{0, 0, 0}

	expectedPaths := []map[string][]Point{
		map[string][]Point{
			"2": []Point{Point{0, 0, 0}, Point{0, 1, 0}},
		},
		map[string][]Point{
			"3": []Point{Point{0, 1, 0}, Point{0, 2, 0}},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
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

	out, _ := buildTree(inputNodes)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string][]Point{
		map[string][]Point{
			"3": []Point{Point{0, 0, 0}, Point{0, 2, 0}},
		},
		map[string][]Point{
			"3": []Point{Point{1, 1, 0}, Point{1, 2, 1}, Point{0, 2, 0}},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
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

	out, _ := buildTree(inputNodes)

	// Expected output
	expectedColumns := []int{0, 1, 0}

	expectedPaths := []map[string][]Point{
		map[string][]Point{
			"3": []Point{Point{0, 0, 0}, Point{0, 2, 0}},
			"2": []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}},
		},
		map[string][]Point{
			"3": []Point{Point{1, 1, 0}, Point{1, 2, 1}, Point{0, 2, 0}},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
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

	out, _ := buildTree(inputNodes)

	// Expected output
	expectedColumns := []int{0, 1, 0, 2, 0}

	expectedPaths := []map[string][]Point{
		map[string][]Point{
			"3": []Point{Point{0, 0, 0}, Point{0, 2, 0}},
			"2": []Point{Point{0, 0, 0}, Point{1, 0, 2}, Point{1, 1, 0}},
		},
		map[string][]Point{
			"5": []Point{Point{1, 1, 0}, Point{1, 4, 1}, Point{0, 4, 0}},
		},
		map[string][]Point{
			"5": []Point{Point{0, 2, 0}, Point{0, 4, 0}},
			"4": []Point{Point{0, 2, 0}, Point{2, 2, 2}, Point{2, 3, 0}},
		},
		map[string][]Point{
			"5": []Point{Point{2, 3, 0}, Point{2, 4, 1}, Point{0, 4, 0}},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
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

	out, _ := buildTree(inputNodes)

	// Expected output
	expectedColumns := []int{0, 1, 2, 0, 1, 0}

	expectedPaths := []map[string][]Point{
		map[string][]Point{
			"4": []Point{Point{0, 0, 0}, Point{0, 3, 0}},
		},
		map[string][]Point{
			"4": []Point{Point{1, 1, 0}, Point{1, 3, 1}, Point{0, 3, 0}},
		},
		map[string][]Point{
			"4": []Point{Point{2, 2, 0}, Point{2, 3, 1}, Point{0, 3, 0}},
		},
		map[string][]Point{
			"6": []Point{Point{0, 3, 0}, Point{0, 5, 0}},
		},
		map[string][]Point{
			"6": []Point{Point{1, 4, 0}, Point{1, 5, 1}, Point{0, 5, 0}},
		},
	}

	// Validation
	validateColumns(t, expectedColumns, out)
	validatePaths(t, expectedPaths, out)
}
