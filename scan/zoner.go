package scan

import (
	"context"
	"fmt"
	"malice-new/plugins"
	"strings"
	"sync"
	"time"
)

func init() {

	startZavService()
}

func startZavService() {
	if plugins.Plugin.Zoner {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		_, err := RunCommand(ctx, "sudo", "/etc/init.d/zavd", "restart", "--no-daemon")

		if err != nil {
			panic("error start zoner service : " + err.Error())
		}
	}

}

func zonerScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "zoner"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()
	output, err := RunCommand(ctx, "zavcli", path)

	if err != nil && err.Error() != "exit status 11" {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}

	resp = ParseZonerOutput(output)

	return resp

}

// ParseZonerOutput convert zoner output into ResultsData struct
func ParseZonerOutput(zonerout string) *ResultsData {

	resp := ResultsData{Infected: false}
	resp.Plugin = "zoner"
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		resp.Engine = getEngine()
		wg.Done()
	}()

	go func() {
		resp.Updated = getUpdatedDate()
		wg.Done()

	}()

	lines := strings.Split(zonerout, "\n")

	// Extract Virus string
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, "INFECTED") {
				result := extractVirusName(line)
				if len(result) != 0 {
					resp.Result = result
					resp.Infected = true
				} else {
					resp.Error = fmt.Sprint("[ERROR] virus name extracted was empty: ", result)
				}
			}
		}
	}
	wg.Wait()
	return &resp
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	keyvalue := strings.Split(line, "INFECTED")
	return strings.Trim(strings.TrimSpace(keyvalue[1]), "[]")
}

func getEngine() string {
	var engine = ""

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	output, err := RunCommand(ctx, "zavcli", "--version-zavd")

	if err != nil {
		return ""
	}

	for _, line := range strings.Split(output, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "ZAVDB version:") {
				engine = strings.TrimSpace(strings.TrimPrefix(line, "ZAVDB version:"))
			}
		}
	}
	return engine
}

func getUpdatedDate() string {
	return ""
}
