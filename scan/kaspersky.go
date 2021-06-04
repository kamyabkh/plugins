package scan

import (
	"context"
	"fmt"
	"malice-new/plugins"
	"os/exec"
	"strings"
	"time"
)

var LicenseKey string

func InitKaspersky() {
	if plugins.Plugin.Kaspersky {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()

		configd := exec.CommandContext(ctx, "/etc/init.d/kav4fs-supervisor", "start")
		_, err := configd.Output()
		if err != nil {
			panic("error start kaspersky : " + err.Error())
		}
		time.Sleep(5)

		err = didLicenseExpire(ctx)
		if err != nil {
			panic("error update license kaspersky : " + err.Error())
		}
	}

}

func didLicenseExpire(ctx context.Context) error {
	// kaspersky needs to have the daemon started first

	license := exec.CommandContext(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--query-status")
	lOut, err := license.Output()
	if err != nil {
		return err
	}
	flag := true
	for _, line := range strings.Split(string(lOut), "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "License status:") {
				flag = strings.Contains(line, "Valid")
			}
		}
	}

	if flag {
		err = updateLicense(ctx)
		if err != nil {
			return err

		}
	}

	return nil
}

func updateLicense(ctx context.Context) error {
	// kaspersky needs to have the daemon started first

	// check for exec context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command updateLicense() timed out")
	}

	if len(LicenseKey) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		(fmt.Println(RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--revoke-active-key", LicenseKey)))
		fmt.Println(RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--install-active-key", LicenseKey))
	}

	return nil
}
func kasperskyScan(path string) *ResultsData {
	var output string
	var sErr error
	resp := new(ResultsData)
	resp.Plugin = "kaspersky"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	// kaspersky needs to have the daemon started first

	output, sErr = RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--scan-file", path)
	if sErr != nil {
		resp.Error = sErr.Error()
		return resp
	}

	virusInfo, err := RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--top-viruses", "1")
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(err)
		resp.Error = err.Error()
		return resp
	}

	results := ParseKasperskyOutput(output, virusInfo, sErr)

	return &results
}

// ParseKasperskyOutput convert kaspersky output into ResultsData struct
func ParseKasperskyOutput(kasperskyOut, virusInfo string, kasperskyErr error) ResultsData {
	resp := ResultsData{}
	resp.Plugin = "kaspersky"

	if kasperskyErr != nil {
		// if kasperskyErr.Error() == "exit status 119" {
		// 	return ResultsData{Error: "ScanEngine is not available"}, kasperskyErr
		// }
		resp.Error = kasperskyErr.Error()
		return resp
	}

	kaspersky := ResultsData{
		Infected: false,
		// Engine:   getKasperskyVersion(),
		// Database: getKasperskyDatabase(),
		// Updated:  getkaspaerskyUpdatedDate(),
		Plugin: "kaspersky",
	}

	for _, line := range strings.Split(kasperskyOut, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "Threats found:       1") {
				kaspersky.Infected = true
				for _, line := range strings.Split(virusInfo, "\n") {
					if len(line) != 0 {
						if strings.Contains(line, "Virus name:") {
							kaspersky.Result = strings.TrimSpace(strings.TrimPrefix(line, "Virus name:"))
						}
					}
				}
			}
		}
	}

	return kaspersky
}
