package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func aviraupdate() echo.Map {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()
	fmt.Println("Updating Avira...")
	output, err := scan.RunCommand(nil, "/opt/malice/update")

	if err != nil {
		return echo.Map{"avira": err.Error()}
	}

	return echo.Map{"avira": output}

}
