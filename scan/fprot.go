package scan

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func fprotScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "fprot"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	output, err := RunCommand(ctx, "/usr/local/bin/fpscan", "-r", path)

	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}

	resp = ParsefprotOutput(output)

	return resp
}

func ParsefprotOutput(fprotout string) *ResultsData {
	fprot := ResultsData{Infected: false}
	fprot.Plugin = "fprot"
	colonSeparated := []string{}

	lines := strings.Split(fprotout, "\n")
	// Extract Virus string and extract colon separated lines into an slice
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
			if strings.Contains(line, "[Found virus]") {
				result := extractfprotVirusName(line)
				if len(result) != 0 {
					fprot.Result = result
					fprot.Infected = true
				}
			}
		}
	}
	// fmt.Println(lines)

	// Extract FPROT Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := strings.Split(line, ":")
				if len(keyvalue) != 0 {
					switch {
					case strings.Contains(keyvalue[0], "Virus signatures"):
						fprot.Updated = parsefprotUpdatedDate(strings.TrimSpace(keyvalue[1]))
					case strings.Contains(line, "Engine version"):
						fprot.Engine = strings.TrimSpace(keyvalue[1])
					}
				}
			}
		}
	} else {
		fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
	}

	// fprot.Updated = getUpdatedDate()

	return &fprot
}

// extractVirusName extracts Virus name from scan results string
func extractfprotVirusName(line string) string {
	r := regexp.MustCompile(`<(.+)>`)
	res := r.FindStringSubmatch(line)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func parsefprotUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}
