package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func updatebitdefender() echo.Map {
	fmt.Println("Updating Bitdefender...")
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	output, err := scan.RunCommand(nil, "bdscan", "--update")
	if err != nil {
		return echo.Map{"bitdefender": err.Error()}
	}

	return echo.Map{"bitdefender": output}
}
