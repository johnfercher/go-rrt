package main

import (
	"errors"
	"fmt"
	"github.com/johnfercher/go-rrt/pkg/rrt"
	"math/rand"
)

func main() {
	stepDistance := 0.1
	r := rrt.New[string](stepDistance)

	r.AddCollisionCondition(func(point *rrt.Point[string]) bool {
		return point.Data == "obstacle"
	})

	r.AddStopCondition(func(testPoint *rrt.Point[string], finish *rrt.Point[string]) bool {
		return testPoint.DistanceTo(finish) <= 5
	})

	space := generateClearSpace(128, 128)
	space = addObstacles(5, 3, space)

	start, finish, err := getStartAndFinishPoints(space)
	if err != nil {
		panic(err)
	}

	start.Println()
	finish.Println()
	points := r.FindPath(start, finish, space)
	for _, point := range points {
		point.Println()
	}
}

func generateClearSpace[T string](x, y int) [][]*rrt.Point[T] {
	var space [][]*rrt.Point[T]
	for i := 0; i < x; i++ {
		var line []*rrt.Point[T]
		for j := 0; j < y; j++ {
			point := &rrt.Point[T]{
				X:    float64(i),
				Y:    float64(j),
				Data: "empty",
			}
			line = append(line, point)
		}
		space = append(space, line)
	}

	return space
}

func addObstacles[T string](qtd, size int, space [][]*rrt.Point[T]) [][]*rrt.Point[T] {
	for i := 0; i < qtd; i++ {
		x := rand.Int()%len(space) - 1
		y := rand.Int()%len(space[0]) - 1
		//fmt.Printf("%d, %d\n", x, y)
		space = addObstacle(x, y, size, space)
	}

	return space
}

func addObstacle[T string](x, y, size int, space [][]*rrt.Point[T]) [][]*rrt.Point[T] {
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

	for i := minXOffset; i <= maxXOffset; i++ {
		for j := minYOffset; j <= maxYOffset; j++ {
			//fmt.Printf("%d - %d\n", i, j)
			space[i][j].Data = "obstacle"
		}
	}

	return space
}

func getStartAndFinishPoints[T string](space [][]*rrt.Point[T]) (*rrt.Point[T], *rrt.Point[T], error) {
	tries := 10
	start := &rrt.Point[T]{}
	finish := &rrt.Point[T]{}

	found := false
	for i := 0; i < tries && !found; i++ {
		x := rand.Int() % len(space)
		y := rand.Int() % len(space[0])

		//fmt.Printf("%d - %d\n", x, y)

		if space[x][y].Data != "obstacle" {
			space[x][y].Data = "start"
			start = space[x][y]
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

		if space[x][y].Data != "obstacle" && space[x][y].Data != "start" {
			space[x][y].Data = "finish"
			finish = space[x][y]
			found = true
		}
	}
	if !found {
		return nil, nil, errors.New("not found finish point")
	}

	return start, finish, nil
}

func print[T string](space [][]*rrt.Point[T]) {
	for i := 0; i < len(space); i++ {
		for j := 0; j < len(space[i]); j++ {
			fmt.Printf("%s ", space[i][j].Data)
		}
		fmt.Println()
	}
}
