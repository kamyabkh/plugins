package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func fprotupdate() echo.Map {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	fmt.Println("Updating F-PROT...")
	output, err := scan.RunCommand(nil, "/opt/f-prot/fpupdate")

	if err != nil {
		return echo.Map{"fprot": err.Error()}

	}

	return echo.Map{"fprot": output}
}
