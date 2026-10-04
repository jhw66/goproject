package main

import (
	"flag"
	"fmt"
)

func main() {
	var name = flag.String("name", "wjh", "the name to greet")
	var age = flag.Int("age", 30, "the age of the person")
	var verbose = flag.Bool("verbose", false, "enbale verbose output")

	flag.Parse()

	fmt.Printf("Hello,%s!You are %d years old.\n", *name, *age)
	if *verbose {
		fmt.Println("Verbose mode enabled.")
	}
}
