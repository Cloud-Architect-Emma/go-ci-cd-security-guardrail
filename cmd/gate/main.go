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
	configPath := flag.String("config", "configs/policies.json", "Path to policy file")
	scanPath := flag.String("path", ".", "Path to scan")
	flag.Parse()

	fmt.Println("🔍 Running Go Security Guardrail...")
	fmt.Println("Config:", *configPath)
	fmt.Println("Scan path:", *scanPath)

	policy, err := scanner.LoadPolicy(*configPath)
	if err != nil {
		fmt.Println("❌ Failed to load policy:", err)
		os.Exit(1)
	}

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

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠️ Warning: grep execution issue:", err)
	}

	input := string(output)

	if input == "" {
		fmt.Println("No suspicious patterns found")
		fmt.Println("Go Security Gate passed")
		return
	}

	result := scanner.Scan(input, policy)

	if !result.Safe {
		fmt.Println("❌ Security violations detected:")
		message := "CI/CD Security Gate Failed\n\nIssues:\n"

		for _, issue := range result.Issues {
			fmt.Println("-", issue)
			message += "- " + issue + "\n"
		}

		err := notify.SendSlackAlert(message)
		if err != nil {
			fmt.Println("Slack notification failed:", err)
		}

		os.Exit(1)
	}

	fmt.Println("Go Security Gate passed")
}