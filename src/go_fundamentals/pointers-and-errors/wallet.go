package pointers_and_errors

type Wallet struct{}

func (w Wallet) Deposit(amount int) {}

func (w Wallet) Balance() int {
	return 0
}
