package scan

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func bitdefenderScan(path string) *ResultsData {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	resp := new(ResultsData)
	results, err := RunCommand(ctx, "bdscan", path)
	resp.Plugin = "bitdefender"

	if err != nil {
		// AVG exits with error status 5 if it finds a virus
		if err.Error() != "exit status 1" {
			resp.Error = err.Error()
			return resp
		}
	}

	resp = ParseBitdefenderOutput(results)
	return resp
}

// ParseBitdefenderOutput convert bitdefender output into ResultsData struct
func ParseBitdefenderOutput(bitdefenderout string) *ResultsData {

	bitdefender := ResultsData{Infected: false, Plugin: "bitdefender"}
	// EXAMPLE OUTPUT:
	// BitDefender Antivirus Scanner for Unices v7.90123 Linux-amd64
	// Copyright (C) 1996-2009 BitDefender. All rights reserved.
	// Trial key found. 30 days remaining.
	//
	// Infected file action: ignore
	// Suspected file action: ignore
	// Loading plugins, please wait
	// Plugins loaded.
	//
	// /malware/EICAR  infected: EICAR-Test-File (not a virus)
	//
	//
	// Results:
	// Folders: 0
	// Files: 1
	// Packed: 0
	// Archives: 0
	// Infected files: 1
	// Suspect files: 0
	// Warnings: 0
	// Identified viruses: 1
	// I/O errors: 0
	lines := strings.Split(bitdefenderout, "\n")

	// Extract Virus string
	for _, line := range lines {
		if len(line) != 0 {
			switch {
			case strings.Contains(line, "infected:"):
				result := extractbitVirusName(line)
				if len(result) != 0 {
					bitdefender.Result = result
					bitdefender.Infected = true
				} else {
					fmt.Println("[ERROR] Virus name extracted was empty: ", result)

				}
			case strings.Contains(line, "Unices v"):
				words := strings.Fields(line)
				for _, word := range words {
					if strings.HasPrefix(word, "v") {
						bitdefender.Engine = strings.TrimPrefix(word, "v")
					}
				}
			}
		}
	}

	bitdefender.Updated = ""

	return &bitdefender
}

// extractbitVirusName extracts Virus name from scan results string
func extractbitVirusName(line string) string {
	keyvalue := strings.Split(line, "infected:")
	return strings.TrimSpace(keyvalue[1])
}
