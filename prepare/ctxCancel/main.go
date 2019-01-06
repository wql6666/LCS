package main

import (
	"context"
	"fmt"
	"time"
)

func test() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("i exit")
					return //return not to leak the goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //cancel when we are finished consuming integers
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}

	}
}

func main() {
	test()
	time.Sleep(time.Hour)
}
