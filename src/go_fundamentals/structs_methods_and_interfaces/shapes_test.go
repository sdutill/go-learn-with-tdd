package shapes

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{Width: 4.0, Height: 7.0}
	got := Perimeter(rectangle)
	want := 22.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	t.Run("calculate the area of a rectangle", func(t *testing.T) {
		rectangle := Rectangle{Width: 9.0, Height: 12.0}
		got := Area(rectangle)
		want := 108.0

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})
	t.Run("calculate the area of a circle", func(t *testing.T) {
		circle := Circle{Radius: 10.0}
		got := Area(circle)
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})
}
