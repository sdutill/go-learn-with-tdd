# Hello, World

It is traditional for your first program in a new language to be Hello, World.

- Create a folder wherever you like
- Put a new file in it called `hello.go` and put the following code inside it

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

To run it, type `go run hello.go`.

## How it works

When you write a program in Go, you will have a `main` package defined with a `main` func inside it. Packages are ways of grouping up related Go code together.

The `func` keyword defined a function with a name and a body.

With `import "fmt"` we are importing a package which contains the `Println` function that we use to print.

## How to test

How do you test this? It is good to separate your "domain" code from the outside world (side-effects). The `fmt.Println` is a side effect (printing to stdout), and the string we send in is our domain.

So let's separate these concerns so it's easier to test.

```go
package main

import "fmt"

func Hello() string {
    return "Hello, world"
}

func main() {
    fmt.Println(Hello())
}
```

We have created a new function with `func`, but this time, we've added another keyword, `string`, to the definition. This means the function returns a `string`.

Now create a new file called `hello_test.go` where we are going to write a test for our `Hello` function

```go
package main

import "testing"

func TestHello(t *testing.T){
    got := Hello()
    want := "Hello, world"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

## Go modules?

The next step is to run the tests. Enter `go test` in your terminal. If the tests pass, then you are probably using an earlier version of Go. However, if you are using Go 1.16 or later, the tests will likely not run. Instead, you will see an error message like this in the terminal:

```sh
$ go test
go: cannot find main module; see 'go help modules'
```

What's the problem? In a word, modules. Luckliy, the problem is easy to fix. Enter `go mod init example.com/hello` in your terminal. That will create a new file with the following contents:

```go
module example.com/hello

go 1.16
```

This file tells the `go` tool essential information about your code. If you planned to distribute your application, you would include where the code was available for download as well as information about dependencies. The name of the module, example.com/hello, usually refers to a URL where the module can be found and downloaded. For compatibility with tools we'll start using soon, make sure your module/s name has a dot somewhere in it, like the dot in .com of example.com/hello. For now, your module file is minimal, and you can leave it that way. To read more information about modules, [you can check out the reference in the Goland documentation](https://golang.org/doc/modules/gomod-ref). We can get back to testing and learning Go now since the tests should run, even on Go 1.16.

In future chapters, you will need to run `go mod init SOMENAME` in each new folder before running commands like `go test` or `go build`.

## Back to Testing

Run `go test` in your terminal. It should've passed! Just to check, try deliberately breaking the test by changing the `want` string.

Notice how you have not had to pick between multiple testing frameworks and then figure out how to install them. Everything you need is built into the language, and the syntax is the same as the rest of the code you will write.

### Writing tests

Writing a test is just like writing a function, with a few rules

- It needs to be in a file with a name like `xxx_test.go`
- The test function must start with the word `Test`
- The test function takes one argument only `t *testing.T`
- To use the `*testing.T` type, you need to `import "testing"`, like we did with `fmt` in the other file

For now, it's enough to know that yout `t` of type `*testing.T` is your "hook" into the testing framework so you can do things like `t.Fail()` when you want to fail.

We've covered some new topics:

### `if`

If statements in Go are very much like other programming languages.

### Declaring variables

We're declaring some variables with the syntax `varName := value`, which lets us reuse some values in our test for readability.

### `t.Errorf`

We are calling the `Errorf` _method_ on our `t`, which will print out a message and fail the test. The `f` stands for format, which allows us to build a string with values inserted into the placeholder values `%q`. When you make the test fail, it should be clear how it works.

You can read more about the placeholder strings in the [fmt documentation](https://pkg.go.dev/fmt#hdr-Printing). For tests, `%q` is veyr useful as it wraps your values in double quotes.

We will later explore the difference between methods and functions.

## Go's documentation

Another quality-of-life feature of Go is the documentation. We just saw the documentation for the fmt package at the official package viewing website, and Go also provides ways for quickly getting at the doumentation offline.

Go has a built-in tool, doc, which lets you examine any package installed on your system, or the module you're currently working on. To view that same documentation for the Printing verbs:

```bash
go doc fmt
package fmt // import "fmt"

Package fmt implements formatted I/O with functions analagous to C's printf and scanf. The format 'verbs' are derived from C's but are simpler.

# Printing

The verbs:

General:
    %v  the value in a default format
        when printing structs, the plus flag (%+v) adds field names.
    %#v a Go-syntax representation of the value
    %T  a Go-syntax representation of the type of the value
    %% a literal percent sign; consumes no value
...

```

Go's second tool for viewing documentatoin is the pkgsite sommand, which powers Go's official package viewing website. You can install pkgsite with `go install golang.org/x/pkgsite/cmd/pkgsite@latest`, then run it with `pkgsite -open .`. Go's install command will download the source files from that repository and build them into an executable binary. For a default installation of Go, that executable will be in `%HOME/go/bin` for Linux and macOS, and `%USERPROFILE%\go\bin` for Windows. If you have not already added those paths to your $PATH var, you might want to do so to make running go-installed commands easier.

The vast majority of the standard library has excellent documentation with examples. Navigating to [https://localhost:8080/testing](https://localhost:8080/testing) would be worthwhile to see what's available to you.

## Hello, YOU

Now that we have a test, we can iterate on our software safely.

In the last example, we wrote the test _after_
the code had been written so that you could get an example of how to write a test and declare a function. From this point on, we will be _writing tests first_.

Our next requirement is to let us specify the recipient of the greeting.

Let's start by capturing these requirements in a test. This is basic test-driven development and allows us ot make sure our test is _actually_ testing what we want. When you retrospectively write tests, there is the risk that your test may continue to pass even if the code doesn't work as intended.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Now run `go test`, you should have a compilation error.

```sh
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

When using a statically typed language like Go it is important to _listen to the compiler_. The compiler understands how your code should snap together and work so you don't have to.

In this case the compiler is telling you what you need to do to continue. We have to change our functoin `Hello` to accept an argument.

Edit the `Hello` function to accept an argument of type string

```go
func Hello(name string) string {
    return "Hello, world"
}
```

If you try and run you tests again your `hello.go` will fail to compile because you're not passing an argment. Send in "world" to make it compile.

```go
func main() {
    fmt.Println(Hello("world"))
}
```

Now when you run your tests, you should see something like

```sh
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

We finally have a compiling program but it's not meeting our requirements according to the test.

Let's make the test pass by using the name argument and concatenate it with `Hello, `

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

When you run the tests, they should now pass. Normally, as part of the TDD cycle, we should now _refactor._

## A note on source control

At this point, if you are using source control (which you should!) I would `commit` the code as it is. We have working software backed by a test.

I _wouldn't_ push it to main though, because I plan to refactor next. It is nice to commit at this point in case you somehow get into a mess with refacgtoring - you can always go back to the working version.

There's not a lot to refactor here, but we can introduce another language feature, _constants_.
