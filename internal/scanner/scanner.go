package scanner

import (
	"strings"
	"sync"
)

type ScanResult struct {
	Safe   bool
	Issues []string
}

func Scan(input string, policy Policy) ScanResult {

	var wg sync.WaitGroup

	issuesChan := make(chan string, len(policy.Rules))

	for _, rule := range policy.Rules {

		wg.Add(1)

		go func(r Rule) {
			defer wg.Done()

			if strings.Contains(input, r.Pattern) {

				issuesChan <- r.Message + " [" + r.Severity + "]"
			}

		}(rule)
	}

	wg.Wait()

	close(issuesChan)

	var issues []string

	for issue := range issuesChan {

		issues = append(issues, issue)
	}

	return ScanResult{
		Safe:   len(issues) == 0,
		Issues: issues,
	}
}
