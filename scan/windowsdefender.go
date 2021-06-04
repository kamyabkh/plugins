package scan

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"
)

func defenderScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "windows-defender"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()
	if err := os.Chdir("/loadlibrary"); err != nil {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}
	// defer os.Chdir("/malware")
	output, err := RunCommand(ctx, "./mpclient", path)

	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}

	resp = ParsedefenderOutput([]byte(output))

	return resp
}

func ParsedefenderOutput(ouput []byte) *ResultsData {
	windef := new(ResultsData)
	windef.Plugin = "windows-defender"

	lines := bytes.Split(ouput, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte("Scanning input")) {
			continue
		}
		if bytes.Contains(line, []byte("EngineScanCallback")) {
			threat := bytes.TrimPrefix(bytes.TrimSpace(line), []byte("EngineScanCallback():"))
			if len(threat) > 0 {
				windef.Infected = true
				threat = bytes.TrimSpace(threat)
				threat = bytes.TrimPrefix(threat, []byte("Threat"))
				threat = bytes.TrimSuffix(threat, []byte("identified."))
				windef.Result = string(bytes.TrimSpace(threat))
			}
		}
	}

	return windef
}
