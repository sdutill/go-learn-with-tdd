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

We will need some kind of _balance_ variable in our struct to store the state

```go
type Wallet struct {
balance float64
}
```

In Go if a symbol (variables, types, functions, et al) starts with a lowercase symbol then it is provate _outside the package it's defined in_.

In our case we want our methods to be able to manipulate this value, but no one else.

Remember we can access the internal `balance` field in the struct using the "receiver" variable.

```go
package pointers_and_errors

type Wallet struct {
	balance float64
}

func (w Wallet) Deposit(amount float64) {
	w.balance += float64(amount)
}

func (w Wallet) Balance() float64 {
	return w.balance
}
```

With our career in fintech secured, run the test suite and bask in the passing test.

```sh
--- FAIL: TestWalet (0.00s)
    wallet_test.go:14: got 0.00 want 10.00
FAIL
exit status 1
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/pointers-and-errors    0.491s
```

## That's not quite right

Well this is confusing, our code looks like it should work. We add the new amount onto our balance and then the balance method should return the current state of it.

In Go, **when you call a function or a method the arguments are** **_copied_**.

When calling `func (w Wallet) Deposit(amount float64)` the `w` is a copy of whatever we called the method from.

Without getting too computer-sciency, when you create a value - likea wallet, it is stored somewhere in memory. You can find out what the _address_ of that bit of memory is with `&myVal`

Experiment by adding some prints to your code

````go
func TestWalet(t *testing.T) {
	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()

	fmt.Printf("address of balance in test is %p \n", &wallet.balance)

	want := 10.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

```go
func (w Wallet) Deposit(amount float64) {
	fmt.Printf("address of balance in Deposit is %p \n", &w.balance)
	w.balance += float64(amount)
}
````

The `%p` placeholder prints memory address in base 16 notation with leading `0x`s and the escape character prints a enw line. Note that we get the pointer (memory address) of somethingby placing an `&` character at the beginning of the symbol.

Now re-run the test

```sh
address of balance in Deposit is 0x2161b48c8160
address of balance in test is 0x2161b48c8158
--- FAIL: TestWalet (0.00s)
    wallet_test.go:20: got 0.00 want 10.00
```

You can see that the addresses of the two balances are different. So when we change the value of teh balance inside the code, we are working on a copy of what came from the test. Therefore the balance in the test is unchanged.

We can fix this with _pointers_. Pointers let us _point_ to some values and then let us change them. So rather than taking a copt of the whole Wallet, we instead take a pointer to that wallet so that we can change the original values within it.

```go
package pointers_and_errors

import "fmt"

type Wallet struct {
	balance float64
}

func (w *Wallet) Deposit(amount float64) {
	fmt.Printf("address of balance in Deposit is %p \n", &w.balance)
	w.balance += float64(amount)
}

func (w *Wallet) Balance() float64 {
	return w.balance
}
```

The difference is the receiver type is `*Wallet` rather than `Wallet` which you can read as "a pointer to a wallet".

Try and re-run the tests and they should pass.

Now you might wonder, why did they pass? We didn't dereference the point in the function, like so:

```go
func (w *Wallet) Balance() float64 {
	return (*w).balance
}
```

and seemingly addressed the object directly. In fact, the code above using (*w) is absolutely valid. Howeer, the makes of Go deemed this notation cumberson, so the language permits us to write `w.balance`, without an explicit dereference. The3se pointers to structs even have their own name: *struct pointers\* and they are automatically dereferenced.

Technically you do not need to change `Balance` to use a pointer receiver as taking a copy of the balance is fine. However, by convention you should keep your method receiver types the same for consistency.

## Refactor

We said we were making a Bitcoin wallet but we have not mentiond them so far. We've been using `int` because they're a good type for counting things! It seems a bit overkill to create a struct for this. `int` is fine in terms of the way it works but it's not descriptive.

Go lets you create new types from existing ones.

The syntax is `type MyName OriginalType`

```go
package pointers_and_errors

import "fmt"

type Bitcoin float64

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	fmt.Printf("address of balance in Deposit is %p \n", &w.balance)
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
```

```go
func TestWalet(t *testing.T) {
	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()

	fmt.Printf("address of balance in test is %p \n", &wallet.balance)

	want := Bitcoin(10.0)

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

To make `Bitcoin` you just use the syntax `Bitcoin(999)`

By doing this we're making a new type and we can declare _methods_ on them. This can be very useful when you try to add some domain specific functionality on top of existing types.

Let's implement Stringer on Bitcoin.

```go
type Stringer interface{
	String() string
}
```

This interface is defined in the `fmt` package and lets you define how your type is printed when used with the `%s` format string in prints.

```go
func (b Bitcoin) String() string {
	return fmt.Sprintf("%.2f BTC", b)
}
```

As you can see, the syntax for creating a method on a type declaration is the same as it is on a struct.

Next we need to update our test format strings so they wil use `String()` instead.

```go
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
```

To see this in action, deliberately break the test so we can see it

```sh
address of balance in Deposit is 0x34e5f0cac158
address of balance in test is 0x34e5f0cac158
--- FAIL: TestWalet (0.00s)
    wallet_test.go:20: got 10.00 BTC want 20.00 BTC
FAIL
```

This makes it clearer what's going on in our test.

The next requirement is for a `Withdraw` function.

## Write the test first

```go
func TestWalet(t *testing.T) {
	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		got := wallet.Balance()

		fmt.Printf("address of balance in test is %p \n", &wallet.balance)

		want := Bitcoin(10.0)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.Withdraw(Bitcoin(10))

		got := wallet.Balance()

		fmt.Printf("address of balance in test is %p \n", &wallet.balance)

		want := Bitcoin(10.0)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
```

## Try to run the test

```sh
.\wallet_test.go:27:10: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/pointers-and-errors [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (w *Wallet) Withdraw(amount Bitcoin) {

}
```

```sh
--- FAIL: TestWalet (0.00s)
    --- FAIL: TestWalet/withdraw (0.00s)
        wallet_test.go:36: got 20.00 BTC want 10.00 BTC
FAIL
exit status 1
```

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
	w.balance -= amount
}
```

## Refactor

There's some duplication in our tests, lets refactor that out.

```go
func TestWalet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()

		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))

	})
	t.Run("withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})
}
```

What should happen if you try to `Withdraw` more than is left in the account? For now, our requirement is to assume there is not an overdraft facility.

How do we signal a problem when using `Withdraw`?

In Go, if you want to indicate an error it is idiomatic for your function to return an `err` for the caller to check and act on.

Lets try this out in a test.

## Write the test first

```go
	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{balance: startingBalance}
		err := walet.withdraw(Bitcoin(100))
		assertBalance(t, wallet, Bitcoin(20))

		if err == nil {
			t.Error("wanted an error but didn't get one")
		}

	})
