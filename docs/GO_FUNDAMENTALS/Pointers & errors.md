# Pointers & errors

We learned about structs in the last section which let us capture a number of values related around a concept.

At some point you ay wish to use structs to manage state, exposing methods to let users change the state in a way you can control.

**Fintech loves Go** and uhhh bitcoins? So let's show what an amazing banking system we can make.

Let's make a `Wallet` struct which lets us deposit `Bitcoin`.

## Write the test first

```go
package pointers_and_errors

import "testing"

func TestWalet(t *testing.T) {
	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

In the previous example we accessed fields directly with the field name, however in our _very secure wallet_ we don't want to expose our inner state to the rest of the world. We want to control access via methods.

## Try to run the test

```sh
.\wallet_test.go:6:12: undefined: Wallet
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/pointers-and-errors [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

The compiler doesn't know what a `Wallet` is so let's tell it.

```go
type Wallet struct {}
```

Now we've made our wallet, try and run the test again

```sh
.\wallet_test.go:8:9: wallet.Deposit undefined (type Wallet has no field or method Deposit)
.\wallet_test.go:10:16: wallet.Balance undefined (type Wallet has no field or method Balance)
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/pointers-and-errors [build failed]
```

We need to define those methods.

Remember to only do enough to make the tests run. We need to make sure our test fails correctly with clear error message.

```go
package pointers_and_errors

type Wallet struct{}

func (w Wallet) Deposit(amount int) {}

func (w Wallet) Balance(amount int) int {
	return 0
}
```

If this syntax is unfamiliar go back and read the structs section.

The tests should now compile and run

```sh
--- FAIL: TestWalet (0.00s)
    wallet_test.go:14: got 0 want 10
FAIL
```

## Write enough code to make it pass
