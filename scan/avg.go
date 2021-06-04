package scan

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Starts the avg deamon. required for scans
func StartAVGDaemon() error {
	// use restart as the daemon could be hanging or smth
	clamd := exec.Command("/etc/init.d/avgd", "restart")
	_, err := clamd.Output()

	if err != nil {
		return err
	}

	return nil
}

func avgScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "avg"

	statusOutput, _ := RunCommand(nil, "/etc/init.d/avgd", "status")
	if !strings.Contains(statusOutput, "is running") {
		err := StartAVGDaemon()
		if err != nil {
			resp.Error = "error start avg : " + err.Error()
			return resp
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()
	output, err := RunCommand(ctx, "/usr/bin/avgscan", path)

	if err != nil {
		resp.Error = err.Error()
		return resp
	}

	resp.Infected = false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		resp.Engine = getAvgVersion()
		wg.Done()
	}()

	lines := bytes.Split([]byte(output), []byte("\n"))

	colonSeparated := [][]byte{}

	for _, line := range lines {
		if len(line) != 0 {
			if bytes.Contains(line, []byte(":")) {
				colonSeparated = append(colonSeparated, line)
			}
			if bytes.Contains(line, []byte(path)) {
				pathVirusString := bytes.Split(line, []byte("  "))
				if len(pathVirusString) >= 2 {
					resp.Result = string(bytes.TrimSpace(pathVirusString[1]))
				}
			}
		}
	}
	// fmt.Println(lines)

	// Extract AVG Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := bytes.Split(line, []byte(":"))
				if len(keyvalue) != 0 {
					switch {
					case bytes.Contains(line, []byte("Virus database version")):
						resp.Database = string(bytes.TrimSpace(keyvalue[1]))
					case bytes.Contains(line, []byte("Virus database release date")):
						date := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Virus database release date:")))
						resp.Updated = parseUpdatedDate(string(date))
					case bytes.Contains(line, []byte("Infections found")):
						if bytes.Contains(keyvalue[1], []byte("1")) {
							resp.Infected = true
						}
					}
				}
			}
		}
	} else {
		fmt.Println("response was : ", output)

		resp.Error = "Unable to parse AVG output"

	}

	wg.Wait()
	return resp

}

func parseUpdatedDate(date string) string {
	layout := "Mon, 02 Jan 2006 15:04:05 +0000"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func getAvgVersion() string {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	output, err := RunCommand(ctx, "/usr/bin/avgscan", "-v")

	if err != nil {
		return ""
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) != 0 {
			keyvalue := strings.Split(line, ":")
			if len(keyvalue) != 0 {
				if strings.Contains(keyvalue[0], "Anti-Virus scanner version") {
					return strings.TrimSpace(keyvalue[1])
				}
			}
		}
	}
	return ""
}
