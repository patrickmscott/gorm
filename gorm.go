// Copyright 2010 Google Inc. All Rights Reserved.
// Author: phanna@google.com (Patrick Scott)

package main

import (
	"fmt"
	"os"
)

var sem = make(chan int, 2000)

func deleteDirectory(result chan int, directory string) {
	sem<-1
	handle, err := os.Open(directory, os.O_RDONLY, 0)
	var count int = 0
	var childResult = make(chan int)
	if err != nil {
		fmt.Printf("Failed to open file %s: %s\n", directory, err)
		os.Exit(1)
	}
	for {
		names, err := handle.Readdirnames(100)
		if err != nil {
			fmt.Printf("Failed to list directory contents: %s\n", err)
			os.Exit(1)
		}
		if len(names) == 0 {
			break
		}
		for _, name := range names {
			fullName := directory + "/" + name
			fileInfo, err := os.Stat(fullName)
			if err != nil {
				fmt.Printf("Failed to stat file %s: %s\n", fullName, err)
				os.Exit(1)
			}
			if fileInfo.IsDirectory() {
				go deleteDirectory(childResult, fullName)
				count++
			} else {
				err := os.Remove(fullName)
				if err != nil {
					fmt.Printf("Failed to remove file %s: %s\n", fullName, err)
					os.Exit(1)
				}
				fmt.Printf("%s removed\n", fullName)
			}
		}
	}
	handle.Close()
	<-sem
	for count != 0 {
		count = count - <-childResult
	}
	os.Remove(directory)
	fmt.Printf("%s removed\n", directory)
	result <- 1
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Missing directory to remove!\n")
		os.Exit(1)
	}
	var dir string = os.Args[1]

	file, err := os.Open(dir, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("Error opening directory: %s\n", err)
		os.Exit(1)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stat'ing directory: %s\n", err)
		os.Exit(1)
	}

	file.Close()
	if !fileInfo.IsDirectory() {
		fmt.Printf("File is not a directory!\n", dir)
		os.Exit(1)
	}

	var c = make(chan int)
	go deleteDirectory(c, dir)
	<-c
}
