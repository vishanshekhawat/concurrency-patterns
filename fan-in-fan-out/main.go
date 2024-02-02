package main

/*
--->What is Fan-in Pattern...
	The fan in patter merge multiple channels values into single channel values.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func readData(file string) <-chan string {
	f, err := os.Open(file) //opens the file for reading
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string) //channel declared

	//returns a scanner to read from f
	fileScanner := bufio.NewScanner(f)

	//loop through the fileScanner based on our token split
	go func() {
		for fileScanner.Scan() {
			val := fileScanner.Text() //returns the recent token
			out <- val                //passed the token value to our channel
		}

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

func fanInMergeData(chs ...<-chan string) chan string {

	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	for _, ch := range chs {
		go func(ch <-chan string) {
			for val := range ch {
				chRes <- val
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()    //waits till the goroutines are completed and wg marked Done
		close(chRes) //close the result channel
	}()

	return chRes
}

func main() {

	//receive data from multiple channels and place it on result channel - FanIn
	chRes := fanInMergeData(readData("text1.txt"), readData("text2.txt"), readData("text3.txt"))

	//some logic with the result channel
	for val := range chRes {
		fmt.Println(val)
	}
}
