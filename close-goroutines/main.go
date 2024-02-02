package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func readDataDONE(done chan struct{}, file string) <-chan string {
	f, err := os.Open(file) //opens the file for reading
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string) //channel declared

	//returns a scanner to read from f
	fileScanner := bufio.NewScanner(f)

	//loop through the fileScanner based on our token split
	go func() {
		defer func() {
			fmt.Println("Done readData")
		}()
		for fileScanner.Scan() {
			val := fileScanner.Text() //returns the recent token

			// select
			select {
			case <-done:
				fmt.Println("DONE--> CHANNEL")
				close(out)
				return
			case out <- val: //passed the token value to our channel
			default:
			}
		}
		fmt.Println("DONE--> readData")
		close(out) //closed the channel when all content of file is read
		//closed the file
		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close an opened file: %v\n", err.Error())
			return
		}
	}()

	return out
}

func fanInMergeDataDONE(done chan struct{}, chs ...<-chan string) chan string {

	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	for _, ch := range chs {
		go func(ch <-chan string) {
			defer wg.Done()
			for val := range ch {
				// select
				select {
				case <-done:
					fmt.Println("Done fanInMergeData")
					return
				case chRes <- val:
				}
			}

		}(ch)
	}

	go func() {

		wg.Wait() //waits till the goroutines are completed and wg marked Done
		fmt.Println("Done Wait Go Routine")
		close(chRes) //close the result channel
	}()

	return chRes
}

func doneChanCloseExample() {

	fmt.Println("Starting Go Routine", runtime.NumGoroutine())
	done := make(chan struct{})
	//receive data from multiple channels and place it on result channel - FanIn
	chRes := fanInMergeDataDONE(done, readDataDONE(done, "text1.txt"))

	fmt.Println("After Merge Go Routines", runtime.NumGoroutine())

	//some logic with the result channel

	timer1 := time.NewTimer(1 * time.Microsecond)
OUTER:
	for {
		select {
		case <-timer1.C:
			fmt.Println("Ticker")
			close(done)
		case v, ok := <-chRes:
			if !ok {
				break OUTER
			}
			fmt.Println(v)
		}
	}

	for {
		fmt.Println("After Complete Go Routines", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}

}

func readDataCTX(ctx context.Context, file string) <-chan string {
	f, err := os.Open(file) //opens the file for reading
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string) //channel declared

	//returns a scanner to read from f
	fileScanner := bufio.NewScanner(f)

	//loop through the fileScanner based on our token split
	go func() {
		defer func() {
			fmt.Println("Done readData")
		}()
		for fileScanner.Scan() {
			val := fileScanner.Text() //returns the recent token

			// select
			select {
			case <-ctx.Done():
				fmt.Println("DONE--> CHANNEL")
				close(out)
				return
			case out <- val: //passed the token value to our channel
			default:
			}
		}
		fmt.Println("DONE--> readData")
		close(out) //closed the channel when all content of file is read
		//closed the file
		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close an opened file: %v\n", err.Error())
			return
		}
	}()

	return out
}

func fanInMergeDataCTX(ctx context.Context, chs ...<-chan string) chan string {

	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	for _, ch := range chs {
		go func(ch <-chan string) {
			defer wg.Done()
			for val := range ch {
				// select
				select {
				case <-ctx.Done():
					fmt.Println("Done fanInMergeData")
					return
				case chRes <- val:
				}
			}

		}(ch)
	}

	go func() {

		wg.Wait() //waits till the goroutines are completed and wg marked Done
		fmt.Println("Done Wait Go Routine")
		close(chRes) //close the result channel
	}()

	return chRes
}

func ctxDoneExample() {

	fmt.Println("Starting Go Routine", runtime.NumGoroutine())

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Microsecond)
	//receive data from multiple channels and place it on result channel - FanIn
	chRes := fanInMergeDataCTX(ctx, readDataCTX(ctx, "text1.txt"))

	fmt.Println("After Merge Go Routines", runtime.NumGoroutine())

	//some logic with the result channel

	timer1 := time.NewTimer(1 * time.Microsecond)
OUTER:
	for {
		select {
		case <-timer1.C:
			fmt.Println("Ticker")
			cancelFunc()
		case v, ok := <-chRes:
			if !ok {
				break OUTER
			}
			fmt.Println(v)
		}
	}

	for {
		fmt.Println("After Complete Go Routines", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ctxDoneExample()
}
