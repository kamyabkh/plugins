package server

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	//EP server echo
	EP = echo.New()
	//CFG config
	CFG = loadConfig()
)

func Run() {
	EP.Use(middleware.CORS())
	EP.Use(middleware.Recover())
	EP.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 8,
	}))

	EP.Logger.Fatal(EP.Start(fmt.Sprintf(":%s", CFG.Port)))

}
