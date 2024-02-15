package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func createAndWriteToFile(textFile string, output chan<- string) {
	//create a textFile
	file, err := os.Create(textFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	sayings := []string{

		"Rise brother, rise, The future is green.",
		"The height of the Iroko is still our dream",
		"Its very peak we will strive till we reach.",
	}

	for _, word := range sayings {
		_, err := writer.WriteString(word + " ")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		// flush buffered data to the file
		if err := writer.Flush(); err != nil {
			fmt.Println("Error writing to file:", err)
		}

	}

	// Send the file name to the channel
	output <- textFile
	close(output)
}

func countWordsAndPrint(textFileCh <-chan string, output chan<- int) {
	// Receive the file name from the channel
	textFile := <-textFileCh

	// Open the file for reading
	file, err := os.Open(textFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new Scanner with the file as the input source
	scanner := bufio.NewScanner(file)
	wordCount := 0

	for scanner.Scan() {
		// Split the line into words
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Send the word count to the channel
	output <- wordCount
	close(output)
}

func main() {

	textFile := "paxscribes.txt"

	// Channels
	fileNameCh := make(chan string)
	wordCountCh := make(chan int)

	// Start goroutines
	go createAndWriteToFile(textFile, fileNameCh)
	go countWordsAndPrint(fileNameCh, wordCountCh)

	// Receive the word count from the channel
	totalWords := <-wordCountCh

	// Print the result
	fmt.Printf("Total number of words in the file: %d\n", totalWords)
}
