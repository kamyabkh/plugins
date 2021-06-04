package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func avgupdate() echo.Map {
	fmt.Println("Updating AVG...")

	output, err := scan.RunCommand(nil, "avgupdate")

	if err != nil {
		return echo.Map{"avg": err.Error()}
	}

	return echo.Map{"avg": output}
}
