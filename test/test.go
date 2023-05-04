package main

import "fmt"

func GenerateNatural() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func PrimeFilter(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

func main() {
	var n int
	fmt.Scanf("%d", &n)
	ch := GenerateNatural()
	for i := 0; i < n; i++ {
		prime := <-ch
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ch, prime)
	}
}
