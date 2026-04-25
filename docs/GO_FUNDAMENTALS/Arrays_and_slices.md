# Arrays and slices

Arrays allow you to store multiple elements of the same type in a variable in a particular order.

When you ahve arrays, it is very common to have to iterate over them. So let's use our new-found-knowledge of `for` to make a `Sum`. `Sum will take an array of numbers and return the total.

Let's use our TDD skills

## Write the test first

Create a new folder to work in. Create a new file called `sum_test.go` and insert the following:

```go
package arrays_and_slices

import "testing"

func TestSum(t *testing.T) {
	numbers := [5]int{1, 2, 3, 4, 5}

	got := Sum(numbers)
	want := 15

	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
```

Arrays have _fixed capacity_ which you define when you declare the variable. We can initialize an array in two ways:

- [N]type{value1, value2, ..., valueN} e.g.`numbers := [5]int{1, 2, 3, 4, 5}
- [...]type{value1, value2, ..., valueN} e.g.`numbers := [...]int{1, 2, 3, 4, 5}

It is sometimes useful to print the inputs to teh function in the error message. Here, we are using the `%v` placeholder to print the "default" format, which works well for arrays.

## Try to run the test

If you had initialized go mod with `go mod init main` you will be presented with an error `_testmain.go:13:2: cannot import "main"`. This is because according to common practice, package main will only contain integration of other packages and not unit-testable code and hence Go will not allow you to import a package with name `main`.

## Write the minimal amount of code for the test to run and check the failing test output

In `sum.go`

```go
package arrays_and_slices

func Sum(numbers [5]int ) int {
	return 0

}
```

Your test should now fail with a _clear error message_

```sh
=== RUN   TestSum
    sum_test.go:12: got 0 want 15 given, [1 2 3 4 5]
--- FAIL: TestSum (0.00s)
```

## Write enough code to make it pass

```go
func Sum(numbers [5]int) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += numbers[i]
	}
	return sum
}
```

To get the value out of an array at a particular index, just use `array[index]` syntax. In this case, we are using `for` to iterate 5 times to work through the array and add each item onto `sum`.

## Refactor

Let's introduce `range` to help clean up our code

```go
func Sum(numbers [5]int ) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

`range` lets you iterate over an array. On each iteration, `range` returns two values - the index and the value. We are choosing to ignore the index by using `_` blank identifier.

### Arrays and their type

An interesting property of arrays is that the size is encoded in its type. If you try to pass an `[4]int` into a function that expects `[5]int`, it won't compile. They are different types so it's just the same as trying to pass a `string` into a function that wants an `int`.

You may be thinking it's quite cumbersome that arrays have a fixed length, and most of the time you probably won't be using them!

Go has _slices_ which do not encode the size of the collection and instead can have any size.

The next requirement will be to sum collections of varying sizes.

## Write the test first

We will now use the slice type which allows us to have collections of any size. The syntax is very simliar to arrays, you just omit the size when declaring them

`mySlice := []int{1, 2, 3}` rather than `myArray := [3]int{1, 2, 3}`

```go
package arrays_and_slices

