package scan

import (
	"bufio"
	"context"
	"fmt"
	"malice-new/plugins"
	"os"
	"strings"
	"time"
)

var comodoversion = ""

func init() {
	if plugins.Plugin.Comodo {
		comodoversion = getComodoVersion()
	}
}

func comodoScan(path string) *ResultsData {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	resp := new(ResultsData)
	resp.Plugin = "comodo"

	output, err := RunCommand(ctx, "/opt/COMODO/cmdscan", "-vs", path)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}

	Results := ParseComodoOutput(output)

	return Results

}

// ParseComodoOutput convert comodo output into ResultsData struct
func ParseComodoOutput(comodoout string) *ResultsData {

	comodo := ResultsData{Infected: false, Engine: comodoversion, Updated: ""}
	comodo.Plugin = "comodo"

	// EXAMPLE OUTPUT:
	// -----== Scan Start ==-----
	// /malware/EICAR ---> Found Virus, Malware Name is Malware
	// -----== Scan End ==-----
	// Number of Scanned Files: 1
	// Number of Found Viruses: 1
	lines := strings.Split(comodoout, "\n")

	// Extract Virus string
	if len(lines[1]) != 0 {
		if strings.Contains(lines[1], "Found Virus") {
			result := extractcomodoVirusName(lines[1])
			if len(result) != 0 {
				comodo.Result = result
				comodo.Infected = true
				return &comodo
			}

		}
	}

	return &comodo
}

// extractcomodoVirusName extracts Virus name from scan results string
func extractcomodoVirusName(line string) string {
	keyvalue := strings.Split(line, "is")
	return strings.TrimSpace(keyvalue[1])
}

func getComodoVersion() string {
	file, err := os.Open("/opt/COMODO/etc/COMODO.xml")
	if err != nil {
		fmt.Println("error open comodo version file")
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "<ProductVersion>") {
			versionOut := strings.TrimSpace(strings.Replace(strings.Replace(line, "<ProductVersion>", "", 1), "</ProductVersion>", "", 1))

			return versionOut
		}
	}
	return ""
}
