package main

import (
	"fmt"
)

func testingFunc(val *int) *int {
	fmt.Printf("%v : is where original int passed was pointing at value of %v\n", val, *val)
	randint := 5095328
	val = &randint
	return val
}

func main() {
	fmt.Println("Hello World?")
	testint := 1994
	var newVal = testingFunc(&testint)
	fmt.Printf("%v : is where original int passed was pointing at value of %v\n", newVal, *newVal)
}
