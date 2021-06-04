package scan

import (
	"context"
	"strings"
	"time"
)

func clamScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "clamav"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	output, err := RunCommand(ctx, "/usr/bin/clamscan", "--stdout", path)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}

	resp = parseclamoutput(output)

	return resp
}

func parseclamoutput(output string) *ResultsData {
	clamAV := ResultsData{}
	clamAV.Plugin = "clamav"
	lines := strings.Split(output, "\n")
	// Extract AV Scan Result
	result := lines[0]
	if len(result) != 0 {
		pathAndResult := strings.Split(result, ":")
		if strings.Contains(pathAndResult[1], "OK") {
			clamAV.Infected = false
		} else {
			clamAV.Infected = true
			clamAV.Result = strings.TrimSpace(strings.TrimRight(pathAndResult[1], "FOUND"))
		}
	} else {
		clamAV.Error = "[ERROR] empty scan result: "
		return &clamAV
	}
	// Extract Clam Details from SCAN SUMMARY
	for _, line := range lines[1:] {
		if len(line) != 0 {
			keyvalue := strings.Split(line, ":")
			if len(keyvalue) != 0 {
				switch {
				case strings.Contains(keyvalue[0], "Known viruses"):

				case strings.Contains(line, "Engine version"):
					clamAV.Engine = strings.TrimSpace(keyvalue[1])
				}
			}
		}
	}

	clamAV.Updated = ""

	return &clamAV
}
