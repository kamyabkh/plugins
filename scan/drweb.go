package scan

import (
	"context"
	"fmt"
	"malice-new/plugins"
	"os/exec"
	"strings"
	"time"
)

var DrWebLicenseKey string

func InitDrweb() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	if plugins.Plugin.DrWeb {
		configd := exec.CommandContext(ctx, "/opt/drweb.com/bin/drweb-configd", "-d")
		_, err := configd.Output()
		if err != nil {
			panic("error start dr web : " + err.Error())
		}

		time.Sleep(1 * time.Second)
		expired, err := diddrwebLicenseExpire(ctx)
		if err != nil {
			panic("error license drweb : " + err.Error())
		}
		if expired {
			err = updatedrwebLicense(ctx)
			if err != nil {
				panic("error update license dr web : " + err.Error())
			}
		}
	}

}

func diddrwebLicenseExpire(ctx context.Context) (bool, error) {
	// drweb needs to have the daemon started first

	license := exec.CommandContext(ctx, "/opt/drweb.com/bin/drweb-ctl", "license")
	lOut, err := license.Output()
	if err != nil {
		return false, err
	}

	if strings.Contains(string(lOut), "No license") {
		return true, nil
	}

	if strings.Contains(string(lOut), "expires") {
		return false, nil
	}

	return true, nil
}

func updatedrwebLicense(ctx context.Context) error {

	// check for exec context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command updateLicense() timed out")
	}

	if len(DrWebLicenseKey) > 0 {
		fmt.Println(RunCommand(ctx, "/opt/drweb.com/bin/drweb-ctl", "license", "--GetRegistered", DrWebLicenseKey))
	} else {
		fmt.Println(RunCommand(ctx, "/opt/drweb.com/bin/drweb-ctl", "license", "--GetDemo"))
	}

	return nil
}

func drwebScan(path string) *ResultsData {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()
	resp := new(ResultsData)
	resp.Plugin = "drweb"
	output, sErr := RunCommand(ctx, "/opt/drweb.com/bin/drweb-ctl", "scan", path)
	if sErr != nil {
		if sErr.Error() == "exit status 119" {
			resp.Error = "ScanEngine is not available"
			return resp
		}

	}

	baseinfo, err := RunCommand(ctx, "/opt/drweb.com/bin/drweb-ctl", "baseinfo")
	if err != nil && err.Error() != "exit status 13" {
		resp.Error = err.Error()
		return resp
	}

	results := ParseDrWEBOutput(output, baseinfo, path, sErr)

	return results
}

// ParseDrWEBOutput convert drweb output into ResultsData struct
func ParseDrWEBOutput(drwebOut, baseInfo, path string, drwebErr error) *ResultsData {

	resp := new(ResultsData)
	resp.Plugin = "drweb"
	if drwebErr != nil {
		if drwebErr.Error() == "exit status 119" {
			resp.Error = "ScanEngine is not available"
			return resp
		}
		resp.Error = drwebErr.Error()
		return resp
	}

	drweb := ResultsData{
		Infected: false,
		Engine:   getDrWebVersion(),
		// Updated:  getUpdatedDate(),
		Plugin: "drweb",
	}

	for _, line := range strings.Split(drwebOut, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "- Ok") {
				break
			}
			if strings.Contains(line, "infected with") {
				drweb.Infected = true
				drweb.Result = strings.TrimSpace(strings.TrimPrefix(line, path+" - infected with"))
			}
		}
	}

	for _, line := range strings.Split(baseInfo, "\n") {
		if len(line) != 0 {
			if strings.Contains(line, "Core engine:") {
				drweb.Engine = strings.TrimSpace(strings.TrimPrefix(line, "Core engine:"))
			}
			if strings.Contains(line, "Virus base records:") {
				drweb.Database = strings.TrimSpace(strings.TrimPrefix(line, "Virus base records:"))
			}
		}
	}

	return &drweb
}

func getDrWebVersion() string {

	versionOut, err := RunCommand(nil, "/opt/drweb.com/bin/drweb-ctl", "--version")

	if err != nil {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(versionOut, "drweb-ctl "))
}
