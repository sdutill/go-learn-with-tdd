# Concurrency

Here's the setup: a colleague has written a function, `CheckWebsites`, that checks the status of a list of URLs.

```go
package main

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = wc(url)
	}

	return results
}
```

It returns a map of each URL checked to a boolean value: `true` for a good respond; `false` for a bad response.

You also have to pass in a `WebsiteChecker` which takes a single URL and returns a boolean. This is used by the function to check all the websites.

Using depedency injection has allowed them to test the function without making real HTTP calls, making it reliable and fast.

Here's the test they've written:

```go
package main

import (
	"reflect"
	"testing"
)

func mockWebsiteChecker(url string) bool {
	return url != "waat://furhurterwe.geds"
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("wanted %v, got %v", want, got)
	}
}
```

The function is in production and being used to check hundreds of websites. But your colleagues has started to get complaints that it's slow, so they've asked you to help speed it up.

## Write a test

Let's use a benchmark to test the speed of `CheckWebsites` so that we can see the effect of our changes.

```go
package main

import (
	"testing"
	"time"
)

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites( b *testing.B) {
	urls :=make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}

	for b.Loop() {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}
```

The benchmark tests `CheckWebsite` using a slice of one hunderd urls and uses a new fake implementation of `WebsiteChecker`. `slowStubWebsiteChecker` is deliberately slow. It uses `time.Sleep` to wait exactly twenty milliseconds and then it returns true.

When we run the benchmark using `go test -bench=.` ( if you're in Windows Powershell `go test -bench="."`):

```sh
goos: windows
goarch: amd64
pkg: github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1
cpu: Intel(R) Core(TM) i7-6600U CPU @ 2.60GHz
BenchmarkCheckWebsites-4               1        2088976600 ns/op
PASS
ok      github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1 2.538s
```

`CheckWebsites` has been benchmarked at `2088976600` nanoseconds - about two and a quarter seconds.

Let's try and make this faster.

### Write enough code to make it pass

Now we can finally talk about concurrency which, for the purposes of the following, means "having more than one thing in progress." This is something that we do natually everyday.

For instance, this morning I made a cup of tea. I put the kettle on and then, while I was waiting for it to boil, I got the milk out of the fridge, got the tea out of the supboard, found my favourite mug, put the teabag into the cup and then, when the kettle had boiled, I put the water in the cup.

What I _didn't_ do was put the kettle on and then stand there blankly staring at the kettle until it boiled, then do everything else once the kettle had boiled.

If you can understand why it's faster to make the tea the first way, then you can understand how we will make `CheckWebsites` faster. Instead of waiting for a website to respond before sending a request to the next website, we will tell our computer to make the next request while it is waiting.

Normally in Go when we call a function `doSomething()` we wait for it to return (even if it has no value to return, we still wait for it to finish) We say that this operation is _blocking_ - it makes us wait for it to finish. An operation that does not block in Go will run in a separate _process_ called a _goroutine_. Think of a process as reading down the page of Go code from top to bottom, going 'inside' each function when it gets called to read what it does. When a separate process starts, it's like anothe reader begins reading inside the function, leaving the original reader to carry on going down the page.

To tell Go to start a new goroutine we turn a function call into a `go` statement by putting the keyword `go` in front of it: `go doSomething()`.

```go
package main

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	return results
}
```

Because the only way to start a goroutine is to put `go` in front of a function call we often use _anyonamous functions_ when we want to start a goroutine. An anonymous function literal looks just the same as a normal function declaration, but without a name (unsurprisingly). You can see one above in the body of the `for` loop.

Anonymous functions have a number of features which make them useful, two of which we're using above. Firstly, they can be executed at the same time that they're declared - this is what the `()` at the end of the anonymous function is doing. Secondly they maintain access to the lexical scope in which they are defined - all the variables that are available at the point when you declare the anonymous function are also available in the body of the function.

The body of the anonymous function above is just the same as the loop body was before. The only difference is that each iteration of the loop will start a new goroutine, concurrent with the current process (the `WebsiteChecker` function). Each goroutine will add its result to the results map.

But when we run `go test`:

```sh
--- FAIL: TestCheckWebsites (0.00s)
    check_websites_test.go:28: wanted map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false], got map[http://blog.gypsydave5.com:true http://google.com:true waat://furhurterwe.geds:false]
FAIL
```

### A quick aside into the concurrency universe

You might not get this result. You might get a panic message that we're going to talk about in a bit. Don't worry if you got that, just keep running the test until you _do_ get the result above. Or pretend that you did. Up to you. Welcome to concurrency: when it's not handled correctly it's hard to predict what's going to happen. Don't worry - that's why we're writing tests, to help us know when we're handling concurrency predictably.

### ... and we're back.

We are caught by the original test `CheckWebsites`, its not returning an empty map. What went wrong?

None of the goroutines that our `for` loop started had enough time to add their result to the `results` map; the `CheckWebsites` function is too fast for them, and it returns the still empty map.
To fix thiis we can just wait while all the goroutines do their work, and then return. Two seconds ought to do it, right?

```go
package main

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
results := make(map[string]bool)

    for _, url := range urls {
    	go func() {
    		results[url] = wc(url)
    	}()
    }

    time.Sleep(2 * time.Second)

    return results

}

```

Now if you're lucky you'll get:

```sh
PASS
ok      github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1 2.539s
```

But if you're unlucky (this is more likely if you run them with the benchmark as you'll get more tries)

