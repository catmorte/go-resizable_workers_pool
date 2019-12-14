package main

import (
	"awesomeProject2/pkg/wpool"
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	wp := wpool.NewPool(ctx, 1)
	var wg sync.WaitGroup
	amount := uint(101)
	wg.Add(int(amount))
	go func() {
		time.Sleep(time.Second)
		_ = wp.Resize(10)
		time.Sleep(time.Second)
		_ = wp.Resize(5)
		time.Sleep(time.Second)
		_ = wp.Resize(100)
	}()
	for i := uint(0); i < amount; i++ {
		wp.Do(toDo(&wg, i, wp))
	}
	wg.Wait()
	fmt.Println("The end")
}

func toDo(wg *sync.WaitGroup, i uint, wp wpool.WorkersPool) func() {
	return func() {
		time.Sleep(time.Second)
		fmt.Printf("do %v workers pool size %v\n", i, wp.GetSize())
		wg.Done()
	}
}
