package api

import (
	"fmt"
	"malice-new/scan"
	"net/http"

	"github.com/labstack/echo"
)

func scanFile(c echo.Context) error {

	file, err := c.FormFile("file")

	if err != nil {
		fmt.Println("go")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	path, err := scan.CreateFile(file)
	if err != nil {
		fmt.Println("error")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	fmt.Println(path)

	resp := scan.Scan(path)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	// }

	return c.JSON(http.StatusOK, echo.Map{"data": resp})

}
