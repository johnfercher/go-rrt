package rrt_test

import (
	"fmt"
	"github.com/johnfercher/go-rrt/pkg/rrt"
	"testing"
)

func TestRapidlyExploringRandomTrees_GetFixedPoint(t *testing.T) {
	// Arrange
	new := &rrt.Point[string]{X: 105, Y: 109}
	min := &rrt.Point[string]{X: 5, Y: 71}

	world := generateClearSpace[string](128, 128)
	sut := rrt.New[string](0.1, 10000, 15)

	// Act
	fixed := sut.GetFixedPoint(min, new, world)

	// Assert
	fmt.Println(fixed.GetString())
}

func generateClearSpace[T string](x, y int) [][]T {
	var space [][]T
	for i := 0; i < x; i++ {
		var line []T
		for j := 0; j < y; j++ {
			line = append(line, "empty")
		}
		space = append(space, line)
	}

	return space
}
