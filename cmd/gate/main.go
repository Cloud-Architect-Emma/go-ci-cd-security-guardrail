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
	configPath := flag.String("config", "configs/policies.json", "Path to policy file")
	scanPath := flag.String("path", ".", "Path to scan")
	flag.Parse()

	fmt.Println("Running Go Security Guardrail...")
	fmt.Println("Config:", *configPath)
	fmt.Println("Scan path:", *scanPath)

	policy, err := scanner.LoadPolicy(*configPath)
	if err != nil {
		fmt.Println("Failed to load policy:", err)
		os.Exit(1)
	}

	cmd := exec.Command(
		"grep",
		"-rEn",
		"--exclude-dir=.git",
		"--exclude-dir=.github",
		"--exclude-dir=configs",
		"--exclude-dir=cmd",       // FIX: exclude the scanner's own source directory
		"--exclude=*.md",
		"--exclude=*.json",
		"--exclude=gate",
		"API_KEY|sk-|token|secret",
		*scanPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Warning: grep execution issue:", err)
	}

	input := string(output)

	if input == "" {
		fmt.Println("No suspicious patterns found")
		fmt.Println("Go Security Gate passed")
		return
	}

	lines := strings.Split(input, "\n")
	result := scanner.Scan(input, policy)

	if !result.Safe {
		fmt.Println("\n Security violations detected:\n")
		message := "CI/CD Security Gate Failed\n\n"

		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			// Expected format: file:line:content
			parts := strings.SplitN(line, ":", 3)
			if len(parts) < 3 {
				continue
			}

			file := parts[0]
			lineNo := parts[1]
			code := parts[2]

			severity := "MEDIUM"
			fix := "Review and remove sensitive data"

			if strings.Contains(code, "sk-") {
				severity = "HIGH"
				fix = "Remove hardcoded secrets and use environment variables or a secret manager (e.g., AWS Secrets Manager)"
			} else if strings.Contains(code, "API_KEY") {
				severity = "MEDIUM"
				fix = "Do not expose API keys in source code. Store them securely in environment variables"
			} else if strings.Contains(strings.ToLower(code), "token") {
				severity = "HIGH"
				fix = "Avoid embedding tokens in code. Use secure authentication flows"
			}

			fmt.Printf("[%s] Issue detected\n", severity)
			fmt.Println("File:", file)
			fmt.Println("Line:", lineNo)
			fmt.Println("Code:", strings.TrimSpace(code))
			fmt.Println("Fix:", fix)
			fmt.Println()

			message += fmt.Sprintf(
				"[%s]\nFile: %s\nLine: %s\nCode: %s\nFix: %s\n\n",
				severity,
				file,
				lineNo,
				strings.TrimSpace(code),
				fix,
			)
		}

		err := notify.SendSlackAlert(message)
		if err != nil {
			fmt.Println(" Slack notification failed:", err)
		}

		os.Exit(1)
	}

	fmt.Println("Go Security Gate passed")
}