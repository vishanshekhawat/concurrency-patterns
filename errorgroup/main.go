package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	withoutCtx()
	withCtx()
}

func withoutCtx() {
	var eg errgroup.Group
	for i := 1; i <= 10; i++ {
		i := i
		eg.Go(func() error {
			return run(i)
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Go Routine:", runtime.NumGoroutine())
	select {}
}

func withCtx() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel the context when the main function exits

	// Create an error group
	eg, ctx := errgroup.WithContext(ctx)

	for i := 1; i <= 10; i++ {
		i := i
		eg.Go(func() error {
			return runCtx(ctx, i)
		})
	}
	cancel()

	// Wait for all goroutines to finish and check for errors
	if err := eg.Wait(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All goroutines completed successfully")
	}

	fmt.Println("Go Routine:", runtime.NumGoroutine())
	select {}
}

func runCtx(ctx context.Context, i int) error {
	fmt.Println("Starting goroutine:", i)

	// Wait for a duration of i seconds or until the context is cancelled
	select {
	case <-time.After(time.Duration(i) * time.Second):
		fmt.Println("Goroutine", i, "completed")
		return nil
	case <-ctx.Done():
		fmt.Println("Goroutine", i, "cancelled")
		return ctx.Err()
	}
}

func run(i int) error {
	defer func() {
		fmt.Println("Ending:", i)
	}()
	if i == 2 {
		return fmt.Errorf("check error:%d", i)
	}

	if i == 4 {
		return fmt.Errorf("check error:%d", i)
	}
	fmt.Println("Running ", i)
	time.Sleep(time.Duration(i) * time.Second)
	return nil
}
