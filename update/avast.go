package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func avastUpdate() echo.Map {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	fmt.Println("starting avast update")

	output, err := scan.RunCommand(nil, "/var/lib/avast/Setup/avast.vpsupdate")
	if err != nil {
		return echo.Map{"avast": err.Error()}
	}

	return echo.Map{"avast": output}

}
