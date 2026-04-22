# Integers

Integers work as you would expect. Let's write an `Add` function to try things out. Create a test file called `adder_test.go` and write this code.

**Note:** Go source files can only have one package per directory. Make sure that your files are organised into their own packages. [Here is a good explanation on this.](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)

Your project directory might look something like this:

```txt
learnGoWithTests
    |
    |-> helloworld
    |    |- hello.go
    |    |- hello_test.go
    |
    |-> integers
    |    |- adder_test.go
    |
    |- go.mod
    |- README.md
```

## Write the test first

```go
package integers

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
```

You will notice that we're using `%d` as our format strings rather than `%q`. That's because we want it to print an integer rather than a string.

Also note that we are no longer using the main package, instead we've defined a package named `integers`, as the name suggests this will group functions for working with integers such as `Add`.

## Try and run the test

Run the `go test`

Inspect the compilation error
`.\adder_test.go:6:9: undefined: Add`

## Write the minimal amount of code for the test to run and check the failing test output

Write enough code to satisfy the compiler _and that's all_ - remember we want to check that our tests fail for the correct reason.

```go
package integers

func Add(x, y int) int {
    return 0
}
```

Remember, when you have more than one argument of the same type (in out case two integers) rather than having `(x int, y int)` you can shorted it to `(x, y int)`.

Now run the tests, and we should be happy that the test is correctly reporting what is wrong.

`adder_test.go:10: expected 4 but got 0`

If you have noticed we learnt about _named return value_ in the last section but aren't using the same here. It should be generally used when the meaning of the result isn't clear from context, in our case it's pretty much clear that `Add` function will add the parameters. You can refer this wiki for more details.

## Write enough code to make it pass

In the stricted sense of TDD we should now write the _minimal amount of code to make the test pass._ A pedantic programmer may do this

```go

func Add(x, y int) int {
    return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, wth some different numbers to forece that test to fail but that feels like a game of cat and mouse.

Once we're more familiar with Go's syntax I will introdce a technique called _Propery Based Testing_, which would stop annoying developers and help you find bugs.

For now, let's fix it properly

```go

func Add(x, y int) int {
    return x + y
}
```

If you re-run the tests they should pass.

## Refactor

There's not a lot in the _actual_ code we can really improve on here.

We explored earlier how by naming the return argument it appears in the documentation but also in most developer's text editors.

This is great because it aids the usability od code you are writing. It is preferable that a user can understand the usage of your code by just looking at the type signature and documentation.

You can add documentation to functions with comments, and these will appear in Go Doc just like when you look at the standard library's documentation.

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
    return x + y
}
```
