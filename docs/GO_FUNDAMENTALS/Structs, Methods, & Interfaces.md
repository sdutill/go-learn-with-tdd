# Structs, methods, & interfaces

Suppose that we need some geometry code to calculate the perimeter of a rectanglt given a height and width. We can write a `Perimeter(width flat64, height float64)` function, where `float64` is for floating-point numbers like `123.45`.

The TDD cycle should be pretty familiar to you by now.

## Write the test first

```go
package shapes

import "testing"

func TestPPerimeter(t *testing.T) {
	got := Perimeter(4.0, 7.0)
	want := 28.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

## Try to run the test

```sh
ethods_and_interfaces.test]
.\shapes_test.go:6:9: undefined: Perimeter
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/structs_methods_and_interfaces [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
package shapes

func Perimeter(length, width float64) float64 {
	return 2.2
	// return 2 * (length + width)
}
```

## Write enough code to make it pass

```go
package shapes

func Perimeter(length, width float64) float64 {
	return 2 * (length + width)
}
```

So far, so easy. Now let's create a function called `Area(width, height float64)` which returns the area of a rectangle. Try to do it yourself, following the TDD cycle.

You should end up with tests like this

```go
package shapes

import "testing"

func TestPerimeter(t *testing.T) {
	got := Perimeter(4.0, 7.0)
	want := 22.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	got := Area(9.0, 12.0)
	want := 108.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

And code like this

```go
package shapes

func Perimeter(length, width float64) float64 {
	return 2 * (length + width)
}

func Area(length, width float64) float64 {
	return length * width
}package shapes

func Perimeter(length, width float64) float64 {
	return 2 * (length + width)
}

func Area(length, width float64) float64 {
	return length * width
}
```

## Refactor

Our code does the job, but it doesn't contain anything explicit about rectangles. An unwary developer might try to supply the width and height of a triangle to these functions without realising they will return the wrong answer.

We could just give the functions more spcific names like `RectangleArea`. A neater solution is to define our own _type_ called `Rectangle` which ncapsulates this concept for us.

We can create a simple type using a **struct**. A struct is just a named collcetion of fields where you can store data.

Declare a struct in your `shapes.go` file like this

```go
type Rectangle struct {
	Width  float64
	Height float64
}
```

Now let's refactor the tests to use `Rectangle instead of plain `float64`s.

```go
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
	rectangle := Rectangle{Width: 9.0, Height: 12.0}
	got := Area(rectangle)
	want := 108.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

Rememberr to run your tests before attempting to fix. The tests should show a helpful error like

```sh
# github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/structs_methods_and_interfaces [github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/structs_methods_and_interfaces.test]
.\shapes_test.go:7:19: not enough arguments in call to Perimeter
        have (Rectangle)
        want (float64, float64)
.\shapes_test.go:17:14: not enough arguments in call to Area
        have (Rectangle)
        want (float64, float64)
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/structs_methods_and_interfaces [build failed]
```

You can access the fields of a struct with the syntax of `myStruct.field`.

Change the two functions to fix the test.

```go
package shapes

type Rectangle struct {
	Width  float64
	Height float64
}

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Height + r.Width)
}

func Area(r Rectangle) float64 {
	return r.Height * r.Width
}
```

I hope you'll agree that passing a `Rectangle` to a function conveys our intent more clearly, but there are more benefits of using structs that we will cover later.

Our next requirement is to write an `Area` function for circles.

## Write the test first

```go
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
```

As you can see, the `f` has been replaced by `g`, with good reason. Use of `g` will print a more precise decimal number in th eerror message. For example, using a radius of 1.5 in a circle area calculation, `f` would show `7.068583` whereas g would show `7.0685834705770345`.

##
