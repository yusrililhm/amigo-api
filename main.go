package main

import (
	"fashion-api/app"
	"runtime"
)

func main() {
	
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.StartApplication()
}
