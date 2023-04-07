package main

import "fmt"

func main() {
	// s1 := make([]int, 3, 10)
	// fmt.Printf("%v\n", s1)
	// fmt.Println(len(s1))
	// fmt.Println(cap(s1))
	// s1 = append(s1, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100)
	// fmt.Printf("%v\n", s1)
	// fmt.Println(len(s1))
	// fmt.Println(cap(s1))

	// var s1 [5]int
	// var s1 = [5]int{0, 0, 0, 0, 0}
	// s1 := [5]int{0, 0, 0, 0, 0}
	// fmt.Printf("%v", s1)

	// fmt.Print("hello")
	fmt.Println(fib(7))
}

// fib sequence
func fib(a int) int {
	if a == 0 {
		return 0
	}
	if a == 1 {
		return 1
	}
	return fib(a-1) + fib(a-2)
}

//
