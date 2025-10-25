package main

import (
	"fmt"
	"math"
)

type Point struct {
	x float64
	y float64
}

func NewPoint(_x float64, _y float64) Point {
	return Point{x: _x, y: _y}
}
func (p Point) Distance(other Point) float64 {
	return math.Sqrt(math.Pow(p.x-other.x, 2) + math.Pow(p.y-other.y, 2))
}

func main() {
	p1 := NewPoint(1.553453, 2)
	p2 := NewPoint(1, 3.9032)

	x := p1.Distance(p2)
	fmt.Println(x)

}
