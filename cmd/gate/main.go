package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/notify"
	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/scanner"
)

func main() {

	// CLI flags
	configPath := flag.String(
		"config",
		"configs/policies.json",
		"Path to policy file",
	)

	scanPath := flag.String(
		"path",
		".",
		"Path to scan",
	)

	flag.Parse()

	fmt.Println("Running Go Security Guardrail...")
	fmt.Println("Config:", *configPath)
	fmt.Println("Scan path:", *scanPath)

	// Load policy
	policy, err := scanner.LoadPolicy(*configPath)
	if err != nil {
		fmt.Println("Failed to load policy:", err)
		os.Exit(1)
	}

	// Run grep scan
	cmd := exec.Command(
		"grep",
		"-rEn",

		// Exclude folders
		"--exclude-dir=.git",
		"--exclude-dir=.github",
		"--exclude-dir=configs",
		"--exclude-dir=internal",

		// Exclude files
		"--exclude=*.md",
		"--exclude=*.json",
		"--exclude=go.sum",
		"--exclude=go.mod",
		"--exclude=main.go",
		"--exclude=gate",

		// Search patterns
		"API_KEY|sk-|token|secret",

		*scanPath,
	)

	output, err := cmd.CombinedOutput()

	// Ignore grep exit code 1
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				fmt.Println("Warning: grep execution issue:", err)
			}
		}
	}

	input := string(output)

	// No findings
	if strings.TrimSpace(input) == "" {
		fmt.Println("No suspicious patterns found")
		fmt.Println("Go Security Gate passed")
		return
	}

	lines := strings.Split(input, "\n")

	// Run scanner
	result := scanner.Scan(input, policy)

	// Violations found
	if !result.Safe {

		fmt.Println("\nSecurity violations detected:\n")

		message := "CI/CD Security Gate Failed\n\n"

		for _, line := range lines {

			if strings.TrimSpace(line) == "" {
				continue
			}

			// file:line:content
			parts := strings.SplitN(line, ":", 3)

			if len(parts) < 3 {
				continue
			}

			file := parts[0]
			lineNo := parts[1]
			code := parts[2]

			severity := "MEDIUM"
			fix := "Review and remove sensitive data"

			// Severity handling
			if strings.Contains(code, "sk-") {

				severity = "HIGH"

				fix = "Remove hardcoded secrets and use environment variables"

			} else if strings.Contains(code, "API_KEY") {

				severity = "MEDIUM"

				fix = "Store API keys securely using environment variables"

			} else if strings.Contains(
				strings.ToLower(code),
				"token",
			) {

				severity = "HIGH"

				fix = "Avoid embedding tokens directly in code"
			}

			// Console output
			fmt.Printf("[%s] Issue detected\n", severity)
			fmt.Println("File:", file)
			fmt.Println("Line:", lineNo)
			fmt.Println("Code:", strings.TrimSpace(code))
			fmt.Println("Fix:", fix)
			fmt.Println()

			// Slack message
			message += fmt.Sprintf(
				"[%s]\nFile: %s\nLine: %s\nCode: %s\nFix: %s\n\n",
				severity,
				file,
				lineNo,
				strings.TrimSpace(code),
				fix,
			)
		}

		// Send Slack alert
		err := notify.SendSlackAlert(message)

		if err != nil {
			fmt.Println("Slack notification failed:", err)
		} else {
			fmt.Println("Slack alert sent")
		}

		os.Exit(1)
	}

	fmt.Println("Go Security Gate passed")
}