```sh
fatal error: concurrent map writes

goroutine 116 [running]:
internal/runtime/maps.fatal({0x7ff644a54c1c?, 0x7ff644a28160?})
        C:/Program Files/Go/src/runtime/panic.go:1181 +0x18
github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.CheckWebsites.func1()
        C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_website.go:12 +0x54
created by github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.CheckWebsites in goroutine 22
        C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_website.go:11 +0x51
exit status 2
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1 2.513s
```

This is long and scary, but all we need to do is take a breath and read the stacktrace: `fatal error: concurrent map writes`. Sometimes, when we run our tests, two of the goroutines write to the results map at exactly the same time. Maps in Go don't like it when more than one thing tries to write to them at once, and so `fatal error`.

This is a _data race_, a bug that occurs when two or more goroutines access the same memory location concurrently, and at least one of those accesses is a write. Because we cannot control exactly when each goroutine executes, we are vulnerable to multiple goroutines trying to write to the `results` map at the exact same time. Go maps are not safer for concurrent writes, so the runtime throws a fatal error to prevent memory corruption.

Go can help us to spot race conditions with its built in [race detector](https://blog.golang.org/race-detector). To enable this feature, run the tests with the `race` flag: `go test -race`.

You should get some output that looks like this:

```sh
==================
WARNING: DATA RACE
Write at 0x00c000024750 by goroutine 9:
  runtime.mapassign_faststr()
      C:/Program Files/Go/src/internal/runtime/maps/runtime_faststr.go:263 +0x0
  github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.CheckWebsites.func1()
      C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_website.go:12 +0x77

Previous write at 0x00c000024750 by goroutine 10:
  runtime.mapassign_faststr()
      C:/Program Files/Go/src/internal/runtime/maps/runtime_faststr.go:263 +0x0
  github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.CheckWebsites.func1()
      C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_website.go:12 +0x77

Goroutine 9 (running) created at:
  github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.CheckWebsites()
      C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_website.go:11 +0x4a
  github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1.TestCheckWebsites()
      C:/Users/shawn/src/go-learn-with-tdd/src/go_fundamentals/concurrency/v1/check_websites_test.go:25 +0x168
  testing.tRunner()
      C:/Program Files/Go/src/testing/testing.go:2036 +0x1ca
  testing.(*T).Run.gowrap1()
      C:/Program Files/Go/src/testing/testing.go:2101 +0x38
```

The details are, again, hard to read - but `WARNING: DATA RACE` is pretty unambiguous. Reading into the body of the error we can see two different goroutines performing writes on a map:

`Write at 0x00c420084d20 by goroutine 8:`
is writing to the same block of memory as
`Previous write at 0x00c420084d20 by goroutine 7:`

On top of that, we can see the line of code where the write is happening:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12`

and the line of code where goroutines 7 and 8 are started:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11`

Everything you need to know is printed to your terminal - all you have to do is be patient enough to read

### Channels

We can solve this data race by coordinateing our goroutines using _channels_. Channels are a Go data structure that can both receive and send values. These operations, along with their details, allow communication between different processes.

In this case we wan to think about the communication between the parent process and each of the goroutines that it makes to do the work of running the `WebsiteChecker` function with the url.

```go
package main

type WebsiteChecker func(string) bool
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {

	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func() {
			resultChannel <- result{url, wc(url)}
		}()
	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}
```

Alongside the `results` map we now have a `resultChannel`, which we `make` in the same way. `chan result` is the type of the channel - a channel of `result`. The new type, `result` has been made to associate the return value of the `WebsiteChecker` with the url being checked - it's a struct of `string` and `bool`. As we don't need either value to be named, each of them is anonymous within the struct; this can be useful when it'ds hard to know what to name a value.

Now when we iterate over the urls, instead of writing to the `map` directly we're sending a `result` struct for each call to `wc` to the `resultchannel` with a _send statement_. this uses the `<-` operator, taking a channel on the left and a value on the right:

```go
// Send statement
resultChannel <- result{url, wc(url)}>
```

The next `for` loop iterates once for each of the urls. Inside we're using a _receive expression_, which assigns a value received from a channel to a variable. This also uses th e`<-` operator, but with the two operands now reverdes: the channel is now on he right and the variable that we're assigning to is on the left:

```go
// Send statement
r := <-resultChannel
```

We then use the `result` received to update the map.

By sending the results into a channel, we can control the timing of each write into the results map, ensuring that it happens once at a time. Although each of the calls of `wc`, and each send to the result channel, is happening concurrently inside its own process, each of the results is being dealt with one at a time as we take values out of the result channel with the receive expression.

We have used concurrency for te part of the coede that we wanted to make faster, while making sure that the part that cannot happen simultaneously still happens linearly. And we have communicated across the multiple processes involved by using channels.

When we run the benchmark:

```sh
goos: windows
goarch: amd64
pkg: github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1
cpu: Intel(R) Core(TM) i7-6600U CPU @ 2.60GHz
BenchmarkCheckWebsites-4              54          20623796 ns/op
PASS
ok      github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/concurrency/v1 1.619s
```

20623796 nanoseconds - 0.023 seconds, about one hundred times as fast as original function. A great success.

## Wrapping up

This exercise has been a little lighted on the TDD than usual. Ina way we've been taking part in one oong refactoring of the `CheckWebsites` function; the inputs and outputs never changes, it just got faster. But the tests we had in place, as well as the benchmark we wrote, allowed us to refactor `CheckWebsites` in a way that maintained confidence that the software was still working, while demonstrating that it had actually become faster.

In making it faster we learned about

- _goroutines_, the basic unit of concurrency in Go, which let us manage more than one website check request.
- _anonymous functions_, which we used to start each of the concurrent processes that check websites.
- _channels_, to help organize and control the communication between the different processes, allowing us to avoid a _race condition_ bug.
- _the race detector_ which helped us debug problems with concurrent code

### Make it fast

On formulation of an agile way of building software, often misattributed to Kent Beck, is:

> [Make it work make it right, make it fast](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)

Where 'work' is making the tests pass, 'right' is refactoring the code, and 'fast' is optimizing the code to make it, for example, run quickly. We can only 'make it fast' once we've made it work and made it right. We were lucky that the code we were given was already demonstrated to be working, and didn't need to be refactored. We should never try to 'make it fast' before the other two steps have been performed because

> [Premature optimization is the root of all evil](https://wiki.c2.com/?PrematureOptimization) -- Donand Knuth
