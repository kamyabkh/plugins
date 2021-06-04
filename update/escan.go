package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func escanupdate() echo.Map {

	fmt.Println("Updating EScan...")
	// EScan needs to have the daemon started first
	output, err := scan.RunCommand(nil, "escan", "--update")
	if err != nil {
		return echo.Map{"escan": err.Error()}
	}

	return echo.Map{"escan": output}
}