```

We want `Withdraw` to return an error _if_ you try to take out more than you have and the balance should stay the same.

We then check an error has returned by failing the test if it is `nil`.

`nil` is synonymous with `null` from other programming languages. Error can be `nil` because the return type of `Withdraw` with be `error`, which is an interface. If you see a function that takes arguments or returns values that are interfaces, they can be nillable.

Like `null` if you try to access a value that is `nil` it will throw a **runtime panic**. This is bad! You should make sure that you check for nils.

## Try and run the test

```sh
.\wallet_test.go:32:10: wallet.Withdraw(Bitcoin(100)) (no value) used as value
FAIL    github.com/sdutill/go-learn-with-tdd/src/go_fundamentals/pointers-and-errors [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.balance -= amount
	return nil
}
```

Again, it is very important to just write enough code to satisfy the compiler. We correct our `Withdraw` method to return `error` and for now we have to return _something_ so let's just return `nil`.

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}
```

Remember to import `errors` into your code.

`errors.New` creates a new `error` with a message of your choosing.

## Refactor

Lets make a quick test for our error check to improve the test's readability

```go
	assertError := func(t *testing.T, got error, want string) {
		t.Helper()

		if got == nil {
			t.Error("wanted an error but didn't get one")
		}

	}
```

And in our test

```go
	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, "cannot withdraw, insufficient funds")
		assertBalance(t, wallet, Bitcoin(20))
	})
```

