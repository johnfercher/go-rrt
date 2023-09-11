package main

import (
	"github.com/johnfercher/go-rrt/pkg/math"
	"github.com/johnfercher/go-rrt/pkg/rrt"
	"math/rand"
)

func main() {
	stepDistance := 5
	maxTries := 1000000
	focusOnFinishEveryTry := 3
	r := rrt.New[string](stepDistance, maxTries, focusOnFinishEveryTry)

	r.AddCollisionCondition(func(point string) bool {
		return point == "obstacle"
	})

	r.AddStopCondition(func(testPoint *math.Point[string], finish *math.Point[string]) bool {
		return testPoint.DistanceTo(finish) <= 2
	})

	space := generateClearSpace(128, 128)
	space = addObstacles(25, 15, space)

	start := math.NewPoint(3, 3)
	finish := math.NewPoint(125, 125)

	_ = r.FindPathAndSavePdf(start, finish, space, "tree.pdf")
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

func addObstacles[T string](qtd, size int, space [][]T) [][]T {
	border := 3 + size
	for i := 0; i < qtd; i++ {
		x := rand.Int()%len(space) - 1
		y := rand.Int()%len(space[0]) - 1

		if x < border {
			x = border
		} else if x > len(space)-border {
			x = len(space) - border
		}

		if y < border {
			y = border
		} else if y > len(space[0])-border {
			y = len(space[0]) - border
		}
		//fmt.Printf("%d, %d\n", x, y)
		space = addObstacle(x, y, size, space)
	}

	return space
}

func addObstacle[T string](x, y, size int, space [][]T) [][]T {
	offset := (size - 1) / 2
	//fmt.Println(offset)

	minXOffset := x - offset
	if minXOffset < 0 {
		minXOffset = 0
	}

	maxXOffset := x + offset
	if maxXOffset > len(space) {
		maxXOffset = len(space)
	}

	minYOffset := y - offset
	if minYOffset < 0 {
		minYOffset = 0
	}

	maxYOffset := y + offset
	if maxYOffset > len(space[0]) {
		maxYOffset = len(space[0])
	}

	//fmt.Printf("XOffset: %d, %d\n", minXOffset, maxXOffset)
	//fmt.Printf("YOffset: %d, %d\n", minYOffset, maxYOffset)

	for i := minXOffset; i < maxXOffset; i++ {
		for j := minYOffset; j < maxYOffset; j++ {
			//fmt.Printf("%d - %d\n", i, j)
			space[i][j] = "obstacle"
		}
	}

	return space
}
