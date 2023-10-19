// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	handler "pdm-plugin.github.com/biz/handler"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.GET("/ping", handler.Ping)

	r.POST("/data", handler.HandleData)

	// your code ...
}
