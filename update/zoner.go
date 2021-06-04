package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func zonerupdate() echo.Map {

	fmt.Println("Updating Zoner...")
	// Zoner needs to have the daemon started first
	output, err := scan.RunCommand(nil, "/etc/init.d/zavd", "update")
	if err != nil {
		return echo.Map{"zoner": err.Error()}
	}

	return echo.Map{"zoner": output}
}
