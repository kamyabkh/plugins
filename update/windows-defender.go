package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func defenderupdate() echo.Map {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()
	fmt.Println("Updating Windows Defender...")
	output, err := scan.RunCommand(nil, "/opt/malice/update")

	if err != nil {
		return echo.Map{"windows-defender": err.Error()}
	}

	return echo.Map{"windows-defender": output}
}
