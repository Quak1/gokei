package main

import (
	"fmt"
	"os"
)

func requireFlag[T comparable](name string, value T) {
	var zero T
	if value == zero {
		fmt.Printf("Error: -%s flag is required\n", name)
		os.Exit(1)
	}
}
