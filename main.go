package main

import (
	"flag"
	"fmt"
	"malice-new/api"
	"malice-new/scan"
	"malice-new/server"
)

var (
	kaspersky, drweb string
)

func init() {
	initFlags()
}

func main() {
	flag.Parse()
	checkFlags()
	scan.InitDrweb()
	scan.InitKaspersky()

	api.Register()

	server.Run()

}

func initFlags() {

	// Strings
	flag.StringVar(&kaspersky, "kaspersky", "", "enter kaspersky license key")
	flag.StringVar(&drweb, "drweb", "", "enter drweb license key")

}

func checkFlags() {

	if kaspersky != "" {
		fmt.Println(kaspersky)
		scan.LicenseKey = kaspersky
	}

	if drweb != "" {
		scan.DrWebLicenseKey = drweb
	}
}
