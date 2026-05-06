package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/notify"
	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/scanner"
)

func main() {
	// CLI flags
	configPath := flag.String("config", "configs/policies.json", "Path to policy file")
	scanPath := flag.String("path", ".", "Path to scan")
	flag.Parse()

	fmt.Println("🔍 Running Go Security Guardrail...")
	fmt.Println("Config:", *configPath)
	fmt.Println("Scan path:", *scanPath)

	// Load policy
	policy, err := scanner.LoadPolicy(*configPath)
	if err != nil {
		fmt.Println("Failed to load policy:", err)
		os.Exit(1)
	}

	// Run grep to collect input dynamically
	cmd := exec.Command(
	"grep",
	"-rE",
	"--exclude-dir=.git",
	"--exclude-dir=.github",
	"--exclude-dir=configs",
	"--exclude=*.md",
	"--exclude=*.json",
	"API_KEY|sk-|token|secret",
	*scanPath,
)
	output, _ := cmd.CombinedOutput() // ignore grep exit code

	input := string(output)

	// Scan
	result := scanner.Scan(input, policy)

	if !result.Safe {
		fmt.Println("❌ Security violations detected:")
		for _, issue := range result.Issues {
			fmt.Println("-", issue)
		}

		// Send Slack alert
		notify.SendSlackAlert(
			"CI/CD Security Gate Failed\n\n" +
				fmt.Sprintf("Issues:\n%v", result.Issues),
		)

		os.Exit(1)
	}

	fmt.Println("Go Security Gate passed")
}