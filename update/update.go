package update

import (
	"fmt"
	"malice-new/plugins"
	"sync"
)

type LicenseKeys struct {
	Kaspersky string `json:"kaspersky" form:"kaspersky"`
	Drweb     string `json:"drweb" form:"drweb"`
}

func Update(licenses *LicenseKeys) []interface{} {

	wg := sync.WaitGroup{}
	plugins := plugins.Enables()
	wg.Add(len(plugins))

	var finalResp []interface{}

	for _, plugin := range plugins {
		switch plugin {

		case "AVG":
			go func() {

				resp := avgupdate()
				finalResp = append(finalResp, resp)
				wg.Done()
				fmt.Println("avg")

			}()
		case "Zoner":
			go func() {

				resp := zonerupdate()
				finalResp = append(finalResp, resp)
				wg.Done()
				fmt.Println("zoner")

			}()
		case "Fprot":
			go func() {

				resp := fprotupdate()
				finalResp = append(finalResp, resp)
				wg.Done()
				fmt.Println("fprot")

			}()
		case "WindowsDefender":

			go func() {

				resp := defenderupdate()
				finalResp = append(finalResp, resp)
				wg.Done()
				fmt.Println("windwosdefender")

			}()
		case "Escan":
			go func() {

				resp := escanupdate()
				finalResp = append(finalResp, resp)
				wg.Done()
				fmt.Println("escan")

			}()
		case "Mcafee":
			go func() {
				fmt.Println("mcafee")

				resp := mcafeeupdate()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()
		case "Clamav":
			go func() {
				fmt.Println("clamav")

				resp := clamupdate()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "Avira":
			go func() {
				fmt.Println("avira")

				resp := aviraupdate()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "Kaspersky":
			go func() {
				fmt.Println("Kaspersky")

				resp := kasperskyupdate(licenses.Kaspersky)
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "DrWeb":
			go func() {
				fmt.Println("DrWeb")

				resp := drwebupdate(licenses.Drweb)
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "Comodo":
			go func() {
				fmt.Println("Comodo")

				resp := comodoUpdate()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "Bitdefender":
			go func() {
				fmt.Println("Bitdefender")

				resp := updatebitdefender()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		case "Avast":
			go func() {
				fmt.Println("Avast")

				resp := avastUpdate()
				finalResp = append(finalResp, resp)
				wg.Done()

			}()

		default:
			fmt.Println("wrong av name : ", plugin)
			wg.Done()

		}

	}

	return finalResp

}
