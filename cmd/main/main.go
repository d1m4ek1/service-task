package main

import (
	"service-task/internal/server"
	"service-task/pkg/seterror"
)

func main() {
	if err := server.InitServer(); err != nil {
		seterror.SetAppError("server.InitServer()", err)
		return
	}
}