import "testing"

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T)	{
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection any slice", func(t *testing.T)	{
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}
```

## Try and run the test

This does not compile

```bash
.\sum_test.go:20:14: cannot use numbers (variable of type []int) as [5]int value in argument to Sum
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/arrays_and_slices [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

The problem here is we can either

- Break the existing API by changing the argument to `Sum` to be a slice rather than an array. When we do this, we will potentially ruin someone's day because our _other_ test will no longer compile!
- Create a new function

In our case, no one else is using our function, so rather than having two functions to maintain, let's have just one.

```go
package arrays_and_slices

func Sum(numbers []int ) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

If you try to run the tests they will still not compile, you will have to change the first test to pass in a slice rather than an array.

## Write enough code to make it pass

It turns out that fixing the compiler problems were all we need to do here and the tests pass!

## Refactor

We already refactored `Sum` - all we did was replace arrays with slices, so no extra changes are required. Remember that we must not neglect our test code in the refactoring stage - we can further improve our `Sum` tests.

```go
package arrays_and_slices

import "testing"

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T)	{
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection any slice", func(t *testing.T)	{
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}
```

It is important to question the value of your tests. It should not be a goal to have as many tests as possible, but rather to have as much _confidence_ as possible in your code base. Having too many tests can turn in to a real problem and it just adds more overhead in maintenance. **Every test has a cost.**

In our case, you can see that having two tests for this functin is redundant. If it works for a slice of one size, it's very likely that it'll work for a slice of any size (within reason).

Go's built-in testing toolkit features a coverage tool. Whilst striving for 100% coverage should not be your end goal, the coverage tool can help identify areas of your code not covered by tests. If you have been strict with TDD, it's quite likely you'll have close to 10% coverage anyways.

Try running

`go test -cover`

You should see

```bash
PASS
coverage: 100.0% of statements
ok      github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/arrays_and_slices
```

Now delete one of the tests and check the coverage again.

Now that we are happy we have a well-tested function you shoudl commit your great work before taking on the next challenge.

We need a new function called `SumAll` which will take a varying number of slices, returning a new slice containing the totals for each slice passed in.

For example

`SumAll([]int{1,2}, []int{0,9})` would return `[]int{3,9}`
or
`SumAll([]int{1,1,1}, []int{0,9})` would return `[]int{3}`

## Write the test first

```go
func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## Try and run the test

`.\sum_test.go:30:9: undefined: SumAll`

## Write the minimal amount of code for the test to run and check the failing test output

We need to define `SumAll` according to what our test wants.

Go can let you write variadic functions that can take a variable number of arguments.

```go
func SumAll(numbersToSum  ...[]int ) []int {
	return nil
}
```

This is valid, but our tests still won't compile!

`.\sum_test.go:33:5: invalid operation: got != want (slice can only be compared to nil)`

Go does not let you use equality operators with slices. You _could_ write a function to iterate over each `got` and `want` slice and check their values, but what if we had a more convenient way to do this?

From Go 1.21, slices standard package is available, which has slices.Equal function to do a simple shallow compare on slices, where you don't need to worry about the types like the above case. Note that this function expects the elements to be comparable. So, it can't be applied to slices with non-comparable elements like 2D slices.

Let's go ahead and put this into practice!

```go
func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

You should have test output like the following:

```bash
--- FAIL: TestSumAll (0.00s)
    sum_test.go:37: got [] want [3 9]
