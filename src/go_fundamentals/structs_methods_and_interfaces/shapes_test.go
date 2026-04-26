package main

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{Width: 4.0, Height: 7.0}
	got := rectangle.Perimeter()
	want := 22.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Height: 9, Width: 12}, hasArea: 108.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Rectangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g hasArea %g", tt.shape, got, tt.hasArea)
			}
		})
	}
}
