package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func mcafeeupdate() echo.Map {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	fmt.Println("Updating McAfee...")
	output, err := scan.RunCommand(nil, "/usr/local/uvscan/update")
	if err != nil {
		return echo.Map{"mcafee": err.Error()}
	}

	return echo.Map{"mcafee": output}
}
