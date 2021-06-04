package scan

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func escanScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "escan"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()
	results, err := RunCommand(ctx, "escan", "-ly", path)
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}

	resp = parseEscan(results)

	return resp
}

func parseEscan(result string) *ResultsData {
	escan := ResultsData{Infected: false}
	escan.Plugin = "escan"

	lines := strings.Split(result, "\n")

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		escan.Engine = getescanEngine()
		wg.Done()
	}()

	go func() {
		escan.Updated = getescanUpdatedDate()
		wg.Done()
	}()

	// Extract Virus string
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, "[INFECTED]") {
				result := extractescanVirusName(line)
				if len(result) != 0 {
					escan.Result = result
					escan.Infected = true
				} else {
					escan.Error = fmt.Sprint("[ERROR] virus name extracted was empty: ", result)
				}
			}
		}
	}

	wg.Wait()

	return &escan

}

// extractVirusName extracts Virus name from scan results string
func extractescanVirusName(line string) string {
	keyvalue := strings.Split(line, "[INFECTED]")
	return strings.Trim(strings.TrimSpace(keyvalue[1]), "[]")
}

func getescanEngine() string {
	var engine = ""

	results, err := RunCommand(nil, "escan", "--version")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(results, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "MicroWorld eScan For Linux Version :") {
				engine = strings.TrimSpace(strings.TrimPrefix(line, "MicroWorld eScan For Linux Version :"))
			}
		}
	}
	return engine
}

func getescanUpdatedDate() string {
	return ""
}
