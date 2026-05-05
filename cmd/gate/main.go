package main

import (
	"fmt"
	"os"

	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/notify"
	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/scanner"
)

func main() {
	input := os.Getenv("CODE_INPUT")

	policy, err := scanner.LoadPolicy("configs/policies.json")
	if err != nil {
		fmt.Println("Failed to load policy:", err)
		os.Exit(1)
	}

	result := scanner.Scan(input, policy)

	if !result.Safe {
		fmt.Println("Security violations detected:")
		for _, issue := range result.Issues {
			fmt.Println("-", issue)
		}

		notify.SendSlackAlert("CI/CD Security Gate Failed:\n" + fmt.Sprint(result.Issues))
		os.Exit(1)
	}

	fmt.Println("Go Security Gate passed")
}