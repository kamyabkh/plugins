package api

import "malice-new/server"

func Register() {
	ep := server.EP

	ep.POST("/scan", scanFile)
	ep.GET("/update", updatee)
}
