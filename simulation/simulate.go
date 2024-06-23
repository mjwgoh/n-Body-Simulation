package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: go run simulate.go {size} {optional: threads} {optional: p or q}")
	}

	if len(os.Args) < 3 {
		inputLink := os.Args[1]
		Sequential(inputLink)
	} else {
		inputLink := os.Args[1]
		numThreadsStr := os.Args[2]
		numThreads, err := strconv.Atoi(numThreadsStr)
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error converting number of threads:", err)
			fmt.Println("Usage: go run simulate.go {size} {optional: threads} {optional: p or q}")
			return
		}

		if len(os.Args) > 3 && os.Args[3] == "p" {
			Parallel(inputLink, numThreads)
		} else if len(os.Args) > 3 && os.Args[3] == "q" {
			WQParallel(inputLink, numThreads)
		}
	}
}
