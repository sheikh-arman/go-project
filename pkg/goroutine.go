package pkg

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func Kitchen(chef int, newItem <-chan int) {
	defer wg.Done()
	for i := range newItem {
		fmt.Println("Order Number: ", i, "Selected chef", chef)
	}
}

func Test() {
	orderedItem := make(chan int)
	for chef := 0; chef < 5; chef++ {
		wg.Add(1)
		go Kitchen(chef, orderedItem)
	}
	for foodItem := 0; foodItem < 15; foodItem++ {
		orderedItem <- foodItem
	}
	close(orderedItem)
	wg.Wait()
}
