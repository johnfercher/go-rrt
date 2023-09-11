package math

import (
	"fmt"
	"math"
)

type Point[T any] struct {
	X    int
	Y    int
	Data T
}

func NewPoint(x int, y int) *Point[any] {
	return &Point[any]{
		X: x,
		Y: y,
	}
}

func (p *Point[T]) DistanceTo(other *Point[T]) float64 {
	pX, pY := float64(p.X), float64(p.Y)
	oX, oY := float64(other.X), float64(other.Y)

	return math.Sqrt(((pX - oX) * (pX - oX)) + ((pY - oY) * (pY - oY)))
}

func (p *Point[T]) Atan2To(other *Point[T]) float64 {
	pX, pY := float64(p.X), float64(p.Y)
	oX, oY := float64(other.X), float64(other.Y)

	return math.Atan2(oX-pX, oY-pY)
}

func (p *Point[T]) GetString() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}
