package scan

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// McAfeeResults is the xml data struct
type McAfeeResults struct {
	XMLName  xml.Name     `xml:"Uvscan"`
	Preamble preamble     `xml:"Preamble"`
	DateTime scanDateTime `xml:"Date_Time"`
	Options  scanOptions  `xml:"Options"`
	File     fileResults  `xml:"File"`
	Time     timeToScan   `xml:"Time"`
}

type productName struct {
	Value string `xml:"value,attr"`
}

type mVersion struct {
	Value string `xml:"value,attr"`
}

type licenseInfo struct {
	Value string `xml:"value,attr"`
}

type aVEngineVersion struct {
	Value string `xml:"value,attr"`
}

type datSetVersion struct {
	Value string `xml:"value,attr"`
}

type preamble struct {
	XMLName         xml.Name        `xml:"Preamble"`
	ProductName     productName     `xml:"Product_name"`
	Version         mVersion        `xml:"Version"`
	LicenseInfo     licenseInfo     `xml:"License_info"`
	AVEngineVersion aVEngineVersion `xml:"AV_Engine_version"`
	DatSetVersion   datSetVersion   `xml:"Dat_set_version"`
}

type scanDateTime struct {
	Value string `xml:"value,attr"`
}

type scanOptions struct {
	Value string `xml:"value,attr"`
}

type fileResults struct {
	Name          string `xml:"name,attr"`
	Status        string `xml:"status,attr"`
	VirusName     string `xml:"virus-name,attr"`
	DetectionType string `xml:"detection-type,attr"`
}

type timeToScan struct {
	Value string `xml:"value,attr"`
}

func mcafeeScan(path string) *ResultsData {
	resp := new(ResultsData)
	resp.Plugin = "mcafee"
	hash := GetSHA256(path)
	if hash == "" {
		resp.Error = "error create mcafee hash"
		return resp
	}
	fmt.Println(hash)
	// defer os.Remove("/tmp/" + hash + ".xml")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	output, err := RunCommand(ctx, "/usr/local/uvscan/uvscan_secure", path, "--xmlpath=/tmp/"+hash+".xml")

	if err != nil && err.Error() != "exit status 1" {
		fmt.Println("first error")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		output, err := RunCommand(ctx, "/usr/local/uvscan/uvscan_secure", path, "--xmlpath=/tmp/"+hash+".xml")
		if err != nil && err.Error() != "exit status 1" {
			fmt.Println(err)
			resp.Error = err.Error()
			return resp
		}

		resp, err = ParseMcAfeeOutput(output, hash)
		if err != nil {
			resp.Plugin = "mcafee"
			resp.Error = err.Error()

			return resp
		}

		return resp
	}

	resp, err = ParseMcAfeeOutput(output, hash)
	if err != nil {
		resp.Plugin = "mcafee"
		resp.Error = err.Error()

		return resp
	}

	return resp

}

// ParseMcAfeeOutput convert mcafee output into ResultsData struct
func ParseMcAfeeOutput(mcafeeout, hash string) (*ResultsData, error) {

	xmlFile, err := os.Open("/tmp/" + hash + ".xml")
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var results McAfeeResults
	err = xml.Unmarshal(byteValue, &results)

	mcafee := ResultsData{
		Infected: strings.EqualFold(results.File.Status, "infected"),
		Engine:   results.Preamble.AVEngineVersion.Value,
		Database: results.Preamble.DatSetVersion.Value,
		Updated:  "",
		Result:   results.File.VirusName,
	}

	return &mcafee, nil
}
