package api

import (
	"malice-new/update"
	"net/http"

	"github.com/labstack/echo"
)

func updatee(c echo.Context) error {
	form := new(update.LicenseKeys)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	resp := update.Update(form)

	return c.JSON(http.StatusOK, echo.Map{"data": resp})
}
