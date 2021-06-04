package scan

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func aviraScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "avira"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	results, err := RunCommand(ctx, "/opt/avira/scancl", path)
	if err != nil {
		if err.Error() == "exit status 219" {
			resp.Error = "ERROR: [No license found] Initialization"
			return resp
		}
		if err.Error() != "exit status 1" {
			resp.Error = err.Error()
			return resp

		}

	}

	resp = parseAviraOutput(results)

	return resp

}

func parseAviraOutput(output string) *ResultsData {
	avira := ResultsData{Infected: false, Plugin: "avira"}
	lines := strings.Split(output, "\n")

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		avira.Engine = ParseAviraEngine(getaviraEngine())

		wg.Done()
	}()

	go func() {
		avira.Updated = getaviraUpdatedDate()

		wg.Done()
	}()

	// Extract Virus string
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, "ALERT:") {
				result, err := extractaviraVirusName(line)
				if err != nil {
					avira.Error = err.Error()
					return &avira

				}
				avira.Result = result
				avira.Infected = true
			}
		}
	}

	wg.Wait()

	return &avira
}

// extractaviraVirusName extracts Virus name from scan results string
func extractaviraVirusName(line string) (string, error) {
	var rgx = regexp.MustCompile(`\[.*?\]`)
	rs := rgx.FindStringSubmatch(line)
	if len(rs) > 0 {
		return strings.Trim(strings.TrimSpace(rs[0]), "[]"), nil
	}
	return "", fmt.Errorf("was not able to extract virus name from: %s", line)
}

func getaviraEngine() string {
	results, _ := RunCommand(nil, "/opt/avira/scancl", "--version")

	return results
}

// ParseAviraEngine convert avira version into engine string
func ParseAviraEngine(aviraVersion string) string {
	var engine = ""
	for _, line := range strings.Split(aviraVersion, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "engine set:") {
				engine = strings.TrimSpace(strings.TrimPrefix(line, "engine set:"))
			}
		}
	}

	return engine
}

func getaviraUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return ""
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	if err != nil {
		return ""
	}
	return string(updated)
}
