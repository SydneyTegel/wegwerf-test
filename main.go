// Package main ist der Einstiegspunkt von wegwerf-test.
package main

import "fmt"

// Hello returns a greeting message.
func Hello() string {
	return "hello"
}

func main() {
	fmt.Println(Hello())
}
