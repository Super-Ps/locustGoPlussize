package main

import "falcon/server/service"

func main() {
	svr := service.New()
	svr.Start()
}