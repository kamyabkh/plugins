package update

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/levigross/grequests"
)

func comodoUpdate() echo.Map {
	fmt.Println("Updating Comodo...")

	response, err := grequests.Get("http://download.comodo.com/av/updates58/sigs/bases/bases.cav", nil)
	if err != nil {
		return echo.Map{"comodo": err.Error()}
	}

	if response.Ok != true {
		return echo.Map{"comodo": "response download comodo not ok"}
	}

	if err = response.DownloadToFile("/opt/COMODO/scanners/bases.cav"); err != nil {
		return echo.Map{"comodo": err.Error()}
	}

	return echo.Map{"comodo": "done"}

}
