package main

import "fmt"

const (
	spanish = "Spanish"
	french  = "French"
	german  = "German"

	englishHelloPrefix = "Hello, "
	frenchHelloPrefix  = "Bonjour, "
	germanHelloPrefix  = "Gutentag, "
	spanishHelloPrefix = "Hola, "
)

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case german:
		prefix = germanHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

func main() {
	fmt.Println(Hello("", ""))
}