```

## Write enough code to make it pass

What we need to do is iterate over the varargs, calculate the sum using our existing `Sum` function, and then add it to the slice we will return.

```go
func SumAll(numbersToSum  ...[]int ) []int {
	lengthOfNumbersToSum := len(numbersToSum)
	sums := make([]int, lengthOfNumbersToSum)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return sums
}
```

Lots of new things to learn!

There's a new way to create a slice. `make` allows you to create a slice with a starting capacityof the `len` of the `numbersToSum` we need to work through. The length of a slice is the number of elements it holds `len(mySlice)`, while the capacity is the number of elements it can hold in the underlying array `cap(mySlice)`, e.g., `make([]int, 0 ,5)` creates a slice with length 0 and capacity 5.

You can index slices like arrays with `mySlice[N]` to get the value out or assign it a new value with `=`

The tests should now pass.

## Refactor

As mentioned, slices have a capacity. If you have a slice with capacity of 2 and try to do `mySlice[10] = 1` you will get a _runtime_ error.

However, you can use the `append` function which takes a slice and a new value, then returns a new slice with all the items in it.

```go
func SumAll(numbersToSum  ...[]int ) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}
```

In this implementation, we are worrying less about capacity. We start with an empty slice `sums` and append it to the result of `Sum` as we work through the varargs.

Our next requirement is to change `SumAll` to `SumAllTails`, where it will calculate the totals of the "tails" of each slice. The tail of a collection is all the items in the collection except for the first one (the "head").

## Write the test first

```go
func TestSumAllTails(t *testing.T) {
	got := SumAllTails([]int{0, 1, 2}, []int{3, 4, 5})
	want := []int{3, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## Try and run the test

`./sum_test.go:26:9: undefined: SumAllTails`

## Write the minimal amount of code to get the test to run and check the failing output

```go
func SumAllTails(numbersToSum  ...[]int ) []int {
	var sums []int

	return sums
}
```

```sh
--- FAIL: TestSumAllTails (0.00s)
    sum_test.go:35: got [] want [3 9]
```

## Write enough code to make it pass

```go
func SumAllTails(numbersToSum  ...[]int ) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}

	return sums
}
```

Slices can be sliced! The syntax is `slice[low:high]`. If you omit the value on one of the sides of the `:` it captures everything to that side of it. In our case, we are saying "take from 1 to the end" with `numbers[1:]`. You may wish to spend some time writing other tests around slices and experiment with the slice operator to get more familiar with it.

## Refactor

Not a lot to refactor this time.

What do you think would happen if you passed in an empty slice into our function? What is the "tail" of an empty slice? What happens when you tell Go to capture all elements from `myEmptySlice[1:]`?

## Write the test first

```go
func TestSumAllTails(t *testing.T) {
	t.Run("make the sum of some slices", func(t *testing.T){	got := SumAllTails([]int{0, 1, 2}, []int{3, 4, 5})
		want := []int{3, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("safely sum empty slices", func(t *testing.T){
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{3, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
```

## Try and run the test

```sh
--- FAIL: TestSumAllTails (0.00s)
    --- FAIL: TestSumAllTails/safely_sum_empty_slices (0.00s)
panic: runtime error: slice bounds out of range [1:0] [recovered, repanicked]
```

Oh no! It's important to note that while the test _has compiled_, it _has a runtime error_.

Compile time errors are our friend because they help us write software that works, runtime errors are our enemies because they affect our users.

## Write enough code to make it pass

```go
func SumAllTails(numbersToSum  ...[]int ) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}
```

## Refactor

Our tests have some repeated code around the assertions again, so let's extract those into a function.

```go
func TestSumAllTails(t *testing.T) {

	checkSums := func(t *testing.T, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sum of some slices", func(t *testing.T){
		got := SumAllTails([]int{0, 1, 2}, []int{3, 4, 5})
		want := []int{3, 9}

		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T){
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{3, 9}

		checkSums(t, got, want)
	})
}
```

We could've created a new function `checkSums` like we normally do, but in this case, we're showing a new technique, assigning a function to a variable. It might look strate but, it's no different to assigning a variable to a `string`, or an `int`, functions in effect are values too.

It's not shown here, but this technique can be useful when you want to bind a function to other local variables in "scope" (e.g between some `{}`). It also allows you to reduce the surface area of your API.

By defining this function inside the test, it cannot be used by other functions in this package. Hiding variables and functions that a don't need to be exported is an important design consideration.

A handy side-effect of this adds a little type-safety to our code. If a developer mistakenly adds a new test with `checkSums(t, got "dave") the compiler will stop them in their tracks.

```sh
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## Wrapping up

We have covered

- Arrays
- Slices
  - The various ways to make them
  - How they have a _fixed_ capacity but you can create new slices from old ones using `append`
  - How to slice, slices!
- `len` to get the length of an array or slice
- Test coverage tool
- `reflectDeepEqual` and why it's useful but can reduce the type-safety of your code

We've used slices and arrays with integers but they work with any other type too, including arrays/slices themselves. So you can declare a variable of `[][]string` if you need to.

Check out the Go blog post on slices for an in-depth look in slices. try writing more tests to solidify what you learn from reading it.

Another handy way to experiment with Go other than writing tests in the Go playground. You can try most things out and you can easily share your code if you need to ask questions.

[Here is an example](https://go.dev/play/p/bTrRmYfNYCp) of slicing an array and how changing the slice affects the original array; but a "copy" of the slice will not affect the original array. Another exapmle of why it's a good idea to make a copy of a slice after slicing a very large slice.
