package main

import (
	"errors"
	"github.com/johnfercher/go-rrt/pkg/rrt"
	"github.com/johnfercher/go-rrt/pkg/rrt/math"
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
	space = addObstacles(15, 20, space)

	start, finish, err := getStartAndFinishCoordinate(space)
	if err != nil {
		panic(err)
	}

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
	for i := 0; i < qtd; i++ {
		x := rand.Int()%len(space) - 1
		y := rand.Int()%len(space[0]) - 1
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

func getStartAndFinishCoordinate[T string](space [][]T) (*math.Coordinate, *math.Coordinate, error) {
	tries := 10
	var start *math.Coordinate
	var finish *math.Coordinate

	found := false
	for i := 0; i < tries && !found; i++ {
		x := rand.Int() % len(space)
		y := rand.Int() % len(space[0])

		//fmt.Printf("%d - %d\n", x, y)

		if space[x][y] != "obstacle" {
			space[x][y] = "start"
			start = &math.Coordinate{X: float64(x), Y: float64(y)}
			found = true
		}
	}
	if !found {
		return nil, nil, errors.New("not found start point")
	}

	found = false
	for i := 0; i < tries && !found; i++ {
		x := rand.Int() % len(space)
		y := rand.Int() % len(space[0])

		//fmt.Printf("%d - %d\n", x, y)

		if space[x][y] != "obstacle" && space[x][y] != "start" {
			space[x][y] = "finish"
			finish = &math.Coordinate{X: float64(x), Y: float64(y)}
			found = true
		}
	}
	if !found {
		return nil, nil, errors.New("not found finish point")
	}

	return start, finish, nil
}
