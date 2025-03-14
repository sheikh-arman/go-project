/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"time"
)

var (
	ck int
)

func main() {
	ck = 0
	//cmd.Execute()
	fmt.Println("Hello World")
	ctx, cancel := context.WithCancel(context.Background())
	go funTest(ctx)
	time.Sleep(1 * time.Second)
	//time.Sleep(4 * time.Second)
	cancel()
	fmt.Println("cancel")
	time.Sleep(6 * time.Second)
	fmt.Println("we wait 6 second so that we proved that go routine finished with out printing test!!")
	fmt.Println(ck)
}

func funTest(ctx context.Context) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine exiting due to cancellation\n\n\n\n\n\n")
			ck = 1
			return
		case <-time.After(3 * time.Second):
			ck = 2
			fmt.Println("This should not print if cancelled early\n\n\n\n")
		}
	}
}
