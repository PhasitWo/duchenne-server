package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PhasitWo/duchenne-server/auth"
)

// main is the entry point of the Go application.
func main() {
	// Create a new buffered reader for standard input (keyboard)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("To exit, type 'exit' or press Ctrl+C.")
	fmt.Println("---------------------------------------")

	// Start an infinite loop to continuously read input
	for {
		// Display a simple prompt
		fmt.Print("Input: ")

		// Read the input until a newline character ('\n') is encountered
		input, err := reader.ReadString('\n')

		// Check for errors during reading (e.g., EOF signaled by Ctrl+D)
		if err != nil {
			fmt.Println("\n\nError reading input or EOF reached. Exiting program.")
			return
		}

		// Clean up the input string: remove leading/trailing whitespace, including the newline
		cleanedInput := strings.TrimSpace(input)

		// Check for the explicit 'exit' command
		if strings.ToLower(cleanedInput) == "exit" {
			fmt.Println("\nExit command received. Goodbye!")
			return
		}
		// Process and output the received input
		if cleanedInput == "" {
			fmt.Printf("can't hash empty string")
			return
		}

		hashed, _ := auth.HashPassword(cleanedInput)
		fmt.Println(hashed)
	}
}
