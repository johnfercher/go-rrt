package shared

import "math"

func Distance(a *Vector3D, b *Vector3D) float64 {
	return math.Sqrt(((a.X - b.X) * (a.X - b.X)) + ((a.Y - b.Y) * (a.Y - b.Y)))
}

func Angle(t1 *Vector3D, t2 *Vector3D) float64 {
	return math.Atan2(t2.X-t1.X, t2.Y-t1.Y) * (180.0 / math.Pi)
}

func Radian(t1 *Vector3D, t2 *Vector3D) float64 {
	return math.Atan2(t2.X-t1.X, t2.Y-t1.Y)
}
