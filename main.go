package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

// readNumber prompts the user for input and securely attempts to parse it as a float64.
// It loops until a valid numeric input is provided.
func readNumber(prompt string) float64 {
   // Create a new reader for standard input.
   reader := bufio.NewReader(os.Stdin)

   for {
      fmt.Print(prompt)

      // ReadString reads input until a newline character is encountered.
      input, err := reader.ReadString('\n')
      if err != nil {
         // Handle potential I/O errors during reading.
         fmt.Println("An unexpected error occurred during input reading:", err)
         os.Exit(1)
      }

      // Clean up the input: remove the trailing newline and any surrounding whitespace.
      cleanedInput := strings.TrimSpace(input)

      if cleanedInput == "" {
         fmt.Println("Error: Input cannot be empty. Please enter a number.")
         continue
      }

      // Attempt to convert the cleaned string input to a float64.
      // Using ParseFloat handles integers and decimals and is crucial for secure validation.
      number, err := strconv.ParseFloat(cleanedInput, 64)
      if err != nil {
         // If conversion fails (e.g., input is "hello"), print an error and loop again.
         fmt.Printf("Error: '%s' is not a valid number. Please try again.\n", cleanedInput)
         continue
      }

      // If conversion is successful and there are no errors, return the valid number.
      return number
   }
}

// main is the entry point of the application.
func main() {
   fmt.Println("Adding two numbers.")

   // Get the first number using the interactive prompt function.
   num1 := readNumber("Enter the first number: ")

   // Get the second number.
   num2 := readNumber("Enter the second number: ")

   // Perform the calculation.
   result := num1 + num2

   // Output the result.
   fmt.Printf("\nCalculation Result:\n")
   // Using %v to print the most natural representation (e.g., 10 instead of 10.0).
   fmt.Printf("%v + %v = %v\n", num1, num2, result)
}