package scanner

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/notify"
	"github.com/Cloud-Architect-Emma/go-ci-cd-security-guardrail/internal/scanner"
)

func main() {
	configPath := flag.String("config", "configs/policies.json", "Path to policy file")
	scanPath := flag.String("path", ".", "Path to scan")
	flag.Parse()

	fmt.Println(" Running Go Security Guardrail...")
	fmt.Println("Config:", *configPath)
	fmt.Println("Scan path:", *scanPath)

	policy, err := scanner.LoadPolicy(*configPath)
	if err != nil {
		fmt.Println(" Failed to load policy:", err)
		os.Exit(1)
	}

	var findings []string

	err = filepath.Walk(*scanPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			if strings.Contains(path, ".git") ||
				strings.Contains(path, ".github") {
				return filepath.SkipDir
			}
			return nil
		}

		// Scan only code files
		if !(strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, ".txt") ||
			strings.HasSuffix(path, ".env")) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scannerReader := bufio.NewScanner(file)
		lineNumber := 0

		for scannerReader.Scan() {
			lineNumber++
			line := scannerReader.Text()

			result := scanner.Scan(line, policy)

			if !result.Safe {
				for _, issue := range result.Issues {

					finding := fmt.Sprintf(
						"%s | Line %d | %s",
						path,
						lineNumber,
						issue,
					)

					findings = append(findings, finding)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Scan error:", err)
		os.Exit(1)
	}

	if len(findings) > 0 {

		fmt.Println(" Security violations detected:")

		message := " CI/CD Security Gate Failed\n\n"

		for _, finding := range findings {
			fmt.Println("-", finding)

			message += "- " + finding + "\n"
		}

		err := notify.SendSlackAlert(message)

		if err != nil {
			fmt.Println("Slack notification failed:", err)
		}

		os.Exit(1)
	}

	fmt.Println(" Go Security Gate passed")
}
