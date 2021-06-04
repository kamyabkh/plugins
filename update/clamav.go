package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func clamupdate() echo.Map {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()
	fmt.Println("Updating ClamAV...")
	output, err := scan.RunCommand(nil, "freshclam")
	// Update UPDATED file
	if err != nil {
		return echo.Map{"clamav": err.Error()}

	}

	return echo.Map{"clamav": output}
}
