package scan

import (
	"bytes"
	"fmt"
	"io"
	"malice-new/plugins"
	"malice-new/server"
	"mime/multipart"
	"os"
	"sync"
)

type ResultsData struct {
	Infected bool   `json:"infected" `
	Result   string `json:"result"`
	Engine   string `json:"engine" `
	Database string `json:"database"`
	Updated  string `json:"updated"`
	Error    string `json:"error"`
	Plugin   string `json:"plugin"`
}

func CreateFile(file *multipart.FileHeader) (string, error) {
	os.MkdirAll(server.CFG.UploadPath, 0777)

	path := server.CFG.UploadPath + "/" + file.Filename
	fmt.Println(path)
	tmp, err := os.Create(path)
	if err != nil {
		fmt.Println("create")
		return path, err
	}

	mlfile, err := file.Open()
	if err != nil {
		fmt.Println("open")
		return path, err
	}

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, mlfile); err != nil {
		fmt.Println("xopy")
		return path, err
	}

	if _, err := tmp.Write(buf.Bytes()); err != nil {
		fmt.Println("Write")
		return path, err

	}

	return path, nil
}

func Scan(path string) []ResultsData {
	wg := sync.WaitGroup{}
	plugins := plugins.Enables()
	wg.Add(len(plugins))
	fmt.Println(len(plugins))
	var finalResp []ResultsData

	for _, plugin := range plugins {
		switch plugin {

		case "AVG":
			go func() {

				resp := avgScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()
				fmt.Println("avg")

			}()
		case "Zoner":
			go func() {

				resp := zonerScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()
				fmt.Println("zoner")

			}()
		case "Fprot":
			go func() {

				resp := fprotScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()
				fmt.Println("fprot")

			}()
		case "WindowsDefender":

			go func() {

				resp := defenderScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()
				fmt.Println("defender")

			}()
		case "Escan":
			go func() {

				resp := escanScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()
				fmt.Println("escan")

			}()
		case "Mcafee":
			go func() {
				fmt.Println("mcafee")

				resp := mcafeeScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()
		case "Clamav":
			go func() {
				fmt.Println("clamav")

				resp := clamScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "Avira":
			go func() {
				fmt.Println("avira")

				resp := aviraScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "Kaspersky":
			go func() {
				fmt.Println("Kaspersky")

				resp := kasperskyScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "DrWeb":
			go func() {
				fmt.Println("DrWeb")

				resp := drwebScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "Comodo":
			go func() {
				fmt.Println("Comodo")

				resp := comodoScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "Bitdefender":
			go func() {
				fmt.Println("Bitdefender")

				resp := bitdefenderScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		case "Avast":
			go func() {
				fmt.Println("Avast")

				resp := avastScan(path)
				finalResp = append(finalResp, *resp)
				wg.Done()

			}()

		default:
			fmt.Println("wrong av name : ", plugin)
			wg.Done()

		}

	}

	fmt.Println("before")

	wg.Wait()
	fmt.Println("after")

	return finalResp

}
