package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

type Myflag struct {
	val int
}

func (f *Myflag) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	f.val = val
	return nil
}
func (f *Myflag) String() string {
	return fmt.Sprintf("%d", f.val)
}

func main() {
	var name = flag.String("name", "wjh", "the name to greet")
	var age = flag.Int("age", 30, "the age of the person")
	var verbose = flag.Bool("verbose", false, "enbale verbose output")
	var tail = flag.Float64("tail", 178, "your hegiht")
	var times = flag.Duration("time", time.Duration(time.Now().Unix())*time.Second, "now time")
	var myFlag Myflag
	flag.Var(&myFlag, "myFlag", "a custom flag")

	flag.Parse()

	fmt.Printf("Hello,%s!You are %d years old and %f tail.\n", *name, *age, *tail)
	if *verbose {
		fmt.Println("Verbose mode enabled.")
	}

	fmt.Println(flag.Lookup("time").Value.String())
	fmt.Println(*times)
	fmt.Println(flag.Lookup("myFlag").Value.String())
	flag.PrintDefaults()
	fmt.Printf("Number of falgs:%d\n", flag.NFlag())
}