Hopefully when returning an error of "oh no" you were thinking hat we _might_ iterate on that because it doesn't seem that useful to return.

Assuming that the error ultimately gets returned to the user, lets update our test to asset on some kind of error message rather than just the existence of an error.

```go
	assertError := func(t *testing.T, got error, want string) {
		t.Helper()

		if got == nil {
			t.Fatal("wanted an error but didn't get one")
		}

		if got.Error() != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}

	}
```

As you can see `Error`s can be converted to a string with the `.Error()` method, which we do in order to compare it with the string we want. We are also making sure that the error is not `nil` to ensure we don't call `.Error()` on `nil`.

And then update the caller

```go
	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, "cannot withdraw, insufficient funds")
		assertBalance(t, wallet, Bitcoin(20))
	})
```

We've introduced `t.Fatal` which will stop the test if it is called. This is because we don't want to make any more assertions on the error returned if there isn't one around. Without this the test would carry on to the next step and panic because of a nil pointer.

## Try to run the test

```sh
--- FAIL: TestWallet (0.00s)
    --- FAIL: TestWallet/withdraw_insufficient_funds (0.00s)
        wallet_test.go:47: got 'oh no', want 'cannot withdraw, insufficient funds'
FAIL
exit status 1
```

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}

	w.balance -= amount
	return nil
}
```

## Refactor

We have dupliaction of the error message in both the test code and the `Withdraw` code.

It would be really annoying for the test to fail if someone wanted to re-word the error and it's just too much detail for our test. We don't _really_ care what the exact wording is, just that some kind of meaningful error around withdrawin is returned given a certain confition.

In Go, errors are values, so we can refactor it out into a variable and have a single source of truth for it.

```go
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount
	return nil
}
```

The `var` keyword allows us to define values global to the package.

This is a positive change in itself because now our `Withdraw` functions looks very clear.

Next we can refactor our test code to use this value instead of specific strings.

```go
func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, Bitcoin(20))
	})
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

And now the test is easier to follow too.

I have moved the helpers out of the main test function just so when someone opens up a file they can start reading our assertions first, rather than some helpers.

Another useful property of tests is that they help us understand the _real_ usage of our code so we can make sympathetic code. We can see here that a developer can simply call our code and do an equals check to `ErrInsufficientFunds` and act accordingly.

## Unchecked errors

Whilst the Go compiler helps you a lot, sometimes there are things you can still miss and error handling can sometimes be tricky.

There is one scenario we have not tested. To find it, run the following in a terminal to install `errcheck`, one of many linters available for Go.

`go install github.com/kisielk/errcheck@latest`

Then, inside the directory with your code run `errcheck .`

You should get something like

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

What this is telling us is that we have not checked the error being returned on that line of code. That line of code on my computer corresponds to our normal withdraw scenario because we have not checked that if the `Withdraw` is successful that en error is _not_ returned.

Here is the final test code that accounts for this.

```go
func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, Bitcoin(20))
	})
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
```

## Wrapping up

### Pointers

- Go copies values when you pass them to functions/methods, so if you're writing a function that needs to mutate state you'll need to take a pointer to the thing you want to change.
- The fact that Go takes a copy of values is useful a lot of the time but cometimes you won't want your system to make a copy of something, in which case you need to pass a reference. Examples incldue referencing very large data structures or things where only one instance is necessary (like database connection pools).

### nil

- Pointers can be nil
- When a function returns a pointer to something, you need to make sure you check if it's nil or you might raise a runtime exception - the compiler won't help you there
- Useful for when you want to describe a value that could be missing

### Errors

- Error are the way to signify failure when alling a function/method.
- By listening to our tests we concluded that checking for a string in an error would result in a flaky test. So we refactored our implementation to use a meaningful value instead and this resulted in easier to test code and concluded this would be easier for users of our API too.
- This is not the end of the story with error handling, you can do more sohpisticated things but this is just an intro. Later sections will cover more strategies
- [Don't jsut check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

## Create new types from existing ones

- Useful for adding more domain specific meaning to values
- Can let you implement interfacees

Pointers and errors are a big part of writing Go that you need to get compfrtable with. Thankfully the compiler will _usually_ help you out if you do something wrong, just take your time and read the error.
