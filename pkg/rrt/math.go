package rrt

import "math"

func Distance[T any](a *Point[T], b *Point[T]) float64 {
	return math.Sqrt(((a.X - b.X) * (a.X - b.X)) + ((a.Y - b.Y) * (a.Y - b.Y)))
}

func Radian[T any](t1 *Point[T], t2 *Point[T]) float64 {
	return math.Atan2(t2.X-t1.X, t2.Y-t1.Y)
}
