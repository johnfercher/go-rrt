package shared

import "fmt"

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

func (v *Vector3D) Println() {
	fmt.Printf("(%f, %f, %f)\n", v.X, v.Y, v.Z)
}

func (v *Vector3D) GetString() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}
