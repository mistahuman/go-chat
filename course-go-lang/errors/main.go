package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func getValue(key string) (string, error) {
	if key != "foo" {
		return "", fmt.Errorf("key %q: %w", key, ErrNotFound)
	}
	return "bar", nil
}

func main() {
	val, err := getValue("baz")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Key not found, continuing...")
			return
		}
		fmt.Println("Unexpected error:", err)
		return
	}
	fmt.Println("Value:", val)
}
