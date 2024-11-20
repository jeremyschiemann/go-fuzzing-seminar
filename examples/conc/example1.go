package main

import (
	"fmt"
	"sync"
)

// CalculateSquare calculates the square of a number and sends it to a channel.
func CalculateSquare(num int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- num * num
}

// SumSquares calculates the sum of squares for a slice of numbers.
func SumSquares(numbers []int) int {
	ch := make(chan int)
	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)
		go CalculateSquare(num, ch, &wg)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	sum := 0
	for square := range ch {
		sum += square
	}
	return sum
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	sum := SumSquares(numbers)
	fmt.Printf("Sum of squares: %d\n", sum)
}
