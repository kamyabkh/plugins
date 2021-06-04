package scan

import (
	"context"
	"malice-new/plugins"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func init() {
	if plugins.Plugin.Avast {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		// Avast needs to have the daemon started first
		avastd := exec.CommandContext(ctx, "/etc/init.d/avast", "start")
		_, err := avastd.Output()
		if err != nil {
			panic("error start avast : " + err.Error())
		}

	}
}

func avastScan(path string) *ResultsData {

	var output string
	var avErr error

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	resp := new(ResultsData)
	resp.Plugin = "avast"

	output, avErr = RunCommand(ctx, "scan", "-abfu", path)
	if avErr != nil {
		if avErr.Error() != "exit status 1" {
			resp.Error = avErr.Error()
			return resp
		}
	}

	return ParseAvastOutput(output, avErr)
}

// ParseAvastOutput convert avast output into ResultsData struct
func ParseAvastOutput(avastout string, avErr error) *ResultsData {

	avast := ResultsData{
		Infected: false,
		Updated:  "",
		Plugin:   "avast",
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		avast.Engine = getAvastVersion()
		wg.Done()
	}()

	go func() {
		avast.Database = getAvastVPS()
		wg.Done()
	}()

	result := strings.Split(avastout, "\t")

	if !strings.Contains(avastout, "[OK]") {
		avast.Infected = true
		avast.Result = strings.TrimSpace(result[1])
	}

	wg.Wait()

	return &avast
}

// Get Anti-Virus scanner version
func getAvastVersion() string {
	versionOut, err := RunCommand(nil, "scan", "-v")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(versionOut)
}

func getAvastVPS() string {
	versionOut, err := RunCommand(nil, "scan", "-V")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(versionOut)
}
