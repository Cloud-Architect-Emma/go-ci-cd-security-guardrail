package main

import (
	"fmt"
	"os"

	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/scanner"
)

func main() {
	input := os.Getenv("CODE_INPUT")

	result := scanner.Scan(input)

	if !result.Safe {
		fmt.Println(" Security violation detected:")
		for _, issue := range result.Issues {
			fmt.Println("-", issue)
		}
		os.Exit(1)
	}

	fmt.Println(" Pipeline passed Go security gate")
}