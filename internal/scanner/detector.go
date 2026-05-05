package scanner

import "strings"

type ScanResult struct {
	Safe   bool
	Issues []string
}

func Scan(input string) ScanResult {
	issues := []string{}

	if strings.Contains(input, "sk-") {
		issues = append(issues, "Hardcoded secret detected")
	}

	if strings.Contains(input, "API_KEY") {
		issues = append(issues, "Potential API key exposure")
	}

	return ScanResult{
		Safe:   len(issues) == 0,
		Issues: issues,
	}
}