package main

import (
	"errors"
	"fmt"
	"github.com/johnfercher/rrt/pkg/rrt"
	"github.com/johnfercher/rrt/pkg/shared"
	"math/rand"
)

func main() {
	stepDistance := 0.1
	rrt := rrt.New(stepDistance)

	rrt.AddCollisionCondition(func(vector3D *shared.Vector3D) bool {
		return vector3D.Z == 1
	})

	rrt.AddStopCondition(func(testPoint *shared.Vector3D, finish *shared.Vector3D) bool {
		return shared.Distance(testPoint, finish) <= 10
	})

	space := generateClearSpace(128, 128)
	space = addObstacles(5, 3, space)

	start, finish, err := getStartAndFinishPoints(space)
	if err != nil {
		panic(err)
	}

	start.Println()
	finish.Println()
	nodes := rrt.FindPath(start, finish, space)
	for _, node := range nodes {
		node.Print("")
	}
}

func generateClearSpace(x, y int) [][]*shared.Vector3D {
	var space [][]*shared.Vector3D
	for i := 0; i < x; i++ {
		var line []*shared.Vector3D
		for j := 0; j < y; j++ {
			line = append(line, &shared.Vector3D{X: float64(i), Y: float64(j), Z: 0})
		}
		space = append(space, line)
	}

	return space
}

func addObstacles(qtd, size int, space [][]*shared.Vector3D) [][]*shared.Vector3D {
	for i := 0; i < qtd; i++ {
		x := rand.Int()%len(space) - 1
		y := rand.Int()%len(space[0]) - 1
		//fmt.Printf("%d, %d\n", x, y)
		space = addObstacle(x, y, size, space)
	}

	return space
}

func addObstacle(x, y, size int, space [][]*shared.Vector3D) [][]*shared.Vector3D {
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
			space[i][j].Z = 1
		}
	}

	return space
}

func getStartAndFinishPoints(space [][]*shared.Vector3D) (*shared.Vector3D, *shared.Vector3D, error) {
	tries := 10
	start := &shared.Vector3D{}
	finish := &shared.Vector3D{}

	found := false
	for i := 0; i < tries && !found; i++ {
		x := rand.Int() % len(space)
		y := rand.Int() % len(space[0])

		//fmt.Printf("%d - %d\n", x, y)

		if space[x][y].Z != 1 {
			space[x][y].Z = 2
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

		if space[x][y].Z != 1 {
			space[x][y].Z = 3
			finish = space[x][y]
			found = true
		}
	}
	if !found {
		return nil, nil, errors.New("not found finish point")
	}

	return start, finish, nil
}

func print(space [][]*shared.Vector3D) {
	for i := 0; i < len(space); i++ {
		for j := 0; j < len(space[i]); j++ {
			fmt.Printf("%0.f ", space[i][j].Z)
		}
		fmt.Println()
	}
}
